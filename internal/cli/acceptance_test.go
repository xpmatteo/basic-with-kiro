package cli

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AcceptanceTest represents a single acceptance test case
type AcceptanceTest struct {
	name        string
	program     string
	inputs      []string
	expected    []string
	wantErr     bool
	errContains string
}

// TestBasicInterpreter_AcceptanceTests runs comprehensive acceptance tests for the BASIC interpreter
func TestBasicInterpreter_AcceptanceTests(t *testing.T) {
	tests := []AcceptanceTest{
		// Basic PRINT statements
		{
			name: "simple print statement",
			program: `10 PRINT "Hello, World!"
20 END`,
			expected: []string{"Hello, World!"},
		},
		{
			name: "print multiple values",
			program: `10 PRINT "Value:", 42
20 END`,
			expected: []string{"Value: 42"},
		},
		{
			name: "print string variable",
			program: `10 NAME$ = "Alice"
20 PRINT NAME$
30 END`,
			expected: []string{"Alice"},
		},
		
		// Variable assignments and arithmetic
		{
			name: "numeric variable assignment",
			program: `10 A = 5
20 PRINT A
30 END`,
			expected: []string{"5"},
		},
		{
			name: "arithmetic operations",
			program: `10 A = 5
20 B = 10
30 C = A + B
40 PRINT "Sum:", C
50 END`,
			expected: []string{"Sum: 15"},
		},
		{
			name: "complex arithmetic",
			program: `10 A = 2
20 B = 3
30 C = A * B + 4
40 PRINT C
50 END`,
			expected: []string{"10"},
		},
		{
			name: "exponentiation",
			program: `10 A = 2
20 B = A ^ 3
30 PRINT B
40 END`,
			expected: []string{"8"},
		},
		
		// FOR loops
		{
			name: "simple for loop",
			program: `10 FOR I = 1 TO 3
20 PRINT "Count:", I
30 NEXT I
40 PRINT "Done!"
50 END`,
			expected: []string{"Count: 1", "Count: 2", "Count: 3", "Done!"},
		},
		{
			name: "for loop with step",
			program: `10 FOR I = 2 TO 10 STEP 2
20 PRINT "Even:", I
30 NEXT I
40 END`,
			expected: []string{"Even: 2", "Even: 4", "Even: 6", "Even: 8", "Even: 10"},
		},
		{
			name: "countdown loop",
			program: `10 FOR I = 5 TO 1 STEP -1
20 PRINT "Countdown:", I
30 NEXT I
40 PRINT "Blast off!"
50 END`,
			expected: []string{"Countdown: 5", "Countdown: 4", "Countdown: 3", "Countdown: 2", "Countdown: 1", "Blast off!"},
		},
		{
			name: "nested loops",
			program: `10 FOR I = 1 TO 2
20 FOR J = 1 TO 2
30 PRINT "I=", I, "J=", J
40 NEXT J
50 NEXT I
60 END`,
			expected: []string{"I= 1 J= 1", "I= 1 J= 2", "I= 2 J= 1", "I= 2 J= 2"},
		},
		
		// IF statements
		{
			name: "if statement true condition",
			program: `10 A = 15
20 IF A > 10 THEN PRINT "A is greater than 10"
30 PRINT "After IF"
40 END`,
			expected: []string{"A is greater than 10", "After IF"},
		},
		{
			name: "if statement false condition",
			program: `10 A = 5
20 IF A > 10 THEN PRINT "A is greater than 10"
30 PRINT "After IF"
40 END`,
			expected: []string{"After IF"},
		},
		{
			name: "if statement with string comparison",
			program: `10 NAME$ = "Alice"
20 IF NAME$ = "Alice" THEN PRINT "Hello Alice!"
30 END`,
			expected: []string{"Hello Alice!"},
		},
		
		// INPUT statements
		{
			name: "simple input",
			program: `10 INPUT X
20 PRINT "You entered:", X
30 END`,
			inputs:   []string{"42"},
			expected: []string{"You entered: 42"},
		},
		{
			name: "string input",
			program: `10 INPUT NAME$
20 PRINT "Hello", NAME$
30 END`,
			inputs:   []string{"World"},
			expected: []string{"Hello World"},
		},
		{
			name: "input with prompt",
			program: `10 INPUT "Enter your age: "; AGE
20 PRINT "You are", AGE, "years old"
30 END`,
			inputs:   []string{"25"},
			expected: []string{"Enter your age: ", "You are 25 years old"},
		},
		
		// Built-in functions
		{
			name: "abs function",
			program: `10 A = -5
20 PRINT ABS(A)
30 END`,
			expected: []string{"5"},
		},
		{
			name: "int function",
			program: `10 A = 3.7
20 PRINT INT(A)
30 END`,
			expected: []string{"3"},
		},
		{
			name: "len function",
			program: `10 S$ = "Hello"
20 PRINT LEN(S$)
30 END`,
			expected: []string{"5"},
		},
		{
			name: "mid function",
			program: `10 S$ = "Hello World"
20 PRINT MID$(S$, 7, 5)
30 END`,
			expected: []string{"World"},
		},
		{
			name: "str function",
			program: `10 A = 42
20 S$ = STR$(A)
30 PRINT S$
40 END`,
			expected: []string{"42"},
		},
		{
			name: "val function",
			program: `10 S$ = "123"
20 A = VAL(S$)
30 PRINT A
40 END`,
			expected: []string{"123"},
		},
		
		// Complex programs
		{
			name: "comprehensive program",
			program: `10 PRINT "BASIC Calculator"
20 A = 10
30 B = 5
40 PRINT "A =", A
50 PRINT "B =", B
60 PRINT "A + B =", A + B
70 PRINT "A - B =", A - B
80 PRINT "A * B =", A * B
90 PRINT "A / B =", A / B
100 END`,
			expected: []string{
				"BASIC Calculator",
				"A = 10",
				"B = 5",
				"A + B = 15",
				"A - B = 5",
				"A * B = 50",
				"A / B = 2",
			},
		},
		{
			name: "fibonacci sequence",
			program: `10 A = 0
20 B = 1
30 PRINT A
40 PRINT B
50 FOR I = 1 TO 5
60 C = A + B
70 PRINT C
80 A = B
90 B = C
100 NEXT I
110 END`,
			expected: []string{"0", "1", "1", "2", "3", "5", "8"},
		},
		{
			name: "interactive calculator",
			program: `10 INPUT "Enter first number: "; A
20 INPUT "Enter second number: "; B
30 PRINT "Sum:", A + B
40 PRINT "Product:", A * B
50 END`,
			inputs:   []string{"6", "7"},
			expected: []string{"Enter first number: ", "Enter second number: ", "Sum: 13", "Product: 42"},
		},
		
		// Comments and empty lines
		{
			name: "program with comments",
			program: `10 REM This is a comment
20 PRINT "Hello"
30 REM Another comment
40 PRINT "World"
50 END`,
			expected: []string{"Hello", "World"},
		},
		
		// Error cases
		{
			name: "syntax error - unterminated string",
			program: `10 PRINT "Hello
20 END`,
			wantErr:     true,
			errContains: "syntax error",
		},
		{
			name: "runtime error - division by zero",
			program: `10 A = 1 / 0
20 PRINT A
30 END`,
			wantErr:     true,
			errContains: "runtime error",
		},
		{
			name: "invalid line number",
			program: `ABC PRINT "Hello"
20 END`,
			wantErr:     true,
			errContains: "invalid line number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runAcceptanceTest(t, tt)
		})
	}
}

