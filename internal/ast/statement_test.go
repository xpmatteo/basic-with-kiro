package ast

import (
	"basic-interpreter/internal/runtime"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAssignmentStatement_Execute_NumericVariable tests assignment of numeric values to variables
func TestAssignmentStatement_Execute_NumericVariable(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Test simple numeric assignment: X = 42
	stmt := &AssignmentStatement{
		Variable:   "X",
		Expression: NewLiteralExpression(runtime.NewNumericValue(42)),
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	value := env.GetVariable("X")
	assert.Equal(t, runtime.NumericValue, value.Type)
	assert.Equal(t, 42.0, value.NumValue)
}

// TestAssignmentStatement_Execute_StringVariable tests assignment of string values to variables
func TestAssignmentStatement_Execute_StringVariable(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Test string assignment: NAME$ = "Hello"
	stmt := &AssignmentStatement{
		Variable:   "NAME$",
		Expression: NewLiteralExpression(runtime.NewStringValue("Hello")),
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	value := env.GetVariable("NAME$")
	assert.Equal(t, runtime.StringValue, value.Type)
	assert.Equal(t, "Hello", value.StrValue)
}

// TestAssignmentStatement_Execute_ComplexExpression tests assignment with complex expressions
func TestAssignmentStatement_Execute_ComplexExpression(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Set up some variables for the expression
	env.SetVariable("A", runtime.NewNumericValue(10))
	env.SetVariable("B", runtime.NewNumericValue(5))
	
	// Test assignment with complex expression: RESULT = A + B * 2
	expr := NewBinaryExpression(
		NewVariableExpression("A"),
		OpAdd,
		NewBinaryExpression(
			NewVariableExpression("B"),
			OpMultiply,
			NewLiteralExpression(runtime.NewNumericValue(2)),
		),
	)
	
	stmt := &AssignmentStatement{
		Variable:   "RESULT",
		Expression: expr,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	value := env.GetVariable("RESULT")
	assert.Equal(t, runtime.NumericValue, value.Type)
	assert.Equal(t, 20.0, value.NumValue) // 10 + (5 * 2) = 20
}

// TestAssignmentStatement_Execute_CaseInsensitiveVariable tests case-insensitive variable names
func TestAssignmentStatement_Execute_CaseInsensitiveVariable(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Assign to lowercase variable
	stmt1 := &AssignmentStatement{
		Variable:   "counter",
		Expression: NewLiteralExpression(runtime.NewNumericValue(1)),
	}
	
	err := stmt1.Execute(env)
	assert.NoError(t, err)
	
	// Assign to uppercase version of same variable
	stmt2 := &AssignmentStatement{
		Variable:   "COUNTER",
		Expression: NewLiteralExpression(runtime.NewNumericValue(2)),
	}
	
	err = stmt2.Execute(env)
	assert.NoError(t, err)
	
	// Both should refer to the same variable
	value1 := env.GetVariable("counter")
	value2 := env.GetVariable("COUNTER")
	assert.Equal(t, 2.0, value1.NumValue)
	assert.Equal(t, 2.0, value2.NumValue)
}

// TestAssignmentStatement_Execute_VariableReassignment tests reassigning variables
func TestAssignmentStatement_Execute_VariableReassignment(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Initial assignment
	stmt1 := &AssignmentStatement{
		Variable:   "X",
		Expression: NewLiteralExpression(runtime.NewNumericValue(10)),
	}
	
	err := stmt1.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, 10.0, env.GetVariable("X").NumValue)
	
	// Reassignment
	stmt2 := &AssignmentStatement{
		Variable:   "X",
		Expression: NewLiteralExpression(runtime.NewNumericValue(20)),
	}
	
	err = stmt2.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, 20.0, env.GetVariable("X").NumValue)
}

// TestAssignmentStatement_Execute_ExpressionError tests error handling when expression evaluation fails
func TestAssignmentStatement_Execute_ExpressionError(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Create an expression that will cause division by zero
	expr := NewBinaryExpression(
		NewLiteralExpression(runtime.NewNumericValue(10)),
		OpDivide,
		NewLiteralExpression(runtime.NewNumericValue(0)),
	)
	
	stmt := &AssignmentStatement{
		Variable:   "X",
		Expression: expr,
	}
	
	err := stmt.Execute(env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "division by zero")
}

// TestAssignmentStatement_Execute_StringToNumericVariable tests type handling
func TestAssignmentStatement_Execute_StringToNumericVariable(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Assign string value to numeric variable (should work in BASIC)
	stmt := &AssignmentStatement{
		Variable:   "X",
		Expression: NewLiteralExpression(runtime.NewStringValue("42")),
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	value := env.GetVariable("X")
	assert.Equal(t, runtime.StringValue, value.Type)
	assert.Equal(t, "42", value.StrValue)
}

// TestAssignmentStatement_Execute_NumericToStringVariable tests type handling
func TestAssignmentStatement_Execute_NumericToStringVariable(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Assign numeric value to string variable (should work in BASIC)
	stmt := &AssignmentStatement{
		Variable:   "NAME$",
		Expression: NewLiteralExpression(runtime.NewNumericValue(123)),
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	value := env.GetVariable("NAME$")
	assert.Equal(t, runtime.NumericValue, value.Type)
	assert.Equal(t, 123.0, value.NumValue)
}

// TestAssignmentStatement_Execute_EmptyVariableName tests error handling for invalid variable names
func TestAssignmentStatement_Execute_EmptyVariableName(t *testing.T) {
	env := runtime.NewEnvironment()
	
	stmt := &AssignmentStatement{
		Variable:   "",
		Expression: NewLiteralExpression(runtime.NewNumericValue(42)),
	}
	
	err := stmt.Execute(env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid variable name")
}

// TestAssignmentStatement_Execute_SelfReference tests variable self-reference in expressions
func TestAssignmentStatement_Execute_SelfReference(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Set initial value
	env.SetVariable("COUNTER", runtime.NewNumericValue(5))
	
	// Test self-reference: COUNTER = COUNTER + 1
	expr := NewBinaryExpression(
		NewVariableExpression("COUNTER"),
		OpAdd,
		NewLiteralExpression(runtime.NewNumericValue(1)),
	)
	
	stmt := &AssignmentStatement{
		Variable:   "COUNTER",
		Expression: expr,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	value := env.GetVariable("COUNTER")
	assert.Equal(t, 6.0, value.NumValue)
}

// MockOutputWriter is a test double for capturing output during tests
type MockOutputWriter struct {
	outputs []string
}

func (m *MockOutputWriter) WriteLine(line string) error {
	m.outputs = append(m.outputs, line)
	return nil
}

func (m *MockOutputWriter) GetOutput() []string {
	return m.outputs
}

func (m *MockOutputWriter) GetLastOutput() string {
	if len(m.outputs) == 0 {
		return ""
	}
	return m.outputs[len(m.outputs)-1]
}

func (m *MockOutputWriter) Clear() {
	m.outputs = nil
}

// TestPrintStatement_Execute_SingleExpression tests printing a single expression
func TestPrintStatement_Execute_SingleExpression(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	// Test printing a numeric literal: PRINT 42
	stmt := &PrintStatement{
		Expressions: []Expression{NewLiteralExpression(runtime.NewNumericValue(42))},
		Output:      output,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, "42", output.GetLastOutput())
}

// TestPrintStatement_Execute_StringExpression tests printing a string expression
func TestPrintStatement_Execute_StringExpression(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	// Test printing a string literal: PRINT "Hello, World!"
	stmt := &PrintStatement{
		Expressions: []Expression{NewLiteralExpression(runtime.NewStringValue("Hello, World!"))},
		Output:      output,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", output.GetLastOutput())
}

// TestPrintStatement_Execute_MultipleExpressions tests printing multiple expressions
func TestPrintStatement_Execute_MultipleExpressions(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	// Test printing multiple expressions: PRINT "Value:", 42, "End"
	stmt := &PrintStatement{
		Expressions: []Expression{
			NewLiteralExpression(runtime.NewStringValue("Value:")),
			NewLiteralExpression(runtime.NewNumericValue(42)),
			NewLiteralExpression(runtime.NewStringValue("End")),
		},
		Output: output,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, "Value: 42 End", output.GetLastOutput())
}

// TestPrintStatement_Execute_VariableExpression tests printing variables
func TestPrintStatement_Execute_VariableExpression(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	// Set up variables
	env.SetVariable("NAME", runtime.NewStringValue("Alice"))
	env.SetVariable("AGE", runtime.NewNumericValue(25))
	
	// Test printing variables: PRINT NAME, AGE
	stmt := &PrintStatement{
		Expressions: []Expression{
			NewVariableExpression("NAME"),
			NewVariableExpression("AGE"),
		},
		Output: output,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, "Alice 25", output.GetLastOutput())
}

// TestPrintStatement_Execute_ComplexExpression tests printing complex expressions
func TestPrintStatement_Execute_ComplexExpression(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	// Set up variables
	env.SetVariable("A", runtime.NewNumericValue(10))
	env.SetVariable("B", runtime.NewNumericValue(5))
	
	// Test printing complex expression: PRINT A + B * 2
	expr := NewBinaryExpression(
		NewVariableExpression("A"),
		OpAdd,
		NewBinaryExpression(
			NewVariableExpression("B"),
			OpMultiply,
			NewLiteralExpression(runtime.NewNumericValue(2)),
		),
	)
	
	stmt := &PrintStatement{
		Expressions: []Expression{expr},
		Output:      output,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, "20", output.GetLastOutput()) // 10 + (5 * 2) = 20
}

// TestPrintStatement_Execute_EmptyPrint tests printing with no expressions
func TestPrintStatement_Execute_EmptyPrint(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	// Test empty PRINT statement
	stmt := &PrintStatement{
		Expressions: []Expression{},
		Output:      output,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, "", output.GetLastOutput())
}

// TestPrintStatement_Execute_NumericFormatting tests numeric value formatting
func TestPrintStatement_Execute_NumericFormatting(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	testCases := []struct {
		name     string
		value    float64
		expected string
	}{
		{"Integer", 42.0, "42"},
		{"Decimal", 3.14, "3.14"},
		{"Zero", 0.0, "0"},
		{"Negative", -5.5, "-5.5"},
		{"Large number", 1000000.0, "1e+06"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output.Clear()
			stmt := &PrintStatement{
				Expressions: []Expression{NewLiteralExpression(runtime.NewNumericValue(tc.value))},
				Output:      output,
			}
			
			err := stmt.Execute(env)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, output.GetLastOutput())
		})
	}
}

// TestPrintStatement_Execute_SpecialCharacters tests printing strings with special characters
func TestPrintStatement_Execute_SpecialCharacters(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Spaces", "  hello  ", "  hello  "},
		{"Quotes", `"quoted"`, `"quoted"`},
		{"Newlines", "line1\nline2", "line1\nline2"},
		{"Tabs", "tab\there", "tab\there"},
		{"Unicode", "Hello 世界", "Hello 世界"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output.Clear()
			stmt := &PrintStatement{
				Expressions: []Expression{NewLiteralExpression(runtime.NewStringValue(tc.input))},
				Output:      output,
			}
			
			err := stmt.Execute(env)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, output.GetLastOutput())
		})
	}
}

// TestPrintStatement_Execute_ExpressionError tests error handling when expression evaluation fails
func TestPrintStatement_Execute_ExpressionError(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	// Create an expression that will cause division by zero
	expr := NewBinaryExpression(
		NewLiteralExpression(runtime.NewNumericValue(10)),
		OpDivide,
		NewLiteralExpression(runtime.NewNumericValue(0)),
	)
	
	stmt := &PrintStatement{
		Expressions: []Expression{expr},
		Output:      output,
	}
	
	err := stmt.Execute(env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "division by zero")
}

// TestPrintStatement_Execute_MixedTypes tests printing mixed numeric and string types
func TestPrintStatement_Execute_MixedTypes(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	// Test mixed types: PRINT "Count:", 42, "items"
	stmt := &PrintStatement{
		Expressions: []Expression{
			NewLiteralExpression(runtime.NewStringValue("Count:")),
			NewLiteralExpression(runtime.NewNumericValue(42)),
			NewLiteralExpression(runtime.NewStringValue("items")),
		},
		Output: output,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, "Count: 42 items", output.GetLastOutput())
}

// TestPrintStatement_Execute_SeparatorHandling tests different separator behaviors
func TestPrintStatement_Execute_SeparatorHandling(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	// Test comma separator (default space)
	stmt := &PrintStatement{
		Expressions: []Expression{
			NewLiteralExpression(runtime.NewStringValue("A")),
			NewLiteralExpression(runtime.NewStringValue("B")),
			NewLiteralExpression(runtime.NewStringValue("C")),
		},
		Output: output,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, "A B C", output.GetLastOutput())
}

// MockInputReader is a test double for providing input during tests
type MockInputReader struct {
	inputs []string
	index  int
}

func (m *MockInputReader) ReadLine() (string, error) {
	if m.index >= len(m.inputs) {
		return "", fmt.Errorf("no more input available")
	}
	result := m.inputs[m.index]
	m.index++
	return result, nil
}

func (m *MockInputReader) SetInputs(inputs []string) {
	m.inputs = inputs
	m.index = 0
}

// TestInputStatement_Execute_NumericInput tests reading numeric input
func TestInputStatement_Execute_NumericInput(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	input := &MockInputReader{}
	input.SetInputs([]string{"42"})
	
	// Test INPUT X
	stmt := &InputStatement{
		Variable: "X",
		Input:    input,
		Output:   output,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	value := env.GetVariable("X")
	assert.Equal(t, runtime.NumericValue, value.Type)
	assert.Equal(t, 42.0, value.NumValue)
}

// TestInputStatement_Execute_StringInput tests reading string input
func TestInputStatement_Execute_StringInput(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	input := &MockInputReader{}
	input.SetInputs([]string{"Hello World"})
	
	// Test INPUT NAME$
	stmt := &InputStatement{
		Variable: "NAME$",
		Input:    input,
		Output:   output,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	value := env.GetVariable("NAME$")
	assert.Equal(t, runtime.StringValue, value.Type)
	assert.Equal(t, "Hello World", value.StrValue)
}

// TestInputStatement_Execute_WithPrompt tests INPUT with a prompt message
func TestInputStatement_Execute_WithPrompt(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	input := &MockInputReader{}
	input.SetInputs([]string{"25"})
	
	// Test INPUT "Enter age: "; AGE
	stmt := &InputStatement{
		Prompt:   "Enter age: ",
		Variable: "AGE",
		Input:    input,
		Output:   output,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	// Check that prompt was displayed
	assert.Equal(t, "Enter age: ", output.GetLastOutput())
	
	// Check that variable was set
	value := env.GetVariable("AGE")
	assert.Equal(t, 25.0, value.NumValue)
}

// TestInputStatement_Execute_NumericConversion tests automatic type conversion for numeric variables
func TestInputStatement_Execute_NumericConversion(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	input := &MockInputReader{}
	
	testCases := []struct {
		name     string
		input    string
		expected float64
	}{
		{"Integer", "42", 42.0},
		{"Decimal", "3.14", 3.14},
		{"Negative", "-5.5", -5.5},
		{"Zero", "0", 0.0},
		{"Leading spaces", "  123  ", 123.0},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input.SetInputs([]string{tc.input})
			
			stmt := &InputStatement{
				Variable: "NUM",
				Input:    input,
				Output:   output,
			}
			
			err := stmt.Execute(env)
			assert.NoError(t, err)
			
			value := env.GetVariable("NUM")
			assert.Equal(t, runtime.NumericValue, value.Type)
			assert.Equal(t, tc.expected, value.NumValue)
		})
	}
}

// TestInputStatement_Execute_StringVariableHandling tests string variable handling
func TestInputStatement_Execute_StringVariableHandling(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	input := &MockInputReader{}
	
	testCases := []struct {
		name     string
		variable string
		input    string
		expected string
	}{
		{"String variable", "NAME$", "Alice", "Alice"},
		{"Empty string", "EMPTY$", "", ""},
		{"Numeric input to string", "STR$", "123", "123"},
		{"Spaces preserved", "SPACES$", "  hello  ", "  hello  "},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input.SetInputs([]string{tc.input})
			
			stmt := &InputStatement{
				Variable: tc.variable,
				Input:    input,
				Output:   output,
			}
			
			err := stmt.Execute(env)
			assert.NoError(t, err)
			
			value := env.GetVariable(tc.variable)
			assert.Equal(t, runtime.StringValue, value.Type)
			assert.Equal(t, tc.expected, value.StrValue)
		})
	}
}

// TestInputStatement_Execute_InvalidNumericInput tests error handling for invalid numeric input
func TestInputStatement_Execute_InvalidNumericInput(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	input := &MockInputReader{}
	input.SetInputs([]string{"not_a_number"})
	
	// Test INPUT to numeric variable with invalid input
	stmt := &InputStatement{
		Variable: "NUM",
		Input:    input,
		Output:   output,
	}
	
	err := stmt.Execute(env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot convert")
}

// TestInputStatement_Execute_InputReadError tests error handling when input reading fails
func TestInputStatement_Execute_InputReadError(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	input := &MockInputReader{}
	// No inputs provided, should cause error
	
	stmt := &InputStatement{
		Variable: "X",
		Input:    input,
		Output:   output,
	}
	
	err := stmt.Execute(env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no more input available")
}

// TestInputStatement_Execute_EmptyVariableName tests error handling for invalid variable names
func TestInputStatement_Execute_EmptyVariableName(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	input := &MockInputReader{}
	input.SetInputs([]string{"42"})
	
	stmt := &InputStatement{
		Variable: "",
		Input:    input,
		Output:   output,
	}
	
	err := stmt.Execute(env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid variable name")
}

// TestInputStatement_Execute_CaseInsensitiveVariable tests case-insensitive variable handling
func TestInputStatement_Execute_CaseInsensitiveVariable(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	input := &MockInputReader{}
	input.SetInputs([]string{"100"})
	
	// Input to lowercase variable
	stmt := &InputStatement{
		Variable: "counter",
		Input:    input,
		Output:   output,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	// Check both lowercase and uppercase access
	value1 := env.GetVariable("counter")
	value2 := env.GetVariable("COUNTER")
	assert.Equal(t, 100.0, value1.NumValue)
	assert.Equal(t, 100.0, value2.NumValue)
}

// TestInputStatement_Execute_NoPrompt tests INPUT without explicit prompt
func TestInputStatement_Execute_NoPrompt(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	input := &MockInputReader{}
	input.SetInputs([]string{"test"})
	
	// Test INPUT without prompt
	stmt := &InputStatement{
		Variable: "VAR$",
		Input:    input,
		Output:   output,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	// Should show default prompt "? "
	assert.Equal(t, "? ", output.GetLastOutput())
	
	value := env.GetVariable("VAR$")
	assert.Equal(t, "test", value.StrValue)
}