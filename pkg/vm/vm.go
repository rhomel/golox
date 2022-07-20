package vm

import "fmt"

var vm *VM
var DebugTraceExecution bool = true

const STACK_MAX = 256

type VM struct {
	Chunk *Chunk

	// instead of using unsafe pointers in Go we will use an array/slice index
	Ip int

	Stack    [STACK_MAX]Value
	StackTop int
}

func resetStack() {
	vm.StackTop = 0
}

func push(value Value) {
	vm.Stack[vm.StackTop] = value
	vm.StackTop++
}

func pop() Value {
	vm.StackTop--
	return vm.Stack[vm.StackTop]
}

func InitVM() {
	vm = &VM{
		Chunk: InitChunk(),
	}
}

type InterpretResult int

const (
	INTERPRET_OK InterpretResult = iota
	INTERPRET_COMPILE_ERROR
	INTERPRET_RUNTIME_ERROR
)

func interpret(source string) InterpretResult {
	compile(source)
	return INTERPRET_OK
}

func run() InterpretResult {
	READ_BYTE := func() uint8 {
		instruction := vm.Chunk.Code[vm.Ip]
		vm.Ip++
		return instruction
	}
	READ_CONSTANT := func() Value {
		return vm.Chunk.Constants.values[READ_BYTE()]
	}

	for {
		if DebugTraceExecution {
			fmt.Printf("          ")
			for i := 0; i < vm.StackTop; i++ {
				fmt.Printf("[ ")
				printValue(vm.Stack[i])
				fmt.Printf(" ]")
			}
			fmt.Println()
			vm.Chunk.DisassembleInstruction(vm.Ip)
		}
		var instruction uint8 = READ_BYTE()
		switch instruction {
		case OP_CONSTANT:
			constant := READ_CONSTANT()
			push(constant)
		case OP_ADD:
			b := pop()
			a := pop()
			push(a + b)
		case OP_SUBTRACT:
			b := pop()
			a := pop()
			push(a - b)
		case OP_MULTIPLY:
			b := pop()
			a := pop()
			push(a * b)
		case OP_DIVIDE:
			b := pop()
			a := pop()
			push(a / b)
		case OP_NEGATE:
			push(-pop())
		case OP_RETURN:
			printValue(pop())
			fmt.Println()
			return INTERPRET_OK
		default:
			// no-op
			// TODO
		}
	}

}
