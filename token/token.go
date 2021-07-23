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
	ILLEGAL    TokenType = "ILLEGAL"
	EOF        TokenType = "EOF"
	IDENT      TokenType = "IDENT"
	INT        TokenType = "INT"
	ASSIGN     TokenType = "ASSIGN"
	PLUS       TokenType = "PLUS"
	MINUS      TokenType = "MINUS"
	BANG       TokenType = "BANG"
	ASTERISK   TokenType = "ASTERISK"
	SLASH      TokenType = "SLASH"
	MODULO     TokenType = "MODULO"
	EQUALS     TokenType = "EQUALS"
	NOT_EQUALS TokenType = "NOT_EQUALS"
	GT         TokenType = "GT"
	LT         TokenType = "LT"
	GTE        TokenType = "GTE"
	LTE        TokenType = "LTE"
	EQ         TokenType = "=="
	NEQ        TokenType = "!="
	AND        TokenType = "AND"
	OR         TokenType = "OR"
	COMMA      TokenType = "COMMA"
	SEMICOLON  TokenType = "SEMICOLON"
	LPAREN     TokenType = "LPAREN"
	RPAREN     TokenType = "RPAREN"
	LBRACE     TokenType = "LBRACE"
	RBRACE     TokenType = "RBRACE"
	FUNCTION   TokenType = "FUNCTION"
	LET        TokenType = "LET"
	TRUE       TokenType = "TRUE"
	FALSE      TokenType = "FALSE"
	IF         TokenType = "IF"
	ELSE       TokenType = "ELSE"
	WHILE      TokenType = "WHILE"
	RETURN     TokenType = "RETURN"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"while":  WHILE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if t, ok := keywords[ident]; ok {
		return t
	}
	return IDENT
}
