package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// assertTokensEqual is a helper function to compare expected and actual tokens
func assertTokensEqual(t *testing.T, expected, actual Token, index int) {
	assert.Equal(t, expected.Type, actual.Type, "token[%d].Type", index)
	assert.Equal(t, expected.Value, actual.Value, "token[%d].Value", index)
	assert.Equal(t, expected.Line, actual.Line, "token[%d].Line", index)
	assert.Equal(t, expected.Column, actual.Column, "token[%d].Column", index)
}

// TestTokenType_String tests that TokenType has proper string representation
func TestTokenType_String(t *testing.T) {
	tests := []struct {
		tokenType TokenType
		expected  string
	}{
		{ILLEGAL, "ILLEGAL"},
		{EOF, "EOF"},
		{NUMBER, "NUMBER"},
		{STRING, "STRING"},
		{IDENTIFIER, "IDENTIFIER"},
		{PRINT, "PRINT"},
		{INPUT, "INPUT"},
		{LET, "LET"},
		{IF, "IF"},
		{THEN, "THEN"},
		{GOTO, "GOTO"},
		{FOR, "FOR"},
		{TO, "TO"},
		{NEXT, "NEXT"},
		{STEP, "STEP"},
		{END, "END"},
		{ASSIGN, "ASSIGN"},
		{PLUS, "PLUS"},
		{MINUS, "MINUS"},
		{MULTIPLY, "MULTIPLY"},
		{DIVIDE, "DIVIDE"},
		{POWER, "POWER"},
		{EQ, "EQ"},
		{LT, "LT"},
		{GT, "GT"},
		{LE, "LE"},
		{GE, "GE"},
		{NE, "NE"},
		{SEMICOLON, "SEMICOLON"},
		{COMMA, "COMMA"},
		{LPAREN, "LPAREN"},
		{RPAREN, "RPAREN"},
		{LINENUMBER, "LINENUMBER"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.tokenType.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestToken_Creation tests basic token creation
func TestToken_Creation(t *testing.T) {
	token := Token{
		Type:   NUMBER,
		Value:  "123",
		Line:   1,
		Column: 1,
	}

	assert.Equal(t, NUMBER, token.Type)
	assert.Equal(t, "123", token.Value)
	assert.Equal(t, 1, token.Line)
	assert.Equal(t, 1, token.Column)
}

// TestLexer_Numbers tests tokenizing various number formats
func TestLexer_Numbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "simple integer",
			input: "123",
			expected: []Token{
				{Type: NUMBER, Value: "123", Line: 1, Column: 1},
				{Type: EOF, Value: "", Line: 1, Column: 4},
			},
		},
		{
			name:  "decimal number",
			input: "123.45",
			expected: []Token{
				{Type: NUMBER, Value: "123.45", Line: 1, Column: 1},
				{Type: EOF, Value: "", Line: 1, Column: 7},
			},
		},
		{
			name:  "number starting with decimal",
			input: ".5",
			expected: []Token{
				{Type: NUMBER, Value: ".5", Line: 1, Column: 1},
				{Type: EOF, Value: "", Line: 1, Column: 3},
			},
		},
		{
			name:  "zero",
			input: "0",
			expected: []Token{
				{Type: NUMBER, Value: "0", Line: 1, Column: 1},
				{Type: EOF, Value: "", Line: 1, Column: 2},
			},
		},
		{
			name:  "multiple numbers with spaces",
			input: "10 20 30",
			expected: []Token{
				{Type: NUMBER, Value: "10", Line: 1, Column: 1},
				{Type: NUMBER, Value: "20", Line: 1, Column: 4},
				{Type: NUMBER, Value: "30", Line: 1, Column: 7},
				{Type: EOF, Value: "", Line: 1, Column: 9},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			
			for i, expectedToken := range tt.expected {
				token := lexer.NextToken()
				assertTokensEqual(t, expectedToken, token, i)
			}
		})
	}
}

// TestLexer_Strings tests tokenizing string literals
func TestLexer_Strings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "simple string",
			input: `"Hello"`,
			expected: []Token{
				{Type: STRING, Value: "Hello", Line: 1, Column: 1},
				{Type: EOF, Value: "", Line: 1, Column: 8},
			},
		},
		{
			name:  "empty string",
			input: `""`,
			expected: []Token{
				{Type: STRING, Value: "", Line: 1, Column: 1},
				{Type: EOF, Value: "", Line: 1, Column: 3},
			},
		},
		{
			name:  "string with spaces",
			input: `"Hello World"`,
			expected: []Token{
				{Type: STRING, Value: "Hello World", Line: 1, Column: 1},
				{Type: EOF, Value: "", Line: 1, Column: 14},
			},
		},
		{
			name:  "string with numbers",
			input: `"Test 123"`,
			expected: []Token{
				{Type: STRING, Value: "Test 123", Line: 1, Column: 1},
				{Type: EOF, Value: "", Line: 1, Column: 11},
			},
		},
		{
			name:  "multiple strings",
			input: `"First" "Second"`,
			expected: []Token{
				{Type: STRING, Value: "First", Line: 1, Column: 1},
				{Type: STRING, Value: "Second", Line: 1, Column: 9},
				{Type: EOF, Value: "", Line: 1, Column: 17},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			
			for i, expectedToken := range tt.expected {
				token := lexer.NextToken()
				assertTokensEqual(t, expectedToken, token, i)
			}
		})
	}
}

