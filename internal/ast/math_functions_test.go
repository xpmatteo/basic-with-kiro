package ast

import (
	"basic-interpreter/internal/runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test ABS function implementation
func TestAbsFunction_Call(t *testing.T) {
	env := runtime.NewEnvironment()
	fn := GetBuiltinFunction("ABS")
	require.NotNil(t, fn)

	testCases := []struct {
		name     string
		input    float64
		expected float64
	}{
		{
			name:     "positive number",
			input:    5.5,
			expected: 5.5,
		},
		{
			name:     "negative number",
			input:    -5.5,
			expected: 5.5,
		},
		{
			name:     "zero",
			input:    0.0,
			expected: 0.0,
		},
		{
			name:     "negative zero",
			input:    -0.0,
			expected: 0.0,
		},
		{
			name:     "large positive number",
			input:    999999.999,
			expected: 999999.999,
		},
		{
			name:     "large negative number",
			input:    -999999.999,
			expected: 999999.999,
		},
		{
			name:     "very small positive number",
			input:    0.000001,
			expected: 0.000001,
		},
		{
			name:     "very small negative number",
			input:    -0.000001,
			expected: 0.000001,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			args := []runtime.Value{runtime.NewNumericValue(tc.input)}
			result, err := fn.Call(args, env)

			require.NoError(t, err)
			assert.Equal(t, runtime.NumericValue, result.Type)
			assert.Equal(t, tc.expected, result.NumValue)
		})
	}
}

