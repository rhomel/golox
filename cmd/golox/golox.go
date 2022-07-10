package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	ast "rhomel.com/crafting-interpreters-go/pkg/ast/gen"
	"rhomel.com/crafting-interpreters-go/pkg/interpreter"
	"rhomel.com/crafting-interpreters-go/pkg/parser"
	"rhomel.com/crafting-interpreters-go/pkg/scanner"
	"rhomel.com/crafting-interpreters-go/pkg/util/ast/printer"
	"rhomel.com/crafting-interpreters-go/pkg/util/exit"
)

func main() {
	lox := NewLox()
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
	hadError        bool
	hadRuntimeError bool

	interpreter *interpreter.Interpreter
}

func NewLox() *Lox {
	lox := &Lox{}
	lox.interpreter = interpreter.NewInterpreter(lox)
	return lox
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
	if l.hadRuntimeError {
		exit.Exitf(exit.ExitRuntimeError, "")
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
	tokens := scanner.ScanTokens()
	//printTokens(tokens) // TODO: make a flag to enable printing scanned tokens
	parser := parser.NewParser(tokens, l)
	statements := parser.Parse()
	//printAst(expr) // TODO: make a flag to enable printing the parsed ast
	if !l.hadError {
		l.interpreter.Interpret(statements)
	}
}

func printAst(expr ast.Expr) {
	printer := &printer.AstPrinter{}
	fmt.Println(printer.Accept(expr))
}

func printTokens(tokens []*scanner.Token) {
	for _, token := range tokens {
		fmt.Printf("line: %d, token: %s\n", token.Line, token.String())
	}
}

func (l *Lox) Error(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) report(line int, where, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
	l.hadError = true
}

func (l *Lox) ParseError(token scanner.Token, message string) {
	if token.Typ == scanner.EOF {
		l.report(token.Line, " at end", message)
	} else {
		l.report(token.Line, " at '"+token.Lexeme+"'", message)
	}
}

func (l *Lox) RuntimeError(token scanner.Token, message string) {
	fmt.Printf("%s\n[line %d]\n", message, token.Line)
	l.hadRuntimeError = true
}
