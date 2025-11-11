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
	debug  bool

	tok token.Token

	rules map[token.TokenType]parseRule
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

type parseRule struct {
	PrefixParseFn prefixParseFn
	InfixParseFn  infixParseFn
	Precedence    precedence
}

func NewParser(input string, debug bool) *Parser {
	l := lexer.NewLexer(input)
	p := &Parser{
		l:      l,
		errors: []string{},
		debug:  debug,
	}
	p.tok = p.peek0()

	p.rules = map[token.TokenType]parseRule{
		token.ILLEGAL: {nil, nil, NONE},
		token.EOF:     {nil, nil, NONE},
		token.IDENT:   {p.parseIdentifier, nil, NONE},
		token.NUMBER:  {p.parseNumberLiteral, nil, NONE},
		token.STRING:  {p.parseStringLiteral, nil, NONE},
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
		token.LPAREN:  {p.parseGroupingExpression, p.parseCallExpression, CALL},
		token.RPAREN:  {nil, nil, NONE},
		token.LBRACE:  {nil, nil, NONE},
		token.RBRACE:  {nil, nil, NONE},
		token.FN:      {p.parseFunctionLiteral, nil, NONE},
		token.LET:     {nil, nil, NONE},
		token.TRUE:    {p.parseBooleanLiteral, nil, NONE},
		token.FALSE:   {p.parseBooleanLiteral, nil, NONE},
		token.IF:      {p.parseIfExpression, nil, NONE},
		token.ELSE:    {nil, nil, NONE},
		token.RETURN:  {nil, nil, NONE},
	}

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	if p.debug {
		defer un(trace("ParseProgram"))
	}

	statements := []ast.Statement{}
	for ; p.tok.Type != token.EOF; p.next() {
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
	if p.debug {
		defer un(trace("ParseStatement"))
	}

	switch p.tok.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetDeclaration {
	if p.debug {
		defer un(trace("ParseLetStatement"))
	}

	stmt := &ast.LetDeclaration{Token: p.tok}
	if !p.expect(token.IDENT) {
		return nil
	}

	stmt.Name = p.parseIdentifier().(*ast.Identifier)
	if !p.expect(token.ASSIGN) {
		return nil
	}

	p.next()

	stmt.Value = p.parseExpression(ASSIGNMENT - 1) // right associativity

	if p.peek1().Type == token.SEMI {
		p.next()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	if p.debug {
		defer un(trace("ParseReturnStatement"))
	}

	stmt := &ast.ReturnStatement{Token: p.tok}

	p.next()

	stmt.ReturnValue = p.parseExpression(ASSIGNMENT - 1) // right associativity

	if p.peek1().Type == token.SEMI {
		p.next()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	if p.debug {
		defer un(trace("ParseExpressionStatement"))
	}

	stmt := &ast.ExpressionStatement{Token: p.tok}
	stmt.Expression = p.parseExpression(ASSIGNMENT - 1) // right associativity

	if p.peek1().Type == token.SEMI {
		p.next()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	if p.debug {
		defer un(trace("ParseBlockStatement"))
	}

	block := &ast.BlockStatement{Token: p.tok}

	p.next()

	statements := []ast.Statement{}
	for ; p.tok.Type != token.RBRACE && p.tok.Type != token.EOF; p.next() {
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}
	block.Statements = statements

	if p.tok.Type != token.RBRACE {
		p.error("expected token to be <RBRACE> but was <%s>", p.tok.Type)
		return nil
	}

	return block
}

func (p *Parser) parseExpression(rightPrecedence precedence) ast.Expression {
	if p.debug {
		defer un(trace("ParseExpression"))
	}

	prefix := p.getRule(p.tok.Type).PrefixParseFn
	if prefix == nil {
		p.error("no prefix parse function defined for %s", p.tok.Type)
		return nil
	}
	left := prefix()

	for rightPrecedence < p.getRule(p.peek1().Type).Precedence {
		p.next()

		infix := p.getRule(p.tok.Type).InfixParseFn
		if infix == nil {
			return left
		}
		left = infix(left)
	}

	return left
}

func (p *Parser) parseGroupingExpression() ast.Expression {
	p.next()

	expr := p.parseExpression(ASSIGNMENT - 1) // right associativity
	if !p.expect(token.RPAREN) {
		return nil
	}

	return expr
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	if p.debug {
		defer un(trace("ParsePrefixExpression"))
	}

	expr := &ast.PrefixExpression{
		Token:    p.tok,
		Operator: p.tok.Literal,
	}

	p.next()

	expr.Right = p.parseExpression(UNARY)

	return expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	if p.debug {
		defer un(trace("ParseInfixExpression"))
	}

	expr := &ast.InfixExpression{
		Token:    p.tok,
		Left:     left,
		Operator: p.tok.Literal,
	}
	rule := p.getRule(p.tok.Type)

	p.next()

	expr.Right = p.parseExpression(rule.Precedence)

	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	if p.debug {
		defer un(trace("ParseIfExpression"))
	}

	expr := &ast.IfExpression{Token: p.tok}
	if !p.expect(token.LPAREN) {
		return nil
	}

	p.next()

	expr.Condition = p.parseExpression(ASSIGNMENT - 1) // right associativity

	if !p.expect(token.RPAREN) {
		return nil
	}

	if !p.expect(token.LBRACE) {
		return nil
	}

	expr.Consequence = p.parseBlockStatement()
	if p.peek1().Type == token.ELSE {
		p.next()

		if !p.expect(token.LBRACE) {
			return nil
		}
		expr.Alternative = p.parseBlockStatement()
	}

	return expr
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	if p.debug {
		defer un(trace("ParseFunctionLiteral"))
	}

	expr := &ast.FunctionLiteral{Token: p.tok}
	if !p.expect(token.LPAREN) {
		return nil
	}

	if p.peek1().Type == token.RPAREN {
		expr.Parameters = []*ast.Identifier{}
	} else {
		p.next()

		expr.Parameters = p.parseParameters()
	}

	if !p.expect(token.RPAREN) {
		return nil
	}

	if !p.expect(token.LBRACE) {
		return nil
	}

	expr.Body = p.parseBlockStatement()

	return expr
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	if p.debug {
		defer un(trace("ParseCallExpression"))
	}

	expr := &ast.CallExpression{
		Token:  p.tok,
		Callee: left,
	}

	if p.peek1().Type == token.RPAREN {
		expr.Arguments = []ast.Expression{}
	} else {
		p.next()

		expr.Arguments = p.parseArguments()
	}

	if !p.expect(token.RPAREN) {
		return nil
	}

	return expr
}

func (p *Parser) parseIdentifier() ast.Expression {
	if p.debug {
		defer un(trace("ParseIdentifier"))
	}

	return &ast.Identifier{Token: p.tok, Value: p.tok.Literal}
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	if p.debug {
		defer un(trace("ParseNumberLiteral"))
	}

	value, err := strconv.ParseFloat(p.tok.Literal, 64)
	if err != nil {
		p.error("cannot parse float %q", p.tok.Literal)
		return nil
	}
	return &ast.NumberLiteral{Token: p.tok, Value: value}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	if p.debug {
		defer un(trace("ParseStringLiteral"))
	}

	return &ast.StringLiteral{
		Token: p.tok,
		Value: p.tok.Literal[1 : len(p.tok.Literal)-1],
	}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	if p.debug {
		defer un(trace("ParseBooleanLiteral"))
	}

	return &ast.BooleanLiteral{Token: p.tok, Value: p.tok.Type == token.TRUE}
}

func (p *Parser) parseParameters() []*ast.Identifier {
	idents := []*ast.Identifier{}

	ident := p.parseIdentifier().(*ast.Identifier)
	if ident != nil {
		idents = append(idents, ident)
	}

	for p.peek1().Type != token.RPAREN {
		if !p.expect(token.COMMA) {
			return nil
		}

		p.next()

		ident := p.parseIdentifier().(*ast.Identifier)
		if ident != nil {
			idents = append(idents, ident)
		}
	}

	return idents
}

func (p *Parser) parseArguments() []ast.Expression {
	args := []ast.Expression{}

	arg := p.parseExpression(ASSIGNMENT - 1) // right associativity
	if arg != nil {
		args = append(args, arg)
	}

	for p.peek1().Type != token.RPAREN {
		if !p.expect(token.COMMA) {
			return nil
		}

		p.next()

		arg := p.parseExpression(ASSIGNMENT - 1) // right associativity
		if arg != nil {
			args = append(args, arg)
		}
	}

	return args
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
	if p.l.NextToken(); p.tok.Type != token.EOF {
		p.tok = p.l.Token(0)
	}
}

func (p *Parser) getRule(ttype token.TokenType) parseRule {
	return p.rules[ttype]
}

func (p *Parser) error(format string, args ...any) {
	p.errors = append(p.errors, fmt.Sprintf(format, args...))
}
