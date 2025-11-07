package lexer

import (
	"testing"

	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/token"
)

func TestNextToken(t *testing.T) {
	tests := []struct {
		input  string
		tokens []struct {
			Type    token.TokenType
			Literal string
		}
	}{
		{
			input: `=+(){},;`,
			tokens: []struct {
				Type    token.TokenType
				Literal string
			}{
				{token.ASSIGN, "="},
				{token.PLUS, "+"},
				{token.LPAREN, "("},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.RBRACE, "}"},
				{token.COMMA, ","},
				{token.SEMI, ";"},
				{token.EOF, ""},
			},
		},
		{
			input: `
			let five = 5;
			let ten = 10;

			let add = fn(x, y) {
				x + y;
			};

			let result = add(five, ten);`,
			tokens: []struct {
				Type    token.TokenType
				Literal string
			}{
				{token.LET, "let"},
				{token.IDENT, "five"},
				{token.ASSIGN, "="},
				{token.NUMBER, "5"},
				{token.SEMI, ";"},
				{token.LET, "let"},
				{token.IDENT, "ten"},
				{token.ASSIGN, "="},
				{token.NUMBER, "10"},
				{token.SEMI, ";"},
				{token.LET, "let"},
				{token.IDENT, "add"},
				{token.ASSIGN, "="},
				{token.FN, "fn"},
				{token.LPAREN, "("},
				{token.IDENT, "x"},
				{token.COMMA, ","},
				{token.IDENT, "y"},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.IDENT, "x"},
				{token.PLUS, "+"},
				{token.IDENT, "y"},
				{token.SEMI, ";"},
				{token.RBRACE, "}"},
				{token.SEMI, ";"},
				{token.LET, "let"},
				{token.IDENT, "result"},
				{token.ASSIGN, "="},
				{token.IDENT, "add"},
				{token.LPAREN, "("},
				{token.IDENT, "five"},
				{token.COMMA, ","},
				{token.IDENT, "ten"},
				{token.RPAREN, ")"},
				{token.SEMI, ";"},
				{token.EOF, ""},
			},
		},
		{
			input: `!-/*5;`,
			tokens: []struct {
				Type    token.TokenType
				Literal string
			}{
				{token.BANG, "!"},
				{token.MINUS, "-"},
				{token.SLASH, "/"},
				{token.STAR, "*"},
				{token.NUMBER, "5"},
				{token.SEMI, ";"},
				{token.EOF, ""},
			},
		},
		{
			input: `5 < 10 > 5;`,
			tokens: []struct {
				Type    token.TokenType
				Literal string
			}{
				{token.NUMBER, "5"},
				{token.LT, "<"},
				{token.NUMBER, "10"},
				{token.GT, ">"},
				{token.NUMBER, "5"},
				{token.SEMI, ";"},
				{token.EOF, ""},
			},
		},
		{
			input: `
			if (5 < 10) {	
				return true;
			} else {
				return false;
			}`,
			tokens: []struct {
				Type    token.TokenType
				Literal string
			}{
				{token.IF, "if"},
				{token.LPAREN, "("},
				{token.NUMBER, "5"},
				{token.LT, "<"},
				{token.NUMBER, "10"},
				{token.RPAREN, ")"},
				{token.LBRACE, "{"},
				{token.RETURN, "return"},
				{token.TRUE, "true"},
				{token.SEMI, ";"},
				{token.RBRACE, "}"},
				{token.ELSE, "else"},
				{token.LBRACE, "{"},
				{token.RETURN, "return"},
				{token.FALSE, "false"},
				{token.SEMI, ";"},
				{token.RBRACE, "}"},
				{token.EOF, ""},
			},
		},
		{
			input: `
			10 == 10;
			10 != 9;`,
			tokens: []struct {
				Type    token.TokenType
				Literal string
			}{
				{token.NUMBER, "10"},
				{token.EQ, "=="},
				{token.NUMBER, "10"},
				{token.SEMI, ";"},
				{token.NUMBER, "10"},
				{token.NOT_EQ, "!="},
				{token.NUMBER, "9"},
				{token.SEMI, ";"},
				{token.EOF, ""},
			},
		},
	}

	for i, test := range tests {
		l := New(test.input)
		for j, expected := range test.tokens {
			actual := l.NextToken()
			if expected.Type != actual.Type {
				t.Errorf("test[%d][%d] - wrong type ==> expected: <%q> but was: <%q>", i, j, expected.Type, actual.Type)
			}

			if expected.Literal != actual.Literal {
				t.Errorf("test[%d][%d] - wrong literal ==> expected: <%q> but was: <%q>", i, j, expected.Literal, actual.Literal)
			}
		}
	}
}
