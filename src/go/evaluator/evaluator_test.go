package evaluator

import (
	"testing"

	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/object"
	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/parser"
)

type EvaluatorTest struct {
	input  string
	object ObjectTest
	error  ErrorTest
}

type ObjectTest interface {
	object()
}

func (f FunctionTest) object() {}
func (n NumberTest) object()   {}
func (b BooleanTest) object()  {}
func (s StringTest) object()   {}
func (n NullTest) object()     {}

type FunctionTest struct {
	Inspect string
}

type (
	NumberTest  float64
	BooleanTest bool
	StringTest  string
)

type NullTest struct{}

type ErrorTest struct {
	Message string
}

var suites = []struct {
	name  string
	tests []EvaluatorTest
}{
	{
		name: "TestNumber",
		tests: []EvaluatorTest{
			{
				input:  `5`,
				object: NumberTest(5),
			},
			{
				input:  `10`,
				object: NumberTest(10),
			},
			{
				input:  `-5`,
				object: NumberTest(-5),
			},
			{
				input:  `-10`,
				object: NumberTest(-10),
			},
			{
				input:  `5 + 5 + 5 + 5 - 10`,
				object: NumberTest(10),
			},
			{
				input:  `2 * 2 * 2 * 2 * 2`,
				object: NumberTest(32),
			},
			{
				input:  `-50 + 100 + -50`,
				object: NumberTest(0),
			},
			{
				input:  `5 * 2 + 10`,
				object: NumberTest(20),
			},
			{
				input:  `5 + 2 * 10`,
				object: NumberTest(25),
			},
			{
				input:  `20 + 2 * -10`,
				object: NumberTest(0),
			},
			{
				input:  `50 / 2 * 2 + 10`,
				object: NumberTest(60),
			},
			{
				input:  `2 * (5 + 10)`,
				object: NumberTest(30),
			},
			{
				input:  `3 * 3 * 3 + 10`,
				object: NumberTest(37),
			},
			{
				input:  `3 * (3 * 3) + 10`,
				object: NumberTest(37),
			},
			{
				input:  `(5 + 10 * 2 + 15 / 3) * 2 + -10`,
				object: NumberTest(50),
			},
			{
				input:  `if (true) { 10 }`,
				object: NumberTest(10),
			},
			{
				input:  `if (1) { 10 }`,
				object: NumberTest(10),
			},
			{
				input:  `if (1 < 2) { 10 }`,
				object: NumberTest(10),
			},
			{
				input:  `if (1 > 2) { 10 } else { 20 }`,
				object: NumberTest(20),
			},
			{
				input:  `if (1 < 2) { 10 } else { 20 }`,
				object: NumberTest(10),
			},
			{
				input:  `return 10;`,
				object: NumberTest(10),
			},
			{
				input:  `return 10; 9;`,
				object: NumberTest(10),
			},
			{
				input:  `return 2 * 5; 9;`,
				object: NumberTest(10),
			},
			{
				input:  `9; return 2 * 5; 9;`,
				object: NumberTest(10),
			},
			{
				input: `
					if (10 > 1) {
						if (10 > 1) {
							return 10;
						}

						return 1;
					}`,
				object: NumberTest(10),
			},
			{
				input:  `let a = 5; a;`,
				object: NumberTest(5),
			},
			{
				input:  `let a = 5 * 5; a;`,
				object: NumberTest(25),
			},
			{
				input:  `let a = 5; let b = a; b;`,
				object: NumberTest(5),
			},
			{
				input:  `let a = 5; let b = a; let c = a + b + 5; c;`,
				object: NumberTest(15),
			},
			{
				input:  `let identity = fn(x) { x; }; identity(5);`,
				object: NumberTest(5),
			},
			{
				input:  `let identity = fn(x) { return x; }; identity(5);`,
				object: NumberTest(5),
			},
			{
				input:  `let double = fn(x) { x * 2; }; double(5);`,
				object: NumberTest(10),
			},
			{
				input:  `let add = fn(x, y) { x + y; }; add(5, 5);`,
				object: NumberTest(10),
			},
			{
				input:  `let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));`,
				object: NumberTest(20),
			},
			{
				input:  `fn(x) { x; }(5)`,
				object: NumberTest(5),
			},
			{
				input: `
					let f = fn(x) {
						return x;
						x + 10;
					};
					f(10);`,
				object: NumberTest(10),
			},
			{
				input: `
					let f = fn(x) {
						let result = x + 10;
						return result;
						return 10;
					};
					f(10);`,
				object: NumberTest(20),
			},
			{
				input: `
					let newAdder = fn(x) {
						fn(y) { x + y };
					};
					let addTwo = newAdder(2);
					addTwo(2);`,
				object: NumberTest(4),
			},
		},
	},
	{
		name: "TestBoolean",
		tests: []EvaluatorTest{
			{
				input:  `true`,
				object: BooleanTest(true),
			},
			{
				input:  `false`,
				object: BooleanTest(false),
			},
			{
				input:  `!true`,
				object: BooleanTest(false),
			},
			{
				input:  `!false`,
				object: BooleanTest(true),
			},
			{
				input:  `!5`,
				object: BooleanTest(false),
			},
			{
				input:  `!!true`,
				object: BooleanTest(true),
			},
			{
				input:  `!!false`,
				object: BooleanTest(false),
			},
			{
				input:  `!!5`,
				object: BooleanTest(true),
			},
			{
				input:  `1 < 2`,
				object: BooleanTest(true),
			},
			{
				input:  `1 > 2`,
				object: BooleanTest(false),
			},
			{
				input:  `1 < 1`,
				object: BooleanTest(false),
			},
			{
				input:  `1 > 1`,
				object: BooleanTest(false),
			},
			{
				input:  `1 == 1`,
				object: BooleanTest(true),
			},
			{
				input:  `1 != 1`,
				object: BooleanTest(false),
			},
			{
				input:  `1 == 2`,
				object: BooleanTest(false),
			},
			{
				input:  `1 != 2`,
				object: BooleanTest(true),
			},
			{
				input:  `true == true`,
				object: BooleanTest(true),
			},
			{
				input:  `false == false`,
				object: BooleanTest(true),
			},
			{
				input:  `true == false`,
				object: BooleanTest(false),
			},
			{
				input:  `true != false`,
				object: BooleanTest(true),
			},
			{
				input:  `false != true`,
				object: BooleanTest(true),
			},
			{
				input:  `(1 < 2) == true`,
				object: BooleanTest(true),
			},
			{
				input:  `(1 < 2) == false`,
				object: BooleanTest(false),
			},
			{
				input:  `(1 > 2) == true`,
				object: BooleanTest(false),
			},
			{
				input:  `(1 > 2) == false`,
				object: BooleanTest(true),
			},
		},
	},
	{
		name: "TestNull",
		tests: []EvaluatorTest{
			{
				input:  `if (false) { 10 }`,
				object: NullTest{},
			},
			{
				input:  `if (1 > 2) { 10 }`,
				object: NullTest{},
			},
		},
	},
	{
		name: "TestError",
		tests: []EvaluatorTest{
			{
				input: `5 + true;`,
				error: ErrorTest{"type mismatch: INTEGER + BOOLEAN"},
			},
			{
				input: `5 + true; 5;`,
				error: ErrorTest{"type mismatch: INTEGER + BOOLEAN"},
			},
			{
				input: `-true`,
				error: ErrorTest{"unknown operator: -BOOLEAN"},
			},
			{
				input: `true + false;`,
				error: ErrorTest{"unknown operator: BOOLEAN + BOOLEAN"},
			},
			{
				input: `5; true + false; 5`,
				error: ErrorTest{"unknown operator: BOOLEAN + BOOLEAN"},
			},
			{
				input: `if (10 > 1) { true + false; }`,
				error: ErrorTest{"unknown operator: BOOLEAN + BOOLEAN"},
			},
			{
				input: `
					if (10 > 1) {
						if (10 > 1) {
							return true + false;
						}

						return 1;
					}`,
				error: ErrorTest{"unknown operator: BOOLEAN + BOOLEAN"},
			},
			{
				input: `foobar`,
				error: ErrorTest{"identifier not found: foobar"},
			},
			{
				input: `"Hello" - "World"`,
				error: ErrorTest{"unknown operator: STRING - STRING"},
			},
		},
	},
	{
		name: "TestFunction",
		tests: []EvaluatorTest{
			{
				input:  `fn(x) { x + 2; };`,
				object: FunctionTest{"<fn (x)>"},
			},
		},
	},
	{
		name: "TestString",
		tests: []EvaluatorTest{
			{
				input:  `"Hello World!"`,
				object: StringTest("Hello World!"),
			},
			{
				input:  `"Hello" + " " + "World!"`,
				object: StringTest("Hello World!"),
			},
		},
	},
}

