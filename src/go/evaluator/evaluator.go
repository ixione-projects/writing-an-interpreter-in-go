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

func Evaluate(node ast.Node, env *object.Environment) (object.Object, object.Interruption) {
	switch node.Type() {
	case ast.PROGRAM:
		return evaluateProgram(node.(*ast.Program), env)
	case ast.LET_DECLARATION:
		return evaluateLetDeclaration(node.(*ast.LetDeclaration), env)
	case ast.RETURN_STATEMENT:
		return evaluateReturnStatement(node.(*ast.ReturnStatement), env)
	case ast.EXPRESSION_STATEMENT:
		return evaluateExpressionStatement(node.(*ast.ExpressionStatement), env)
	case ast.BLOCK_STATEMENT:
		return evaluateBlockStatement(node.(*ast.BlockStatement), env)
	case ast.PREFIX_EXPRESSION:
		return evaluatePrefixExpression(node.(*ast.PrefixExpression), env)
	case ast.INFIX_EXPRESSION:
		return evaluateInfixExpression(node.(*ast.InfixExpression), env)
	case ast.IF_EXPRESSION:
		return evaluateIfExpression(node.(*ast.IfExpression), env)
	case ast.FUNCTION_LITERAL:
		return evaluateFunctionLiteral(node.(*ast.FunctionLiteral), env)
	case ast.CALL_EXPRESSION:
		return evaluateCallExpression(node.(*ast.CallExpression), env)
	case ast.IDENTIFIER:
		return evaluateIdentifier(node.(*ast.Identifier), env)
	case ast.NUMBER_LITERAL:
		return evaluateNumberLiteral(node.(*ast.NumberLiteral))
	case ast.BOOLEAN_LITERAL:
		return evaluateBooleanLiteral(node.(*ast.BooleanLiteral))
	}
	return nil, nil
}

func evaluateProgram(node *ast.Program, env *object.Environment) (object.Object, object.Interruption) {
	var result object.Object
	var interrupt object.Interruption
	for _, stmt := range node.Statements {
		result, interrupt = Evaluate(stmt, env)
		if interrupt != nil {
			if interrupt.Type() == object.RETURN_VALUE {
				return interrupt.(*object.ReturnValue).Value, nil
			}
			return nil, interrupt
		}
	}
	return result, nil
}

func evaluateLetDeclaration(node *ast.LetDeclaration, env *object.Environment) (object.Object, object.Interruption) {
	value, interrupt := Evaluate(node.Value, env)
	if interrupt != nil {
		return nil, interrupt
	}
	env.Set(node.Name.Value, value)
	return NULL, nil
}

func evaluateReturnStatement(node *ast.ReturnStatement, env *object.Environment) (object.Object, object.Interruption) {
	value, interrupt := Evaluate(node.ReturnValue, env)
	if interrupt != nil {
		return nil, interrupt
	}
	return nil, &object.ReturnValue{Value: value}
}

func evaluateExpressionStatement(node *ast.ExpressionStatement, env *object.Environment) (object.Object, object.Interruption) {
	return Evaluate(node.Expression, env)
}

func evaluateBlockStatement(node *ast.BlockStatement, env *object.Environment) (object.Object, object.Interruption) {
	var result object.Object
	var interrupt object.Interruption
	for _, stmt := range node.Statements {
		result, interrupt = Evaluate(stmt, env)
		if interrupt != nil {
			return nil, interrupt
		}
	}
	return result, nil
}

func evaluatePrefixExpression(node *ast.PrefixExpression, env *object.Environment) (object.Object, object.Interruption) {
	right, interrupt := Evaluate(node.Right, env)
	if interrupt != nil {
		return nil, interrupt
	}
	switch node.Operator {
	case "!":
		if isTruthy(right) {
			return FALSE, nil
		}
		return TRUE, nil
	case "-":
		switch right.Type() {
		case object.NUMBER:
			return object.Number(-right.(object.Number)), nil
		}
	}

	return nil, toError("unknown operator: %s%s", node.Operator, right.Type())
}

func evaluateInfixExpression(node *ast.InfixExpression, env *object.Environment) (object.Object, object.Interruption) {
	left, interrupt := Evaluate(node.Left, env)
	if interrupt != nil {
		return nil, interrupt
	}

	right, interrupt := Evaluate(node.Right, env)
	if interrupt != nil {
		return nil, interrupt
	}
	switch {
	case left.Type() == object.NUMBER && right.Type() == object.NUMBER:
		switch node.Operator {
		case "+":
			return object.Number(left.(object.Number) + right.(object.Number)), nil
		case "-":
			return object.Number(left.(object.Number) - right.(object.Number)), nil
		case "*":
			return object.Number(left.(object.Number) * right.(object.Number)), nil
		case "/":
			return object.Number(left.(object.Number) / right.(object.Number)), nil
		case "<":
			return toBoolean(left.(object.Number) < right.(object.Number)), nil
		case ">":
			return toBoolean(left.(object.Number) > right.(object.Number)), nil
		case "==":
			return toBoolean(left.(object.Number) == right.(object.Number)), nil
		case "!=":
			return toBoolean(left.(object.Number) != right.(object.Number)), nil
		}
	case left.Type() != right.Type():
		return nil, toError("type mismatch: %s %s %s", left.Type(), node.Operator, right.Type())
	case node.Operator == "==":
		return toBoolean(left == right), nil
	case node.Operator == "!=":
		return toBoolean(left != right), nil
	}

	return nil, toError("unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
}

func evaluateIfExpression(node *ast.IfExpression, env *object.Environment) (object.Object, object.Interruption) {
	condition, interrupt := Evaluate(node.Condition, env)
	if interrupt != nil {
		return nil, interrupt
	}
	if isTruthy(condition) {
		return Evaluate(node.Consequence, env)
	} else if node.Alternative != nil {
		return Evaluate(node.Alternative, env)
	}
	return NULL, nil
}

func evaluateCallExpression(node *ast.CallExpression, env *object.Environment) (object.Object, object.Interruption) {
	value, interrupt := Evaluate(node.Callee, env)
	if interrupt != nil {
		return nil, interrupt
	}

	function, ok := value.(*object.Function)
	if !ok {
		return nil, toError("unknown operator: %s()", function.Type())
	}

	environment := object.NewEnvironment(function.Closure)
	for i, parameter := range function.Declaration.Parameters {
		value, interrupt := Evaluate(node.Arguments[i], env)
		if interrupt != nil {
			return nil, interrupt
		}
		environment.Set(parameter.Value, value)
	}

	result, interrupt := Evaluate(function.Declaration.Body, environment)
	if interrupt != nil {
		if interrupt.Type() == object.RETURN_VALUE {
			return interrupt.(*object.ReturnValue).Value, nil
		}
		return nil, interrupt
	}
	return result, nil
}

func evaluateFunctionLiteral(node *ast.FunctionLiteral, env *object.Environment) (object.Object, object.Interruption) {
	return &object.Function{Declaration: node, Closure: env}, nil
}

func evaluateIdentifier(node *ast.Identifier, env *object.Environment) (object.Object, object.Interruption) {
	value, found := env.Get(node.Value)
	if !found {
		return nil, toError("identifier not found: %s", node.Value)
	}
	return value, nil
}

func evaluateNumberLiteral(node *ast.NumberLiteral) (object.Object, object.Interruption) {
	return object.Number(node.Value), nil
}

func evaluateBooleanLiteral(node *ast.BooleanLiteral) (object.Object, object.Interruption) {
	return toBoolean(node.Value), nil
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
