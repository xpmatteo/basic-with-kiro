package cli

import "fmt"

// parseLineNumber parses a string as a line number
func parseLineNumber(s string) (int, error) {
	var lineNum int
	num, err := fmt.Sscanf(s, "%d", &lineNum)
	if err != nil || num != 1 {
		return 0, fmt.Errorf("invalid line number: %s", s)
	}
	
	if err := validateLineNumber(lineNum); err != nil {
		return 0, err
	}
	
	return lineNum, nil
}

// validateLineNumber validates that a line number is within acceptable range
func validateLineNumber(lineNum int) error {
	if lineNum < MinLineNumber || lineNum > MaxLineNumber {
		return fmt.Errorf("line number out of range: %d (must be between %d and %d)", 
			lineNum, MinLineNumber, MaxLineNumber)
	}
	return nil
}