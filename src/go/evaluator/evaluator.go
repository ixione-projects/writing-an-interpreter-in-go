package evaluator

import (
	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/ast"
	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/object"
)

const (
	ZERO = object.Number(0)

	TRUE  = object.Boolean(true)
	FALSE = object.Boolean(false)
)

var NULL = &object.Null{}

func Evaluate(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evaluateProgram(node)
	case *ast.ExpressionStatement:
		return evaluateExpressionStatement(node)
	case *ast.NumberLiteral:
		return evaluateNumberLiteral(node)
	case *ast.BooleanLiteral:
		return evaluateBooleanLiteral(node)
	}
	return nil
}

func evaluateProgram(node *ast.Program) object.Object {
	var result object.Object
	for _, stmt := range node.Statements {
		result = Evaluate(stmt)
	}
	return result
}

func evaluateExpressionStatement(node *ast.ExpressionStatement) object.Object {
	return Evaluate(node.Expression)
}

func evaluateNumberLiteral(node *ast.NumberLiteral) object.Object {
	if node.Value == 0 {
		return ZERO
	}
	return object.Number(node.Value)
}

func evaluateBooleanLiteral(node *ast.BooleanLiteral) object.Object {
	if node.Value {
		return TRUE
	}
	return FALSE
}
