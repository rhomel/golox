package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"rhomel.com/crafting-interpreters-go/pkg/scanner"
)

const (
	ExitCodeOK         = 0
	ExitCodeUsageError = 1
	ExitIOError        = 100
)

func exitf(code int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintln(os.Stderr)
	os.Exit(code)
}

func main() {
	lox := &Lox{}
	args := &Args{}
	l := args.len()
	switch {
	case l > 1:
		exitf(ExitCodeUsageError, "usage: golox [file]")
	case l == 1:
		lox.runFile(args.get()[0])
	default:
		lox.runPrompt()
	}
}

type Args struct{}

func (a *Args) len() int {
	return len(a.get())
}

func (a *Args) get() []string {
	return os.Args[1:]
}

type Lox struct{}

func (l *Lox) runFile(file string) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		exitf(ExitIOError, "error reading file '%s': %v", file, err)
	}
	l.run(string(b))
}

func (l *Lox) runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			exitf(ExitCodeOK, "// #quit")
		}
		if err != nil {
			exitf(ExitIOError, "error reading from stdin: %v", err)
		}
		l.run(line)
	}
}

func (l *Lox) run(line string) {
	scanner := scanner.NewScanner(line)
	for token := range scanner.ScanTokens() {
		// TODO
		fmt.Println(token)
	}
}
