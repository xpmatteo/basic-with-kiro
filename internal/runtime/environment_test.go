package runtime

import (
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