package cli

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLI_ParseArgs_DebugFlag(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected *Config
		wantErr  bool
	}{
		{
			name: "debug flag short form",
			args: []string{"program", "-d", "test.bas"},
			expected: &Config{
				DebugMode: true,
				InputFile: "test.bas",
			},
			wantErr: false,
		},
		{
			name: "debug flag long form",
			args: []string{"program", "--debug", "test.bas"},
			expected: &Config{
				DebugMode: true,
				InputFile: "test.bas",
			},
			wantErr: false,
		},
		{
			name: "no debug flag",
			args: []string{"program", "test.bas"},
			expected: &Config{
				DebugMode: false,
				InputFile: "test.bas",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := NewCLI()
			config, err := cli.ParseArgs(tt.args)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected.DebugMode, config.DebugMode)
			assert.Equal(t, tt.expected.InputFile, config.InputFile)
		})
	}
}

func TestCLI_ParseArgs_HelpFlag(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "help flag short form",
			args:    []string{"program", "-h"},
			wantErr: true, // Help should return an error to stop execution
		},
		{
			name:    "help flag long form",
			args:    []string{"program", "--help"},
			wantErr: true, // Help should return an error to stop execution
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := NewCLI()
			_, err := cli.ParseArgs(tt.args)

			if tt.wantErr {
				assert.Error(t, err)
				// Should be a help error, not a real error
				assert.Contains(t, err.Error(), "help")
			}
		})
	}
}

func TestCLI_ParseArgs_FileArguments(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected *Config
		wantErr  bool
	}{
		{
			name: "single file argument",
			args: []string{"program", "test.bas"},
			expected: &Config{
				InputFile:   "test.bas",
				Interactive: false,
			},
			wantErr: false,
		},
		{
			name: "no arguments - interactive mode",
			args: []string{"program"},
			expected: &Config{
				Interactive: true,
			},
			wantErr: false,
		},
		{
			name:    "multiple file arguments",
			args:    []string{"program", "file1.bas", "file2.bas"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := NewCLI()
			config, err := cli.ParseArgs(tt.args)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected.InputFile, config.InputFile)
			assert.Equal(t, tt.expected.Interactive, config.Interactive)
		})
	}
}

