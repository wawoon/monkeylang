package parser

import (
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
		{
			input:    `a + add(b * c) + d`,
			expected: "((a + add((b * c))) + d)",
		},
		{
			input:    `add(a, b, 1, 2 * 3, 4 + 5, add(6 + 7 + 8))`,
			expected: "add(a, b, 1, (2 * 3), (4 + 5), add(((6 + 7) + 8)))",
		},
		{
			input:    `add(a + b + c * d / f + g)`,
			expected: "add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			input:    `a * [1,2,3,4][b*c] * d`,
			expected: "((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			input:    `add(a * b[2], b[1], 2 * [1, 2][1])`,
			expected: `add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))`,
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

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("ParseProgram: expected 1 statements, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ParseProgram: expected a ExpressionStatement, got %T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("ParseProgram: expected an IfExpression, got %T", stmt.Expression)
	}

	testInfixExpression(t, exp.Condition, "x", "<", "y")

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("ParseProgram: expected 1 consequence statements, got %d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ParseProgram: expected a ExpressionStatement, got %T", exp.Consequence.Statements[0])
	}

	testIdentifier(t, consequence.Expression, "x")

	if exp.Alternative != nil {
		t.Fatalf("ParseProgram: expected no alternative, got %T", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("ParseProgram: expected 1 statements, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ParseProgram: expected a ExpressionStatement, got %T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("ParseProgram: expected an IfExpression, got %T", stmt.Expression)
	}

	testInfixExpression(t, exp.Condition, "x", "<", "y")

	if len(exp.Consequence.Statements) != 1 {
		t.Fatalf("ParseProgram: expected 1 consequence statements, got %d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ParseProgram: expected an ExpressionStatement, got %T", exp.Consequence.Statements[0])
	}

	testIdentifier(t, consequence.Expression, "x")

	if len(exp.Alternative.Statements) != 1 {
		t.Fatalf("ParseProgram: expected 1 alternative statements, got %d", len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ParseProgram: expected an ExpressionStatement, got %T", exp.Alternative.Statements[0])
	}
	testIdentifier(t, alternative.Expression, "y")
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`
	l := lexer.New(input)

	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("ParseProgram: expected 1 statements, got %d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ParseProgram: expected a ExpressionStatement, got %T", program.Statements[0])
	}
	fl, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("ParseProgram: expected a FunctionLiteral, got %T", stmt.Expression)
	}
	if len(fl.Parameters) != 2 {
		t.Fatalf("ParseProgram: expected 2 parameters, got %d", len(fl.Parameters))
	}
	if !testIdentifier(t, fl.Parameters[0], "x") {
		return
	}
	if !testIdentifier(t, fl.Parameters[1], "y") {
		return
	}
	if len(fl.Body.Statements) != 1 {
		t.Fatalf("ParseProgram: expected 1 statements, got %d", len(fl.Body.Statements))
	}
	testInfixExpression(t, fl.Body.Statements[0].(*ast.ExpressionStatement).Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input:    "fn() {}",
			expected: []string{},
		},
		{
			input:    "fn(x) {}",
			expected: []string{"x"},
		},
		{
			input:    "fn(x, y, z) {}",
			expected: []string{"x", "y", "z"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserError(t, p)
		fl, ok := program.Statements[0].(*ast.ExpressionStatement).Expression.(*ast.FunctionLiteral)
		if !ok {
			t.Fatalf("ParseProgram: expected a FunctionLiteral, got %T", program.Statements[0].(*ast.ExpressionStatement).Expression)
		}
		if len(fl.Parameters) != len(tt.expected) {
			t.Fatalf("ParseProgram: expected %d parameters, got %d", len(tt.expected), len(fl.Parameters))
		}
		for i, p := range fl.Parameters {
			testIdentifier(t, p, tt.expected[i])
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("ParseProgram: expected 1 statements, got %d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ParseProgram: expected a ExpressionStatement, got %T", program.Statements[0])
	}
	ce, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("ParseProgram: expected a CallExpression, got %T", stmt.Expression)
	}
	if len(ce.Arguments) != 3 {
		t.Fatalf("ParseProgram: expected 3 arguments, got %d", len(ce.Arguments))
	}

	testIntegerLiteral(t, ce.Arguments[0], 1)
	testInfixExpression(t, ce.Arguments[1], 2, "*", 3)
	testInfixExpression(t, ce.Arguments[2], 4, "+", 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("ParseProgram: expected 1 statements, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ParseProgram: expected a ExpressionStatement, got %T", program.Statements[0])
	}
	testStringLiteral(t, stmt.Expression, "hello world")
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1 , 2 * 2, 3 + 3]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("ParseProgram: expected 1 statements, got %d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ParseProgram: expected a ExpressionStatement, got %T", program.Statements[0])
	}
	ar, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("ParseProgram: expected a ArrayLiteral, got %T", stmt.Expression)
	}
	if len(ar.Elements) != 3 {
		t.Fatalf("ParseProgram: expected 3 elements, got %d", len(ar.Elements))
	}
	testIntegerLiteral(t, ar.Elements[0], 1)
	testInfixExpression(t, ar.Elements[1], 2, "*", 2)
	testInfixExpression(t, ar.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserError(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("ParseProgram: expected 1 statements, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ParseProgram: expected a ExpressionStatement, got %T", program.Statements[0])
	}

	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("ParseProgram: expected a IndexExpression, got %T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserError(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("ParseProgram: expected 1 statements, got %d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ParseProgram: expected a ExpressionStatement, got %T", program.Statements[0])
	}
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("ParseProgram: expected a HashLiteral, got %T", stmt.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Fatalf("ParseProgram: expected 3 pairs, got %d", len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Fatalf("ParseProgram: expected a StringLiteral, got %T", key)
		}

		expectedValue := expected[literal.Value]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyHashLiteralsStringKeys(t *testing.T) {
	input := `{}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("ParseProgram: expected 1 statements, got %d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ParseProgram: expected a ExpressionStatement, got %T", program.Statements[0])
	}
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("ParseProgram: expected a HashLiteral, got %T", stmt.Expression)
	}
	if len(hash.Pairs) != 0 {
		t.Fatalf("ParseProgram: expected 0 pairs, got %d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserError(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("ParseProgram: expected 1 statements, got %d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ParseProgram: expected a ExpressionStatement, got %T", program.Statements[0])
	}
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("ParseProgram: expected a HashLiteral, got %T", stmt.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Fatalf("ParseProgram: expected 3 pairs, got %d", len(hash.Pairs))
	}

	tests := map[string]func(e ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Fatalf("ParseProgram: expected a StringLiteral, got %T", key)
		}

		test := tests[literal.Value]
		test(value)
	}
}
