//ff:func feature=chain type=util control=sequence
//ff:what test: buildTestGraph
package chain

func buildTestGraph() *CallGraph {
	// Caller -> HelperA, HelperB
	// HelperA -> Leaf
	g := &CallGraph{
		Children: map[string][]string{
			"testdata.Caller":  {"testdata.HelperA", "testdata.HelperB"},
			"testdata.HelperA": {"testdata.Leaf"},
		},
		Parents: map[string][]string{
			"testdata.HelperA": {"testdata.Caller"},
			"testdata.HelperB": {"testdata.Caller"},
			"testdata.Leaf":    {"testdata.HelperA"},
		},
	}
	return g
}
