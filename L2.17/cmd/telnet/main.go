package main

import (
	"fmt"
	"os"
	"telnet/internal/app"
	"telnet/internal/config"
)

func main() {
	arg := os.Args
	cfg, err := config.Parse(arg[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(cfg)
	if err := app.Run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
