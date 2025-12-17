package assert

import (
	"fmt"
	"sync"
	"testing"
)

type Reporter interface {
	Report(error)
}

type (
	MainReporter struct{}
	TestReporter struct {
		harness testing.TB
	}
)

func (mr *MainReporter) Report(err error) {
	panic(err)
}

func (tr *TestReporter) Report(err error) {
	tr.harness.Helper()
	tr.harness.Error(err)
}

var main_reporter *MainReporter
var once sync.Once

func GetMainReporter() *MainReporter {
	once.Do(func() {
		main_reporter = &MainReporter{}
	})
	return main_reporter
}

func GetTestReporter(harness testing.TB) *TestReporter {
	return &TestReporter{
		harness: harness,
	}
}

func Equals(expected, actual any, r Reporter) bool {
	return EqualsWithMessage(expected, actual, "", r)
}

func EqualsWithMessage(expected, actual any, msg string, r Reporter) bool {
	if tr, ok := r.(*TestReporter); ok {
		tr.harness.Helper()
	}

	if expected != actual {
		r.Report(buildEqualsError(expected, actual, msg))
		return false
	}

	return true
}

func buildEqualsError(expected, actual any, msg string) error {
	reason := fmt.Sprintf("expected: <%#v> but was: <%#v>", expected, actual)
	if msg != "" {
		return fmt.Errorf("%s ==> %s", msg, reason)
	}
	return fmt.Errorf("%s", reason)
}
