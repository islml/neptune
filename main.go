package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/islml/neptune/scanner"
)

type Lox struct {
	hadError bool
}

func main() {
	l := Lox{}
	if len(os.Args) > 2 {
		fmt.Println("Usage: neptune [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		l.RunFile(os.Args[1])
	} else {
		l.RunPrompt();
	}
}

func (l *Lox) RunFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	l.Run(string(bytes))
	if l.hadError {
		os.Exit(65)
	}
}

func (l *Lox) RunPrompt() {
	input := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if ok := input.Scan(); !ok {
			break
		}
		line := input.Text()
		l.Run(line)
		l.hadError = false
	}
}

func (l *Lox) Run(source string) {
	scanner := scanner.Scanner{
		Source: []rune(source),
	}
	tokens := scanner.ScanTokens()

	for _, token := range tokens {
		fmt.Println(token.String())
	}
}

func (l *Lox) Report(line int, where string, message string) {
	fmt.Printf("[line %d] Error %s: %s", line, where, message)
	l.hadError = true
}