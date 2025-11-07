package lexer

import (
	"slices"

	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/token"
)

type Lexer struct {
	input   string
	start   int
	current int
	ch      byte
}

func New(input string) *Lexer {
	return &Lexer{input: input}
}

func (l *Lexer) NextToken() *token.Token {
	l.start = l.current

	for !l.isEOF() {
		l.ch = l.peek0()

		switch l.ch {
		case ' ', '\t', '\r', '\n':
			l.skip(' ', '\t', '\r', '\n')
			l.start = l.current
		case '=':
			if l.match('=') {
				return l.emit(token.EQ)
			}
			return l.emit(token.ASSIGN)
		case '+':
			return l.emit(token.PLUS)
		case '-':
			return l.emit(token.MINUS)
		case '!':
			if l.match('=') {
				return l.emit(token.NOT_EQ)
			}
			return l.emit(token.BANG)
		case '/':
			return l.emit(token.SLASH)
		case '*':
			return l.emit(token.STAR)
		case '<':
			return l.emit(token.LT)
		case '>':
			return l.emit(token.GT)
		case ';':
			return l.emit(token.SEMI)
		case ',':
			return l.emit(token.COMMA)
		case '(':
			return l.emit(token.LPAREN)
		case ')':
			return l.emit(token.RPAREN)
		case '{':
			return l.emit(token.LBRACE)
		case '}':
			return l.emit(token.RBRACE)
		default:
			if isAlpha(l.ch) {
				return l.ident()
			} else if isNumber(l.ch) {
				return l.number()
			} else {
				return l.emit(token.ILLEGAL)
			}
		}
	}

	return l.emit(token.EOF)
}

func (l *Lexer) ident() *token.Token {
	for isAlphaNumeric(l.ch) {
		l.next()
	}
	ident := l.input[l.start:l.current]
	return &token.Token{
		Type:    token.LookupIdent(ident),
		Literal: ident,
	}
}

func (l *Lexer) number() *token.Token {
	for isNumber(l.ch) {
		l.next()
	}

	if l.ch == '.' {
		if !isNumber(l.peek1()) {
			return &token.Token{
				Type:    token.ILLEGAL,
				Literal: l.input[l.start:l.current],
			}
		}

		l.next()
		for isNumber(l.ch) {
			l.next()
		}
	}

	return &token.Token{
		Type:    token.NUMBER,
		Literal: l.input[l.start:l.current],
	}
}

func (l *Lexer) emit(ttype token.TokenType) *token.Token {
	l.next()
	return &token.Token{
		Type:    ttype,
		Literal: l.input[l.start:l.current],
	}
}

func (l *Lexer) skip(chars ...byte) {
	for slices.Contains(chars, l.ch) {
		l.next()
	}
}

func (l *Lexer) match(char byte) bool {
	if l.peek1() == char {
		l.next()
		return true
	}
	return false
}

func (l *Lexer) peek0() byte {
	if l.isEOF() {
		return 0
	}
	return l.input[l.current]
}

func (l *Lexer) peek1() byte {
	if l.current+1 >= len(l.input) {
		return 0
	}
	return l.input[l.current+1]
}

func (l *Lexer) next() {
	if l.isEOF() {
		return
	}

	l.current += 1
	if l.isEOF() {
		l.ch = 0
	} else {
		l.ch = l.input[l.current]
	}
}

func (l *Lexer) isEOF() bool {
	return l.current >= len(l.input)
}

func isAlphaNumeric(ch byte) bool {
	return isAlpha(ch) || isNumber(ch)
}

func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isNumber(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
