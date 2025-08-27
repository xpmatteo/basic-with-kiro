# BASIC Interpreter

A BASIC interpreter written in Go that supports core BASIC language constructs including variables, control flow, arithmetic operations, and basic I/O operations.

## Project Structure

```
basic-interpreter/
├── cmd/basic-interpreter/     # Main application entry point
├── internal/
│   ├── ast/                   # Abstract Syntax Tree definitions
│   ├── errors/                # Error types and handling
│   ├── interpreter/           # Interpreter execution engine
│   ├── lexer/                 # Lexical analysis (tokenization)
│   ├── parser/                # Syntax analysis (parsing)
│   └── runtime/               # Runtime environment and values
├── go.mod                     # Go module definition
└── README.md                  # This file
```

## Building

```bash
go build -o basic-interpreter cmd/basic-interpreter/main.go
```

## Running Tests

```bash
go test ./...
```

## Current Status

✅ Project structure and core interfaces defined
🚧 Implementation in progress

## Planned Features

- Line-numbered BASIC programs
- Variables (numeric and string)
- Arithmetic operations (+, -, *, /, ^)
- Control flow (GOTO, IF-THEN, FOR-NEXT)
- I/O operations (PRINT, INPUT)
- Built-in functions (ABS, INT, RND, LEN, etc.)
- Command-line interface with file execution and interactive modes
- Error handling and debugging support