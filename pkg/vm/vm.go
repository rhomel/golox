package vm

import (
	"fmt"
	"math"
	"os"
)

var vm *VM
var DebugTraceExecution bool = true

const UINT8_COUNT = math.MaxUint8 + 1
const STACK_MAX = 256

type VM struct {
	Chunk *Chunk

	// instead of using unsafe pointers in Go we will use an array/slice index
	Ip int

	Stack    [STACK_MAX]Value
	StackTop int
	Strings  *Table
	Globals  *Table
	Objects  Obj
}

func resetStack() {
	vm.StackTop = 0
}

func runtimeError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	// the book does some C pointer math here that I don't remmeber so this
	// might not be correct.
	line := vm.Chunk.Lines[vm.Ip]
	fmt.Fprintf(os.Stderr, "[line %d] in script\n", line)
}

func push(value Value) {
	vm.Stack[vm.StackTop] = value
	vm.StackTop++
}

func pop() Value {
	vm.StackTop--
	return vm.Stack[vm.StackTop]
}

func peek(distance int) Value {
	return vm.Stack[vm.StackTop-1-distance]
}

func isFalsey(value Value) bool {
	return value.IsNil() || (value.IsBool() && !value.AsBool())
}

func concatenate() {
	b := AsString(pop())
	a := AsString(pop())
	push(ObjVal(takeString(a.String + b.String)))
}

func InitVM() {
	vm = &VM{
		Chunk:   InitChunk(),
		Globals: &Table{},
		Strings: &Table{},
	}
	vm.Globals.initTable()
	vm.Strings.initTable()
}

func FreeVM() {
	vm.Globals.freeTable()
	vm.Strings.freeTable()
	freeObjects()
}

type InterpretResult int

const (
	INTERPRET_OK InterpretResult = iota
	INTERPRET_COMPILE_ERROR
	INTERPRET_RUNTIME_ERROR
)

func interpret(source string) InterpretResult {
	chunk := InitChunk()
	if !compile(source, chunk) {
		return INTERPRET_COMPILE_ERROR
	}
	vm.Chunk = chunk
	vm.Ip = 0

	result := run()
	return result
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
	READ_STRING := func() *ObjectString {
		return AsString(READ_CONSTANT())
	}
	// we have to deviate from the book because we don't have C macros. So
	// instead of calling `push` here we simply check if both arguments are
	// numbers and if so return them.
	BINARY_OP := func() (float64, float64, InterpretResult) {
		if !peek(0).IsNumber() || !peek(1).IsNumber() {
			runtimeError("Operands must be numbers.")
			return 0, 0, INTERPRET_RUNTIME_ERROR
		}
		b := pop().AsNumber()
		a := pop().AsNumber()
		return a, b, INTERPRET_OK
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
		case OP_NIL:
			push(NilValue())
		case OP_TRUE:
			push(BooleanValue(true))
		case OP_FALSE:
			push(BooleanValue(false))
		case OP_POP:
			pop()
		case OP_GET_GLOBAL:
			name := READ_STRING()
			value := Value{}
			if !vm.Globals.Get(name, &value) {
				runtimeError("Undefined variable '%s'.", name.String)
				return INTERPRET_RUNTIME_ERROR
			}
			push(value)
		case OP_DEFINE_GLOBAL:
			name := READ_STRING()
			vm.Globals.Set(name, peek(0))
			pop()
		case OP_SET_GLOBAL:
			name := READ_STRING()
			if vm.Globals.Set(name, peek(0)) {
				vm.Globals.Delete(name)
				runtimeError("Undefined variable '%s'.", name.String)
				return INTERPRET_RUNTIME_ERROR
			}
		case OP_EQUAL:
			b := pop()
			a := pop()
			push(BooleanValue(ValuesEqual(a, b)))
		case OP_GREATER:
			a, b, i := BINARY_OP()
			if i != INTERPRET_OK {
				return i
			}
			push(BooleanValue(a > b))
		case OP_LESS:
			a, b, i := BINARY_OP()
			if i != INTERPRET_OK {
				return i
			}
			push(BooleanValue(a < b))
		case OP_ADD:
			if IsString(peek(0)) && IsString(peek(1)) {
				concatenate()
			} else if peek(0).IsNumber() && peek(1).IsNumber() {
				a, b, i := BINARY_OP()
				if i != INTERPRET_OK {
					return i
				}
				push(NumberValue(a + b))
			} else {
				runtimeError("Operands must be two numbers or two strings.")
				return INTERPRET_RUNTIME_ERROR
			}
		case OP_SUBTRACT:
			a, b, i := BINARY_OP()
			if i != INTERPRET_OK {
				return i
			}
			push(NumberValue(a - b))
		case OP_MULTIPLY:
			a, b, i := BINARY_OP()
			if i != INTERPRET_OK {
				return i
			}
			push(NumberValue(a * b))
		case OP_DIVIDE:
			a, b, i := BINARY_OP()
			if i != INTERPRET_OK {
				return i
			}
			push(NumberValue(a / b))
		case OP_NOT:
			push(BooleanValue(isFalsey(pop())))
		case OP_NEGATE:
			if !peek(0).IsNumber() {
				runtimeError("Operand must be a number.")
				return INTERPRET_RUNTIME_ERROR
			}
			push(NumberValue(-pop().AsNumber()))
		case OP_PRINT:
			printValue(pop())
			fmt.Println()
		case OP_RETURN:
			return INTERPRET_OK
		default:
			// no-op
			// TODO
		}
	}

}
