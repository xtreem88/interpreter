package interpreter

import (
	"github.com/codecrafters-io/interpreter-starter-go/cmd/scanner"
)

type Environment struct {
	values map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{values: make(map[string]interface{})}
}

func (env *Environment) Define(name string, value interface{}) {
	env.values[name] = value
}

func (env *Environment) Get(name scanner.Token) (interface{}, error) {
	if value, ok := env.values[name.Lexeme]; ok {
		return value, nil
	}

	return nil, &RuntimeError{token: name, message: "Undefined variable '" + name.Lexeme + "'."}
}

func (env *Environment) Assign(name scanner.Token, value interface{}) error {
	if _, ok := env.values[name.Lexeme]; ok {
		env.values[name.Lexeme] = value
		return nil
	}

	return &RuntimeError{token: name, message: "Undefined variable '" + name.Lexeme + "'."}
}
