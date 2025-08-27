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

// AssignmentStatement represents a variable assignment statement
// Assigns the result of an expression to a variable
type AssignmentStatement struct {
	Variable   string
	Expression Expression
}

// Execute performs the variable assignment by evaluating the expression and storing the result
func (a *AssignmentStatement) Execute(env *runtime.Environment) error {
	// Validate variable name
	if err := ValidateVariableName(a.Variable); err != nil {
		return err
	}

	// Evaluate the expression
	value, err := a.Expression.Evaluate(env)
	if err != nil {
		return fmt.Errorf("error evaluating expression for assignment: %w", err)
	}

	// Store the value in the environment
	env.SetVariable(a.Variable, value)
	return nil
}

// OutputWriter interface for output operations (allows mocking in tests)
type OutputWriter interface {
	WriteLine(line string) error
}

// InputReader interface for input operations (allows mocking in tests)
type InputReader interface {
	ReadLine() (string, error)
}

// PrintStatement represents a PRINT statement that outputs expressions
type PrintStatement struct {
	Expressions []Expression
	Output      OutputWriter
}

// Execute performs the print operation by evaluating expressions and outputting them
func (p *PrintStatement) Execute(env *runtime.Environment) error {
	if len(p.Expressions) == 0 {
		// Empty PRINT statement outputs empty line
		return p.Output.WriteLine("")
	}

	// Evaluate all expressions and format output
	output, err := p.evaluateAndFormatExpressions(env)
	if err != nil {
		return err
	}

	return p.Output.WriteLine(output)
}

// evaluateAndFormatExpressions evaluates all expressions and formats them for output
func (p *PrintStatement) evaluateAndFormatExpressions(env *runtime.Environment) (string, error) {
	var parts []string
	for _, expr := range p.Expressions {
		value, err := expr.Evaluate(env)
		if err != nil {
			return "", fmt.Errorf("error evaluating expression for print: %w", err)
		}
		parts = append(parts, value.ToString())
	}

	return p.formatOutput(parts), nil
}

// formatOutput formats the evaluated parts into a single output string
func (p *PrintStatement) formatOutput(parts []string) string {
	// Join with spaces (default separator for comma-separated expressions)
	output := ""
	for i, part := range parts {
		if i > 0 {
			output += " "
		}
		output += part
	}
	return output
}

// InputStatement represents an INPUT statement that reads user input into a variable
type InputStatement struct {
	Prompt   string
	Variable string
	Input    InputReader
	Output   OutputWriter
}

// Execute performs the input operation by displaying prompt and reading input
func (i *InputStatement) Execute(env *runtime.Environment) error {
	// Validate variable name
	if err := ValidateVariableName(i.Variable); err != nil {
		return err
	}

	// Display prompt
	if err := i.displayPrompt(); err != nil {
		return fmt.Errorf("error displaying prompt: %w", err)
	}

	// Read and process input
	value, err := i.readAndConvertInput()
	if err != nil {
		return err
	}

	// Store the value in the environment
	env.SetVariable(i.Variable, value)
	return nil
}

// displayPrompt displays the input prompt to the user
func (i *InputStatement) displayPrompt() error {
	prompt := i.Prompt
	if prompt == "" {
		prompt = "? " // Default BASIC prompt
	}
	return i.Output.WriteLine(prompt)
}

// readAndConvertInput reads input and converts it to the appropriate type
func (i *InputStatement) readAndConvertInput() (runtime.Value, error) {
	input, err := i.Input.ReadLine()
	if err != nil {
		return runtime.Value{}, fmt.Errorf("error reading input: %w", err)
	}

	// Determine variable type and convert input accordingly
	if IsStringVariable(i.Variable) {
		// String variable - store input as-is
		return runtime.NewStringValue(input), nil
	}

	// Numeric variable - try to convert input to number
	if numValue, err := runtime.NewStringValue(input).ToNumber(); err == nil {
		return runtime.NewNumericValue(numValue), nil
	} else {
		return runtime.Value{}, fmt.Errorf("cannot convert input '%s' to number: %w", input, err)
	}
}

// isStringVariable checks if a variable name indicates a string variable (ends with $)
func (i *InputStatement) isStringVariable(name string) bool {
	return len(name) > 0 && name[len(name)-1] == '$'
}

// Helper functions for creating statements programmatically

// NewAssignmentStatement creates a new assignment statement
func NewAssignmentStatement(variable string, expression Expression) *AssignmentStatement {
	return &AssignmentStatement{
		Variable:   variable,
		Expression: expression,
	}
}

// NewPrintStatement creates a new print statement with the given expressions and output writer
func NewPrintStatement(expressions []Expression, output OutputWriter) *PrintStatement {
	return &PrintStatement{
		Expressions: expressions,
		Output:      output,
	}
}

// NewInputStatement creates a new input statement with the given parameters
func NewInputStatement(variable string, input InputReader, output OutputWriter) *InputStatement {
	return &InputStatement{
		Variable: variable,
		Input:    input,
		Output:   output,
	}
}

// NewInputStatementWithPrompt creates a new input statement with a custom prompt
func NewInputStatementWithPrompt(prompt, variable string, input InputReader, output OutputWriter) *InputStatement {
	return &InputStatement{
		Prompt:   prompt,
		Variable: variable,
		Input:    input,
		Output:   output,
	}
}

// Validation helper functions

// ValidateVariableName validates that a variable name is not empty
func ValidateVariableName(name string) error {
	if name == "" {
		return fmt.Errorf("invalid variable name: empty variable name")
	}
	return nil
}

// IsStringVariable checks if a variable name indicates a string variable (ends with $)
func IsStringVariable(name string) bool {
	return len(name) > 0 && name[len(name)-1] == '$'
}