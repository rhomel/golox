package vm

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/rhomel/golox/pkg/args"
	"github.com/rhomel/golox/pkg/util/exit"
)

func sampleChunk() {
	var constant int
	chunk := InitChunk()
	constant = chunk.AddConstant(1.2)
	chunk.Write(OP_CONSTANT, 123)
	chunk.Write(uint8(constant), 123)

	constant = chunk.AddConstant(3.4)
	chunk.Write(OP_CONSTANT, 123)
	chunk.Write(uint8(constant), 123)

	chunk.Write(OP_ADD, 123)

	constant = chunk.AddConstant(5.6)
	chunk.Write(OP_CONSTANT, 123)
	chunk.Write(uint8(constant), 123)

	chunk.Write(OP_DIVIDE, 123)
	chunk.Write(OP_NEGATE, 123)
	chunk.Write(OP_RETURN, 123)
	//chunk.Disassemble("test chunk")

	//interpret(chunk)
}

// Main is the entry point for the VM implementation of golox.
func Main(args *args.Args) {
	InitVM()

	if args.Len() == 0 {
		repl()
	} else if args.Len() == 1 {
		file := args.Get()[0]
		runFile(file)
	} else {
		exit.Exitf(exit.ExitCodeUsageError, "usage: golox <flags> [file]")
	}
}

func repl() {
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

		interpret(line)
	}
}

func runFile(file string) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		exit.Exitf(74, "error reading file '%s': %v", file, err)
	}
	result := interpret(string(b))
	if result == INTERPRET_COMPILE_ERROR {
		exit.Exitf(65, "compile error")
	}
	if result == INTERPRET_RUNTIME_ERROR {
		exit.Exitf(70, "runtime error")
	}

}
