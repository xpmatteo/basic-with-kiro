# Project Structure

## Directory Organization

```
basic-interpreter/
├── cmd/basic-interpreter/     # Main application entry point
├── internal/                  # Private implementation packages
│   ├── ast/                   # Abstract Syntax Tree definitions
│   ├── cli/                   # Command-line interface
│   ├── errors/                # Error types and handling
│   ├── interpreter/           # Interpreter execution engine
│   ├── lexer/                 # Lexical analysis (tokenization)
│   ├── parser/                # Syntax analysis (parsing)
│   └── runtime/               # Runtime environment and values
├── go.mod                     # Go module definition
└── *.bas                      # BASIC program test files
```

## Package Responsibilities

### `cmd/basic-interpreter/`
- Application entry point and CLI argument parsing
- Orchestrates CLI components for interactive/file execution modes

### `internal/ast/`
- AST node definitions for statements and expressions
- Built-in function implementations (math, string functions)
- Statement execution logic and expression evaluation
- Helper functions for AST construction

### `internal/cli/`
- CLI configuration and argument parsing
- Interactive mode implementation
- File execution engine
- I/O interfaces and implementations

### `internal/lexer/`
- Token definitions and lexical analysis
- Converts BASIC source code into token stream
- Handles keywords, operators, literals, and line numbers

### `internal/parser/`
- Recursive descent parser implementation
- Converts token stream into AST
- Handles operator precedence and syntax validation

### `internal/runtime/`
- Runtime environment and variable storage
- Value types (numeric/string) and operations
- FOR loop state management
- Program counter and execution context

### `internal/errors/`
- Custom error types for interpreter errors
- Error formatting and context information

## Coding Conventions

### Interface Design
- Use interfaces for testability (OutputWriter, InputReader, etc.)
- Dependency injection for I/O operations
- Clear separation between parsing and execution

### Error Handling
- Wrap errors with context using `fmt.Errorf`
- Validate inputs at package boundaries
- Return meaningful error messages for BASIC programmers

### Testing
- Comprehensive unit tests for each package
- Use testify for assertions and test structure
- Mock interfaces for isolated testing
- Integration tests for end-to-end scenarios

### Naming Conventions
- Go standard naming (PascalCase for exported, camelCase for private)
- BASIC keywords in UPPERCASE in comments/strings
- Clear, descriptive function and variable names
- Package names reflect their primary responsibility