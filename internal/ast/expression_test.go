package ast

import (
	"basic-interpreter/internal/runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLiteralExpression_Evaluate(t *testing.T) {
	env := runtime.NewEnvironment()

	t.Run("numeric literal", func(t *testing.T) {
		expr := &LiteralExpression{Value: runtime.NewNumericValue(42.5)}
		result, err := expr.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, runtime.NumericValue, result.Type)
		assert.Equal(t, 42.5, result.NumValue)
	})

	t.Run("string literal", func(t *testing.T) {
		expr := &LiteralExpression{Value: runtime.NewStringValue("hello")}
		result, err := expr.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, runtime.StringValue, result.Type)
		assert.Equal(t, "hello", result.StrValue)
	})
}

func TestVariableExpression_Evaluate(t *testing.T) {
	env := runtime.NewEnvironment()

	t.Run("existing numeric variable", func(t *testing.T) {
		env.SetVariable("X", runtime.NewNumericValue(10))
		expr := &VariableExpression{Name: "X"}
		
		result, err := expr.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, runtime.NumericValue, result.Type)
		assert.Equal(t, 10.0, result.NumValue)
	})

	t.Run("existing string variable", func(t *testing.T) {
		env.SetVariable("NAME$", runtime.NewStringValue("test"))
		expr := &VariableExpression{Name: "NAME$"}
		
		result, err := expr.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, runtime.StringValue, result.Type)
		assert.Equal(t, "test", result.StrValue)
	})

	t.Run("uninitialized numeric variable defaults to zero", func(t *testing.T) {
		expr := &VariableExpression{Name: "Y"}
		
		result, err := expr.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, runtime.NumericValue, result.Type)
		assert.Equal(t, 0.0, result.NumValue)
	})

	t.Run("uninitialized string variable defaults to empty string", func(t *testing.T) {
		expr := &VariableExpression{Name: "EMPTY$"}
		
		result, err := expr.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, runtime.StringValue, result.Type)
		assert.Equal(t, "", result.StrValue)
	})

	t.Run("case insensitive variable names", func(t *testing.T) {
		env.SetVariable("test", runtime.NewNumericValue(5))
		expr := &VariableExpression{Name: "TEST"}
		
		result, err := expr.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, 5.0, result.NumValue)
	})
}

func TestBinaryExpression_Addition(t *testing.T) {
	env := runtime.NewEnvironment()

	t.Run("numeric addition", func(t *testing.T) {
		left := &LiteralExpression{Value: runtime.NewNumericValue(5)}
		right := &LiteralExpression{Value: runtime.NewNumericValue(3)}
		expr := &BinaryExpression{Left: left, Operator: "+", Right: right}
		
		result, err := expr.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, runtime.NumericValue, result.Type)
		assert.Equal(t, 8.0, result.NumValue)
	})

	t.Run("string concatenation", func(t *testing.T) {
		left := &LiteralExpression{Value: runtime.NewStringValue("hello")}
		right := &LiteralExpression{Value: runtime.NewStringValue(" world")}
		expr := &BinaryExpression{Left: left, Operator: "+", Right: right}
		
		result, err := expr.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, runtime.StringValue, result.Type)
		assert.Equal(t, "hello world", result.StrValue)
	})

	t.Run("mixed type addition - number and convertible string", func(t *testing.T) {
		left := &LiteralExpression{Value: runtime.NewNumericValue(5)}
		right := &LiteralExpression{Value: runtime.NewStringValue("3")}
		expr := &BinaryExpression{Left: left, Operator: "+", Right: right}
		
		result, err := expr.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, runtime.NumericValue, result.Type)
		assert.Equal(t, 8.0, result.NumValue)
	})
}

