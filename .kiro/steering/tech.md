# Technology Stack

## Language & Runtime
- **Go 1.24.4** - Primary implementation language
- Standard Go toolchain for building and testing

## Dependencies
- `github.com/stretchr/testify v1.11.1` - Testing framework for unit tests
- No external runtime dependencies - pure Go implementation

## Architecture Pattern
- **Interpreter Architecture** with clear separation of concerns:
  - Lexical analysis (tokenization)
  - Syntax analysis (parsing to AST)
  - Runtime execution with environment management
  - CLI interface layer

## Project Structure
- `cmd/` - Application entry points
- `internal/` - Private implementation packages
- Interface-driven design for testability and modularity

## Common Commands

### Building
```bash
go build -o basic-interpreter cmd/basic-interpreter/main.go
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/lexer
go test ./internal/parser
go test ./internal/ast
```

### Running
```bash
# Interactive mode
./basic-interpreter

# Execute BASIC file
./basic-interpreter program.bas

# Debug mode
./basic-interpreter -debug program.bas
```

### Development
```bash
# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Run linter (if golangci-lint installed)
golangci-lint run
```