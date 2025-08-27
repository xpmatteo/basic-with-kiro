package ast

import "basic-interpreter/internal/runtime"

// Statement represents any executable statement in BASIC
type Statement interface {
	Execute(env *runtime.Environment) error
}

// Expression represents any expression that can be evaluated to a value
type Expression interface {
	Evaluate(env *runtime.Environment) (runtime.Value, error)
}

// Program represents a complete BASIC program
type Program struct {
	Lines map[int]Statement // Line number -> Statement mapping
	Order []int             // Ordered list of line numbers
}