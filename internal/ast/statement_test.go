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

// GOTO Statement Tests - Following TDD: Write failing tests first

// TestGotoStatement_Execute_ValidLineNumber tests GOTO to a valid line number
func TestGotoStatement_Execute_ValidLineNumber(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Create a program with line numbers 10, 20, 30
	program := &Program{
		Lines: map[int]Statement{
			10: NewAssignmentStatement("X", NewLiteralExpression(runtime.NewNumericValue(1))),
			20: NewAssignmentStatement("Y", NewLiteralExpression(runtime.NewNumericValue(2))),
			30: NewAssignmentStatement("Z", NewLiteralExpression(runtime.NewNumericValue(3))),
		},
		Order: []int{10, 20, 30},
	}
	
	// Test GOTO 30
	stmt := &GotoStatement{
		LineNumber: 30,
		Program:    program,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	// Program counter should be set to 30
	assert.Equal(t, 30, env.ProgramCounter)
}

// TestGotoStatement_Execute_ForwardJump tests GOTO jumping forward in the program
func TestGotoStatement_Execute_ForwardJump(t *testing.T) {
	env := runtime.NewEnvironment()
	env.ProgramCounter = 10 // Currently at line 10
	
	program := &Program{
		Lines: map[int]Statement{
			10: NewAssignmentStatement("X", NewLiteralExpression(runtime.NewNumericValue(1))),
			20: NewAssignmentStatement("Y", NewLiteralExpression(runtime.NewNumericValue(2))),
			30: NewAssignmentStatement("Z", NewLiteralExpression(runtime.NewNumericValue(3))),
		},
		Order: []int{10, 20, 30},
	}
	
	// Test forward jump: GOTO 30 (from line 10)
	stmt := &GotoStatement{
		LineNumber: 30,
		Program:    program,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, 30, env.ProgramCounter)
}

// TestGotoStatement_Execute_BackwardJump tests GOTO jumping backward in the program
func TestGotoStatement_Execute_BackwardJump(t *testing.T) {
	env := runtime.NewEnvironment()
	env.ProgramCounter = 30 // Currently at line 30
	
	program := &Program{
		Lines: map[int]Statement{
			10: NewAssignmentStatement("X", NewLiteralExpression(runtime.NewNumericValue(1))),
			20: NewAssignmentStatement("Y", NewLiteralExpression(runtime.NewNumericValue(2))),
			30: NewAssignmentStatement("Z", NewLiteralExpression(runtime.NewNumericValue(3))),
		},
		Order: []int{10, 20, 30},
	}
	
	// Test backward jump: GOTO 10 (from line 30)
	stmt := &GotoStatement{
		LineNumber: 10,
		Program:    program,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, 10, env.ProgramCounter)
}

// TestGotoStatement_Execute_InvalidLineNumber tests GOTO to non-existent line number
func TestGotoStatement_Execute_InvalidLineNumber(t *testing.T) {
	env := runtime.NewEnvironment()
	
	program := &Program{
		Lines: map[int]Statement{
			10: NewAssignmentStatement("X", NewLiteralExpression(runtime.NewNumericValue(1))),
			20: NewAssignmentStatement("Y", NewLiteralExpression(runtime.NewNumericValue(2))),
		},
		Order: []int{10, 20},
	}
	
	// Test GOTO to non-existent line 50
	stmt := &GotoStatement{
		LineNumber: 50,
		Program:    program,
	}
	
	err := stmt.Execute(env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "line number 50 does not exist")
}

// TestGotoStatement_Execute_ZeroLineNumber tests GOTO to line number 0
func TestGotoStatement_Execute_ZeroLineNumber(t *testing.T) {
	env := runtime.NewEnvironment()
	
	program := &Program{
		Lines: map[int]Statement{
			10: NewAssignmentStatement("X", NewLiteralExpression(runtime.NewNumericValue(1))),
		},
		Order: []int{10},
	}
	
	// Test GOTO 0 (invalid line number)
	stmt := &GotoStatement{
		LineNumber: 0,
		Program:    program,
	}
	
	err := stmt.Execute(env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "line number 0 does not exist")
}

// TestGotoStatement_Execute_NegativeLineNumber tests GOTO to negative line number
func TestGotoStatement_Execute_NegativeLineNumber(t *testing.T) {
	env := runtime.NewEnvironment()
	
	program := &Program{
		Lines: map[int]Statement{
			10: NewAssignmentStatement("X", NewLiteralExpression(runtime.NewNumericValue(1))),
		},
		Order: []int{10},
	}
	
	// Test GOTO -10 (invalid line number)
	stmt := &GotoStatement{
		LineNumber: -10,
		Program:    program,
	}
	
	err := stmt.Execute(env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "line number -10 does not exist")
}

// TestGotoStatement_Execute_EmptyProgram tests GOTO in empty program
func TestGotoStatement_Execute_EmptyProgram(t *testing.T) {
	env := runtime.NewEnvironment()
	
	program := &Program{
		Lines: map[int]Statement{},
		Order: []int{},
	}
	
	// Test GOTO 10 in empty program
	stmt := &GotoStatement{
		LineNumber: 10,
		Program:    program,
	}
	
	err := stmt.Execute(env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "line number 10 does not exist")
}

// TestGotoStatement_Execute_SelfReference tests GOTO to current line (infinite loop potential)
func TestGotoStatement_Execute_SelfReference(t *testing.T) {
	env := runtime.NewEnvironment()
	env.ProgramCounter = 20 // Currently at line 20
	
	program := &Program{
		Lines: map[int]Statement{
			10: NewAssignmentStatement("X", NewLiteralExpression(runtime.NewNumericValue(1))),
			20: &GotoStatement{LineNumber: 20}, // This will be the statement itself
		},
		Order: []int{10, 20},
	}
	
	// Test GOTO 20 (self-reference)
	stmt := &GotoStatement{
		LineNumber: 20,
		Program:    program,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, 20, env.ProgramCounter)
}

// TestGotoStatement_Execute_ProgramCounterPreservation tests that other environment state is preserved
func TestGotoStatement_Execute_ProgramCounterPreservation(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Set up some environment state
	env.SetVariable("TEST", runtime.NewNumericValue(42))
	env.CallStack = []int{5, 15}
	originalCallStack := make([]int, len(env.CallStack))
	copy(originalCallStack, env.CallStack)
	
	program := &Program{
		Lines: map[int]Statement{
			10: NewAssignmentStatement("X", NewLiteralExpression(runtime.NewNumericValue(1))),
			20: NewAssignmentStatement("Y", NewLiteralExpression(runtime.NewNumericValue(2))),
		},
		Order: []int{10, 20},
	}
	
	// Test GOTO 20
	stmt := &GotoStatement{
		LineNumber: 20,
		Program:    program,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	// Program counter should change
	assert.Equal(t, 20, env.ProgramCounter)
	
	// Other state should be preserved
	assert.Equal(t, 42.0, env.GetVariable("TEST").NumValue)
	assert.Equal(t, originalCallStack, env.CallStack)
}

// TestGotoStatement_Execute_LargeLineNumbers tests GOTO with large line numbers
func TestGotoStatement_Execute_LargeLineNumbers(t *testing.T) {
	env := runtime.NewEnvironment()
	
	program := &Program{
		Lines: map[int]Statement{
			1000: NewAssignmentStatement("X", NewLiteralExpression(runtime.NewNumericValue(1))),
			9999: NewAssignmentStatement("Y", NewLiteralExpression(runtime.NewNumericValue(2))),
		},
		Order: []int{1000, 9999},
	}
	
	// Test GOTO 9999
	stmt := &GotoStatement{
		LineNumber: 9999,
		Program:    program,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, 9999, env.ProgramCounter)
}

// IF-THEN Statement Tests - Following TDD: Write failing tests first

// TestIfStatement_Execute_TrueCondition tests IF-THEN with true condition
func TestIfStatement_Execute_TrueCondition(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	// Test IF 5 > 3 THEN PRINT "True"
	condition := &ComparisonExpression{
		Left:     NewLiteralExpression(runtime.NewNumericValue(5)),
		Operator: ">",
		Right:    NewLiteralExpression(runtime.NewNumericValue(3)),
	}
	
	thenStatement := NewPrintStatement(
		[]Expression{NewLiteralExpression(runtime.NewStringValue("True"))},
		output,
	)
	
	stmt := &IfStatement{
		Condition:     condition,
		ThenStatement: thenStatement,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, "True", output.GetLastOutput())
}

// TestIfStatement_Execute_FalseCondition tests IF-THEN with false condition
func TestIfStatement_Execute_FalseCondition(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	// Test IF 3 > 5 THEN PRINT "False"
	condition := &ComparisonExpression{
		Left:     NewLiteralExpression(runtime.NewNumericValue(3)),
		Operator: ">",
		Right:    NewLiteralExpression(runtime.NewNumericValue(5)),
	}
	
	thenStatement := NewPrintStatement(
		[]Expression{NewLiteralExpression(runtime.NewStringValue("False"))},
		output,
	)
	
	stmt := &IfStatement{
		Condition:     condition,
		ThenStatement: thenStatement,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	// Should not execute THEN statement, so no output
	assert.Equal(t, 0, len(output.GetOutput()))
}

// TestIfStatement_Execute_EqualityOperator tests IF-THEN with equality comparison
func TestIfStatement_Execute_EqualityOperator(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	testCases := []struct {
		name      string
		left      float64
		right     float64
		operator  string
		shouldRun bool
	}{
		{"Equal values", 5, 5, "=", true},
		{"Unequal values", 5, 3, "=", false},
		{"Not equal - different values", 5, 3, "<>", true},
		{"Not equal - same values", 5, 5, "<>", false},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output.Clear()
			
			condition := &ComparisonExpression{
				Left:     NewLiteralExpression(runtime.NewNumericValue(tc.left)),
				Operator: tc.operator,
				Right:    NewLiteralExpression(runtime.NewNumericValue(tc.right)),
			}
			
			thenStatement := NewPrintStatement(
				[]Expression{NewLiteralExpression(runtime.NewStringValue("Executed"))},
				output,
			)
			
			stmt := &IfStatement{
				Condition:     condition,
				ThenStatement: thenStatement,
			}
			
			err := stmt.Execute(env)
			assert.NoError(t, err)
			
			if tc.shouldRun {
				assert.Equal(t, "Executed", output.GetLastOutput())
			} else {
				assert.Equal(t, 0, len(output.GetOutput()))
			}
		})
	}
}

// TestIfStatement_Execute_ComparisonOperators tests all comparison operators
func TestIfStatement_Execute_ComparisonOperators(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	testCases := []struct {
		name      string
		left      float64
		right     float64
		operator  string
		shouldRun bool
	}{
		{"Less than - true", 3, 5, "<", true},
		{"Less than - false", 5, 3, "<", false},
		{"Less than - equal", 5, 5, "<", false},
		{"Greater than - true", 5, 3, ">", true},
		{"Greater than - false", 3, 5, ">", false},
		{"Greater than - equal", 5, 5, ">", false},
		{"Less than or equal - less", 3, 5, "<=", true},
		{"Less than or equal - equal", 5, 5, "<=", true},
		{"Less than or equal - greater", 5, 3, "<=", false},
		{"Greater than or equal - greater", 5, 3, ">=", true},
		{"Greater than or equal - equal", 5, 5, ">=", true},
		{"Greater than or equal - less", 3, 5, ">=", false},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output.Clear()
			
			condition := &ComparisonExpression{
				Left:     NewLiteralExpression(runtime.NewNumericValue(tc.left)),
				Operator: tc.operator,
				Right:    NewLiteralExpression(runtime.NewNumericValue(tc.right)),
			}
			
			thenStatement := NewPrintStatement(
				[]Expression{NewLiteralExpression(runtime.NewStringValue("Executed"))},
				output,
			)
			
			stmt := &IfStatement{
				Condition:     condition,
				ThenStatement: thenStatement,
			}
			
			err := stmt.Execute(env)
			assert.NoError(t, err)
			
			if tc.shouldRun {
				assert.Equal(t, "Executed", output.GetLastOutput())
			} else {
				assert.Equal(t, 0, len(output.GetOutput()))
			}
		})
	}
}

