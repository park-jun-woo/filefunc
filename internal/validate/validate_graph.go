//ff:func feature=validate type=command control=sequence
//ff:what toulmin defeats graph — 전체 validate 룰과 예외 관계를 선언
package validate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// ValidateGraph declares all validation rules and their defeat relationships.
var ValidateGraph = toulmin.NewGraph("validate").
	// F rules (file structure)
	Warrant(RuleF1, nil, 1.0).
	Warrant(RuleF2, nil, 1.0).
	Warrant(RuleF3, nil, 1.0).
	Warrant(RuleF4, nil, 1.0).
	// Q rules (code quality)
	Warrant(RuleQ1, nil, 1.0).
	Warrant(RuleQ2Q3, nil, 1.0).
	// A rules (annotation)
	Warrant(RuleA1, nil, 1.0).
	Warrant(RuleA2, nil, 1.0).
	Warrant(RuleA3, nil, 1.0).
	Warrant(RuleA6, nil, 1.0).
	Warrant(RuleA7, nil, 1.0).
	Warrant(RuleA8, nil, 1.0).
	Warrant(RuleA9, nil, 1.0).
	Warrant(RuleA10, nil, 1.0).
	Warrant(RuleA11, nil, 1.0).
	Warrant(RuleA12, nil, 1.0).
	Warrant(RuleA13, nil, 1.0).
	Warrant(RuleA14, nil, 1.0).
	Warrant(RuleA15, nil, 1.0).
	Warrant(RuleA16, nil, 1.0).
	// Defeaters
	Defeater(DefeaterTestFile, nil, 1.0).
	Defeater(DefeaterConstOnly, nil, 1.0).
	Defeater(DefeaterNoFunc, nil, 1.0).
	// Defeat edges: test files defeat F/A rules
	Defeat(DefeaterTestFile, RuleF1).
	Defeat(DefeaterTestFile, RuleF2).
	Defeat(DefeaterTestFile, RuleF3).
	Defeat(DefeaterTestFile, RuleA1).
	Defeat(DefeaterTestFile, RuleA2).
	Defeat(DefeaterTestFile, RuleA3).
	Defeat(DefeaterTestFile, RuleA6).
	Defeat(DefeaterTestFile, RuleA7).
	Defeat(DefeaterTestFile, RuleA8).
	Defeat(DefeaterTestFile, RuleA9).
	Defeat(DefeaterTestFile, RuleA10).
	Defeat(DefeaterTestFile, RuleA11).
	Defeat(DefeaterTestFile, RuleA12).
	Defeat(DefeaterTestFile, RuleA13).
	Defeat(DefeaterTestFile, RuleA14).
	Defeat(DefeaterTestFile, RuleA15).
	Defeat(DefeaterTestFile, RuleA16).
	// Defeat edges: const-only files defeat F1
	Defeat(DefeaterConstOnly, RuleF1).
	// Defeat edges: no-func files defeat control/dimension/annotation rules
	Defeat(DefeaterNoFunc, RuleA9).
	Defeat(DefeaterNoFunc, RuleA10).
	Defeat(DefeaterNoFunc, RuleA11).
	Defeat(DefeaterNoFunc, RuleA12).
	Defeat(DefeaterNoFunc, RuleA13).
	Defeat(DefeaterNoFunc, RuleA14).
	Defeat(DefeaterNoFunc, RuleA15).
	Defeat(DefeaterNoFunc, RuleA16)
