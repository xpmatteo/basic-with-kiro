package ast

import (
	"basic-interpreter/internal/runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test LEN function implementation
func TestLenFunction_Call(t *testing.T) {
	env := runtime.NewEnvironment()
	fn := GetBuiltinFunction("LEN")
	require.NotNil(t, fn)

	testCases := []struct {
		name     string
		input    string
		expected float64
	}{
		{
			name:     "empty string",
			input:    "",
			expected: 0.0,
		},
		{
			name:     "single character",
			input:    "A",
			expected: 1.0,
		},
		{
			name:     "normal string",
			input:    "HELLO",
			expected: 5.0,
		},
		{
			name:     "string with spaces",
			input:    "HELLO WORLD",
			expected: 11.0,
		},
		{
			name:     "string with special characters",
			input:    "Hello, World!",
			expected: 13.0,
		},
		{
			name:     "string with numbers",
			input:    "123456",
			expected: 6.0,
		},
		{
			name:     "string with mixed content",
			input:    "ABC123!@#",
			expected: 9.0,
		},
		{
			name:     "string with tabs and newlines",
			input:    "A\tB\nC",
			expected: 5.0,
		},
		{
			name:     "very long string",
			input:    "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
			expected: 36.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			args := []runtime.Value{runtime.NewStringValue(tc.input)}
			result, err := fn.Call(args, env)

			require.NoError(t, err)
			assert.Equal(t, runtime.NumericValue, result.Type)
			assert.Equal(t, tc.expected, result.NumValue)
		})
	}
}

