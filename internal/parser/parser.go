package parser

import (
	"basic-interpreter/internal/ast"
	"basic-interpreter/internal/lexer"
	"basic-interpreter/internal/runtime"
	"fmt"
	"strconv"
)

// Parser interface defines the contract for parsing BASIC source code
type Parser interface {
	ParseProgram() (*ast.Program, error)
	ParseStatement() (ast.Statement, error)
	ParseExpression() (ast.Expression, error)
}

// BasicParser implements the Parser interface
type BasicParser struct {
	lexer           lexer.Lexer
	curToken        lexer.Token
	peekToken       lexer.Token
	currentLineNumber int
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

// ParseProgram parses a complete BASIC program
func (p *BasicParser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{
		Lines: make(map[int]ast.Statement),
		Order: []int{},
	}
	
	// Parse all statements in the program
	for p.curToken.Type != lexer.EOF {
		// Skip empty lines and whitespace
		if p.curToken.Type == lexer.EOF {
			break
		}
		
		// Expect line number
		if p.curToken.Type != lexer.LINENUMBER && p.curToken.Type != lexer.NUMBER {
			return nil, fmt.Errorf("expected line number at start of statement")
		}
		
		// Parse line number
		lineNumber, err := p.parseLineNumber()
		if err != nil {
			return nil, err
		}
		
		// Check for duplicate line numbers
		if _, exists := program.Lines[lineNumber]; exists {
			return nil, fmt.Errorf("duplicate line number: %d", lineNumber)
		}
		
		// Set current line number for statement parsing
		p.currentLineNumber = lineNumber
		
		// Parse the statement for this line
		stmt, err := p.ParseStatement()
		if err != nil {
			return nil, fmt.Errorf("error parsing statement at line %d: %w", lineNumber, err)
		}
		
		// Add to program
		program.Lines[lineNumber] = stmt
		program.Order = append(program.Order, lineNumber)
	}
	
	// Sort line numbers
	p.sortLineNumbers(program)
	
	return program, nil
}

// parseLineNumber parses and validates a line number
func (p *BasicParser) parseLineNumber() (int, error) {
	if !p.isLineNumberToken() {
		return 0, fmt.Errorf("expected line number")
	}
	
	lineNumber, err := p.convertToLineNumber(p.curToken.Value)
	if err != nil {
		return 0, err
	}
	
	p.nextToken() // consume line number
	return lineNumber, nil
}

// isLineNumberToken checks if current token can be a line number
func (p *BasicParser) isLineNumberToken() bool {
	return p.curToken.Type == lexer.LINENUMBER || p.curToken.Type == lexer.NUMBER
}

// convertToLineNumber converts a string to a valid line number
func (p *BasicParser) convertToLineNumber(value string) (int, error) {
	lineNumber, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid line number: %s", value)
	}
	
	return p.validateLineNumberRange(lineNumber)
}

// validateLineNumberRange validates that a line number is within valid range
func (p *BasicParser) validateLineNumberRange(lineNumber int) (int, error) {
	const (
		minLineNumber = 1
		maxLineNumber = 99999
	)
	
	if lineNumber < minLineNumber {
		return 0, fmt.Errorf("invalid line number: %d (must be >= %d)", lineNumber, minLineNumber)
	}
	
	if lineNumber > maxLineNumber {
		return 0, fmt.Errorf("invalid line number: %d (must be <= %d)", lineNumber, maxLineNumber)
	}
	
	return lineNumber, nil
}



// sortLineNumbers sorts the line numbers in the program order
func (p *BasicParser) sortLineNumbers(program *ast.Program) {
	// Simple insertion sort for line numbers
	for i := 1; i < len(program.Order); i++ {
		key := program.Order[i]
		j := i - 1
		
		for j >= 0 && program.Order[j] > key {
			program.Order[j+1] = program.Order[j]
			j--
		}
		program.Order[j+1] = key
	}
}

