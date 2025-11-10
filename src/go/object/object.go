package object

import "fmt"

type ObjectType int

const (
	NUMBER ObjectType = iota
	BOOLEAN
	NULL
	RETURN_VALUE
	ERROR
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type (
	Number  float64
	Boolean bool
)

func (n Number) Type() ObjectType {
	return NUMBER
}

func (b Boolean) Type() ObjectType {
	return BOOLEAN
}

func (n Number) Inspect() string {
	return fmt.Sprintf("%g", float64(n))
}

func (b Boolean) Inspect() string {
	return fmt.Sprintf("%t", bool(b))
}

type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL
}

func (n *Null) Inspect() string {
	return "null"
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

var objects = map[ObjectType]string{
	NUMBER:       "INTEGER",
	BOOLEAN:      "BOOLEAN",
	NULL:         "NULL",
	RETURN_VALUE: "RETURN_VALUE",
	ERROR:        "ERROR",
}

func (ot ObjectType) String() string {
	return objects[ot]
}
