package interpreter

import (
	"basic-interpreter/internal/ast"
	"basic-interpreter/internal/runtime"
)

// Interpreter represents the BASIC interpreter
type Interpreter struct {
	debugMode bool
}

// NewInterpreter creates a new interpreter instance
func NewInterpreter(debugMode bool) *Interpreter {
	return &Interpreter{
		debugMode: debugMode,
	}
}

// Execute executes a BASIC program
func (i *Interpreter) Execute(program *ast.Program, env *runtime.Environment) error {
	// Implementation will be added in later tasks
	return nil
}