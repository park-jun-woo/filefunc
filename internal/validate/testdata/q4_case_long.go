package testdata

// Q4 violation: case body PURE > 10 lines
func Q4CaseLong(kind string) {
	switch kind {
	case "a":
		_ = 1
		_ = 2
		_ = 3
		_ = 4
		_ = 5
		_ = 6
		_ = 7
		_ = 8
		_ = 9
		_ = 10
		_ = 11
	case "b":
		_ = 1
	}
}
