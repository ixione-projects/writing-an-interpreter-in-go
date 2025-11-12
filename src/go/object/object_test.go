package object

import "testing"

func TestNumberHashKey(t *testing.T) {
	expected := Number(1)
	actual := Number(1)

	if expected.HashKey() != actual.HashKey() {
		t.Errorf("Number.HashKey() ==> expected: <%d> but was: <%d>", expected.HashKey(), actual.HashKey())
	}
}

func TestBooleanHashKey(t *testing.T) {
	expected := Boolean(true)
	actual := Boolean(true)

	if expected.HashKey() != actual.HashKey() {
		t.Errorf("Boolean.HashKey() ==> expected: <%d> but was: <%d>", expected.HashKey(), actual.HashKey())
	}
}

func TestStringHashKey(t *testing.T) {
	expected := String("Hello World")
	actual := String("Hello World")

	if expected.HashKey() != actual.HashKey() {
		t.Errorf("String.HashKey() ==> expected: <%d> but was: <%d>", expected.HashKey(), actual.HashKey())
	}
}