// parseBinaryExpression is a generic helper for parsing left-associative binary expressions
func (p *BasicParser) parseBinaryExpression(
	parseNext func() (ast.Expression, error),
	operators []lexer.TokenType,
) (ast.Expression, error) {
	left, err := parseNext()
	if err != nil {
		return nil, err
	}
	
	for p.isTokenInList(p.curToken.Type, operators) {
		operator := p.getArithmeticOperator(p.curToken.Type)
		p.nextToken() // consume operator
		
		// Check for end of input
		if p.curToken.Type == lexer.EOF {
			return nil, fmt.Errorf("unexpected end of input after operator %s", operator)
		}
		
		right, err := parseNext()
		if err != nil {
			return nil, err
		}
		
		left = ast.NewBinaryExpression(left, operator, right)
	}
	
	return left, nil
}

// isTokenInList checks if a token type is in the given list
func (p *BasicParser) isTokenInList(tokenType lexer.TokenType, list []lexer.TokenType) bool {
	for _, t := range list {
		if tokenType == t {
			return true
		}
	}
	return false
}

// expectToken checks if current token matches expected type
func (p *BasicParser) expectToken(expected lexer.TokenType, name string) error {
	if p.curToken.Type != expected {
		return fmt.Errorf("expected %s", name)
	}
	return nil
}

// isEndOfStatement checks if we're at the end of a statement
func (p *BasicParser) isEndOfStatement() bool {
	return p.curToken.Type == lexer.EOF || p.curToken.Type == lexer.LINENUMBER
}

// parseExpressionList parses a comma-separated list of expressions
func (p *BasicParser) parseExpressionList() ([]ast.Expression, error) {
	var expressions []ast.Expression
	
	// Parse first expression
	expr, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	expressions = append(expressions, expr)
	
	// Parse additional expressions separated by commas
	for p.curToken.Type == lexer.COMMA {
		p.nextToken() // consume comma
		
		expr, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
	}
	
	return expressions, nil
}

// ParseStatement parses a single BASIC statement (without line number)
func (p *BasicParser) ParseStatement() (ast.Statement, error) {
	// Parse based on statement type
	switch p.curToken.Type {
	case lexer.PRINT:
		return p.parsePrintStatement()
	case lexer.INPUT:
		return p.parseInputStatement()
	case lexer.GOTO:
		return p.parseGotoStatement()
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.FOR:
		return p.parseForStatement()
	case lexer.NEXT:
		return p.parseNextStatement()
	case lexer.IDENTIFIER:
		return p.parseAssignmentStatement()
	case lexer.NUMBER:
		// Treat numbers as identifiers for assignment (like variable names that are numbers)
		return p.parseAssignmentStatement()
	case lexer.EOF:
		return nil, fmt.Errorf("unexpected end of input")
	default:
		return nil, fmt.Errorf("unknown statement type: %s", p.curToken.Value)
	}
}

// parsePrintStatement parses a PRINT statement
func (p *BasicParser) parsePrintStatement() (ast.Statement, error) {
	if err := p.expectToken(lexer.PRINT, "PRINT"); err != nil {
		return nil, err
	}
	
	p.nextToken() // consume PRINT
	
	// Handle empty PRINT statement
	if p.isEndOfStatement() {
		return ast.NewPrintStatement([]ast.Expression{}, nil), nil
	}
	
	// Parse comma-separated expressions
	expressions, err := p.parseExpressionList()
	if err != nil {
		return nil, fmt.Errorf("error parsing PRINT expressions: %w", err)
	}
	
	return ast.NewPrintStatement(expressions, nil), nil
}

// parseInputStatement parses an INPUT statement
func (p *BasicParser) parseInputStatement() (ast.Statement, error) {
	if p.curToken.Type != lexer.INPUT {
		return nil, fmt.Errorf("expected INPUT")
	}
	
	p.nextToken() // consume INPUT
	
	var prompt string
	
	// Check for optional prompt string
	if p.curToken.Type == lexer.STRING {
		prompt = p.curToken.Value
		p.nextToken() // consume string
		
		// Expect semicolon after prompt
		if p.curToken.Type == lexer.SEMICOLON {
			p.nextToken() // consume semicolon
		}
	}
	
	// Expect variable name
	if p.curToken.Type != lexer.IDENTIFIER {
		return nil, fmt.Errorf("expected variable name in INPUT statement")
	}
	
	variable := p.curToken.Value
	p.nextToken() // consume variable
	
	if prompt != "" {
		return ast.NewInputStatementWithPrompt(prompt, variable, nil, nil), nil
	}
	return ast.NewInputStatement(variable, nil, nil), nil
}

