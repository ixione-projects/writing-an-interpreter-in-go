package lexer

import (
	"testing"

	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/token"
)

func TestNextToken(t *testing.T) {
	tests := []struct {
		input  string
		tokens []token.Token
	}{
		{
			input: `=+(){},;`,
			tokens: []token.Token{
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.SEMI, Literal: ";"},
				{Type: token.EOF, Literal: ""},
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
			tokens: []token.Token{
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENT, Literal: "five"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.NUMBER, Literal: "5"},
				{Type: token.SEMI, Literal: ";"},
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENT, Literal: "ten"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.NUMBER, Literal: "10"},
				{Type: token.SEMI, Literal: ";"},
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENT, Literal: "add"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.FN, Literal: "fn"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.IDENT, Literal: "x"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.IDENT, Literal: "y"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.IDENT, Literal: "x"},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.IDENT, Literal: "y"},
				{Type: token.SEMI, Literal: ";"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.SEMI, Literal: ";"},
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENT, Literal: "result"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.IDENT, Literal: "add"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.IDENT, Literal: "five"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.IDENT, Literal: "ten"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.SEMI, Literal: ";"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			input: `!-/*5;`,
			tokens: []token.Token{
				{Type: token.BANG, Literal: "!"},
				{Type: token.MINUS, Literal: "-"},
				{Type: token.SLASH, Literal: "/"},
				{Type: token.STAR, Literal: "*"},
				{Type: token.NUMBER, Literal: "5"},
				{Type: token.SEMI, Literal: ";"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			input: `5 < 10 > 5;`,
			tokens: []token.Token{
				{Type: token.NUMBER, Literal: "5"},
				{Type: token.LT, Literal: "<"},
				{Type: token.NUMBER, Literal: "10"},
				{Type: token.GT, Literal: ">"},
				{Type: token.NUMBER, Literal: "5"},
				{Type: token.SEMI, Literal: ";"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			input: `
			if (5 < 10) {	
				return true;
			} else {
				return false;
			}`,
			tokens: []token.Token{
				{Type: token.IF, Literal: "if"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.NUMBER, Literal: "5"},
				{Type: token.LT, Literal: "<"},
				{Type: token.NUMBER, Literal: "10"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.RETURN, Literal: "return"},
				{Type: token.TRUE, Literal: "true"},
				{Type: token.SEMI, Literal: ";"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.ELSE, Literal: "else"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.RETURN, Literal: "return"},
				{Type: token.FALSE, Literal: "false"},
				{Type: token.SEMI, Literal: ";"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			input: `
			10 == 10;
			10 != 9;`,
			tokens: []token.Token{
				{Type: token.NUMBER, Literal: "10"},
				{Type: token.EQ, Literal: "=="},
				{Type: token.NUMBER, Literal: "10"},
				{Type: token.SEMI, Literal: ";"},
				{Type: token.NUMBER, Literal: "10"},
				{Type: token.NOT_EQ, Literal: "!="},
				{Type: token.NUMBER, Literal: "9"},
				{Type: token.SEMI, Literal: ";"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			input: `
			"foobar"
			"foo bar"`,
			tokens: []token.Token{
				{Type: token.STRING, Literal: "\"foobar\""},
				{Type: token.STRING, Literal: "\"foo bar\""},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			input: `[1, 2];`,
			tokens: []token.Token{
				{Type: token.LBRACK, Literal: "["},
				{Type: token.NUMBER, Literal: "1"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.NUMBER, Literal: "2"},
				{Type: token.RBRACK, Literal: "]"},
				{Type: token.SEMI, Literal: ";"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			input: `{"foo": "bar"}`,
			tokens: []token.Token{
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.STRING, Literal: "\"foo\""},
				{Type: token.COLON, Literal: ":"},
				{Type: token.STRING, Literal: "\"bar\""},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.EOF, Literal: ""},
			},
		},
	}

	for i, test := range tests {
		l := NewLexer(test.input)
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

func TestToken(t *testing.T) {
	tests := []struct {
		input     string
		index     int
		lookahead token.Token
		tokens    []token.Token
	}{
		{
			input: `=+(){},;`,
			index: 0,
			lookahead: token.Token{
				Type:    token.ASSIGN,
				Literal: "=",
			},
			tokens: []token.Token{
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.SEMI, Literal: ";"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			input: `=+(){},;`,
			index: 4,
			lookahead: token.Token{
				Type:    token.LBRACE,
				Literal: "{",
			},
			tokens: []token.Token{
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.SEMI, Literal: ";"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			input: `=+(){},;`,
			index: 8,
			lookahead: token.Token{
				Type:    token.EOF,
				Literal: "",
			},
			tokens: []token.Token{
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.SEMI, Literal: ";"},
				{Type: token.EOF, Literal: ""},
			},
		},
	}

	for i, test := range tests {
		l := NewLexer(test.input)
		token := l.Token(test.index)
		if test.lookahead.Type != token.Type {
			t.Errorf("test[%d] - wrong type ==> expected: <%q> but was: <%q>", i, test.lookahead.Type, token.Type)
		}

		if test.lookahead.Literal != token.Literal {
			t.Errorf("test[%d] - wrong literal ==> expected: <%q> but was: <%q>", i, test.lookahead.Literal, token.Literal)
		}

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
