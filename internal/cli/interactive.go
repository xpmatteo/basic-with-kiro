package cli

import (
	"fmt"
	"strings"
)

// InteractiveMode handles the REPL functionality
type InteractiveMode struct {
	input     InputReader
	output    OutputWriter
	program   map[int]string // Store program lines as strings
	order     []int          // Track line order
	variables map[string]interface{} // Store variables
}

// NewInteractiveMode creates a new interactive mode instance
func NewInteractiveMode(input InputReader, output OutputWriter) *InteractiveMode {
	return &InteractiveMode{
		input:     input,
		output:    output,
		program:   make(map[int]string),
		order:     []int{},
		variables: make(map[string]interface{}),
	}
}

// Run starts the interactive REPL
func (im *InteractiveMode) Run() error {
	im.displayWelcomeMessage()

	for {
		im.output.WriteLine(ReadyPrompt)
		
		line, err := im.input.ReadLine()
		if err != nil {
			break // EOF or error, exit gracefully
		}
		
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Handle special commands
		handled, shouldExit := im.handleCommand(line)
		if shouldExit {
			break
		}
		if handled {
			continue // Command was handled, don't process as program line
		}
		
		// Try to parse as a numbered line
		if err := im.processLine(line); err != nil {
			im.displayError(err)
		}
	}
	
	return nil
}

// displayWelcomeMessage shows the welcome message and instructions
func (im *InteractiveMode) displayWelcomeMessage() {
	im.output.WriteLine(InteractiveModeHeader)
	im.output.WriteLine(InteractiveModeInstructions)
	im.output.WriteLine("")
}

// displayError displays an error message with consistent formatting
func (im *InteractiveMode) displayError(err error) {
	im.output.WriteLine(fmt.Sprintf("Error: %s", err.Error()))
}

// handleCommand handles special interactive commands
// Returns true if the command was handled (including exit)
func (im *InteractiveMode) handleCommand(line string) (bool, bool) {
	switch strings.ToUpper(line) {
	case "EXIT", "QUIT":
		im.output.WriteLine(GoodbyeMessage)
		return true, true // Command handled, should exit
	case "LIST":
		im.listProgram()
		return true, false // Command handled, continue running
	case "RUN":
		im.runProgram()
		return true, false // Command handled, continue running
	case "CLEAR":
		im.clearProgram()
		return true, false // Command handled, continue running
	default:
		return false, false // Not a command
	}
}

// processLine processes a line of input (either a program line or immediate command)
func (im *InteractiveMode) processLine(line string) error {
	// Check if line starts with a number (program line)
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return nil
	}
	
	// Try to parse first part as line number
	if lineNum, err := parseLineNumber(parts[0]); err == nil {
		// This is a program line
		if len(parts) == 1 {
			// Just a line number - delete the line
			im.deleteLine(lineNum)
		} else {
			// Line number with statement - add/replace line
			statement := strings.Join(parts[1:], " ")
			
			// Basic syntax validation
			if err := im.validateStatement(statement); err != nil {
				return err
			}
			
			im.addLine(lineNum, statement)
		}
		return nil
	}
	
	// Not a program line, treat as immediate command
	return fmt.Errorf("immediate commands not supported yet: %s", line)
}

// validateStatement performs basic syntax validation
func (im *InteractiveMode) validateStatement(statement string) error {
	// Check for unterminated strings
	quoteCount := strings.Count(statement, "\"")
	if quoteCount%2 != 0 {
		return fmt.Errorf("unterminated string")
	}
	
	return nil
}

// addLine adds or replaces a program line
func (im *InteractiveMode) addLine(lineNum int, statement string) {
	im.program[lineNum] = statement
	
	if !im.lineExists(lineNum) {
		im.insertLineInOrder(lineNum)
	}
}

// lineExists checks if a line number already exists in the program
func (im *InteractiveMode) lineExists(lineNum int) bool {
	for _, num := range im.order {
		if num == lineNum {
			return true
		}
	}
	return false
}

// insertLineInOrder inserts a line number in the correct sorted position
func (im *InteractiveMode) insertLineInOrder(lineNum int) {
	im.order = append(im.order, lineNum)
	
	// Keep order sorted using insertion sort for the last element
	for i := len(im.order) - 1; i > 0; i-- {
		if im.order[i] < im.order[i-1] {
			im.order[i], im.order[i-1] = im.order[i-1], im.order[i]
		} else {
			break
		}
	}
}

// deleteLine removes a program line
func (im *InteractiveMode) deleteLine(lineNum int) {
	delete(im.program, lineNum)
	im.removeLineFromOrder(lineNum)
}

// removeLineFromOrder removes a line number from the order slice
func (im *InteractiveMode) removeLineFromOrder(lineNum int) {
	for i, num := range im.order {
		if num == lineNum {
			im.order = append(im.order[:i], im.order[i+1:]...)
			break
		}
	}
}

// listProgram displays the current program
func (im *InteractiveMode) listProgram() {
	if len(im.program) == 0 {
		im.output.WriteLine(NoProgramLoadedMessage)
		return
	}
	
	for _, lineNum := range im.order {
		if statement, exists := im.program[lineNum]; exists {
			im.output.WriteLine(fmt.Sprintf("%d %s", lineNum, statement))
		}
	}
}

// runProgram executes the current program
func (im *InteractiveMode) runProgram() {
	if len(im.program) == 0 {
		im.output.WriteLine(NoProgramMessage)
		return
	}
	
	im.output.WriteLine(RunningProgramMessage)
	
	// Create a file executor to handle the execution logic
	fileExecutor := NewFileExecutor(im.output)
	
	// Execute the program using the same logic as file execution
	if err := fileExecutor.ExecuteProgram(im.program, false); err != nil {
		im.output.WriteLine(fmt.Sprintf("Runtime error: %s", err.Error()))
		return
	}
	
	im.output.WriteLine(ProgramCompletedMessage)
}



// clearProgram clears the current program
func (im *InteractiveMode) clearProgram() {
	im.program = make(map[int]string)
	im.order = []int{}
	im.variables = make(map[string]interface{})
	im.output.WriteLine(ProgramClearedMessage)
}