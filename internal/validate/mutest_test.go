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

func ruleViolations(fn func(any, any) (bool, any), gf *model.GoFile, cb *model.Codebook) []model.Violation {
	g := &ValidateGround{File: gf, Codebook: cb, HasChecked: true}
	ok, ev := fn(gf.Path, g)
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

// F1
func TestMutest_F1(t *testing.T) {
	expectViolation(t, ruleViolations(RuleF1, mustParse(t, "testdata/multi_func.go"), nil), "F1")
	expectNoViolation(t, ruleViolations(RuleF1, mustParse(t, "testdata/clean.go"), nil))
}

// F2
func TestMutest_F2(t *testing.T) {
	expectViolation(t, ruleViolations(RuleF2, mustParse(t, "testdata/multi_type.go"), nil), "F2")
}

// F3
func TestMutest_F3(t *testing.T) {
	expectViolation(t, ruleViolations(RuleF3, mustParse(t, "testdata/multi_method.go"), nil), "F3")
}

// F4
func TestMutest_F4(t *testing.T) {
	expectViolation(t, ruleViolations(RuleF4, mustParse(t, "testdata/init_alone.go"), nil), "F4")
	expectNoViolation(t, ruleViolations(RuleF4, mustParse(t, "testdata/clean.go"), nil))
}

// Q1
func TestMutest_Q1(t *testing.T) {
	expectViolation(t, ruleViolations(RuleQ1, mustParse(t, "testdata/deep_nesting.go"), nil), "Q1")
	expectNoViolation(t, ruleViolations(RuleQ1, mustParse(t, "testdata/clean.go"), nil))
}

// Q2
func TestMutest_Q2(t *testing.T) {
	expectViolation(t, ruleViolations(RuleQ2Q3, mustParse(t, "testdata/long_func.go"), nil), "Q2")
}

// Q3
func TestMutest_Q3(t *testing.T) {
	expectViolation(t, ruleViolations(RuleQ2Q3, mustParse(t, "testdata/medium_func.go"), nil), "Q3")
}

// A1
func TestMutest_A1(t *testing.T) {
	expectViolation(t, ruleViolations(RuleA1, mustParse(t, "testdata/no_annotation.go"), nil), "A1")
}

// A2
func TestMutest_A2(t *testing.T) {
	cb := &model.Codebook{
		Required: map[string]map[string]string{
			"feature": {"validate": "", "parse": ""},
			"type":    {"rule": "", "parser": ""},
		},
	}
	expectViolation(t, ruleViolations(RuleA2, mustParse(t, "testdata/bad_codebook_value.go"), cb), "A2")
}

// A3
func TestMutest_A3(t *testing.T) {
	expectViolation(t, ruleViolations(RuleA3, mustParse(t, "testdata/no_what.go"), nil), "A3")
}

// A6
func TestMutest_A6(t *testing.T) {
	expectViolation(t, ruleViolations(RuleA6, mustParse(t, "testdata/annotation_after_func.go"), nil), "A6")
}

// A13
func TestMutest_A13(t *testing.T) {
	expectViolation(t, ruleViolations(RuleA13, mustParse(t, "testdata/control_selection_with_loop.go"), nil), "A13")
	expectNoViolation(t, ruleViolations(RuleA13, mustParse(t, "testdata/clean.go"), nil))
}

// A14
func TestMutest_A14(t *testing.T) {
	expectViolation(t, ruleViolations(RuleA14, mustParse(t, "testdata/control_iteration_with_switch.go"), nil), "A14")
	expectNoViolation(t, ruleViolations(RuleA14, mustParse(t, "testdata/clean.go"), nil))
}

// A15
func TestMutest_A15(t *testing.T) {
	expectViolation(t, ruleViolations(RuleA15, mustParse(t, "testdata/iter_no_dimension.go"), nil), "A15")
	expectNoViolation(t, ruleViolations(RuleA15, mustParse(t, "testdata/clean.go"), nil))
}

// A16
func TestMutest_A16(t *testing.T) {
	expectViolation(t, ruleViolations(RuleA16, mustParse(t, "testdata/bad_dimension_value.go"), nil), "A16")
	expectNoViolation(t, ruleViolations(RuleA16, mustParse(t, "testdata/clean.go"), nil))
}

// Q1 dimension
func TestMutest_Q1_Dimension(t *testing.T) {
	expectNoViolation(t, ruleViolations(RuleQ1, mustParse(t, "testdata/dimension2_depth3.go"), nil))
}

// Q3 backtick hint
func TestMutest_Q3_Backtick(t *testing.T) {
	violations := ruleViolations(RuleQ2Q3, mustParse(t, "testdata/q3_backtick.go"), nil)
	expectViolation(t, violations, "Q3")
	if len(violations) > 0 && !strings.Contains(violations[0].Message, "var-only file") {
		t.Errorf("expected backtick hint in message, got %q", violations[0].Message)
	}
}

// A9
func TestMutest_A9(t *testing.T) {
	expectViolation(t, ruleViolations(RuleA9, mustParse(t, "testdata/no_control.go"), nil), "A9")
}

// A10
func TestMutest_A10(t *testing.T) {
	expectViolation(t, ruleViolations(RuleA10, mustParse(t, "testdata/selection_no_switch.go"), nil), "A10")
}

// A11
func TestMutest_A11(t *testing.T) {
	expectViolation(t, ruleViolations(RuleA11, mustParse(t, "testdata/iteration_no_loop.go"), nil), "A11")
}

// A12
func TestMutest_A12(t *testing.T) {
	expectViolation(t, ruleViolations(RuleA12, mustParse(t, "testdata/sequence_with_loop.go"), nil), "A12")
}

// A7
func TestMutest_A7(t *testing.T) {
	expectViolation(t, ruleViolations(RuleA7, mustParse(t, "testdata/checked_hash_mismatch.go"), nil), "A7")
}

// A8
func TestMutest_A8(t *testing.T) {
	cb := &model.Codebook{
		Required: map[string]map[string]string{
			"feature": {"validate": "", "parse": ""},
			"type":    {"rule": "", "parser": ""},
		},
	}
	expectViolation(t, ruleViolations(RuleA8, mustParse(t, "testdata/missing_required_key.go"), cb), "A8")
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
