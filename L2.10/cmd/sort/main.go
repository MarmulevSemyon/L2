package main

import (
	"fmt"
	"main/internal"
	"os"
	"runtime"
)

func main() {
	arg := os.Args

	lineArgs, err := internal.ParseLine(arg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	less, err := internal.BuildLess(lineArgs)

	fileName := arg[len(arg)-1]
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
}
