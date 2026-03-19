//ff:func feature=validate type=command control=sequence
//ff:what toulmin defeats graph — 전체 validate 룰과 예외 관계를 선언
package validate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// ValidateGraph declares all validation rules and their defeat relationships.
var ValidateGraph = newValidateGraph()

func newValidateGraph() *toulmin.Graph {
	g := toulmin.NewGraph("validate")
	// F rules (file structure)
	wF1 := g.Warrant(RuleF1, nil, 1.0)
	wF2 := g.Warrant(RuleF2, nil, 1.0)
	wF3 := g.Warrant(RuleF3, nil, 1.0)
	_ = g.Warrant(RuleF4, nil, 1.0)
	// Q rules (code quality)
	_ = g.Warrant(RuleQ1, nil, 1.0)
	_ = g.Warrant(RuleQ2Q3, nil, 1.0)
	// A rules (annotation)
	wA1 := g.Warrant(RuleA1, nil, 1.0)
	wA2 := g.Warrant(RuleA2, nil, 1.0)
	wA3 := g.Warrant(RuleA3, nil, 1.0)
	wA6 := g.Warrant(RuleA6, nil, 1.0)
	wA7 := g.Warrant(RuleA7, nil, 1.0)
	wA8 := g.Warrant(RuleA8, nil, 1.0)
	wA9 := g.Warrant(RuleA9, nil, 1.0)
	wA10 := g.Warrant(RuleA10, nil, 1.0)
	wA11 := g.Warrant(RuleA11, nil, 1.0)
	wA12 := g.Warrant(RuleA12, nil, 1.0)
	wA13 := g.Warrant(RuleA13, nil, 1.0)
	wA14 := g.Warrant(RuleA14, nil, 1.0)
	wA15 := g.Warrant(RuleA15, nil, 1.0)
	wA16 := g.Warrant(RuleA16, nil, 1.0)
	// Defeaters
	dTestFile := g.Defeater(DefeaterTestFile, nil, 1.0)
	dConstOnly := g.Defeater(DefeaterConstOnly, nil, 1.0)
	dNoFunc := g.Defeater(DefeaterNoFunc, nil, 1.0)
	// Defeat edges: test files defeat F/A rules
	g.Defeat(dTestFile, wF1)
	g.Defeat(dTestFile, wF2)
	g.Defeat(dTestFile, wF3)
	g.Defeat(dTestFile, wA1)
	g.Defeat(dTestFile, wA2)
	g.Defeat(dTestFile, wA3)
	g.Defeat(dTestFile, wA6)
	g.Defeat(dTestFile, wA7)
	g.Defeat(dTestFile, wA8)
	g.Defeat(dTestFile, wA9)
	g.Defeat(dTestFile, wA10)
	g.Defeat(dTestFile, wA11)
	g.Defeat(dTestFile, wA12)
	g.Defeat(dTestFile, wA13)
	g.Defeat(dTestFile, wA14)
	g.Defeat(dTestFile, wA15)
	g.Defeat(dTestFile, wA16)
	// Defeat edges: const-only files defeat F1
	g.Defeat(dConstOnly, wF1)
	// Defeat edges: no-func files defeat control/dimension/annotation rules
	g.Defeat(dNoFunc, wA9)
	g.Defeat(dNoFunc, wA10)
	g.Defeat(dNoFunc, wA11)
	g.Defeat(dNoFunc, wA12)
	g.Defeat(dNoFunc, wA13)
	g.Defeat(dNoFunc, wA14)
	g.Defeat(dNoFunc, wA15)
	g.Defeat(dNoFunc, wA16)
	return g
}
