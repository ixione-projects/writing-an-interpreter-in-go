package evaluator

import (
	"strings"

	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/object"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) (object.Object, object.Interruption) {
			if len(args) != 1 {
				return nil, toBuiltinError("len", args)
			}

			switch args[0].Type() {
			case object.STRING:
				return object.Number(len(args[0].(object.String))), nil
			default:
				return nil, toBuiltinError("len", args)
			}
		},
	},
}

func toBuiltinError(name string, args []object.Object) *object.Error {
	types := []string{}
	for _, arg := range args {
		types = append(types, arg.Type().String())
	}
	return toError("argument(s) to `%s` not supported: (%s)", name, strings.Join(types, ", "))
}