func TestLenFunction_ErrorCases(t *testing.T) {
	env := runtime.NewEnvironment()
	fn := GetBuiltinFunction("LEN")
	require.NotNil(t, fn)

	t.Run("wrong argument count - none", func(t *testing.T) {
		args := []runtime.Value{}
		_, err := fn.Call(args, env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected 1 argument")
	})

	t.Run("wrong argument count - too many", func(t *testing.T) {
		args := []runtime.Value{
			runtime.NewStringValue("hello"),
			runtime.NewStringValue("world"),
		}
		_, err := fn.Call(args, env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected 1 argument")
	})

	t.Run("wrong argument type - numeric", func(t *testing.T) {
		args := []runtime.Value{runtime.NewNumericValue(123)}
		_, err := fn.Call(args, env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "argument must be string")
	})
}

// Test MID$ function implementation
func TestMidFunction_Call(t *testing.T) {
	env := runtime.NewEnvironment()
	fn := GetBuiltinFunction("MID$")
	require.NotNil(t, fn)

	testCases := []struct {
		name     string
		str      string
		start    float64
		length   float64
		expected string
	}{
		{
			name:     "normal substring",
			str:      "HELLO",
			start:    2.0,
			length:   3.0,
			expected: "ELL",
		},
		{
			name:     "substring from beginning",
			str:      "HELLO",
			start:    1.0,
			length:   3.0,
			expected: "HEL",
		},
		{
			name:     "substring to end",
			str:      "HELLO",
			start:    3.0,
			length:   10.0, // Length longer than remaining string
			expected: "LLO",
		},
		{
			name:     "single character",
			str:      "HELLO",
			start:    1.0,
			length:   1.0,
			expected: "H",
		},
		{
			name:     "last character",
			str:      "HELLO",
			start:    5.0,
			length:   1.0,
			expected: "O",
		},
		{
			name:     "empty string input",
			str:      "",
			start:    1.0,
			length:   1.0,
			expected: "",
		},
		{
			name:     "zero length",
			str:      "HELLO",
			start:    2.0,
			length:   0.0,
			expected: "",
		},
		{
			name:     "start beyond string length",
			str:      "HELLO",
			start:    10.0,
			length:   3.0,
			expected: "",
		},
		{
			name:     "start at string length",
			str:      "HELLO",
			start:    5.0,
			length:   3.0,
			expected: "O",
		},
		{
			name:     "string with spaces",
			str:      "HELLO WORLD",
			start:    6.0,
			length:   6.0,
			expected: " WORLD",
		},
		{
			name:     "string with special characters",
			str:      "Hello, World!",
			start:    7.0,
			length:   5.0,
			expected: " Worl",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			args := []runtime.Value{
				runtime.NewStringValue(tc.str),
				runtime.NewNumericValue(tc.start),
				runtime.NewNumericValue(tc.length),
			}
			result, err := fn.Call(args, env)

			require.NoError(t, err)
			assert.Equal(t, runtime.StringValue, result.Type)
			assert.Equal(t, tc.expected, result.StrValue)
		})
	}
}

func TestMidFunction_EdgeCases(t *testing.T) {
	env := runtime.NewEnvironment()
	fn := GetBuiltinFunction("MID$")
	require.NotNil(t, fn)

	t.Run("negative start index", func(t *testing.T) {
		args := []runtime.Value{
			runtime.NewStringValue("HELLO"),
			runtime.NewNumericValue(-1.0),
			runtime.NewNumericValue(3.0),
		}
		result, err := fn.Call(args, env)

		require.NoError(t, err)
		assert.Equal(t, runtime.StringValue, result.Type)
		assert.Equal(t, "", result.StrValue) // Should return empty string for invalid start
	})

	t.Run("zero start index", func(t *testing.T) {
		args := []runtime.Value{
			runtime.NewStringValue("HELLO"),
			runtime.NewNumericValue(0.0),
			runtime.NewNumericValue(3.0),
		}
		result, err := fn.Call(args, env)

		require.NoError(t, err)
		assert.Equal(t, runtime.StringValue, result.Type)
		assert.Equal(t, "", result.StrValue) // Should return empty string for invalid start (BASIC uses 1-based indexing)
	})

	t.Run("negative length", func(t *testing.T) {
		args := []runtime.Value{
			runtime.NewStringValue("HELLO"),
			runtime.NewNumericValue(2.0),
			runtime.NewNumericValue(-1.0),
		}
		result, err := fn.Call(args, env)

		require.NoError(t, err)
		assert.Equal(t, runtime.StringValue, result.Type)
		assert.Equal(t, "", result.StrValue) // Should return empty string for negative length
	})
}

func TestMidFunction_ErrorCases(t *testing.T) {
	env := runtime.NewEnvironment()
	fn := GetBuiltinFunction("MID$")
	require.NotNil(t, fn)

	t.Run("wrong argument count - too few", func(t *testing.T) {
		args := []runtime.Value{
			runtime.NewStringValue("HELLO"),
			runtime.NewNumericValue(2.0),
		}
		_, err := fn.Call(args, env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected 3 arguments")
	})

	t.Run("wrong argument count - too many", func(t *testing.T) {
		args := []runtime.Value{
			runtime.NewStringValue("HELLO"),
			runtime.NewNumericValue(2.0),
			runtime.NewNumericValue(3.0),
			runtime.NewNumericValue(4.0),
		}
		_, err := fn.Call(args, env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected 3 arguments")
	})

	t.Run("first argument not string", func(t *testing.T) {
		args := []runtime.Value{
			runtime.NewNumericValue(123),
			runtime.NewNumericValue(2.0),
			runtime.NewNumericValue(3.0),
		}
		_, err := fn.Call(args, env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "first argument must be string")
	})

	t.Run("second argument not numeric", func(t *testing.T) {
		args := []runtime.Value{
			runtime.NewStringValue("HELLO"),
			runtime.NewStringValue("2"),
			runtime.NewNumericValue(3.0),
		}
		_, err := fn.Call(args, env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "second argument must be numeric")
	})

	t.Run("third argument not numeric", func(t *testing.T) {
		args := []runtime.Value{
			runtime.NewStringValue("HELLO"),
			runtime.NewNumericValue(2.0),
			runtime.NewStringValue("3"),
		}
		_, err := fn.Call(args, env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "third argument must be numeric")
	})
}

// Test STR$ function implementation
func TestStrFunction_Call(t *testing.T) {
	env := runtime.NewEnvironment()
	fn := GetBuiltinFunction("STR$")
	require.NotNil(t, fn)

	testCases := []struct {
		name     string
		input    float64
		expected string
	}{
		{
			name:     "positive integer",
			input:    42.0,
			expected: "42",
		},
		{
			name:     "negative integer",
			input:    -42.0,
			expected: "-42",
		},
		{
			name:     "zero",
			input:    0.0,
			expected: "0",
		},
		{
			name:     "positive decimal",
			input:    3.14,
			expected: "3.14",
		},
		{
			name:     "negative decimal",
			input:    -3.14,
			expected: "-3.14",
		},
		{
			name:     "very small number",
			input:    0.001,
			expected: "0.001",
		},
		{
			name:     "very large number",
			input:    999999.0,
			expected: "999999",
		},
		{
			name:     "number with many decimal places",
			input:    1.23456789,
			expected: "1.23456789",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			args := []runtime.Value{runtime.NewNumericValue(tc.input)}
			result, err := fn.Call(args, env)

			require.NoError(t, err)
			assert.Equal(t, runtime.StringValue, result.Type)
			assert.Equal(t, tc.expected, result.StrValue)
		})
	}
}

func TestStrFunction_ErrorCases(t *testing.T) {
	env := runtime.NewEnvironment()
	fn := GetBuiltinFunction("STR$")
	require.NotNil(t, fn)

	t.Run("wrong argument count - none", func(t *testing.T) {
		args := []runtime.Value{}
		_, err := fn.Call(args, env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected 1 argument")
	})

	t.Run("wrong argument count - too many", func(t *testing.T) {
		args := []runtime.Value{
			runtime.NewNumericValue(42),
			runtime.NewNumericValue(24),
		}
		_, err := fn.Call(args, env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected 1 argument")
	})

	t.Run("wrong argument type - string", func(t *testing.T) {
		args := []runtime.Value{runtime.NewStringValue("42")}
		_, err := fn.Call(args, env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "argument must be numeric")
	})
}

// Test VAL function implementation
func TestValFunction_Call(t *testing.T) {
	env := runtime.NewEnvironment()
	fn := GetBuiltinFunction("VAL")
	require.NotNil(t, fn)

	testCases := []struct {
		name     string
		input    string
		expected float64
	}{
		{
			name:     "positive integer string",
			input:    "42",
			expected: 42.0,
		},
		{
			name:     "negative integer string",
			input:    "-42",
			expected: -42.0,
		},
		{
			name:     "zero string",
			input:    "0",
			expected: 0.0,
		},
		{
			name:     "positive decimal string",
			input:    "3.14",
			expected: 3.14,
		},
		{
			name:     "negative decimal string",
			input:    "-3.14",
			expected: -3.14,
		},
		{
			name:     "string with leading spaces",
			input:    "  42",
			expected: 42.0,
		},
		{
			name:     "string with trailing spaces",
			input:    "42  ",
			expected: 42.0,
		},
		{
			name:     "string with leading and trailing spaces",
			input:    "  42  ",
			expected: 42.0,
		},
		{
			name:     "very small number string",
			input:    "0.001",
			expected: 0.001,
		},
		{
			name:     "very large number string",
			input:    "999999",
			expected: 999999.0,
		},
		{
			name:     "scientific notation",
			input:    "1e3",
			expected: 1000.0,
		},
		{
			name:     "negative scientific notation",
			input:    "-1e3",
			expected: -1000.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			args := []runtime.Value{runtime.NewStringValue(tc.input)}
			result, err := fn.Call(args, env)

			require.NoError(t, err)
			assert.Equal(t, runtime.NumericValue, result.Type)
			assert.Equal(t, tc.expected, result.NumValue)
		})
	}
}

func TestValFunction_InvalidConversions(t *testing.T) {
	env := runtime.NewEnvironment()
	fn := GetBuiltinFunction("VAL")
	require.NotNil(t, fn)

	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "empty string",
			input: "",
		},
		{
			name:  "non-numeric string",
			input: "hello",
		},
		{
			name:  "string with letters",
			input: "abc123",
		},
		{
			name:  "string with special characters",
			input: "!@#$%",
		},
		{
			name:  "mixed alphanumeric",
			input: "12abc34",
		},
		{
			name:  "string with only spaces",
			input: "   ",
		},
		{
			name:  "invalid decimal format",
			input: "3.14.15",
		},
		{
			name:  "multiple signs",
			input: "--42",
		},
		{
			name:  "sign in wrong position",
			input: "4-2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			args := []runtime.Value{runtime.NewStringValue(tc.input)}
			_, err := fn.Call(args, env)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "cannot convert")
		})
	}
}

