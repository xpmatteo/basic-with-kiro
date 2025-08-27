package runtime

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEnvironment(t *testing.T) {
	env := NewEnvironment()
	require.NotNil(t, env)
	assert.NotNil(t, env.Variables)
	assert.Equal(t, 0, env.ProgramCounter)
	assert.NotNil(t, env.CallStack)
	assert.NotNil(t, env.ForLoops)
	assert.NotZero(t, env.RandomSeed)
}

func TestVariableOperations(t *testing.T) {
	env := NewEnvironment()

	// Test setting and getting a numeric variable
	numVal := NewNumericValue(42)
	env.SetVariable("X", numVal)
	
	result := env.GetVariable("X")
	assert.Equal(t, NumericValue, result.Type)
	assert.Equal(t, 42.0, result.NumValue)

	// Test case insensitivity
	result = env.GetVariable("x")
	assert.Equal(t, NumericValue, result.Type)
	assert.Equal(t, 42.0, result.NumValue)

	// Test string variable
	strVal := NewStringValue("hello")
	env.SetVariable("NAME$", strVal)
	
	result = env.GetVariable("name$")
	assert.Equal(t, StringValue, result.Type)
	assert.Equal(t, "hello", result.StrValue)
}

func TestUndefinedVariables(t *testing.T) {
	env := NewEnvironment()

	// Test undefined numeric variable (should default to 0)
	result := env.GetVariable("UNDEFINED")
	assert.Equal(t, NumericValue, result.Type)
	assert.Equal(t, 0.0, result.NumValue)

	// Test undefined string variable (should default to empty string)
	result = env.GetVariable("UNDEFINED$")
	assert.Equal(t, StringValue, result.Type)
	assert.Equal(t, "", result.StrValue)
}

func TestRandom(t *testing.T) {
	env := NewEnvironment()
	
	// Test that random returns a value between 0 and 1
	for i := 0; i < 10; i++ {
		val := env.Random()
		assert.GreaterOrEqual(t, val, 0.0)
		assert.Less(t, val, 1.0)
	}

	// Test that different calls return different values (with high probability)
	val1 := env.Random()
	val2 := env.Random()
	assert.NotEqual(t, val1, val2, "Random should return different values")
}

// Additional comprehensive tests for Environment as required by task 3.3

func TestEnvironmentVariableStorage_CaseInsensitive(t *testing.T) {
	env := NewEnvironment()

	t.Run("Case insensitive variable names", func(t *testing.T) {
		// Set variable with lowercase name
		env.SetVariable("myvar", NewNumericValue(100))
		
		// Retrieve with different cases
		assert.Equal(t, 100.0, env.GetVariable("myvar").NumValue)
		assert.Equal(t, 100.0, env.GetVariable("MYVAR").NumValue)
		assert.Equal(t, 100.0, env.GetVariable("MyVar").NumValue)
		assert.Equal(t, 100.0, env.GetVariable("myVAR").NumValue)
	})

	t.Run("String variables case insensitive", func(t *testing.T) {
		env.SetVariable("name$", NewStringValue("test"))
		
		assert.Equal(t, "test", env.GetVariable("name$").StrValue)
		assert.Equal(t, "test", env.GetVariable("NAME$").StrValue)
		assert.Equal(t, "test", env.GetVariable("Name$").StrValue)
	})

	t.Run("Overwrite variable with different case", func(t *testing.T) {
		env.SetVariable("counter", NewNumericValue(1))
		env.SetVariable("COUNTER", NewNumericValue(2))
		
		// Should overwrite the same variable
		assert.Equal(t, 2.0, env.GetVariable("counter").NumValue)
		assert.Equal(t, 2.0, env.GetVariable("COUNTER").NumValue)
	})
}

func TestEnvironmentDefaultValues(t *testing.T) {
	env := NewEnvironment()

	t.Run("Numeric variables default to zero", func(t *testing.T) {
		testVars := []string{"X", "Y", "COUNTER", "index", "MyVar"}
		for _, varName := range testVars {
			val := env.GetVariable(varName)
			assert.Equal(t, NumericValue, val.Type, "Variable %s should be numeric", varName)
			assert.Equal(t, 0.0, val.NumValue, "Variable %s should default to 0", varName)
		}
	})

	t.Run("String variables default to empty string", func(t *testing.T) {
		testVars := []string{"NAME$", "TEXT$", "message$", "MyString$"}
		for _, varName := range testVars {
			val := env.GetVariable(varName)
			assert.Equal(t, StringValue, val.Type, "Variable %s should be string", varName)
			assert.Equal(t, "", val.StrValue, "Variable %s should default to empty string", varName)
		}
	})

	t.Run("Mixed case string variable defaults", func(t *testing.T) {
		val := env.GetVariable("MixedCase$")
		assert.Equal(t, StringValue, val.Type)
		assert.Equal(t, "", val.StrValue)
	})
}

