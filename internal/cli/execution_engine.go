package cli

import (
	"basic-interpreter/internal/ast"
	"basic-interpreter/internal/interpreter"
	"basic-interpreter/internal/lexer"
	"basic-interpreter/internal/parser"
	"basic-interpreter/internal/runtime"
	"fmt"
	"strings"
)

// ForLoopInfo tracks FOR loop state
type ForLoopInfo struct {
	Variable  string
	Current   int
	End       int
	Step      int
	StartLine int
}

// ExecuteProgram executes a parsed BASIC program using the real interpreter
func (fe *FileExecutor) ExecuteProgram(program map[int]string, debugMode bool) error {
	if len(program) == 0 {
		return nil // Empty program is valid
	}
	
	// Convert the string-based program to source code
	sourceCode := fe.programToSourceCode(program)
	
	// Create lexer and parser
	lex := lexer.NewLexer(sourceCode)
	p := parser.NewParser(lex)
	
	// Parse the program
	astProgram, err := p.ParseProgram()
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}
	
	// Set output writer for all PRINT statements
	fe.setPrintOutputWriters(astProgram)
	
	// Set input/output writers for all INPUT statements
	fe.setInputOutputWriters(astProgram)
	
	// Create runtime environment
	env := runtime.NewEnvironment()
	
	// Create interpreter with debug output if needed
	var interpreterInstance *interpreter.Interpreter
	if debugMode {
		interpreterInstance = interpreter.NewInterpreter(interpreter.InterpreterConfig{
			DebugMode:   true,
			DebugOutput: fe.output,
			MaxSteps:    -1,
		})
	} else {
		interpreterInstance = interpreter.NewBasicInterpreter(false)
	}
	
	// Execute the program
	return interpreterInstance.Execute(astProgram, env)
}

// programToSourceCode converts a map[int]string program to source code
func (fe *FileExecutor) programToSourceCode(program map[int]string) string {
	// Get sorted line numbers
	lineNumbers := fe.getSortedLineNumbers(program)
	
	var lines []string
	for _, lineNum := range lineNumbers {
		if statement, exists := program[lineNum]; exists {
			lines = append(lines, fmt.Sprintf("%d %s", lineNum, statement))
		}
	}
	
	return strings.Join(lines, "\n")
}

// ExecutionContext holds the state for program execution
type ExecutionContext struct {
	lineNumbers []int
	program     map[int]string
	variables   map[string]interface{}
	forLoops    map[string]ForLoopInfo
	pc          int
}

// createExecutionContext creates a new execution context
func (fe *FileExecutor) createExecutionContext(program map[int]string) *ExecutionContext {
	return &ExecutionContext{
		lineNumbers: fe.getSortedLineNumbers(program),
		program:     program,
		variables:   make(map[string]interface{}),
		forLoops:    make(map[string]ForLoopInfo),
		pc:          0,
	}
}

// getCurrentStatement returns the current statement being executed
func (ctx *ExecutionContext) getCurrentStatement() string {
	if ctx.pc >= len(ctx.lineNumbers) {
		return ""
	}
	lineNum := ctx.lineNumbers[ctx.pc]
	return ctx.program[lineNum]
}

// getCurrentLineNumber returns the current line number being executed
func (ctx *ExecutionContext) getCurrentLineNumber() int {
	if ctx.pc >= len(ctx.lineNumbers) {
		return 0
	}
	return ctx.lineNumbers[ctx.pc]
}

// executeNextStatement executes the next statement in the program
func (fe *FileExecutor) executeNextStatement(ctx *ExecutionContext, debugMode bool) error {
	lineNum := ctx.getCurrentLineNumber()
	statement := ctx.getCurrentStatement()
	
	if debugMode {
		fe.output.WriteLine(fmt.Sprintf("Debug: Executing line %d: %s", lineNum, statement))
	}
	
	// Execute statement
	nextPC, err := fe.executeStatementWithControl(statement, ctx.variables, ctx.forLoops, ctx.pc, ctx.lineNumbers, ctx.program)
	if err != nil {
		return fmt.Errorf("error at line %d: %w", lineNum, err)
	}
	
	ctx.pc = nextPC
	return nil
}

// isEndStatement checks if a statement is an END statement
func (fe *FileExecutor) isEndStatement(statement string) bool {
	return strings.ToUpper(strings.TrimSpace(statement)) == "END"
}

// getSortedLineNumbers returns line numbers sorted in ascending order
func (fe *FileExecutor) getSortedLineNumbers(program map[int]string) []int {
	var lineNumbers []int
	for lineNum := range program {
		lineNumbers = append(lineNumbers, lineNum)
	}
	
	// Sort line numbers
	for i := 0; i < len(lineNumbers)-1; i++ {
		for j := i + 1; j < len(lineNumbers); j++ {
			if lineNumbers[i] > lineNumbers[j] {
				lineNumbers[i], lineNumbers[j] = lineNumbers[j], lineNumbers[i]
			}
		}
	}
	
	return lineNumbers
}

