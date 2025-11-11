package object

import (
	"fmt"
	"strings"

	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/ast"
)

type ObjectType int

const (
	FUNCTION ObjectType = iota
	NUMBER
	BOOLEAN
	STRING
	NULL
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Function struct {
	Declaration *ast.FunctionLiteral
	Closure     *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION
}

func (f *Function) Inspect() string {
	params := []string{}
	for _, param := range f.Declaration.Parameters {
		params = append(params, param.Value)
	}
	return "<fn (" + strings.Join(params, ",") + ")>"
}

type (
	Number  float64
	Boolean bool
	String  string
)

func (n Number) Type() ObjectType {
	return NUMBER
}

func (b Boolean) Type() ObjectType {
	return BOOLEAN
}

func (s String) Type() ObjectType {
	return STRING
}

func (n Number) Inspect() string {
	return fmt.Sprintf("%g", float64(n))
}

func (b Boolean) Inspect() string {
	return fmt.Sprintf("%t", bool(b))
}

func (s String) Inspect() string {
	return "\"" + string(s) + "\""
}

type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL
}

func (n *Null) Inspect() string {
	return "null"
}

type InterruptionType int

const (
	RETURN_VALUE InterruptionType = iota
	ERROR
)

type Interruption interface {
	Type() InterruptionType
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() InterruptionType {
	return RETURN_VALUE
}

type Error struct {
	Message string
}

func (e *Error) Type() InterruptionType {
	return ERROR
}

var objects = map[ObjectType]string{
	NUMBER:  "INTEGER",
	BOOLEAN: "BOOLEAN",
	STRING:  "STRING",
	NULL:    "NULL",
}

func (ot ObjectType) String() string {
	return objects[ot]
}

var interruptions = map[InterruptionType]string{
	RETURN_VALUE: "RETURN_VALUE",
	ERROR:        "ERROR",
}

func (it InterruptionType) String() string {
	return interruptions[it]
}
