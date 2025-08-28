package parser

import (
	"basic-interpreter/internal/ast"
	"basic-interpreter/internal/lexer"
	"basic-interpreter/internal/runtime"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockOutputWriter for testing statements that require output
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

// MockInputReader for testing statements that require input
type MockInputReader struct {
	inputs []string
	index  int
}

func (m *MockInputReader) ReadLine() (string, error) {
	if m.index >= len(m.inputs) {
		return "", assert.AnError
	}
	result := m.inputs[m.index]
	m.index++
	return result, nil
}

func (m *MockInputReader) SetInputs(inputs []string) {
	m.inputs = inputs
	m.index = 0
}

// Helper function to create a parser from source code
func createParser(source string) *BasicParser {
	l := lexer.NewLexer(source)
	return NewParser(l)
}

// Test ParseStatement method - Assignment statements

func TestParser_ParseStatement_AssignmentNumeric(t *testing.T) {
	parser := createParser("X = 42")
	
	stmt, err := parser.ParseStatement()
	require.NoError(t, err)
	require.NotNil(t, stmt)
	
	assignment, ok := stmt.(*ast.AssignmentStatement)
	require.True(t, ok, "Expected AssignmentStatement")
	assert.Equal(t, "X", assignment.Variable)
	
	// Test that expression evaluates correctly
	env := runtime.NewEnvironment()
	value, err := assignment.Expression.Evaluate(env)
	require.NoError(t, err)
	assert.Equal(t, runtime.NumericValue, value.Type)
	assert.Equal(t, 42.0, value.NumValue)
}

func TestParser_ParseStatement_AssignmentString(t *testing.T) {
	parser := createParser(`NAME$ = "Hello"`)
	
	stmt, err := parser.ParseStatement()
	require.NoError(t, err)
	require.NotNil(t, stmt)
	
	assignment, ok := stmt.(*ast.AssignmentStatement)
	require.True(t, ok, "Expected AssignmentStatement")
	assert.Equal(t, "NAME$", assignment.Variable)
	
	// Test that expression evaluates correctly
	env := runtime.NewEnvironment()
	value, err := assignment.Expression.Evaluate(env)
	require.NoError(t, err)
	assert.Equal(t, runtime.StringValue, value.Type)
	assert.Equal(t, "Hello", value.StrValue)
}

func TestParser_ParseStatement_AssignmentExpression(t *testing.T) {
	parser := createParser("RESULT = A + B * 2")
	
	stmt, err := parser.ParseStatement()
	require.NoError(t, err)
	require.NotNil(t, stmt)
	
	assignment, ok := stmt.(*ast.AssignmentStatement)
	require.True(t, ok, "Expected AssignmentStatement")
	assert.Equal(t, "RESULT", assignment.Variable)
	
	// Test that expression evaluates correctly with variables
	env := runtime.NewEnvironment()
	env.SetVariable("A", runtime.NewNumericValue(10))
	env.SetVariable("B", runtime.NewNumericValue(5))
	
	value, err := assignment.Expression.Evaluate(env)
	require.NoError(t, err)
	assert.Equal(t, runtime.NumericValue, value.Type)
	assert.Equal(t, 20.0, value.NumValue) // 10 + (5 * 2) = 20
}

// Test ParseStatement method - PRINT statements

func TestParser_ParseStatement_PrintSingle(t *testing.T) {
	parser := createParser(`PRINT "Hello"`)
	
	stmt, err := parser.ParseStatement()
	require.NoError(t, err)
	require.NotNil(t, stmt)
	
	printStmt, ok := stmt.(*ast.PrintStatement)
	require.True(t, ok, "Expected PrintStatement")
	assert.Len(t, printStmt.Expressions, 1)
	
	// Test execution
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	printStmt.Output = output
	
	err = printStmt.Execute(env)
	require.NoError(t, err)
	assert.Equal(t, "Hello", output.GetLastOutput())
}

func TestParser_ParseStatement_PrintMultiple(t *testing.T) {
	parser := createParser(`PRINT "Value:", 42, "End"`)
	
	stmt, err := parser.ParseStatement()
	require.NoError(t, err)
	require.NotNil(t, stmt)
	
	printStmt, ok := stmt.(*ast.PrintStatement)
	require.True(t, ok, "Expected PrintStatement")
	assert.Len(t, printStmt.Expressions, 3)
	
	// Test execution
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	printStmt.Output = output
	
	err = printStmt.Execute(env)
	require.NoError(t, err)
	assert.Equal(t, "Value: 42 End", output.GetLastOutput())
}

func TestParser_ParseStatement_PrintVariable(t *testing.T) {
	parser := createParser("PRINT X")
	
	stmt, err := parser.ParseStatement()
	require.NoError(t, err)
	require.NotNil(t, stmt)
	
	printStmt, ok := stmt.(*ast.PrintStatement)
	require.True(t, ok, "Expected PrintStatement")
	assert.Len(t, printStmt.Expressions, 1)
	
	// Test execution with variable
	env := runtime.NewEnvironment()
	env.SetVariable("X", runtime.NewNumericValue(123))
	output := &MockOutputWriter{}
	printStmt.Output = output
	
	err = printStmt.Execute(env)
	require.NoError(t, err)
	assert.Equal(t, "123", output.GetLastOutput())
}

func TestParser_ParseStatement_PrintEmpty(t *testing.T) {
	parser := createParser("PRINT")
	
	stmt, err := parser.ParseStatement()
	require.NoError(t, err)
	require.NotNil(t, stmt)
	
	printStmt, ok := stmt.(*ast.PrintStatement)
	require.True(t, ok, "Expected PrintStatement")
	assert.Len(t, printStmt.Expressions, 0)
	
	// Test execution
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	printStmt.Output = output
	
	err = printStmt.Execute(env)
	require.NoError(t, err)
	assert.Equal(t, "", output.GetLastOutput())
}

// Test ParseStatement method - INPUT statements

func TestParser_ParseStatement_InputNumeric(t *testing.T) {
	parser := createParser("INPUT X")
	
	stmt, err := parser.ParseStatement()
	require.NoError(t, err)
	require.NotNil(t, stmt)
	
	inputStmt, ok := stmt.(*ast.InputStatement)
	require.True(t, ok, "Expected InputStatement")
	assert.Equal(t, "X", inputStmt.Variable)
	assert.Equal(t, "", inputStmt.Prompt)
	
	// Test execution
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	input := &MockInputReader{}
	input.SetInputs([]string{"42"})
	inputStmt.Output = output
	inputStmt.Input = input
	
	err = inputStmt.Execute(env)
	require.NoError(t, err)
	
	value := env.GetVariable("X")
	assert.Equal(t, runtime.NumericValue, value.Type)
	assert.Equal(t, 42.0, value.NumValue)
}

func TestParser_ParseStatement_InputString(t *testing.T) {
	parser := createParser("INPUT NAME$")
	
	stmt, err := parser.ParseStatement()
	require.NoError(t, err)
	require.NotNil(t, stmt)
	
	inputStmt, ok := stmt.(*ast.InputStatement)
	require.True(t, ok, "Expected InputStatement")
	assert.Equal(t, "NAME$", inputStmt.Variable)
	
	// Test execution
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	input := &MockInputReader{}
	input.SetInputs([]string{"Alice"})
	inputStmt.Output = output
	inputStmt.Input = input
	
	err = inputStmt.Execute(env)
	require.NoError(t, err)
	
	value := env.GetVariable("NAME$")
	assert.Equal(t, runtime.StringValue, value.Type)
	assert.Equal(t, "Alice", value.StrValue)
}

func TestParser_ParseStatement_InputWithPrompt(t *testing.T) {
	parser := createParser(`INPUT "Enter age: "; AGE`)
	
	stmt, err := parser.ParseStatement()
	require.NoError(t, err)
	require.NotNil(t, stmt)
	
	inputStmt, ok := stmt.(*ast.InputStatement)
	require.True(t, ok, "Expected InputStatement")
	assert.Equal(t, "AGE", inputStmt.Variable)
	assert.Equal(t, "Enter age: ", inputStmt.Prompt)
	
	// Test execution
	env := runtime.NewEnvironment()
	output := &MockOutputWriter{}
	input := &MockInputReader{}
	input.SetInputs([]string{"25"})
	inputStmt.Output = output
	inputStmt.Input = input
	
	err = inputStmt.Execute(env)
	require.NoError(t, err)
	
	assert.Equal(t, "Enter age: ", output.GetLastOutput())
	value := env.GetVariable("AGE")
	assert.Equal(t, 25.0, value.NumValue)
}

// Test ParseStatement method - GOTO statements

func TestParser_ParseStatement_Goto(t *testing.T) {
	parser := createParser("GOTO 100")
	
	stmt, err := parser.ParseStatement()
	require.NoError(t, err)
	require.NotNil(t, stmt)
	
	gotoStmt, ok := stmt.(*ast.GotoStatement)
	require.True(t, ok, "Expected GotoStatement")
	assert.Equal(t, 100, gotoStmt.LineNumber)
}

func TestParser_ParseStatement_GotoVariable(t *testing.T) {
	parser := createParser("GOTO X")
	
	_, err := parser.ParseStatement()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected line number after GOTO")
}

// Test ParseStatement method - IF-THEN statements

func TestParser_ParseStatement_IfThenSimple(t *testing.T) {
	parser := createParser(`IF X = 5 THEN PRINT "Five"`)
	
	stmt, err := parser.ParseStatement()
	require.NoError(t, err)
	require.NotNil(t, stmt)
	
	ifStmt, ok := stmt.(*ast.IfStatement)
	require.True(t, ok, "Expected IfStatement")
	
	// Test condition
	compExpr, ok := ifStmt.Condition.(*ast.ComparisonExpression)
	require.True(t, ok, "Expected ComparisonExpression")
	assert.Equal(t, "=", compExpr.Operator)
	
	// Test THEN statement
	printStmt, ok := ifStmt.ThenStatement.(*ast.PrintStatement)
	require.True(t, ok, "Expected PrintStatement in THEN clause")
	assert.Len(t, printStmt.Expressions, 1)
}

func TestParser_ParseStatement_IfThenGoto(t *testing.T) {
	parser := createParser("IF X > 10 THEN GOTO 200")
	
	stmt, err := parser.ParseStatement()
	require.NoError(t, err)
	require.NotNil(t, stmt)
	
	ifStmt, ok := stmt.(*ast.IfStatement)
	require.True(t, ok, "Expected IfStatement")
	
	// Test THEN statement is GOTO
	gotoStmt, ok := ifStmt.ThenStatement.(*ast.GotoStatement)
	require.True(t, ok, "Expected GotoStatement in THEN clause")
	assert.Equal(t, 200, gotoStmt.LineNumber)
}

func TestParser_ParseStatement_IfThenAssignment(t *testing.T) {
	parser := createParser("IF A < B THEN C = A + B")
	
	stmt, err := parser.ParseStatement()
	require.NoError(t, err)
	require.NotNil(t, stmt)
	
	ifStmt, ok := stmt.(*ast.IfStatement)
	require.True(t, ok, "Expected IfStatement")
	
	// Test THEN statement is assignment
	assignStmt, ok := ifStmt.ThenStatement.(*ast.AssignmentStatement)
	require.True(t, ok, "Expected AssignmentStatement in THEN clause")
	assert.Equal(t, "C", assignStmt.Variable)
}

// Test ParseStatement method - FOR statements

func TestParser_ParseStatement_ForBasic(t *testing.T) {
	parser := createParser("FOR I = 1 TO 10")
	
	stmt, err := parser.ParseStatement()
	require.NoError(t, err)
	require.NotNil(t, stmt)
	
	forStmt, ok := stmt.(*ast.ForStatement)
	require.True(t, ok, "Expected ForStatement")
	assert.Equal(t, "I", forStmt.Variable)
	
	// Test start expression
	env := runtime.NewEnvironment()
	startValue, err := forStmt.StartExpr.Evaluate(env)
	require.NoError(t, err)
	assert.Equal(t, 1.0, startValue.NumValue)
	
	// Test end expression
	endValue, err := forStmt.EndExpr.Evaluate(env)
	require.NoError(t, err)
	assert.Equal(t, 10.0, endValue.NumValue)
	
	// Test step expression (should default to 1)
	stepValue, err := forStmt.StepExpr.Evaluate(env)
	require.NoError(t, err)
	assert.Equal(t, 1.0, stepValue.NumValue)
}

func TestParser_ParseStatement_ForWithStep(t *testing.T) {
	parser := createParser("FOR J = 10 TO 1 STEP -1")
	
	stmt, err := parser.ParseStatement()
	require.NoError(t, err)
	require.NotNil(t, stmt)
	
	forStmt, ok := stmt.(*ast.ForStatement)
	require.True(t, ok, "Expected ForStatement")
	assert.Equal(t, "J", forStmt.Variable)
	
	// Test expressions
	env := runtime.NewEnvironment()
	startValue, err := forStmt.StartExpr.Evaluate(env)
	require.NoError(t, err)
	assert.Equal(t, 10.0, startValue.NumValue)
	
	endValue, err := forStmt.EndExpr.Evaluate(env)
	require.NoError(t, err)
	assert.Equal(t, 1.0, endValue.NumValue)
	
	stepValue, err := forStmt.StepExpr.Evaluate(env)
	require.NoError(t, err)
	assert.Equal(t, -1.0, stepValue.NumValue)
}

func TestParser_ParseStatement_ForWithExpressions(t *testing.T) {
	parser := createParser("FOR K = A + 1 TO B * 2 STEP C")
	
	stmt, err := parser.ParseStatement()
	require.NoError(t, err)
	require.NotNil(t, stmt)
	
	forStmt, ok := stmt.(*ast.ForStatement)
	require.True(t, ok, "Expected ForStatement")
	assert.Equal(t, "K", forStmt.Variable)
	
	// Test with variables set
	env := runtime.NewEnvironment()
	env.SetVariable("A", runtime.NewNumericValue(5))
	env.SetVariable("B", runtime.NewNumericValue(10))
	env.SetVariable("C", runtime.NewNumericValue(2))
	
	startValue, err := forStmt.StartExpr.Evaluate(env)
	require.NoError(t, err)
	assert.Equal(t, 6.0, startValue.NumValue) // A + 1 = 5 + 1 = 6
	
	endValue, err := forStmt.EndExpr.Evaluate(env)
	require.NoError(t, err)
	assert.Equal(t, 20.0, endValue.NumValue) // B * 2 = 10 * 2 = 20
	
	stepValue, err := forStmt.StepExpr.Evaluate(env)
	require.NoError(t, err)
	assert.Equal(t, 2.0, stepValue.NumValue) // C = 2
}

// Test ParseStatement method - NEXT statements

func TestParser_ParseStatement_NextWithVariable(t *testing.T) {
	parser := createParser("NEXT I")
	
	stmt, err := parser.ParseStatement()
	require.NoError(t, err)
	require.NotNil(t, stmt)
	
	nextStmt, ok := stmt.(*ast.NextStatement)
	require.True(t, ok, "Expected NextStatement")
	assert.Equal(t, "I", nextStmt.Variable)
}

func TestParser_ParseStatement_NextWithoutVariable(t *testing.T) {
	parser := createParser("NEXT")
	
	stmt, err := parser.ParseStatement()
	require.NoError(t, err)
	require.NotNil(t, stmt)
	
	nextStmt, ok := stmt.(*ast.NextStatement)
	require.True(t, ok, "Expected NextStatement")
	assert.Equal(t, "", nextStmt.Variable)
}

// Test ParseStatement method - Line number handling

func TestParser_ParseStatement_WithLineNumber(t *testing.T) {
	parser := createParser("10 PRINT \"Hello\"")
	
	// Skip the line number token manually since ParseStatement doesn't handle line numbers
	parser.nextToken() // skip line number
	
	stmt, err := parser.ParseStatement()
	require.NoError(t, err)
	require.NotNil(t, stmt)
	
	printStmt, ok := stmt.(*ast.PrintStatement)
	require.True(t, ok, "Expected PrintStatement")
	assert.Len(t, printStmt.Expressions, 1)
}

func TestParser_ParseStatement_LineNumberOnly(t *testing.T) {
	parser := createParser("100")
	
	stmt, err := parser.ParseStatement()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected assignment operator")
	assert.Nil(t, stmt)
}

// Test ParseStatement method - Error cases

func TestParser_ParseStatement_InvalidSyntax(t *testing.T) {
	testCases := []struct {
		name   string
		source string
		error  string
	}{
		{"Missing assignment operator", "X 42", "expected assignment operator"},
		{"Missing THEN in IF", "IF X = 5 PRINT \"Hello\"", "expected THEN after IF condition"},
		{"Missing TO in FOR", "FOR I = 1 10", "expected TO in FOR statement"},
		{"Invalid GOTO target", "GOTO \"hello\"", "expected line number after GOTO"},
		{"Empty statement", "", "unexpected end of input"},
		{"Invalid keyword", "INVALID X", "expected assignment operator"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := createParser(tc.source)
			
			stmt, err := parser.ParseStatement()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.error)
			assert.Nil(t, stmt)
		})
	}
}