// TestLexer_Identifiers tests tokenizing identifiers and variable names
func TestLexer_Identifiers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "simple identifier",
			input: "X",
			expected: []Token{
				{Type: IDENTIFIER, Value: "X", Line: 1, Column: 1},
				{Type: EOF, Value: "", Line: 1, Column: 2},
			},
		},
		{
			name:  "multi-character identifier",
			input: "VAR1",
			expected: []Token{
				{Type: IDENTIFIER, Value: "VAR1", Line: 1, Column: 1},
				{Type: EOF, Value: "", Line: 1, Column: 5},
			},
		},
		{
			name:  "string variable",
			input: "NAME$",
			expected: []Token{
				{Type: IDENTIFIER, Value: "NAME$", Line: 1, Column: 1},
				{Type: EOF, Value: "", Line: 1, Column: 6},
			},
		},
		{
			name:  "mixed case identifier",
			input: "MyVar",
			expected: []Token{
				{Type: IDENTIFIER, Value: "MyVar", Line: 1, Column: 1},
				{Type: EOF, Value: "", Line: 1, Column: 6},
			},
		},
		{
			name:  "multiple identifiers",
			input: "A B C",
			expected: []Token{
				{Type: IDENTIFIER, Value: "A", Line: 1, Column: 1},
				{Type: IDENTIFIER, Value: "B", Line: 1, Column: 3},
				{Type: IDENTIFIER, Value: "C", Line: 1, Column: 5},
				{Type: EOF, Value: "", Line: 1, Column: 6},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			
			for i, expectedToken := range tt.expected {
				token := lexer.NextToken()
				assertTokensEqual(t, expectedToken, token, i)
			}
		})
	}
}

// TestLexer_Keywords tests tokenizing BASIC keywords
func TestLexer_Keywords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "PRINT keyword",
			input: "PRINT",
			expected: []Token{
				{Type: PRINT, Value: "PRINT", Line: 1, Column: 1},
				{Type: EOF, Value: "", Line: 1, Column: 6},
			},
		},
		{
			name:  "INPUT keyword",
			input: "INPUT",
			expected: []Token{
				{Type: INPUT, Value: "INPUT", Line: 1, Column: 1},
				{Type: EOF, Value: "", Line: 1, Column: 6},
			},
		},
		{
			name:  "LET keyword",
			input: "LET",
			expected: []Token{
				{Type: LET, Value: "LET", Line: 1, Column: 1},
				{Type: EOF, Value: "", Line: 1, Column: 4},
			},
		},
		{
			name:  "IF THEN keywords",
			input: "IF THEN",
			expected: []Token{
				{Type: IF, Value: "IF", Line: 1, Column: 1},
				{Type: THEN, Value: "THEN", Line: 1, Column: 4},
				{Type: EOF, Value: "", Line: 1, Column: 8},
			},
		},
		{
			name:  "FOR TO NEXT STEP keywords",
			input: "FOR TO NEXT STEP",
			expected: []Token{
				{Type: FOR, Value: "FOR", Line: 1, Column: 1},
				{Type: TO, Value: "TO", Line: 1, Column: 5},
				{Type: NEXT, Value: "NEXT", Line: 1, Column: 8},
				{Type: STEP, Value: "STEP", Line: 1, Column: 13},
				{Type: EOF, Value: "", Line: 1, Column: 17},
			},
		},
		{
			name:  "GOTO and END keywords",
			input: "GOTO END",
			expected: []Token{
				{Type: GOTO, Value: "GOTO", Line: 1, Column: 1},
				{Type: END, Value: "END", Line: 1, Column: 6},
				{Type: EOF, Value: "", Line: 1, Column: 9},
			},
		},
		{
			name:  "lowercase keywords should be recognized",
			input: "print input",
			expected: []Token{
				{Type: PRINT, Value: "print", Line: 1, Column: 1},
				{Type: INPUT, Value: "input", Line: 1, Column: 7},
				{Type: EOF, Value: "", Line: 1, Column: 12},
			},
		},
		{
			name:  "mixed case keywords should be recognized",
			input: "Print Input",
			expected: []Token{
				{Type: PRINT, Value: "Print", Line: 1, Column: 1},
				{Type: INPUT, Value: "Input", Line: 1, Column: 7},
				{Type: EOF, Value: "", Line: 1, Column: 12},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			
			for i, expectedToken := range tt.expected {
				token := lexer.NextToken()
				assertTokensEqual(t, expectedToken, token, i)
			}
		})
	}
}

