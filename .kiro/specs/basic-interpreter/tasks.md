# Implementation Plan

- [x] 1. Set up project structure and core interfaces
  - Create Go module and directory structure for lexer, parser, interpreter, and CLI components
  - Define core interfaces for Token, Statement, Expression, and Value types
  - Create basic error types and handling structures
  - Verify project compiles successfully and all interfaces are properly defined
  - _Requirements: 6.1, 6.2_

- [x] 2. Implement lexer/tokenizer using TDD
- [x] 2.1 Write lexer tests first
  - Write failing tests for TokenType enumeration and basic token creation
  - Write tests for tokenizing numbers, strings, identifiers, keywords, and operators
  - Write tests for error cases: invalid characters, unterminated strings
  - Write tests for position tracking (line and column numbers)
  - _Requirements: 1.3, 6.2, TDD Requirements 1.1_

- [x] 2.2 Implement minimal lexer to pass tests
  - Create Lexer struct with NextToken() and HasMoreTokens() methods
  - Implement just enough functionality to make the failing tests pass
  - Add support for tokenizing basic elements as defined by tests
  - Verify all lexer tests pass
  - _Requirements: 1.3, 6.2, TDD Requirements 1.1_

- [x] 2.3 Refactor lexer implementation
  - Improve lexer code structure while keeping tests green
  - Extract common tokenization patterns into helper methods
  - Optimize performance and readability
  - Ensure all tests still pass after refactoring
  - _Requirements: 1.3, 6.2, TDD Requirements 1.1_

- [x] 3. Implement value system and variables using TDD
- [x] 3.1 Write value system tests first
  - Write failing tests for Value struct supporting numeric and string types
  - Write tests for type conversion methods between numeric and string values
  - Write tests for edge cases: empty strings, zero values, invalid conversions
  - Write tests for value comparison and arithmetic operations
  - _Requirements: 2.1, 2.3, 2.4, TDD Requirements 2.3_

- [x] 3.2 Implement minimal value system to pass tests
  - Create Value struct with minimal functionality to pass tests
  - Implement basic type conversion methods as defined by tests
  - Add value operations (comparison, arithmetic) to satisfy test requirements
  - Verify all value system tests pass
  - _Requirements: 2.1, 2.3, 2.4, TDD Requirements 2.3_

- [x] 3.3 Write environment tests first
  - Write failing tests for Environment struct with variable storage
  - Write tests for case-insensitive variable lookup and assignment
  - Write tests for uninitialized variable behavior (default values)
  - Write tests for variable scoping and state management
  - _Requirements: 2.1, 2.4, TDD Requirements 2.4_

- [x] 3.4 Implement environment to pass tests
  - Create Environment struct with variable storage functionality
  - Implement case-insensitive lookup and proper default value initialization
  - Add state management features as required by tests
  - Verify all environment tests pass
  - _Requirements: 2.1, 2.4, TDD Requirements 2.4_

- [x] 3.5 Refactor value and environment code
  - Improve code structure while keeping all tests green
  - Extract common patterns and optimize performance
  - Ensure clean separation between value operations and storage
  - Verify all tests still pass after refactoring
  - _Requirements: 2.1, 2.3, 2.4_

- [x] 4. Implement basic expressions and arithmetic using TDD
- [x] 4.1 Write expression tests first
  - Write failing tests for Expression interface and basic expression types
  - Write tests for binary operators (+, -, *, /, ^) with proper precedence
  - Write tests for variable reference expressions and literal expressions
  - Write tests for complex nested expressions and parentheses handling
  - _Requirements: 2.2, TDD Requirements 2.2_

- [x] 4.2 Implement minimal expression system to pass tests
  - Create Expression interface and basic expression implementations
  - Implement arithmetic expression evaluation with correct operator precedence
  - Add variable reference and literal expression support
  - Verify all expression tests pass
  - _Requirements: 2.2, TDD Requirements 2.2_

