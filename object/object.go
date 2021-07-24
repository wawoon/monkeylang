package object

import "strconv"

type ObjectType string

const (
	INTEGER_OBJECT ObjectType = "INTEGER"
	BOOLEAN_OBJECT ObjectType = "BOOLEAN"
	NULL_OBJECT    ObjectType = "NULL"
	RETURN_OBJECT  ObjectType = "RETURN"
	ERROR_OBJECT   ObjectType = "ERROR"
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