func TestParser_ParseStatement_MalformedStatements(t *testing.T) {
	testCases := []struct {
		name   string
		source string
	}{
		{"Unterminated string", `PRINT "Hello`},
		{"Missing closing parenthesis", "PRINT (1 + 2"},
		{"Invalid operator", "X = @"},
		{"Missing variable in INPUT", "INPUT"},
		{"Missing expression in assignment", "X ="},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := createParser(tc.source)
			
			_, err := parser.ParseStatement()
			assert.Error(t, err)
		})
	}
}

// Test statement recognition

func TestParser_ParseStatement_StatementRecognition(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		expected string
	}{
		{"Assignment", "X = 5", "*ast.AssignmentStatement"},
		{"Print", "PRINT X", "*ast.PrintStatement"},
		{"Input", "INPUT Y", "*ast.InputStatement"},
		{"Goto", "GOTO 100", "*ast.GotoStatement"},
		{"If-Then", "IF X = 5 THEN PRINT X", "*ast.IfStatement"},
		{"For", "FOR I = 1 TO 10", "*ast.ForStatement"},
		{"Next", "NEXT I", "*ast.NextStatement"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := createParser(tc.source)
			
			stmt, err := parser.ParseStatement()
			require.NoError(t, err)
			require.NotNil(t, stmt)
			
			stmtType := fmt.Sprintf("%T", stmt)
			assert.Equal(t, tc.expected, stmtType)
		})
	}
}