// parseGotoStatement parses a GOTO statement
func (p *BasicParser) parseGotoStatement() (ast.Statement, error) {
	if p.curToken.Type != lexer.GOTO {
		return nil, fmt.Errorf("expected GOTO")
	}
	
	p.nextToken() // consume GOTO
	
	// Expect line number
	if p.curToken.Type != lexer.NUMBER {
		return nil, fmt.Errorf("expected line number after GOTO")
	}
	
	lineNumber, err := strconv.Atoi(p.curToken.Value)
	if err != nil {
		return nil, fmt.Errorf("invalid line number: %s", p.curToken.Value)
	}
	
	p.nextToken() // consume line number
	
	return ast.NewGotoStatement(lineNumber, nil), nil
}

// parseIfStatement parses an IF-THEN statement
func (p *BasicParser) parseIfStatement() (ast.Statement, error) {
	if p.curToken.Type != lexer.IF {
		return nil, fmt.Errorf("expected IF")
	}
	
	p.nextToken() // consume IF
	
	// Parse condition expression
	condition, err := p.ParseExpression()
	if err != nil {
		return nil, fmt.Errorf("error parsing IF condition: %w", err)
	}
	
	// Expect THEN
	if p.curToken.Type != lexer.THEN {
		return nil, fmt.Errorf("expected THEN after IF condition")
	}
	
	p.nextToken() // consume THEN
	
	// Parse THEN statement
	thenStatement, err := p.ParseStatement()
	if err != nil {
		return nil, fmt.Errorf("error parsing THEN statement: %w", err)
	}
	
	return ast.NewIfStatement(condition, thenStatement), nil
}

// parseForStatement parses a FOR statement
func (p *BasicParser) parseForStatement() (ast.Statement, error) {
	if p.curToken.Type != lexer.FOR {
		return nil, fmt.Errorf("expected FOR")
	}
	
	p.nextToken() // consume FOR
	
	// Expect variable name
	if p.curToken.Type != lexer.IDENTIFIER {
		return nil, fmt.Errorf("expected variable name in FOR statement")
	}
	
	variable := p.curToken.Value
	p.nextToken() // consume variable
	
	// Expect assignment operator
	if p.curToken.Type != lexer.ASSIGN {
		return nil, fmt.Errorf("expected = in FOR statement")
	}
	
	p.nextToken() // consume =
	
	// Parse start expression
	startExpr, err := p.ParseExpression()
	if err != nil {
		return nil, fmt.Errorf("error parsing FOR start expression: %w", err)
	}
	
	// Expect TO
	if p.curToken.Type != lexer.TO {
		return nil, fmt.Errorf("expected TO in FOR statement")
	}
	
	p.nextToken() // consume TO
	
	// Parse end expression
	endExpr, err := p.ParseExpression()
	if err != nil {
		return nil, fmt.Errorf("error parsing FOR end expression: %w", err)
	}
	
	// Check for optional STEP
	var stepExpr ast.Expression = ast.NewLiteralExpression(runtime.NewNumericValue(1)) // default step
	
	if p.curToken.Type == lexer.STEP {
		p.nextToken() // consume STEP
		
		stepExpr, err = p.ParseExpression()
		if err != nil {
			return nil, fmt.Errorf("error parsing FOR step expression: %w", err)
		}
	}
	
	return ast.NewForStatement(variable, startExpr, endExpr, stepExpr, p.currentLineNumber), nil
}

// parseNextStatement parses a NEXT statement
func (p *BasicParser) parseNextStatement() (ast.Statement, error) {
	if p.curToken.Type != lexer.NEXT {
		return nil, fmt.Errorf("expected NEXT")
	}
	
	p.nextToken() // consume NEXT
	
	var variable string
	
	// Check for optional variable name
	if p.curToken.Type == lexer.IDENTIFIER {
		variable = p.curToken.Value
		p.nextToken() // consume variable
	}
	
	return ast.NewNextStatement(variable), nil
}

