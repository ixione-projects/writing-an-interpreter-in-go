package ast

import (
	"testing"

	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/assert"
	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/token"
)

func TestString(t *testing.T) {
	r := assert.GetTestReporter(t)

	program := &Program{
		Statements: []Statement{
			&LetDeclaration{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "ident"},
					Value: "ident",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "value"},
					Value: "value",
				},
			},
		},
	}

	assert.EqualsWithMessage("let ident=value;", program.String(), "program.String()", r)
}
