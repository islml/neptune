package main

import (
	"fmt"
	"os"

	"github.com/islml/neptune/internal/lox"
)

func main() {
	app := lox.New(os.Stdin, os.Stdout, os.Stderr)

	exitCode, err := app.Execute(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	if exitCode != 0 {
		os.Exit(exitCode)
	}
}