package interpreter

import (
	"fmt"

	"rhomel.com/crafting-interpreters-go/pkg/scanner"
)

type Environment struct {
	values    map[string]interface{}
	enclosing *Environment
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		values:    make(map[string]interface{}),
		enclosing: enclosing,
	}
}

func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) Get(name scanner.Token) interface{} {
	if value, ok := e.values[name.Lexeme]; ok {
		return value
	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	panic(&RuntimeError{name, fmt.Sprintf("Undefined variable '%s'", name.Lexeme)})
}

func (e *Environment) Assign(name scanner.Token, value interface{}) {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return
	}
	if e.enclosing != nil {
		e.enclosing.Assign(name, value)
		return
	}
	panic(&RuntimeError{name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme)})
}
