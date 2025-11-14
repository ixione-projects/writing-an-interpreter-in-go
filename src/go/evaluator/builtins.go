package evaluator

import (
	"fmt"
	"strings"

	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/ast"
	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/object"
	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/token"
)

var builtins map[string]object.Builtin

func init() {
	builtins = map[string]object.Builtin{
		"len": &object.BuiltinFunction{
			Fn: func(ctx *object.Environment, args ...object.Object) (object.Object, object.Interruption) {
				if len(args) != 1 {
					return nil, toBuiltinError("len", args)
				}

				switch args[0].Type() {
				case object.STRING:
					return object.Number(len(args[0].(object.String))), nil
				case object.ARRAY:
					return object.Number(len(args[0].(*object.Array).Elements)), nil
				default:
					return nil, toBuiltinError("len", args)
				}
			},
		},
		"puts": &object.BuiltinFunction{
			Fn: func(ctx *object.Environment, args ...object.Object) (object.Object, object.Interruption) {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}
				return NULL, nil
			},
		},
		"first": &object.BuiltinFunction{
			Fn: func(ctx *object.Environment, args ...object.Object) (object.Object, object.Interruption) {
				if len(args) != 1 {
					return nil, toBuiltinError("first", args)
				}

				switch args[0].Type() {
				case object.ARRAY:
					elements := args[0].(*object.Array).Elements
					if len(elements) == 0 {
						return NULL, nil
					}
					return elements[0], nil
				default:
					return nil, toBuiltinError("first", args)
				}
			},
		},
		"last": &object.BuiltinFunction{
			Fn: func(ctx *object.Environment, args ...object.Object) (object.Object, object.Interruption) {
				if len(args) != 1 {
					return nil, toBuiltinError("last", args)
				}

				switch args[0].Type() {
				case object.ARRAY:
					elements := args[0].(*object.Array).Elements
					if len(elements) == 0 {
						return NULL, nil
					}
					return elements[len(elements)-1], nil
				default:
					return nil, toBuiltinError("last", args)
				}
			},
		},
		"rest": &object.BuiltinFunction{
			Fn: func(ctx *object.Environment, args ...object.Object) (object.Object, object.Interruption) {
				if len(args) != 1 {
					return nil, toBuiltinError("rest", args)
				}

				switch args[0].Type() {
				case object.ARRAY:
					elements := args[0].(*object.Array).Elements
					length := len(elements)
					if length == 0 {
						return NULL, nil
					}
					result := make([]object.Object, length-1)
					copy(result, elements[1:])
					return &object.Array{Elements: result}, nil
				default:
					return nil, toBuiltinError("rest", args)
				}
			},
		},
		"push": &object.BuiltinFunction{
			Fn: func(ctx *object.Environment, args ...object.Object) (object.Object, object.Interruption) {
				if len(args) != 2 {
					return nil, toBuiltinError("push", args)
				}

				switch args[0].Type() {
				case object.ARRAY:
					elements := args[0].(*object.Array).Elements
					length := len(elements)
					result := make([]object.Object, length+1)
					copy(result, elements)
					result[length] = args[1]
					return &object.Array{Elements: result}, nil
				default:
					return nil, toBuiltinError("push", args)
				}
			},
		},
		"quote": &object.BuiltinMacro{
			Fn: func(ctx *object.Environment, args ...object.Object) (object.Object, object.Interruption) {
				if len(args) != 1 {
					return nil, toBuiltinError("quote", args)
				}

				quote, ok := args[0].(*object.Quote)
				if !ok {
					return nil, toBuiltinError("quote", args)
				}

				node := ast.Modify(quote.Node, func(node ast.Node) ast.Node {
					call, ok := node.(*ast.CallExpression)
					if !ok {
						return node
					}

					if call.Callee.TokenLiteral() != "unquote" {
						return node
					}

					ctx.Quoting = true
					result, interrupt := Evaluate(call, ctx)
					ctx.Quoting = false
					if interrupt != nil {
						switch interrupt := interrupt.(type) {
						case *object.ReturnValue:
							result = interrupt.Value
						case *object.Error:
							return &ast.Error{
								Message: interrupt.Message,
							}
						}
					}

					return toNode(result)
				})

				return &object.Quote{Node: node.(ast.Expression)}, nil
			},
		},
		"unquote": &object.BuiltinMacro{
			Fn: func(ctx *object.Environment, args ...object.Object) (object.Object, object.Interruption) {
				if !ctx.Quoting {
					return nil, &object.Error{Message: "`unquote` can only be invoked during quoting"}
				}

				if len(args) != 1 {
					return nil, toBuiltinError("unquote", args)
				}

				quote, ok := args[0].(*object.Quote)
				if !ok {
					return nil, toBuiltinError("quote", args)
				}

				return Evaluate(quote.Node, ctx)
			},
		},
	}
}

func toNode(o object.Object) ast.Expression {
	switch o := o.(type) {
	case *object.Function:
		return o.Literal
	case object.Number:
		return &ast.NumberLiteral{
			Token: token.Token{
				Type:    token.NUMBER,
				Literal: fmt.Sprintf("%g", o),
			},
			Value: float64(o),
		}
	case object.Boolean:
		var tok token.Token
		if o {
			tok = token.Token{
				Type:    token.TRUE,
				Literal: "true",
			}
		} else {
			tok = token.Token{
				Type:    token.FALSE,
				Literal: "false",
			}
		}

		return &ast.BooleanLiteral{
			Token: tok,
			Value: bool(o),
		}
	case object.String:
		return &ast.StringLiteral{
			Token: token.Token{
				Type:    token.STRING,
				Literal: "\"" + string(o) + "\"",
			},
			Value: string(o),
		}
	case *object.Array:
		node := &ast.ArrayLiteral{
			Token: token.Token{
				Type:    token.LBRACK,
				Literal: "[",
			},
			Elements: []ast.Expression{},
		}

		for i, elem := range o.Elements {
			node.Elements[i] = toNode(elem)
		}
		return node
	case *object.Hash:
		node := &ast.HashLiteral{
			Token: token.Token{
				Type:    token.LBRACK,
				Literal: "[",
			},
			Keys:  []ast.Expression{},
			Pairs: map[ast.Expression]ast.Expression{},
		}

		for _, pair := range o.Pairs {
			key := toNode(pair.Key)
			value := toNode(pair.Value)

			node.Keys = append(node.Keys, key)
			node.Pairs[key] = value
		}

		return node
	case *object.Null:
		return &ast.NullLiteral{
			Token: token.Token{
				Type:    token.NULL,
				Literal: "null",
			},
		}
	case *object.Quote:
		return o.Node
	default:
		panic(fmt.Errorf("unexpected object type: %s", o.Type()))
	}
}

func toBuiltinError(name string, args []object.Object) *object.Error {
	types := []string{}
	for _, arg := range args {
		types = append(types, arg.Type().String())
	}
	return toError("argument(s) to `%s` not supported: (%s)", name, strings.Join(types, ", "))
}
