package lexer

import "strings"

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

// String returns the string representation of a TokenType
func (t TokenType) String() string {
	switch t {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case NUMBER:
		return "NUMBER"
	case STRING:
		return "STRING"
	case IDENTIFIER:
		return "IDENTIFIER"
	case PRINT:
		return "PRINT"
	case INPUT:
		return "INPUT"
	case LET:
		return "LET"
	case IF:
		return "IF"
	case THEN:
		return "THEN"
	case GOTO:
		return "GOTO"
	case FOR:
		return "FOR"
	case TO:
		return "TO"
	case NEXT:
		return "NEXT"
	case STEP:
		return "STEP"
	case END:
		return "END"
	case ASSIGN:
		return "ASSIGN"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case MULTIPLY:
		return "MULTIPLY"
	case DIVIDE:
		return "DIVIDE"
	case POWER:
		return "POWER"
	case EQ:
		return "EQ"
	case LT:
		return "LT"
	case GT:
		return "GT"
	case LE:
		return "LE"
	case GE:
		return "GE"
	case NE:
		return "NE"
	case SEMICOLON:
		return "SEMICOLON"
	case COMMA:
		return "COMMA"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case LINENUMBER:
		return "LINENUMBER"
	default:
		return "UNKNOWN"
	}
}

// Lexer interface defines the contract for tokenizing BASIC source code
type Lexer interface {
	NextToken() Token
	HasMoreTokens() bool
}

// BasicLexer implements the Lexer interface for BASIC language tokenization
type BasicLexer struct {
	input         string // the input source code
	position      int    // current position in input (points to current char)
	readPosition  int    // current reading position in input (after current char)
	ch            byte   // current char under examination
	line          int    // current line number (1-based)
	column        int    // current column number (1-based)
	atStartOfLine bool   // true if we're at the start of a line (for line number detection)
}

// NewLexer creates a new lexer instance
func NewLexer(input string) Lexer {
	l := &BasicLexer{
		input:         input,
		line:          1,
		column:        0,
		atStartOfLine: true,
	}
	l.readChar()
	return l
}

// readChar reads the next character and advances position
func (l *BasicLexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII NUL character represents "EOF"
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	
	if l.ch == '\n' {
		l.line++
		l.column = 0
		l.atStartOfLine = true
	} else {
		l.column++
	}
}

// peekChar returns the next character without advancing position
func (l *BasicLexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// skipWhitespace skips whitespace characters except newlines
func (l *BasicLexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

// readString reads a string literal
func (l *BasicLexer) readString() (string, bool) {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	
	if l.ch == 0 {
		// Unterminated string
		return "unterminated string", false
	}
	
	result := l.input[position:l.position]
	l.readChar() // consume the closing quote and position after it
	return result, true
}

// readNumber reads a numeric literal (integer or decimal)
func (l *BasicLexer) readNumber() string {
	position := l.position
	
	// Read integer part
	for isDigit(l.ch) {
		l.readChar()
	}
	
	// Handle decimal point if followed by digits
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar() // consume '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	
	return l.input[position:l.position]
}

// readIdentifier reads an identifier (letters, digits, and $ for string variables)
func (l *BasicLexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '$' {
		l.readChar()
	}
	return l.input[position:l.position]
}

// isLetter checks if character is a letter
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit checks if character is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// keywords maps keyword strings to their token types (case-insensitive)
var keywords = map[string]TokenType{
	"PRINT": PRINT,
	"INPUT": INPUT,
	"LET":   LET,
	"IF":    IF,
	"THEN":  THEN,
	"GOTO":  GOTO,
	"FOR":   FOR,
	"TO":    TO,
	"NEXT":  NEXT,
	"STEP":  STEP,
	"END":   END,
}

// lookupIdent checks if identifier is a keyword (case-insensitive)
func lookupIdent(ident string) TokenType {
	// Convert to uppercase for case-insensitive lookup
	upperIdent := strings.ToUpper(ident)
	if tok, ok := keywords[upperIdent]; ok {
		return tok
	}
	return IDENTIFIER
}

// makeSingleCharToken creates a token for a single character
func (l *BasicLexer) makeSingleCharToken(tokenType TokenType, line, column int) Token {
	l.atStartOfLine = false
	return Token{Type: tokenType, Value: string(l.ch), Line: line, Column: column}
}

// makeTwoCharToken creates a token for a two-character operator
func (l *BasicLexer) makeTwoCharToken(tokenType TokenType, line, column int) Token {
	ch := l.ch
	l.readChar()
	l.atStartOfLine = false
	return Token{Type: tokenType, Value: string(ch) + string(l.ch), Line: line, Column: column}
}

// NextToken scans the input and returns the next token
func (l *BasicLexer) NextToken() Token {
	var tok Token
	
	l.skipWhitespace()
	
	startColumn := l.column
	startLine := l.line
	
	switch l.ch {
	case '=':
		tok = l.makeSingleCharToken(ASSIGN, startLine, startColumn)
	case '+':
		tok = l.makeSingleCharToken(PLUS, startLine, startColumn)
	case '-':
		tok = l.makeSingleCharToken(MINUS, startLine, startColumn)
	case '*':
		tok = l.makeSingleCharToken(MULTIPLY, startLine, startColumn)
	case '/':
		tok = l.makeSingleCharToken(DIVIDE, startLine, startColumn)
	case '^':
		tok = l.makeSingleCharToken(POWER, startLine, startColumn)
	case '<':
		if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(LE, startLine, startColumn)
		} else if l.peekChar() == '>' {
			tok = l.makeTwoCharToken(NE, startLine, startColumn)
		} else {
			tok = l.makeSingleCharToken(LT, startLine, startColumn)
		}
	case '>':
		if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(GE, startLine, startColumn)
		} else {
			tok = l.makeSingleCharToken(GT, startLine, startColumn)
		}
	case ';':
		tok = l.makeSingleCharToken(SEMICOLON, startLine, startColumn)
	case ',':
		tok = l.makeSingleCharToken(COMMA, startLine, startColumn)
	case '(':
		tok = l.makeSingleCharToken(LPAREN, startLine, startColumn)
	case ')':
		tok = l.makeSingleCharToken(RPAREN, startLine, startColumn)
	case '"':
		return l.readStringToken(startLine, startColumn)
	case '\n':
		return l.handleNewline()
	case 0:
		tok = Token{Type: EOF, Value: "", Line: startLine, Column: startColumn}
	default:
		if isLetter(l.ch) {
			return l.readIdentifierToken(startLine, startColumn)
		} else if isDigit(l.ch) || (l.ch == '.' && isDigit(l.peekChar())) {
			return l.readNumberToken(startLine, startColumn)
		} else {
			tok = Token{Type: ILLEGAL, Value: string(l.ch), Line: startLine, Column: startColumn}
			l.atStartOfLine = false
		}
	}
	
	l.readChar()
	return tok
}

