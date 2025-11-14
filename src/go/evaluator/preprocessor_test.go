package evaluator

import (
	"testing"

	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/ast"
	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/object"
	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/parser"
)

type MacroTest struct {
	name   string
	object string
}

func TestDefineMacros(t *testing.T) {
	tests := []struct {
		input  string
		macros []MacroTest
	}{
		{
			input: `
			let number = 1;
			let function = fn(x, y) { x + y; };
			macro add(x, y) { x + y; };`,
			macros: []MacroTest{
				{
					"add", "<macro add(x, y)>",
				},
			},
		},
	}

	for i, test := range tests {
		env := object.NewEnvironment(nil)
		p := parser.NewParser(test.input, false)
		program := p.ParseProgram()

		if 0 != len(p.Errors()) {
			t.Errorf("test[%d] - len(p.Errors()) ==> expected: <%d> but was: <%d>", i, 0, len(p.Errors()))
			for j, msg := range p.Errors() {
				t.Errorf("--------- p.Errors()[%d]: %s", j, msg)
			}
			t.Fatalf("test[%d] - %s", i, test.input)
		}

		DefineMacros(program, env)

		if len(test.macros) != env.Length() {
			t.Errorf("test[%d] - env.Length() ==> expected: <%d> but was: <%d>", i, len(test.macros), env.Length())
		}

		for j, expected := range test.macros {
			actual, ok := env.Get(expected.name)
			if !ok {
				t.Errorf("test[%d][%d] - env.Get(macro.name) ==> expected: <%t> but was: <%t>", i, j, true, ok)
			}

			if expected.object != actual.Inspect() {
				t.Errorf("test[%d][%d] - env.Length() ==> expected: <%d> but was: <%d>", i, j, len(test.macros), env.Length())
			}
		}
	}
}

func TestExpandMacros(t *testing.T) {
	tests := []struct {
		input  string
		object string
	}{
		{
			input: `
			macro infix() { quote(1 + 2); };

			infix();`,
			object: "(1+2);",
		},
		{
			input: `
			macro reverse(a, b) { quote(unquote(b) - unquote(a)); };

			reverse(2 + 2, 10 - 5);`,
			object: "((10-5)-(2+2));",
		},
		{
			input: `
            macro unless(condition, consequence, alternative) {
                quote(if (!(unquote(condition))) {
                    unquote(consequence);
                } else {
                    unquote(alternative);
                });
            };

            unless(10 > 5, puts("not greater"), puts("greater"));
            `,
			object: "if (!(10>5)) {puts(\"not greater\");} else {puts(\"greater\");};",
		},
	}

	for i, test := range tests {
		env := object.NewEnvironment(nil)
		p := parser.NewParser(test.input, false)
		program := p.ParseProgram()

		if 0 != len(p.Errors()) {
			t.Errorf("test[%d] - len(p.Errors()) ==> expected: <%d> but was: <%d>", i, 0, len(p.Errors()))
			for j, msg := range p.Errors() {
				t.Errorf("--------- p.Errors()[%d]: %s", j, msg)
			}
			t.Fatalf("test[%d] - %s", i, test.input)
		}

		program = DefineMacros(program, env)
		program = ExpandMacros(program, env).(*ast.Program)

		if test.object != program.String() {
			t.Errorf("test[%d] - program.String() ==> expected: <%s> but was: <%s>", i, test.object, program.String())
		}
	}
}
