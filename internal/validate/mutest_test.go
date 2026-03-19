package validate

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

func mustParse(t *testing.T, path string) *model.GoFile {
	t.Helper()
	gf, err := parse.ParseGoFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return gf
}

func ruleViolations(fn func(any, any, any) (bool, any), gf *model.GoFile, cb *model.Codebook, backing any) []model.Violation {
	g := &ValidateGround{File: gf, Codebook: cb, HasChecked: true}
	ok, ev := fn(gf.Path, g, backing)
	if !ok || ev == nil {
		return nil
	}
	return ev.([]model.Violation)
}

func expectViolation(t *testing.T, violations []model.Violation, rule string) {
	t.Helper()
	if len(violations) == 0 {
		t.Errorf("expected violation %s, got none", rule)
		return
	}
	for _, v := range violations {
		if v.Rule == rule {
			return
		}
	}
	t.Errorf("expected rule %s, got %v", rule, violations)
}

func expectNoViolation(t *testing.T, violations []model.Violation) {
	t.Helper()
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

// --- backing definitions for tests ---

var (
	backingF1 = &CountMaxBacking{Field: "Funcs", Max: 1, Rule: "F1",
		Message: "file contains multiple funcs; expected 1 file 1 func"}
	backingF2 = &CountMaxBacking{Field: "Types", Max: 1, Rule: "F2",
		Message: "file contains multiple types; expected 1 file 1 type"}
	backingF3 = &CountMaxBacking{Field: "Methods", Max: 1, Rule: "F3",
		Message: "file contains multiple methods; expected 1 file 1 method"}
	backingF4 = &ExistsWhenBacking{When: "HasInit", Need: "companion", Rule: "F4",
		Level: "ERROR", Message: "init() must not exist alone; requires accompanying var or func"}
	backingA1f = &ExistsWhenBacking{When: "HasFuncs", Need: "ff:func", Rule: "A1",
		Level: "ERROR", Message: "file with func must have //ff:func annotation"}
	backingA1t = &ExistsWhenBacking{When: "HasTypes", Need: "ff:type", Rule: "A1",
		Level: "ERROR", Message: "file with type must have //ff:type annotation"}
	backingA3 = &ExistsWhenBacking{When: "HasFuncOrType", Need: "ff:what", Rule: "A3",
		Level: "ERROR", Message: "file with func or type must have //ff:what annotation"}
	backingA9 = &ExistsWhenBacking{When: "HasFuncs", Need: "control", Rule: "A9",
		Level: "ERROR", Message: "func file must have control= annotation (sequence, selection, or iteration)"}
	backingA15 = &ExistsWhenBacking{When: "ControlIteration", Need: "dimension", Rule: "A15",
		Level: "ERROR", Message: "control=iteration requires dimension= annotation"}
	backingA2 = &InCodebookBacking{Direction: "value→codebook", Rule: "A2"}
	backingA8 = &InCodebookBacking{Direction: "codebook→annotation", Rule: "A8"}
	backingA10 = &ControlMatchBacking{Control: "selection", MustHave: "switch", Rule: "A10",
		Message: "control=selection but no switch found at depth 1"}
	backingA11 = &ControlMatchBacking{Control: "iteration", MustHave: "loop", Rule: "A11",
		Message: "control=iteration but no loop found at depth 1"}
	backingA12 = &ControlMatchBacking{Control: "sequence", MustNotHave: "switch|loop", Rule: "A12",
		Message: "control=sequence but %s found at depth 1; add control=%s or extract to separate func"}
	backingA13 = &ControlMatchBacking{Control: "selection", MustNotHave: "loop", Rule: "A13",
		Message: "control=selection but loop found at depth 1; extract loop to separate func"}
	backingA14 = &ControlMatchBacking{Control: "iteration", MustNotHave: "switch", Rule: "A14",
		Message: "control=iteration but switch found at depth 1; extract switch to separate func"}
)

// F1
func TestMutest_F1(t *testing.T) {
	expectViolation(t, ruleViolations(CountMax, mustParse(t, "testdata/multi_func.go"), nil, backingF1), "F1")
	expectNoViolation(t, ruleViolations(CountMax, mustParse(t, "testdata/clean.go"), nil, backingF1))
}

// F2
func TestMutest_F2(t *testing.T) {
	expectViolation(t, ruleViolations(CountMax, mustParse(t, "testdata/multi_type.go"), nil, backingF2), "F2")
}

// F3
func TestMutest_F3(t *testing.T) {
	expectViolation(t, ruleViolations(CountMax, mustParse(t, "testdata/multi_method.go"), nil, backingF3), "F3")
}

// F4
func TestMutest_F4(t *testing.T) {
	expectViolation(t, ruleViolations(ExistsWhen, mustParse(t, "testdata/init_alone.go"), nil, backingF4), "F4")
	expectNoViolation(t, ruleViolations(ExistsWhen, mustParse(t, "testdata/clean.go"), nil, backingF4))
}

// Q1
func TestMutest_Q1(t *testing.T) {
	expectViolation(t, ruleViolations(CheckDepthLimit, mustParse(t, "testdata/deep_nesting.go"), nil, nil), "Q1")
	expectNoViolation(t, ruleViolations(CheckDepthLimit, mustParse(t, "testdata/clean.go"), nil, nil))
}

// Q2
func TestMutest_Q2(t *testing.T) {
	expectViolation(t, ruleViolations(CheckFuncLines, mustParse(t, "testdata/long_func.go"), nil, nil), "Q2")
}

// Q3
func TestMutest_Q3(t *testing.T) {
	expectViolation(t, ruleViolations(CheckFuncLines, mustParse(t, "testdata/medium_func.go"), nil, nil), "Q3")
}

// A1
func TestMutest_A1(t *testing.T) {
	expectViolation(t, ruleViolations(ExistsWhen, mustParse(t, "testdata/no_annotation.go"), nil, backingA1f), "A1")
}

// A2
func TestMutest_A2(t *testing.T) {
	cb := &model.Codebook{
		Required: map[string]map[string]string{
			"feature": {"validate": "", "parse": ""},
			"type":    {"rule": "", "parser": ""},
		},
	}
	expectViolation(t, ruleViolations(InCodebook, mustParse(t, "testdata/bad_codebook_value.go"), cb, backingA2), "A2")
}

// A3
func TestMutest_A3(t *testing.T) {
	expectViolation(t, ruleViolations(ExistsWhen, mustParse(t, "testdata/no_what.go"), nil, backingA3), "A3")
}

// A6
func TestMutest_A6(t *testing.T) {
	expectViolation(t, ruleViolations(AnnotationAtTop, mustParse(t, "testdata/annotation_after_func.go"), nil, nil), "A6")
}

// A13
func TestMutest_A13(t *testing.T) {
	expectViolation(t, ruleViolations(ControlMatch, mustParse(t, "testdata/control_selection_with_loop.go"), nil, backingA13), "A13")
	expectNoViolation(t, ruleViolations(ControlMatch, mustParse(t, "testdata/clean.go"), nil, backingA13))
}

// A14
func TestMutest_A14(t *testing.T) {
	expectViolation(t, ruleViolations(ControlMatch, mustParse(t, "testdata/control_iteration_with_switch.go"), nil, backingA14), "A14")
	expectNoViolation(t, ruleViolations(ControlMatch, mustParse(t, "testdata/clean.go"), nil, backingA14))
}

// A15
func TestMutest_A15(t *testing.T) {
	expectViolation(t, ruleViolations(ExistsWhen, mustParse(t, "testdata/iter_no_dimension.go"), nil, backingA15), "A15")
	expectNoViolation(t, ruleViolations(ExistsWhen, mustParse(t, "testdata/clean.go"), nil, backingA15))
}

// A16
func TestMutest_A16(t *testing.T) {
	expectViolation(t, ruleViolations(ValidDimension, mustParse(t, "testdata/bad_dimension_value.go"), nil, nil), "A16")
	expectNoViolation(t, ruleViolations(ValidDimension, mustParse(t, "testdata/clean.go"), nil, nil))
}

// Q1 dimension
func TestMutest_Q1_Dimension(t *testing.T) {
	expectNoViolation(t, ruleViolations(CheckDepthLimit, mustParse(t, "testdata/dimension2_depth3.go"), nil, nil))
}

// Q3 backtick hint
func TestMutest_Q3_Backtick(t *testing.T) {
	violations := ruleViolations(CheckFuncLines, mustParse(t, "testdata/q3_backtick.go"), nil, nil)
	expectViolation(t, violations, "Q3")
	if len(violations) > 0 && !strings.Contains(violations[0].Message, "var-only file") {
		t.Errorf("expected backtick hint in message, got %q", violations[0].Message)
	}
}

// A9
func TestMutest_A9(t *testing.T) {
	expectViolation(t, ruleViolations(ExistsWhen, mustParse(t, "testdata/no_control.go"), nil, backingA9), "A9")
}

// A10
func TestMutest_A10(t *testing.T) {
	expectViolation(t, ruleViolations(ControlMatch, mustParse(t, "testdata/selection_no_switch.go"), nil, backingA10), "A10")
}

// A11
func TestMutest_A11(t *testing.T) {
	expectViolation(t, ruleViolations(ControlMatch, mustParse(t, "testdata/iteration_no_loop.go"), nil, backingA11), "A11")
}

// A12
func TestMutest_A12(t *testing.T) {
	expectViolation(t, ruleViolations(ControlMatch, mustParse(t, "testdata/sequence_with_loop.go"), nil, backingA12), "A12")
}

// A12: control="" should NOT fire A12 (A9 handles it)
func TestMutest_A12_NoControl(t *testing.T) {
	expectNoViolation(t, ruleViolations(ControlMatch, mustParse(t, "testdata/no_control.go"), nil, backingA12))
}

// A1: file with func+type should check both annotations
func TestMutest_A1_FuncAndType(t *testing.T) {
	gf := mustParse(t, "testdata/sample_with_func_and_type.go")
	// has //ff:func but not //ff:type — should fire A1 for type
	violations := ruleViolations(ExistsWhen, gf, nil, backingA1t)
	expectViolation(t, violations, "A1")
}

// A7
func TestMutest_A7(t *testing.T) {
	expectViolation(t, ruleViolations(CheckedHashMatch, mustParse(t, "testdata/checked_hash_mismatch.go"), nil, nil), "A7")
}

// A8
func TestMutest_A8(t *testing.T) {
	cb := &model.Codebook{
		Required: map[string]map[string]string{
			"feature": {"validate": "", "parse": ""},
			"type":    {"rule": "", "parser": ""},
		},
	}
	expectViolation(t, ruleViolations(InCodebook, mustParse(t, "testdata/missing_required_key.go"), cb, backingA8), "A8")
}

// C1
func TestMutest_C1(t *testing.T) {
	cb, err := parse.ParseCodebook("testdata/codebook_empty_required.yaml")
	if err != nil {
		t.Fatal(err)
	}
	expectViolation(t, CheckCodebookRequiredKeys(cb), "C1")
}

// C2
func TestMutest_C2(t *testing.T) {
	expectViolation(t, CheckCodebookDuplicates("testdata/codebook_duplicate_key.yaml"), "C2")
}

// C3
func TestMutest_C3(t *testing.T) {
	cb, err := parse.ParseCodebook("testdata/codebook_bad_format.yaml")
	if err != nil {
		t.Fatal(err)
	}
	expectViolation(t, CheckCodebookValueFormat(cb), "C3")
}

// C4
func TestMutest_C4(t *testing.T) {
	cb, err := parse.ParseCodebook("testdata/codebook_no_description.yaml")
	if err != nil {
		t.Fatal(err)
	}
	expectViolation(t, CheckCodebookDescription(cb), "C4")
}
