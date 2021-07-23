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
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("ParseProgram: expected a ReturnStatement, got %T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %T", returnStmt.TokenLiteral())
		}
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

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("ParseProgram: expected an Identifier, got %T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Fatalf("ParseProgram: expected an Identifier foobar, got %s", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("ParseProgram: expected an Identifier foobar, got %s", ident.TokenLiteral())
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
