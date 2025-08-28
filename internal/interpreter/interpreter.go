package interpreter

import (
	"basic-interpreter/internal/ast"
	"basic-interpreter/internal/runtime"
	"fmt"
	"strings"
)

// OutputWriter interface for debug output
type OutputWriter interface {
	WriteLine(line string) error
}

// Interpreter represents the BASIC interpreter
type Interpreter struct {
	debugMode    bool
	debugOutput  OutputWriter
	maxSteps     int
	stepCount    int
}

// InterpreterConfig holds configuration options for the interpreter
type InterpreterConfig struct {
	DebugMode   bool
	DebugOutput OutputWriter
	MaxSteps    int // -1 for no limit
}

// NewInterpreter creates a new interpreter instance with the given configuration
func NewInterpreter(config InterpreterConfig) *Interpreter {
	if config.MaxSteps == 0 {
		config.MaxSteps = -1 // Default to no limit
	}
	
	return &Interpreter{
		debugMode:   config.DebugMode,
		debugOutput: config.DebugOutput,
		maxSteps:    config.MaxSteps,
		stepCount:   0,
	}
}

// NewBasicInterpreter creates a new interpreter with default settings
func NewBasicInterpreter(debugMode bool) *Interpreter {
	return NewInterpreter(InterpreterConfig{
		DebugMode: debugMode,
		MaxSteps:  -1,
	})
}

// GetStepCount returns the number of execution steps performed
func (i *Interpreter) GetStepCount() int {
	return i.stepCount
}

// Execute executes a BASIC program
func (i *Interpreter) Execute(program *ast.Program, env *runtime.Environment) error {
	if program == nil {
		return nil
	}

	// Reset step counter
	i.stepCount = 0

	// Start execution from the first line
	currentIndex := 0

	for currentIndex < len(program.Order) {
		// Check execution step limit
		if err := i.checkStepLimit(); err != nil {
			return err
		}

		lineNumber := program.Order[currentIndex]
		
		// Set current program counter
		env.ProgramCounter = lineNumber

		// Get the statement for this line
		statement, exists := program.Lines[lineNumber]
		if !exists {
			currentIndex++
			continue // Skip missing lines
		}

		// Debug output: show line before execution
		i.outputDebugMessage(lineNumber, statement)

		// Increment step counter
		i.stepCount++

		// Store original program counter to detect changes
		originalPC := env.ProgramCounter

		// Execute the statement
		err := statement.Execute(env)
		if err != nil {
			// Wrap error with line number information
			return fmt.Errorf("runtime error at line %d: %w", lineNumber, err)
		}

		// Handle program counter changes and determine next execution position
		nextIndex, shouldBreak := i.handleProgramCounterChange(program, statement, lineNumber, originalPC, currentIndex, env)
		if shouldBreak {
			break
		}
		currentIndex = nextIndex
	}

	return nil
}

// findNextLineIndex finds the index in program.Order for the given line number
func (i *Interpreter) findNextLineIndex(program *ast.Program, lineNumber int) int {
	for idx, line := range program.Order {
		if line == lineNumber {
			return idx
		}
	}
	return -1
}

// isForStatement checks if the statement at the given line number is a FOR statement
func (i *Interpreter) isForStatement(program *ast.Program, lineNumber int) bool {
	statement, exists := program.Lines[lineNumber]
	if !exists {
		return false
	}
	_, isFor := statement.(*ast.ForStatement)
	return isFor
}

// isComingFromNext checks if we're coming from a NEXT statement
func (i *Interpreter) isComingFromNext(program *ast.Program, previousLineNumber int) bool {
	statement, exists := program.Lines[previousLineNumber]
	if !exists {
		return false
	}
	_, isNext := statement.(*ast.NextStatement)
	return isNext
}

// isGotoStatement checks if a statement is a GOTO statement
func (i *Interpreter) isGotoStatement(statement ast.Statement) bool {
	_, isGoto := statement.(*ast.GotoStatement)
	return isGoto
}

// outputDebugMessage outputs debug information if debug mode is enabled
func (i *Interpreter) outputDebugMessage(lineNumber int, statement ast.Statement) {
	if i.debugMode && i.debugOutput != nil {
		debugMsg := i.formatDebugMessage(lineNumber, statement)
		i.debugOutput.WriteLine(debugMsg)
	}
}

