package parser

import (
	"strconv"
	"testing"

	"github.com/wawoon/monkeylang/ast"
)

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Fatalf("testLetStatement: statement token should be let, got %q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Fatalf("testLetStatement: statement should be of type LetStatement, got %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Fatalf("testLetStatement: let statement name should be %q, got %q", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Fatalf("testLetStatement: let statement name token should be %q, got %q", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func testReturnStatement(t *testing.T, stmt ast.Statement) bool {
	returnStmt, ok := stmt.(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("ParseProgram: expected a ReturnStatement, got %T", stmt)
		return false
	}

	if returnStmt.TokenLiteral() != "return" {
		t.Fatalf("returnStmt.TokenLiteral not 'return', got %T", returnStmt.TokenLiteral())
		return false
	}
	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	if il == nil {
		t.Fatalf("testIntegerLiteral: nil integer literal")
		return false
	}

	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("testIntegerLiteral: integer literal should be of type IntegerLiteral, got %T", il)
		return false
	}

	if integer.Value != value {
		t.Fatalf("testIntegerLiteral: integer literal value should be %d, got %d", value, integer.Value)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Fatalf("testIdentifier: expression should be of type Identifier, got %T", exp)
		return false
	}

	if ident.Value != value {
		t.Fatalf("testIdentifier: identifier value should be %q, got %q", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Fatalf("testIdentifier: identifier token should be %q, got %q", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testBoolean(t *testing.T, exp ast.Expression, value bool) bool {
	b, ok := exp.(*ast.Boolean)
	if !ok {
		t.Fatalf("testBoolean: expression should be of type Boolean, got %T", exp)
		return false
	}

	if b.Value != value {
		t.Fatalf("testBoolean: value should be %s, got %s", strconv.FormatBool(value), strconv.FormatBool(b.Value))
		return false
	}

	if b.TokenLiteral() != strconv.FormatBool(value) {
		t.Fatalf("testBoolean: token should be %s, got %s", strconv.FormatBool(value), b.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolean(t, exp, v)
	}

	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testPrefixExpression(t *testing.T, exp ast.Expression, operator string, value int64) bool {
	prefix, ok := exp.(*ast.PrefixExpression)
	if !ok {
		t.Fatalf("testPrefixExpression: expression should be of type PrefixExpression, got %T", exp)
		return false
	}

	if prefix.Operator != operator {
		t.Fatalf("testPrefixExpression: prefix expression operator should be %q, got %q", operator, prefix.Operator)
		return false
	}

	return testIntegerLiteral(t, prefix.Right, value)
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("testInfixExpression: expression should be of type InfixExpression, got %T", exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Fatalf("testInfixExpression: operator should be %q, got %q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, exp ast.Expression, value string) bool {
	str, ok := exp.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("testStringLiteral: expression should be of type StringLiteral, got %T", exp)
		return false
	}

	if str.Value != value {
		t.Fatalf("testStringLiteral: string literal value should be %q, got %q", value, str.Value)
		return false
	}

	return true
}

func checkParserError(t *testing.T, p *Parser) {
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errs))
	for _, msg := range errs {
		t.Errorf("parser error: %s", msg)
	}
	t.FailNow()
}
