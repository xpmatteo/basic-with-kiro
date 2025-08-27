package runtime

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNumericValue(t *testing.T) {
	val := NewNumericValue(42.5)
	assert.Equal(t, NumericValue, val.Type)
	assert.Equal(t, 42.5, val.NumValue)
}

func TestNewStringValue(t *testing.T) {
	val := NewStringValue("hello")
	assert.Equal(t, StringValue, val.Type)
	assert.Equal(t, "hello", val.StrValue)
}

func TestValueString(t *testing.T) {
	numVal := NewNumericValue(42)
	assert.Equal(t, "42", numVal.String())

	strVal := NewStringValue("test")
	assert.Equal(t, "test", strVal.String())
}

func TestValueToNumber(t *testing.T) {
	// Test numeric value
	numVal := NewNumericValue(42.5)
	result, err := numVal.ToNumber()
	require.NoError(t, err)
	assert.Equal(t, 42.5, result)

	// Test string value that can be converted
	strVal := NewStringValue("123.45")
	result, err = strVal.ToNumber()
	require.NoError(t, err)
	assert.Equal(t, 123.45, result)

	// Test string value that cannot be converted
	invalidStrVal := NewStringValue("hello")
	_, err = invalidStrVal.ToNumber()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot convert string 'hello' to number")
}

func TestValueToString(t *testing.T) {
	// Test numeric value conversion
	numVal := NewNumericValue(42.5)
	assert.Equal(t, "42.5", numVal.ToString())

	// Test string value conversion
	strVal := NewStringValue("test")
	assert.Equal(t, "test", strVal.ToString())

	// Test integer conversion (should not show decimal)
	intVal := NewNumericValue(42)
	assert.Equal(t, "42", intVal.ToString())
}

// Additional tests for comprehensive value system coverage as required by task 3.1

func TestValueTypeConversions_EdgeCases(t *testing.T) {
	t.Run("Empty string to number", func(t *testing.T) {
		emptyStr := NewStringValue("")
		_, err := emptyStr.ToNumber()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot convert string '' to number")
	})

	t.Run("Zero value conversions", func(t *testing.T) {
		zeroNum := NewNumericValue(0)
		result, err := zeroNum.ToNumber()
		require.NoError(t, err)
		assert.Equal(t, 0.0, result)
		assert.Equal(t, "0", zeroNum.ToString())

		zeroStr := NewStringValue("0")
		result, err = zeroStr.ToNumber()
		require.NoError(t, err)
		assert.Equal(t, 0.0, result)
	})

	t.Run("Whitespace string conversions", func(t *testing.T) {
		whitespaceStr := NewStringValue("  123.45  ")
		result, err := whitespaceStr.ToNumber()
		require.NoError(t, err)
		assert.Equal(t, 123.45, result)
	})

	t.Run("Invalid numeric strings", func(t *testing.T) {
		testCases := []string{"abc", "12.34.56", "1e", "++123", "--456"}
		for _, tc := range testCases {
			invalidStr := NewStringValue(tc)
			_, err := invalidStr.ToNumber()
			assert.Error(t, err, "Expected error for string: %s", tc)
		}
	})

	t.Run("Special numeric values", func(t *testing.T) {
		// Test infinity
		infVal := NewNumericValue(math.Inf(1))
		assert.Equal(t, "+Inf", infVal.ToString())

		// Test negative infinity
		negInfVal := NewNumericValue(math.Inf(-1))
		assert.Equal(t, "-Inf", negInfVal.ToString())

		// Test NaN
		nanVal := NewNumericValue(math.NaN())
		assert.Equal(t, "NaN", nanVal.ToString())
	})
}

