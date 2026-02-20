package main

import (
	"fmt"
	"main/internal"
	"os"
)

func main() {

	// fmt.Println(runtime.GOMAXPROCS(4))
	lineArgs, err := internal.ParseLine(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	less, err := internal.BuildLess(lineArgs)
	// fmt.Println(lineArgs)
	if lineArgs.C {
		res, err := internal.PrintIsSorted(os.Args[len(os.Args)-1], less)
		if err != nil {
			fmt.Printf("Произошла ошибка:\n%w", err)
		} else {
			fmt.Println(res)
		}
		return
	}

}