// Test ParseExpression method - Operator precedence

func TestParser_ParseExpression_OperatorPrecedence(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		expected float64
	}{
		{"Addition and multiplication", "2 + 3 * 4", 14.0},        // 2 + (3 * 4) = 14
		{"Multiplication and addition", "3 * 4 + 2", 14.0},        // (3 * 4) + 2 = 14
		{"Division and subtraction", "10 - 6 / 2", 7.0},           // 10 - (6 / 2) = 7
		{"Power and multiplication", "2 * 3 ^ 2", 18.0},           // 2 * (3 ^ 2) = 18
		{"Complex precedence", "2 + 3 * 4 ^ 2 - 1", 49.0},        // 2 + (3 * (4 ^ 2)) - 1 = 49
		{"Left associativity", "10 - 5 - 2", 3.0},                 // (10 - 5) - 2 = 3
		{"Right associativity power", "2 ^ 3 ^ 2", 512.0},         // 2 ^ (3 ^ 2) = 512
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := createParser(tc.source)
			
			expr, err := parser.ParseExpression()
			require.NoError(t, err)
			require.NotNil(t, expr)
			
			// Test evaluation
			env := runtime.NewEnvironment()
			value, err := expr.Evaluate(env)
			require.NoError(t, err)
			assert.Equal(t, runtime.NumericValue, value.Type)
			assert.Equal(t, tc.expected, value.NumValue)
		})
	}
}

