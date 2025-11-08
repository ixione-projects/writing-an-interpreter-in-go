package parser

import (
	"fmt"
	"testing"

	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/ast"
)

type ParserTest struct {
	input   string
	errors  []string
	debug   bool
	program ProgramTest
}

type NodeTest interface {
	node()
}

func (p ProgramTest) node()              {}
func (ls LetStatementTest) node()        {}
func (rs ReturnStatementTest) node()     {}
func (es ExpressionStatementTest) node() {}
func (bs BlockStatementTest) node()      {}
func (pe PrefixExpressionTest) node()    {}
func (ie InfixExpressionTest) node()     {}
func (ie IfExpressionTest) node()        {}
func (i IdentifierTest) node()           {}
func (nl NumberLiteralTest) node()       {}
func (bl BooleanLiteralTest) node()      {}

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
func (bs BlockStatementTest) statementNode()      {}

type LetStatementTest struct {
	Name   IdentifierTest
	Value  ExpressionTest
	String string
}

type ReturnStatementTest struct {
	ReturnValue ExpressionTest
	String      string
}

type ExpressionStatementTest struct {
	Expression ExpressionTest
	String     string
}

type BlockStatementTest struct {
	Statements []StatementTest
}

type ExpressionTest interface {
	NodeTest
	expressionNode()
}

func (pe PrefixExpressionTest) expressionNode() {}
func (ie InfixExpressionTest) expressionNode()  {}
func (ie IfExpressionTest) expressionNode()     {}
func (i IdentifierTest) expressionNode()        {}
func (nl NumberLiteralTest) expressionNode()    {}
func (nl BooleanLiteralTest) expressionNode()   {}

type PrefixExpressionTest struct {
	Operator string
	Right    ExpressionTest
}

type InfixExpressionTest struct {
	Left     ExpressionTest
	Operator string
	Right    ExpressionTest
}

type IfExpressionTest struct {
	Condition   ExpressionTest
	Consequence BlockStatementTest
	Alternative BlockStatementTest
}

type (
	IdentifierTest     string
	NumberLiteralTest  float64
	BooleanLiteralTest bool
)

func TestLetStatement(t *testing.T) {
	tests := []ParserTest{
		{
			input: `
			let x = 5;
			let y = 10;
			let foobar = 838383;`,
			program: ProgramTest{
				[]StatementTest{
					LetStatementTest{
						"x", NumberLiteralTest(5),
						"let x=5;",
					},
					LetStatementTest{
						"y", NumberLiteralTest(10),
						"let y=10;",
					},
					LetStatementTest{
						"foobar", NumberLiteralTest(838383),
						"let foobar=838383;",
					},
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
					ReturnStatementTest{
						NumberLiteralTest(5),
						"return 5;",
					},
					ReturnStatementTest{
						NumberLiteralTest(10),
						"return 10;",
					},
					ReturnStatementTest{
						NumberLiteralTest(993322),
						"return 993322;",
					},
				},
			},
		},
	}

	for i, test := range tests {
		testParser(t, i, test)
	}
}