// parseAssignmentStatement parses an assignment statement
func (p *BasicParser) parseAssignmentStatement() (ast.Statement, error) {
	if p.curToken.Type != lexer.IDENTIFIER && p.curToken.Type != lexer.NUMBER {
		return nil, fmt.Errorf("expected variable name")
	}
	
	variable := p.curToken.Value
	p.nextToken() // consume variable
	
	// Expect assignment operator
	if p.curToken.Type != lexer.ASSIGN {
		return nil, fmt.Errorf("expected assignment operator")
	}
	
	p.nextToken() // consume =
	
	// Parse expression
	expr, err := p.ParseExpression()
	if err != nil {
		return nil, fmt.Errorf("error parsing assignment expression: %w", err)
	}
	
	return ast.NewAssignmentStatement(variable, expr), nil
}

// ParseExpression parses an expression (placeholder implementation)
func (p *BasicParser) ParseExpression() (ast.Expression, error) {
	return p.parseComparison()
}

// parseComparison parses comparison expressions
func (p *BasicParser) parseComparison() (ast.Expression, error) {
	left, err := p.parseArithmetic()
	if err != nil {
		return nil, err
	}
	
	// Check for comparison operators
	if p.isComparisonOperator(p.curToken.Type) {
		operator := p.getComparisonOperator(p.curToken.Type)
		p.nextToken() // consume operator
		
		right, err := p.parseArithmetic()
		if err != nil {
			return nil, err
		}
		
		return ast.NewComparisonExpression(left, operator, right), nil
	}
	
	return left, nil
}

// parseArithmetic parses arithmetic expressions with precedence
func (p *BasicParser) parseArithmetic() (ast.Expression, error) {
	return p.parseAddition()
}

// parseAddition parses addition and subtraction (lowest precedence)
func (p *BasicParser) parseAddition() (ast.Expression, error) {
	return p.parseBinaryExpression(
		p.parseMultiplication,
		[]lexer.TokenType{lexer.PLUS, lexer.MINUS},
	)
}

// parseMultiplication parses multiplication and division (medium precedence)
func (p *BasicParser) parseMultiplication() (ast.Expression, error) {
	return p.parseBinaryExpression(
		p.parsePower,
		[]lexer.TokenType{lexer.MULTIPLY, lexer.DIVIDE},
	)
}

// parsePower parses power expressions (highest precedence)
func (p *BasicParser) parsePower() (ast.Expression, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	
	if p.curToken.Type == lexer.POWER {
		operator := p.getArithmeticOperator(p.curToken.Type)
		p.nextToken() // consume operator
		
		// Check for end of input
		if p.curToken.Type == lexer.EOF {
			return nil, fmt.Errorf("unexpected end of input after operator %s", operator)
		}
		
		// Power is right-associative
		right, err := p.parsePower()
		if err != nil {
			return nil, err
		}
		
		return ast.NewBinaryExpression(left, operator, right), nil
	}
	
	return left, nil
}

// parsePrimary parses primary expressions (literals, variables, parentheses, function calls)
func (p *BasicParser) parsePrimary() (ast.Expression, error) {
	switch p.curToken.Type {
	case lexer.NUMBER:
		return p.parseNumberLiteral()
	case lexer.STRING:
		return p.parseStringLiteral()
	case lexer.IDENTIFIER:
		return p.parseIdentifierOrFunction()
	case lexer.LPAREN:
		return p.parseParentheses()
	case lexer.MINUS:
		// Handle unary minus
		return p.parseUnaryMinus()
	case lexer.EOF:
		return nil, fmt.Errorf("unexpected end of input")
	default:
		return nil, fmt.Errorf("unexpected token in expression: %s", p.curToken.Value)
	}
}

// parseUnaryMinus parses unary minus expressions
func (p *BasicParser) parseUnaryMinus() (ast.Expression, error) {
	if p.curToken.Type != lexer.MINUS {
		return nil, fmt.Errorf("expected -")
	}
	
	p.nextToken() // consume -
	
	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	
	// Create a binary expression: 0 - expr
	zero := ast.NewLiteralExpression(runtime.NewNumericValue(0))
	return ast.NewBinaryExpression(zero, "-", expr), nil
}

// parseNumberLiteral parses a numeric literal
func (p *BasicParser) parseNumberLiteral() (ast.Expression, error) {
	if p.curToken.Type != lexer.NUMBER {
		return nil, fmt.Errorf("expected number")
	}
	
	value, err := strconv.ParseFloat(p.curToken.Value, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number: %s", p.curToken.Value)
	}
	
	p.nextToken() // consume number
	
	return ast.NewLiteralExpression(runtime.NewNumericValue(value)), nil
}

