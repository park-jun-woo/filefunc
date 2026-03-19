package validate

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// --- Exception tests: rules must NOT fire ---

// F5: _test.go with multiple funcs → F1 should not fire (defeated by IsTestFile)
func TestException_F5_TestFile(t *testing.T) {
	gf, err := parse.ParseGoFile("testdata/test_file_test.go")
	if err != nil {
		t.Fatal(err)
	}
	violations := RunAll([]*model.GoFile{gf}, nil)
	expectNoViolation(t, violations)
}

// F6: func + unexported param type → F2 should not fire
func TestException_F6_ParamType(t *testing.T) {
	gf, err := parse.ParseGoFile("testdata/func_with_param_type.go")
	if err != nil {
		t.Fatal(err)
	}
	expectNoViolation(t, ruleViolations(CountMax, gf, nil, backingF2))
}

// F7: const-only file → F1 should not fire (defeated by IsConstOnlyDefeater)
func TestException_F7_ConstOnly(t *testing.T) {
	gf, err := parse.ParseGoFile("testdata/const_only.go")
	if err != nil {
		t.Fatal(err)
	}
	violations := RunAll([]*model.GoFile{gf}, nil)
	expectNoViolation(t, violations)
}

// F4 exception: var + init() → should not fire
func TestException_F4_VarWithInit(t *testing.T) {
	gf, err := parse.ParseGoFile("testdata/var_with_init.go")
	if err != nil {
		t.Fatal(err)
	}
	expectNoViolation(t, ruleViolations(ExistsWhen, gf, nil, backingF4))
}

// --- Clean: all rules pass ---

func TestClean_AllRules(t *testing.T) {
	gf, err := parse.ParseGoFile("testdata/clean.go")
	if err != nil {
		t.Fatal(err)
	}
	violations := RunAll([]*model.GoFile{gf}, nil)
	expectNoViolation(t, violations)
}

// --- //ff:type path tests ---

// A1: type-only file without //ff:type → should fire
func TestType_A1_Missing(t *testing.T) {
	gf, err := parse.ParseGoFile("testdata/type_no_annotation.go")
	if err != nil {
		t.Fatal(err)
	}
	expectViolation(t, ruleViolations(ExistsWhen, gf, nil, backingA1t), "A1")
}

// A1: type-only file with //ff:type → should pass
func TestType_A1_Present(t *testing.T) {
	gf, err := parse.ParseGoFile("testdata/type_with_annotation.go")
	if err != nil {
		t.Fatal(err)
	}
	expectNoViolation(t, ruleViolations(ExistsWhen, gf, nil, backingA1t))
}

// A2: //ff:type with bad codebook value → should fire
func TestType_A2_BadCodebook(t *testing.T) {
	gf, err := parse.ParseGoFile("testdata/type_bad_codebook.go")
	if err != nil {
		t.Fatal(err)
	}
	cb := &model.Codebook{
		Required: map[string]map[string]string{
			"feature": {"validate": "", "parse": ""},
			"type":    {"rule": "", "model": ""},
		},
	}
	expectViolation(t, ruleViolations(InCodebook, gf, cb, backingA2), "A2")
}
