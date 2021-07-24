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
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testIntegerObject(t, evaluated, test.expected)
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
