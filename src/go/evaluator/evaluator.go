package evaluator

import (
	"fmt"

	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/ast"
	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/object"
)

const (
	TRUE  = object.Boolean(true)
	FALSE = object.Boolean(false)
)

var NULL = &object.Null{}

func Evaluate(node ast.Node) object.Object {
	switch node.Type() {
	case ast.PROGRAM:
		return evaluateProgram(node.(*ast.Program))
	case ast.RETURN_STATEMENT:
		return evaluateReturnStatement(node.(*ast.ReturnStatement))
	case ast.EXPRESSION_STATEMENT:
		return evaluateExpressionStatement(node.(*ast.ExpressionStatement))
	case ast.BLOCK_STATEMENT:
		return evaluateBlockStatement(node.(*ast.BlockStatement))
	case ast.PREFIX_EXPRESSION:
		return evaluatePrefixExpression(node.(*ast.PrefixExpression))
	case ast.INFIX_EXPRESSION:
		return evaluateInfixExpression(node.(*ast.InfixExpression))
	case ast.IF_EXPRESSION:
		return evaluateIfExpression(node.(*ast.IfExpression))
	case ast.NUMBER_LITERAL:
		return evaluateNumberLiteral(node.(*ast.NumberLiteral))
	case ast.BOOLEAN_LITERAL:
		return evaluateBooleanLiteral(node.(*ast.BooleanLiteral))
	}
	return nil
}

func evaluateProgram(node *ast.Program) object.Object {
	var result object.Object
	for _, stmt := range node.Statements {
		result = Evaluate(stmt)
		switch result.Type() {
		case object.RETURN_VALUE:
			return result.(*object.ReturnValue).Value
		case object.ERROR:
			return result
		}
	}
	return result
}

func evaluateReturnStatement(node *ast.ReturnStatement) object.Object {
	value := Evaluate(node.ReturnValue)
	if isError(value) {
		return value
	}
	return &object.ReturnValue{Value: value}
}

func evaluateExpressionStatement(node *ast.ExpressionStatement) object.Object {
	return Evaluate(node.Expression)
}

func evaluateBlockStatement(node *ast.BlockStatement) object.Object {
	var result object.Object
	for _, stmt := range node.Statements {
		result = Evaluate(stmt)
		if result != nil && (result.Type() == object.RETURN_VALUE || result.Type() == object.ERROR) {
			return result
		}
	}
	return result
}

func evaluatePrefixExpression(node *ast.PrefixExpression) object.Object {
	right := Evaluate(node.Right)
	if isError(right) {
		return right
	}
	switch node.Operator {
	case "!":
		if isTruthy(right) {
			return FALSE
		}
		return TRUE
	case "-":
		switch right.Type() {
		case object.NUMBER:
			return object.Number(-right.(object.Number))
		}
	}

	return toError("unknown operator: %s%s", node.Operator, right.Type())
}

func evaluateInfixExpression(node *ast.InfixExpression) object.Object {
	left := Evaluate(node.Left)
	if isError(left) {
		return left
	}

	right := Evaluate(node.Right)
	if isError(right) {
		return right
	}
	switch {
	case left.Type() == object.NUMBER && right.Type() == object.NUMBER:
		switch node.Operator {
		case "+":
			return object.Number(left.(object.Number) + right.(object.Number))
		case "-":
			return object.Number(left.(object.Number) - right.(object.Number))
		case "*":
			return object.Number(left.(object.Number) * right.(object.Number))
		case "/":
			return object.Number(left.(object.Number) / right.(object.Number))
		case "<":
			return toBoolean(left.(object.Number) < right.(object.Number))
		case ">":
			return toBoolean(left.(object.Number) > right.(object.Number))
		case "==":
			return toBoolean(left.(object.Number) == right.(object.Number))
		case "!=":
			return toBoolean(left.(object.Number) != right.(object.Number))
		}
	case left.Type() != right.Type():
		return toError("type mismatch: %s %s %s", left.Type(), node.Operator, right.Type())
	case node.Operator == "==":
		return toBoolean(left == right)
	case node.Operator == "!=":
		return toBoolean(left != right)
	}

	return toError("unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
}

func evaluateIfExpression(node *ast.IfExpression) object.Object {
	condition := Evaluate(node.Condition)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Evaluate(node.Consequence)
	} else if node.Alternative != nil {
		return Evaluate(node.Alternative)
	}
	return NULL
}

func evaluateNumberLiteral(node *ast.NumberLiteral) object.Object {
	return object.Number(node.Value)
}

func evaluateBooleanLiteral(node *ast.BooleanLiteral) object.Object {
	return toBoolean(node.Value)
}

func toBoolean(value bool) object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
}

func toError(format string, args ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, args...)}
}

func isError(o object.Object) bool {
	if o != nil {
		return o.Type() == object.ERROR
	}
	return false
}

func isTruthy(o object.Object) object.Boolean {
	switch {
	case o == FALSE:
		return FALSE
	case o == NULL:
		return FALSE
	case o.Type() == object.NUMBER && o.(object.Number) == 0:
		return FALSE
	default:
		return TRUE
	}
}