// TestIfStatement_Execute_StringComparison tests string comparison in conditions
func TestIfStatement_Execute_StringComparison(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	testCases := []struct {
		name      string
		left      string
		right     string
		operator  string
		shouldRun bool
	}{
		{"String equality - true", "hello", "hello", "=", true},
		{"String equality - false", "hello", "world", "=", false},
		{"String inequality - true", "hello", "world", "<>", true},
		{"String inequality - false", "hello", "hello", "<>", false},
		{"String less than - true", "apple", "banana", "<", true},
		{"String less than - false", "banana", "apple", "<", false},
		{"String greater than - true", "banana", "apple", ">", true},
		{"String greater than - false", "apple", "banana", ">", false},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output.Clear()
			
			condition := &ComparisonExpression{
				Left:     NewLiteralExpression(runtime.NewStringValue(tc.left)),
				Operator: tc.operator,
				Right:    NewLiteralExpression(runtime.NewStringValue(tc.right)),
			}
			
			thenStatement := NewPrintStatement(
				[]Expression{NewLiteralExpression(runtime.NewStringValue("Executed"))},
				output,
			)
			
			stmt := &IfStatement{
				Condition:     condition,
				ThenStatement: thenStatement,
			}
			
			err := stmt.Execute(env)
			assert.NoError(t, err)
			
			if tc.shouldRun {
				assert.Equal(t, "Executed", output.GetLastOutput())
			} else {
				assert.Equal(t, 0, len(output.GetOutput()))
			}
		})
	}
}

