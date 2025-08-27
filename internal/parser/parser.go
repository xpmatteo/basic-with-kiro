package parser

import (
	"basic-interpreter/internal/ast"
	"basic-interpreter/internal/lexer"
)

// Parser interface defines the contract for parsing BASIC source code
type Parser interface {
	ParseProgram() (*ast.Program, error)
	ParseStatement() (ast.Statement, error)
	ParseExpression() (ast.Expression, error)
}

// BasicParser implements the Parser interface
type BasicParser struct {
	lexer  lexer.Lexer
	curToken lexer.Token
	peekToken lexer.Token
}

// NewParser creates a new parser instance
func NewParser(l lexer.Lexer) *BasicParser {
	p := &BasicParser{
		lexer: l,
	}
	
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	
	return p
}

// nextToken advances the parser to the next token
func (p *BasicParser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}