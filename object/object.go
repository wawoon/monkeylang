package object

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/wawoon/monkeylang/ast"
)

type ObjectType string

const (
	INTEGER_OBJECT  ObjectType = "INTEGER"
	BOOLEAN_OBJECT  ObjectType = "BOOLEAN"
	NULL_OBJECT     ObjectType = "NULL"
	RETURN_OBJECT   ObjectType = "RETURN"
	ERROR_OBJECT    ObjectType = "ERROR"
	FUNCTION_OBJECT ObjectType = "FUNCTION"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i Integer) Type() ObjectType {
	return INTEGER_OBJECT
}

func (i Integer) Inspect() string {
	return strconv.FormatInt(i.Value, 10)
}

type Boolean struct {
	Value bool
}

func (b Boolean) Type() ObjectType {
	return BOOLEAN_OBJECT
}

func (b Boolean) Inspect() string {
	return strconv.FormatBool(b.Value)
}

type Null struct{}

func (n Null) Type() ObjectType {
	return NULL_OBJECT
}
func (n Null) Inspect() string {
	return "null"
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_OBJECT
}
func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type Error struct {
	Message string
}

func (e Error) Type() ObjectType {
	return ERROR_OBJECT
}
func (e Error) Inspect() string {
	return "Error: " + e.Message
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJECT
}
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("}\n")
	return out.String()
}
