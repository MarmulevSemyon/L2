package main

import (
	"fmt"
	"main/internal"
	"os"
)

func main() {

	fmt.Println(os.Args)
	lineArgs, err := internal.ParseLine(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(lineArgs)

}