// checkStepLimit checks if the execution step limit has been exceeded
func (i *Interpreter) checkStepLimit() error {
	if i.maxSteps > 0 && i.stepCount >= i.maxSteps {
		return fmt.Errorf("execution limit exceeded: maximum %d steps reached", i.maxSteps)
	}
	return nil
}

// handleProgramCounterChange handles program counter modifications and returns the next execution index
func (i *Interpreter) handleProgramCounterChange(program *ast.Program, statement ast.Statement, lineNumber int, originalPC int, currentIndex int, env *runtime.Environment) (int, bool) {
	// Check if program counter was modified by control flow statements
	if !i.hasProgramCounterChanged(originalPC, lineNumber, statement, env) {
		return currentIndex + 1, false // Normal sequential execution
	}
	
	// Find the next line to execute based on the new program counter
	nextIndex := i.findNextLineIndex(program, env.ProgramCounter)
	if nextIndex == -1 {
		return currentIndex, true // Program counter points to non-existent line, end execution
	}
	
	// Handle special case for FOR-NEXT loops
	if i.isForStatement(program, env.ProgramCounter) && i.isComingFromNext(program, lineNumber) {
		return nextIndex + 1, false // Continue from the line AFTER the FOR statement
	}
	
	return nextIndex, false // Continue execution from the new position (GOTO case)
}

// hasProgramCounterChanged checks if the program counter was modified by a control flow statement
func (i *Interpreter) hasProgramCounterChanged(originalPC, currentLine int, statement ast.Statement, env *runtime.Environment) bool {
	// GOTO statements always change program counter (even to same line)
	if i.isGotoStatement(statement) {
		return true
	}
	
	// Other statements change PC if it differs from current line
	return env.ProgramCounter != currentLine
}

// formatDebugMessage formats a debug message for a statement
func (i *Interpreter) formatDebugMessage(lineNumber int, statement ast.Statement) string {
	switch stmt := statement.(type) {
	case *ast.AssignmentStatement:
		return fmt.Sprintf("Executing line %d: %s = %s", lineNumber, stmt.Variable, i.formatExpression(stmt.Expression))
	case *ast.PrintStatement:
		return fmt.Sprintf("Executing line %d: PRINT %s", lineNumber, i.formatExpressionList(stmt.Expressions))
	case *ast.InputStatement:
		return fmt.Sprintf("Executing line %d: INPUT %s", lineNumber, stmt.Variable)
	case *ast.GotoStatement:
		return fmt.Sprintf("Executing line %d: GOTO %d", lineNumber, stmt.LineNumber)
	case *ast.IfStatement:
		return fmt.Sprintf("Executing line %d: IF %s THEN ...", lineNumber, i.formatExpression(stmt.Condition))
	case *ast.ForStatement:
		return fmt.Sprintf("Executing line %d: FOR %s = %s TO %s STEP %s", lineNumber, stmt.Variable, 
			i.formatExpression(stmt.StartExpr), i.formatExpression(stmt.EndExpr), i.formatExpression(stmt.StepExpr))
	case *ast.NextStatement:
		return fmt.Sprintf("Executing line %d: NEXT %s", lineNumber, stmt.Variable)
	default:
		return fmt.Sprintf("Executing line %d", lineNumber)
	}
}

// formatExpression formats an expression for debug output
func (i *Interpreter) formatExpression(expr ast.Expression) string {
	switch e := expr.(type) {
	case *ast.LiteralExpression:
		if e.Value.Type == runtime.StringValue {
			return fmt.Sprintf("\"%s\"", e.Value.StrValue)
		}
		return fmt.Sprintf("%.0f", e.Value.NumValue)
	case *ast.VariableExpression:
		return e.Name
	case *ast.BinaryExpression:
		return fmt.Sprintf("%s %s %s", i.formatExpression(e.Left), e.Operator, i.formatExpression(e.Right))
	default:
		return "..."
	}
}

// formatExpressionList formats a list of expressions for debug output
func (i *Interpreter) formatExpressionList(expressions []ast.Expression) string {
	if len(expressions) == 0 {
		return ""
	}
	
	var parts []string
	for _, expr := range expressions {
		parts = append(parts, i.formatExpression(expr))
	}
	return strings.Join(parts, ", ")
}