// parseStringLiteral parses a string literal
func (p *BasicParser) parseStringLiteral() (ast.Expression, error) {
	if p.curToken.Type != lexer.STRING {
		return nil, fmt.Errorf("expected string")
	}
	
	value := p.curToken.Value
	p.nextToken() // consume string
	
	return ast.NewLiteralExpression(runtime.NewStringValue(value)), nil
}

// parseIdentifierOrFunction parses a variable reference or function call
func (p *BasicParser) parseIdentifierOrFunction() (ast.Expression, error) {
	if p.curToken.Type != lexer.IDENTIFIER {
		return nil, fmt.Errorf("expected identifier")
	}
	
	name := p.curToken.Value
	p.nextToken() // consume identifier
	
	// Check if this is a function call with parentheses (even for unknown functions)
	if p.curToken.Type == lexer.LPAREN {
		return p.parseFunctionCall(name)
	}
	
	// Check if this is a known function without parentheses (like RND)
	if ast.IsFunctionRegistered(name) {
		// Create function call with no arguments
		return ast.NewFunctionCallExpression(name, []ast.Expression{}), nil
	}
	
	// It's a variable reference
	return ast.NewVariableExpression(name), nil
}

// parseFunctionCall parses a function call
func (p *BasicParser) parseFunctionCall(name string) (ast.Expression, error) {
	if p.curToken.Type != lexer.LPAREN {
		return nil, fmt.Errorf("expected ( in function call")
	}
	
	p.nextToken() // consume (
	
	var args []ast.Expression
	
	// Handle empty argument list
	if p.curToken.Type == lexer.RPAREN {
		p.nextToken() // consume )
		return ast.NewFunctionCallExpression(name, args), nil
	}
	
	// Parse arguments
	for {
		arg, err := p.ParseExpression()
		if err != nil {
			return nil, fmt.Errorf("error parsing function argument: %w", err)
		}
		args = append(args, arg)
		
		if p.curToken.Type == lexer.COMMA {
			p.nextToken() // consume comma
			continue
		}
		
		break
	}
	
	// Expect closing parenthesis
	if p.curToken.Type != lexer.RPAREN {
		return nil, fmt.Errorf("expected ) in function call")
	}
	
	p.nextToken() // consume )
	
	return ast.NewFunctionCallExpression(name, args), nil
}

// parseParentheses parses parenthesized expressions
func (p *BasicParser) parseParentheses() (ast.Expression, error) {
	if p.curToken.Type != lexer.LPAREN {
		return nil, fmt.Errorf("expected (")
	}
	
	p.nextToken() // consume (
	
	expr, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	
	if p.curToken.Type != lexer.RPAREN {
		return nil, fmt.Errorf("expected )")
	}
	
	p.nextToken() // consume )
	
	return ast.NewParenthesesExpression(expr), nil
}

// Helper functions for operator handling

// isComparisonOperator checks if token is a comparison operator
func (p *BasicParser) isComparisonOperator(tokenType lexer.TokenType) bool {
	switch tokenType {
	case lexer.EQ, lexer.ASSIGN, lexer.LT, lexer.GT, lexer.LE, lexer.GE, lexer.NE:
		return true
	default:
		return false
	}
}

// getComparisonOperator converts token type to operator string
func (p *BasicParser) getComparisonOperator(tokenType lexer.TokenType) string {
	switch tokenType {
	case lexer.EQ, lexer.ASSIGN:
		return "="
	case lexer.LT:
		return "<"
	case lexer.GT:
		return ">"
	case lexer.LE:
		return "<="
	case lexer.GE:
		return ">="
	case lexer.NE:
		return "<>"
	default:
		return ""
	}
}

// getArithmeticOperator converts token type to operator string
func (p *BasicParser) getArithmeticOperator(tokenType lexer.TokenType) string {
	switch tokenType {
	case lexer.PLUS:
		return "+"
	case lexer.MINUS:
		return "-"
	case lexer.MULTIPLY:
		return "*"
	case lexer.DIVIDE:
		return "/"
	case lexer.POWER:
		return "^"
	default:
		return ""
	}
}