package parser

import (
	"fmt"
	"testing"

	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/ast"
)

type ParserTest struct {
	input   string
	errors  []string
	program ProgramTest
}

type NodeTest interface {
	node()
}

func (p ProgramTest) node()              {}
func (ls LetStatementTest) node()        {}
func (rs ReturnStatementTest) node()     {}
func (es ExpressionStatementTest) node() {}
func (pe PrefixExpressionTest) node()    {}
func (i IdentifierTest) node()           {}
func (n NumberTest) node()               {}

type ProgramTest struct {
	Statements []StatementTest
}

type StatementTest interface {
	NodeTest
	statementNode()
}

func (ls LetStatementTest) statementNode()        {}
func (rs ReturnStatementTest) statementNode()     {}
func (es ExpressionStatementTest) statementNode() {}

type LetStatementTest struct {
	Name  IdentifierTest
	Value ExpressionTest
}

type ReturnStatementTest struct {
	ReturnValue ExpressionTest
}

type ExpressionStatementTest struct {
	TokenLiteral string
	Expression   ExpressionTest
}

type ExpressionTest interface {
	NodeTest
	expressionNode()
}

func (pe PrefixExpressionTest) expressionNode() {}
func (i IdentifierTest) expressionNode()        {}
func (n NumberTest) expressionNode()            {}

type PrefixExpressionTest struct {
	Operator string
	Right    ExpressionTest
}

type IdentifierTest string
type NumberTest float64

func TestLetStatement(t *testing.T) {
	tests := []ParserTest{
		{
			input: `
			let x = 5;
			let y = 10;
			let foobar = 838383;`,
			program: ProgramTest{
				[]StatementTest{
					LetStatementTest{"x", NumberTest(5)},
					LetStatementTest{"y", NumberTest(10)},
					LetStatementTest{"foobar", NumberTest(838383)},
				},
			},
		},
		// {
		// 	input: `
		// 	let x 5;
		// 	let y 10;
		// 	let 838383;`,
		// 	errors: []string{
		// 		"expected next token to be <ASSIGN> but was <NUMBER>",
		// 		"expected next token to be <ASSIGN> but was <NUMBER>",
		// 		"expected next token to be <IDENT> but was <NUMBER>",
		// 	},
		// },
	}

	for i, test := range tests {
		testParser(t, i, test)
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []ParserTest{
		{
			input: `
			return 5;
			return 10;
			return 993322;`,
			program: ProgramTest{
				[]StatementTest{
					ReturnStatementTest{NumberTest(5)},
					ReturnStatementTest{NumberTest(10)},
					ReturnStatementTest{NumberTest(993322)},
				},
			},
		},
	}

	for i, test := range tests {
		testParser(t, i, test)
	}
}

func TestExpressionStatement(t *testing.T) {
	tests := []ParserTest{
		{
			input: `foobar;`,
			program: ProgramTest{
				[]StatementTest{
					ExpressionStatementTest{
						"foobar",
						IdentifierTest("foobar"),
					},
				},
			},
		},
		{
			input: `5;`,
			program: ProgramTest{
				[]StatementTest{
					ExpressionStatementTest{
						"5",
						NumberTest(5),
					},
				},
			},
		},
		{
			input: `!5;`,
			program: ProgramTest{
				[]StatementTest{
					ExpressionStatementTest{
						"!",
						PrefixExpressionTest{"!", NumberTest(5)},
					},
				},
			},
		},
		{
			input: `-15;`,
			program: ProgramTest{
				[]StatementTest{
					ExpressionStatementTest{
						"-",
						PrefixExpressionTest{"-", NumberTest(15)},
					},
				},
			},
		},
	}

	for i, test := range tests {
		testParser(t, i, test)
	}
}

func testParser(t *testing.T, i int, test ParserTest) {
	p := New(test.input)
	program := p.ParseProgram()

	fmt.Printf("%s\n", program.String())
	if len(test.errors) != len(p.Errors()) {
		t.Errorf("test[%d] - len(p.Errors()) ==> expected: <%d> but was: <%d>", i, len(test.errors), len(p.Errors()))
		for j, msg := range p.Errors() {
			t.Errorf("--------- p.Errors()[%d]: %s", j, msg)
		}
		t.FailNow()
	}

	for j, msg := range test.errors {
		fail := false
		if msg != p.Errors()[j] {
			t.Errorf("test[%d] - p.Errors()[%d] ==> expected: <%s> but was: <%s>", i, j, msg, p.Errors()[j])
			fail = true
		}

		if fail {
			t.FailNow()
		}
	}

	if len(test.program.Statements) > 0 {
		testProgram(t, i, test.program, program)
	}
}

func testProgram(t *testing.T, i int, expected ProgramTest, actual *ast.Program) {
	if actual == nil {
		t.Fatalf("test[%d] - ParseProgram() ==> expected: not <%#v>", i, actual)
	}

	if len(expected.Statements) != len(actual.Statements) {
		t.Fatalf("test[%d] - len(program.Statements) ==> expected: <%d> but was: <%d>", i, len(expected.Statements), len(actual.Statements))
	}

	for j, stmt := range expected.Statements {
		switch stmt := stmt.(type) {
		case LetStatementTest:
			if !testLetStatement(t, i, j, stmt, actual.Statements[j]) {
				t.FailNow()
			}
		case ReturnStatementTest:
			if !testReturnStatement(t, i, j, stmt, actual.Statements[j]) {
				t.FailNow()
			}
		case ExpressionStatementTest:
			if !testExpressionStatement(t, i, j, stmt, actual.Statements[j]) {
				t.FailNow()
			}
		default:
			t.Fatalf("test[%d][%d] - unexpected type <%T>", i, j, stmt)
		}
	}
}

func testLetStatement(t *testing.T, i, j int, expected LetStatementTest, actual ast.Statement) bool {
	if "let" != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.LetStatement.TokenLiteral() ==> expected: <%s> but was: <%s>", i, j, "let", actual.TokenLiteral())
		return false
	}

	stmt, ok := actual.(*ast.LetStatement)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.LetStatement) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.LetStatement{}, actual)
		return false
	}

	if !testIdentifier(t, i, j, expected.Name, stmt.Name) {
		return false
	}

	if !testExpression(t, i, j, expected.Value, stmt.Value) {
		return false
	}

	return true
}

