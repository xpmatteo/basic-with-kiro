package interpreter

import (
	"basic-interpreter/internal/ast"
	"basic-interpreter/internal/runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockOutputWriter for testing output operations
type MockOutputWriter struct {
	Lines []string
}

func (m *MockOutputWriter) WriteLine(line string) error {
	m.Lines = append(m.Lines, line)
	return nil
}

// MockInputReader for testing input operations
type MockInputReader struct {
	Inputs []string
	Index  int
}

func (m *MockInputReader) ReadLine() (string, error) {
	if m.Index >= len(m.Inputs) {
		return "", nil
	}
	result := m.Inputs[m.Index]
	m.Index++
	return result, nil
}

// Test program execution with basic statements
func TestInterpreter_Execute_BasicProgram(t *testing.T) {
	// Create a simple program with assignment and print
	program := &ast.Program{
		Lines: map[int]ast.Statement{
			10: ast.NewAssignmentStatement("X", ast.NewLiteralExpression(runtime.NewNumericValue(42))),
			20: ast.NewPrintStatement([]ast.Expression{ast.NewVariableExpression("X")}, &MockOutputWriter{}),
		},
		Order: []int{10, 20},
	}

	env := runtime.NewEnvironment()
	interpreter := NewBasicInterpreter(false)

	err := interpreter.Execute(program, env)
	assert.NoError(t, err)

	// Verify variable was set
	value := env.GetVariable("X")
	assert.Equal(t, runtime.NumericValue, value.Type)
	assert.Equal(t, 42.0, value.NumValue)
}

// Test program counter management during sequential execution
func TestInterpreter_Execute_ProgramCounterSequential(t *testing.T) {
	output := &MockOutputWriter{}
	program := &ast.Program{
		Lines: map[int]ast.Statement{
			10: ast.NewPrintStatement([]ast.Expression{ast.NewLiteralExpression(runtime.NewStringValue("Line 10"))}, output),
			20: ast.NewPrintStatement([]ast.Expression{ast.NewLiteralExpression(runtime.NewStringValue("Line 20"))}, output),
			30: ast.NewPrintStatement([]ast.Expression{ast.NewLiteralExpression(runtime.NewStringValue("Line 30"))}, output),
		},
		Order: []int{10, 20, 30},
	}

	env := runtime.NewEnvironment()
	interpreter := NewBasicInterpreter(false)

	err := interpreter.Execute(program, env)
	assert.NoError(t, err)

	// Verify all lines were executed in order
	assert.Equal(t, []string{"Line 10", "Line 20", "Line 30"}, output.Lines)
}

// Test program counter management with GOTO statement
func TestInterpreter_Execute_ProgramCounterWithGoto(t *testing.T) {
	output := &MockOutputWriter{}
	program := &ast.Program{
		Lines: map[int]ast.Statement{
			10: ast.NewPrintStatement([]ast.Expression{ast.NewLiteralExpression(runtime.NewStringValue("Line 10"))}, output),
			20: ast.NewGotoStatement(40, nil), // Will be set after program creation
			30: ast.NewPrintStatement([]ast.Expression{ast.NewLiteralExpression(runtime.NewStringValue("Line 30"))}, output),
			40: ast.NewPrintStatement([]ast.Expression{ast.NewLiteralExpression(runtime.NewStringValue("Line 40"))}, output),
		},
		Order: []int{10, 20, 30, 40},
	}

	// Set program reference for GOTO statement
	program.Lines[20].(*ast.GotoStatement).Program = program

	env := runtime.NewEnvironment()
	interpreter := NewBasicInterpreter(false)

	err := interpreter.Execute(program, env)
	assert.NoError(t, err)

	// Verify line 30 was skipped due to GOTO
	assert.Equal(t, []string{"Line 10", "Line 40"}, output.Lines)
}

// Test execution state management with variables
func TestInterpreter_Execute_VariablePersistence(t *testing.T) {
	program := &ast.Program{
		Lines: map[int]ast.Statement{
			10: ast.NewAssignmentStatement("A", ast.NewLiteralExpression(runtime.NewNumericValue(10))),
			20: ast.NewAssignmentStatement("B", ast.NewVariableExpression("A")),
			30: ast.NewAssignmentStatement("A", ast.NewBinaryExpression(
				ast.NewVariableExpression("A"),
				"+",
				ast.NewLiteralExpression(runtime.NewNumericValue(5)),
			)),
		},
		Order: []int{10, 20, 30},
	}

	env := runtime.NewEnvironment()
	interpreter := NewBasicInterpreter(false)

	err := interpreter.Execute(program, env)
	assert.NoError(t, err)

	// Verify variables were updated correctly
	valueA := env.GetVariable("A")
	valueB := env.GetVariable("B")
	assert.Equal(t, 15.0, valueA.NumValue)
	assert.Equal(t, 10.0, valueB.NumValue)
}

// Test control flow with IF-THEN statement
func TestInterpreter_Execute_ControlFlowIfThen(t *testing.T) {
	output := &MockOutputWriter{}
	program := &ast.Program{
		Lines: map[int]ast.Statement{
			10: ast.NewAssignmentStatement("X", ast.NewLiteralExpression(runtime.NewNumericValue(5))),
			20: ast.NewIfStatement(
				ast.NewComparisonExpression(
					ast.NewVariableExpression("X"),
					">",
					ast.NewLiteralExpression(runtime.NewNumericValue(3)),
				),
				ast.NewPrintStatement([]ast.Expression{ast.NewLiteralExpression(runtime.NewStringValue("X is greater than 3"))}, output),
			),
			30: ast.NewPrintStatement([]ast.Expression{ast.NewLiteralExpression(runtime.NewStringValue("End"))}, output),
		},
		Order: []int{10, 20, 30},
	}

	env := runtime.NewEnvironment()
	interpreter := NewBasicInterpreter(false)

	err := interpreter.Execute(program, env)
	assert.NoError(t, err)

	// Verify conditional statement was executed
	assert.Equal(t, []string{"X is greater than 3", "End"}, output.Lines)
}

// Test FOR-NEXT loop execution
func TestInterpreter_Execute_ForNextLoop(t *testing.T) {
	output := &MockOutputWriter{}
	program := &ast.Program{
		Lines: map[int]ast.Statement{
			10: ast.NewForStatement("I", 
				ast.NewLiteralExpression(runtime.NewNumericValue(1)),
				ast.NewLiteralExpression(runtime.NewNumericValue(3)),
				ast.NewLiteralExpression(runtime.NewNumericValue(1)),
				10,
			),
			20: ast.NewPrintStatement([]ast.Expression{ast.NewVariableExpression("I")}, output),
			30: ast.NewNextStatement("I"),
		},
		Order: []int{10, 20, 30},
	}

	env := runtime.NewEnvironment()
	interpreter := NewBasicInterpreter(false)

	err := interpreter.Execute(program, env)
	assert.NoError(t, err)

	// Verify loop executed correctly
	assert.Equal(t, []string{"1", "2", "3"}, output.Lines)
}

// Test empty program execution
func TestInterpreter_Execute_EmptyProgram(t *testing.T) {
	program := &ast.Program{
		Lines: map[int]ast.Statement{},
		Order: []int{},
	}

	env := runtime.NewEnvironment()
	interpreter := NewBasicInterpreter(false)

	err := interpreter.Execute(program, env)
	assert.NoError(t, err)
}

// Test program with single statement
func TestInterpreter_Execute_SingleStatement(t *testing.T) {
	output := &MockOutputWriter{}
	program := &ast.Program{
		Lines: map[int]ast.Statement{
			10: ast.NewPrintStatement([]ast.Expression{ast.NewLiteralExpression(runtime.NewStringValue("Hello"))}, output),
		},
		Order: []int{10},
	}

	env := runtime.NewEnvironment()
	interpreter := NewBasicInterpreter(false)

	err := interpreter.Execute(program, env)
	assert.NoError(t, err)

	assert.Equal(t, []string{"Hello"}, output.Lines)
}

// Test program execution with unordered line numbers
func TestInterpreter_Execute_UnorderedLines(t *testing.T) {
	output := &MockOutputWriter{}
	program := &ast.Program{
		Lines: map[int]ast.Statement{
			30: ast.NewPrintStatement([]ast.Expression{ast.NewLiteralExpression(runtime.NewStringValue("Third"))}, output),
			10: ast.NewPrintStatement([]ast.Expression{ast.NewLiteralExpression(runtime.NewStringValue("First"))}, output),
			20: ast.NewPrintStatement([]ast.Expression{ast.NewLiteralExpression(runtime.NewStringValue("Second"))}, output),
		},
		Order: []int{10, 20, 30}, // Order should be respected
	}

	env := runtime.NewEnvironment()
	interpreter := NewBasicInterpreter(false)

	err := interpreter.Execute(program, env)
	assert.NoError(t, err)

	// Verify execution followed the order, not the line number sequence
	assert.Equal(t, []string{"First", "Second", "Third"}, output.Lines)
}

// Test execution state management with nested FOR loops
func TestInterpreter_Execute_NestedForLoops(t *testing.T) {
	output := &MockOutputWriter{}
	program := &ast.Program{
		Lines: map[int]ast.Statement{
			10: ast.NewForStatement("I", 
				ast.NewLiteralExpression(runtime.NewNumericValue(1)),
				ast.NewLiteralExpression(runtime.NewNumericValue(2)),
				ast.NewLiteralExpression(runtime.NewNumericValue(1)),
				10,
			),
			20: ast.NewForStatement("J", 
				ast.NewLiteralExpression(runtime.NewNumericValue(1)),
				ast.NewLiteralExpression(runtime.NewNumericValue(2)),
				ast.NewLiteralExpression(runtime.NewNumericValue(1)),
				20,
			),
			30: ast.NewPrintStatement([]ast.Expression{
				ast.NewVariableExpression("I"),
				ast.NewLiteralExpression(runtime.NewStringValue(",")),
				ast.NewVariableExpression("J"),
			}, output),
			40: ast.NewNextStatement("J"),
			50: ast.NewNextStatement("I"),
		},
		Order: []int{10, 20, 30, 40, 50},
	}

	env := runtime.NewEnvironment()
	interpreter := NewBasicInterpreter(false)

	err := interpreter.Execute(program, env)
	assert.NoError(t, err)

	// Verify nested loops executed correctly
	expected := []string{"1 , 1", "1 , 2", "2 , 1", "2 , 2"}
	assert.Equal(t, expected, output.Lines)
}

// Error handling and debugging tests

// Test runtime error reporting with line numbers
func TestInterpreter_Execute_RuntimeErrorWithLineNumber(t *testing.T) {
	// Create a program that will cause a runtime error (division by zero)
	program := &ast.Program{
		Lines: map[int]ast.Statement{
			10: ast.NewAssignmentStatement("X", ast.NewLiteralExpression(runtime.NewNumericValue(10))),
			20: ast.NewAssignmentStatement("Y", ast.NewLiteralExpression(runtime.NewNumericValue(0))),
			30: ast.NewAssignmentStatement("Z", ast.NewBinaryExpression(
				ast.NewVariableExpression("X"),
				"/",
				ast.NewVariableExpression("Y"),
			)),
		},
		Order: []int{10, 20, 30},
	}

	env := runtime.NewEnvironment()
	interpreter := NewBasicInterpreter(false)

	err := interpreter.Execute(program, env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "line 30") // Should include line number in error
}

// Test debug mode showing each line before execution
func TestInterpreter_Execute_DebugMode(t *testing.T) {
	output := &MockOutputWriter{}
	debugOutput := &MockOutputWriter{}
	
	program := &ast.Program{
		Lines: map[int]ast.Statement{
			10: ast.NewAssignmentStatement("X", ast.NewLiteralExpression(runtime.NewNumericValue(5))),
			20: ast.NewPrintStatement([]ast.Expression{ast.NewVariableExpression("X")}, output),
		},
		Order: []int{10, 20},
	}

	env := runtime.NewEnvironment()
	interpreter := NewInterpreter(InterpreterConfig{DebugMode: true, DebugOutput: debugOutput})

	err := interpreter.Execute(program, env)
	assert.NoError(t, err)

	// Verify debug output shows line numbers before execution
	assert.Len(t, debugOutput.Lines, 2)
	assert.Contains(t, debugOutput.Lines[0], "Executing line 10")
	assert.Contains(t, debugOutput.Lines[1], "Executing line 20")
	
	// Verify normal output still works
	assert.Equal(t, []string{"5"}, output.Lines)
}

// Test interrupt handling for infinite loops
func TestInterpreter_Execute_InfiniteLoopProtection(t *testing.T) {
	program := &ast.Program{
		Lines: make(map[int]ast.Statement),
		Order: []int{10},
	}
	
	// Create GOTO statement with program reference
	program.Lines[10] = ast.NewGotoStatement(10, program) // Infinite loop

	env := runtime.NewEnvironment()
	interpreter := NewInterpreter(InterpreterConfig{DebugMode: false, MaxSteps: 1000}) // Limit execution steps

	err := interpreter.Execute(program, env)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "execution limit exceeded") // Should detect infinite loop
}