// runAcceptanceTest executes a single acceptance test
func runAcceptanceTest(t *testing.T, test AcceptanceTest) {
	// Create mock input/output
	mockInput := &MockInputReader{inputs: test.inputs}
	mockOutput := &MockOutputWriter{}
	
	// Create file executor
	fileExecutor := NewFileExecutor(mockInput, mockOutput)
	
	// Parse the program
	program, err := fileExecutor.parseProgram(test.program)
	if test.wantErr {
		if err != nil {
			// Expected error during parsing
			assert.Contains(t, err.Error(), test.errContains)
			return
		}
		// Continue to execution to catch runtime errors
	} else {
		require.NoError(t, err, "Failed to parse program")
	}
	
	// Execute the program
	err = fileExecutor.ExecuteProgram(program, false)
	
	if test.wantErr {
		assert.Error(t, err)
		if test.errContains != "" {
			assert.Contains(t, err.Error(), test.errContains)
		}
		return
	}
	
	require.NoError(t, err, "Program execution failed")
	
	// Verify expected outputs
	for _, expected := range test.expected {
		found := false
		for _, output := range mockOutput.outputs {
			if strings.Contains(output, expected) {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected output '%s' not found in %v", expected, mockOutput.outputs)
	}
}

// TestBasicInterpreter_DebugMode tests debug mode functionality
func TestBasicInterpreter_DebugMode(t *testing.T) {
	program := `10 PRINT "Line 1"
20 A = 5
30 PRINT "A =", A
40 END`
	
	mockInput := &MockInputReader{}
	mockOutput := &MockOutputWriter{}
	fileExecutor := NewFileExecutor(mockInput, mockOutput)
	
	parsedProgram, err := fileExecutor.parseProgram(program)
	require.NoError(t, err)
	
	err = fileExecutor.ExecuteProgram(parsedProgram, true) // Debug mode enabled
	require.NoError(t, err)
	
	// Check that debug output is present
	debugFound := false
	for _, output := range mockOutput.outputs {
		if strings.Contains(output, "Debug:") || strings.Contains(output, "Executing") {
			debugFound = true
			break
		}
	}
	assert.True(t, debugFound, "Debug output should be present when debug mode is enabled")
	
	// Check that program output is still present
	outputFound := false
	for _, output := range mockOutput.outputs {
		if strings.Contains(output, "Line 1") {
			outputFound = true
			break
		}
	}
	assert.True(t, outputFound, "Program output should be present in debug mode")
}

// TestBasicInterpreter_EmptyProgram tests handling of empty programs
func TestBasicInterpreter_EmptyProgram(t *testing.T) {
	mockInput := &MockInputReader{}
	mockOutput := &MockOutputWriter{}
	fileExecutor := NewFileExecutor(mockInput, mockOutput)
	
	// Test completely empty program
	program, err := fileExecutor.parseProgram("")
	require.NoError(t, err)
	
	err = fileExecutor.ExecuteProgram(program, false)
	assert.NoError(t, err, "Empty program should execute without error")
	
	// Test program with only comments
	commentProgram := `10 REM This is a comment
20 REM Another comment`
	
	program, err = fileExecutor.parseProgram(commentProgram)
	require.NoError(t, err)
	
	err = fileExecutor.ExecuteProgram(program, false)
	assert.NoError(t, err, "Comment-only program should execute without error")
}

// TestBasicInterpreter_LineNumberOrdering tests that line numbers are executed in correct order
func TestBasicInterpreter_LineNumberOrdering(t *testing.T) {
	// Program with line numbers out of order
	program := `30 PRINT "Third"
10 PRINT "First"
20 PRINT "Second"
40 END`
	
	mockInput := &MockInputReader{}
	mockOutput := &MockOutputWriter{}
	fileExecutor := NewFileExecutor(mockInput, mockOutput)
	
	parsedProgram, err := fileExecutor.parseProgram(program)
	require.NoError(t, err)
	
	err = fileExecutor.ExecuteProgram(parsedProgram, false)
	require.NoError(t, err)
	
	// Check that outputs appear in correct order
	expectedOrder := []string{"First", "Second", "Third"}
	outputIndex := 0
	
	for _, output := range mockOutput.outputs {
		for _, expected := range expectedOrder {
			if strings.Contains(output, expected) {
				assert.Equal(t, expectedOrder[outputIndex], expected, 
					"Output should appear in line number order, expected %s at position %d", 
					expectedOrder[outputIndex], outputIndex)
				outputIndex++
				break
			}
		}
	}
	
	assert.Equal(t, len(expectedOrder), outputIndex, "All expected outputs should be found in correct order")
}

// TestBasicInterpreter_VariableScoping tests variable scoping and persistence
func TestBasicInterpreter_VariableScoping(t *testing.T) {
	program := `10 A = 5
20 FOR I = 1 TO 3
30 A = A + I
40 PRINT "A =", A
50 NEXT I
60 PRINT "Final A =", A
70 END`
	
	mockInput := &MockInputReader{}
	mockOutput := &MockOutputWriter{}
	fileExecutor := NewFileExecutor(mockInput, mockOutput)
	
	parsedProgram, err := fileExecutor.parseProgram(program)
	require.NoError(t, err)
	
	err = fileExecutor.ExecuteProgram(parsedProgram, false)
	require.NoError(t, err)
	
	// Check that variable A is modified correctly through the loop
	expectedValues := []string{"A = 6", "A = 8", "A = 11", "Final A = 11"}
	
	for _, expected := range expectedValues {
		found := false
		for _, output := range mockOutput.outputs {
			if strings.Contains(output, expected) {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected output '%s' not found in %v", expected, mockOutput.outputs)
	}
}