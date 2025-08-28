package cli

import (
	"errors"
	"fmt"
	"strings"
)

// Config holds the parsed command line configuration
type Config struct {
	DebugMode   bool
	Interactive bool
	InputFile   string
}

// CLI handles command line argument parsing
type CLI struct {
	version string
}

// NewCLI creates a new CLI instance
func NewCLI() *CLI {
	return &CLI{
		version: InterpreterVersion,
	}
}

// ParseArgs parses command line arguments and returns a Config
func (c *CLI) ParseArgs(args []string) (*Config, error) {
	if len(args) == 0 {
		return nil, errors.New("no arguments provided")
	}

	config := &Config{}
	
	// Skip program name (first argument)
	args = args[1:]
	
	var fileArgs []string
	
	for i := 0; i < len(args); i++ {
		arg := args[i]
		
		switch arg {
		case "-h", "--help":
			return nil, errors.New("help requested")
		case "-d", "--debug":
			config.DebugMode = true
		default:
			if strings.HasPrefix(arg, "-") {
				return nil, fmt.Errorf("unknown flag: %s", arg)
			}
			fileArgs = append(fileArgs, arg)
		}
	}
	
	// Handle file arguments
	if len(fileArgs) == 0 {
		config.Interactive = true
	} else if len(fileArgs) == 1 {
		config.InputFile = fileArgs[0]
		config.Interactive = false
	} else {
		return nil, errors.New("too many file arguments")
	}
	
	// Validate debug mode requirements
	if config.DebugMode && config.Interactive {
		return nil, errors.New("debug mode requires a file")
	}
	
	return config, nil
}

// GetHelpMessage returns the help message
func (c *CLI) GetHelpMessage() string {
	return `BASIC Interpreter

Usage:
  basic-interpreter [options] [file]

Options:
  -d, --debug    Enable debug mode (shows each line before execution)
  -h, --help     Show this help message

Arguments:
  file           BASIC program file to execute

Examples:
  basic-interpreter                    # Interactive mode
  basic-interpreter program.bas       # Execute file
  basic-interpreter -d program.bas    # Execute file with debug output
`
}

// GetVersionInfo returns version information
func (c *CLI) GetVersionInfo() string {
	return fmt.Sprintf("BASIC Interpreter version %s", c.version)
}

// Validate checks if the configuration is valid
func (config *Config) Validate() error {
	if config.InputFile != "" && config.Interactive {
		return errors.New("cannot specify both file and interactive mode")
	}
	
	if config.DebugMode && config.Interactive {
		return errors.New("debug mode is only available for file execution")
	}
	
	return nil
}