package evaluator

import (
	"testing"

	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/object"
	"github.com/ixione-projects/writing-an-interpreter-in-go/src/go/parser"
)

type EvaluatorTest struct {
	input  string
	object ObjectTest
}

type ObjectTest interface {
	object()
}

func (n NumberTest) object()  {}
func (b BooleanTest) object() {}
func (n NullTest) object()    {}
func (e ErrorTest) object()   {}

type (
	NumberTest  float64
	BooleanTest bool
)

type NullTest struct{}

type ErrorTest struct {
	Message string
}

func TestEvaluateNumber(t *testing.T) {
	suites := []struct {
		name  string
		tests []EvaluatorTest
	}{
		{
			name: "TestEvaluateNumber",
			tests: []EvaluatorTest{
				{
					input:  "5",
					object: NumberTest(5),
				},
				{
					input:  "10",
					object: NumberTest(10),
				},
				{
					input:  "-5",
					object: NumberTest(-5),
				},
				{
					input:  "-10",
					object: NumberTest(-10),
				},
				{
					input:  "5 + 5 + 5 + 5 - 10",
					object: NumberTest(10),
				},
				{
					input:  "2 * 2 * 2 * 2 * 2",
					object: NumberTest(32),
				},
				{
					input:  "-50 + 100 + -50",
					object: NumberTest(0),
				},
				{
					input:  "5 * 2 + 10",
					object: NumberTest(20),
				},
				{
					input:  "5 + 2 * 10",
					object: NumberTest(25),
				},
				{
					input:  "20 + 2 * -10",
					object: NumberTest(0),
				},
				{
					input:  "50 / 2 * 2 + 10",
					object: NumberTest(60),
				},
				{
					input:  "2 * (5 + 10)",
					object: NumberTest(30),
				},
				{
					input:  "3 * 3 * 3 + 10",
					object: NumberTest(37),
				},
				{
					input:  "3 * (3 * 3) + 10",
					object: NumberTest(37),
				},
				{
					input:  "(5 + 10 * 2 + 15 / 3) * 2 + -10",
					object: NumberTest(50),
				},
				{
					input:  "if (true) { 10 }",
					object: NumberTest(10),
				},
				{
					input:  "if (1) { 10 }",
					object: NumberTest(10),
				},
				{
					input:  "if (1 < 2) { 10 }",
					object: NumberTest(10),
				},
				{
					input:  "if (1 > 2) { 10 } else { 20 }",
					object: NumberTest(20),
				},
				{
					input:  "if (1 < 2) { 10 } else { 20 }",
					object: NumberTest(10),
				},
			},
		},
		{
			name: "TestEvaluateBoolean",
			tests: []EvaluatorTest{
				{
					input:  "true",
					object: BooleanTest(true),
				},
				{
					input:  "false",
					object: BooleanTest(false),
				},
				{
					input:  "!true",
					object: BooleanTest(false),
				},
				{
					input:  "!false",
					object: BooleanTest(true),
				},
				{
					input:  "!5",
					object: BooleanTest(false),
				},
				{
					input:  "!!true",
					object: BooleanTest(true),
				},
				{
					input:  "!!false",
					object: BooleanTest(false),
				},
				{
					input:  "!!5",
					object: BooleanTest(true),
				},
				{
					input:  "1 < 2",
					object: BooleanTest(true),
				},
				{
					input:  "1 > 2",
					object: BooleanTest(false),
				},
				{
					input:  "1 < 1",
					object: BooleanTest(false),
				},
				{
					input:  "1 > 1",
					object: BooleanTest(false),
				},
				{
					input:  "1 == 1",
					object: BooleanTest(true),
				},
				{
					input:  "1 != 1",
					object: BooleanTest(false),
				},
				{
					input:  "1 == 2",
					object: BooleanTest(false),
				},
				{
					input:  "1 != 2",
					object: BooleanTest(true),
				},
				{
					input:  "true == true",
					object: BooleanTest(true),
				},
				{
					input:  "false == false",
					object: BooleanTest(true),
				},
				{
					input:  "true == false",
					object: BooleanTest(false),
				},
				{
					input:  "true != false",
					object: BooleanTest(true),
				},
				{
					input:  "false != true",
					object: BooleanTest(true),
				},
				{
					input:  "(1 < 2) == true",
					object: BooleanTest(true),
				},
				{
					input:  "(1 < 2) == false",
					object: BooleanTest(false),
				},
				{
					input:  "(1 > 2) == true",
					object: BooleanTest(false),
				},
				{
					input:  "(1 > 2) == false",
					object: BooleanTest(true),
				},
			},
		},
		{
			name: "TestEvaluateNull",
			tests: []EvaluatorTest{
				{
					input:  "if (false) { 10 }",
					object: NullTest{},
				},
				{
					input:  "if (1 > 2) { 10 }",
					object: NullTest{},
				},
			},
		},
		{
			name: "TestEvaluateReturnValue",
			tests: []EvaluatorTest{
				{
					input:  "return 10;",
					object: NumberTest(10),
				},
				{
					input:  "return 10; 9;",
					object: NumberTest(10),
				},
				{
					input:  "return 2 * 5; 9;",
					object: NumberTest(10),
				},
				{
					input:  "9; return 2 * 5; 9;",
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
			},
		},
		{
			name: "TestEvaluateError",
			tests: []EvaluatorTest{
				{
					input:  "5 + true;",
					object: ErrorTest{"type mismatch: INTEGER + BOOLEAN"},
				},
				{
					input:  "5 + true; 5;",
					object: ErrorTest{"type mismatch: INTEGER + BOOLEAN"},
				},
				{
					input:  "-true",
					object: ErrorTest{"unknown operator: -BOOLEAN"},
				},
				{
					input:  "true + false;",
					object: ErrorTest{"unknown operator: BOOLEAN + BOOLEAN"},
				},
				{
					input:  "5; true + false; 5",
					object: ErrorTest{"unknown operator: BOOLEAN + BOOLEAN"},
				},
				{
					input:  "if (10 > 1) { true + false; }",
					object: ErrorTest{"unknown operator: BOOLEAN + BOOLEAN"},
				},
				{
					input: `
					if (10 > 1) {
						if (10 > 1) {
							return true + false;
						}

						return 1;
					}`,
					object: ErrorTest{"unknown operator: BOOLEAN + BOOLEAN"},
				},
			},
		},
	}

	for _, suite := range suites {
		t.Run(suite.name, func(t *testing.T) {
			for i, test := range suite.tests {
				testEvaluator(t, i, test)
			}
		})
	}
}

func testEvaluator(t *testing.T, i int, test EvaluatorTest) {
	p := parser.New(test.input, false)
	program := p.ParseProgram()

	if 0 != len(p.Errors()) {
		t.Errorf("test[%d] - len(p.Errors()) ==> expected: <%d> but was: <%d>", i, 0, len(p.Errors()))
		for j, msg := range p.Errors() {
			t.Errorf("--------- p.Errors()[%d]: %s", j, msg)
		}
		t.FailNow()
	}

	if !testObject(t, i, test.object, Evaluate(program)) {
		t.Fatalf("test[%d] - %s", i, program.String())
	}
}

func testObject(t *testing.T, i int, expected ObjectTest, actual object.Object) bool {
	switch expected := expected.(type) {
	case NumberTest:
		if !testNumber(t, i, expected, actual) {
			return false
		}
	case BooleanTest:
		if !testBoolean(t, i, expected, actual) {
			return false
		}
	case NullTest:
		if !testNull(t, i, expected, actual) {
			return false
		}
	case ErrorTest:
		if !testError(t, i, expected, actual) {
			return false
		}
	default:
		t.Fatalf("test[%d] - unexpected type <%T>", i, expected)
	}
	return true
}

func testNumber(t *testing.T, i int, expected NumberTest, actual object.Object) bool {
	value, ok := actual.(object.Number)
	if !ok {
		t.Errorf("test[%d] - actual.(object.Number) ==> unexpected type, expected: <%T> but was: <%T>", i, object.Number(0.0), actual)
		return false
	}

	if float64(expected) != float64(value) {
		t.Errorf("test[%d] - object.Number ==> expected: <%f> but was: <%f>", i, float64(expected), float64(value))
		return false
	}

	return true
}

func testBoolean(t *testing.T, i int, expected BooleanTest, actual object.Object) bool {
	value, ok := actual.(object.Boolean)
	if !ok {
		t.Errorf("test[%d] - actual.(object.Boolean) ==> unexpected type, expected: <%T> but was: <%T>", i, object.Boolean(false), actual)
		return false
	}

	if bool(expected) != bool(value) {
		t.Errorf("test[%d] - object.Boolean ==> expected: <%t> but was: <%t>", i, bool(expected), bool(value))
		return false
	}

	return true
}

func testNull(t *testing.T, i int, expected NullTest, actual object.Object) bool {
	_, ok := actual.(*object.Null)
	if !ok {
		t.Errorf("test[%d] - actual.(*object.Null) ==> unexpected type, expected: <%T> but was: <%T>", i, &object.Null{}, actual)
		return false
	}

	return true
}

func testError(t *testing.T, i int, expected ErrorTest, actual object.Object) bool {
	value, ok := actual.(*object.Error)
	if !ok {
		t.Errorf("test[%d] - actual.(*object.Error) ==> unexpected type, expected: <%T> but was: <%T>", i, &object.Error{}, actual)
		return false
	}

	if expected.Message != value.Message {
		t.Errorf("test[%d] - *object.Error.Message ==> expected: <%s> but was: <%s>", i, expected.Message, value.Message)
		return false
	}

	return true
}
