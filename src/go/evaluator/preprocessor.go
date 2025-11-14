package evaluator

import (
	"fmt"

	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/ast"
	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/object"
)

func DefineMacros(program *ast.Program, env *object.Environment) *ast.Program {
	defs := []int{}

	for i, stmt := range program.Statements {
		if macro, ok := stmt.(*ast.MacroStatement); ok {
			defs = append(defs, i)
			object := &object.Macro{
				Declaration: macro,
				Environment: env,
			}
			env.Set(macro.Name.Value, object)
		}
	}

	for i := len(defs) - 1; i >= 0; i -= 1 {
		program.Statements = append(program.Statements[:defs[i]], program.Statements[defs[i]+1:]...)
	}

	return program
}

func ExpandMacros(program ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(program, func(node ast.Node) ast.Node {
		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		ident, ok := call.Callee.(*ast.Identifier)
		if !ok {
			return node
		}

		value, ok := env.Get(ident.Value)
		if !ok {
			return node
		}

		macro, ok := value.(*object.Macro)
		if !ok {
			return node
		}

		args := []*object.Quote{}
		for _, arg := range call.Arguments {
			args = append(args, &object.Quote{Node: arg})
		}

		environemnt := object.NewEnvironment(env)
		for i, param := range macro.Declaration.Parameters {
			environemnt.Set(param.Value, args[i])
		}

		result, interrput := Evaluate(macro.Declaration.Body, environemnt)
		if interrput != nil {
			panic("unexpected interruption during macro expansion")
		}

		quote, ok := result.(*object.Quote)
		if !ok {
			panic(fmt.Errorf("unsupported type returned from macro expansion: %s", result.Type()))
		}

		return quote.Node
	})
}