func TestEvaluateNumber(t *testing.T) {
	for _, suite := range suites {
		t.Run(suite.name, func(t *testing.T) {
			for i, test := range suite.tests {
				testEvaluator(t, i, test)
			}
		})
	}
}

func testEvaluator(tb testing.TB, i int, test EvaluatorTest) {
	p := parser.NewParser(test.input, false)
	program := p.ParseProgram()

	if 0 != len(p.Errors()) {
		tb.Errorf("test[%d] - len(p.Errors()) ==> expected: <%d> but was: <%d>", i, 0, len(p.Errors()))
		for j, msg := range p.Errors() {
			tb.Errorf("--------- p.Errors()[%d]: %s", j, msg)
		}
		tb.FailNow()
	}

	object, interruption := Evaluate(program, object.NewEnvironment(nil))
	if test.error.Message == "" {
		if !testObject(tb, i, test.object, object) {
			tb.Fatalf("test[%d] - %s", i, program.String())
		}
		return
	}

	if !testError(tb, i, test.error, interruption) {
		tb.Fatalf("test[%d] - %s", i, program.String())
	}
}

func testObject(tb testing.TB, i int, expected ObjectTest, actual object.Object) bool {
	switch expected := expected.(type) {
	case FunctionTest:
		if !testFunction(tb, i, expected, actual) {
			return false
		}
	case NumberTest:
		if !testNumber(tb, i, expected, actual) {
			return false
		}
	case BooleanTest:
		if !testBoolean(tb, i, expected, actual) {
			return false
		}
	case StringTest:
		if !testString(tb, i, expected, actual) {
			return false
		}
	case NullTest:
		if !testNull(tb, i, expected, actual) {
			return false
		}
	default:
		tb.Fatalf("test[%d] - unexpected type <%T>", i, expected)
	}
	return true
}

