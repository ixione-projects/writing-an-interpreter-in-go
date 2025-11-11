package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/ast"
)

type ObjectType int

const (
	FUNCTION ObjectType = iota
	BUILTIN
	NUMBER
	BOOLEAN
	STRING
	ARRAY
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
	return "<fn (" + strings.Join(params, ", ") + ")>"
}

type BuiltinFunction func(args ...Object) (Object, Interruption)

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
	return BUILTIN
}

func (b *Builtin) Inspect() string {
	return "<fn builtin>"
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

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType {
	return ARRAY
}

func (a *Array) Inspect() string {
	var out bytes.Buffer

	elems := []string{}
	for _, elem := range a.Elements {
		elems = append(elems, elem.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elems, ", "))
	out.WriteString("]")

	return out.String()
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
	FUNCTION: "FUNCTION",
	BUILTIN:  "BUILTIN",
	NUMBER:   "INTEGER",
	BOOLEAN:  "BOOLEAN",
	STRING:   "STRING",
	ARRAY:    "ARRAY",
	NULL:     "NULL",
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
