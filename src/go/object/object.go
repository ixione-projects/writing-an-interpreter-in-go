package object

import "fmt"

type ObjectType int

const (
	INTEGER ObjectType = iota
	BOOLEAN
	NULL
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
	return INTEGER
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

var objects = map[ObjectType]string{
	INTEGER: "INTEGER",
	BOOLEAN: "BOOLEAN",
	NULL:    "NULL",
}

func (ot ObjectType) String() string {
	return objects[ot]
}
