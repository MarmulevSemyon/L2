package main

import (
	"fmt"
	"runtime"
)

func main() {

	fmt.Println(runtime.GOMAXPROCS(4))
	// lineArgs, err := internal.ParseLine(os.Args)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err)
	// 	os.Exit(1)
	// }

	// fmt.Println(lineArgs)

}
