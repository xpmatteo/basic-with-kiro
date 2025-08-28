package cli

const (
	// Version information
	InterpreterVersion = "1.0.0"
	
	// Line number constraints
	MinLineNumber = 1
	MaxLineNumber = 99999
	
	// Interactive mode messages
	InteractiveModeHeader = "BASIC Interpreter - Interactive Mode"
	InteractiveModeInstructions = "Type EXIT to quit, LIST to show program, RUN to execute, CLEAR to clear program"
	ReadyPrompt = "READY"
	GoodbyeMessage = "Goodbye!"
	
	// Program execution messages
	RunningProgramMessage = "Running program..."
	ProgramCompletedMessage = "Program completed"
	NoProgramMessage = "No program to run"
	NoProgramLoadedMessage = "No program loaded"
	ProgramClearedMessage = "Program cleared"
)