- [x] 4.3 Refactor expression implementation
  - Improve expression evaluation code while keeping tests green
  - Extract common evaluation patterns and optimize precedence handling
  - Ensure clean separation between different expression types
  - Verify all tests still pass after refactoring
  - _Requirements: 2.2_

- [x] 5. Implement core statements using TDD
- [x] 5.1 Write assignment statement tests first
  - Write failing tests for AssignmentStatement struct and Execute method
  - Write tests for both numeric and string variable assignments
  - Write tests for assignment with complex expressions on right-hand side
  - Write tests for error cases: invalid variable names, type mismatches
  - _Requirements: 2.1, 2.3, TDD Requirements 2.1_

- [x] 5.2 Implement assignment statement to pass tests
  - Create AssignmentStatement struct with minimal Execute method
  - Add support for variable assignments as defined by tests
  - Implement error handling for invalid assignments
  - Verify all assignment statement tests pass
  - _Requirements: 2.1, 2.3, TDD Requirements 2.1_

- [x] 5.3 Write PRINT statement tests first
  - Write failing tests for PrintStatement struct with mocked output
  - Write tests for multiple expressions and different separators
  - Write tests for numeric and string formatting
  - Write tests for edge cases: empty print, special characters
  - _Requirements: 4.1, 4.3, TDD Requirements 4.1_

- [x] 5.4 Implement PRINT statement to pass tests
  - Create PrintStatement struct with output formatting functionality
  - Implement expression evaluation and output generation as required by tests
  - Add proper separator handling and value formatting
  - Verify all PRINT statement tests pass
  - _Requirements: 4.1, 4.3, TDD Requirements 4.1_

- [x] 5.5 Write INPUT statement tests first
  - Write failing tests for InputStatement struct with mocked input
  - Write tests for input validation and type conversion
  - Write tests for prompt display and variable assignment
  - Write tests for error cases: invalid input, conversion failures
  - _Requirements: 4.2, 4.4, TDD Requirements 4.2_

- [x] 5.6 Implement INPUT statement to pass tests
  - Create InputStatement struct with input reading functionality
  - Implement input validation and type conversion as defined by tests
  - Add proper error handling for invalid input scenarios
  - Verify all INPUT statement tests pass
  - _Requirements: 4.2, 4.4, TDD Requirements 4.2_

- [x] 5.7 Refactor core statements
  - Improve statement implementations while keeping all tests green
  - Extract common patterns and optimize code structure
  - Ensure clean separation between statement logic and I/O operations
  - Verify all tests still pass after refactoring
  - _Requirements: 2.1, 2.3, 4.1, 4.2, 4.3, 4.4_

- [x] 6. Implement control flow statements using TDD
- [x] 6.1 Write GOTO statement tests first
  - Write failing tests for GotoStatement struct and program counter modification
  - Write tests for validation of target line number existence
  - Write tests for error cases: invalid line numbers, missing targets
  - Write tests for forward and backward jumps
  - _Requirements: 3.1, 3.4, TDD Requirements 3.1_

- [x] 6.2 Implement GOTO statement to pass tests
  - Create GotoStatement struct with program counter modification
  - Implement line number validation as defined by tests
  - Add proper error handling for invalid jump targets
  - Verify all GOTO statement tests pass
  - _Requirements: 3.1, 3.4, TDD Requirements 3.1_

- [x] 6.3 Write IF-THEN statement tests first
  - Write failing tests for IfStatement struct with condition evaluation
  - Write tests for all comparison operators (=, <, >, <=, >=, <>)
  - Write tests for conditional execution logic and branching
  - Write tests for edge cases: type mismatches, complex conditions
  - _Requirements: 3.2, TDD Requirements 3.2_

- [x] 6.4 Implement IF-THEN statement to pass tests
  - Create IfStatement struct with condition evaluation functionality
  - Implement comparison operators and conditional execution as required by tests
  - Add proper type handling for condition evaluation
  - Verify all IF-THEN statement tests pass
  - _Requirements: 3.2, TDD Requirements 3.2_