func TestParser_ParseExpression_ArithmeticOperators(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		expected float64
	}{
		{"Addition", "5 + 3", 8.0},
		{"Subtraction", "10 - 4", 6.0},
		{"Multiplication", "6 * 7", 42.0},
		{"Division", "15 / 3", 5.0},
		{"Power", "2 ^ 3", 8.0},
		{"Negative numbers", "-5 + 3", -2.0},
		{"Decimal numbers", "3.5 + 2.1", 5.6},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := createParser(tc.source)
			
			expr, err := parser.ParseExpression()
			require.NoError(t, err)
			require.NotNil(t, expr)
			
			// Test evaluation
			env := runtime.NewEnvironment()
			value, err := expr.Evaluate(env)
			require.NoError(t, err)
			assert.Equal(t, runtime.NumericValue, value.Type)
			assert.InDelta(t, tc.expected, value.NumValue, 0.0001)
		})
	}
}

func TestParser_ParseExpression_ComparisonOperators(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		expected float64 // -1 for true, 0 for false in BASIC
	}{
		{"Equal true", "5 = 5", -1.0},
		{"Equal false", "5 = 3", 0.0},
		{"Not equal true", "5 <> 3", -1.0},
		{"Not equal false", "5 <> 5", 0.0},
		{"Less than true", "3 < 5", -1.0},
		{"Less than false", "5 < 3", 0.0},
		{"Greater than true", "5 > 3", -1.0},
		{"Greater than false", "3 > 5", 0.0},
		{"Less equal true", "3 <= 5", -1.0},
		{"Less equal equal", "5 <= 5", -1.0},
		{"Less equal false", "5 <= 3", 0.0},
		{"Greater equal true", "5 >= 3", -1.0},
		{"Greater equal equal", "5 >= 5", -1.0},
		{"Greater equal false", "3 >= 5", 0.0},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := createParser(tc.source)
			
			expr, err := parser.ParseExpression()
			require.NoError(t, err)
			require.NotNil(t, expr)
			
			// Test evaluation
			env := runtime.NewEnvironment()
			value, err := expr.Evaluate(env)
			require.NoError(t, err)
			assert.Equal(t, runtime.NumericValue, value.Type)
			assert.Equal(t, tc.expected, value.NumValue)
		})
	}
}

func TestParser_ParseExpression_StringComparisons(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		expected float64 // -1 for true, 0 for false in BASIC
	}{
		{"String equal true", `"hello" = "hello"`, -1.0},
		{"String equal false", `"hello" = "world"`, 0.0},
		{"String not equal true", `"hello" <> "world"`, -1.0},
		{"String not equal false", `"hello" <> "hello"`, 0.0},
		{"String less than true", `"apple" < "banana"`, -1.0},
		{"String less than false", `"banana" < "apple"`, 0.0},
		{"String greater than true", `"banana" > "apple"`, -1.0},
		{"String greater than false", `"apple" > "banana"`, 0.0},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := createParser(tc.source)
			
			expr, err := parser.ParseExpression()
			require.NoError(t, err)
			require.NotNil(t, expr)
			
			// Test evaluation
			env := runtime.NewEnvironment()
			value, err := expr.Evaluate(env)
			require.NoError(t, err)
			assert.Equal(t, runtime.NumericValue, value.Type)
			assert.Equal(t, tc.expected, value.NumValue)
		})
	}
}

// Test ParseExpression method - Parentheses handling

