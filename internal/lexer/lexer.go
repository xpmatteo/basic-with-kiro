package lexer

// TokenType represents the type of a token
type TokenType int

const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF

	// Literals
	NUMBER
	STRING
	IDENTIFIER

	// Keywords
	PRINT
	INPUT
	LET
	IF
	THEN
	GOTO
	FOR
	TO
	NEXT
	STEP
	END

	// Operators
	ASSIGN  // =
	PLUS    // +
	MINUS   // -
	MULTIPLY // *
	DIVIDE  // /
	POWER   // ^

	// Comparison operators
	EQ // =
	LT // <
	GT // >
	LE // <=
	GE // >=
	NE // <>

	// Delimiters
	SEMICOLON // ;
	COMMA     // ,
	LPAREN    // (
	RPAREN    // )

	// Line number
	LINENUMBER
)

// Token represents a single token
type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

// Lexer interface defines the contract for tokenizing BASIC source code
type Lexer interface {
	NextToken() Token
	HasMoreTokens() bool
}