- [x] 6.5 Write FOR-NEXT loop tests first
  - Write failing tests for ForStatement and NextStatement structs
  - Write tests for ForLoopState tracking and nested loop handling
  - Write tests for loop variable increment and termination logic
  - Write tests for error cases: missing NEXT, invalid loop bounds
  - _Requirements: 3.3, 3.5, TDD Requirements 3.3_

- [x] 6.6 Implement FOR-NEXT loops to pass tests
  - Create ForStatement and NextStatement structs with loop logic
  - Implement ForLoopState tracking in Environment for nested loops
  - Add loop variable management and termination as defined by tests
  - Verify all FOR-NEXT loop tests pass
  - _Requirements: 3.3, 3.5, TDD Requirements 3.3_

- [x] 6.7 Refactor control flow statements
  - Improve control flow implementations while keeping all tests green
  - Extract common patterns and optimize program counter management
  - Ensure clean separation between control flow logic and execution state
  - Verify all tests still pass after refactoring
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5_

- [x] 7. Implement built-in functions using TDD
- [x] 7.1 Write function call system tests first
  - Write failing tests for FunctionCallExpression struct and BuiltinFunction interface
  - Write tests for function registry and built-in function lookup
  - Write tests for function call parsing and argument validation
  - Write tests for error cases: unknown functions, wrong argument counts
  - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, TDD Requirements 5.1_

- [x] 7.2 Implement function call system to pass tests
  - Create FunctionCallExpression struct and BuiltinFunction interface
  - Implement function registry for built-in function lookup
  - Add function call parsing and execution as defined by tests
  - Verify all function call system tests pass
  - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, TDD Requirements 5.1_

- [x] 7.3 Write mathematical function tests first
  - Write failing tests for ABS, INT, and RND function implementations
  - Write tests for random number generator state management
  - Write tests for edge cases: negative numbers, zero, boundary values
  - Write tests for function argument validation and error handling
  - _Requirements: 5.1, 5.2, 5.3, TDD Requirements 5.1_

- [x] 7.4 Implement mathematical functions to pass tests
  - Create ABS, INT, and RND function implementations
  - Add random number generator state management in Environment
  - Implement argument validation and error handling as required by tests
  - Verify all mathematical function tests pass
  - _Requirements: 5.1, 5.2, 5.3, TDD Requirements 5.1_

- [x] 7.5 Write string function tests first
  - Write failing tests for LEN, MID$, STR$, and VAL function implementations
  - Write tests for string manipulation and conversion logic
  - Write tests for edge cases: empty strings, invalid indices, conversion errors
  - Write tests for function argument validation and type checking
  - _Requirements: 5.4, 5.5, TDD Requirements 5.2_

- [x] 7.6 Implement string functions to pass tests
  - Create LEN, MID$, STR$, and VAL function implementations
  - Add string manipulation and conversion logic as defined by tests
  - Implement proper error handling for invalid operations
  - Verify all string function tests pass
  - _Requirements: 5.4, 5.5, TDD Requirements 5.2_

- [x] 7.7 Refactor built-in functions
  - Improve function implementations while keeping all tests green
  - Extract common patterns and optimize function dispatch
  - Ensure clean separation between function logic and argument handling
  - Verify all tests still pass after refactoring
  - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5_

- [x] 8. Implement parser using TDD
- [x] 8.1 Write statement parser tests first
  - Write failing tests for Parser struct with ParseStatement method
  - Write tests for parsing each statement type (PRINT, INPUT, GOTO, IF, FOR, NEXT, assignment)
  - Write tests for line number handling and statement recognition
  - Write tests for error cases: invalid syntax, malformed statements
  - _Requirements: 1.1, 1.3, 6.2, TDD Requirements 1.2_

- [x] 8.2 Implement statement parser to pass tests
  - Create Parser struct with ParseStatement method for each statement type
  - Implement line number handling and statement recognition as defined by tests
  - Add proper error reporting for syntax errors
  - Verify all statement parser tests pass
  - _Requirements: 1.1, 1.3, 6.2, TDD Requirements 1.2_

