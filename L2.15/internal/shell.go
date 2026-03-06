package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func Run() {
	reader := bufio.NewReader(os.Stdin)
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

	err = executePipeline(pipeline)
	if err != nil {
		return fmt.Errorf("executePipeline: %v", err)
	}
	return nil
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
	execCommands[0].Stdin = os.Stdin
	for i := 0; i < len(execCommands)-1; i++ {
		r, w, err := os.Pipe()
		if err != nil {
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

	// запускаем все команды
	for _, cmd := range execCommands {
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("cmd.Start(): %v", err)
		}
	}
	for _, rw := range readersWriters {
		if err := rw.Close(); err != nil {
			return fmt.Errorf("rw.Close(): %v", err)
		}
	}
	// Ждём все команды
	for _, cmd := range execCommands {
		if err := cmd.Wait(); err != nil {
			if strings.Contains(err.Error(), "broken pipe") {
				continue
			}
			return fmt.Errorf("cmd.Wait(): %v", err)
		}
	}

	return nil
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

	return cmd.Run()
}

func runCd(cmd Command) (bool, error) {
	if len(cmd.Args) == 0 {
		return false, nil
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