func testFunction(tb testing.TB, i int, expected FunctionTest, actual object.Object) bool {
	function, ok := actual.(*object.Function)
	if !ok {
		tb.Errorf("test[%d] - actual.(*object.Function) ==> unexpected type, expected: <%T> but was: <%T>", i, &object.Function{}, actual)
		return false
	}

	if expected.Inspect != function.Inspect() {
		tb.Errorf("test[%d] - function.Inspect() ==> expected: <%s> but was: <%s>", i, expected.Inspect, function.Inspect())
		return false
	}

	return true
}

func testNumber(tb testing.TB, i int, expected NumberTest, actual object.Object) bool {
	value, ok := actual.(object.Number)
	if !ok {
		tb.Errorf("test[%d] - actual.(object.Number) ==> unexpected type, expected: <%T> but was: <%T>", i, object.Number(0.0), actual)
		return false
	}

	if float64(expected) != float64(value) {
		tb.Errorf("test[%d] - object.Number ==> expected: <%f> but was: <%f>", i, float64(expected), float64(value))
		return false
	}

	return true
}

func testBoolean(tb testing.TB, i int, expected BooleanTest, actual object.Object) bool {
	value, ok := actual.(object.Boolean)
	if !ok {
		tb.Errorf("test[%d] - actual.(object.Boolean) ==> unexpected type, expected: <%T> but was: <%T>", i, object.Boolean(false), actual)
		return false
	}

	if bool(expected) != bool(value) {
		tb.Errorf("test[%d] - object.Boolean ==> expected: <%t> but was: <%t>", i, bool(expected), bool(value))
		return false
	}

	return true
}

func testString(tb testing.TB, i int, expected StringTest, actual object.Object) bool {
	value, ok := actual.(object.String)
	if !ok {
		tb.Errorf("test[%d] - actual.(object.String) ==> unexpected type, expected: <%T> but was: <%T>", i, object.String(""), actual)
		return false
	}

	if string(expected) != string(value) {
		tb.Errorf("test[%d] - object.String ==> expected: <%s> but was: <%s>", i, string(expected), string(value))
		return false
	}

	return true
}

func testNull(tb testing.TB, i int, expected NullTest, actual object.Object) bool {
	_, ok := actual.(*object.Null)
	if !ok {
		tb.Errorf("test[%d] - actual.(*object.Null) ==> unexpected type, expected: <%T> but was: <%T>", i, &object.Null{}, actual)
		return false
	}

	return true
}

func testError(tb testing.TB, i int, expected ErrorTest, actual object.Interruption) bool {
	value, ok := actual.(*object.Error)
	if !ok {
		tb.Errorf("test[%d] - actual.(*object.Error) ==> unexpected type, expected: <%T> but was: <%T>", i, &object.Error{}, actual)
		return false
	}

	if expected.Message != value.Message {
		tb.Errorf("test[%d] - *object.Error.Message ==> expected: <%s> but was: <%s>", i, expected.Message, value.Message)
		return false
	}

	return true
}