func TestParser_ParseExpression_Parentheses(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		expected float64
	}{
		{"Simple parentheses", "(5 + 3)", 8.0},
		{"Precedence override", "(2 + 3) * 4", 20.0},        // (2 + 3) * 4 = 20
		{"Nested parentheses", "((2 + 3) * 4)", 20.0},
		{"Complex nesting", "2 * (3 + (4 * 5))", 46.0},     // 2 * (3 + 20) = 46
		{"Multiple groups", "(2 + 3) * (4 + 5)", 45.0},     // 5 * 9 = 45
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := createParser(tc.source)
			
			expr, err := parser.ParseExpression()
			require.NoError(t, err)
			require.NotNil(t, expr)
			
			// Test evaluation
			env := runtime.NewEnvironment()
			value, err := expr.Evaluate(env)
			require.NoError(t, err)
			assert.Equal(t, runtime.NumericValue, value.Type)
			assert.Equal(t, tc.expected, value.NumValue)
		})
	}
}

func TestParser_ParseExpression_UnbalancedParentheses(t *testing.T) {
	testCases := []struct {
		name   string
		source string
	}{
		{"Missing closing", "(5 + 3"},
		{"Missing opening", ")"},
		{"Extra closing", ")5"},
		{"Nested unbalanced", "((5 + 3)"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := createParser(tc.source)
			
			_, err := parser.ParseExpression()
			assert.Error(t, err)
		})
	}
}

// Test ParseExpression method - Function calls

func TestParser_ParseExpression_FunctionCallsNoArgs(t *testing.T) {
	parser := createParser("RND")
	
	expr, err := parser.ParseExpression()
	require.NoError(t, err)
	require.NotNil(t, expr)
	
	funcCall, ok := expr.(*ast.FunctionCallExpression)
	require.True(t, ok, "Expected FunctionCallExpression")
	assert.Equal(t, "RND", funcCall.Name)
	assert.Len(t, funcCall.Args, 0)
	
	// Test evaluation
	env := runtime.NewEnvironment()
	value, err := funcCall.Evaluate(env)
	require.NoError(t, err)
	assert.Equal(t, runtime.NumericValue, value.Type)
	assert.True(t, value.NumValue >= 0 && value.NumValue <= 1)
}

func TestParser_ParseExpression_FunctionCallsWithArgs(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		funcName string
		argCount int
		expected float64
	}{
		{"ABS function", "ABS(-5)", "ABS", 1, 5.0},
		{"INT function", "INT(3.7)", "INT", 1, 3.0},
		{"LEN function", `LEN("hello")`, "LEN", 1, 5.0},
		{"MID$ function", `MID$("hello", 2, 3)`, "MID$", 3, 0.0}, // We'll check the string result separately
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := createParser(tc.source)
			
			expr, err := parser.ParseExpression()
			require.NoError(t, err)
			require.NotNil(t, expr)
			
			funcCall, ok := expr.(*ast.FunctionCallExpression)
			require.True(t, ok, "Expected FunctionCallExpression")
			assert.Equal(t, tc.funcName, funcCall.Name)
			assert.Len(t, funcCall.Args, tc.argCount)
			
			// Test evaluation for numeric functions
			if tc.funcName != "MID$" {
				env := runtime.NewEnvironment()
				value, err := funcCall.Evaluate(env)
				require.NoError(t, err)
				assert.Equal(t, runtime.NumericValue, value.Type)
				assert.Equal(t, tc.expected, value.NumValue)
			}
		})
	}
}

func TestParser_ParseExpression_FunctionCallsInExpressions(t *testing.T) {
	parser := createParser("ABS(-5) + INT(3.7)")
	
	expr, err := parser.ParseExpression()
	require.NoError(t, err)
	require.NotNil(t, expr)
	
	// Test evaluation
	env := runtime.NewEnvironment()
	value, err := expr.Evaluate(env)
	require.NoError(t, err)
	assert.Equal(t, runtime.NumericValue, value.Type)
	assert.Equal(t, 8.0, value.NumValue) // ABS(-5) + INT(3.7) = 5 + 3 = 8
}

func TestParser_ParseExpression_NestedFunctionCalls(t *testing.T) {
	parser := createParser("ABS(INT(-3.7))")
	
	expr, err := parser.ParseExpression()
	require.NoError(t, err)
	require.NotNil(t, expr)
	
	// Test evaluation
	env := runtime.NewEnvironment()
	value, err := expr.Evaluate(env)
	require.NoError(t, err)
	assert.Equal(t, runtime.NumericValue, value.Type)
	assert.Equal(t, 3.0, value.NumValue) // ABS(INT(-3.7)) = ABS(-3) = 3
}

// Test ParseExpression method - Variables

func TestParser_ParseExpression_Variables(t *testing.T) {
	parser := createParser("X + Y * Z")
	
	expr, err := parser.ParseExpression()
	require.NoError(t, err)
	require.NotNil(t, expr)
	
	// Test evaluation with variables
	env := runtime.NewEnvironment()
	env.SetVariable("X", runtime.NewNumericValue(10))
	env.SetVariable("Y", runtime.NewNumericValue(5))
	env.SetVariable("Z", runtime.NewNumericValue(2))
	
	value, err := expr.Evaluate(env)
	require.NoError(t, err)
	assert.Equal(t, runtime.NumericValue, value.Type)
	assert.Equal(t, 20.0, value.NumValue) // 10 + (5 * 2) = 20
}

func TestParser_ParseExpression_StringVariables(t *testing.T) {
	parser := createParser(`NAME$ + " " + TITLE$`)
	
	expr, err := parser.ParseExpression()
	require.NoError(t, err)
	require.NotNil(t, expr)
	
	// Test evaluation with string variables
	env := runtime.NewEnvironment()
	env.SetVariable("NAME$", runtime.NewStringValue("John"))
	env.SetVariable("TITLE$", runtime.NewStringValue("Doe"))
	
	value, err := expr.Evaluate(env)
	require.NoError(t, err)
	assert.Equal(t, runtime.StringValue, value.Type)
	assert.Equal(t, "John Doe", value.StrValue)
}