- [x] 8.3 Write expression parser tests first
  - Write failing tests for ParseExpression method with operator precedence
  - Write tests for parentheses handling and function calls in expressions
  - Write tests for complex nested expressions and edge cases
  - Write tests for error cases: unbalanced parentheses, invalid operators
  - _Requirements: 2.2, 5.1, 5.2, 5.3, 5.4, 5.5, TDD Requirements 2.2_

- [x] 8.4 Implement expression parser to pass tests
  - Create ParseExpression method with operator precedence handling
  - Add support for parentheses and function calls as defined by tests
  - Implement proper error handling for expression parsing errors
  - Verify all expression parser tests pass
  - _Requirements: 2.2, 5.1, 5.2, 5.3, 5.4, 5.5, TDD Requirements 2.2_

- [x] 8.5 Write program parser tests first
  - Write failing tests for ParseProgram method building complete Program struct
  - Write tests for line number ordering and duplicate line detection
  - Write tests for complete program parsing with multiple statements
  - Write tests for error cases: duplicate lines, invalid program structure
  - _Requirements: 1.1, 1.2, TDD Requirements 1.2_

- [x] 8.6 Implement program parser to pass tests
  - Create ParseProgram method that builds complete Program struct
  - Implement line number ordering and duplicate line detection as required by tests
  - Add comprehensive error handling for program-level parsing errors
  - Verify all program parser tests pass
  - _Requirements: 1.1, 1.2, TDD Requirements 1.2_

- [x] 8.7 Refactor parser implementation
  - Improve parser code structure while keeping all tests green
  - Extract common parsing patterns and optimize error recovery
  - Ensure clean separation between different parsing responsibilities
  - Verify all tests still pass after refactoring
  - _Requirements: 1.1, 1.2, 1.3, 2.2, 5.1, 5.2, 5.3, 5.4, 5.5, 6.2_

- [x] 9. Implement interpreter execution engine using TDD
- [x] 9.1 Write program execution tests first
  - Write failing tests for Interpreter struct with Execute method for Program
  - Write tests for program counter management and statement execution sequencing
  - Write tests for basic program execution flow and control flow changes
  - Write tests for execution state management and variable persistence
  - _Requirements: 1.2, 6.1, TDD Requirements 1.2_

- [x] 9.2 Implement program execution to pass tests
  - Create Interpreter struct with Execute method for Program
  - Implement program counter management and statement execution sequencing
  - Add execution state management as defined by tests
  - Verify all program execution tests pass
  - _Requirements: 1.2, 6.1, TDD Requirements 1.2_

- [x] 9.3 Write error handling and debugging tests first
  - Write failing tests for runtime error reporting with line numbers
  - Write tests for optional debug mode showing each line before execution
  - Write tests for interrupt handling and infinite loop protection
  - Write tests for error message formatting and debug output
  - _Requirements: 6.1, 6.2, 6.3, 6.4, TDD Requirements 6.1_

- [x] 9.4 Implement error handling and debugging to pass tests
  - Add runtime error reporting with line numbers as defined by tests
  - Implement optional debug mode that shows execution progress
  - Create interrupt handling for infinite loop protection
  - Verify all error handling and debugging tests pass
  - _Requirements: 6.1, 6.2, 6.3, 6.4, TDD Requirements 6.1_

- [x] 9.5 Refactor interpreter execution engine
  - Improve interpreter code structure while keeping all tests green
  - Extract common execution patterns and optimize performance
  - Ensure clean separation between execution logic and error handling
  - Verify all tests still pass after refactoring
  - _Requirements: 1.2, 6.1, 6.2, 6.3, 6.4_

- [-] 10. Implement command-line interface using TDD
- [x] 10.1 Write CLI argument parsing tests first
  - Write failing tests for CLI struct with flag parsing for debug mode and help
  - Write tests for file argument handling and validation
  - Write tests for error cases: invalid flags, missing files, conflicting options
  - Write tests for help message formatting and version information
  - _Requirements: 7.1, 7.4, TDD Requirements 7.1_