func TestValueComparison(t *testing.T) {
	t.Run("Numeric equality", func(t *testing.T) {
		val1 := NewNumericValue(42.5)
		val2 := NewNumericValue(42.5)
		val3 := NewNumericValue(43.0)

		assert.True(t, val1.Equals(val2))
		assert.False(t, val1.Equals(val3))
	})

	t.Run("String equality", func(t *testing.T) {
		val1 := NewStringValue("hello")
		val2 := NewStringValue("hello")
		val3 := NewStringValue("world")

		assert.True(t, val1.Equals(val2))
		assert.False(t, val1.Equals(val3))
	})

	t.Run("Cross-type equality", func(t *testing.T) {
		numVal := NewNumericValue(42)
		strVal := NewStringValue("42")

		// Should be equal when converted
		assert.True(t, numVal.Equals(strVal))
	})

	t.Run("Numeric comparison", func(t *testing.T) {
		val1 := NewNumericValue(10)
		val2 := NewNumericValue(20)
		val3 := NewNumericValue(10)

		assert.Equal(t, -1, val1.Compare(val2)) // 10 < 20
		assert.Equal(t, 1, val2.Compare(val1))  // 20 > 10
		assert.Equal(t, 0, val1.Compare(val3))  // 10 == 10
	})

	t.Run("String comparison", func(t *testing.T) {
		val1 := NewStringValue("apple")
		val2 := NewStringValue("banana")
		val3 := NewStringValue("apple")

		assert.Equal(t, -1, val1.Compare(val2)) // "apple" < "banana"
		assert.Equal(t, 1, val2.Compare(val1))  // "banana" > "apple"
		assert.Equal(t, 0, val1.Compare(val3))  // "apple" == "apple"
	})
}

func TestValueArithmetic(t *testing.T) {
	t.Run("Addition", func(t *testing.T) {
		val1 := NewNumericValue(10)
		val2 := NewNumericValue(5)
		result, err := val1.Add(val2)
		require.NoError(t, err)
		assert.Equal(t, 15.0, result.NumValue)
		assert.Equal(t, NumericValue, result.Type)

		// String concatenation
		str1 := NewStringValue("Hello")
		str2 := NewStringValue(" World")
		result, err = str1.Add(str2)
		require.NoError(t, err)
		assert.Equal(t, "Hello World", result.StrValue)
		assert.Equal(t, StringValue, result.Type)
	})

	t.Run("Subtraction", func(t *testing.T) {
		val1 := NewNumericValue(10)
		val2 := NewNumericValue(3)
		result, err := val1.Subtract(val2)
		require.NoError(t, err)
		assert.Equal(t, 7.0, result.NumValue)
		assert.Equal(t, NumericValue, result.Type)

		// String subtraction should error
		str1 := NewStringValue("Hello")
		str2 := NewStringValue("World")
		_, err = str1.Subtract(str2)
		assert.Error(t, err)
	})

	t.Run("Multiplication", func(t *testing.T) {
		val1 := NewNumericValue(4)
		val2 := NewNumericValue(3)
		result, err := val1.Multiply(val2)
		require.NoError(t, err)
		assert.Equal(t, 12.0, result.NumValue)
		assert.Equal(t, NumericValue, result.Type)

		// String multiplication should error
		str1 := NewStringValue("Hello")
		str2 := NewStringValue("World")
		_, err = str1.Multiply(str2)
		assert.Error(t, err)
	})

	t.Run("Division", func(t *testing.T) {
		val1 := NewNumericValue(15)
		val2 := NewNumericValue(3)
		result, err := val1.Divide(val2)
		require.NoError(t, err)
		assert.Equal(t, 5.0, result.NumValue)
		assert.Equal(t, NumericValue, result.Type)

		// Division by zero
		val3 := NewNumericValue(0)
		_, err = val1.Divide(val3)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "division by zero")
	})

	t.Run("Power", func(t *testing.T) {
		val1 := NewNumericValue(2)
		val2 := NewNumericValue(3)
		result, err := val1.Power(val2)
		require.NoError(t, err)
		assert.Equal(t, 8.0, result.NumValue)
		assert.Equal(t, NumericValue, result.Type)

		// String power should error
		str1 := NewStringValue("Hello")
		str2 := NewStringValue("World")
		_, err = str1.Power(str2)
		assert.Error(t, err)
	})

	t.Run("Mixed type arithmetic", func(t *testing.T) {
		numVal := NewNumericValue(10)
		strVal := NewStringValue("5")

		// Should convert string to number for arithmetic
		result, err := numVal.Add(strVal)
		require.NoError(t, err)
		assert.Equal(t, 15.0, result.NumValue)
		assert.Equal(t, NumericValue, result.Type)

		// Invalid string conversion should error
		invalidStr := NewStringValue("hello")
		_, err = numVal.Add(invalidStr)
		assert.Error(t, err)
	})
}