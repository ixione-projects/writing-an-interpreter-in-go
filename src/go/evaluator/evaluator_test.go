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

type (
	NumberTest  float64
	BooleanTest bool
)

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
