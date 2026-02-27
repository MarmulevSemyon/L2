package main

import (
	"fmt"
	"main/internal"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
)

func main() {
	lineArgs, remainingArgs, err := internal.ParseLine(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(remainingArgs) == 0 {
		fmt.Fprintln(os.Stderr, "Не указан входной файл")
		os.Exit(1)
	}
	// ---- CPU profile ----
	var cpuFile *os.File
	if lineArgs.CPUProfile != "" {
		cpuFile, err = os.Create(lineArgs.CPUProfile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка создания cpuprofile:", err)
			os.Exit(1)
		}
		if err := pprof.StartCPUProfile(cpuFile); err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка запуска CPU профиля:", err)
			_ = cpuFile.Close()
			os.Exit(1)
		}
		defer func() {
			pprof.StopCPUProfile()
			_ = cpuFile.Close()
		}()
	}

	// ---- Trace ----
	var traceFile *os.File
	if lineArgs.Trace != "" {
		traceFile, err = os.Create(lineArgs.Trace)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка создания trace:", err)
			os.Exit(1)
		}
		if err := trace.Start(traceFile); err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка запуска trace:", err)
			_ = traceFile.Close()
			os.Exit(1)
		}
		defer func() {
			trace.Stop()
			_ = traceFile.Close()
		}()
	}

	less, err := internal.BuildLess(lineArgs)

	fileName := remainingArgs[0]
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка открытия файла:", fileName, "\n", err)
		os.Exit(1)
	}
	defer file.Close()

	if lineArgs.C {
		res, err := internal.PrintIsSorted(file, less)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка:", err)
			os.Exit(1)
		}
		fmt.Println(res)
		return
	}

	// получаем назвния файлов с отсортированными строками
	fileNames, err := internal.MakeSortedChunks(file, less, runtime.GOMAXPROCS(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка создания файлов:", err)
		os.Exit(1)
	}

	// сливаем всё по порядку в бинарную кучу и оттуда во writer
	err = internal.MergeSortedFilesToWriterHeap(fileNames, os.Stdout, less, lineArgs.U)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка слияния файлов:", err)
		os.Exit(1)
	}

	writeMemProfile(lineArgs.MemProfile)
}

func writeMemProfile(path string) {
	if path == "" {
		return
	}
	// чтобы профайл отражал "после сборки мусора"
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
