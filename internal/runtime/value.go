package runtime

import (
	"fmt"
	"strconv"
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
		if val, err := strconv.ParseFloat(v.StrValue, 64); err == nil {
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