// TestIfStatement_Execute_VariableComparison tests comparison with variables
func TestIfStatement_Execute_VariableComparison(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	// Set up variables
	env.SetVariable("X", runtime.NewNumericValue(10))
	env.SetVariable("Y", runtime.NewNumericValue(5))
	
	// Test IF X > Y THEN PRINT "X is greater"
	condition := &ComparisonExpression{
		Left:     NewVariableExpression("X"),
		Operator: ">",
		Right:    NewVariableExpression("Y"),
	}
	
	thenStatement := NewPrintStatement(
		[]Expression{NewLiteralExpression(runtime.NewStringValue("X is greater"))},
		output,
	)
	
	stmt := &IfStatement{
		Condition:     condition,
		ThenStatement: thenStatement,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, "X is greater", output.GetLastOutput())
}

// TestIfStatement_Execute_ComplexCondition tests complex expressions in conditions
func TestIfStatement_Execute_ComplexCondition(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	// Set up variables
	env.SetVariable("A", runtime.NewNumericValue(10))
	env.SetVariable("B", runtime.NewNumericValue(5))
	
	// Test IF A + B > 12 THEN PRINT "Sum is greater than 12"
	leftExpr := NewBinaryExpression(
		NewVariableExpression("A"),
		OpAdd,
		NewVariableExpression("B"),
	)
	
	condition := &ComparisonExpression{
		Left:     leftExpr,
		Operator: ">",
		Right:    NewLiteralExpression(runtime.NewNumericValue(12)),
	}
	
	thenStatement := NewPrintStatement(
		[]Expression{NewLiteralExpression(runtime.NewStringValue("Sum is greater than 12"))},
		output,
	)
	
	stmt := &IfStatement{
		Condition:     condition,
		ThenStatement: thenStatement,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, "Sum is greater than 12", output.GetLastOutput()) // 10 + 5 = 15 > 12
}