func TestParser_ParseExpression_MixedVariableTypes(t *testing.T) {
	parser := createParser(`STR$(X) + Y$`)
	
	expr, err := parser.ParseExpression()
	require.NoError(t, err)
	require.NotNil(t, expr)
	
	// Test evaluation with mixed types
	env := runtime.NewEnvironment()
	env.SetVariable("X", runtime.NewNumericValue(42))
	env.SetVariable("Y$", runtime.NewStringValue(" is the answer"))
	
	value, err := expr.Evaluate(env)
	require.NoError(t, err)
	assert.Equal(t, runtime.StringValue, value.Type)
	assert.Equal(t, "42 is the answer", value.StrValue)
}

// Test ParseExpression method - Complex nested expressions

func TestParser_ParseExpression_ComplexNested(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		expected float64
	}{
		{"Arithmetic with functions", "ABS(-5) * (3 + 2)", 25.0},
		{"Comparison with arithmetic", "(5 + 3) > (2 * 3)", -1.0}, // true
		{"Functions in comparisons", "ABS(-10) = INT(10.5)", -1.0}, // true: 10 = 10
		{"Complex precedence", "2 + 3 * ABS(-4) ^ 2 - 1", 49.0},   // 2 + (3 * (4 ^ 2)) - 1
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := createParser(tc.source)
			
			expr, err := parser.ParseExpression()
			require.NoError(t, err)
			require.NotNil(t, expr)
			
			// Test evaluation
			env := runtime.NewEnvironment()
			value, err := expr.Evaluate(env)
			require.NoError(t, err)
			assert.Equal(t, runtime.NumericValue, value.Type)
			assert.Equal(t, tc.expected, value.NumValue)
		})
	}
}

func TestParser_ParseExpression_DeepNesting(t *testing.T) {
	parser := createParser("((((5 + 3) * 2) - 1) / 3)")
	
	expr, err := parser.ParseExpression()
	require.NoError(t, err)
	require.NotNil(t, expr)
	
	// Test evaluation
	env := runtime.NewEnvironment()
	value, err := expr.Evaluate(env)
	require.NoError(t, err)
	assert.Equal(t, runtime.NumericValue, value.Type)
	assert.Equal(t, 5.0, value.NumValue) // ((((5 + 3) * 2) - 1) / 3) = (((16) - 1) / 3) = 15 / 3 = 5
}

// Test ParseExpression method - Error cases

func TestParser_ParseExpression_ErrorCases(t *testing.T) {
	testCases := []struct {
		name   string
		source string
		error  string
	}{
		{"Invalid token", "@", "unexpected token in expression"},
		{"Missing operand", "5 + ", "unexpected end of input"},
		{"Unbalanced parentheses", "(5 + 3", "expected )"},
		{"Empty parentheses in arithmetic", "5 + ()", "unexpected token in expression"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := createParser(tc.source)
			
			_, err := parser.ParseExpression()
			assert.Error(t, err)
			if tc.error != "" {
				assert.Contains(t, err.Error(), tc.error)
			}
		})
	}
}

func TestParser_ParseExpression_EdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		expected interface{}
	}{
		{"Single number", "42", 42.0},
		{"Single string", `"hello"`, "hello"},
		{"Single variable", "X", 0.0}, // uninitialized numeric variable
		{"Empty string", `""`, ""},
		{"Zero", "0", 0.0},
		{"Negative zero", "-0", 0.0},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := createParser(tc.source)
			
			expr, err := parser.ParseExpression()
			require.NoError(t, err)
			require.NotNil(t, expr)
			
			// Test evaluation
			env := runtime.NewEnvironment()
			value, err := expr.Evaluate(env)
			require.NoError(t, err)
			
			switch expected := tc.expected.(type) {
			case float64:
				assert.Equal(t, runtime.NumericValue, value.Type)
				assert.Equal(t, expected, value.NumValue)
			case string:
				assert.Equal(t, runtime.StringValue, value.Type)
				assert.Equal(t, expected, value.StrValue)
			}
		})
	}
}

// Test ParseProgram method - Complete program parsing

func TestParser_ParseProgram_SingleStatement(t *testing.T) {
	parser := createParser("10 PRINT \"Hello\"")
	
	program, err := parser.ParseProgram()
	require.NoError(t, err)
	require.NotNil(t, program)
	
	// Check program structure
	assert.Len(t, program.Lines, 1)
	assert.Len(t, program.Order, 1)
	assert.Equal(t, 10, program.Order[0])
	
	// Check that line 10 exists
	stmt, exists := program.Lines[10]
	require.True(t, exists)
	
	// Check that it's a PRINT statement
	printStmt, ok := stmt.(*ast.PrintStatement)
	require.True(t, ok)
	assert.Len(t, printStmt.Expressions, 1)
}

func TestParser_ParseProgram_MultipleStatements(t *testing.T) {
	source := `10 PRINT "Hello"
20 X = 42
30 PRINT X`
	
	parser := createParser(source)
	
	program, err := parser.ParseProgram()
	require.NoError(t, err)
	require.NotNil(t, program)
	
	// Check program structure
	assert.Len(t, program.Lines, 3)
	assert.Len(t, program.Order, 3)
	assert.Equal(t, []int{10, 20, 30}, program.Order)
	
	// Check each line exists
	for _, lineNum := range []int{10, 20, 30} {
		_, exists := program.Lines[lineNum]
		assert.True(t, exists, "Line %d should exist", lineNum)
	}
	
	// Check statement types
	_, ok := program.Lines[10].(*ast.PrintStatement)
	assert.True(t, ok, "Line 10 should be PRINT")
	
	_, ok = program.Lines[20].(*ast.AssignmentStatement)
	assert.True(t, ok, "Line 20 should be assignment")
	
	_, ok = program.Lines[30].(*ast.PrintStatement)
	assert.True(t, ok, "Line 30 should be PRINT")
}

