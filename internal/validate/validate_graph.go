//ff:func feature=validate type=command control=iteration dimension=1
//ff:what toulmin defeats graph — 전체 validate 룰과 예외 관계를 선언
package validate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// ValidateGraph declares all validation rules and their defeat relationships.
var ValidateGraph = newValidateGraph()

func newValidateGraph() *toulmin.Graph {
	g := toulmin.NewGraph("validate")

	// ── 파일 구조: 파일당 하나 ──
	wF1 := g.Warrant(CountMax, &CountMaxBacking{Field: "Funcs", Max: 1, Rule: "F1",
		Message: "file contains multiple funcs; expected 1 file 1 func"}, 1.0)
	_ = g.Warrant(CountMax, &CountMaxBacking{Field: "Types", Max: 1, Rule: "F2",
		Message: "file contains multiple types; expected 1 file 1 type"}, 1.0)
	_ = g.Warrant(CountMax, &CountMaxBacking{Field: "Methods", Max: 1, Rule: "F3",
		Message: "file contains multiple methods; expected 1 file 1 method"}, 1.0)

	// ── 파일 구조: init 단독 불허 ──
	_ = g.Warrant(ExistsWhen, &ExistsWhenBacking{When: "HasInit", Need: "companion", Rule: "F4",
		Level: "ERROR", Message: "init() must not exist alone; requires accompanying var or func"}, 1.0)

	// ── 코드 품질 ──
	_ = g.Warrant(CheckDepthLimit, nil, 1.0)
	_ = g.Warrant(CheckFuncLines, nil, 1.0)

	// ── 어노테이션: 존재 필수 ──
	_ = g.Warrant(ExistsWhen, &ExistsWhenBacking{When: "HasFuncs", Need: "ff:func", Rule: "A1",
		Level: "ERROR", Message: "file with func must have //ff:func annotation"}, 1.0)
	_ = g.Warrant(ExistsWhen, &ExistsWhenBacking{When: "HasTypes", Need: "ff:type", Rule: "A1",
		Level: "ERROR", Message: "file with type must have //ff:type annotation"}, 1.0)
	_ = g.Warrant(ExistsWhen, &ExistsWhenBacking{When: "HasFuncOrType", Need: "ff:what", Rule: "A3",
		Level: "ERROR", Message: "file with func or type must have //ff:what annotation"}, 1.0)
	wA9 := g.Warrant(ExistsWhen, &ExistsWhenBacking{When: "HasFuncs", Need: "control", Rule: "A9",
		Level: "ERROR", Message: "func file must have control= annotation (sequence, selection, or iteration)"}, 1.0)
	wA15 := g.Warrant(ExistsWhen, &ExistsWhenBacking{When: "ControlIteration", Need: "dimension", Rule: "A15",
		Level: "ERROR", Message: "control=iteration requires dimension= annotation"}, 1.0)

	// ── 어노테이션: 제어 구조 일치 ──
	wA10 := g.Warrant(ControlMatch, &ControlMatchBacking{Control: "selection", MustHave: "switch", Rule: "A10",
		Message: "control=selection but no switch found at depth 1"}, 1.0)
	wA11 := g.Warrant(ControlMatch, &ControlMatchBacking{Control: "iteration", MustHave: "loop", Rule: "A11",
		Message: "control=iteration but no loop found at depth 1"}, 1.0)
	wA12 := g.Warrant(ControlMatch, &ControlMatchBacking{Control: "sequence", MustNotHave: "switch|loop", Rule: "A12",
		Message: "control=sequence but %s found at depth 1; add control=%s or extract to separate func"}, 1.0)
	wA13 := g.Warrant(ControlMatch, &ControlMatchBacking{Control: "selection", MustNotHave: "loop", Rule: "A13",
		Message: "control=selection but loop found at depth 1; extract loop to separate func"}, 1.0)
	wA14 := g.Warrant(ControlMatch, &ControlMatchBacking{Control: "iteration", MustNotHave: "switch", Rule: "A14",
		Message: "control=iteration but switch found at depth 1; extract switch to separate func"}, 1.0)

	// ── 어노테이션: 코드북 적합 ──
	_ = g.Warrant(InCodebook, &InCodebookBacking{Direction: "value→codebook", Rule: "A2"}, 1.0)
	_ = g.Warrant(InCodebook, &InCodebookBacking{Direction: "codebook→annotation", Rule: "A8"}, 1.0)

	// ── 고유 검사 ──
	_ = g.Warrant(AnnotationAtTop, nil, 1.0)
	_ = g.Warrant(CheckedHashMatch, nil, 1.0)
	wA16 := g.Warrant(ValidDimension, nil, 1.0)

	// ══ 예외: const 전용 파일은 F1 면제 ══
	dConst := g.Defeater(IsConstOnlyDefeater, nil, 1.0)
	g.Defeat(dConst, wF1)

	// ══ 예외: func 없는 파일은 제어 구조 룰 면제 ══
	dNoFunc := g.Defeater(HasNoFunc, nil, 1.0)
	for _, w := range []*toulmin.Rule{wA9, wA10, wA11, wA12, wA13, wA14, wA15, wA16} {
		g.Defeat(dNoFunc, w)
	}

	return g
}
