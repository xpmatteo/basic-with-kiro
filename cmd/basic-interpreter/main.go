package main

import (
	"fmt"
	"os"
	"strings"

	"basic-interpreter/internal/cli"
)

func main() {
	// Create CLI instance
	cliInstance := cli.NewCLI()
	
	// Parse command line arguments
	config, err := cliInstance.ParseArgs(os.Args)
	if err != nil {
		// Handle help request
		if strings.Contains(err.Error(), "help") {
			fmt.Print(cliInstance.GetHelpMessage())
			os.Exit(0)
		}
		
		// Handle other errors
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		fmt.Print(cliInstance.GetHelpMessage())
		os.Exit(1)
	}
	
	// Validate configuration
	if err := config.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	
	// Create I/O interfaces
	input := cli.NewStdInputReader()
	output := cli.NewStdOutputWriter()
	
	// Execute based on mode
	if config.Interactive {
		// Interactive mode
		interactive := cli.NewInteractiveMode(input, output)
		if err := interactive.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error in interactive mode: %s\n", err.Error())
			os.Exit(1)
		}
	} else {
		// File execution mode
		fileExecutor := cli.NewFileExecutor(input, output)
		if err := fileExecutor.ExecuteFile(config.InputFile, config.DebugMode); err != nil {
			fmt.Fprintf(os.Stderr, "Error executing file: %s\n", err.Error())
			os.Exit(1)
		}
	}
}