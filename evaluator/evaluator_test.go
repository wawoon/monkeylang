package evaluator

import (
	"testing"

	"github.com/wawoon/monkeylang/lexer"
	"github.com/wawoon/monkeylang/object"
	"github.com/wawoon/monkeylang/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{input: "5", expected: 5},
		{input: "10", expected: 10},
		{input: "-5", expected: -5},
		{input: "-10", expected: -10},
		{input: "5 + 5 + 5 + 5 - 10", expected: 10},
		{input: "2 * 2 * 2 * 2 * 2", expected: 32},
		{input: "-50 + 100 + -50", expected: 0},
		{input: "5 * 2 + 10", expected: 20},
		{input: "5 + 2 * 10", expected: 25},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{input: "true", expected: true},
		{input: "false", expected: false},
		{input: "1 < 2", expected: true},
		{input: "1 > 2", expected: false},
		{input: "1 < 1", expected: false},
		{input: "1 > 1", expected: false},
		{input: "1 == 1", expected: true},
		{input: "1 != 1", expected: false},
		{input: "1 == 2", expected: false},
		{input: "1 != 2", expected: true},
		{input: "true == true", expected: true},
		{input: "true != true", expected: false},
		{input: "true == false", expected: false},
		{input: "true != false", expected: true},
		{input: "false == false", expected: true},
		{input: "false != false", expected: false},
		{input: "false == true", expected: false},
		{input: "false != true", expected: true},
		{input: "(1 < 2) == true", expected: true},
		{input: "(1 < 2) == false", expected: false},
		{input: "(1 > 2) == true", expected: false},
		{input: "(1 > 2) == false", expected: true},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testBooleanObject(t, evaluated, test.expected)
	}
}

func TestEvalBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{input: "!true", expected: false},
		{input: "!false", expected: true},
		{input: "!5", expected: false},
		{input: "!!true", expected: true},
		{input: "!!false", expected: false},
		{input: "!!5", expected: true},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testBooleanObject(t, evaluated, test.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{input: `if (true) { 10 }`, expected: 10},
		{input: `if (false) { 10 }`, expected: nil},
		{input: `if (1) { 10 }`, expected: 10},
		{input: `if (1 < 2) { 10 }`, expected: 10},
		{input: `if (1 > 2) { 10 }`, expected: nil},
		{input: `if (1 < 2) { 10 } else { 20 }`, expected: 10},
		{input: `if (1 > 2) { 10 } else { 20 }`, expected: 20},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		integer, ok := test.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{input: "return 10;", expected: 10},
		{input: "return 10; return 20;", expected: 10},
		{input: "return 2 * 5; 9", expected: 10},
		{input: "9; return 2 * 5; return 9", expected: 10},
		{input: `
		if (10 > 1) {
			if (10 > 2) {
				return 10;
			}
			return 1;
		}
		`, expected: 10},
	}

	for _, test := range tests {
		val := testEval(test.input)
		testIntegerObject(t, val, test.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "5 + true;",
			expected: "type mismatch: INTEGER + BOOLEAN",
		},
		{
			input:    "5 + true; 5;",
			expected: "type mismatch: INTEGER + BOOLEAN",
		},
		{
			input:    "-true;",
			expected: "unknown operator: -BOOLEAN",
		},
		{
			input:    "true + true;",
			expected: "unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			input:    "5; true + false; 5;",
			expected: "unknown operator: BOOLEAN + BOOLEAN"},
		{
			input:    "if (10 > 1) { true + false; }",
			expected: "unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			input: `
		if (10 > 1) {
			if (10 > 1) {
				return true + false;
			}
			return 1;
		}
		`,
			expected: "unknown operator: BOOLEAN + BOOLEAN",
		},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testErrorObject(t, evaluated, test.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program)
}

func testIntegerObject(t *testing.T, evaluated object.Object, expected int64) {
	if evaluated.Type() != object.INTEGER_OBJECT {
		t.Errorf("Expected an integer, but got %s", evaluated.Type())
	}

	if evaluated.(*object.Integer).Value != expected {
		t.Errorf("Expected %d, but got %d", expected, evaluated.(*object.Integer).Value)
	}
}

func testBooleanObject(t *testing.T, evaluated object.Object, expected bool) {
	if evaluated.Type() != object.BOOLEAN_OBJECT {
		t.Errorf("Expected a boolean, but got %s", evaluated.Type())
	}

	if evaluated.(*object.Boolean).Value != expected {
		t.Errorf("Expected %t, but got %t", expected, evaluated.(*object.Boolean).Value)
	}
}

func testNullObject(t *testing.T, obj object.Object) {
	if obj.Type() != object.NULL_OBJECT {
		t.Errorf("Expected a null, but got %s", obj.Type())
	}
}

func testErrorObject(t *testing.T, obj object.Object, expected string) {
	if obj.Type() != object.ERROR_OBJECT {
		t.Errorf("Expected an error, but got %s", obj.Type())
	}

	if obj.(*object.Error).Message != expected {
		t.Errorf("Expected %s, but got %s", expected, obj.(*object.Error).Message)
	}
}