func TestBinaryExpression_Subtraction(t *testing.T) {
	env := runtime.NewEnvironment()

	t.Run("numeric subtraction", func(t *testing.T) {
		left := &LiteralExpression{Value: runtime.NewNumericValue(10)}
		right := &LiteralExpression{Value: runtime.NewNumericValue(3)}
		expr := &BinaryExpression{Left: left, Operator: "-", Right: right}
		
		result, err := expr.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, runtime.NumericValue, result.Type)
		assert.Equal(t, 7.0, result.NumValue)
	})

	t.Run("string subtraction should fail", func(t *testing.T) {
		left := &LiteralExpression{Value: runtime.NewStringValue("hello")}
		right := &LiteralExpression{Value: runtime.NewStringValue("world")}
		expr := &BinaryExpression{Left: left, Operator: "-", Right: right}
		
		_, err := expr.Evaluate(env)
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot subtract strings")
	})
}

func TestBinaryExpression_Multiplication(t *testing.T) {
	env := runtime.NewEnvironment()

	t.Run("numeric multiplication", func(t *testing.T) {
		left := &LiteralExpression{Value: runtime.NewNumericValue(4)}
		right := &LiteralExpression{Value: runtime.NewNumericValue(3)}
		expr := &BinaryExpression{Left: left, Operator: "*", Right: right}
		
		result, err := expr.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, runtime.NumericValue, result.Type)
		assert.Equal(t, 12.0, result.NumValue)
	})

	t.Run("string multiplication should fail", func(t *testing.T) {
		left := &LiteralExpression{Value: runtime.NewStringValue("hello")}
		right := &LiteralExpression{Value: runtime.NewNumericValue(3)}
		expr := &BinaryExpression{Left: left, Operator: "*", Right: right}
		
		_, err := expr.Evaluate(env)
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot multiply strings")
	})
}

func TestBinaryExpression_Division(t *testing.T) {
	env := runtime.NewEnvironment()

	t.Run("numeric division", func(t *testing.T) {
		left := &LiteralExpression{Value: runtime.NewNumericValue(15)}
		right := &LiteralExpression{Value: runtime.NewNumericValue(3)}
		expr := &BinaryExpression{Left: left, Operator: "/", Right: right}
		
		result, err := expr.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, runtime.NumericValue, result.Type)
		assert.Equal(t, 5.0, result.NumValue)
	})

	t.Run("division by zero", func(t *testing.T) {
		left := &LiteralExpression{Value: runtime.NewNumericValue(10)}
		right := &LiteralExpression{Value: runtime.NewNumericValue(0)}
		expr := &BinaryExpression{Left: left, Operator: "/", Right: right}
		
		_, err := expr.Evaluate(env)
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "division by zero")
	})
}

func TestBinaryExpression_Power(t *testing.T) {
	env := runtime.NewEnvironment()

	t.Run("numeric exponentiation", func(t *testing.T) {
		left := &LiteralExpression{Value: runtime.NewNumericValue(2)}
		right := &LiteralExpression{Value: runtime.NewNumericValue(3)}
		expr := &BinaryExpression{Left: left, Operator: "^", Right: right}
		
		result, err := expr.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, runtime.NumericValue, result.Type)
		assert.Equal(t, 8.0, result.NumValue)
	})

	t.Run("string power should fail", func(t *testing.T) {
		left := &LiteralExpression{Value: runtime.NewStringValue("hello")}
		right := &LiteralExpression{Value: runtime.NewNumericValue(2)}
		expr := &BinaryExpression{Left: left, Operator: "^", Right: right}
		
		_, err := expr.Evaluate(env)
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot raise strings to power")
	})
}