// TestLexer_Operators tests tokenizing operators
func TestLexer_Operators(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "arithmetic operators",
			input: "+ - * / ^",
			expected: []Token{
				{Type: PLUS, Value: "+", Line: 1, Column: 1},
				{Type: MINUS, Value: "-", Line: 1, Column: 3},
				{Type: MULTIPLY, Value: "*", Line: 1, Column: 5},
				{Type: DIVIDE, Value: "/", Line: 1, Column: 7},
				{Type: POWER, Value: "^", Line: 1, Column: 9},
				{Type: EOF, Value: "", Line: 1, Column: 10},
			},
		},
		{
			name:  "assignment operator",
			input: "=",
			expected: []Token{
				{Type: ASSIGN, Value: "=", Line: 1, Column: 1},
				{Type: EOF, Value: "", Line: 1, Column: 2},
			},
		},
		{
			name:  "comparison operators",
			input: "< > <= >= <>",
			expected: []Token{
				{Type: LT, Value: "<", Line: 1, Column: 1},
				{Type: GT, Value: ">", Line: 1, Column: 3},
				{Type: LE, Value: "<=", Line: 1, Column: 5},
				{Type: GE, Value: ">=", Line: 1, Column: 8},
				{Type: NE, Value: "<>", Line: 1, Column: 11},
				{Type: EOF, Value: "", Line: 1, Column: 13},
			},
		},
		{
			name:  "delimiters",
			input: "( ) , ;",
			expected: []Token{
				{Type: LPAREN, Value: "(", Line: 1, Column: 1},
				{Type: RPAREN, Value: ")", Line: 1, Column: 3},
				{Type: COMMA, Value: ",", Line: 1, Column: 5},
				{Type: SEMICOLON, Value: ";", Line: 1, Column: 7},
				{Type: EOF, Value: "", Line: 1, Column: 8},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			
			for i, expectedToken := range tt.expected {
				token := lexer.NextToken()
				assertTokensEqual(t, expectedToken, token, i)
			}
		})
	}
}

// TestLexer_ErrorCases tests error handling for invalid input
func TestLexer_ErrorCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "unterminated string",
			input: `"Hello`,
			expected: []Token{
				{Type: ILLEGAL, Value: "unterminated string", Line: 1, Column: 1},
				{Type: EOF, Value: "", Line: 1, Column: 7},
			},
		},
		{
			name:  "invalid character",
			input: "@",
			expected: []Token{
				{Type: ILLEGAL, Value: "@", Line: 1, Column: 1},
				{Type: EOF, Value: "", Line: 1, Column: 2},
			},
		},
		{
			name:  "invalid character in middle",
			input: "A @ B",
			expected: []Token{
				{Type: IDENTIFIER, Value: "A", Line: 1, Column: 1},
				{Type: ILLEGAL, Value: "@", Line: 1, Column: 3},
				{Type: IDENTIFIER, Value: "B", Line: 1, Column: 5},
				{Type: EOF, Value: "", Line: 1, Column: 6},
			},
		},
		{
			name:  "multiple invalid characters",
			input: "# $ %",
			expected: []Token{
				{Type: ILLEGAL, Value: "#", Line: 1, Column: 1},
				{Type: ILLEGAL, Value: "$", Line: 1, Column: 3},
				{Type: ILLEGAL, Value: "%", Line: 1, Column: 5},
				{Type: EOF, Value: "", Line: 1, Column: 6},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			
			for i, expectedToken := range tt.expected {
				token := lexer.NextToken()
				assertTokensEqual(t, expectedToken, token, i)
			}
		})
	}
}

