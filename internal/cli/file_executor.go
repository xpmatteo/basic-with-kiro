package cli

import (
	"fmt"
	"os"
	"strings"
)

// FileExecutor handles file-based program execution
type FileExecutor struct {
	input  InputReader
	output OutputWriter
}

// NewFileExecutor creates a new file executor instance
func NewFileExecutor(input InputReader, output OutputWriter) *FileExecutor {
	return &FileExecutor{
		input:  input,
		output: output,
	}
}

// ExecuteFile loads and executes a BASIC program from a file
func (fe *FileExecutor) ExecuteFile(filename string, debugMode bool) error {
	// Read file content
	content, err := fe.readFile(filename)
	if err != nil {
		return fe.wrapFileError("failed to read file", filename, err)
	}
	
	// Parse program
	program, err := fe.parseProgram(content)
	if err != nil {
		return fe.wrapFileError("syntax error in", filename, err)
	}
	
	// Execute program
	if err := fe.ExecuteProgram(program, debugMode); err != nil {
		return fe.wrapFileError("runtime error in", filename, err)
	}
	
	return nil
}

// wrapFileError wraps an error with file context information
func (fe *FileExecutor) wrapFileError(prefix, filename string, err error) error {
	return fmt.Errorf("%s %s: %w", prefix, filename, err)
}

// readFile reads the content of a file
func (fe *FileExecutor) readFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// parseProgram parses a BASIC program from source code
func (fe *FileExecutor) parseProgram(content string) (map[int]string, error) {
	program := make(map[int]string)
	lines := strings.Split(content, "\n")
	
	for lineIdx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue // Skip empty lines
		}
		
		// Parse line number and statement
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		
		// Try to parse line number
		lineNum, err := parseLineNumber(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid line number at line %d: %s", lineIdx+1, parts[0])
		}
		
		// Get statement (everything after line number)
		if len(parts) > 1 {
			statement := strings.Join(parts[1:], " ")
			
			// Basic syntax validation
			if err := fe.validateStatement(statement); err != nil {
				return nil, fmt.Errorf("syntax error at line %d: %w", lineNum, err)
			}
			
			program[lineNum] = statement
		}
	}
	
	return program, nil
}

// validateStatement performs basic syntax validation
func (fe *FileExecutor) validateStatement(statement string) error {
	// Check for unterminated strings
	quoteCount := strings.Count(statement, "\"")
	if quoteCount%2 != 0 {
		return fmt.Errorf("unterminated string")
	}
	
	return nil
}