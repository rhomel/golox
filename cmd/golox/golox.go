package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime/pprof"

	"github.com/rhomel/golox/pkg/args"
	ast "github.com/rhomel/golox/pkg/ast/gen"
	"github.com/rhomel/golox/pkg/interpreter"
	"github.com/rhomel/golox/pkg/parser"
	"github.com/rhomel/golox/pkg/resolver"
	"github.com/rhomel/golox/pkg/scanner"
	"github.com/rhomel/golox/pkg/util/ast/printer"
	"github.com/rhomel/golox/pkg/util/exit"
	"github.com/rhomel/golox/pkg/vm"
)

func main() {
	implementation := flag.String("implementation", "treewalk", "interpreter implementation to use")
	cpuProfileFile := flag.String("cpu-profile", "", "file to output cpu profile")
	flag.Parse()
	args := args.New()
	switch *implementation {
	case "treewalk":
		treewalkMain(args, *cpuProfileFile)
	case "vm":
		vm.Main(args)
	default:
		exit.Exitf(exit.ExitCodeUsageError, fmt.Sprintf("%s is not a valid implementation flag value", *implementation))
	}
}

func treewalkMain(args *args.Args, cpuProfileFile string) {
	lox := NewLox()
	l := args.Len()
	switch {
	case l > 1:
		exit.Exitf(exit.ExitCodeUsageError, "usage: golox <flags> [file]")
	case l == 1:
		if cpuProfileFile != "" {
			cleanup := profile(cpuProfileFile)
			defer cleanup()
		}
		lox.runFile(args.Get()[0])
	default:
		lox.runPrompt()
	}
}

func profile(file string) func() {
	f, err := os.Create(file)
	if err != nil {
		exit.Exitf(exit.ExitIOError, "failed to create file %s", file)
	}
	pprof.StartCPUProfile(f)
	return pprof.StopCPUProfile
}

type Lox struct {
	hadError        bool
	hadRuntimeError bool

	interpreter interpreter.Interpreter
}

func NewLox() *Lox {
	lox := &Lox{}
	lox.interpreter = interpreter.NewTreeWalkInterpreter(lox)
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
	fmt.Println("Welcome to golox REPL. Use ctrl+d to exit.")
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
	if l.hadError {
		return
	}
	resolver := resolver.NewResolver(l.interpreter, l)
	resolver.ResolveStmts(statements)
	if l.hadError {
		return
	}
	l.interpreter.Interpret(statements)
}

func printAst(expr ast.Expr) {
	printer := &printer.AstPrinter{}
	fmt.Println(printer.Accept(expr))
}

func printTokens(tokens []*scanner.Token) {
	for _, token := range tokens {
		fmt.Fprintf(os.Stderr, "line: %d, token: %s\n", token.Line, token.String())
	}
}

func (l *Lox) Error(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) report(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	l.hadError = true
}

func (l *Lox) ParseError(token scanner.Token, message string) {
	if token.Typ == scanner.EOF {
		l.report(token.Line, " at end", message)
	} else {
		l.report(token.Line, " at '"+token.Lexeme+"'", message)
	}
}

func (l *Lox) ResolveError(token scanner.Token, message string) {
	l.ParseError(token, message)
}

func (l *Lox) RuntimeError(token scanner.Token, message string) {
	fmt.Fprintf(os.Stderr, "%s\n[line %d]\n", message, token.Line)
	l.hadRuntimeError = true
}
