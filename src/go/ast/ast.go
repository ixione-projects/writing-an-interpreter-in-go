package ast

import "github.com/ixione-projects/writing-an-interpreter-in-go/src/go/token"

type Node interface {
	TokenLiteral() string
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

type Statement interface {
	Node
	statementNode()
}

func (ls *LetStatement) statementNode() {}

type LetStatement struct {
	Token token.Token
	name  *Identifier
	Value Expression
}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type Expression interface {
	Node
	expressionNode()
}

func (i *Identifier) expressionNode() {}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