func TestExpressionStatement(t *testing.T) {
	suites := []struct {
		name  string
		tests []ParserTest
	}{
		{
			name: "TestIdentifierExpression",
			tests: []ParserTest{
				{
					input: `foobar;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								IdentifierTest("foobar"),
								"foobar;",
							},
						},
					},
				},
			},
		},
		{
			name: "TestNumberLiteralExpression",
			tests: []ParserTest{
				{
					input: `5;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								NumberLiteralTest(5),
								"5;",
							},
						},
					},
				},
			},
		},
		{
			name: "TestPrefixExpression",
			tests: []ParserTest{
				{
					input: `!5;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								PrefixExpressionTest{"!", NumberLiteralTest(5)},
								"(!5);",
							},
						},
					},
				},
				{
					input: `-15;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								PrefixExpressionTest{"-", NumberLiteralTest(15)},
								"(-15);",
							},
						},
					},
				},
				{
					input: `!true;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								PrefixExpressionTest{"!", BooleanLiteralTest(true)},
								"(!true);",
							},
						},
					},
				},
				{
					input: `!false;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								PrefixExpressionTest{"!", BooleanLiteralTest(false)},
								"(!false);",
							},
						},
					},
				},
			},
		},
		{
			name: "TestInfixExpression",
			tests: []ParserTest{
				{
					input: `5 + 5;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{NumberLiteralTest(5), "+", NumberLiteralTest(5)},
								"(5+5);",
							},
						},
					},
				},
				{
					input: `5 - 5;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{NumberLiteralTest(5), "-", NumberLiteralTest(5)},
								"(5-5);",
							},
						},
					},
				},
				{
					input: `5 * 5;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{NumberLiteralTest(5), "*", NumberLiteralTest(5)},
								"(5*5);",
							},
						},
					},
				},
				{
					input: `5 / 5;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{NumberLiteralTest(5), "/", NumberLiteralTest(5)},
								"(5/5);",
							},
						},
					},
				},
				{
					input: `5 > 5;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{NumberLiteralTest(5), ">", NumberLiteralTest(5)},
								"(5>5);",
							},
						},
					},
				},
				{
					input: `5 < 5;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{NumberLiteralTest(5), "<", NumberLiteralTest(5)},
								"(5<5);",
							},
						},
					},
				},
				{
					input: `5 == 5;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{NumberLiteralTest(5), "==", NumberLiteralTest(5)},
								"(5==5);",
							},
						},
					},
				},
				{
					input: `5 != 5;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{NumberLiteralTest(5), "!=", NumberLiteralTest(5)},
								"(5!=5);",
							},
						},
					},
				},
				{
					input: `true == true`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{BooleanLiteralTest(true), "==", BooleanLiteralTest(true)},
								"(true==true);",
							},
						},
					},
				},
				{
					input: `true != false`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{BooleanLiteralTest(true), "!=", BooleanLiteralTest(false)},
								"(true!=false);",
							},
						},
					},
				},
				{
					input: `false == false`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{BooleanLiteralTest(false), "==", BooleanLiteralTest(false)},
								"(false==false);",
							},
						},
					},
				},
			},
		},
		{
			name: "TestOperatorPrecedence",
			tests: []ParserTest{
				{
					input: `-a * b;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{
									PrefixExpressionTest{"-", IdentifierTest("a")},
									"*",
									IdentifierTest("b"),
								},
								"((-a)*b);",
							},
						},
					},
				},
				{
					input: `!-a;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								PrefixExpressionTest{
									"!",
									PrefixExpressionTest{"-", IdentifierTest("a")},
								},
								"(!(-a));",
							},
						},
					},
				},
				{
					input: `a + b + c;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{
									InfixExpressionTest{IdentifierTest("a"), "+", IdentifierTest("b")},
									"+",
									IdentifierTest("c"),
								},
								"((a+b)+c);",
							},
						},
					},
				},
				{
					input: `a + b - c;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{
									InfixExpressionTest{IdentifierTest("a"), "+", IdentifierTest("b")},
									"-",
									IdentifierTest("c"),
								},
								"((a+b)-c);",
							},
						},
					},
				},
				{
					input: `a * b * c;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{
									InfixExpressionTest{IdentifierTest("a"), "*", IdentifierTest("b")},
									"*",
									IdentifierTest("c"),
								},
								"((a*b)*c);",
							},
						},
					},
				},
				{
					input: `a * b / c;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{
									InfixExpressionTest{IdentifierTest("a"), "*", IdentifierTest("b")},
									"/",
									IdentifierTest("c"),
								},
								"((a*b)/c);",
							},
						},
					},
				},
				{
					input: `a + b / c;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{
									IdentifierTest("a"),
									"+",
									InfixExpressionTest{IdentifierTest("b"), "/", IdentifierTest("c")},
								},
								"(a+(b/c));",
							},
						},
					},
				},
				{
					input: `a + b * c + d / e - f;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{
									InfixExpressionTest{
										InfixExpressionTest{
											IdentifierTest("a"),
											"+",
											InfixExpressionTest{IdentifierTest("b"), "*", IdentifierTest("c")},
										},
										"+",
										InfixExpressionTest{IdentifierTest("d"), "/", IdentifierTest("e")},
									},
									"-",
									IdentifierTest("f"),
								},
								"(((a+(b*c))+(d/e))-f);",
							},
						},
					},
				},
				{
					input: `3 + 4; -5 * 5`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{
									NumberLiteralTest(3),
									"+",
									NumberLiteralTest(4),
								},
								"(3+4);",
							},
							ExpressionStatementTest{
								InfixExpressionTest{
									PrefixExpressionTest{"-", NumberLiteralTest(5)},
									"*",
									NumberLiteralTest(5),
								},
								"((-5)*5);",
							},
						},
					},
				},
				{
					input: `5 > 4 == 3 < 4`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{
									InfixExpressionTest{NumberLiteralTest(5), ">", NumberLiteralTest(4)},
									"==",
									InfixExpressionTest{NumberLiteralTest(3), "<", NumberLiteralTest(4)},
								},
								"((5>4)==(3<4));",
							},
						},
					},
				},
				{
					input: `5 < 4 != 3 > 4`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{
									InfixExpressionTest{NumberLiteralTest(5), "<", NumberLiteralTest(4)},
									"!=",
									InfixExpressionTest{NumberLiteralTest(3), ">", NumberLiteralTest(4)},
								},
								"((5<4)!=(3>4));",
							},
						},
					},
				},
				{
					input: `3 + 4 * 5 == 3 * 1 + 4 * 5`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{
									InfixExpressionTest{
										NumberLiteralTest(3),
										"+",
										InfixExpressionTest{NumberLiteralTest(4), "*", NumberLiteralTest(5)},
									},
									"==",
									InfixExpressionTest{
										InfixExpressionTest{NumberLiteralTest(3), "*", NumberLiteralTest(1)},
										"+",
										InfixExpressionTest{NumberLiteralTest(4), "*", NumberLiteralTest(5)},
									},
								},
								"((3+(4*5))==((3*1)+(4*5)));",
							},
						},
					},
				},
				{
					input: `true`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BooleanLiteralTest(true),
								"true;",
							},
						},
					},
				},
				{
					input: `false`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BooleanLiteralTest(false),
								"false;",
							},
						},
					},
				},
				{
					input: `3 > 5 == false`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{
									InfixExpressionTest{NumberLiteralTest(3), ">", NumberLiteralTest(5)},
									"==",
									BooleanLiteralTest(false),
								},
								"((3>5)==false);",
							},
						},
					},
				},
				{
					input: `3 < 5 == true`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{
									InfixExpressionTest{NumberLiteralTest(3), "<", NumberLiteralTest(5)},
									"==",
									BooleanLiteralTest(true),
								},
								"((3<5)==true);",
							},
						},
					},
				},
				{
					input: `1 + (2 + 3) + 4`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{
									InfixExpressionTest{
										NumberLiteralTest(1),
										"+",
										InfixExpressionTest{NumberLiteralTest(2), "+", NumberLiteralTest(3)},
									},
									"+",
									NumberLiteralTest(4),
								},
								"((1+(2+3))+4);",
							},
						},
					},
				},
				{
					input: `(5 + 5) * 2`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{
									InfixExpressionTest{
										NumberLiteralTest(5),
										"+",
										NumberLiteralTest(5),
									},
									"*",
									NumberLiteralTest(2),
								},
								"((5+5)*2);",
							},
						},
					},
				},
				{
					input: `2 / (5 + 5)`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								InfixExpressionTest{
									NumberLiteralTest(2),
									"/",
									InfixExpressionTest{
										NumberLiteralTest(5),
										"+",
										NumberLiteralTest(5),
									},
								},
								"(2/(5+5));",
							},
						},
					},
				},
				{
					input: `-(5 + 5)`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								PrefixExpressionTest{
									"-",
									InfixExpressionTest{
										NumberLiteralTest(5),
										"+",
										NumberLiteralTest(5),
									},
								},
								"(-(5+5));",
							},
						},
					},
				},
				{
					input: `!(true == true)`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								PrefixExpressionTest{
									"!",
									InfixExpressionTest{
										BooleanLiteralTest(true),
										"==",
										BooleanLiteralTest(true),
									},
								},
								"(!(true==true));",
							},
						},
					},
				},
			},
		},
		{
			name: "TestBooleanLiteralExpression",
			tests: []ParserTest{
				{
					input: `true;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BooleanLiteralTest(true),
								"true;",
							},
						},
					},
				},
				{
					input: `false;`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								BooleanLiteralTest(false),
								"false;",
							},
						},
					},
				},
			},
		},
		{
			name: "TestIfExpression",
			tests: []ParserTest{
				{
					input: `if (x < y) { x }`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								IfExpressionTest{
									InfixExpressionTest{IdentifierTest("x"), "<", IdentifierTest("y")},
									BlockStatementTest{
										[]StatementTest{
											ExpressionStatementTest{
												IdentifierTest("x"),
												"x;",
											},
										},
									},
									BlockStatementTest{},
								},
								"if (x<y) {x;};",
							},
						},
					},
				},
				{
					input: `if (x < y) { x } else { y }`,
					program: ProgramTest{
						[]StatementTest{
							ExpressionStatementTest{
								IfExpressionTest{
									InfixExpressionTest{IdentifierTest("x"), "<", IdentifierTest("y")},
									BlockStatementTest{
										[]StatementTest{
											ExpressionStatementTest{
												IdentifierTest("x"),
												"x;",
											},
										},
									},
									BlockStatementTest{
										[]StatementTest{
											ExpressionStatementTest{
												IdentifierTest("y"),
												"y;",
											},
										},
									},
								},
								"if (x<y) {x;} else {y;};",
							},
						},
					},
				},
			},
		},
	}

	for _, suite := range suites {
		t.Run(suite.name, func(t *testing.T) {
			for i, test := range suite.tests {
				testParser(t, i, test)
			}
		})
	}
}

func testParser(t *testing.T, i int, test ParserTest) {
	p := New(test.input, test.debug)
	program := p.ParseProgram()

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
		if !testStatement(t, i, j, stmt, actual.Statements[j]) {
			t.Fatalf("test[%d][%d] - %s", i, j, actual.Statements[j].String())
		}
	}
}

func testStatement(t *testing.T, i, j int, expected StatementTest, actual ast.Statement) bool {
	switch expected := expected.(type) {
	case LetStatementTest:
		if !testLetStatement(t, i, j, expected, actual) {
			return false
		}
	case ReturnStatementTest:
		if !testReturnStatement(t, i, j, expected, actual) {
			return false
		}
	case ExpressionStatementTest:
		if !testExpressionStatement(t, i, j, expected, actual) {
			return false
		}
	case BlockStatementTest:
		if !testBlockStatement(t, i, j, expected, actual) {
			return false
		}
	default:
		t.Fatalf("test[%d][%d] - unexpected type <%T>", i, j, expected)
	}
	return true
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

	if expected.String != actual.String() {
		t.Errorf("test[%d][%d] - *ast.LetStatement.String() ==> expected: <%s> but was: <%s>", i, j, expected.String, actual.String())
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

	if expected.String != actual.String() {
		t.Errorf("test[%d][%d] - *ast.ReturnStatement.String() ==> expected: <%s> but was: <%s>", i, j, expected.String, actual.String())
		return false
	}

	return true
}

func testExpressionStatement(t *testing.T, i, j int, expected ExpressionStatementTest, actual ast.Statement) bool {
	stmt, ok := actual.(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.ExpressionStatement) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.ExpressionStatement{}, actual)
		return false
	}

	if !testExpression(t, i, j, expected.Expression, stmt.Expression) {
		return false
	}

	if expected.String != actual.String() {
		t.Errorf("test[%d][%d] - *ast.ExpressionStatement.String() ==> expected: <%s> but was: <%s>", i, j, expected.String, actual.String())
		return false
	}

	return true
}

func testBlockStatement(t *testing.T, i, j int, expected BlockStatementTest, actual ast.Statement) bool {
	block, ok := actual.(*ast.BlockStatement)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.BlockStatement) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.BlockStatement{}, actual)
		return false
	}

	if len(expected.Statements) != len(block.Statements) {
		t.Fatalf("test[%d] - len(block.Statements) ==> expected: <%d> but was: <%d>", i, len(expected.Statements), len(block.Statements))
	}

	for k, stmt := range expected.Statements {
		if !testStatement(t, i, j, stmt, block.Statements[k]) {
			return false
		}
	}

	return true
}

func testExpression(t *testing.T, i, j int, expected ExpressionTest, actual ast.Expression) bool {
	switch expected := expected.(type) {
	case PrefixExpressionTest:
		if !testPrefixExpression(t, i, j, expected, actual) {
			return false
		}
	case InfixExpressionTest:
		if !testInfixExpression(t, i, j, expected, actual) {
			return false
		}
	case IfExpressionTest:
		if !testIfExpression(t, i, j, expected, actual) {
			return false
		}
	case IdentifierTest:
		if !testIdentifier(t, i, j, expected, actual) {
			return false
		}
	case NumberLiteralTest:
		if !testNumberLiteral(t, i, j, expected, actual) {
			return false
		}
	case BooleanLiteralTest:
		if !testBooleanLiteral(t, i, j, expected, actual) {
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

func testInfixExpression(t *testing.T, i, j int, expected InfixExpressionTest, actual ast.Expression) bool {
	if expected.Operator != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.InfixExpression.TokenLiteral() ==> expected: <%s> but was: <%s>", i, j, expected.Operator, actual.TokenLiteral())
		return false
	}

	expr, ok := actual.(*ast.InfixExpression)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.InfixExpression) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.InfixExpression{}, actual)
		return false
	}

	if !testExpression(t, i, j, expected.Left, expr.Left) {
		return false
	}

	if !testExpression(t, i, j, expected.Right, expr.Right) {
		return false
	}

	return true
}

func testIfExpression(t *testing.T, i, j int, expected IfExpressionTest, actual ast.Expression) bool {
	if "if" != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.IfExpression.TokenLiteral() ==> expected: <%s> but was: <%s>", i, j, "if", actual.TokenLiteral())
		return false
	}

	expr, ok := actual.(*ast.IfExpression)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.IfExpression) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.IfExpression{}, actual)
		return false
	}

	if !testExpression(t, i, j, expected.Condition, expr.Condition) {
		return false
	}

	if !testBlockStatement(t, i, j, expected.Consequence, expr.Consequence) {
		return false
	}

	if len(expected.Alternative.Statements) > 0 && !testBlockStatement(t, i, j, expected.Alternative, expr.Alternative) {
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

func testNumberLiteral(t *testing.T, i, j int, expected NumberLiteralTest, actual ast.Expression) bool {
	if fmt.Sprintf("%.0f", float64(expected)) != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.NumberLiteral.TokenLiteral() ==> expected: <%f> but was: <%s>", i, j, float64(expected), actual.TokenLiteral())
		return false
	}

	expr, ok := actual.(*ast.NumberLiteral)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.NumberLiteral) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.NumberLiteral{}, actual)
		return false
	}

	if float64(expected) != expr.Value {
		t.Errorf("test[%d][%d] - *ast.NumberLiteral.Value ==> expected: <%f> but was: <%f>", i, j, float64(expected), expr.Value)
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, i, j int, expected BooleanLiteralTest, actual ast.Expression) bool {
	if fmt.Sprintf("%t", bool(expected)) != actual.TokenLiteral() {
		t.Errorf("test[%d][%d] - *ast.BooleanLiteral.TokenLiteral() ==> expected: <%t> but was: <%s>", i, j, bool(expected), actual.TokenLiteral())
		return false
	}

	expr, ok := actual.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("test[%d][%d] - actual.(*ast.BooleanLiteral) ==> unexpected type, expected: <%T> but was: <%T>", i, j, &ast.NumberLiteral{}, actual)
		return false
	}

	if bool(expected) != expr.Value {
		t.Errorf("test[%d][%d] - *ast.BooleanLiteral.Value ==> expected: <%t> but was: <%t>", i, j, bool(expected), expr.Value)
		return false
	}

	return true
}
