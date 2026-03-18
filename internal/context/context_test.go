package context

import (
	"testing"

	"github.com/park-jun-woo/filefunc/internal/model"
)

func TestParseFeatures(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{`["validate", "chain"]`, 2},
		{`some text ["validate"]`, 1},
		{`invalid`, 0},
		{`[]`, 0},
	}
	for _, tt := range tests {
		got := ParseFeatures(tt.input)
		if len(got) != tt.want {
			t.Errorf("ParseFeatures(%q) len = %d, want %d", tt.input, len(got), tt.want)
		}
	}
}

func TestParseScores(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"0.85", 1},
		{"1. 0.85\n2. 0.70", 2},
		{"<think>reasoning</think>\n0.85", 1},
		{"not a score", 0},
	}
	for _, tt := range tests {
		got := ParseScores(tt.input)
		if len(got) != tt.want {
			t.Errorf("ParseScores(%q) len = %d, want %d", tt.input, len(got), tt.want)
		}
	}
}

func TestParseSingleScore(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"0.85", 0.85},
		{"1. 0.70", 0.70},
		{"abc", -1},
		{"", -1},
	}
	for _, tt := range tests {
		got := parseSingleScore(tt.input)
		if got != tt.want {
			t.Errorf("parseSingleScore(%q) = %f, want %f", tt.input, got, tt.want)
		}
	}
}

func TestParseSearch(t *testing.T) {
	got := ParseSearch("feature=validate type=rule")
	if got["feature"] != "validate" || got["type"] != "rule" {
		t.Errorf("ParseSearch = %v", got)
	}
	empty := ParseSearch("")
	if len(empty) != 0 {
		t.Errorf("ParseSearch empty = %v", empty)
	}
}

func TestFilterFeature(t *testing.T) {
	files := []*model.GoFile{
		{Annotation: &model.Annotation{Func: map[string]string{"feature": "validate"}}},
		{Annotation: &model.Annotation{Func: map[string]string{"feature": "chain"}}},
		{Annotation: nil},
	}
	got := FilterFeature(files, []string{"validate"})
	if len(got) != 1 {
		t.Errorf("FilterFeature len = %d, want 1", len(got))
	}
}

func TestMatchAnnotation(t *testing.T) {
	ann := &model.Annotation{
		Func: map[string]string{"feature": "validate", "type": "rule"},
	}
	if !matchAnnotation(ann, map[string]string{"feature": "validate"}) {
		t.Error("expected match")
	}
	if matchAnnotation(ann, map[string]string{"feature": "chain"}) {
		t.Error("expected no match")
	}
	if !matchAnnotation(ann, map[string]string{"feature": "validate", "type": "rule"}) {
		t.Error("expected multi-key match")
	}
}