// TestLexer_PositionTracking tests line and column number tracking
func TestLexer_PositionTracking(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "single line with multiple tokens",
			input: "10 PRINT X",
			expected: []Token{
				{Type: LINENUMBER, Value: "10", Line: 1, Column: 1},
				{Type: PRINT, Value: "PRINT", Line: 1, Column: 4},
				{Type: IDENTIFIER, Value: "X", Line: 1, Column: 10},
				{Type: EOF, Value: "", Line: 1, Column: 11},
			},
		},
		{
			name:  "multiple lines",
			input: "10 PRINT\n20 INPUT",
			expected: []Token{
				{Type: LINENUMBER, Value: "10", Line: 1, Column: 1},
				{Type: PRINT, Value: "PRINT", Line: 1, Column: 4},
				{Type: LINENUMBER, Value: "20", Line: 2, Column: 1},
				{Type: INPUT, Value: "INPUT", Line: 2, Column: 4},
				{Type: EOF, Value: "", Line: 2, Column: 9},
			},
		},
		{
			name:  "line with string spanning multiple columns",
			input: `10 PRINT "Hello World"`,
			expected: []Token{
				{Type: LINENUMBER, Value: "10", Line: 1, Column: 1},
				{Type: PRINT, Value: "PRINT", Line: 1, Column: 4},
				{Type: STRING, Value: "Hello World", Line: 1, Column: 10},
				{Type: EOF, Value: "", Line: 1, Column: 23},
			},
		},
		{
			name:  "empty lines should increment line counter",
			input: "10 PRINT\n\n20 INPUT",
			expected: []Token{
				{Type: LINENUMBER, Value: "10", Line: 1, Column: 1},
				{Type: PRINT, Value: "PRINT", Line: 1, Column: 4},
				{Type: LINENUMBER, Value: "20", Line: 3, Column: 1},
				{Type: INPUT, Value: "INPUT", Line: 3, Column: 4},
				{Type: EOF, Value: "", Line: 3, Column: 9},
			},
		},
		{
			name:  "tabs and spaces should be handled correctly",
			input: "10\tPRINT\t\"Test\"",
			expected: []Token{
				{Type: LINENUMBER, Value: "10", Line: 1, Column: 1},
				{Type: PRINT, Value: "PRINT", Line: 1, Column: 4},
				{Type: STRING, Value: "Test", Line: 1, Column: 10},
				{Type: EOF, Value: "", Line: 1, Column: 16},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			
			for i, expectedToken := range tt.expected {
				token := lexer.NextToken()
				assertTokensEqual(t, expectedToken, token, i)
			}
		})
	}
}

// TestLexer_HasMoreTokens tests the HasMoreTokens method
func TestLexer_HasMoreTokens(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty input",
			input: "",
		},
		{
			name:  "single token",
			input: "123",
		},
		{
			name:  "multiple tokens",
			input: "10 PRINT X",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			
			// Should have tokens initially (at least EOF)
			assert.True(t, lexer.HasMoreTokens(), "HasMoreTokens() should return true initially")
			
			// Consume all tokens
			for {
				token := lexer.NextToken()
				if token.Type == EOF {
					break
				}
			}
			
			// Should still have EOF token available
			assert.True(t, lexer.HasMoreTokens(), "HasMoreTokens() should return true when EOF is available")
			
			// Get EOF token
			token := lexer.NextToken()
			assert.Equal(t, EOF, token.Type, "Expected EOF token")
			
			// After EOF, should still return true (EOF is always available)
			assert.True(t, lexer.HasMoreTokens(), "HasMoreTokens() should return true even after EOF")
		})
	}
}

// TestLexer_ComplexProgram tests tokenizing a complete BASIC program
func TestLexer_ComplexProgram(t *testing.T) {
	input := `10 FOR I = 1 TO 10
20 PRINT "Number: "; I
30 NEXT I
40 END`

	expected := []Token{
		{Type: LINENUMBER, Value: "10", Line: 1, Column: 1},
		{Type: FOR, Value: "FOR", Line: 1, Column: 4},
		{Type: IDENTIFIER, Value: "I", Line: 1, Column: 8},
		{Type: ASSIGN, Value: "=", Line: 1, Column: 10},
		{Type: NUMBER, Value: "1", Line: 1, Column: 12},
		{Type: TO, Value: "TO", Line: 1, Column: 14},
		{Type: NUMBER, Value: "10", Line: 1, Column: 17},
		{Type: LINENUMBER, Value: "20", Line: 2, Column: 1},
		{Type: PRINT, Value: "PRINT", Line: 2, Column: 4},
		{Type: STRING, Value: "Number: ", Line: 2, Column: 10},
		{Type: SEMICOLON, Value: ";", Line: 2, Column: 20},
		{Type: IDENTIFIER, Value: "I", Line: 2, Column: 22},
		{Type: LINENUMBER, Value: "30", Line: 3, Column: 1},
		{Type: NEXT, Value: "NEXT", Line: 3, Column: 4},
		{Type: IDENTIFIER, Value: "I", Line: 3, Column: 9},
		{Type: LINENUMBER, Value: "40", Line: 4, Column: 1},
		{Type: END, Value: "END", Line: 4, Column: 4},
		{Type: EOF, Value: "", Line: 4, Column: 7},
	}

	lexer := NewLexer(input)
	
	for i, expectedToken := range expected {
		token := lexer.NextToken()
		assertTokensEqual(t, expectedToken, token, i)
	}
}