package runtime

import (
	"math/rand"
	"strings"
	"time"
)

// ForLoopState represents the state of a FOR loop
type ForLoopState struct {
	Variable string
	Current  float64
	End      float64
	Step     float64
	LineNum  int
}

// Environment represents the runtime environment for BASIC program execution
type Environment struct {
	Variables      map[string]Value    // Case-insensitive variable storage
	ProgramCounter int                 // Current line number being executed
	CallStack      []int               // Stack for nested control structures
	ForLoops       []ForLoopState      // Stack for nested FOR loops
	RandomSeed     int64               // Seed for random number generation
	rng            *rand.Rand          // Random number generator
}

// NewEnvironment creates a new runtime environment
func NewEnvironment() *Environment {
	seed := time.Now().UnixNano()
	return &Environment{
		Variables:      make(map[string]Value),
		ProgramCounter: 0,
		CallStack:      make([]int, 0),
		ForLoops:       make([]ForLoopState, 0),
		RandomSeed:     seed,
		rng:            rand.New(rand.NewSource(seed)),
	}
}

// GetVariable retrieves a variable value (case-insensitive)
func (env *Environment) GetVariable(name string) Value {
	key := strings.ToUpper(name)
	if val, exists := env.Variables[key]; exists {
		return val
	}
	
	// BASIC initializes undefined numeric variables to 0 and string variables to empty string
	if strings.HasSuffix(key, "$") {
		return NewStringValue("")
	}
	return NewNumericValue(0)
}

// SetVariable sets a variable value (case-insensitive)
func (env *Environment) SetVariable(name string, value Value) {
	key := strings.ToUpper(name)
	env.Variables[key] = value
}

// Random returns a random number between 0 and 1
func (env *Environment) Random() float64 {
	return env.rng.Float64()
}