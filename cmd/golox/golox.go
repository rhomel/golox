package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"rhomel.com/crafting-interpreters-go/pkg/scanner"
	"rhomel.com/crafting-interpreters-go/pkg/util/exit"
)

func main() {
	lox := &Lox{}
	args := &Args{}
	l := args.len()
	switch {
	case l > 1:
		exit.Exitf(exit.ExitCodeUsageError, "usage: golox [file]")
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

type Lox struct {
	hadError bool
}

func (l *Lox) runFile(file string) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		exit.Exitf(exit.ExitIOError, "error reading file '%s': %v", file, err)
	}
	l.run(string(b))
	if l.hadError {
		exit.Exitf(exit.ExitSyntaxError, "")
	}
}

func (l *Lox) runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			exit.Exitf(exit.ExitCodeOK, "// #quit")
		}
		if err != nil {
			exit.Exitf(exit.ExitIOError, "error reading from stdin: %v", err)
		}
		l.run(line)
		l.hadError = false // don't exit the repl on syntax errors, just ignore the input
	}
}

func (l *Lox) run(line string) {
	scanner := scanner.NewScanner(line, l)
	for _, token := range scanner.ScanTokens() {
		// TODO
		fmt.Printf("line: %d, token: %s\n", token.Line, token.String())
	}
}

func (l *Lox) Error(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) report(line int, where, message string) {
	fmt.Printf("[line %d] Error%s: %s", line, where, message)
	l.hadError = true
}
