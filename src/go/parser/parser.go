package parser

import (
	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/ast"
	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/lexer"
	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/token"
)

type Parser struct {
	l *lexer.Lexer

	tok token.Token
}

func New(input string) *Parser {
	l := lexer.New(input)
	return &Parser{
		l:   l,
		tok: l.Token(0),
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	statements := []ast.Statement{}
	for p.tok.Type != token.EOF {
		p.tok = p.l.Token(0)
	}
	return &ast.Program{Statements: statements}
}
