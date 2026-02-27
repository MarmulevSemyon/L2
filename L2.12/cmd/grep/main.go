package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"

	"l2.12/internal"
)

func main() {
	// парсим флаги
	flags, remainingArgs, err := internal.ParseLine(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// ---- CPU profile ----
	cpuCleanup, err := setupCPUProfile(flags.CPUProfile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "CPU profile:", err)
		os.Exit(1)
	}
	defer cpuCleanup()
	// ---- Trace ----
	traceCleanup, err := setupTrace(flags.Trace)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Trace:", err)
		os.Exit(1)
	}
	defer traceCleanup()
	// mem profile
	defer writeMemProfile(flags.MemProfile)

	// создаём функцию, определяющую наличие паттерна в строке
	pattern := remainingArgs[0]
	match, err := internal.BildMatcher(flags, pattern)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// открываем файл
	fileName := remainingArgs[1]
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка открытия файла:", fileName, "\n", err)
		os.Exit(1)
	}
	defer file.Close()
	// если есть флаг -с, то просто печатаем число строк
	if flags.Count {
		count, err := internal.CountOfMatch(match, file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(count)
		return
	}
	// печатает нужные строки из файла с нужным форматом
	err = internal.PrintGrep(file, match, flags)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func writeMemProfile(path string) {
	if path == "" {
		return
	}
	// чтобы профайл отражал после сборки мусора
	runtime.GC()

	f, err := os.Create(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка создания memprofile:", err)
		return
	}
	defer f.Close()

	if err := pprof.WriteHeapProfile(f); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка записи heap профиля:", err)
	}
}

func setupCPUProfile(path string) (cleanup func(), err error) {
	if path == "" {
		return func() {}, nil
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("create cpuprofile: %w", err)
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("start cpu profile: %w", err)
	}

	return func() {
		pprof.StopCPUProfile()
		_ = f.Close()
	}, nil
}
func setupTrace(path string) (cleanup func(), err error) {
	if path == "" {
		return func() {}, nil
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("create trace: %w", err)
	}

	if err := trace.Start(f); err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("start trace: %w", err)
	}

	return func() {
		trace.Stop()
		_ = f.Close()
	}, nil
}
