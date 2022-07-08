package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	ExitCodeUsageError = 1
	ExitIOError        = 100
)

func exitf(code int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintln(os.Stderr)
	os.Exit(code)
}

func main() {
	args := &Args{}
	l := args.len()
	switch {
	case l > 1:
		exitf(ExitCodeUsageError, "usage: golox [file]")
	case l == 1:
		runFile(args.get()[0])
	default:
		runPrompt()
	}
}

type Args struct{}

func (a *Args) len() int {
	return len(a.get())
}

func (a *Args) get() []string {
	return os.Args[1:]
}

func runFile(file string) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		exitf(ExitIOError, "error reading file '%s': %v", file, err)
	}
	run(string(b))
}

func runPrompt() {
	// TODO
	fmt.Println("TODO")
}

func run(line string) {
}
