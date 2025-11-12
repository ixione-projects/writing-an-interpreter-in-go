package ast

import (
	"bytes"
	"strings"

	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/token"
)

type NodeType int

const (
	PROGRAM NodeType = iota
	LET_DECLARATION
	RETURN_STATEMENT
	EXPRESSION_STATEMENT
	BLOCK_STATEMENT
	PREFIX_EXPRESSION
	INFIX_EXPRESSION
	IF_EXPRESSION
	FUNCTION_LITERAL
	CALL_EXPRESSION
	ASSIGNMENT_EXPRESSION
	SUBSCRIPT_EXPRESSION
	IDENTIFIER
	NUMBER_LITERAL
	BOOLEAN_LITERAL
	STRING_LITERAL
	ARRAY_LITERAL
	HASH_LITERAL
	NULL_LITERAL
)

type Node interface {
	TokenLiteral() string

	Type() NodeType
	String() string
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) <= 0 {
		return ""
	}
	return p.Statements[0].TokenLiteral()
}

func (p *Program) Type() NodeType {
	return PROGRAM
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type Statement interface {
	Node
	statementNode()
}

func (ls *LetDeclaration) statementNode()      {}
func (rs *ReturnStatement) statementNode()     {}
func (es *ExpressionStatement) statementNode() {}
func (bs *BlockStatement) statementNode()      {}

type LetDeclaration struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetDeclaration) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetDeclaration) Type() NodeType {
	return LET_DECLARATION
}

func (ls *LetDeclaration) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ls.Name.String())
	out.WriteString("=")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) Type() NodeType {
	return RETURN_STATEMENT
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral())
	out.WriteString(" ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) Type() NodeType {
	return EXPRESSION_STATEMENT
}

func (es *ExpressionStatement) String() string {
	var out bytes.Buffer

	if es.Expression != nil {
		out.WriteString(es.Expression.String())
	}
	out.WriteString(";")

	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) Type() NodeType {
	return BLOCK_STATEMENT
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString("}")
	return out.String()
}

type Expression interface {
	Node
	expressionNode()
}

func (pe *PrefixExpression) expressionNode()     {}
func (ie *InfixExpression) expressionNode()      {}
func (ie *IfExpression) expressionNode()         {}
func (fl *FunctionLiteral) expressionNode()      {}
func (ce *CallExpression) expressionNode()       {}
func (ae *AssignmentExpression) expressionNode() {}
func (ie *SubscriptExpression) expressionNode()  {}
func (i *Identifier) expressionNode()            {}
func (nl *NumberLiteral) expressionNode()        {}
func (bl *BooleanLiteral) expressionNode()       {}
func (sl *StringLiteral) expressionNode()        {}
func (al *ArrayLiteral) expressionNode()         {}
func (hl *HashLiteral) expressionNode()          {}
func (nl *NullLiteral) expressionNode()          {}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) Type() NodeType {
	return PREFIX_EXPRESSION
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InfixExpression) Type() NodeType {
	return INFIX_EXPRESSION
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(ie.Operator)
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IfExpression) Type() NodeType {
	return IF_EXPRESSION
}

func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ie.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) Type() NodeType {
	return FUNCTION_LITERAL
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, param := range fl.Parameters {
		params = append(params, param.Value)
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token
	Callee    Expression
	Arguments []Expression
}

func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExpression) Type() NodeType {
	return CALL_EXPRESSION
}

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString(ce.Callee.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ","))
	out.WriteString(")")

	return out.String()
}

type AssignmentExpression struct {
	Token  token.Token
	LValue Expression
	RValue Expression
}

func (ae *AssignmentExpression) TokenLiteral() string {
	return ae.Token.Literal
}

func (ae *AssignmentExpression) Type() NodeType {
	return ASSIGNMENT_EXPRESSION
}

func (ae *AssignmentExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ae.LValue.String())
	out.WriteString("=")
	out.WriteString(ae.RValue.String())
	out.WriteString(")")

	return out.String()
}

type SubscriptExpression struct {
	Token     token.Token
	Base      Expression
	Subscript Expression
}

func (ie *SubscriptExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *SubscriptExpression) Type() NodeType {
	return SUBSCRIPT_EXPRESSION
}

func (ie *SubscriptExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Base.String())
	out.WriteString("[")
	out.WriteString(ie.Subscript.String())
	out.WriteString("])")

	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) Type() NodeType {
	return IDENTIFIER
}

func (i *Identifier) String() string {
	return i.Value
}

type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (nl *NumberLiteral) TokenLiteral() string {
	return nl.Token.Literal
}

func (nl *NumberLiteral) Type() NodeType {
	return NUMBER_LITERAL
}

func (nl *NumberLiteral) String() string {
	return nl.Token.Literal
}

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) TokenLiteral() string {
	return bl.Token.Literal
}

func (bl *BooleanLiteral) Type() NodeType {
	return BOOLEAN_LITERAL
}

func (bl *BooleanLiteral) String() string {
	return bl.Token.Literal
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) Type() NodeType {
	return STRING_LITERAL
}

func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

func (al *ArrayLiteral) Type() NodeType {
	return ARRAY_LITERAL
}

func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elems := []string{}
	for _, elem := range al.Elements {
		elems = append(elems, elem.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elems, ","))
	out.WriteString("]")

	return out.String()
}

type HashLiteral struct {
	Token token.Token
	Keys  []Expression
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) TokenLiteral() string {
	return hl.Token.Literal
}

func (hl *HashLiteral) Type() NodeType {
	return HASH_LITERAL
}

func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, key := range hl.Keys {
		pairs = append(pairs, key.String()+":"+hl.Pairs[key].String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ","))
	out.WriteString("}")

	return out.String()
}

type NullLiteral struct {
	Token token.Token
}

func (nl *NullLiteral) TokenLiteral() string {
	return nl.Token.Literal
}

func (nl *NullLiteral) Type() NodeType {
	return NULL_LITERAL
}

func (nl *NullLiteral) String() string {
	return nl.Token.Literal
}

var nodes = map[NodeType]string{
	PROGRAM:               "PROGRAM",
	LET_DECLARATION:       "LET_DECLARATION",
	RETURN_STATEMENT:      "RETURN_STATEMENT",
	EXPRESSION_STATEMENT:  "EXPRESSION_STATEMENT",
	BLOCK_STATEMENT:       "BLOCK_STATEMENT",
	PREFIX_EXPRESSION:     "PREFIX_EXPRESSION",
	INFIX_EXPRESSION:      "INFIX_EXPRESSION",
	IF_EXPRESSION:         "IF_EXPRESSION",
	FUNCTION_LITERAL:      "FUNCTION_LITERAL",
	CALL_EXPRESSION:       "CALL_EXPRESSION",
	ASSIGNMENT_EXPRESSION: "ASSIGNMENT_EXPRESSION",
	SUBSCRIPT_EXPRESSION:  "SUBSCRIPT_EXPRESSION",
	IDENTIFIER:            "IDENTIFIER",
	NUMBER_LITERAL:        "NUMBER_LITERAL",
	BOOLEAN_LITERAL:       "BOOLEAN_LITERAL",
	STRING_LITERAL:        "STRING_LITERAL",
	ARRAY_LITERAL:         "ARRAY_LITERAL",
	HASH_LITERAL:          "HASH_LITERAL",
	NULL_LITERAL:          "NULL_LITERAL",
}

func (nt NodeType) String() string {
	return nodes[nt]
}