func TestAbsFunction_ErrorCases(t *testing.T) {
	env := runtime.NewEnvironment()
	fn := GetBuiltinFunction("ABS")
	require.NotNil(t, fn)

	t.Run("wrong argument count - none", func(t *testing.T) {
		args := []runtime.Value{}
		_, err := fn.Call(args, env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected 1 argument")
	})

	t.Run("wrong argument count - too many", func(t *testing.T) {
		args := []runtime.Value{
			runtime.NewNumericValue(5),
			runtime.NewNumericValue(3),
		}
		_, err := fn.Call(args, env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected 1 argument")
	})

	t.Run("wrong argument type - string", func(t *testing.T) {
		args := []runtime.Value{runtime.NewStringValue("hello")}
		_, err := fn.Call(args, env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "argument must be numeric")
	})
}

// Test INT function implementation
func TestIntFunction_Call(t *testing.T) {
	env := runtime.NewEnvironment()
	fn := GetBuiltinFunction("INT")
	require.NotNil(t, fn)

	testCases := []struct {
		name     string
		input    float64
		expected float64
	}{
		{
			name:     "positive integer",
			input:    5.0,
			expected: 5.0,
		},
		{
			name:     "positive decimal",
			input:    5.7,
			expected: 5.0,
		},
		{
			name:     "positive decimal close to next integer",
			input:    5.99,
			expected: 5.0,
		},
		{
			name:     "negative integer",
			input:    -5.0,
			expected: -5.0,
		},
		{
			name:     "negative decimal",
			input:    -5.7,
			expected: -5.0,
		},
		{
			name:     "negative decimal close to next integer",
			input:    -5.99,
			expected: -5.0,
		},
		{
			name:     "zero",
			input:    0.0,
			expected: 0.0,
		},
		{
			name:     "positive small decimal",
			input:    0.9,
			expected: 0.0,
		},
		{
			name:     "negative small decimal",
			input:    -0.9,
			expected: 0.0,
		},
		{
			name:     "large positive number",
			input:    999999.999,
			expected: 999999.0,
		},
		{
			name:     "large negative number",
			input:    -999999.999,
			expected: -999999.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			args := []runtime.Value{runtime.NewNumericValue(tc.input)}
			result, err := fn.Call(args, env)

			require.NoError(t, err)
			assert.Equal(t, runtime.NumericValue, result.Type)
			assert.Equal(t, tc.expected, result.NumValue)
		})
	}
}

func TestIntFunction_ErrorCases(t *testing.T) {
	env := runtime.NewEnvironment()
	fn := GetBuiltinFunction("INT")
	require.NotNil(t, fn)

	t.Run("wrong argument count - none", func(t *testing.T) {
		args := []runtime.Value{}
		_, err := fn.Call(args, env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected 1 argument")
	})

	t.Run("wrong argument count - too many", func(t *testing.T) {
		args := []runtime.Value{
			runtime.NewNumericValue(5.5),
			runtime.NewNumericValue(3.3),
		}
		_, err := fn.Call(args, env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected 1 argument")
	})

	t.Run("wrong argument type - string", func(t *testing.T) {
		args := []runtime.Value{runtime.NewStringValue("5.5")}
		_, err := fn.Call(args, env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "argument must be numeric")
	})
}

// Test RND function implementation and random number generator state management
func TestRndFunction_Call(t *testing.T) {
	env := runtime.NewEnvironment()
	fn := GetBuiltinFunction("RND")
	require.NotNil(t, fn)

	t.Run("returns value in correct range", func(t *testing.T) {
		args := []runtime.Value{}
		result, err := fn.Call(args, env)

		require.NoError(t, err)
		assert.Equal(t, runtime.NumericValue, result.Type)
		assert.GreaterOrEqual(t, result.NumValue, 0.0)
		assert.Less(t, result.NumValue, 1.0)
	})

	t.Run("generates different values on multiple calls", func(t *testing.T) {
		args := []runtime.Value{}
		
		// Generate multiple random numbers
		var results []float64
		for i := 0; i < 10; i++ {
			result, err := fn.Call(args, env)
			require.NoError(t, err)
			results = append(results, result.NumValue)
		}

		// Check that not all values are the same (very unlikely with proper randomness)
		allSame := true
		firstValue := results[0]
		for _, value := range results[1:] {
			if value != firstValue {
				allSame = false
				break
			}
		}
		assert.False(t, allSame, "RND should generate different values on multiple calls")
	})

	t.Run("uses environment random number generator", func(t *testing.T) {
		// Create two environments with the same seed
		env1 := runtime.NewEnvironment()
		env2 := runtime.NewEnvironment()
		
		// Set the same seed for both environments
		env1.SetRandomSeed(12345)
		env2.SetRandomSeed(12345)

		args := []runtime.Value{}
		
		// Generate numbers from both environments
		result1, err1 := fn.Call(args, env1)
		result2, err2 := fn.Call(args, env2)

		require.NoError(t, err1)
		require.NoError(t, err2)
		
		// With the same seed, the first random number should be the same
		assert.Equal(t, result1.NumValue, result2.NumValue)
	})

	t.Run("different seeds produce different sequences", func(t *testing.T) {
		// Create two environments with different seeds
		env1 := runtime.NewEnvironment()
		env2 := runtime.NewEnvironment()
		
		env1.SetRandomSeed(12345)
		env2.SetRandomSeed(54321)

		args := []runtime.Value{}
		
		// Generate sequences from both environments
		var seq1, seq2 []float64
		for i := 0; i < 5; i++ {
			result1, err1 := fn.Call(args, env1)
			result2, err2 := fn.Call(args, env2)
			
			require.NoError(t, err1)
			require.NoError(t, err2)
			
			seq1 = append(seq1, result1.NumValue)
			seq2 = append(seq2, result2.NumValue)
		}

		// Sequences should be different
		assert.NotEqual(t, seq1, seq2)
	})
}

func TestRndFunction_ErrorCases(t *testing.T) {
	env := runtime.NewEnvironment()
	fn := GetBuiltinFunction("RND")
	require.NotNil(t, fn)

	t.Run("wrong argument count - too many", func(t *testing.T) {
		args := []runtime.Value{runtime.NewNumericValue(5)}
		_, err := fn.Call(args, env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected 0 arguments")
	})
}

// Test random number generator state management in Environment
func TestRandomNumberGeneratorStateManagement(t *testing.T) {
	t.Run("environment maintains random state", func(t *testing.T) {
		env := runtime.NewEnvironment()
		fn := GetBuiltinFunction("RND")
		require.NotNil(t, fn)

		args := []runtime.Value{}
		
		// Generate a sequence of random numbers
		var sequence1 []float64
		for i := 0; i < 5; i++ {
			result, err := fn.Call(args, env)
			require.NoError(t, err)
			sequence1 = append(sequence1, result.NumValue)
		}

		// Reset the environment with the same seed
		originalSeed := env.GetRandomSeed()
		env.SetRandomSeed(originalSeed)

		// Generate the same sequence again
		var sequence2 []float64
		for i := 0; i < 5; i++ {
			result, err := fn.Call(args, env)
			require.NoError(t, err)
			sequence2 = append(sequence2, result.NumValue)
		}

		// Sequences should be identical when using the same seed
		assert.Equal(t, sequence1, sequence2)
	})

	t.Run("seed changes affect subsequent random numbers", func(t *testing.T) {
		env := runtime.NewEnvironment()
		fn := GetBuiltinFunction("RND")
		require.NotNil(t, fn)

		args := []runtime.Value{}
		
		// Generate a number with initial seed
		result1, err := fn.Call(args, env)
		require.NoError(t, err)

		// Change the seed
		env.SetRandomSeed(99999)

		// Generate another number
		result2, err := fn.Call(args, env)
		require.NoError(t, err)

		// Results should be different due to seed changes
		assert.NotEqual(t, result1.NumValue, result2.NumValue)
	})
}

// Test boundary values and edge cases for mathematical functions
func TestMathematicalFunctionsBoundaryValues(t *testing.T) {
	env := runtime.NewEnvironment()

	t.Run("ABS with boundary values", func(t *testing.T) {
		fn := GetBuiltinFunction("ABS")
		require.NotNil(t, fn)

		testCases := []struct {
			name     string
			input    float64
			expected float64
		}{
			{"positive infinity", 1e308, 1e308},
			{"negative infinity", -1e308, 1e308},
			{"smallest positive", 1e-308, 1e-308},
			{"smallest negative", -1e-308, 1e-308},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				args := []runtime.Value{runtime.NewNumericValue(tc.input)}
				result, err := fn.Call(args, env)

				require.NoError(t, err)
				assert.Equal(t, tc.expected, result.NumValue)
			})
		}
	})

	t.Run("INT with boundary values", func(t *testing.T) {
		fn := GetBuiltinFunction("INT")
		require.NotNil(t, fn)

		testCases := []struct {
			name     string
			input    float64
			expected float64
		}{
			{"large positive", 1e10, 1e10},
			{"large negative", -1e10, -1e10},
			{"just below 1", 0.9999999, 0.0},
			{"just above -1", -0.9999999, 0.0},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				args := []runtime.Value{runtime.NewNumericValue(tc.input)}
				result, err := fn.Call(args, env)

				require.NoError(t, err)
				assert.Equal(t, tc.expected, result.NumValue)
			})
		}
	})
}

// Test function argument validation for mathematical functions
func TestMathematicalFunctionArgumentValidation(t *testing.T) {
	env := runtime.NewEnvironment()

	testCases := []struct {
		functionName string
		expectedArgs int
	}{
		{"ABS", 1},
		{"INT", 1},
		{"RND", 0},
	}

	for _, tc := range testCases {
		t.Run(tc.functionName+" argument validation", func(t *testing.T) {
			fn := GetBuiltinFunction(tc.functionName)
			require.NotNil(t, fn)

			// Test correct argument count
			assert.Equal(t, tc.expectedArgs, fn.ArgCount())

			// Test with wrong number of arguments
			if tc.expectedArgs > 0 {
				// Test with no arguments when some are expected
				args := []runtime.Value{}
				_, err := fn.Call(args, env)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "expected")

				// Test with too many arguments
				args = make([]runtime.Value, tc.expectedArgs+1)
				for i := range args {
					args[i] = runtime.NewNumericValue(1.0)
				}
				_, err = fn.Call(args, env)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "expected")
			} else {
				// Test RND with arguments when none are expected
				args := []runtime.Value{runtime.NewNumericValue(1.0)}
				_, err := fn.Call(args, env)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "expected 0 arguments")
			}
		})
	}
}