// Test error message formatting
func TestInterpreter_Execute_ErrorMessageFormatting(t *testing.T) {
	// Create a program with an invalid variable name
	program := &ast.Program{
		Lines: map[int]ast.Statement{
			10: ast.NewAssignmentStatement("", ast.NewLiteralExpression(runtime.NewNumericValue(5))), // Invalid empty variable name
		},
		Order: []int{10},
	}

	env := runtime.NewEnvironment()
	interpreter := NewBasicInterpreter(false)

	err := interpreter.Execute(program, env)
	assert.Error(t, err)
	
	// Error should be well-formatted with context
	assert.Contains(t, err.Error(), "line 10")
	assert.Contains(t, err.Error(), "variable name")
}

// Test debug output formatting
func TestInterpreter_Execute_DebugOutputFormatting(t *testing.T) {
	debugOutput := &MockOutputWriter{}
	program := &ast.Program{
		Lines: map[int]ast.Statement{
			10: ast.NewAssignmentStatement("X", ast.NewLiteralExpression(runtime.NewNumericValue(42))),
			30: ast.NewAssignmentStatement("Y", ast.NewLiteralExpression(runtime.NewStringValue("Hello"))),
		},
		Order: []int{10, 30},
	}

	env := runtime.NewEnvironment()
	interpreter := NewInterpreter(InterpreterConfig{DebugMode: true, DebugOutput: debugOutput})

	err := interpreter.Execute(program, env)
	assert.NoError(t, err)

	// Verify debug output is properly formatted
	assert.Len(t, debugOutput.Lines, 2)
	assert.Equal(t, "Executing line 10: X = 42", debugOutput.Lines[0])
	assert.Equal(t, "Executing line 30: Y = \"Hello\"", debugOutput.Lines[1])
}

