package parse

import (
	"testing"
)

func TestParseCodebook(t *testing.T) {
	cb, err := ParseCodebook("testdata/codebook.yaml")
	if err != nil {
		t.Fatalf("ParseCodebook failed: %v", err)
	}

	if len(cb.Required["feature"]) != 2 {
		t.Errorf("Required feature count = %d, want 2", len(cb.Required["feature"]))
	}

	if cb.Required["feature"][0] != "validate" {
		t.Errorf("Required feature[0] = %q, want %q", cb.Required["feature"][0], "validate")
	}

	if len(cb.Required["type"]) != 2 {
		t.Errorf("Required type count = %d, want 2", len(cb.Required["type"]))
	}

	if len(cb.Optional["level"]) != 2 {
		t.Errorf("Optional level count = %d, want 2", len(cb.Optional["level"]))
	}
}
