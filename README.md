# An experiment in AI-assisted coding

I wanted to try [Harper Reed's process](https://harper.blog/2025/05/08/basic-claude-code/ "Basic Claude Code | Harper Reed's Blog"). I also wanted to try Kiro.  

**Status**: done -- I do not expect to work anymore on this


# My experience with Kiro

TL;DR: it's much easier to use than most current tools, at least for greenfield small projects.  It's fun.  It still needs expert guidance.  Great potential for onboarding teams to AI-assisted coding

The good: 
it incorporates emerging best practices; it does a more or less decent job of planning an MVP.  
It helps you with creating "steering" documents eg style rules etc.  
It looks like it would do a good job to get a team started with working with AI assistance.
UI is streamlined; less confusing than cursor or VS Code with Copilot
it will do a great job of fixing its own bugs

The bad:
usage cost is out of control.  It is based on "requests", but one single task out of the task list will consume an unpredictable amount of requests. I ran out of bonus requests in an afternoon
the steps in the plan it builds are mostly too big for the context window
you must explain to it how to do TDD
poor context management.  All of a sudden it will stop and tell me that the context is exhausted and I should compact

The experience: emboldened by Harper Reed's report on being able to generate a BASIC interpreter on the spot, I tried the same with Kiro.  I started a new project, it generated a plan.

The plan was mostly waterfall, eg step 2 is create the lexer, step 3 is create the parser, step 10 is make it work from the CLI.  I let it slide because I wanted to see if it would work in the end.  Given that the whole "waterfall" was completed in 2-3 hours of mostly unattended work, I guess it can be a viable approach at this scale.

Unfortunately, the second step in the plan was too difficult for the agent to complete.  I then asked Kiro to rewrite the plan applying TDD.  It rewrote all the steps, breaking them down in three phases 1. create tests 2. make them pass 3. refactor.

Of course, human TDD guidance says you write one test and then you make it pass; you don't write all the tests for a feature upfront. But this is an agent, not a human, and it seems it works well for agents to write all the tests for a feature.  Compare with classic TDD guidance, where you are encouraged to write a test list as the first step: what the agent does in the end is writing the test list, and given that it has way more energy and speed than a human, it's not unconceivable that writing all the tests for a feature upfront is a good way to go for agents.
It turns out that writing the tests first helps the AI to work better; with this new plan, it was able to complete the lexer autonomously, and it did most of the other steps unattended (I was running another experiment with Claude Code in another terminal, so I did not pay much attention to Kiro)

Towards the end I asked it to run a refactoring, using the classic four rules of simple design from Kent Beck: a very simple prompt that worked fairly well.  The generated plan included refactoring steps, and they seemed to help, as it "extracted helper functions" at the end of most steps, but asking repeatedly it to go deeper looking for improvements seems to work.

Then I asked it to test the interpreter by hand, and by doing this it fixed several bugs.

The end result: I mostly did not pay attention to it, and it built a working prototype in an afternoon.  It's testing and fixing bugs as I write this

And the code it wrote is not bad at all.



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