func TestParser_ParseProgram_LineNumberOrdering(t *testing.T) {
	// Test with unordered line numbers
	source := `30 PRINT "Third"
10 PRINT "First"
20 PRINT "Second"`
	
	parser := createParser(source)
	
	program, err := parser.ParseProgram()
	require.NoError(t, err)
	require.NotNil(t, program)
	
	// Check that lines are ordered correctly
	assert.Equal(t, []int{10, 20, 30}, program.Order)
	
	// Check that all lines exist
	assert.Len(t, program.Lines, 3)
	for _, lineNum := range []int{10, 20, 30} {
		_, exists := program.Lines[lineNum]
		assert.True(t, exists, "Line %d should exist", lineNum)
	}
}

func TestParser_ParseProgram_LargeLineNumbers(t *testing.T) {
	source := `1000 PRINT "Large"
9999 PRINT "Larger"`
	
	parser := createParser(source)
	
	program, err := parser.ParseProgram()
	require.NoError(t, err)
	require.NotNil(t, program)
	
	// Check program structure
	assert.Len(t, program.Lines, 2)
	assert.Equal(t, []int{1000, 9999}, program.Order)
	
	// Check that lines exist
	_, exists := program.Lines[1000]
	assert.True(t, exists)
	_, exists = program.Lines[9999]
	assert.True(t, exists)
}

func TestParser_ParseProgram_ComplexProgram(t *testing.T) {
	source := `10 FOR I = 1 TO 10
20 PRINT I
30 NEXT I
40 PRINT "Done"`
	
	parser := createParser(source)
	
	program, err := parser.ParseProgram()
	require.NoError(t, err)
	require.NotNil(t, program)
	
	// Check program structure
	assert.Len(t, program.Lines, 4)
	assert.Equal(t, []int{10, 20, 30, 40}, program.Order)
	
	// Check statement types
	_, ok := program.Lines[10].(*ast.ForStatement)
	assert.True(t, ok, "Line 10 should be FOR")
	
	_, ok = program.Lines[20].(*ast.PrintStatement)
	assert.True(t, ok, "Line 20 should be PRINT")
	
	_, ok = program.Lines[30].(*ast.NextStatement)
	assert.True(t, ok, "Line 30 should be NEXT")
	
	_, ok = program.Lines[40].(*ast.PrintStatement)
	assert.True(t, ok, "Line 40 should be PRINT")
}

func TestParser_ParseProgram_WithControlFlow(t *testing.T) {
	source := `10 X = 5
20 IF X > 3 THEN GOTO 50
30 PRINT "Small"
40 GOTO 60
50 PRINT "Large"
60 PRINT "End"`
	
	parser := createParser(source)
	
	program, err := parser.ParseProgram()
	require.NoError(t, err)
	require.NotNil(t, program)
	
	// Check program structure
	assert.Len(t, program.Lines, 6)
	assert.Equal(t, []int{10, 20, 30, 40, 50, 60}, program.Order)
	
	// Check specific statement types
	_, ok := program.Lines[10].(*ast.AssignmentStatement)
	assert.True(t, ok, "Line 10 should be assignment")
	
	_, ok = program.Lines[20].(*ast.IfStatement)
	assert.True(t, ok, "Line 20 should be IF")
	
	_, ok = program.Lines[40].(*ast.GotoStatement)
	assert.True(t, ok, "Line 40 should be GOTO")
}

func TestParser_ParseProgram_EmptyProgram(t *testing.T) {
	parser := createParser("")
	
	program, err := parser.ParseProgram()
	require.NoError(t, err)
	require.NotNil(t, program)
	
	// Check empty program structure
	assert.Len(t, program.Lines, 0)
	assert.Len(t, program.Order, 0)
}

func TestParser_ParseProgram_WhitespaceAndComments(t *testing.T) {
	source := `
10 PRINT "Hello"

20 X = 42

30 PRINT X
`
	
	parser := createParser(source)
	
	program, err := parser.ParseProgram()
	require.NoError(t, err)
	require.NotNil(t, program)
	
	// Should parse correctly despite whitespace
	assert.Len(t, program.Lines, 3)
	assert.Equal(t, []int{10, 20, 30}, program.Order)
}

// Test ParseProgram method - Duplicate line detection

func TestParser_ParseProgram_DuplicateLines(t *testing.T) {
	source := `10 PRINT "First"
10 PRINT "Second"`
	
	parser := createParser(source)
	
	program, err := parser.ParseProgram()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate line number")
	assert.Contains(t, err.Error(), "10")
	assert.Nil(t, program)
}

func TestParser_ParseProgram_MultipleDuplicates(t *testing.T) {
	source := `10 PRINT "First"
20 X = 1
10 PRINT "Duplicate"
30 Y = 2
20 Z = 3`
	
	parser := createParser(source)
	
	program, err := parser.ParseProgram()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate line number")
	// Should catch the first duplicate (line 10)
	assert.Contains(t, err.Error(), "10")
	assert.Nil(t, program)
}

// Test ParseProgram method - Error cases

func TestParser_ParseProgram_InvalidSyntax(t *testing.T) {
	testCases := []struct {
		name   string
		source string
		error  string
	}{
		{"Invalid statement", "10 INVALID", "expected assignment operator"},
		{"Missing line number", "PRINT \"Hello\"", "expected line number"},
		{"Invalid line number", "ABC PRINT \"Hello\"", "expected line number"},
		{"Malformed statement", "10 PRINT", ""},  // Should parse successfully (empty PRINT)
		{"Unterminated string", "10 PRINT \"Hello", "unterminated string"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := createParser(tc.source)
			
			program, err := parser.ParseProgram()
			
			if tc.error == "" {
				// Should succeed
				assert.NoError(t, err)
				assert.NotNil(t, program)
			} else {
				// Should fail
				assert.Error(t, err)
				if tc.error != "" {
					assert.Contains(t, err.Error(), tc.error)
				}
				assert.Nil(t, program)
			}
		})
	}
}