// TestIfStatement_Execute_TypeMismatch tests error handling for type mismatches
func TestIfStatement_Execute_TypeMismatch(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	// Test comparing number with string: IF 5 > "hello" THEN PRINT "Test"
	condition := &ComparisonExpression{
		Left:     NewLiteralExpression(runtime.NewNumericValue(5)),
		Operator: ">",
		Right:    NewLiteralExpression(runtime.NewStringValue("hello")),
	}
	
	thenStatement := NewPrintStatement(
		[]Expression{NewLiteralExpression(runtime.NewStringValue("Test"))},
		output,
	)
	
	stmt := &IfStatement{
		Condition:     condition,
		ThenStatement: thenStatement,
	}
	
	err := stmt.Execute(env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type mismatch")
}

// TestIfStatement_Execute_InvalidOperator tests error handling for invalid comparison operators
func TestIfStatement_Execute_InvalidOperator(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	// Test invalid operator: IF 5 ?? 3 THEN PRINT "Test"
	condition := &ComparisonExpression{
		Left:     NewLiteralExpression(runtime.NewNumericValue(5)),
		Operator: "??",
		Right:    NewLiteralExpression(runtime.NewNumericValue(3)),
	}
	
	thenStatement := NewPrintStatement(
		[]Expression{NewLiteralExpression(runtime.NewStringValue("Test"))},
		output,
	)
	
	stmt := &IfStatement{
		Condition:     condition,
		ThenStatement: thenStatement,
	}
	
	err := stmt.Execute(env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported comparison operator")
}

// TestIfStatement_Execute_ConditionEvaluationError tests error handling when condition evaluation fails
func TestIfStatement_Execute_ConditionEvaluationError(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	// Create a condition that will cause division by zero
	leftExpr := NewBinaryExpression(
		NewLiteralExpression(runtime.NewNumericValue(10)),
		OpDivide,
		NewLiteralExpression(runtime.NewNumericValue(0)),
	)
	
	condition := &ComparisonExpression{
		Left:     leftExpr,
		Operator: ">",
		Right:    NewLiteralExpression(runtime.NewNumericValue(5)),
	}
	
	thenStatement := NewPrintStatement(
		[]Expression{NewLiteralExpression(runtime.NewStringValue("Test"))},
		output,
	)
	
	stmt := &IfStatement{
		Condition:     condition,
		ThenStatement: thenStatement,
	}
	
	err := stmt.Execute(env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "division by zero")
}

// TestIfStatement_Execute_ThenStatementError tests error handling when THEN statement execution fails
func TestIfStatement_Execute_ThenStatementError(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Create a true condition
	condition := &ComparisonExpression{
		Left:     NewLiteralExpression(runtime.NewNumericValue(5)),
		Operator: ">",
		Right:    NewLiteralExpression(runtime.NewNumericValue(3)),
	}
	
	// Create a THEN statement that will fail (assignment to empty variable name)
	thenStatement := NewAssignmentStatement("", NewLiteralExpression(runtime.NewNumericValue(42)))
	
	stmt := &IfStatement{
		Condition:     condition,
		ThenStatement: thenStatement,
	}
	
	err := stmt.Execute(env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid variable name")
}

// TestIfStatement_Execute_GotoInThen tests GOTO statement in THEN clause
func TestIfStatement_Execute_GotoInThen(t *testing.T) {
	env := runtime.NewEnvironment()
	
	program := &Program{
		Lines: map[int]Statement{
			10: NewAssignmentStatement("X", NewLiteralExpression(runtime.NewNumericValue(1))),
			20: NewAssignmentStatement("Y", NewLiteralExpression(runtime.NewNumericValue(2))),
		},
		Order: []int{10, 20},
	}
	
	// Test IF 5 > 3 THEN GOTO 20
	condition := &ComparisonExpression{
		Left:     NewLiteralExpression(runtime.NewNumericValue(5)),
		Operator: ">",
		Right:    NewLiteralExpression(runtime.NewNumericValue(3)),
	}
	
	thenStatement := NewGotoStatement(20, program)
	
	stmt := &IfStatement{
		Condition:     condition,
		ThenStatement: thenStatement,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, 20, env.ProgramCounter)
}

// TestIfStatement_Execute_AssignmentInThen tests assignment statement in THEN clause
func TestIfStatement_Execute_AssignmentInThen(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Test IF 5 > 3 THEN X = 42
	condition := &ComparisonExpression{
		Left:     NewLiteralExpression(runtime.NewNumericValue(5)),
		Operator: ">",
		Right:    NewLiteralExpression(runtime.NewNumericValue(3)),
	}
	
	thenStatement := NewAssignmentStatement("X", NewLiteralExpression(runtime.NewNumericValue(42)))
	
	stmt := &IfStatement{
		Condition:     condition,
		ThenStatement: thenStatement,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, 42.0, env.GetVariable("X").NumValue)
}

// TestIfStatement_Execute_ZeroAsCondition tests zero and non-zero values as conditions
func TestIfStatement_Execute_ZeroAsCondition(t *testing.T) {
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	
	testCases := []struct {
		name      string
		value     float64
		shouldRun bool
	}{
		{"Zero is false", 0, false},
		{"Positive number is true", 5, true},
		{"Negative number is true", -3, true},
		{"Decimal is true", 0.1, true},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output.Clear()
			
			// Test IF value THEN PRINT "Executed" (direct value as condition)
			condition := &ComparisonExpression{
				Left:     NewLiteralExpression(runtime.NewNumericValue(tc.value)),
				Operator: "<>",
				Right:    NewLiteralExpression(runtime.NewNumericValue(0)),
			}
			
			thenStatement := NewPrintStatement(
				[]Expression{NewLiteralExpression(runtime.NewStringValue("Executed"))},
				output,
			)
			
			stmt := &IfStatement{
				Condition:     condition,
				ThenStatement: thenStatement,
			}
			
			err := stmt.Execute(env)
			assert.NoError(t, err)
			
			if tc.shouldRun {
				assert.Equal(t, "Executed", output.GetLastOutput())
			} else {
				assert.Equal(t, 0, len(output.GetOutput()))
			}
		})
	}
}

// FOR-NEXT Loop Tests - Following TDD: Write failing tests first

// TestForStatement_Execute_SimpleLoop tests basic FOR loop initialization
func TestForStatement_Execute_SimpleLoop(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Test FOR I = 1 TO 5
	stmt := &ForStatement{
		Variable:  "I",
		StartExpr: NewLiteralExpression(runtime.NewNumericValue(1)),
		EndExpr:   NewLiteralExpression(runtime.NewNumericValue(5)),
		StepExpr:  NewLiteralExpression(runtime.NewNumericValue(1)), // Default step
		LineNum:   10,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	// Check that loop variable is initialized
	assert.Equal(t, 1.0, env.GetVariable("I").NumValue)
	
	// Check that loop state is pushed onto stack
	assert.Equal(t, 1, len(env.ForLoops))
	assert.Equal(t, "I", env.ForLoops[0].Variable)
	assert.Equal(t, 1.0, env.ForLoops[0].Current)
	assert.Equal(t, 5.0, env.ForLoops[0].End)
	assert.Equal(t, 1.0, env.ForLoops[0].Step)
	assert.Equal(t, 10, env.ForLoops[0].LineNum)
}

// TestForStatement_Execute_WithStep tests FOR loop with custom step
func TestForStatement_Execute_WithStep(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Test FOR I = 0 TO 10 STEP 2
	stmt := &ForStatement{
		Variable:  "I",
		StartExpr: NewLiteralExpression(runtime.NewNumericValue(0)),
		EndExpr:   NewLiteralExpression(runtime.NewNumericValue(10)),
		StepExpr:  NewLiteralExpression(runtime.NewNumericValue(2)),
		LineNum:   20,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	// Check that loop variable is initialized
	assert.Equal(t, 0.0, env.GetVariable("I").NumValue)
	
	// Check that loop state has correct step
	assert.Equal(t, 1, len(env.ForLoops))
	assert.Equal(t, 2.0, env.ForLoops[0].Step)
}

// TestForStatement_Execute_NegativeStep tests FOR loop with negative step
func TestForStatement_Execute_NegativeStep(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Test FOR I = 10 TO 1 STEP -1
	stmt := &ForStatement{
		Variable:  "I",
		StartExpr: NewLiteralExpression(runtime.NewNumericValue(10)),
		EndExpr:   NewLiteralExpression(runtime.NewNumericValue(1)),
		StepExpr:  NewLiteralExpression(runtime.NewNumericValue(-1)),
		LineNum:   30,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	// Check that loop variable is initialized
	assert.Equal(t, 10.0, env.GetVariable("I").NumValue)
	
	// Check that loop state has negative step
	assert.Equal(t, 1, len(env.ForLoops))
	assert.Equal(t, -1.0, env.ForLoops[0].Step)
}

// TestForStatement_Execute_VariableExpressions tests FOR loop with variable expressions
func TestForStatement_Execute_VariableExpressions(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Set up variables for loop bounds
	env.SetVariable("START", runtime.NewNumericValue(5))
	env.SetVariable("END", runtime.NewNumericValue(15))
	env.SetVariable("STEP", runtime.NewNumericValue(3))
	
	// Test FOR I = START TO END STEP STEP
	stmt := &ForStatement{
		Variable:  "I",
		StartExpr: NewVariableExpression("START"),
		EndExpr:   NewVariableExpression("END"),
		StepExpr:  NewVariableExpression("STEP"),
		LineNum:   40,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	// Check that expressions were evaluated correctly
	assert.Equal(t, 5.0, env.GetVariable("I").NumValue)
	assert.Equal(t, 1, len(env.ForLoops))
	assert.Equal(t, 5.0, env.ForLoops[0].Current)
	assert.Equal(t, 15.0, env.ForLoops[0].End)
	assert.Equal(t, 3.0, env.ForLoops[0].Step)
}

// TestForStatement_Execute_ComplexExpressions tests FOR loop with complex expressions
func TestForStatement_Execute_ComplexExpressions(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Set up variables
	env.SetVariable("A", runtime.NewNumericValue(2))
	env.SetVariable("B", runtime.NewNumericValue(3))
	
	// Test FOR I = A * B TO A * B + 10 STEP A
	startExpr := NewBinaryExpression(
		NewVariableExpression("A"),
		OpMultiply,
		NewVariableExpression("B"),
	)
	
	endExpr := NewBinaryExpression(
		NewBinaryExpression(
			NewVariableExpression("A"),
			OpMultiply,
			NewVariableExpression("B"),
		),
		OpAdd,
		NewLiteralExpression(runtime.NewNumericValue(10)),
	)
	
	stmt := &ForStatement{
		Variable:  "I",
		StartExpr: startExpr,
		EndExpr:   endExpr,
		StepExpr:  NewVariableExpression("A"),
		LineNum:   50,
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	// Check that complex expressions were evaluated correctly
	assert.Equal(t, 6.0, env.GetVariable("I").NumValue) // 2 * 3 = 6
	assert.Equal(t, 1, len(env.ForLoops))
	assert.Equal(t, 6.0, env.ForLoops[0].Current)
	assert.Equal(t, 16.0, env.ForLoops[0].End) // 6 + 10 = 16
	assert.Equal(t, 2.0, env.ForLoops[0].Step)
}

// TestForStatement_Execute_ExpressionError tests error handling when expressions fail
func TestForStatement_Execute_ExpressionError(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Create expression that will cause division by zero
	startExpr := NewBinaryExpression(
		NewLiteralExpression(runtime.NewNumericValue(10)),
		OpDivide,
		NewLiteralExpression(runtime.NewNumericValue(0)),
	)
	
	stmt := &ForStatement{
		Variable:  "I",
		StartExpr: startExpr,
		EndExpr:   NewLiteralExpression(runtime.NewNumericValue(5)),
		StepExpr:  NewLiteralExpression(runtime.NewNumericValue(1)),
		LineNum:   60,
	}
	
	err := stmt.Execute(env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "division by zero")
	
	// Should not have created loop state
	assert.Equal(t, 0, len(env.ForLoops))
}

// TestForStatement_Execute_ZeroStep tests error handling for zero step
func TestForStatement_Execute_ZeroStep(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Test FOR I = 1 TO 5 STEP 0 (invalid)
	stmt := &ForStatement{
		Variable:  "I",
		StartExpr: NewLiteralExpression(runtime.NewNumericValue(1)),
		EndExpr:   NewLiteralExpression(runtime.NewNumericValue(5)),
		StepExpr:  NewLiteralExpression(runtime.NewNumericValue(0)),
		LineNum:   70,
	}
	
	err := stmt.Execute(env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "step cannot be zero")
	
	// Should not have created loop state
	assert.Equal(t, 0, len(env.ForLoops))
}

// TestForStatement_Execute_EmptyVariableName tests error handling for invalid variable names
func TestForStatement_Execute_EmptyVariableName(t *testing.T) {
	env := runtime.NewEnvironment()
	
	stmt := &ForStatement{
		Variable:  "",
		StartExpr: NewLiteralExpression(runtime.NewNumericValue(1)),
		EndExpr:   NewLiteralExpression(runtime.NewNumericValue(5)),
		StepExpr:  NewLiteralExpression(runtime.NewNumericValue(1)),
		LineNum:   80,
	}
	
	err := stmt.Execute(env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid variable name")
}

// TestNextStatement_Execute_SimpleNext tests basic NEXT statement
func TestNextStatement_Execute_SimpleNext(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Set up a FOR loop state
	env.ForLoops = []runtime.ForLoopState{
		{
			Variable: "I",
			Current:  1.0,
			End:      5.0,
			Step:     1.0,
			LineNum:  10,
		},
	}
	env.SetVariable("I", runtime.NewNumericValue(1))
	
	// Test NEXT I
	stmt := &NextStatement{
		Variable: "I",
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	// Check that loop variable was incremented
	assert.Equal(t, 2.0, env.GetVariable("I").NumValue)
	
	// Check that loop state was updated
	assert.Equal(t, 1, len(env.ForLoops))
	assert.Equal(t, 2.0, env.ForLoops[0].Current)
	
	// Program counter should be set to continue loop
	assert.Equal(t, 10, env.ProgramCounter)
}

// TestNextStatement_Execute_LoopCompletion tests NEXT when loop is complete
func TestNextStatement_Execute_LoopCompletion(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Set up a FOR loop state at the end
	env.ForLoops = []runtime.ForLoopState{
		{
			Variable: "I",
			Current:  5.0,
			End:      5.0,
			Step:     1.0,
			LineNum:  10,
		},
	}
	env.SetVariable("I", runtime.NewNumericValue(5))
	
	// Test NEXT I
	stmt := &NextStatement{
		Variable: "I",
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	// Loop should be completed and removed from stack
	assert.Equal(t, 0, len(env.ForLoops))
	
	// Variable should have final value
	assert.Equal(t, 6.0, env.GetVariable("I").NumValue) // 5 + 1
	
	// Program counter should not be changed (continue to next line)
	assert.Equal(t, 0, env.ProgramCounter) // Default value, not set to loop start
}

// TestNextStatement_Execute_NegativeStep tests NEXT with negative step
func TestNextStatement_Execute_NegativeStep(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Set up a FOR loop state with negative step
	env.ForLoops = []runtime.ForLoopState{
		{
			Variable: "I",
			Current:  10.0,
			End:      1.0,
			Step:     -2.0,
			LineNum:  20,
		},
	}
	env.SetVariable("I", runtime.NewNumericValue(10))
	
	// Test NEXT I
	stmt := &NextStatement{
		Variable: "I",
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	// Check that loop variable was decremented
	assert.Equal(t, 8.0, env.GetVariable("I").NumValue) // 10 - 2
	
	// Check that loop state was updated
	assert.Equal(t, 1, len(env.ForLoops))
	assert.Equal(t, 8.0, env.ForLoops[0].Current)
	
	// Program counter should be set to continue loop
	assert.Equal(t, 20, env.ProgramCounter)
}

// TestNextStatement_Execute_NegativeStepCompletion tests NEXT completion with negative step
func TestNextStatement_Execute_NegativeStepCompletion(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Set up a FOR loop state at the end with negative step
	env.ForLoops = []runtime.ForLoopState{
		{
			Variable: "I",
			Current:  1.0,
			End:      1.0,
			Step:     -1.0,
			LineNum:  30,
		},
	}
	env.SetVariable("I", runtime.NewNumericValue(1))
	
	// Test NEXT I
	stmt := &NextStatement{
		Variable: "I",
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	// Loop should be completed and removed from stack
	assert.Equal(t, 0, len(env.ForLoops))
	
	// Variable should have final value
	assert.Equal(t, 0.0, env.GetVariable("I").NumValue) // 1 - 1
}

// TestNextStatement_Execute_NestedLoops tests nested FOR loops
func TestNextStatement_Execute_NestedLoops(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Set up nested FOR loop states
	env.ForLoops = []runtime.ForLoopState{
		{
			Variable: "I",
			Current:  1.0,
			End:      3.0,
			Step:     1.0,
			LineNum:  10,
		},
		{
			Variable: "J",
			Current:  1.0,
			End:      2.0,
			Step:     1.0,
			LineNum:  20,
		},
	}
	env.SetVariable("I", runtime.NewNumericValue(1))
	env.SetVariable("J", runtime.NewNumericValue(1))
	
	// Test NEXT J (inner loop)
	stmt := &NextStatement{
		Variable: "J",
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	// Inner loop variable should be incremented
	assert.Equal(t, 2.0, env.GetVariable("J").NumValue)
	
	// Both loops should still be on stack
	assert.Equal(t, 2, len(env.ForLoops))
	assert.Equal(t, 2.0, env.ForLoops[1].Current) // Inner loop updated
	assert.Equal(t, 1.0, env.ForLoops[0].Current) // Outer loop unchanged
	
	// Program counter should be set to inner loop start
	assert.Equal(t, 20, env.ProgramCounter)
}

// TestNextStatement_Execute_NestedLoopCompletion tests completion of nested loop
func TestNextStatement_Execute_NestedLoopCompletion(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Set up nested FOR loop states with inner loop at end
	env.ForLoops = []runtime.ForLoopState{
		{
			Variable: "I",
			Current:  1.0,
			End:      3.0,
			Step:     1.0,
			LineNum:  10,
		},
		{
			Variable: "J",
			Current:  2.0,
			End:      2.0,
			Step:     1.0,
			LineNum:  20,
		},
	}
	env.SetVariable("I", runtime.NewNumericValue(1))
	env.SetVariable("J", runtime.NewNumericValue(2))
	
	// Test NEXT J (inner loop completion)
	stmt := &NextStatement{
		Variable: "J",
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	// Inner loop should be completed and removed
	assert.Equal(t, 1, len(env.ForLoops))
	assert.Equal(t, "I", env.ForLoops[0].Variable) // Only outer loop remains
	
	// Inner loop variable should have final value
	assert.Equal(t, 3.0, env.GetVariable("J").NumValue) // 2 + 1
}

// TestNextStatement_Execute_NoMatchingFor tests error when NEXT has no matching FOR
func TestNextStatement_Execute_NoMatchingFor(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// No FOR loops on stack
	
	// Test NEXT I
	stmt := &NextStatement{
		Variable: "I",
	}
	
	err := stmt.Execute(env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "NEXT without FOR")
}

// TestNextStatement_Execute_WrongVariable tests error when NEXT variable doesn't match
func TestNextStatement_Execute_WrongVariable(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Set up a FOR loop state for variable I
	env.ForLoops = []runtime.ForLoopState{
		{
			Variable: "I",
			Current:  1.0,
			End:      5.0,
			Step:     1.0,
			LineNum:  10,
		},
	}
	
	// Test NEXT J (wrong variable)
	stmt := &NextStatement{
		Variable: "J",
	}
	
	err := stmt.Execute(env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "NEXT J without matching FOR J")
}

// TestNextStatement_Execute_CaseInsensitiveVariable tests case-insensitive variable matching
func TestNextStatement_Execute_CaseInsensitiveVariable(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Set up a FOR loop state with lowercase variable
	env.ForLoops = []runtime.ForLoopState{
		{
			Variable: "counter",
			Current:  1.0,
			End:      3.0,
			Step:     1.0,
			LineNum:  40,
		},
	}
	env.SetVariable("counter", runtime.NewNumericValue(1))
	
	// Test NEXT COUNTER (uppercase)
	stmt := &NextStatement{
		Variable: "COUNTER",
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	// Should work with case-insensitive matching
	assert.Equal(t, 2.0, env.GetVariable("counter").NumValue)
	assert.Equal(t, 1, len(env.ForLoops))
}

// TestNextStatement_Execute_EmptyVariable tests NEXT without variable (should match innermost loop)
func TestNextStatement_Execute_EmptyVariable(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Set up nested FOR loop states
	env.ForLoops = []runtime.ForLoopState{
		{
			Variable: "I",
			Current:  1.0,
			End:      3.0,
			Step:     1.0,
			LineNum:  10,
		},
		{
			Variable: "J",
			Current:  1.0,
			End:      2.0,
			Step:     1.0,
			LineNum:  20,
		},
	}
	env.SetVariable("I", runtime.NewNumericValue(1))
	env.SetVariable("J", runtime.NewNumericValue(1))
	
	// Test NEXT (no variable - should match innermost loop)
	stmt := &NextStatement{
		Variable: "", // Empty variable name
	}
	
	err := stmt.Execute(env)
	assert.NoError(t, err)
	
	// Should increment innermost loop variable (J)
	assert.Equal(t, 2.0, env.GetVariable("J").NumValue)
	assert.Equal(t, 1.0, env.GetVariable("I").NumValue) // Outer loop unchanged
	
	// Both loops should still be on stack
	assert.Equal(t, 2, len(env.ForLoops))
	assert.Equal(t, 2.0, env.ForLoops[1].Current) // Inner loop updated
}

// TestForNextIntegration_SimpleLoop tests complete FOR-NEXT loop execution
func TestForNextIntegration_SimpleLoop(t *testing.T) {
	env := runtime.NewEnvironment()
	
	// Test complete loop: FOR I = 1 TO 3 ... NEXT I
	forStmt := &ForStatement{
		Variable:  "I",
		StartExpr: NewLiteralExpression(runtime.NewNumericValue(1)),
		EndExpr:   NewLiteralExpression(runtime.NewNumericValue(3)),
		StepExpr:  NewLiteralExpression(runtime.NewNumericValue(1)),
		LineNum:   10,
	}
	
	nextStmt := &NextStatement{
		Variable: "I",
	}
	
	// Execute FOR statement
	err := forStmt.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, 1.0, env.GetVariable("I").NumValue)
	assert.Equal(t, 1, len(env.ForLoops))
	
	// Execute NEXT statement (iteration 1)
	err = nextStmt.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, 2.0, env.GetVariable("I").NumValue)
	assert.Equal(t, 1, len(env.ForLoops))
	assert.Equal(t, 10, env.ProgramCounter) // Should jump back to loop start
	
	// Execute NEXT statement (iteration 2)
	err = nextStmt.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, 3.0, env.GetVariable("I").NumValue)
	assert.Equal(t, 1, len(env.ForLoops))
	assert.Equal(t, 10, env.ProgramCounter) // Should jump back to loop start
	
	// Execute NEXT statement (iteration 3 - completion)
	err = nextStmt.Execute(env)
	assert.NoError(t, err)
	assert.Equal(t, 4.0, env.GetVariable("I").NumValue) // Final value after increment
	assert.Equal(t, 0, len(env.ForLoops)) // Loop completed and removed
	assert.Equal(t, 10, env.ProgramCounter) // Program counter from last iteration
}