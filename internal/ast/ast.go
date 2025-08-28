package ast

import (
	"basic-interpreter/internal/runtime"
	"fmt"
	"strings"
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

// BuiltinFunction interface represents a built-in function that can be called
type BuiltinFunction interface {
	Call(args []runtime.Value, env *runtime.Environment) (runtime.Value, error)
	Name() string
	ArgCount() int
}

// FunctionCallExpression represents a function call in an expression
type FunctionCallExpression struct {
	Name string
	Args []Expression
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

// Evaluate evaluates a function call by looking up the function and calling it with evaluated arguments
func (f *FunctionCallExpression) Evaluate(env *runtime.Environment) (runtime.Value, error) {
	// Look up the built-in function
	builtinFunc := GetBuiltinFunction(f.Name)
	if builtinFunc == nil {
		return runtime.Value{}, fmt.Errorf("unknown function: %s", f.Name)
	}

	// Validate argument count
	if len(f.Args) != builtinFunc.ArgCount() {
		return runtime.Value{}, fmt.Errorf("function %s expects %d argument(s), got %d", 
			f.Name, builtinFunc.ArgCount(), len(f.Args))
	}

	// Evaluate all arguments
	args := make([]runtime.Value, len(f.Args))
	for i, argExpr := range f.Args {
		value, err := argExpr.Evaluate(env)
		if err != nil {
			return runtime.Value{}, fmt.Errorf("error evaluating argument %d for function %s: %w", i, f.Name, err)
		}
		args[i] = value
	}

	// Call the function
	result, err := builtinFunc.Call(args, env)
	if err != nil {
		return runtime.Value{}, fmt.Errorf("error calling function %s: %w", f.Name, err)
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

// NewFunctionCallExpression creates a new function call expression with the given name and arguments
func NewFunctionCallExpression(name string, args []Expression) *FunctionCallExpression {
	return &FunctionCallExpression{
		Name: name,
		Args: args,
	}
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

// NewGotoStatement creates a new GOTO statement with the given line number and program reference
func NewGotoStatement(lineNumber int, program *Program) *GotoStatement {
	return &GotoStatement{
		LineNumber: lineNumber,
		Program:    program,
	}
}

// ComparisonExpression represents a comparison between two expressions
type ComparisonExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

// Evaluate performs the comparison operation and returns a boolean result as a numeric value
// In BASIC, true is represented as -1 and false as 0
func (c *ComparisonExpression) Evaluate(env *runtime.Environment) (runtime.Value, error) {
	// Evaluate left operand
	leftVal, err := c.Left.Evaluate(env)
	if err != nil {
		return runtime.Value{}, fmt.Errorf("error evaluating left operand in comparison: %w", err)
	}

	// Evaluate right operand
	rightVal, err := c.Right.Evaluate(env)
	if err != nil {
		return runtime.Value{}, fmt.Errorf("error evaluating right operand in comparison: %w", err)
	}

	// Perform the comparison
	result, err := c.performComparison(leftVal, rightVal)
	if err != nil {
		return runtime.Value{}, err
	}

	// Convert boolean result to BASIC numeric value (true = -1, false = 0)
	if result {
		return runtime.NewNumericValue(-1), nil
	}
	return runtime.NewNumericValue(0), nil
}

// performComparison performs the actual comparison based on the operator
func (c *ComparisonExpression) performComparison(left, right runtime.Value) (bool, error) {
	// Check for type compatibility
	if left.Type != right.Type {
		return false, fmt.Errorf("type mismatch in comparison: cannot compare %v with %v", left.Type, right.Type)
	}

	// Use a more efficient approach for comparison operations
	switch c.Operator {
	case "=":
		return c.compareEqual(left, right), nil
	case "<>":
		return !c.compareEqual(left, right), nil
	case "<":
		return c.compareValues(left, right, func(a, b float64) bool { return a < b }, func(a, b string) bool { return a < b })
	case ">":
		return c.compareValues(left, right, func(a, b float64) bool { return a > b }, func(a, b string) bool { return a > b })
	case "<=":
		return c.compareValues(left, right, func(a, b float64) bool { return a <= b }, func(a, b string) bool { return a <= b })
	case ">=":
		return c.compareValues(left, right, func(a, b float64) bool { return a >= b }, func(a, b string) bool { return a >= b })
	default:
		return false, fmt.Errorf("unsupported comparison operator: %s", c.Operator)
	}
}

// compareEqual checks if two values are equal
func (c *ComparisonExpression) compareEqual(left, right runtime.Value) bool {
	if left.Type == runtime.NumericValue {
		return left.NumValue == right.NumValue
	}
	return left.StrValue == right.StrValue
}

// compareValues is a generic comparison function that handles both numeric and string comparisons
func (c *ComparisonExpression) compareValues(left, right runtime.Value, numericCompare func(float64, float64) bool, stringCompare func(string, string) bool) (bool, error) {
	if left.Type == runtime.NumericValue {
		return numericCompare(left.NumValue, right.NumValue), nil
	}
	return stringCompare(left.StrValue, right.StrValue), nil
}

// IfStatement represents an IF-THEN conditional statement
type IfStatement struct {
	Condition     Expression
	ThenStatement Statement
}

// Execute performs the conditional execution by evaluating the condition and executing the THEN statement if true
func (i *IfStatement) Execute(env *runtime.Environment) error {
	// Evaluate the condition
	conditionValue, err := i.Condition.Evaluate(env)
	if err != nil {
		return fmt.Errorf("error evaluating IF condition: %w", err)
	}

	// Check if condition is true (non-zero for numeric values, non-empty for strings)
	if i.isConditionTrue(conditionValue) {
		// Execute the THEN statement
		if err := i.ThenStatement.Execute(env); err != nil {
			return fmt.Errorf("error executing THEN statement: %w", err)
		}
	}

	return nil
}

// isConditionTrue determines if a condition value should be considered true
// In BASIC, zero is false, non-zero is true for numbers; empty string is false, non-empty is true for strings
func (i *IfStatement) isConditionTrue(value runtime.Value) bool {
	if value.Type == runtime.NumericValue {
		return value.NumValue != 0
	}
	return value.StrValue != ""
}

// NewComparisonExpression creates a new comparison expression with the given operands and operator
func NewComparisonExpression(left Expression, operator string, right Expression) *ComparisonExpression {
	return &ComparisonExpression{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

// NewIfStatement creates a new IF-THEN statement with the given condition and THEN statement
func NewIfStatement(condition Expression, thenStatement Statement) *IfStatement {
	return &IfStatement{
		Condition:     condition,
		ThenStatement: thenStatement,
	}
}

// ForStatement represents a FOR loop statement
type ForStatement struct {
	Variable  string
	StartExpr Expression
	EndExpr   Expression
	StepExpr  Expression
	LineNum   int
}

// Execute initializes a FOR loop by evaluating expressions and setting up loop state
func (f *ForStatement) Execute(env *runtime.Environment) error {
	// Validate variable name
	if err := ValidateVariableName(f.Variable); err != nil {
		return err
	}

	// Evaluate expressions using helper function
	startValue, err := EvaluateNumericExpression(f.StartExpr, env, "FOR start")
	if err != nil {
		return err
	}

	endValue, err := EvaluateNumericExpression(f.EndExpr, env, "FOR end")
	if err != nil {
		return err
	}

	stepValue, err := EvaluateNumericExpression(f.StepExpr, env, "FOR step")
	if err != nil {
		return err
	}

	// Validate step is not zero
	if stepValue == 0 {
		return fmt.Errorf("FOR step cannot be zero")
	}

	// Set the loop variable to the start value
	env.SetVariable(f.Variable, runtime.NewNumericValue(startValue))

	// Create and push loop state onto the stack
	loopState := runtime.ForLoopState{
		Variable: f.Variable,
		Current:  startValue,
		End:      endValue,
		Step:     stepValue,
		LineNum:  f.LineNum,
	}

	env.ForLoops = append(env.ForLoops, loopState)

	return nil
}

// NextStatement represents a NEXT statement that continues or terminates a FOR loop
type NextStatement struct {
	Variable string
}

// Execute processes a NEXT statement by incrementing the loop variable and checking termination
func (n *NextStatement) Execute(env *runtime.Environment) error {
	// Check if there are any FOR loops on the stack
	if len(env.ForLoops) == 0 {
		return fmt.Errorf("NEXT without FOR")
	}

	// Find the matching FOR loop
	loopIndex := n.findMatchingLoop(env)
	if loopIndex == -1 {
		if n.Variable == "" {
			return fmt.Errorf("NEXT without FOR")
		}
		return fmt.Errorf("NEXT %s without matching FOR %s", n.Variable, n.Variable)
	}

	// Get the loop state
	loop := &env.ForLoops[loopIndex]

	// Increment the loop variable
	newValue := loop.Current + loop.Step
	env.SetVariable(loop.Variable, runtime.NewNumericValue(newValue))
	loop.Current = newValue

	// Check if loop should continue
	if n.shouldContinueLoop(loop) {
		// Continue loop - set program counter to loop start
		SetProgramCounter(env, loop.LineNum)
	} else {
		// Loop completed - remove from stack
		env.ForLoops = append(env.ForLoops[:loopIndex], env.ForLoops[loopIndex+1:]...)
	}

	return nil
}

// findMatchingLoop finds the FOR loop that matches this NEXT statement
func (n *NextStatement) findMatchingLoop(env *runtime.Environment) int {
	if n.Variable == "" {
		// Empty variable name matches the innermost (last) loop
		return len(env.ForLoops) - 1
	}

	// Search from innermost to outermost for matching variable name
	normalizedVar := NormalizeVariableName(n.Variable)
	for i := len(env.ForLoops) - 1; i >= 0; i-- {
		if NormalizeVariableName(env.ForLoops[i].Variable) == normalizedVar {
			return i
		}
	}

	return -1 // No matching loop found
}

// shouldContinueLoop determines if the loop should continue based on step direction and bounds
func (n *NextStatement) shouldContinueLoop(loop *runtime.ForLoopState) bool {
	if loop.Step > 0 {
		// Positive step: continue if current <= end
		return loop.Current <= loop.End
	} else {
		// Negative step: continue if current >= end
		return loop.Current >= loop.End
	}
}

// NewForStatement creates a new FOR statement with the given parameters
func NewForStatement(variable string, startExpr, endExpr, stepExpr Expression, lineNum int) *ForStatement {
	return &ForStatement{
		Variable:  variable,
		StartExpr: startExpr,
		EndExpr:   endExpr,
		StepExpr:  stepExpr,
		LineNum:   lineNum,
	}
}

// NewNextStatement creates a new NEXT statement with the given variable
func NewNextStatement(variable string) *NextStatement {
	return &NextStatement{
		Variable: variable,
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

// Program counter management helper functions

// SetProgramCounter sets the program counter to the specified line number
// This centralizes program counter management for control flow statements
func SetProgramCounter(env *runtime.Environment, lineNumber int) {
	env.ProgramCounter = lineNumber
}

// ValidateLineNumber checks if a line number exists in the program
func ValidateLineNumber(program *Program, lineNumber int) error {
	if program == nil {
		return fmt.Errorf("program reference is nil")
	}
	
	if _, exists := program.Lines[lineNumber]; !exists {
		return fmt.Errorf("line number %d does not exist", lineNumber)
	}
	
	return nil
}

// Expression evaluation helper functions

// EvaluateNumericExpression evaluates an expression and ensures it returns a numeric value
func EvaluateNumericExpression(expr Expression, env *runtime.Environment, context string) (float64, error) {
	value, err := expr.Evaluate(env)
	if err != nil {
		return 0, fmt.Errorf("error evaluating %s expression: %w", context, err)
	}
	
	if value.Type != runtime.NumericValue {
		return 0, fmt.Errorf("%s value must be numeric", context)
	}
	
	return value.NumValue, nil
}

// Variable name normalization helper

// NormalizeVariableName converts variable names to uppercase for case-insensitive operations
func NormalizeVariableName(name string) string {
	return strings.ToUpper(name)
}

// GotoStatement represents a GOTO statement that jumps to a specific line number
type GotoStatement struct {
	LineNumber int
	Program    *Program
}

// Execute performs the GOTO operation by setting the program counter to the target line
func (g *GotoStatement) Execute(env *runtime.Environment) error {
	// Validate that the target line number exists in the program
	if err := ValidateLineNumber(g.Program, g.LineNumber); err != nil {
		return err
	}
	
	// Set the program counter to the target line number
	SetProgramCounter(env, g.LineNumber)
	return nil
}

// Function registry for built-in functions
var builtinFunctions map[string]BuiltinFunction

// init initializes the built-in function registry
func init() {
	builtinFunctions = make(map[string]BuiltinFunction)
	
	// Register mathematical functions
	registerFunction(&AbsFunction{})
	registerFunction(&IntFunction{})
	registerFunction(&RndFunction{})
	
	// Register string functions
	registerFunction(&LenFunction{})
	registerFunction(&MidFunction{})
	registerFunction(&StrFunction{})
	registerFunction(&ValFunction{})
}

// registerFunction registers a built-in function in the registry
func registerFunction(fn BuiltinFunction) {
	builtinFunctions[fn.Name()] = fn
}

// GetRegisteredFunctionNames returns a list of all registered function names
func GetRegisteredFunctionNames() []string {
	names := make([]string, 0, len(builtinFunctions))
	for name := range builtinFunctions {
		names = append(names, name)
	}
	return names
}

// IsFunctionRegistered checks if a function is registered
func IsFunctionRegistered(name string) bool {
	normalizedName := strings.ToUpper(name)
	_, exists := builtinFunctions[normalizedName]
	return exists
}

// GetBuiltinFunction retrieves a built-in function by name (case-insensitive)
func GetBuiltinFunction(name string) BuiltinFunction {
	normalizedName := strings.ToUpper(name)
	return builtinFunctions[normalizedName]
}

// Helper functions for common validation patterns

// validateArgumentCount validates that the correct number of arguments is provided
func validateArgumentCount(functionName string, expected int, actual int) error {
	if actual != expected {
		if expected == 0 {
			return fmt.Errorf("%s function expected 0 arguments, got %d", functionName, actual)
		} else if expected == 1 {
			return fmt.Errorf("%s function expected 1 argument, got %d", functionName, actual)
		} else {
			return fmt.Errorf("%s function expected %d arguments, got %d", functionName, expected, actual)
		}
	}
	return nil
}

// validateNumericArgument validates that an argument is numeric
func validateNumericArgument(functionName string, argIndex int, arg runtime.Value) error {
	if arg.Type != runtime.NumericValue {
		if argIndex == 0 {
			return fmt.Errorf("%s function argument must be numeric", functionName)
		}
		return fmt.Errorf("%s function %s argument must be numeric", functionName, getOrdinal(argIndex+1))
	}
	return nil
}

// validateStringArgument validates that an argument is a string
func validateStringArgument(functionName string, argIndex int, arg runtime.Value) error {
	if arg.Type != runtime.StringValue {
		if argIndex == 0 {
			// Special case for functions with multiple arguments - be more specific
			if functionName == "MID$" || functionName == "VAL" {
				return fmt.Errorf("%s function %s argument must be string", functionName, getOrdinal(argIndex+1))
			}
			return fmt.Errorf("%s function argument must be string", functionName)
		}
		return fmt.Errorf("%s function %s argument must be string", functionName, getOrdinal(argIndex+1))
	}
	return nil
}

// getOrdinal returns the ordinal form of a number (1st, 2nd, 3rd, etc.)
func getOrdinal(n int) string {
	switch n {
	case 1:
		return "first"
	case 2:
		return "second"
	case 3:
		return "third"
	default:
		return fmt.Sprintf("%dth", n)
	}
}

// Built-in function implementations

// Mathematical Functions

// AbsFunction implements the ABS function (absolute value)
type AbsFunction struct{}

func (f *AbsFunction) Name() string { return "ABS" }
func (f *AbsFunction) ArgCount() int { return 1 }

func (f *AbsFunction) Call(args []runtime.Value, env *runtime.Environment) (runtime.Value, error) {
	if err := validateArgumentCount("ABS", 1, len(args)); err != nil {
		return runtime.Value{}, err
	}
	
	if err := validateNumericArgument("ABS", 0, args[0]); err != nil {
		return runtime.Value{}, err
	}
	
	result := args[0].NumValue
	if result < 0 {
		result = -result
	}
	
	return runtime.NewNumericValue(result), nil
}

// IntFunction implements the INT function (integer part)
type IntFunction struct{}

func (f *IntFunction) Name() string { return "INT" }
func (f *IntFunction) ArgCount() int { return 1 }

func (f *IntFunction) Call(args []runtime.Value, env *runtime.Environment) (runtime.Value, error) {
	if err := validateArgumentCount("INT", 1, len(args)); err != nil {
		return runtime.Value{}, err
	}
	
	if err := validateNumericArgument("INT", 0, args[0]); err != nil {
		return runtime.Value{}, err
	}
	
	result := float64(int(args[0].NumValue))
	return runtime.NewNumericValue(result), nil
}

// RndFunction implements the RND function (random number 0-1)
type RndFunction struct{}

func (f *RndFunction) Name() string { return "RND" }
func (f *RndFunction) ArgCount() int { return 0 }

func (f *RndFunction) Call(args []runtime.Value, env *runtime.Environment) (runtime.Value, error) {
	if err := validateArgumentCount("RND", 0, len(args)); err != nil {
		return runtime.Value{}, err
	}
	
	// Use the environment's random number generator
	result := env.Random()
	return runtime.NewNumericValue(result), nil
}

// String Functions

// LenFunction implements the LEN function (string length)
type LenFunction struct{}

func (f *LenFunction) Name() string { return "LEN" }
func (f *LenFunction) ArgCount() int { return 1 }

func (f *LenFunction) Call(args []runtime.Value, env *runtime.Environment) (runtime.Value, error) {
	if err := validateArgumentCount("LEN", 1, len(args)); err != nil {
		return runtime.Value{}, err
	}
	
	if err := validateStringArgument("LEN", 0, args[0]); err != nil {
		return runtime.Value{}, err
	}
	
	result := float64(len(args[0].StrValue))
	return runtime.NewNumericValue(result), nil
}

// MidFunction implements the MID$ function (substring extraction)
type MidFunction struct{}

func (f *MidFunction) Name() string { return "MID$" }
func (f *MidFunction) ArgCount() int { return 3 }

func (f *MidFunction) Call(args []runtime.Value, env *runtime.Environment) (runtime.Value, error) {
	if err := validateArgumentCount("MID$", 3, len(args)); err != nil {
		return runtime.Value{}, err
	}
	
	if err := validateStringArgument("MID$", 0, args[0]); err != nil {
		return runtime.Value{}, err
	}
	if err := validateNumericArgument("MID$", 1, args[1]); err != nil {
		return runtime.Value{}, err
	}
	if err := validateNumericArgument("MID$", 2, args[2]); err != nil {
		return runtime.Value{}, err
	}
	
	return f.extractSubstring(args[0].StrValue, args[1].NumValue, args[2].NumValue), nil
}

// extractSubstring extracts a substring using BASIC's 1-based indexing
func (f *MidFunction) extractSubstring(str string, start, length float64) runtime.Value {
	startIdx := int(start) - 1 // BASIC uses 1-based indexing
	lengthVal := int(length)
	
	// Handle invalid start index or negative length
	if startIdx < 0 || startIdx >= len(str) || lengthVal <= 0 {
		return runtime.NewStringValue("")
	}
	
	endIdx := startIdx + lengthVal
	if endIdx > len(str) {
		endIdx = len(str)
	}
	
	result := str[startIdx:endIdx]
	return runtime.NewStringValue(result)
}

// StrFunction implements the STR$ function (number to string conversion)
type StrFunction struct{}

func (f *StrFunction) Name() string { return "STR$" }
func (f *StrFunction) ArgCount() int { return 1 }

func (f *StrFunction) Call(args []runtime.Value, env *runtime.Environment) (runtime.Value, error) {
	if err := validateArgumentCount("STR$", 1, len(args)); err != nil {
		return runtime.Value{}, err
	}
	
	if err := validateNumericArgument("STR$", 0, args[0]); err != nil {
		return runtime.Value{}, err
	}
	
	result := args[0].ToString()
	return runtime.NewStringValue(result), nil
}

// ValFunction implements the VAL function (string to number conversion)
type ValFunction struct{}

func (f *ValFunction) Name() string { return "VAL" }
func (f *ValFunction) ArgCount() int { return 1 }

func (f *ValFunction) Call(args []runtime.Value, env *runtime.Environment) (runtime.Value, error) {
	if err := validateArgumentCount("VAL", 1, len(args)); err != nil {
		return runtime.Value{}, err
	}
	
	if err := validateStringArgument("VAL", 0, args[0]); err != nil {
		return runtime.Value{}, err
	}
	
	numValue, err := args[0].ToNumber()
	if err != nil {
		return runtime.Value{}, fmt.Errorf("VAL function cannot convert '%s' to number: %w", args[0].StrValue, err)
	}
	
	return runtime.NewNumericValue(numValue), nil
}