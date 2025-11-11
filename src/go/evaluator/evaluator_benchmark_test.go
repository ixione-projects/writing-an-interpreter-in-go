package evaluator

import "testing"

func BenchmarkEvaluate(b *testing.B) {
	for _, suite := range suites {
		b.Run(suite.name, func(b *testing.B) {
			for b.Loop() {
				for i, test := range suite.tests {
					testEvaluator(b, i, test)
				}
			}
		})
	}
}
