package runtime

import (
	"testing"
)

func TestNewEnvironment(t *testing.T) {
	env := NewEnvironment()
	if env == nil {
		t.Error("Expected non-nil environment")
	}
	if env.Variables == nil {
		t.Error("Expected initialized Variables map")
	}
	if env.ProgramCounter != 0 {
		t.Errorf("Expected ProgramCounter to be 0, got %d", env.ProgramCounter)
	}
}

func TestVariableOperations(t *testing.T) {
	env := NewEnvironment()

	// Test setting and getting a numeric variable
	numVal := NewNumericValue(42)
	env.SetVariable("X", numVal)
	
	result := env.GetVariable("X")
	if result.Type != NumericValue || result.NumValue != 42 {
		t.Errorf("Expected numeric value 42, got %v", result)
	}

	// Test case insensitivity
	result = env.GetVariable("x")
	if result.Type != NumericValue || result.NumValue != 42 {
		t.Errorf("Expected case-insensitive access to work, got %v", result)
	}

	// Test string variable
	strVal := NewStringValue("hello")
	env.SetVariable("NAME$", strVal)
	
	result = env.GetVariable("name$")
	if result.Type != StringValue || result.StrValue != "hello" {
		t.Errorf("Expected string value 'hello', got %v", result)
	}
}

func TestUndefinedVariables(t *testing.T) {
	env := NewEnvironment()

	// Test undefined numeric variable (should default to 0)
	result := env.GetVariable("UNDEFINED")
	if result.Type != NumericValue || result.NumValue != 0 {
		t.Errorf("Expected undefined numeric variable to be 0, got %v", result)
	}

	// Test undefined string variable (should default to empty string)
	result = env.GetVariable("UNDEFINED$")
	if result.Type != StringValue || result.StrValue != "" {
		t.Errorf("Expected undefined string variable to be empty, got %v", result)
	}
}

func TestRandom(t *testing.T) {
	env := NewEnvironment()
	
	// Test that random returns a value between 0 and 1
	for i := 0; i < 10; i++ {
		val := env.Random()
		if val < 0 || val >= 1 {
			t.Errorf("Expected random value between 0 and 1, got %v", val)
		}
	}
}