func TestBinaryExpression_OperatorPrecedence(t *testing.T) {
	env := runtime.NewEnvironment()

	t.Run("multiplication before addition", func(t *testing.T) {
		// 2 + 3 * 4 should be 2 + 12 = 14, not (2 + 3) * 4 = 20
		// This will be tested through parser, but we test the evaluation here
		
		// Create: 3 * 4
		mult := &BinaryExpression{
			Left:     &LiteralExpression{Value: runtime.NewNumericValue(3)},
			Operator: "*",
			Right:    &LiteralExpression{Value: runtime.NewNumericValue(4)},
		}
		
		// Create: 2 + (3 * 4)
		expr := &BinaryExpression{
			Left:     &LiteralExpression{Value: runtime.NewNumericValue(2)},
			Operator: "+",
			Right:    mult,
		}
		
		result, err := expr.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, 14.0, result.NumValue)
	})

	t.Run("power before multiplication", func(t *testing.T) {
		// 2 * 3 ^ 2 should be 2 * 9 = 18, not (2 * 3) ^ 2 = 36
		
		// Create: 3 ^ 2
		power := &BinaryExpression{
			Left:     &LiteralExpression{Value: runtime.NewNumericValue(3)},
			Operator: "^",
			Right:    &LiteralExpression{Value: runtime.NewNumericValue(2)},
		}
		
		// Create: 2 * (3 ^ 2)
		expr := &BinaryExpression{
			Left:     &LiteralExpression{Value: runtime.NewNumericValue(2)},
			Operator: "*",
			Right:    power,
		}
		
		result, err := expr.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, 18.0, result.NumValue)
	})
}

func TestBinaryExpression_WithVariables(t *testing.T) {
	env := runtime.NewEnvironment()
	env.SetVariable("X", runtime.NewNumericValue(5))
	env.SetVariable("Y", runtime.NewNumericValue(3))

	t.Run("variable addition", func(t *testing.T) {
		left := &VariableExpression{Name: "X"}
		right := &VariableExpression{Name: "Y"}
		expr := &BinaryExpression{Left: left, Operator: "+", Right: right}
		
		result, err := expr.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, 8.0, result.NumValue)
	})

	t.Run("variable and literal", func(t *testing.T) {
		left := &VariableExpression{Name: "X"}
		right := &LiteralExpression{Value: runtime.NewNumericValue(10)}
		expr := &BinaryExpression{Left: left, Operator: "*", Right: right}
		
		result, err := expr.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, 50.0, result.NumValue)
	})
}

func TestParenthesesExpression_Evaluate(t *testing.T) {
	env := runtime.NewEnvironment()

	t.Run("simple parentheses", func(t *testing.T) {
		inner := &LiteralExpression{Value: runtime.NewNumericValue(42)}
		expr := &ParenthesesExpression{Expression: inner}
		
		result, err := expr.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, 42.0, result.NumValue)
	})

	t.Run("parentheses change precedence", func(t *testing.T) {
		// (2 + 3) * 4 should be 5 * 4 = 20
		
		// Create: 2 + 3
		add := &BinaryExpression{
			Left:     &LiteralExpression{Value: runtime.NewNumericValue(2)},
			Operator: "+",
			Right:    &LiteralExpression{Value: runtime.NewNumericValue(3)},
		}
		
		// Create: (2 + 3)
		paren := &ParenthesesExpression{Expression: add}
		
		// Create: (2 + 3) * 4
		expr := &BinaryExpression{
			Left:     paren,
			Operator: "*",
			Right:    &LiteralExpression{Value: runtime.NewNumericValue(4)},
		}
		
		result, err := expr.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, 20.0, result.NumValue)
	})

	t.Run("nested parentheses", func(t *testing.T) {
		// ((5 + 3) * 2) should be (8 * 2) = 16
		
		// Create: 5 + 3
		add := &BinaryExpression{
			Left:     &LiteralExpression{Value: runtime.NewNumericValue(5)},
			Operator: "+",
			Right:    &LiteralExpression{Value: runtime.NewNumericValue(3)},
		}
		
		// Create: (5 + 3)
		innerParen := &ParenthesesExpression{Expression: add}
		
		// Create: (5 + 3) * 2
		mult := &BinaryExpression{
			Left:     innerParen,
			Operator: "*",
			Right:    &LiteralExpression{Value: runtime.NewNumericValue(2)},
		}
		
		// Create: ((5 + 3) * 2)
		outerParen := &ParenthesesExpression{Expression: mult}
		
		result, err := outerParen.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, 16.0, result.NumValue)
	})
}

