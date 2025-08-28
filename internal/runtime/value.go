package runtime

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ValueType represents the type of a value in BASIC
type ValueType int

const (
	NumericValue ValueType = iota
	StringValue
)

// Value represents a value in the BASIC interpreter
type Value struct {
	Type     ValueType
	NumValue float64
	StrValue string
}

// NewNumericValue creates a new numeric value
func NewNumericValue(val float64) Value {
	return Value{
		Type:     NumericValue,
		NumValue: val,
	}
}

// NewStringValue creates a new string value
func NewStringValue(val string) Value {
	return Value{
		Type:     StringValue,
		StrValue: val,
	}
}

// String returns the string representation of the value
func (v Value) String() string {
	switch v.Type {
	case NumericValue:
		return fmt.Sprintf("%g", v.NumValue)
	case StringValue:
		return v.StrValue
	default:
		return ""
	}
}

// ToNumber converts the value to a numeric value
func (v Value) ToNumber() (float64, error) {
	switch v.Type {
	case NumericValue:
		return v.NumValue, nil
	case StringValue:
		// Trim whitespace before parsing
		trimmed := strings.TrimSpace(v.StrValue)
		if val, err := strconv.ParseFloat(trimmed, 64); err == nil {
			return val, nil
		}
		return 0, fmt.Errorf("cannot convert string '%s' to number", v.StrValue)
	default:
		return 0, fmt.Errorf("unknown value type")
	}
}

// ToString converts the value to a string value
func (v Value) ToString() string {
	switch v.Type {
	case NumericValue:
		return fmt.Sprintf("%g", v.NumValue)
	case StringValue:
		return v.StrValue
	default:
		return ""
	}
}

// Equals compares two values for equality
func (v Value) Equals(other Value) bool {
	// If both are same type, compare directly
	if v.Type == other.Type {
		switch v.Type {
		case NumericValue:
			return v.NumValue == other.NumValue
		case StringValue:
			return v.StrValue == other.StrValue
		}
	}

	// Cross-type comparison: try to convert both to numbers
	vNum, vErr := v.ToNumber()
	otherNum, otherErr := other.ToNumber()
	if vErr == nil && otherErr == nil {
		return vNum == otherNum
	}

	// If number conversion fails, compare as strings
	return v.ToString() == other.ToString()
}

// Compare compares two values and returns -1, 0, or 1
func (v Value) Compare(other Value) int {
	// If both are numeric, compare as numbers
	if v.Type == NumericValue && other.Type == NumericValue {
		return compareFloat64(v.NumValue, other.NumValue)
	}

	// Try to convert both to numbers for comparison
	vNum, vErr := v.ToNumber()
	otherNum, otherErr := other.ToNumber()
	if vErr == nil && otherErr == nil {
		return compareFloat64(vNum, otherNum)
	}

	// Compare as strings
	vStr := v.ToString()
	otherStr := other.ToString()
	return compareString(vStr, otherStr)
}

// compareFloat64 compares two float64 values
func compareFloat64(a, b float64) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

// compareString compares two string values
func compareString(a, b string) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

// Add performs addition operation
func (v Value) Add(other Value) (Value, error) {
	// If both are explicitly strings, do string concatenation
	if v.Type == StringValue && other.Type == StringValue {
		return NewStringValue(v.ToString() + other.ToString()), nil
	}

	// Try numeric addition
	vNum, vErr := v.ToNumber()
	otherNum, otherErr := other.ToNumber()
	
	// If both can be converted to numbers, do numeric addition
	if vErr == nil && otherErr == nil {
		return NewNumericValue(vNum + otherNum), nil
	}

	// If one is a number and the other can't be converted, return error
	if (v.Type == NumericValue && otherErr != nil) || (other.Type == NumericValue && vErr != nil) {
		if vErr != nil {
			return Value{}, vErr
		}
		return Value{}, otherErr
	}

	// Otherwise, do string concatenation
	return NewStringValue(v.ToString() + other.ToString()), nil
}

// Subtract performs subtraction operation
func (v Value) Subtract(other Value) (Value, error) {
	return v.performNumericOperation(other, "subtract", func(a, b float64) (float64, error) {
		return a - b, nil
	})
}

// Multiply performs multiplication operation
func (v Value) Multiply(other Value) (Value, error) {
	return v.performNumericOperation(other, "multiply", func(a, b float64) (float64, error) {
		return a * b, nil
	})
}

// Divide performs division operation
func (v Value) Divide(other Value) (Value, error) {
	return v.performNumericOperation(other, "divide", func(a, b float64) (float64, error) {
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	})
}

// Power performs exponentiation operation
func (v Value) Power(other Value) (Value, error) {
	return v.performNumericOperation(other, "raise strings to power", func(a, b float64) (float64, error) {
		return math.Pow(a, b), nil
	})
}

// performNumericOperation is a helper function for numeric operations
func (v Value) performNumericOperation(other Value, operation string, op func(float64, float64) (float64, error)) (Value, error) {
	// Only numeric operations are allowed
	if v.Type == StringValue || other.Type == StringValue {
		return Value{}, fmt.Errorf("cannot %s strings", operation)
	}

	vNum, err := v.ToNumber()
	if err != nil {
		return Value{}, err
	}
	otherNum, err := other.ToNumber()
	if err != nil {
		return Value{}, err
	}

	result, err := op(vNum, otherNum)
	if err != nil {
		return Value{}, err
	}

	return NewNumericValue(result), nil
}