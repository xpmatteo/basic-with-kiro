package runtime

import (
	"testing"
)

func TestNewNumericValue(t *testing.T) {
	val := NewNumericValue(42.5)
	if val.Type != NumericValue {
		t.Errorf("Expected NumericValue, got %v", val.Type)
	}
	if val.NumValue != 42.5 {
		t.Errorf("Expected 42.5, got %v", val.NumValue)
	}
}

func TestNewStringValue(t *testing.T) {
	val := NewStringValue("hello")
	if val.Type != StringValue {
		t.Errorf("Expected StringValue, got %v", val.Type)
	}
	if val.StrValue != "hello" {
		t.Errorf("Expected 'hello', got %v", val.StrValue)
	}
}

func TestValueString(t *testing.T) {
	numVal := NewNumericValue(42)
	if numVal.String() != "42" {
		t.Errorf("Expected '42', got '%s'", numVal.String())
	}

	strVal := NewStringValue("test")
	if strVal.String() != "test" {
		t.Errorf("Expected 'test', got '%s'", strVal.String())
	}
}

func TestValueToNumber(t *testing.T) {
	// Test numeric value
	numVal := NewNumericValue(42.5)
	result, err := numVal.ToNumber()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 42.5 {
		t.Errorf("Expected 42.5, got %v", result)
	}

	// Test string value that can be converted
	strVal := NewStringValue("123.45")
	result, err = strVal.ToNumber()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 123.45 {
		t.Errorf("Expected 123.45, got %v", result)
	}

	// Test string value that cannot be converted
	invalidStrVal := NewStringValue("hello")
	_, err = invalidStrVal.ToNumber()
	if err == nil {
		t.Error("Expected error for invalid string conversion")
	}
}