- [x] 10.2 Implement CLI argument parsing to pass tests
  - Create CLI struct with flag parsing functionality
  - Implement file argument handling and validation as defined by tests
  - Add proper error handling and help message generation
  - Verify all CLI argument parsing tests pass
  - _Requirements: 7.1, 7.4, TDD Requirements 7.1_

- [x] 10.3 Write interactive mode tests first
  - Write failing tests for interactive REPL with mocked I/O
  - Write tests for line-by-line BASIC input processing
  - Write tests for command history and program state maintenance
  - Write tests for error handling in interactive mode
  - _Requirements: 7.2, 7.3, TDD Requirements 7.2_

- [x] 10.4 Implement interactive mode to pass tests
  - Create interactive REPL that accepts line-by-line BASIC input
  - Implement command processing and program state maintenance as defined by tests
  - Add basic line editing and history functionality
  - Verify all interactive mode tests pass
  - _Requirements: 7.2, 7.3, TDD Requirements 7.2_

- [x] 10.5 Write file execution mode tests first
  - Write failing tests for file loading and complete program execution
  - Write tests for file I/O error handling and validation
  - Write tests with sample BASIC program files and expected outputs
  - Write tests for integration between file loading and interpreter execution
  - _Requirements: 7.1, 1.1, TDD Requirements 7.3_

- [x] 10.6 Implement file execution mode to pass tests
  - Create file loading and complete program execution functionality
  - Implement proper error handling for file I/O operations as defined by tests
  - Add integration between file loading and interpreter execution
  - Verify all file execution mode tests pass
  - _Requirements: 7.1, 1.1, TDD Requirements 7.3_

- [x] 10.7 Refactor command-line interface
  - Improve CLI code structure while keeping all tests green
  - Extract common patterns and optimize user experience
  - Ensure clean separation between CLI logic and interpreter functionality
  - Verify all tests still pass after refactoring
  - _Requirements: 7.1, 7.2, 7.3, 7.4_

- [ ] 11. Create comprehensive integration tests using TDD
- [ ] 11.1 Write integration test suite first
  - Write failing tests for complete BASIC programs covering all language features
  - Write tests for complex programs with nested control structures and multiple statements
  - Write tests for error handling with various malformed programs
  - Write performance tests for reasonable program execution limits
  - _Requirements: All requirements, comprehensive integration testing_

- [ ] 11.2 Ensure all integration tests pass
  - Run comprehensive integration test suite and verify all tests pass
  - Fix any integration issues discovered by tests
  - Add additional test cases for edge cases discovered during integration
  - Verify all BASIC language features work together correctly in realistic programs
  - _Requirements: All requirements, comprehensive integration testing_

- [ ] 11.3 Refactor based on integration test feedback
  - Improve code structure based on integration test results while keeping tests green
  - Optimize performance for complex programs as identified by performance tests
  - Enhance error handling based on integration test scenarios
  - Verify all tests still pass after integration-driven refactoring
  - _Requirements: All requirements_

- [ ] 12. Create main application entry point using TDD
- [ ] 12.1 Write main application tests first
  - Write failing tests for main.go that ties together CLI, file loading, and interpreter execution
  - Write tests for proper application initialization and cleanup
  - Write tests for complete application workflow in all supported usage modes
  - Write tests for end-to-end application behavior with real file I/O
  - _Requirements: 7.1, 7.2, 7.3, 7.4, TDD Requirements 7.4_

- [ ] 12.2 Implement main application to pass tests
  - Create main.go that integrates CLI, file loading, and interpreter execution
  - Implement proper application initialization and cleanup as defined by tests
  - Add complete workflow support for all usage modes
  - Verify all main application tests pass
  - _Requirements: 7.1, 7.2, 7.3, 7.4, TDD Requirements 7.4_

- [ ] 12.3 Final refactoring and cleanup
  - Perform final code review and refactoring while keeping all tests green
  - Ensure clean architecture and optimal performance
  - Verify complete application works end-to-end for all supported usage modes
  - Run full test suite to confirm all functionality works correctly
  - _Requirements: 7.1, 7.2, 7.3, 7.4_