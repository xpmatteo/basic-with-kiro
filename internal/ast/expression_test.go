package ast

import (
	"basic-interpreter/internal/runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLiteralExpression_Evaluate(t *testing.T) {
	env := runtime.NewEnvironment()

	testCases := []struct {
		name          string
		value         runtime.Value
		expectedType  runtime.ValueType
		expectedNum   float64
		expectedStr   string
	}{
		{
			name:         "numeric literal",
			value:        runtime.NewNumericValue(42.5),
			expectedType: runtime.NumericValue,
			expectedNum:  42.5,
		},
		{
			name:         "string literal",
			value:        runtime.NewStringValue("hello"),
			expectedType: runtime.StringValue,
			expectedStr:  "hello",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expr := &LiteralExpression{Value: tc.value}
			result, err := expr.Evaluate(env)
			
			require.NoError(t, err)
			assert.Equal(t, tc.expectedType, result.Type)
			
			if tc.expectedType == runtime.NumericValue {
				assert.Equal(t, tc.expectedNum, result.NumValue)
			} else {
				assert.Equal(t, tc.expectedStr, result.StrValue)
			}
		})
	}
}

func TestVariableExpression_Evaluate(t *testing.T) {
	testCases := []struct {
		name         string
		setupVar     string
		setupValue   runtime.Value
		queryVar     string
		expectedType runtime.ValueType
		expectedNum  float64
		expectedStr  string
	}{
		{
			name:         "existing numeric variable",
			setupVar:     "X",
			setupValue:   runtime.NewNumericValue(10),
			queryVar:     "X",
			expectedType: runtime.NumericValue,
			expectedNum:  10.0,
		},
		{
			name:         "existing string variable",
			setupVar:     "NAME$",
			setupValue:   runtime.NewStringValue("test"),
			queryVar:     "NAME$",
			expectedType: runtime.StringValue,
			expectedStr:  "test",
		},
		{
			name:         "uninitialized numeric variable defaults to zero",
			queryVar:     "Y",
			expectedType: runtime.NumericValue,
			expectedNum:  0.0,
		},
		{
			name:         "uninitialized string variable defaults to empty string",
			queryVar:     "EMPTY$",
			expectedType: runtime.StringValue,
			expectedStr:  "",
		},
		{
			name:         "case insensitive variable names",
			setupVar:     "test",
			setupValue:   runtime.NewNumericValue(5),
			queryVar:     "TEST",
			expectedType: runtime.NumericValue,
			expectedNum:  5.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			env := runtime.NewEnvironment()
			
			// Setup variable if specified
			if tc.setupVar != "" {
				env.SetVariable(tc.setupVar, tc.setupValue)
			}
			
			expr := &VariableExpression{Name: tc.queryVar}
			result, err := expr.Evaluate(env)
			
			require.NoError(t, err)
			assert.Equal(t, tc.expectedType, result.Type)
			
			if tc.expectedType == runtime.NumericValue {
				assert.Equal(t, tc.expectedNum, result.NumValue)
			} else {
				assert.Equal(t, tc.expectedStr, result.StrValue)
			}
		})
	}
}

