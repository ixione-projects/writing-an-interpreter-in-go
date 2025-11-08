package parser

import (
	"fmt"
	"strconv"

	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/ast"
	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/lexer"
	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/token"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	tok token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	EQUALITY
	COMPARISON
	TERM
	FACTOR
	UNARY
	CALL
)

func New(input string) *Parser {
	l := lexer.New(input)
	p := &Parser{
		l:      l,
		errors: []string{},
		tok:    l.Token(0),
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.prefixParseFns[token.IDENT] = p.parseIdentifier
	p.prefixParseFns[token.NUMBER] = p.parseNumber
	p.prefixParseFns[token.BANG] = p.parsePrefixExpression
	p.prefixParseFns[token.MINUS] = p.parsePrefixExpression

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	statements := []ast.Statement{}
	var i = 0
	for p.tok.Type != token.EOF && i < 100 {
		p.tok = p.peek0()

		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
		i += 1
	}
	return &ast.Program{Statements: statements}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.tok.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.tok}
	if !p.expect(token.IDENT) {
		return nil
	}

	stmt.Name = p.parseIdentifier().(*ast.Identifier)
	if !p.expect(token.ASSIGN) {
		return nil
	}

	p.next()

	stmt.Value = p.parseExpression(LOWEST)

	p.next()

	if p.tok.Type == token.SEMI {
		p.next()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.tok}

	p.next()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	p.next()

	if p.tok.Type == token.SEMI {
		p.next()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.tok}
	stmt.Expression = p.parseExpression(LOWEST)

	p.next()

	if p.tok.Type == token.SEMI {
		p.next()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.tok.Type]
	if prefix == nil {
		p.error("no prefix parse function defined for %s", p.tok.Type)
		return nil
	}
	left := prefix()

	return left
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Token:    p.tok,
		Operator: p.tok.Literal,
	}

	p.next()

	expr.Right = p.parseExpression(UNARY)

	return expr
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.tok, Value: p.tok.Literal}
}

func (p *Parser) parseNumber() ast.Expression {
	value, err := strconv.ParseFloat(p.tok.Literal, 64)
	if err != nil {
		p.error("cannot parse float %q", p.tok.Literal)
		return nil
	}
	return &ast.Number{Token: p.tok, Value: value}
}

func (p *Parser) expect(ttype token.TokenType) bool {
	if p.peek1().Type == ttype {
		p.next()
		return true
	}
	p.error(
		"expected next token to be <%s> but was <%s>",
		ttype,
		p.peek1().Type,
	)
	return false
}

func (p *Parser) peek0() token.Token {
	return p.l.Token(0)
}

func (p *Parser) peek1() token.Token {
	return p.l.Token(1)
}

func (p *Parser) next() {
	p.l.NextToken()
	p.tok = p.l.Token(0)
}

func (p *Parser) error(format string, args ...any) {
	p.errors = append(p.errors, fmt.Sprintf(format, args...))
}
