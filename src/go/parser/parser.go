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

	rules map[token.TokenType]parserRule
}

type precedence int

const (
	_ precedence = iota
	NONE
	ASSIGNMENT
	EQUALITY
	COMPARISON
	TERM
	FACTOR
	UNARY
	CALL
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type parserRule struct {
	PrefixParseFn prefixParseFn
	InfixParseFn  infixParseFn
	Precedence    precedence
}

func New(input string) *Parser {
	l := lexer.New(input)
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	p.tok = p.peek0()

	p.rules = map[token.TokenType]parserRule{
		token.ILLEGAL: {nil, nil, NONE},
		token.EOF:     {nil, nil, NONE},
		token.IDENT:   {p.parseIdentifier, nil, NONE},
		token.NUMBER:  {p.parseNumberLiteral, nil, NONE},
		token.ASSIGN:  {nil, nil, NONE},
		token.PLUS:    {nil, p.parseInfixExpression, TERM},
		token.MINUS:   {p.parsePrefixExpression, p.parseInfixExpression, TERM},
		token.BANG:    {p.parsePrefixExpression, nil, NONE},
		token.STAR:    {nil, p.parseInfixExpression, FACTOR},
		token.SLASH:   {nil, p.parseInfixExpression, FACTOR},
		token.LT:      {nil, p.parseInfixExpression, COMPARISON},
		token.GT:      {nil, p.parseInfixExpression, COMPARISON},
		token.EQ:      {nil, p.parseInfixExpression, EQUALITY},
		token.NOT_EQ:  {nil, p.parseInfixExpression, EQUALITY},
		token.COMMA:   {nil, nil, NONE},
		token.SEMI:    {nil, nil, NONE},
		token.LPAREN:  {nil, nil, NONE},
		token.RPAREN:  {nil, nil, NONE},
		token.LBRACE:  {nil, nil, NONE},
		token.RBRACE:  {nil, nil, NONE},
		token.FN:      {nil, nil, NONE},
		token.LET:     {nil, nil, NONE},
		token.TRUE:    {nil, nil, NONE},
		token.FALSE:   {nil, nil, NONE},
		token.IF:      {nil, nil, NONE},
		token.ELSE:    {nil, nil, NONE},
		token.RETURN:  {nil, nil, NONE},
	}

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	statements := []ast.Statement{}
	for p.tok.Type != token.EOF {
		p.tok = p.peek0()

		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
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

	stmt.Value = p.parseExpression(ASSIGNMENT)

	p.next()

	if p.tok.Type == token.SEMI {
		p.next()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.tok}

	p.next()

	stmt.ReturnValue = p.parseExpression(ASSIGNMENT)

	p.next()

	if p.tok.Type == token.SEMI {
		p.next()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.tok}
	stmt.Expression = p.parseExpression(ASSIGNMENT)

	p.next()

	if p.tok.Type == token.SEMI {
		p.next()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence precedence) ast.Expression {
	prefix := p.getRule(p.tok.Type).PrefixParseFn
	if prefix == nil {
		p.error("no prefix parse function defined for %s", p.tok.Type)
		return nil
	}
	left := prefix()

	for precedence <= p.getRule(p.peek1().Type).Precedence {
		p.next()

		infix := p.getRule(p.tok.Type).InfixParseFn
		if infix == nil {
			return left
		}
		left = infix(left)
	}

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

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.tok,
		Left:     left,
		Operator: p.tok.Literal,
	}
	rule := p.getRule(p.tok.Type)

	p.next()

	expr.Right = p.parseExpression(rule.Precedence + 1) // left associativity

	return expr
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.tok, Value: p.tok.Literal}
}

func (p *Parser) parseNumberLiteral() ast.Expression {
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
	if p.tok = p.l.NextToken(); p.tok.Type != token.EOF {
		p.tok = p.l.Token(0)
	}
}

func (p *Parser) getRule(ttype token.TokenType) parserRule {
	return p.rules[ttype]
}

func (p *Parser) error(format string, args ...any) {
	p.errors = append(p.errors, fmt.Sprintf(format, args...))
}
