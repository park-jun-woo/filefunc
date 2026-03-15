package parse

import (
	"testing"
)

func TestParseCodebook(t *testing.T) {
	cb, err := ParseCodebook("testdata/codebook.yaml")
	if err != nil {
		t.Fatalf("ParseCodebook failed: %v", err)
	}

	if len(cb.Feature) != 2 {
		t.Errorf("Feature count = %d, want 2", len(cb.Feature))
	}

	if cb.Feature[0] != "validate" {
		t.Errorf("Feature[0] = %q, want %q", cb.Feature[0], "validate")
	}

	if len(cb.Type) != 2 {
		t.Errorf("Type count = %d, want 2", len(cb.Type))
	}

	if len(cb.Level) != 2 {
		t.Errorf("Level count = %d, want 2", len(cb.Level))
	}
}