func TestParser_ParseProgram_InvalidProgramStructure(t *testing.T) {
	testCases := []struct {
		name   string
		source string
		error  string
	}{
		{"Zero line number", "0 PRINT \"Hello\"", "invalid line number"},
		{"Negative line number", "-10 PRINT \"Hello\"", "expected line number"},
		{"Line number too large", "99999 PRINT \"Hello\"", ""}, // Should be valid
		{"Mixed valid/invalid", "10 PRINT \"OK\"\n-5 PRINT \"Bad\"", "expected line number"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := createParser(tc.source)
			
			program, err := parser.ParseProgram()
			
			if tc.error == "" {
				// Should succeed
				assert.NoError(t, err)
				assert.NotNil(t, program)
			} else {
				// Should fail
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.error)
				assert.Nil(t, program)
			}
		})
	}
}

// Test ParseProgram method - Integration with statements

func TestParser_ParseProgram_AllStatementTypes(t *testing.T) {
	source := `10 X = 42
20 PRINT "Value:", X
30 INPUT "Enter Y: "; Y
40 IF X > Y THEN GOTO 70
50 PRINT "X is not greater"
60 GOTO 80
70 PRINT "X is greater"
80 FOR I = 1 TO 3
90 PRINT I
100 NEXT I
110 PRINT "Done"`
	
	parser := createParser(source)
	
	program, err := parser.ParseProgram()
	require.NoError(t, err)
	require.NotNil(t, program)
	
	// Check program structure
	expectedLines := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110}
	assert.Len(t, program.Lines, len(expectedLines))
	assert.Equal(t, expectedLines, program.Order)
	
	// Check that all lines exist and have correct types
	expectedTypes := map[int]string{
		10:  "*ast.AssignmentStatement",
		20:  "*ast.PrintStatement",
		30:  "*ast.InputStatement",
		40:  "*ast.IfStatement",
		50:  "*ast.PrintStatement",
		60:  "*ast.GotoStatement",
		70:  "*ast.PrintStatement",
		80:  "*ast.ForStatement",
		90:  "*ast.PrintStatement",
		100: "*ast.NextStatement",
		110: "*ast.PrintStatement",
	}
	
	for lineNum, expectedType := range expectedTypes {
		stmt, exists := program.Lines[lineNum]
		require.True(t, exists, "Line %d should exist", lineNum)
		
		actualType := fmt.Sprintf("%T", stmt)
		assert.Equal(t, expectedType, actualType, "Line %d should be %s", lineNum, expectedType)
	}
}

func TestParser_ParseProgram_NestedStructures(t *testing.T) {
	source := `10 FOR I = 1 TO 3
20 FOR J = 1 TO 2
30 PRINT I, J
40 NEXT J
50 NEXT I`
	
	parser := createParser(source)
	
	program, err := parser.ParseProgram()
	require.NoError(t, err)
	require.NotNil(t, program)
	
	// Check program structure
	assert.Len(t, program.Lines, 5)
	assert.Equal(t, []int{10, 20, 30, 40, 50}, program.Order)
	
	// Check nested FOR loops
	forStmt1, ok := program.Lines[10].(*ast.ForStatement)
	require.True(t, ok)
	assert.Equal(t, "I", forStmt1.Variable)
	
	forStmt2, ok := program.Lines[20].(*ast.ForStatement)
	require.True(t, ok)
	assert.Equal(t, "J", forStmt2.Variable)
	
	nextStmt1, ok := program.Lines[40].(*ast.NextStatement)
	require.True(t, ok)
	assert.Equal(t, "J", nextStmt1.Variable)
	
	nextStmt2, ok := program.Lines[50].(*ast.NextStatement)
	require.True(t, ok)
	assert.Equal(t, "I", nextStmt2.Variable)
}

// Test ParseProgram method - Edge cases

func TestParser_ParseProgram_SingleLineProgram(t *testing.T) {
	parser := createParser("100 PRINT \"Single line\"")
	
	program, err := parser.ParseProgram()
	require.NoError(t, err)
	require.NotNil(t, program)
	
	assert.Len(t, program.Lines, 1)
	assert.Equal(t, []int{100}, program.Order)
	
	stmt, exists := program.Lines[100]
	require.True(t, exists)
	
	printStmt, ok := stmt.(*ast.PrintStatement)
	require.True(t, ok)
	assert.Len(t, printStmt.Expressions, 1)
}

func TestParser_ParseProgram_LargeProgram(t *testing.T) {
	// Create a program with many lines
	var lines []string
	for i := 10; i <= 100; i += 10 {
		lines = append(lines, fmt.Sprintf("%d PRINT %d", i, i))
	}
	source := strings.Join(lines, "\n")
	
	parser := createParser(source)
	
	program, err := parser.ParseProgram()
	require.NoError(t, err)
	require.NotNil(t, program)
	
	// Should have 10 lines (10, 20, 30, ..., 100)
	assert.Len(t, program.Lines, 10)
	assert.Len(t, program.Order, 10)
	
	// Check that all lines are in order
	expectedOrder := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	assert.Equal(t, expectedOrder, program.Order)
	
	// Check that all lines exist
	for _, lineNum := range expectedOrder {
		_, exists := program.Lines[lineNum]
		assert.True(t, exists, "Line %d should exist", lineNum)
	}
}

func TestParser_ParseProgram_NonSequentialLines(t *testing.T) {
	source := `100 PRINT "Hundred"
5 PRINT "Five"
50 PRINT "Fifty"
1 PRINT "One"`
	
	parser := createParser(source)
	
	program, err := parser.ParseProgram()
	require.NoError(t, err)
	require.NotNil(t, program)
	
	// Lines should be ordered correctly
	assert.Equal(t, []int{1, 5, 50, 100}, program.Order)
	
	// All lines should exist
	for _, lineNum := range []int{1, 5, 50, 100} {
		_, exists := program.Lines[lineNum]
		assert.True(t, exists, "Line %d should exist", lineNum)
	}
}