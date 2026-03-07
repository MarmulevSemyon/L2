package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	// "golang.org/x/sys/unix"
)

var currentPGID int
var jobMu sync.Mutex

func setCurrentPGID(pgid int) {
	jobMu.Lock()
	defer jobMu.Unlock()
	currentPGID = pgid
}

func clearCurrentPGID() {
	jobMu.Lock()
	defer jobMu.Unlock()
	currentPGID = 0
}

func getCurrentPGID() int {
	jobMu.Lock()
	defer jobMu.Unlock()
	return currentPGID
}

func Run() {
	reader := bufio.NewReader(os.Stdin)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	// перехватываем ctrl+c и убиваем группу процессов (весь pipeline) по пиду группы
	go func() {
		for range sigCh {
			pgid := getCurrentPGID()
			if pgid != 0 {
				// _ = unix.Kill(-pgid, unix.SIGINT)
				_ = syscall.Kill(-pgid, syscall.SIGINT)
				fmt.Print("\n")
			} else {
				fmt.Print("\n$ ")
			}
		}
	}()

	for {
		fmt.Print("$ ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println()
				return
			}
			fmt.Fprintf(os.Stderr, "read: %v\n", err)
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if err := executeLine(line); err != nil {
			fmt.Fprintf(os.Stderr, "run: %v\n", err)
		}
	}
}

func executeLine(line string) error {
	pipeline, err := parsePipeline(line)
	if err != nil {
		return fmt.Errorf("parsePipeline: %v", err)
	}

	return executePipeline(pipeline)
}

func parsePipeline(line string) (Pipeline, error) {
	pipe := Pipeline{Commands: []Command{}}
	parts := strings.Split(line, "|")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return Pipeline{}, fmt.Errorf("minishell: syntax error near `|'")
		}
		args := strings.Fields(part)
		if len(args) == 0 {
			return Pipeline{}, fmt.Errorf("minishell: syntax error near `|'")
		}

		pipe.Commands = append(pipe.Commands, Command{
			Args: args,
		})
	}

	return pipe, nil
}

func executePipeline(p Pipeline) error {
	if len(p.Commands) == 1 {
		return executeCommand(p.Commands[0])
	}
	// собираем исполняемые команды (если cd, то error)
	execCommands := []*exec.Cmd{}
	for _, cmd := range p.Commands {
		execCmd, err := buildExecCmd(cmd)
		if err != nil {
			return fmt.Errorf("buildExecCmd: %v", err)
		}

		execCommands = append(execCommands, execCmd)
	}

	// связываем команды между собой (подаём выход первого на вход последующего)
	readersWriters := []*os.File{}
	defer closeFiles(readersWriters)

	execCommands[0].Stdin = os.Stdin
	for i := 0; i < len(execCommands)-1; i++ {
		r, w, err := os.Pipe()
		if err != nil {
			closeFiles(readersWriters)
			return fmt.Errorf("os.Pipe(): %v", err)
		}
		readersWriters = append(readersWriters, r, w)

		execCommands[i].Stderr = os.Stderr
		execCommands[i].Stdout = w
		execCommands[i+1].Stdin = r
	}
	last := len(execCommands) - 1
	execCommands[last].Stdout = os.Stdout
	execCommands[last].Stderr = os.Stderr

	// первая команда создаёт process group
	execCommands[0].SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	if err := execCommands[0].Start(); err != nil {
		return fmt.Errorf("cmd.Start(): %v", err)
	}

	pgid := execCommands[0].Process.Pid
	setCurrentPGID(pgid)
	defer clearCurrentPGID()

	// запускаем все команды и добавляем в группу к первой команде
	for i := 1; i < len(execCommands); i++ {
		execCommands[i].SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
			Pgid:    pgid,
		}
		if err := execCommands[i].Start(); err != nil {
			closeFiles(readersWriters)
			return fmt.Errorf("cmd.Start(): %v", err)
		}
	}

	// закрывем читателя и писателя (для обработки yes | head -n 5)
	closeFiles(readersWriters)

	// Ждём все команды
	var waitErr error
	for i, cmd := range execCommands {
		err := cmd.Wait()
		if isIgnorablePipelineErr(err, p.Commands[i]) {
			continue
		}
		// возвращаем первую ошибку
		if waitErr == nil {
			waitErr = fmt.Errorf("cmd.Wait(): %v", err)
		}
	}

	return waitErr
}
func isInterruptErr(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "interrupt")
}
func closeFiles(files []*os.File) {
	for _, f := range files {
		if f != nil {
			_ = f.Close()
		}
	}
}
func isIgnorablePipelineErr(err error, cmd Command) bool {
	if err == nil {
		return true
	}

	s := err.Error()

	if strings.Contains(s, "broken pipe") {
		return true
	}

	if len(cmd.Args) > 0 && cmd.Args[0] == "grep" && strings.Contains(s, "exit status 1") {
		return true
	}

	if strings.Contains(s, "interrupt") {
		return true
	}

	return false
}
func buildExecCmd(command Command) (*exec.Cmd, error) {
	if command.Args[0] == "cd" {
		return nil, fmt.Errorf("cd cannot be used in pipeline")
	}
	return exec.Command(command.Args[0], command.Args[1:]...), nil
}

func executeCommand(command Command) error {
	handled, err := runCd(command)
	if handled {
		return err
	}

	cmd := exec.Command(command.Args[0], command.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	setCurrentPGID(cmd.Process.Pid)
	defer clearCurrentPGID()

	err = cmd.Wait()
	if err != nil {
		if isInterruptErr(err) {
			return nil
		}
	}
	return err
}

func runCd(cmd Command) (bool, error) {
	if len(cmd.Args) == 0 {
		return true, nil
	}
	if cmd.Args[0] == "cd" {
		if len(cmd.Args) < 2 {
			return true, fmt.Errorf("cd: missing path")
		}
		if len(cmd.Args) > 2 {
			return true, fmt.Errorf("cd: too many arguments")
		}
		return true, os.Chdir(cmd.Args[1])
	}
	return false, nil
}