func TestEnvironmentStateManagement(t *testing.T) {
	env := NewEnvironment()

	t.Run("Program counter management", func(t *testing.T) {
		assert.Equal(t, 0, env.ProgramCounter)
		
		env.ProgramCounter = 10
		assert.Equal(t, 10, env.ProgramCounter)
		
		env.ProgramCounter = 100
		assert.Equal(t, 100, env.ProgramCounter)
	})

	t.Run("Call stack operations", func(t *testing.T) {
		assert.Empty(t, env.CallStack)
		
		// Push values onto call stack
		env.CallStack = append(env.CallStack, 10)
		env.CallStack = append(env.CallStack, 20)
		env.CallStack = append(env.CallStack, 30)
		
		assert.Len(t, env.CallStack, 3)
		assert.Equal(t, []int{10, 20, 30}, env.CallStack)
		
		// Pop from call stack
		env.CallStack = env.CallStack[:len(env.CallStack)-1]
		assert.Len(t, env.CallStack, 2)
		assert.Equal(t, []int{10, 20}, env.CallStack)
	})

	t.Run("FOR loop state management", func(t *testing.T) {
		assert.Empty(t, env.ForLoops)
		
		// Add FOR loop state
		forState := ForLoopState{
			Variable: "I",
			Current:  1,
			End:      10,
			Step:     1,
			LineNum:  20,
		}
		env.ForLoops = append(env.ForLoops, forState)
		
		assert.Len(t, env.ForLoops, 1)
		assert.Equal(t, "I", env.ForLoops[0].Variable)
		assert.Equal(t, 1.0, env.ForLoops[0].Current)
		assert.Equal(t, 10.0, env.ForLoops[0].End)
		assert.Equal(t, 1.0, env.ForLoops[0].Step)
		assert.Equal(t, 20, env.ForLoops[0].LineNum)
	})
}

func TestEnvironmentVariableScoping(t *testing.T) {
	env := NewEnvironment()

	t.Run("Variable persistence across operations", func(t *testing.T) {
		// Set multiple variables
		env.SetVariable("X", NewNumericValue(10))
		env.SetVariable("Y", NewNumericValue(20))
		env.SetVariable("NAME$", NewStringValue("test"))
		
		// Modify program counter and call stack
		env.ProgramCounter = 50
		env.CallStack = append(env.CallStack, 100)
		
		// Variables should still be accessible
		assert.Equal(t, 10.0, env.GetVariable("X").NumValue)
		assert.Equal(t, 20.0, env.GetVariable("Y").NumValue)
		assert.Equal(t, "test", env.GetVariable("NAME$").StrValue)
	})

	t.Run("Variable updates", func(t *testing.T) {
		// Set initial value
		env.SetVariable("COUNTER", NewNumericValue(0))
		assert.Equal(t, 0.0, env.GetVariable("COUNTER").NumValue)
		
		// Update value multiple times
		env.SetVariable("COUNTER", NewNumericValue(1))
		assert.Equal(t, 1.0, env.GetVariable("COUNTER").NumValue)
		
		env.SetVariable("COUNTER", NewNumericValue(5))
		assert.Equal(t, 5.0, env.GetVariable("COUNTER").NumValue)
		
		// Change type (numeric to string with same name - different variable)
		env.SetVariable("COUNTER$", NewStringValue("five"))
		assert.Equal(t, 5.0, env.GetVariable("COUNTER").NumValue) // Original numeric still exists
		assert.Equal(t, "five", env.GetVariable("COUNTER$").StrValue) // New string variable
	})
}

func TestEnvironmentRandomNumberGeneration(t *testing.T) {
	t.Run("Random seed consistency", func(t *testing.T) {
		// Create two environments with same seed
		env1 := NewEnvironment()
		env2 := NewEnvironment()
		
		// Set same seed
		env1.RandomSeed = 12345
		env1.rng = rand.New(rand.NewSource(12345))
		env2.RandomSeed = 12345
		env2.rng = rand.New(rand.NewSource(12345))
		
		// Should generate same sequence
		for i := 0; i < 5; i++ {
			val1 := env1.Random()
			val2 := env2.Random()
			assert.Equal(t, val1, val2, "Same seed should produce same random sequence")
		}
	})

	t.Run("Random number range", func(t *testing.T) {
		env := NewEnvironment()
		
		// Generate many random numbers and verify range
		for i := 0; i < 1000; i++ {
			val := env.Random()
			assert.GreaterOrEqual(t, val, 0.0, "Random value should be >= 0")
			assert.Less(t, val, 1.0, "Random value should be < 1")
		}
	})

	t.Run("Random distribution", func(t *testing.T) {
		env := NewEnvironment()
		
		// Generate random numbers and check basic distribution
		var sum float64
		count := 1000
		for i := 0; i < count; i++ {
			sum += env.Random()
		}
		
		average := sum / float64(count)
		// Average should be around 0.5 for uniform distribution
		assert.Greater(t, average, 0.4, "Average should be > 0.4")
		assert.Less(t, average, 0.6, "Average should be < 0.6")
	})
}

func TestEnvironmentEdgeCases(t *testing.T) {
	env := NewEnvironment()

	t.Run("Empty variable name", func(t *testing.T) {
		// Should handle empty variable names gracefully
		env.SetVariable("", NewNumericValue(42))
		val := env.GetVariable("")
		assert.Equal(t, 42.0, val.NumValue)
	})

	t.Run("Variable names with special characters", func(t *testing.T) {
		// BASIC typically only allows alphanumeric + $ for strings
		env.SetVariable("VAR123", NewNumericValue(123))
		assert.Equal(t, 123.0, env.GetVariable("VAR123").NumValue)
		
		env.SetVariable("STR123$", NewStringValue("test"))
		assert.Equal(t, "test", env.GetVariable("STR123$").StrValue)
	})

	t.Run("Very long variable names", func(t *testing.T) {
		longName := "VERYLONGVARIABLENAME123456789"
		env.SetVariable(longName, NewNumericValue(999))
		assert.Equal(t, 999.0, env.GetVariable(longName).NumValue)
	})

	t.Run("Multiple environments isolation", func(t *testing.T) {
		env1 := NewEnvironment()
		env2 := NewEnvironment()
		
		env1.SetVariable("X", NewNumericValue(100))
		env2.SetVariable("X", NewNumericValue(200))
		
		// Environments should be isolated
		assert.Equal(t, 100.0, env1.GetVariable("X").NumValue)
		assert.Equal(t, 200.0, env2.GetVariable("X").NumValue)
	})
}