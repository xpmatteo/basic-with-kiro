package ast

import (
	"basic-interpreter/internal/runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test for FunctionCallExpression struct and BuiltinFunction interface
func TestFunctionCallExpression_Evaluate(t *testing.T) {
	env := runtime.NewEnvironment()

	t.Run("function call with no arguments", func(t *testing.T) {
		// Test RND function which takes no arguments
		expr := &FunctionCallExpression{
			Name: "RND",
			Args: []Expression{},
		}

		result, err := expr.Evaluate(env)

		require.NoError(t, err)
		assert.Equal(t, runtime.NumericValue, result.Type)
		assert.GreaterOrEqual(t, result.NumValue, 0.0)
		assert.Less(t, result.NumValue, 1.0)
	})

	t.Run("function call with single argument", func(t *testing.T) {
		// Test ABS function with single numeric argument
		expr := &FunctionCallExpression{
			Name: "ABS",
			Args: []Expression{
				NewLiteralExpression(runtime.NewNumericValue(-5.5)),
			},
		}

		result, err := expr.Evaluate(env)

		require.NoError(t, err)
		assert.Equal(t, runtime.NumericValue, result.Type)
		assert.Equal(t, 5.5, result.NumValue)
	})

	t.Run("function call with multiple arguments", func(t *testing.T) {
		// Test MID$ function with multiple arguments
		expr := &FunctionCallExpression{
			Name: "MID$",
			Args: []Expression{
				NewLiteralExpression(runtime.NewStringValue("HELLO")),
				NewLiteralExpression(runtime.NewNumericValue(2)),
				NewLiteralExpression(runtime.NewNumericValue(3)),
			},
		}

		result, err := expr.Evaluate(env)

		require.NoError(t, err)
		assert.Equal(t, runtime.StringValue, result.Type)
		assert.Equal(t, "ELL", result.StrValue)
	})

	t.Run("function call with variable arguments", func(t *testing.T) {
		// Set up variables
		env.SetVariable("X", runtime.NewNumericValue(-10))
		
		expr := &FunctionCallExpression{
			Name: "ABS",
			Args: []Expression{
				NewVariableExpression("X"),
			},
		}

		result, err := expr.Evaluate(env)

		require.NoError(t, err)
		assert.Equal(t, runtime.NumericValue, result.Type)
		assert.Equal(t, 10.0, result.NumValue)
	})

	t.Run("function call with expression arguments", func(t *testing.T) {
		// Test function call with complex expression as argument
		expr := &FunctionCallExpression{
			Name: "ABS",
			Args: []Expression{
				NewBinaryExpression(
					NewLiteralExpression(runtime.NewNumericValue(3)),
					OpSubtract,
					NewLiteralExpression(runtime.NewNumericValue(8)),
				),
			},
		}

		result, err := expr.Evaluate(env)

		require.NoError(t, err)
		assert.Equal(t, runtime.NumericValue, result.Type)
		assert.Equal(t, 5.0, result.NumValue)
	})
}

func TestFunctionCallExpression_ErrorCases(t *testing.T) {
	env := runtime.NewEnvironment()

	t.Run("unknown function", func(t *testing.T) {
		expr := &FunctionCallExpression{
			Name: "UNKNOWN",
			Args: []Expression{},
		}

		_, err := expr.Evaluate(env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown function: UNKNOWN")
	})

	t.Run("wrong argument count - too few", func(t *testing.T) {
		expr := &FunctionCallExpression{
			Name: "ABS",
			Args: []Expression{}, // ABS requires 1 argument
		}

		_, err := expr.Evaluate(env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "function ABS expects 1 argument(s), got 0")
	})

	t.Run("wrong argument count - too many", func(t *testing.T) {
		expr := &FunctionCallExpression{
			Name: "ABS",
			Args: []Expression{
				NewLiteralExpression(runtime.NewNumericValue(5)),
				NewLiteralExpression(runtime.NewNumericValue(3)),
			}, // ABS requires only 1 argument
		}

		_, err := expr.Evaluate(env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "function ABS expects 1 argument(s), got 2")
	})

	t.Run("error evaluating argument", func(t *testing.T) {
		// Create an argument expression that will fail
		badExpr := NewBinaryExpression(
			NewLiteralExpression(runtime.NewStringValue("hello")),
			OpDivide,
			NewLiteralExpression(runtime.NewNumericValue(0)),
		)

		expr := &FunctionCallExpression{
			Name: "ABS",
			Args: []Expression{badExpr},
		}

		_, err := expr.Evaluate(env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error evaluating argument 0")
	})

	t.Run("function execution error", func(t *testing.T) {
		// Test VAL function with invalid string
		expr := &FunctionCallExpression{
			Name: "VAL",
			Args: []Expression{
				NewLiteralExpression(runtime.NewStringValue("not_a_number")),
			},
		}

		_, err := expr.Evaluate(env)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error calling function VAL")
	})
}

// Test for function registry and built-in function lookup
func TestFunctionRegistry(t *testing.T) {
	t.Run("get registered function", func(t *testing.T) {
		fn := GetBuiltinFunction("ABS")
		
		require.NotNil(t, fn)
		assert.Equal(t, "ABS", fn.Name())
		assert.Equal(t, 1, fn.ArgCount())
	})

	t.Run("get unregistered function", func(t *testing.T) {
		fn := GetBuiltinFunction("NONEXISTENT")
		
		assert.Nil(t, fn)
	})

	t.Run("case insensitive function lookup", func(t *testing.T) {
		testCases := []string{"abs", "ABS", "Abs", "aBs"}
		
		for _, name := range testCases {
			fn := GetBuiltinFunction(name)
			require.NotNil(t, fn, "function lookup should be case insensitive for %s", name)
			assert.Equal(t, "ABS", fn.Name())
		}
	})

	t.Run("all expected functions are registered", func(t *testing.T) {
		expectedFunctions := []string{
			"ABS", "INT", "RND", "LEN", "MID$", "STR$", "VAL",
		}
		
		for _, name := range expectedFunctions {
			fn := GetBuiltinFunction(name)
			assert.NotNil(t, fn, "function %s should be registered", name)
		}
	})
}

// Test for BuiltinFunction interface implementation
func TestBuiltinFunctionInterface(t *testing.T) {
	env := runtime.NewEnvironment()
	
	t.Run("ABS function interface", func(t *testing.T) {
		fn := GetBuiltinFunction("ABS")
		require.NotNil(t, fn)
		
		assert.Equal(t, "ABS", fn.Name())
		assert.Equal(t, 1, fn.ArgCount())
		
		// Test function call
		args := []runtime.Value{runtime.NewNumericValue(-5)}
		result, err := fn.Call(args, env)
		
		require.NoError(t, err)
		assert.Equal(t, runtime.NumericValue, result.Type)
		assert.Equal(t, 5.0, result.NumValue)
	})

	t.Run("RND function interface", func(t *testing.T) {
		fn := GetBuiltinFunction("RND")
		require.NotNil(t, fn)
		
		assert.Equal(t, "RND", fn.Name())
		assert.Equal(t, 0, fn.ArgCount())
		
		// Test function call
		args := []runtime.Value{}
		result, err := fn.Call(args, env)
		
		require.NoError(t, err)
		assert.Equal(t, runtime.NumericValue, result.Type)
		assert.GreaterOrEqual(t, result.NumValue, 0.0)
		assert.Less(t, result.NumValue, 1.0)
	})

	t.Run("LEN function interface", func(t *testing.T) {
		fn := GetBuiltinFunction("LEN")
		require.NotNil(t, fn)
		
		assert.Equal(t, "LEN", fn.Name())
		assert.Equal(t, 1, fn.ArgCount())
		
		// Test function call
		args := []runtime.Value{runtime.NewStringValue("HELLO")}
		result, err := fn.Call(args, env)
		
		require.NoError(t, err)
		assert.Equal(t, runtime.NumericValue, result.Type)
		assert.Equal(t, 5.0, result.NumValue)
	})
}

// Test for function argument validation
func TestFunctionArgumentValidation(t *testing.T) {
	env := runtime.NewEnvironment()
	
	t.Run("validate argument count", func(t *testing.T) {
		fn := GetBuiltinFunction("ABS")
		require.NotNil(t, fn)
		
		// Test with correct argument count
		args := []runtime.Value{runtime.NewNumericValue(5)}
		_, err := fn.Call(args, env)
		assert.NoError(t, err)
		
		// Test with wrong argument count
		wrongArgs := []runtime.Value{}
		_, err = fn.Call(wrongArgs, env)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected 1 argument")
	})

	t.Run("validate argument types", func(t *testing.T) {
		fn := GetBuiltinFunction("ABS")
		require.NotNil(t, fn)
		
		// Test with wrong argument type
		args := []runtime.Value{runtime.NewStringValue("hello")}
		_, err := fn.Call(args, env)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "argument must be numeric")
	})
}

// Test helper function for creating function call expressions
func TestNewFunctionCallExpression(t *testing.T) {
	t.Run("create function call expression", func(t *testing.T) {
		args := []Expression{
			NewLiteralExpression(runtime.NewNumericValue(5)),
		}
		
		expr := NewFunctionCallExpression("ABS", args)
		
		assert.NotNil(t, expr)
		assert.Equal(t, "ABS", expr.Name)
		assert.Equal(t, args, expr.Args)
	})
}