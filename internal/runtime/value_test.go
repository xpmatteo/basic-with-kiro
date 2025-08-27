package runtime

import (
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