func TestComplexNestedExpressions(t *testing.T) {
	env := runtime.NewEnvironment()
	env.SetVariable("A", runtime.NewNumericValue(2))
	env.SetVariable("B", runtime.NewNumericValue(3))
	env.SetVariable("C", runtime.NewNumericValue(4))

	t.Run("complex arithmetic expression", func(t *testing.T) {
		// A * (B + C) ^ 2 - 5 should be 2 * (3 + 4) ^ 2 - 5 = 2 * 49 - 5 = 93
		
		// Create: B + C
		add := &BinaryExpression{
			Left:     &VariableExpression{Name: "B"},
			Operator: "+",
			Right:    &VariableExpression{Name: "C"},
		}
		
		// Create: (B + C)
		paren := &ParenthesesExpression{Expression: add}
		
		// Create: (B + C) ^ 2
		power := &BinaryExpression{
			Left:     paren,
			Operator: "^",
			Right:    &LiteralExpression{Value: runtime.NewNumericValue(2)},
		}
		
		// Create: A * (B + C) ^ 2
		mult := &BinaryExpression{
			Left:     &VariableExpression{Name: "A"},
			Operator: "*",
			Right:    power,
		}
		
		// Create: A * (B + C) ^ 2 - 5
		expr := &BinaryExpression{
			Left:     mult,
			Operator: "-",
			Right:    &LiteralExpression{Value: runtime.NewNumericValue(5)},
		}
		
		result, err := expr.Evaluate(env)
		
		require.NoError(t, err)
		assert.Equal(t, 93.0, result.NumValue)
	})
}

func TestInvalidOperator(t *testing.T) {
	env := runtime.NewEnvironment()

	t.Run("unsupported operator", func(t *testing.T) {
		left := &LiteralExpression{Value: runtime.NewNumericValue(5)}
		right := &LiteralExpression{Value: runtime.NewNumericValue(3)}
		expr := &BinaryExpression{Left: left, Operator: "%", Right: right}
		
		_, err := expr.Evaluate(env)
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported operator")
	})
}

// Tests for helper functions and refactored functionality

func TestHelperFunctions(t *testing.T) {
	t.Run("NewLiteralExpression", func(t *testing.T) {
		value := runtime.NewNumericValue(42)
		expr := NewLiteralExpression(value)
		
		assert.NotNil(t, expr)
		assert.Equal(t, value, expr.Value)
	})

	t.Run("NewVariableExpression", func(t *testing.T) {
		expr := NewVariableExpression("TEST")
		
		assert.NotNil(t, expr)
		assert.Equal(t, "TEST", expr.Name)
	})

	t.Run("NewBinaryExpression", func(t *testing.T) {
		left := NewLiteralExpression(runtime.NewNumericValue(5))
		right := NewLiteralExpression(runtime.NewNumericValue(3))
		expr := NewBinaryExpression(left, OpAdd, right)
		
		assert.NotNil(t, expr)
		assert.Equal(t, left, expr.Left)
		assert.Equal(t, OpAdd, expr.Operator)
		assert.Equal(t, right, expr.Right)
	})

	t.Run("NewParenthesesExpression", func(t *testing.T) {
		inner := NewLiteralExpression(runtime.NewNumericValue(42))
		expr := NewParenthesesExpression(inner)
		
		assert.NotNil(t, expr)
		assert.Equal(t, inner, expr.Expression)
	})
}

func TestOperatorConstants(t *testing.T) {
	t.Run("operator constants are correct", func(t *testing.T) {
		assert.Equal(t, "+", OpAdd)
		assert.Equal(t, "-", OpSubtract)
		assert.Equal(t, "*", OpMultiply)
		assert.Equal(t, "/", OpDivide)
		assert.Equal(t, "^", OpPower)
	})

	t.Run("binary expression with operator constants", func(t *testing.T) {
		env := runtime.NewEnvironment()
		
		// Test each operator constant
		testCases := []struct {
			operator string
			left     float64
			right    float64
			expected float64
		}{
			{OpAdd, 5, 3, 8},
			{OpSubtract, 10, 4, 6},
			{OpMultiply, 6, 7, 42},
			{OpDivide, 15, 3, 5},
			{OpPower, 2, 3, 8},
		}

		for _, tc := range testCases {
			left := NewLiteralExpression(runtime.NewNumericValue(tc.left))
			right := NewLiteralExpression(runtime.NewNumericValue(tc.right))
			expr := NewBinaryExpression(left, tc.operator, right)
			
			result, err := expr.Evaluate(env)
			
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result.NumValue, "operator %s failed", tc.operator)
		}
	})
}