// Test nested error handling (error in nested control structure)
func TestInterpreter_Execute_NestedErrorHandling(t *testing.T) {
	program := &ast.Program{
		Lines: map[int]ast.Statement{
			10: ast.NewForStatement("I", 
				ast.NewLiteralExpression(runtime.NewNumericValue(1)),
				ast.NewLiteralExpression(runtime.NewNumericValue(2)),
				ast.NewLiteralExpression(runtime.NewNumericValue(1)),
				10,
			),
			20: ast.NewAssignmentStatement("X", ast.NewBinaryExpression(
				ast.NewLiteralExpression(runtime.NewNumericValue(10)),
				"/",
				ast.NewLiteralExpression(runtime.NewNumericValue(0)), // Division by zero
			)),
			30: ast.NewNextStatement("I"),
		},
		Order: []int{10, 20, 30},
	}

	env := runtime.NewEnvironment()
	interpreter := NewBasicInterpreter(false)

	err := interpreter.Execute(program, env)
	assert.Error(t, err)
	
	// Should report error with line number and context
	assert.Contains(t, err.Error(), "line 20")
	assert.Contains(t, err.Error(), "division by zero")
}

// Test execution step counting for performance monitoring
func TestInterpreter_Execute_StepCounting(t *testing.T) {
	program := &ast.Program{
		Lines: map[int]ast.Statement{
			10: ast.NewAssignmentStatement("X", ast.NewLiteralExpression(runtime.NewNumericValue(1))),
			20: ast.NewAssignmentStatement("Y", ast.NewLiteralExpression(runtime.NewNumericValue(2))),
			30: ast.NewAssignmentStatement("Z", ast.NewLiteralExpression(runtime.NewNumericValue(3))),
		},
		Order: []int{10, 20, 30},
	}

	env := runtime.NewEnvironment()
	interpreter := NewBasicInterpreter(false)

	err := interpreter.Execute(program, env)
	assert.NoError(t, err)

	// Should have executed exactly 3 steps
	assert.Equal(t, 3, interpreter.GetStepCount())
}