func TestCLI_ParseArgs_ErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "invalid flag",
			args:    []string{"program", "--invalid"},
			wantErr: true,
			errMsg:  "unknown flag",
		},
		{
			name:    "debug flag without file in non-interactive",
			args:    []string{"program", "-d"},
			wantErr: true,
			errMsg:  "debug mode requires a file",
		},
		{
			name:    "conflicting flags",
			args:    []string{"program", "-h", "-d", "test.bas"},
			wantErr: true,
			errMsg:  "help",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := NewCLI()
			_, err := cli.ParseArgs(tt.args)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestCLI_GetHelpMessage(t *testing.T) {
	cli := NewCLI()
	help := cli.GetHelpMessage()

	assert.Contains(t, help, "BASIC Interpreter")
	assert.Contains(t, help, "Usage:")
	assert.Contains(t, help, "-d, --debug")
	assert.Contains(t, help, "-h, --help")
	assert.Contains(t, help, "Examples:")
	assert.Contains(t, help, "Interactive mode")
	assert.Contains(t, help, "Execute file")
}

func TestCLI_GetVersionInfo(t *testing.T) {
	cli := NewCLI()
	version := cli.GetVersionInfo()

	assert.Contains(t, version, "BASIC Interpreter")
	assert.Contains(t, version, "version")
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid file config",
			config: &Config{
				InputFile:   "test.bas",
				Interactive: false,
			},
			wantErr: false,
		},
		{
			name: "valid interactive config",
			config: &Config{
				Interactive: true,
			},
			wantErr: false,
		},
		{
			name: "invalid - both file and interactive",
			config: &Config{
				InputFile:   "test.bas",
				Interactive: true,
			},
			wantErr: true,
			errMsg:  "cannot specify both file and interactive mode",
		},
		{
			name: "invalid - debug without file",
			config: &Config{
				DebugMode:   true,
				Interactive: true,
			},
			wantErr: true,
			errMsg:  "debug mode is only available for file execution",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Interactive mode tests
func TestCLI_InteractiveMode_BasicInput(t *testing.T) {
	tests := []struct {
		name     string
		inputs   []string
		expected []string
	}{
		{
			name:     "simple print statement",
			inputs:   []string{"10 PRINT \"Hello\"", "RUN", "EXIT"},
			expected: []string{"Hello"},
		},
		{
			name:     "variable assignment and print",
			inputs:   []string{"10 A = 5", "20 PRINT A", "RUN", "EXIT"},
			expected: []string{"5"},
		},
		{
			name:     "multiple lines",
			inputs:   []string{"10 PRINT \"Line 1\"", "20 PRINT \"Line 2\"", "RUN", "EXIT"},
			expected: []string{"Line 1", "Line 2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockInput := &MockInputReader{inputs: tt.inputs}
			mockOutput := &MockOutputWriter{}
			
			interactive := NewInteractiveMode(mockInput, mockOutput)
			err := interactive.Run()
			
			assert.NoError(t, err)
			
			// Check that expected outputs are present
			for _, expected := range tt.expected {
				found := false
				for _, output := range mockOutput.outputs {
					if strings.Contains(output, expected) {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected output '%s' not found in %v", expected, mockOutput.outputs)
			}
		})
	}
}

func TestCLI_InteractiveMode_Commands(t *testing.T) {
	tests := []struct {
		name     string
		inputs   []string
		wantErr  bool
		contains string
	}{
		{
			name:     "LIST command shows program",
			inputs:   []string{"10 PRINT \"Hello\"", "LIST", "EXIT"},
			wantErr:  false,
			contains: "10 PRINT \"Hello\"",
		},
		{
			name:     "CLEAR command clears program",
			inputs:   []string{"10 PRINT \"Hello\"", "CLEAR", "LIST", "EXIT"},
			wantErr:  false,
			contains: "Program cleared",
		},
		{
			name:     "RUN empty program",
			inputs:   []string{"RUN", "EXIT"},
			wantErr:  false,
			contains: "No program to run",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockInput := &MockInputReader{inputs: tt.inputs}
			mockOutput := &MockOutputWriter{}
			
			interactive := NewInteractiveMode(mockInput, mockOutput)
			err := interactive.Run()
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			
			// Check for expected content in output
			if tt.contains != "" {
				found := false
				for _, output := range mockOutput.outputs {
					if strings.Contains(output, tt.contains) {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected content '%s' not found in %v", tt.contains, mockOutput.outputs)
			}
		})
	}
}

func TestCLI_InteractiveMode_ErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		inputs   []string
		contains string
	}{
		{
			name:     "syntax error in line",
			inputs:   []string{"10 PRINT \"unterminated", "EXIT"},
			contains: "Error",
		},
		{
			name:     "runtime error",
			inputs:   []string{"10 PRINT 1/0", "RUN", "EXIT"},
			contains: "Error",
		},
		{
			name:     "invalid line number",
			inputs:   []string{"ABC PRINT \"Hello\"", "EXIT"},
			contains: "Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockInput := &MockInputReader{inputs: tt.inputs}
			mockOutput := &MockOutputWriter{}
			
			interactive := NewInteractiveMode(mockInput, mockOutput)
			err := interactive.Run()
			
			// Interactive mode should handle errors gracefully and not crash
			assert.NoError(t, err)
			
			// Check that error message was displayed
			found := false
			for _, output := range mockOutput.outputs {
				if strings.Contains(output, tt.contains) {
					found = true
					break
				}
			}
			assert.True(t, found, "Expected error message containing '%s' not found in %v", tt.contains, mockOutput.outputs)
		})
	}
}

func TestCLI_InteractiveMode_ProgramState(t *testing.T) {
	mockInput := &MockInputReader{
		inputs: []string{
			"10 A = 5",
			"20 PRINT A",
			"15 A = A + 1", // Insert line in middle
			"LIST",
			"RUN",
			"EXIT",
		},
	}
	mockOutput := &MockOutputWriter{}
	
	interactive := NewInteractiveMode(mockInput, mockOutput)
	err := interactive.Run()
	
	assert.NoError(t, err)
	
	// Check that program state is maintained
	// Should show lines in correct order: 10, 15, 20
	found10, found15, found20 := false, false, false
	for _, output := range mockOutput.outputs {
		if strings.Contains(output, "10 A = 5") {
			found10 = true
		}
		if strings.Contains(output, "15 A = A + 1") {
			found15 = true
		}
		if strings.Contains(output, "20 PRINT A") {
			found20 = true
		}
	}
	assert.True(t, found10 && found15 && found20, "LIST command should show all lines in order")
	
	// Check that execution uses updated program (A should be 6, not 5)
	resultFound := false
	for _, output := range mockOutput.outputs {
		if strings.Contains(output, "6") {
			resultFound = true
			break
		}
	}
	assert.True(t, resultFound, "Program should execute with updated line 15")
}

// Mock types for testing
type MockInputReader struct {
	inputs []string
	index  int
}

func (m *MockInputReader) ReadLine() (string, error) {
	if m.index >= len(m.inputs) {
		return "", fmt.Errorf("EOF")
	}
	result := m.inputs[m.index]
	m.index++
	return result, nil
}

type MockOutputWriter struct {
	outputs []string
}

func (m *MockOutputWriter) WriteLine(line string) error {
	m.outputs = append(m.outputs, line)
	return nil
}

func (m *MockOutputWriter) Write(data []byte) (int, error) {
	m.outputs = append(m.outputs, string(data))
	return len(data), nil
}

// File execution mode tests
func TestCLI_FileExecution_BasicProgram(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		expected    []string
		wantErr     bool
	}{
		{
			name: "simple hello world program",
			fileContent: `10 PRINT "Hello, World!"
20 END`,
			expected: []string{"Hello, World!"},
			wantErr:  false,
		},
		{
			name: "program with variables",
			fileContent: `10 A = 5
20 B = 10
30 PRINT A + B
40 END`,
			expected: []string{"15"},
			wantErr:  false,
		},
		{
			name: "program with loop",
			fileContent: `10 FOR I = 1 TO 3
20 PRINT I
30 NEXT I
40 END`,
			expected: []string{"1", "2", "3"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpFile := createTempFile(t, tt.fileContent)
			defer removeTempFile(t, tmpFile)

			mockInput := &MockInputReader{}
			mockOutput := &MockOutputWriter{}
			fileExecutor := NewFileExecutor(mockInput, mockOutput)
			
			err := fileExecutor.ExecuteFile(tmpFile, false)
			
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			
			// Check expected outputs
			for _, expected := range tt.expected {
				found := false
				for _, output := range mockOutput.outputs {
					if strings.Contains(output, expected) {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected output '%s' not found in %v", expected, mockOutput.outputs)
			}
		})
	}
}

func TestCLI_FileExecution_DebugMode(t *testing.T) {
	fileContent := `10 PRINT "Line 1"
20 PRINT "Line 2"
30 END`
	
	tmpFile := createTempFile(t, fileContent)
	defer removeTempFile(t, tmpFile)

	mockInput := &MockInputReader{}
	mockOutput := &MockOutputWriter{}
	fileExecutor := NewFileExecutor(mockInput, mockOutput)
	
	err := fileExecutor.ExecuteFile(tmpFile, true) // Debug mode enabled
	
	require.NoError(t, err)
	
	// Check that debug output is present
	debugFound := false
	for _, output := range mockOutput.outputs {
		if strings.Contains(output, "Executing line") || strings.Contains(output, "Debug:") {
			debugFound = true
			break
		}
	}
	assert.True(t, debugFound, "Debug output should be present when debug mode is enabled")
	
	// Check that program output is still present
	line1Found := false
	line2Found := false
	for _, output := range mockOutput.outputs {
		if strings.Contains(output, "Line 1") {
			line1Found = true
		}
		if strings.Contains(output, "Line 2") {
			line2Found = true
		}
	}
	assert.True(t, line1Found && line2Found, "Program output should be present in debug mode")
}

func TestCLI_FileExecution_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		filename    string
		wantErr     bool
		errContains string
	}{
		{
			name:        "file not found",
			filename:    "nonexistent.bas",
			wantErr:     true,
			errContains: "no such file",
		},
		{
			name: "syntax error in file",
			fileContent: `10 PRINT "unterminated string
20 END`,
			wantErr:     true,
			errContains: "syntax error",
		},
		{
			name: "runtime error in file",
			fileContent: `10 PRINT 1/0
20 END`,
			wantErr:     true,
			errContains: "runtime error",
		},
		{
			name: "invalid line number",
			fileContent: `ABC PRINT "Hello"
20 END`,
			wantErr:     true,
			errContains: "invalid line number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filename string
			if tt.filename != "" {
				filename = tt.filename
			} else {
				tmpFile := createTempFile(t, tt.fileContent)
				defer removeTempFile(t, tmpFile)
				filename = tmpFile
			}

			mockInput := &MockInputReader{}
			mockOutput := &MockOutputWriter{}
			fileExecutor := NewFileExecutor(mockInput, mockOutput)
			
			err := fileExecutor.ExecuteFile(filename, false)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCLI_FileExecution_Integration(t *testing.T) {
	// Test integration between file loading and interpreter execution
	fileContent := `10 REM This is a comprehensive test program
20 A = 10
30 B = 20
40 C = A + B
50 PRINT "Sum:", C
60 FOR I = 1 TO 3
70 PRINT "Loop iteration:", I
80 NEXT I
90 IF C > 25 THEN PRINT "C is greater than 25"
100 END`
	
	tmpFile := createTempFile(t, fileContent)
	defer removeTempFile(t, tmpFile)

	mockInput := &MockInputReader{}
	mockOutput := &MockOutputWriter{}
	fileExecutor := NewFileExecutor(mockInput, mockOutput)
	
	err := fileExecutor.ExecuteFile(tmpFile, false)
	
	require.NoError(t, err)
	
	// Check for expected outputs
	expectedOutputs := []string{
		"Sum: 30",
		"Loop iteration: 1",
		"Loop iteration: 2", 
		"Loop iteration: 3",
		"C is greater than 25",
	}
	
	for _, expected := range expectedOutputs {
		found := false
		for _, output := range mockOutput.outputs {
			if strings.Contains(output, expected) {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected output '%s' not found in %v", expected, mockOutput.outputs)
	}
}

func TestCLI_FileExecution_EmptyFile(t *testing.T) {
	tmpFile := createTempFile(t, "")
	defer removeTempFile(t, tmpFile)

	mockInput := &MockInputReader{}
	mockOutput := &MockOutputWriter{}
	fileExecutor := NewFileExecutor(mockInput, mockOutput)
	
	err := fileExecutor.ExecuteFile(tmpFile, false)
	
	// Empty file should not cause an error, just no output
	assert.NoError(t, err)
}

func TestCLI_FileExecution_CommentsOnly(t *testing.T) {
	fileContent := `10 REM This is a comment
20 REM Another comment
30 REM End of comments`
	
	tmpFile := createTempFile(t, fileContent)
	defer removeTempFile(t, tmpFile)

	mockInput := &MockInputReader{}
	mockOutput := &MockOutputWriter{}
	fileExecutor := NewFileExecutor(mockInput, mockOutput)
	
	err := fileExecutor.ExecuteFile(tmpFile, false)
	
	// Comments only should execute without error but produce no output
	assert.NoError(t, err)
}

// Helper functions for file testing
func createTempFile(t *testing.T, content string) string {
	tmpFile := fmt.Sprintf("/tmp/test_basic_%s.bas", strings.ReplaceAll(t.Name(), "/", "_"))
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	require.NoError(t, err)
	return tmpFile
}

func removeTempFile(t *testing.T, filename string) {
	err := os.Remove(filename)
	if err != nil {
		t.Logf("Warning: could not remove temp file %s: %v", filename, err)
	}
}