//ff:func feature=validate type=util control=sequence
//ff:what test: detectCycle returns empty for acyclic graph
package validate

import "testing"

func TestDetectCycleNone(t *testing.T) {
	graph := map[string][]string{
		"A": {"B", "C"},
		"B": {},
		"C": {},
	}
	cycles := detectCycle(graph)
	if len(cycles) != 0 {
		t.Errorf("expected no cycles, got %v", cycles)
	}
}
