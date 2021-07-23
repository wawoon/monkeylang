package token

import "fmt"

type TokenType string

func (tt TokenType) String() string {
	return string(tt)
}

type Token struct {
	Type    TokenType
	Literal string
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%s)", t.Type, t.Literal)
}

const (
	ILLEGAL        TokenType = "ILLEGAL"
	EOF            TokenType = "EOF"
	IDENT          TokenType = "IDENT"
	INT            TokenType = "INT"
	ASSIGN         TokenType = "ASSIGN"
	PLUS           TokenType = "PLUS"
	MINUS          TokenType = "MINUS"
	MULTIPLY       TokenType = "MULTIPLY"
	DIVIDE         TokenType = "DIVIDE"
	MODULO         TokenType = "MODULO"
	EQUALS         TokenType = "EQUALS"
	NOT_EQUALS     TokenType = "NOT_EQUALS"
	GREATER        TokenType = "GREATER"
	LESS           TokenType = "LESS"
	GREATER_EQUALS TokenType = "GREATER_EQUALS"
	LESS_EQUALS    TokenType = "LESS_EQUALS"
	AND            TokenType = "AND"
	OR             TokenType = "OR"
	COMMA          TokenType = "COMMA"
	SEMICOLON      TokenType = "SEMICOLON"
	LPAREN         TokenType = "LPAREN"
	RPAREN         TokenType = "RPAREN"
	LBRACE         TokenType = "LBRACE"
	RBRACE         TokenType = "RBRACE"
	FUNCTION       TokenType = "FUNCTION"
	LET            TokenType = "LET"
)

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func LookupIdent(ident string) TokenType {
	if t, ok := keywords[ident]; ok {
		return t
	}
	return IDENT
}
