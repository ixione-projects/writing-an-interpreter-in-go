package ast

import (
	"testing"

	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/token"
)

var one = func() *NumberLiteral {
	return &NumberLiteral{
		token.Token{Type: token.NUMBER, Literal: "1"},
		1,
	}
}
var toTwo = func(node Node) Node {
	num, ok := node.(*NumberLiteral)
	if !ok {
		return node
	}

	if num.Value == 2 {
		return node
	}

	num.Value = 2
	num.Token.Literal = "2"
	return node
}

func TestModify(t *testing.T) {
	tests := []struct {
		input    Node
		modifier func(Node) Node
		output   string
	}{
		{
			input:    one(),
			modifier: toTwo,
			output:   "2",
		},
		{
			input: &Program{
				[]Statement{
					&ExpressionStatement{
						token.Token{Type: token.NUMBER, Literal: "1"},
						one(),
					},
				},
			},
			modifier: toTwo,
			output:   "2;",
		},
		{
			input: &BinaryExpression{
				token.Token{Type: token.PLUS, Literal: "+"},
				one(),
				"+",
				one(),
			},
			modifier: toTwo,
			output:   "(2+2)",
		},
		{
			input: &BinaryExpression{
				token.Token{Type: token.PLUS, Literal: "+"},
				one(),
				"+",
				one(),
			},
			modifier: toTwo,
			output:   "(2+2)",
		},
		{
			input: &UnaryExpression{
				token.Token{Type: token.MINUS, Literal: "-"},
				"-",
				one(),
			},
			modifier: toTwo,
			output:   "(-2)",
		},
		{
			input: &SubscriptExpression{
				token.Token{Type: token.LBRACK, Literal: "["},
				one(),
				one(),
			},
			modifier: toTwo,
			output:   "(2[2])",
		},
		{
			input: &ConditionalExpression{
				token.Token{Type: token.IF, Literal: "if"},
				one(),
				&BlockStatement{
					token.Token{Type: token.LBRACE, Literal: "{"},
					[]Statement{
						&ExpressionStatement{
							token.Token{Type: token.NUMBER, Literal: "1"},
							one(),
						},
					},
				},
				&BlockStatement{
					token.Token{Type: token.LBRACE, Literal: "{"},
					[]Statement{
						&ExpressionStatement{
							token.Token{Type: token.NUMBER, Literal: "1"},
							one(),
						},
					},
				},
			},
			modifier: toTwo,
			output:   "if 2 {2;} else {2;}",
		},
		{
			input: &ReturnStatement{
				token.Token{Type: token.RETURN, Literal: "return"},
				one(),
			},
			modifier: toTwo,
			output:   "return 2;",
		},
		{
			input: &LetDeclaration{
				token.Token{Type: token.LET, Literal: "let"},
				&Identifier{
					token.Token{Type: token.IDENT, Literal: "ident"},
					"ident",
				},
				one(),
			},
			modifier: toTwo,
			output:   "let ident=2;",
		},
		{
			input: &FunctionLiteral{
				token.Token{Type: token.FN, Literal: "fn"},
				[]*Identifier{},
				&BlockStatement{
					token.Token{Type: token.LBRACE, Literal: "{"},
					[]Statement{
						&ExpressionStatement{
							token.Token{Type: token.NUMBER, Literal: "1"},
							one(),
						},
					},
				},
			},
			modifier: toTwo,
			output:   "fn(){2;}",
		},
		{
			input: &ArrayLiteral{
				token.Token{Type: token.LBRACK, Literal: "["},
				[]Expression{
					one(),
					one(),
				},
			},
			modifier: toTwo,
			output:   "[2,2]",
		},
	}

	for i, test := range tests {
		actual := Modify(test.input, test.modifier)
		if test.output != actual.String() {
			t.Fatalf("test[%d] - Modify() ==> expected: <%s> but was: <%s>", i, test.output, actual.String())
		}
	}
}

func TestModifyHashLiteral(t *testing.T) {
	tests := []struct {
		input    map[Expression]Expression
		modifier func(Node) Node
		output   string
	}{
		{
			input: map[Expression]Expression{
				one(): one(),
				one(): one(),
			},
			modifier: toTwo,
			output:   "{2:2,2:2}",
		},
	}

	for i, test := range tests {
		input := &HashLiteral{
			Token: token.Token{Type: token.LBRACE, Literal: "{"},
			Keys:  []Expression{},
			Pairs: map[Expression]Expression{},
		}
		for key, value := range test.input {
			input.Keys = append(input.Keys, key)
			input.Pairs[key] = value
		}
		actual := Modify(input, test.modifier)
		if test.output != actual.String() {
			t.Fatalf("test[%d] - Modify() ==> expected: <%s> but was: <%s>", i, test.output, actual.String())
		}
	}
}
