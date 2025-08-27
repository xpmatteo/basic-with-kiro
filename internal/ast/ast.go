package ast

import (
	"basic-interpreter/internal/runtime"
	"fmt"
)

// Operator constants for better maintainability and type safety
const (
	OpAdd      = "+"
	OpSubtract = "-"
	OpMultiply = "*"
	OpDivide   = "/"
	OpPower    = "^"
)

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

// LiteralExpression represents a literal value (number or string)
// This is the simplest form of expression that directly holds a value
type LiteralExpression struct {
	Value runtime.Value
}

// Evaluate returns the literal value without any computation
func (l *LiteralExpression) Evaluate(env *runtime.Environment) (runtime.Value, error) {
	return l.Value, nil
}

// VariableExpression represents a variable reference in an expression
// Variables are resolved from the runtime environment
type VariableExpression struct {
	Name string
}

// Evaluate retrieves the variable value from the environment
// Uninitialized variables return appropriate default values (0 for numeric, "" for string)
func (v *VariableExpression) Evaluate(env *runtime.Environment) (runtime.Value, error) {
	return env.GetVariable(v.Name), nil
}

// BinaryExpression represents a binary operation between two expressions
// Supports arithmetic operations: +, -, *, /, ^ with proper operator precedence
type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

// Evaluate performs the binary operation by evaluating both operands and applying the operator
func (b *BinaryExpression) Evaluate(env *runtime.Environment) (runtime.Value, error) {
	// Evaluate left operand
	leftVal, err := b.Left.Evaluate(env)
	if err != nil {
		return runtime.Value{}, fmt.Errorf("error evaluating left operand: %w", err)
	}

	// Evaluate right operand
	rightVal, err := b.Right.Evaluate(env)
	if err != nil {
		return runtime.Value{}, fmt.Errorf("error evaluating right operand: %w", err)
	}

	// Apply the operator using the extracted operation logic
	return b.applyOperation(leftVal, rightVal)
}

// applyOperation applies the binary operator to two values
// This method encapsulates the operator dispatch logic for better maintainability
func (b *BinaryExpression) applyOperation(left, right runtime.Value) (runtime.Value, error) {
	switch b.Operator {
	case OpAdd:
		return left.Add(right)
	case OpSubtract:
		return left.Subtract(right)
	case OpMultiply:
		return left.Multiply(right)
	case OpDivide:
		return left.Divide(right)
	case OpPower:
		return left.Power(right)
	default:
		return runtime.Value{}, fmt.Errorf("unsupported operator: %s", b.Operator)
	}
}

// ParenthesesExpression represents an expression wrapped in parentheses
// Parentheses are used to override operator precedence in arithmetic expressions
type ParenthesesExpression struct {
	Expression Expression
}

// Evaluate evaluates the wrapped expression
// Parentheses don't change the evaluation logic, only precedence during parsing
func (p *ParenthesesExpression) Evaluate(env *runtime.Environment) (runtime.Value, error) {
	result, err := p.Expression.Evaluate(env)
	if err != nil {
		return runtime.Value{}, fmt.Errorf("error evaluating parenthesized expression: %w", err)
	}
	return result, nil
}

// Helper functions for creating expressions programmatically

// NewLiteralExpression creates a new literal expression with the given value
func NewLiteralExpression(value runtime.Value) *LiteralExpression {
	return &LiteralExpression{Value: value}
}

// NewVariableExpression creates a new variable expression with the given name
func NewVariableExpression(name string) *VariableExpression {
	return &VariableExpression{Name: name}
}

// NewBinaryExpression creates a new binary expression with the given operands and operator
func NewBinaryExpression(left Expression, operator string, right Expression) *BinaryExpression {
	return &BinaryExpression{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

// NewParenthesesExpression creates a new parentheses expression wrapping the given expression
func NewParenthesesExpression(expr Expression) *ParenthesesExpression {
	return &ParenthesesExpression{Expression: expr}
}

// IsValidOperator checks if the given operator is supported
func IsValidOperator(op string) bool {
	switch op {
	case OpAdd, OpSubtract, OpMultiply, OpDivide, OpPower:
		return true
	default:
		return false
	}
}

// GetOperatorPrecedence returns the precedence level of an operator
// Higher numbers indicate higher precedence (evaluated first)
func GetOperatorPrecedence(op string) int {
	switch op {
	case OpPower:
		return 3 // Highest precedence
	case OpMultiply, OpDivide:
		return 2 // Medium precedence
	case OpAdd, OpSubtract:
		return 1 // Lowest precedence
	default:
		return 0 // Unknown operator
	}
}