func TestValFunction_ErrorCases(t *testing.T) {
	env := runtime.NewEnvironment()
	fn := GetBuiltinFunction("VAL")
	require.NotNil(t, fn)

	t.Run("wrong argument count - none", func(t *testing.T) {
		args := []runtime.Value{}
		_, err := fn.Call(args, env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected 1 argument")
	})

	t.Run("wrong argument count - too many", func(t *testing.T) {
		args := []runtime.Value{
			runtime.NewStringValue("42"),
			runtime.NewStringValue("24"),
		}
		_, err := fn.Call(args, env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected 1 argument")
	})

	t.Run("wrong argument type - numeric", func(t *testing.T) {
		args := []runtime.Value{runtime.NewNumericValue(42)}
		_, err := fn.Call(args, env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "argument must be string")
	})
}

// Test string function argument validation and type checking
func TestStringFunctionArgumentValidation(t *testing.T) {
	env := runtime.NewEnvironment()

	testCases := []struct {
		functionName string
		expectedArgs int
		validArgs    []runtime.Value
		invalidArgs  [][]runtime.Value
	}{
		{
			functionName: "LEN",
			expectedArgs: 1,
			validArgs:    []runtime.Value{runtime.NewStringValue("test")},
			invalidArgs: [][]runtime.Value{
				{}, // no args
				{runtime.NewStringValue("test"), runtime.NewStringValue("extra")}, // too many
				{runtime.NewNumericValue(123)}, // wrong type
			},
		},
		{
			functionName: "MID$",
			expectedArgs: 3,
			validArgs: []runtime.Value{
				runtime.NewStringValue("test"),
				runtime.NewNumericValue(1),
				runtime.NewNumericValue(2),
			},
			invalidArgs: [][]runtime.Value{
				{}, // no args
				{runtime.NewStringValue("test")}, // too few
				{runtime.NewStringValue("test"), runtime.NewNumericValue(1)}, // too few
				{runtime.NewStringValue("test"), runtime.NewNumericValue(1), runtime.NewNumericValue(2), runtime.NewNumericValue(3)}, // too many
			},
		},
		{
			functionName: "STR$",
			expectedArgs: 1,
			validArgs:    []runtime.Value{runtime.NewNumericValue(123)},
			invalidArgs: [][]runtime.Value{
				{}, // no args
				{runtime.NewNumericValue(123), runtime.NewNumericValue(456)}, // too many
				{runtime.NewStringValue("123")}, // wrong type
			},
		},
		{
			functionName: "VAL",
			expectedArgs: 1,
			validArgs:    []runtime.Value{runtime.NewStringValue("123")},
			invalidArgs: [][]runtime.Value{
				{}, // no args
				{runtime.NewStringValue("123"), runtime.NewStringValue("456")}, // too many
				{runtime.NewNumericValue(123)}, // wrong type
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.functionName+" argument validation", func(t *testing.T) {
			fn := GetBuiltinFunction(tc.functionName)
			require.NotNil(t, fn)

			// Test correct argument count
			assert.Equal(t, tc.expectedArgs, fn.ArgCount())

			// Test valid arguments
			_, err := fn.Call(tc.validArgs, env)
			if tc.functionName == "VAL" && tc.validArgs[0].StrValue == "123" {
				// VAL with valid string should work
				assert.NoError(t, err)
			} else if tc.functionName != "VAL" {
				// Other functions should work with valid args
				assert.NoError(t, err)
			}

			// Test invalid arguments
			for i, invalidArgs := range tc.invalidArgs {
				_, err := fn.Call(invalidArgs, env)
				assert.Error(t, err, "invalid args case %d should fail", i)
			}
		})
	}
}

// Test edge cases for string manipulation and conversion logic
func TestStringFunctionEdgeCases(t *testing.T) {
	env := runtime.NewEnvironment()

	t.Run("LEN with unicode characters", func(t *testing.T) {
		fn := GetBuiltinFunction("LEN")
		require.NotNil(t, fn)

		args := []runtime.Value{runtime.NewStringValue("Hello 世界")}
		result, err := fn.Call(args, env)

		require.NoError(t, err)
		assert.Equal(t, runtime.NumericValue, result.Type)
		// Length should count bytes, not runes in basic implementation
		assert.Equal(t, float64(len("Hello 世界")), result.NumValue)
	})

	t.Run("MID$ with unicode characters", func(t *testing.T) {
		fn := GetBuiltinFunction("MID$")
		require.NotNil(t, fn)

		args := []runtime.Value{
			runtime.NewStringValue("Hello 世界"),
			runtime.NewNumericValue(1.0),
			runtime.NewNumericValue(5.0),
		}
		result, err := fn.Call(args, env)

		require.NoError(t, err)
		assert.Equal(t, runtime.StringValue, result.Type)
		// Should extract first 5 bytes
		expected := "Hello 世界"[:5]
		assert.Equal(t, expected, result.StrValue)
	})

	t.Run("STR$ with special numeric values", func(t *testing.T) {
		fn := GetBuiltinFunction("STR$")
		require.NotNil(t, fn)

		testCases := []struct {
			name     string
			input    float64
			expected string
		}{
			{"positive zero", 0.0, "0"},
			{"negative zero", -0.0, "0"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				args := []runtime.Value{runtime.NewNumericValue(tc.input)}
				result, err := fn.Call(args, env)

				require.NoError(t, err)
				assert.Equal(t, runtime.StringValue, result.Type)
				assert.Equal(t, tc.expected, result.StrValue)
			})
		}
	})
}