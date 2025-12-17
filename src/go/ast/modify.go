package ast

import "fmt"

func Modify(node Node, modifier func(Node) Node) Node {
	switch node := node.(type) {
	case *Program:
		for i, stmt := range node.Statements {
			modified, ok := Modify(stmt, modifier).(Statement)
			if !ok {
				return toErrorNode(Statement(nil), modified)
			}
			node.Statements[i] = modified
		}
	case *LetDeclaration:
		mname, ok := Modify(node.Name, modifier).(*Identifier)
		if !ok {
			return toErrorNode(&Identifier{}, mname)
		}
		node.Name = mname
		mvalue, ok := Modify(node.Value, modifier).(Expression)
		if !ok {
			return toErrorNode(Expression(nil), mvalue)
		}
		node.Value = mvalue
	case *ReturnStatement:
		if node.ReturnValue != nil {
			modified, ok := Modify(node.ReturnValue, modifier).(Expression)
			if !ok {
				return toErrorNode(Expression(nil), modified)
			}
			node.ReturnValue = modified
		}
	case *ExpressionStatement:
		modified, ok := Modify(node.Expression, modifier).(Expression)
		if !ok {
			return toErrorNode(Expression(nil), modified)
		}
		node.Expression = modified
	case *BlockStatement:
		for i, stmt := range node.Statements {
			modified, ok := Modify(stmt, modifier).(Statement)
			if !ok {
				return toErrorNode(Statement(nil), modified)
			}
			node.Statements[i] = modified
		}
	case *UnaryExpression:
		modified, ok := Modify(node.Right, modifier).(Expression)
		if !ok {
			return toErrorNode(Expression(nil), modified)
		}
		node.Right = modified
	case *BinaryExpression:
		modified, ok := Modify(node.Left, modifier).(Expression)
		if !ok {
			return toErrorNode(Expression(nil), modified)
		}
		node.Left = modified
		modified, ok = Modify(node.Right, modifier).(Expression)
		if !ok {
			return toErrorNode(Expression(nil), modified)
		}
		node.Right = modified
	case *LogicalExpression:
		modified, ok := Modify(node.Left, modifier).(Expression)
		if !ok {
			return toErrorNode(Expression(nil), modified)
		}
		node.Left = modified
		modified, ok = Modify(node.Right, modifier).(Expression)
		if !ok {
			return toErrorNode(Expression(nil), modified)
		}
		node.Right = modified
	case *ConditionalExpression:
		condition, ok := Modify(node.Condition, modifier).(Expression)
		if !ok {
			return toErrorNode(Expression(nil), condition)
		}
		node.Condition = condition
		modified, ok := Modify(node.Consequence, modifier).(*BlockStatement)
		if !ok {
			return toErrorNode(&BlockStatement{}, modified)
		}
		node.Consequence = modified
		if node.Alternative != nil {
			modified, ok = Modify(node.Alternative, modifier).(*BlockStatement)
			if !ok {
				return toErrorNode(&BlockStatement{}, modified)
			}
			node.Alternative = modified
		}
	case *FunctionLiteral:
		for i, param := range node.Parameters {
			modified, ok := Modify(param, modifier).(*Identifier)
			if !ok {
				return toErrorNode(&Identifier{}, modified)
			}
			node.Parameters[i] = modified
		}
		modified, ok := Modify(node.Body, modifier).(*BlockStatement)
		if !ok {
			return toErrorNode(&BlockStatement{}, modified)
		}
		node.Body = modified
	case *AssignmentExpression:
		modified, ok := Modify(node.LValue, modifier).(Expression)
		if !ok {
			return toErrorNode(Expression(nil), modified)
		}
		node.LValue = modified
		modified, ok = Modify(node.RValue, modifier).(Expression)
		if !ok {
			return toErrorNode(Expression(nil), modified)
		}
		node.RValue = modified
	case *CallExpression:
		modified, ok := Modify(node.Callee, modifier).(Expression)
		if !ok {
			return toErrorNode(Expression(nil), modified)
		}
		node.Callee = modified
		for i, expr := range node.Arguments {
			modified, ok = Modify(expr, modifier).(Expression)
			if !ok {
				return toErrorNode(Expression(nil), modified)
			}
			node.Arguments[i] = modified
		}
	case *SubscriptExpression:
		modified, ok := Modify(node.Base, modifier).(Expression)
		if !ok {
			return toErrorNode(Expression(nil), modified)
		}
		node.Base = modified
		modified, ok = Modify(node.Subscript, modifier).(Expression)
		if !ok {
			return toErrorNode(Expression(nil), modified)
		}
		node.Subscript = modified
	case *ArrayLiteral:
		for i, expr := range node.Elements {
			modified, ok := Modify(expr, modifier).(Expression)
			if !ok {
				return toErrorNode(Expression(nil), modified)
			}
			node.Elements[i] = modified
		}
	case *HashLiteral:
		pairs := map[Expression]Expression{}
		for i, key := range node.Keys {
			mkey, ok := Modify(key, modifier).(Expression)
			if !ok {
				return toErrorNode(Expression(nil), mkey)
			}
			node.Keys[i] = mkey
			mvalue, ok := Modify(node.Pairs[key], modifier).(Expression)
			if !ok {
				return toErrorNode(Expression(nil), mvalue)
			}
			pairs[mkey] = mvalue
		}
		node.Pairs = pairs
	}

	return modifier(node)
}

func toErrorNode(expected, actual Node) Node {
	if _, ok := actual.(*Error); ok {
		return actual
	}
	return &Error{
		Message: fmt.Sprintf("Modify() ==> unexpected type, expected: <%T> but was: <%T>", expected, actual),
	}
}