func testReturnStatement(t *testing.T, i, j int, expected ReturnStatementTest, actual ast.Statement) bool {
	if "return" != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.ReturnStatement.TokenLiteral() ==> expected: <%s> but was: <%s>", i, j, "return", actual.TokenLiteral())
		return false
	}

	stmt, ok := actual.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.ReturnStatement) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.ReturnStatement{}, actual)
		return false
	}

	if !testExpression(t, i, j, expected.ReturnValue, stmt.ReturnValue) {
		return false
	}

	return true
}

func testExpressionStatement(t *testing.T, i, j int, expected ExpressionStatementTest, actual ast.Statement) bool {
	fmt.Printf("%#v\n", actual)
	if expected.TokenLiteral != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.ExpressionStatement.TokenLiteral() ==> expected: <%s> but was: <%s>", i, j, expected.TokenLiteral, actual.TokenLiteral())
		return false
	}

	stmt, ok := actual.(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.ExpressionStatement) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.ExpressionStatement{}, actual)
		return false
	}

	if !testExpression(t, i, j, expected.Expression, stmt.Expression) {
		return false
	}

	return true
}

func testExpression(t *testing.T, i, j int, expected ExpressionTest, actual ast.Expression) bool {
	switch expected := expected.(type) {
	case PrefixExpressionTest:
		if !testPrefixExpression(t, i, j, expected, actual) {
			return false
		}
	case IdentifierTest:
		if !testIdentifier(t, i, j, expected, actual) {
			return false
		}
	case NumberTest:
		if !testNumber(t, i, j, expected, actual) {
			return false
		}
	default:
		t.Fatalf("test[%d][%d] - unexpected type <%T>", i, j, expected)
	}
	return true
}

func testPrefixExpression(t *testing.T, i, j int, expected PrefixExpressionTest, actual ast.Expression) bool {
	if expected.Operator != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.PrefixExpression.TokenLiteral() ==> expected: <%s> but was: <%s>", i, j, expected.Operator, actual.TokenLiteral())
		return false
	}

	expr, ok := actual.(*ast.PrefixExpression)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.PrefixExpression) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.PrefixExpression{}, actual)
		return false
	}

	if !testExpression(t, i, j, expected.Right, expr.Right) {
		return false
	}

	return true
}

func testIdentifier(t *testing.T, i, j int, expected IdentifierTest, actual ast.Expression) bool {
	if string(expected) != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.Identifier.TokenLiteral() ==> expected: <%s> but was: <%s>", i, j, string(expected), actual.TokenLiteral())
		return false
	}

	expr, ok := actual.(*ast.Identifier)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.Identifier) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.Identifier{}, actual)
		return false
	}

	if string(expected) != expr.Value {
		t.Errorf("test[%d][%d] - *ast.Identifier.Value ==> expected: <%s> but was: <%s>", i, j, string(expected), expr.Value)
		return false
	}

	return true
}

func testNumber(t *testing.T, i, j int, expected NumberTest, actual ast.Expression) bool {
	if fmt.Sprintf("%.0f", float64(expected)) != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.Number.TokenLiteral() ==> expected: <%f> but was: <%s>", i, j, float64(expected), actual.TokenLiteral())
		return false
	}

	expr, ok := actual.(*ast.Number)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.Number) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.Number{}, actual)
		return false
	}

	if float64(expected) != expr.Value {
		t.Errorf("test[%d][%d] - *ast.Number.Value ==> expected: <%f> but was: <%f>", i, j, float64(expected), expr.Value)
		return false
	}

	return true
}
