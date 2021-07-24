package parser

import (
	"strconv"
	"testing"

	"github.com/wawoon/monkeylang/ast"
	"github.com/wawoon/monkeylang/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram: returned nil")
	}
	checkParserError(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("ParseProgram: expected 3 statements, got %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdent string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdent) {
			t.FailNow()
		}
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram: returned nil")
	}
	checkParserError(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("ParseProgram: expected 3 statements, got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		testReturnStatement(t, stmt)
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram: returned nil")
	}
	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("ParseProgram: expected 1 statements, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ParseProgram: expected a ExpressionStatement, got %T", program.Statements[0])
	}

	testIdentifier(t, stmt.Expression, "foobar")
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		exp   bool
	}{
		{`true;`, true},
		{`false;`, false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		if program == nil {
			t.Fatalf("ParseProgram: returned nil")
		}
		checkParserError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("ParseProgram: expected 1 statements, got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("ParseProgram: expected a ExpressionStatement, got %T", program.Statements[0])
		}

		testBoolean(t, stmt.Expression, tt.exp)
	}
}

func TestIntegerExpression(t *testing.T) {
	input := `5;`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram: returned nil")
	}
	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("ParseProgram: expected 1 statements, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ParseProgram: expected a ExpressionStatement, got %T", program.Statements[0])
	}

	il, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("ParseProgram: expected an IntegerLiteral, got %T", stmt.Expression)
	}

	if il.Value != 5 {
		t.Fatalf("ParseProgram: expected an IntegerLiteral 5, got %d", il.Value)
	}

	if il.TokenLiteral() != "5" {
		t.Fatalf("ParseProgram: expected an Identifier foobar, got %s", il.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{`!5;`, "!", 5},
		{`-15;`, "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		if program == nil {
			t.Fatalf("ParseProgram: returned nil")
		}

		checkParserError(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("ParseProgram: expected 1 statements, got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("ParseProgram: expected a ExpressionStatement, got %T", program.Statements[0])
		}

		testPrefixExpression(t, stmt.Expression, tt.operator, tt.integerValue)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTest := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{`5 + 5;`, 5, "+", 5},
		{`5 - 5;`, 5, "-", 5},
		{`5 * 5;`, 5, "*", 5},
		{`5 / 5;`, 5, "/", 5},
		{`5 > 5;`, 5, ">", 5},
		{`5 < 5;`, 5, "<", 5},
		{`5 == 5;`, 5, "==", 5},
		{`5 != 5;`, 5, "!=", 5},
		{`true == true;`, true, "==", true},
		{`true != false;`, true, "!=", false},
		{`false == false;`, false, "==", false},
	}

	for _, tt := range infixTest {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		if program == nil {
			t.Fatalf("ParseProgram: returned nil")
		}

		checkParserError(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("ParseProgram: expected 1 statements, got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("ParseProgram: expected a ExpressionStatement, got %T", program.Statements[0])
		}

		testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    `-a * b`,
			expected: "((-a) * b)",
		},
		{
			input:    `!-a`,
			expected: "(!(-a))",
		},
		{
			input:    `a + b + c`,
			expected: "((a + b) + c)",
		},
		{
			input:    `a + b - c`,
			expected: "((a + b) - c)",
		},
		{
			input:    `a * b * c`,
			expected: "((a * b) * c)",
		},
		{
			input:    `a * b / c`,
			expected: "((a * b) / c)",
		},
		{
			input:    `a + b / c`,
			expected: "(a + (b / c))",
		},
		{
			input:    `a + b * c + d / e - f`,
			expected: `(((a + (b * c)) + (d / e)) - f)`,
		},
		{
			input:    `3 + 4; -5 * 5`,
			expected: `(3 + 4)((-5) * 5)`,
		},
		{
			input:    `5 > 4 == 3 < 4`,
			expected: `((5 > 4) == (3 < 4))`,
		},
		{
			input:    `5 < 4 != 3 > 4`,
			expected: `((5 < 4) != (3 > 4))`,
		},
		{
			input:    `3 + 4 * 5 == 3 * 1 + 4 * 5`,
			expected: `((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))`,
		},
		{
			input:    `true`,
			expected: "true",
		},
		{
			input:    `false`,
			expected: "false",
		},
		{
			input:    `3 > 5 == false`,
			expected: `((3 > 5) == false)`,
		},
		{
			input:    `3 < 5 == true`,
			expected: `((3 < 5) == true)`,
		},
		{
			input:    `1 + (2 + 3) + 4`,
			expected: "((1 + (2 + 3)) + 4)",
		},
		{
			input:    `(5 + 5) * 2`,
			expected: "((5 + 5) * 2)",
		},
		{
			input:    `2 / (5 + 5)`,
			expected: "(2 / (5 + 5))",
		},
		{
			input:    `-(5 + 5)`,
			expected: "(-(5 + 5))",
		},
		{
			input:    `!(true == true)`,
			expected: "(!(true == true))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		if program == nil {
			t.Fatalf("ParseProgram: returned nil")
		}
		checkParserError(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("ParseProgram: expected %q, got %q", tt.expected, actual)
		}
	}
}

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