func TestIsValidOperator(t *testing.T) {
	t.Run("valid operators", func(t *testing.T) {
		validOps := []string{OpAdd, OpSubtract, OpMultiply, OpDivide, OpPower}
		
		for _, op := range validOps {
			assert.True(t, IsValidOperator(op), "operator %s should be valid", op)
		}
	})

	t.Run("invalid operators", func(t *testing.T) {
		invalidOps := []string{"%", "&", "|", "<<", ">>", "&&", "||", "==", "!="}
		
		for _, op := range invalidOps {
			assert.False(t, IsValidOperator(op), "operator %s should be invalid", op)
		}
	})
}

func TestGetOperatorPrecedence(t *testing.T) {
	t.Run("precedence levels", func(t *testing.T) {
		// Power has highest precedence
		assert.Equal(t, 3, GetOperatorPrecedence(OpPower))
		
		// Multiply and divide have medium precedence
		assert.Equal(t, 2, GetOperatorPrecedence(OpMultiply))
		assert.Equal(t, 2, GetOperatorPrecedence(OpDivide))
		
		// Add and subtract have lowest precedence
		assert.Equal(t, 1, GetOperatorPrecedence(OpAdd))
		assert.Equal(t, 1, GetOperatorPrecedence(OpSubtract))
		
		// Unknown operators have zero precedence
		assert.Equal(t, 0, GetOperatorPrecedence("%"))
		assert.Equal(t, 0, GetOperatorPrecedence("unknown"))
	})

	t.Run("precedence ordering", func(t *testing.T) {
		// Verify precedence ordering is correct
		assert.Greater(t, GetOperatorPrecedence(OpPower), GetOperatorPrecedence(OpMultiply))
		assert.Greater(t, GetOperatorPrecedence(OpMultiply), GetOperatorPrecedence(OpAdd))
		assert.Greater(t, GetOperatorPrecedence(OpDivide), GetOperatorPrecedence(OpSubtract))
		assert.Equal(t, GetOperatorPrecedence(OpMultiply), GetOperatorPrecedence(OpDivide))
		assert.Equal(t, GetOperatorPrecedence(OpAdd), GetOperatorPrecedence(OpSubtract))
	})
}

func TestImprovedErrorMessages(t *testing.T) {
	env := runtime.NewEnvironment()

	t.Run("error in left operand", func(t *testing.T) {
		// Create an expression that will fail on the left side
		left := NewBinaryExpression(
			NewLiteralExpression(runtime.NewStringValue("hello")),
			OpSubtract,
			NewLiteralExpression(runtime.NewStringValue("world")),
		)
		right := NewLiteralExpression(runtime.NewNumericValue(5))
		expr := NewBinaryExpression(left, OpAdd, right)
		
		_, err := expr.Evaluate(env)
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error evaluating left operand")
	})

	t.Run("error in right operand", func(t *testing.T) {
		// Create an expression that will fail on the right side
		left := NewLiteralExpression(runtime.NewNumericValue(5))
		right := NewBinaryExpression(
			NewLiteralExpression(runtime.NewStringValue("hello")),
			OpMultiply,
			NewLiteralExpression(runtime.NewNumericValue(3)),
		)
		expr := NewBinaryExpression(left, OpAdd, right)
		
		_, err := expr.Evaluate(env)
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error evaluating right operand")
	})

	t.Run("error in parenthesized expression", func(t *testing.T) {
		// Create a parenthesized expression that will fail
		inner := NewBinaryExpression(
			NewLiteralExpression(runtime.NewStringValue("hello")),
			OpDivide,
			NewLiteralExpression(runtime.NewNumericValue(0)),
		)
		expr := NewParenthesesExpression(inner)
		
		_, err := expr.Evaluate(env)
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error evaluating parenthesized expression")
	})
}