func TestBinaryExpression_Arithmetic(t *testing.T) {
	env := runtime.NewEnvironment()

	testCases := []struct {
		name         string
		left         runtime.Value
		operator     string
		right        runtime.Value
		expectedType runtime.ValueType
		expectedNum  float64
		expectedStr  string
		shouldError  bool
		errorContains string
	}{
		// Addition tests
		{
			name:         "numeric addition",
			left:         runtime.NewNumericValue(5),
			operator:     "+",
			right:        runtime.NewNumericValue(3),
			expectedType: runtime.NumericValue,
			expectedNum:  8.0,
		},
		{
			name:         "string concatenation",
			left:         runtime.NewStringValue("hello"),
			operator:     "+",
			right:        runtime.NewStringValue(" world"),
			expectedType: runtime.StringValue,
			expectedStr:  "hello world",
		},
		{
			name:         "mixed type addition - number and convertible string",
			left:         runtime.NewNumericValue(5),
			operator:     "+",
			right:        runtime.NewStringValue("3"),
			expectedType: runtime.NumericValue,
			expectedNum:  8.0,
		},
		// Subtraction tests
		{
			name:         "numeric subtraction",
			left:         runtime.NewNumericValue(10),
			operator:     "-",
			right:        runtime.NewNumericValue(3),
			expectedType: runtime.NumericValue,
			expectedNum:  7.0,
		},
		{
			name:          "string subtraction should fail",
			left:          runtime.NewStringValue("hello"),
			operator:      "-",
			right:         runtime.NewStringValue("world"),
			shouldError:   true,
			errorContains: "cannot subtract strings",
		},
		// Multiplication tests
		{
			name:         "numeric multiplication",
			left:         runtime.NewNumericValue(4),
			operator:     "*",
			right:        runtime.NewNumericValue(3),
			expectedType: runtime.NumericValue,
			expectedNum:  12.0,
		},
		{
			name:          "string multiplication should fail",
			left:          runtime.NewStringValue("hello"),
			operator:      "*",
			right:         runtime.NewNumericValue(3),
			shouldError:   true,
			errorContains: "cannot multiply strings",
		},
		// Division tests
		{
			name:         "numeric division",
			left:         runtime.NewNumericValue(15),
			operator:     "/",
			right:        runtime.NewNumericValue(3),
			expectedType: runtime.NumericValue,
			expectedNum:  5.0,
		},
		{
			name:          "division by zero",
			left:          runtime.NewNumericValue(10),
			operator:      "/",
			right:         runtime.NewNumericValue(0),
			shouldError:   true,
			errorContains: "division by zero",
		},
		// Power tests
		{
			name:         "numeric exponentiation",
			left:         runtime.NewNumericValue(2),
			operator:     "^",
			right:        runtime.NewNumericValue(3),
			expectedType: runtime.NumericValue,
			expectedNum:  8.0,
		},
		{
			name:          "string power should fail",
			left:          runtime.NewStringValue("hello"),
			operator:      "^",
			right:         runtime.NewNumericValue(2),
			shouldError:   true,
			errorContains: "cannot raise strings to power",
		},
		// Invalid operator test
		{
			name:          "unsupported operator",
			left:          runtime.NewNumericValue(5),
			operator:      "%",
			right:         runtime.NewNumericValue(3),
			shouldError:   true,
			errorContains: "unsupported operator",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			left := &LiteralExpression{Value: tc.left}
			right := &LiteralExpression{Value: tc.right}
			expr := &BinaryExpression{Left: left, Operator: tc.operator, Right: right}
			
			result, err := expr.Evaluate(env)
			
			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedType, result.Type)
				
				if tc.expectedType == runtime.NumericValue {
					assert.Equal(t, tc.expectedNum, result.NumValue)
				} else {
					assert.Equal(t, tc.expectedStr, result.StrValue)
				}
			}
		})
	}
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



// Tests for helper functions and refactored functionality

func TestHelperFunctions(t *testing.T) {
	testCases := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "NewLiteralExpression",
			test: func(t *testing.T) {
				value := runtime.NewNumericValue(42)
				expr := NewLiteralExpression(value)
				
				assert.NotNil(t, expr)
				assert.Equal(t, value, expr.Value)
			},
		},
		{
			name: "NewVariableExpression",
			test: func(t *testing.T) {
				expr := NewVariableExpression("TEST")
				
				assert.NotNil(t, expr)
				assert.Equal(t, "TEST", expr.Name)
			},
		},
		{
			name: "NewBinaryExpression",
			test: func(t *testing.T) {
				left := NewLiteralExpression(runtime.NewNumericValue(5))
				right := NewLiteralExpression(runtime.NewNumericValue(3))
				expr := NewBinaryExpression(left, OpAdd, right)
				
				assert.NotNil(t, expr)
				assert.Equal(t, left, expr.Left)
				assert.Equal(t, OpAdd, expr.Operator)
				assert.Equal(t, right, expr.Right)
			},
		},
		{
			name: "NewParenthesesExpression",
			test: func(t *testing.T) {
				inner := NewLiteralExpression(runtime.NewNumericValue(42))
				expr := NewParenthesesExpression(inner)
				
				assert.NotNil(t, expr)
				assert.Equal(t, inner, expr.Expression)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.test)
	}
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