// executeStatementWithControl executes a single BASIC statement with control flow
func (fe *FileExecutor) executeStatementWithControl(statement string, variables map[string]interface{}, forLoops map[string]ForLoopInfo, pc int, lineNumbers []int, program map[int]string) (int, error) {
	statement = strings.TrimSpace(statement)
	upper := strings.ToUpper(statement)
	
	// Handle different statement types using a more structured approach
	switch {
	case strings.HasPrefix(upper, "PRINT"):
		err := fe.executePrint(statement, variables)
		return pc + 1, err
		
	case strings.HasPrefix(upper, "REM"):
		// Comment - do nothing
		return pc + 1, nil
		
	case strings.HasPrefix(upper, "FOR"):
		// FOR loop start - check this BEFORE assignment
		return fe.executeForStart(statement, variables, forLoops, pc)
		
	case fe.isAssignmentStatement(statement):
		// Assignment
		err := fe.executeAssignment(statement, variables)
		return pc + 1, err
		
	case strings.HasPrefix(upper, "NEXT"):
		// NEXT - handle loop continuation
		return fe.executeNext(statement, variables, forLoops, pc, lineNumbers)
		
	case strings.HasPrefix(upper, "IF"):
		// IF statement
		err := fe.executeIf(statement, variables)
		return pc + 1, err
		
	case upper == "END":
		// END statement
		return pc + 1, nil
		
	default:
		// Unknown statement - ignore for now
		return pc + 1, nil
	}
}

// isAssignmentStatement checks if a statement is an assignment
func (fe *FileExecutor) isAssignmentStatement(statement string) bool {
	return strings.Contains(statement, "=") && !strings.Contains(statement, "==")
}

// executeStatement executes a single BASIC statement (legacy method)
func (fe *FileExecutor) executeStatement(statement string, variables map[string]interface{}) error {
	statement = strings.TrimSpace(statement)
	upper := strings.ToUpper(statement)
	
	// Handle different statement types
	if strings.HasPrefix(upper, "PRINT") {
		return fe.executePrint(statement, variables)
	} else if strings.HasPrefix(upper, "REM") {
		// Comment - do nothing
		return nil
	} else if strings.Contains(statement, "=") && !strings.Contains(statement, "==") {
		// Assignment
		return fe.executeAssignment(statement, variables)
	} else if strings.HasPrefix(upper, "FOR") {
		// FOR loop - not supported in legacy mode
		return nil
	} else if strings.HasPrefix(upper, "NEXT") {
		// NEXT - not supported in legacy mode
		return nil
	} else if strings.HasPrefix(upper, "IF") {
		// IF statement - simplified simulation
		return fe.executeIf(statement, variables)
	} else if upper == "END" {
		// END statement
		return nil
	}
	
	return nil // Unknown statement - ignore for now
}

// setPrintOutputWriters sets the output writer for all PRINT statements in the program
func (fe *FileExecutor) setPrintOutputWriters(program *ast.Program) {
	for _, statement := range program.Lines {
		fe.setPrintOutputWriterForStatement(statement)
	}
}

// setPrintOutputWriterForStatement recursively sets output writers for PRINT statements
func (fe *FileExecutor) setPrintOutputWriterForStatement(statement ast.Statement) {
	switch stmt := statement.(type) {
	case *ast.PrintStatement:
		stmt.Output = fe.output
	case *ast.IfStatement:
		// Handle PRINT statements in IF-THEN clauses
		if stmt.ThenStatement != nil {
			fe.setPrintOutputWriterForStatement(stmt.ThenStatement)
		}
	// Add other statement types that might contain PRINT statements as needed
	}
}

// setInputOutputWriters sets the input/output writers for all INPUT statements in the program
func (fe *FileExecutor) setInputOutputWriters(program *ast.Program) {
	for _, statement := range program.Lines {
		fe.setInputOutputWriterForStatement(statement)
	}
}

// setInputOutputWriterForStatement recursively sets input/output writers for INPUT statements
func (fe *FileExecutor) setInputOutputWriterForStatement(statement ast.Statement) {
	switch stmt := statement.(type) {
	case *ast.InputStatement:
		stmt.Input = fe.input
		stmt.Output = fe.output
	case *ast.IfStatement:
		// Handle INPUT statements in IF-THEN clauses
		if stmt.ThenStatement != nil {
			fe.setInputOutputWriterForStatement(stmt.ThenStatement)
		}
	// Add other statement types that might contain INPUT statements as needed
	}
}