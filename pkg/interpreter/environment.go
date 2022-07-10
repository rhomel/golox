package interpreter

import (
	"fmt"

	"rhomel.com/crafting-interpreters-go/pkg/scanner"
)

type Environment struct {
	values map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]interface{}),
	}
}

func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) Get(name scanner.Token) interface{} {
	if value, ok := e.values[name.Lexeme]; ok {
		return value
	}
	panic(&RuntimeError{name, fmt.Sprintf("Undefined variable '%s'", name.Lexeme)})
}

func (e *Environment) Assign(name scanner.Token, value interface{}) {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return
	}
	panic(&RuntimeError{name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme)})
}
