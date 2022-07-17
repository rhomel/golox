package interpreter

import (
	"fmt"

	"github.com/rhomel/golox/pkg/scanner"
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
	panic(&RuntimeError{name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme)})
}

func (e *Environment) GetAt(distance int, name string) interface{} {
	if value, ok := e.ancestor(distance).values[name]; ok {
		return value
	}
	panic(fmt.Errorf("Resolution data error: undefined ancestor at distance %d, name '%s'", distance, name))
}

func (e *Environment) ancestor(distance int) *Environment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}
	return environment
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

func (e *Environment) AssignAt(distance int, name scanner.Token, value interface{}) {
	e.ancestor(distance).values[name.Lexeme] = value
}