// readStringToken handles string literal tokenization
func (l *BasicLexer) readStringToken(line, column int) Token {
	value, ok := l.readString()
	l.atStartOfLine = false
	if !ok {
		return Token{Type: ILLEGAL, Value: value, Line: line, Column: column}
	}
	return Token{Type: STRING, Value: value, Line: line, Column: column}
}

// readIdentifierToken handles identifier and keyword tokenization
func (l *BasicLexer) readIdentifierToken(line, column int) Token {
	value := l.readIdentifier()
	tokenType := lookupIdent(value)
	l.atStartOfLine = false
	return Token{Type: tokenType, Value: value, Line: line, Column: column}
}

// readNumberToken handles number and line number tokenization
func (l *BasicLexer) readNumberToken(line, column int) Token {
	value := l.readNumber()
	var tokenType TokenType
	if l.atStartOfLine && l.isFollowedByKeyword() {
		tokenType = LINENUMBER
	} else {
		tokenType = NUMBER
	}
	l.atStartOfLine = false
	return Token{Type: tokenType, Value: value, Line: line, Column: column}
}

// handleNewline processes newline characters
func (l *BasicLexer) handleNewline() Token {
	l.readChar()
	l.atStartOfLine = true
	return l.NextToken()
}

// isLineNumber determines if a number token should be treated as a line number
func (l *BasicLexer) isLineNumber(value string) bool {
	// A number is a line number if it appears at the start of a line
	// We need to track if we've seen any non-whitespace on this line
	return l.atStartOfLine
}

// isFollowedByKeyword checks if the current position is followed by whitespace and then a keyword
func (l *BasicLexer) isFollowedByKeyword() bool {
	// Look ahead in the input without modifying lexer state
	pos := l.readPosition
	
	// Skip whitespace
	for pos < len(l.input) && (l.input[pos] == ' ' || l.input[pos] == '\t') {
		pos++
	}
	
	// Check if we have a letter (start of identifier/keyword)
	if pos >= len(l.input) || !isLetter(l.input[pos]) {
		return false
	}
	
	// Read the identifier
	start := pos
	for pos < len(l.input) && (isLetter(l.input[pos]) || isDigit(l.input[pos]) || l.input[pos] == '$') {
		pos++
	}
	
	if start == pos {
		return false
	}
	
	// Check if it's a keyword
	ident := l.input[start:pos]
	return lookupIdent(ident) != IDENTIFIER
}

// HasMoreTokens returns true if there are more tokens to read
func (l *BasicLexer) HasMoreTokens() bool {
	return true // Always return true as EOF is always available
}