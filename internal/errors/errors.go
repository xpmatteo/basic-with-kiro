package errors

import "fmt"

// ErrorType represents the type of error that occurred
type ErrorType int

const (
	LexicalError ErrorType = iota
	SyntaxError
	RuntimeError
	LogicError
)

// BasicError represents an error in the BASIC interpreter
type BasicError struct {
	Type    ErrorType
	Message string
	Line    int
	Column  int
}

// Error implements the error interface
func (e BasicError) Error() string {
	var errorType string
	switch e.Type {
	case LexicalError:
		errorType = "Lexical Error"
	case SyntaxError:
		errorType = "Syntax Error"
	case RuntimeError:
		errorType = "Runtime Error"
	case LogicError:
		errorType = "Logic Error"
	default:
		errorType = "Unknown Error"
	}

	if e.Line > 0 {
		return fmt.Sprintf("%s at line %d: %s", errorType, e.Line, e.Message)
	}
	return fmt.Sprintf("%s: %s", errorType, e.Message)
}

// NewLexicalError creates a new lexical error
func NewLexicalError(message string, line, column int) *BasicError {
	return &BasicError{
		Type:    LexicalError,
		Message: message,
		Line:    line,
		Column:  column,
	}
}

// NewSyntaxError creates a new syntax error
func NewSyntaxError(message string, line int) *BasicError {
	return &BasicError{
		Type:    SyntaxError,
		Message: message,
		Line:    line,
	}
}

// NewRuntimeError creates a new runtime error
func NewRuntimeError(message string, line int) *BasicError {
	return &BasicError{
		Type:    RuntimeError,
		Message: message,
		Line:    line,
	}
}

// NewLogicError creates a new logic error
func NewLogicError(message string, line int) *BasicError {
	return &BasicError{
		Type:    LogicError,
		Message: message,
		Line:    line,
	}
}