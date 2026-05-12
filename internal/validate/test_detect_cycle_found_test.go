//ff:func feature=validate type=util control=sequence
//ff:what test: detectCycle finds cycle in A->B->C->A graph
package validate

import "testing"

func TestDetectCycleFound(t *testing.T) {
	graph := map[string][]string{
		"A": {"B"},
		"B": {"C"},
		"C": {"A"},
	}
	cycles := detectCycle(graph)
	if len(cycles) == 0 {
		t.Error("expected at least one cycle, got none")
	}
}
