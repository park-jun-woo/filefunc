package validate

import (
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
	expectViolation(t, CheckOneFileOneFunc(mustParse(t, "testdata/multi_func.go")), "F1")
	expectNoViolation(t, CheckOneFileOneFunc(mustParse(t, "testdata/clean.go")))
}

// F2
func TestMutest_F2(t *testing.T) {
	expectViolation(t, CheckOneFileOneType(mustParse(t, "testdata/multi_type.go")), "F2")
}

// F3
func TestMutest_F3(t *testing.T) {
	expectViolation(t, CheckOneFileOneMethod(mustParse(t, "testdata/multi_method.go")), "F3")
}

// F4
func TestMutest_F4(t *testing.T) {
	expectViolation(t, CheckInitStandalone(mustParse(t, "testdata/init_alone.go")), "F4")
	expectNoViolation(t, CheckInitStandalone(mustParse(t, "testdata/clean.go")))
}

// Q1
func TestMutest_Q1(t *testing.T) {
	expectViolation(t, CheckNestingDepth(mustParse(t, "testdata/deep_nesting.go")), "Q1")
	expectNoViolation(t, CheckNestingDepth(mustParse(t, "testdata/clean.go")))
}

// Q2
func TestMutest_Q2(t *testing.T) {
	expectViolation(t, CheckFuncLines(mustParse(t, "testdata/long_func.go")), "Q2")
}

// Q3
func TestMutest_Q3(t *testing.T) {
	expectViolation(t, CheckFuncLines(mustParse(t, "testdata/medium_func.go")), "Q3")
}

// A1
func TestMutest_A1(t *testing.T) {
	expectViolation(t, CheckAnnotationRequired(mustParse(t, "testdata/no_annotation.go")), "A1")
}

// A2
func TestMutest_A2(t *testing.T) {
	cb := &model.Codebook{
		Required: map[string][]string{
			"feature": {"validate", "parse"},
			"type":    {"rule", "parser"},
		},
	}
	expectViolation(t, CheckCodebookValues(mustParse(t, "testdata/bad_codebook_value.go"), cb), "A2")
}

// A3
func TestMutest_A3(t *testing.T) {
	expectViolation(t, CheckWhatRequired(mustParse(t, "testdata/no_what.go")), "A3")
}

// A6
func TestMutest_A6(t *testing.T) {
	expectViolation(t, CheckAnnotationPosition(mustParse(t, "testdata/annotation_after_func.go")), "A6")
}

// A13
func TestMutest_A13(t *testing.T) {
	expectViolation(t, CheckControlSelectionNoLoop(mustParse(t, "testdata/control_selection_with_loop.go")), "A13")
	expectNoViolation(t, CheckControlSelectionNoLoop(mustParse(t, "testdata/clean.go")))
}

// A14
func TestMutest_A14(t *testing.T) {
	expectViolation(t, CheckControlIterationNoSwitch(mustParse(t, "testdata/control_iteration_with_switch.go")), "A14")
	expectNoViolation(t, CheckControlIterationNoSwitch(mustParse(t, "testdata/clean.go")))
}

// A15
func TestMutest_A15(t *testing.T) {
	expectViolation(t, CheckDimensionRequired(mustParse(t, "testdata/iter_no_dimension.go")), "A15")
	expectNoViolation(t, CheckDimensionRequired(mustParse(t, "testdata/clean.go")))
}

// A16
func TestMutest_A16(t *testing.T) {
	expectViolation(t, CheckDimensionValue(mustParse(t, "testdata/bad_dimension_value.go")), "A16")
	expectNoViolation(t, CheckDimensionValue(mustParse(t, "testdata/clean.go")))
}

// Q1 dimension
func TestMutest_Q1_Dimension(t *testing.T) {
	expectNoViolation(t, CheckNestingDepth(mustParse(t, "testdata/dimension2_depth3.go")))
}
