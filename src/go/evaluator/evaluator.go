package evaluator

import (
	"fmt"
	"math"
	"slices"

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
	case ast.ASSIGNMENT_EXPRESSION:
		return evaluateAssignmentExpression(node.(*ast.AssignmentExpression), env)
	case ast.CALL_EXPRESSION:
		return evaluateCallExpression(node.(*ast.CallExpression), env)
	case ast.SUBSCRIPT_EXPRESSION:
		return evaluateSubscriptExpression(node.(*ast.SubscriptExpression), env)
	case ast.IDENTIFIER:
		return evaluateIdentifier(node.(*ast.Identifier), env)
	case ast.NUMBER_LITERAL:
		return evaluateNumberLiteral(node.(*ast.NumberLiteral))
	case ast.BOOLEAN_LITERAL:
		return evaluateBooleanLiteral(node.(*ast.BooleanLiteral))
	case ast.STRING_LITERAL:
		return evaluateStringLiteral(node.(*ast.StringLiteral))
	case ast.ARRAY_LITERAL:
		return evaluateArrayLiteral(node.(*ast.ArrayLiteral), env)
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
	case left.Type() == object.STRING && right.Type() == object.STRING:
		switch node.Operator {
		case "+":
			return object.String(left.(object.String) + right.(object.String)), nil
		case "==":
			return toBoolean(left.(object.String) == right.(object.String)), nil
		case "!=":
			return toBoolean(left.(object.String) != right.(object.String)), nil
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

func evaluateAssignmentExpression(node *ast.AssignmentExpression, env *object.Environment) (object.Object, object.Interruption) {
	rvalue, interrupt := Evaluate(node.RValue, env)
	if interrupt != nil {
		return nil, interrupt
	}

	switch lvalue := node.LValue.(type) {
	case *ast.Identifier:
		if _, found := env.Get(lvalue.Value); found {
			env.Set(lvalue.Value, rvalue)
			return rvalue, nil
		}

		return nil, toError("identifier not found: %s", lvalue.Value)
	case *ast.SubscriptExpression:
		baseValue, interrupt := Evaluate(lvalue.Base, env)
		if interrupt != nil {
			return nil, interrupt
		}

		subscriptValue, interrupt := Evaluate(lvalue.Subscript, env)
		if interrupt != nil {
			return nil, interrupt
		}

		subscript, ok := subscriptValue.(object.Number)
		if !ok {
			return nil, toError("unknown operator: %s[%s]", baseValue.Type(), subscriptValue.Type())
		}

		switch base := baseValue.(type) {
		case *object.Array:
			index, valid := toNativeInt(subscript)
			if !valid || index < 0 {
				return nil, toError("unexpected subscript value: %s", subscript.Inspect())
			}

			if index >= cap(base.Elements) {
				current := cap(base.Elements)
				base.Elements = slices.Grow(base.Elements, index+1-current)
				base.Elements = base.Elements[:cap(base.Elements)]
				for i := current; i < index; i += 1 {
					base.Elements[i] = NULL
				}
			} else {
				base.Elements = base.Elements[:cap(base.Elements)]
			}
			base.Elements[index] = rvalue

			return rvalue, nil
		default:
			return nil, toError("unknown operator: %s[%s]", base.Type(), subscript.Type())
		}
	default:
		panic(fmt.Errorf("unknown lvalue type: %s", node.LValue.Type()))
	}
}

func evaluateCallExpression(node *ast.CallExpression, env *object.Environment) (object.Object, object.Interruption) {
	value, interrupt := Evaluate(node.Callee, env)
	if interrupt != nil {
		return nil, interrupt
	}

	switch function := value.(type) {
	case *object.Function:
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
	case *object.Builtin:
		args := []object.Object{}
		for _, arg := range node.Arguments {
			value, interrupt := Evaluate(arg, env)
			if interrupt != nil {
				return nil, interrupt
			}
			args = append(args, value)
		}
		return function.Fn(args...)
	default:
		return nil, toError("unknown operator: %s()", function.Type())
	}
}

func evaluateSubscriptExpression(node *ast.SubscriptExpression, env *object.Environment) (object.Object, object.Interruption) {
	baseValue, interrupt := Evaluate(node.Base, env)
	if interrupt != nil {
		return nil, interrupt
	}

	subscriptValue, interrupt := Evaluate(node.Subscript, env)
	if interrupt != nil {
		return nil, interrupt
	}

	subscript, ok := subscriptValue.(object.Number)
	if !ok {
		return nil, toError("unknown operator: %s[%s]", baseValue.Type(), subscriptValue.Type())
	}

	switch base := baseValue.(type) {
	case *object.Array:
		index, valid := toNativeInt(subscript)
		if !valid || index < 0 || index >= len(base.Elements) {
			return NULL, nil
		}
		return base.Elements[index], nil
	default:
		return nil, toError("unknown operator: %s[%s]", base.Type(), subscript.Type())
	}
}

func evaluateFunctionLiteral(node *ast.FunctionLiteral, env *object.Environment) (object.Object, object.Interruption) {
	return &object.Function{Declaration: node, Closure: env}, nil
}

func evaluateIdentifier(node *ast.Identifier, env *object.Environment) (object.Object, object.Interruption) {
	if value, found := env.Get(node.Value); found {
		return value, nil
	}

	if builtin, found := builtins[node.Value]; found {
		return builtin, nil
	}

	return nil, toError("identifier not found: %s", node.Value)
}

func evaluateNumberLiteral(node *ast.NumberLiteral) (object.Object, object.Interruption) {
	return object.Number(node.Value), nil
}

func evaluateBooleanLiteral(node *ast.BooleanLiteral) (object.Object, object.Interruption) {
	return toBoolean(node.Value), nil
}

func evaluateStringLiteral(node *ast.StringLiteral) (object.Object, object.Interruption) {
	return object.String(node.Value), nil
}

func evaluateArrayLiteral(node *ast.ArrayLiteral, env *object.Environment) (object.Object, object.Interruption) {
	elements := []object.Object{}
	for _, element := range node.Elements {
		value, interrupt := Evaluate(element, env)
		if interrupt != nil {
			return nil, interrupt
		}
		elements = append(elements, value)
	}
	return &object.Array{Elements: elements}, nil
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

func toNativeInt(subscript object.Number) (int, bool) {
	if subscript == object.Number(math.Trunc(float64(subscript))) {
		return int(subscript), true
	}
	return 0, false
}

func isTruthy(o object.Object) object.Boolean {
	switch {
	case o == FALSE:
		return FALSE
	case o == NULL:
		return FALSE
	case o.Type() == object.NUMBER && o.(object.Number) == 0.0:
		return FALSE
	case o.Type() == object.STRING && o.(object.String) == "":
		return FALSE
	case o.Type() == object.ARRAY && len(o.(*object.Array).Elements) == 0:
		return FALSE
	default:
		return TRUE
	}
}
