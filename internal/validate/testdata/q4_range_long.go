package testdata

// Q4 violation: range body PURE > 10 lines
func Q4RangeLong() {
	items := []int{1, 2, 3}
	for _, item := range items {
		_ = item
		_ = item + 1
		_ = item + 2
		_ = item + 3
		_ = item + 4
		_ = item + 5
		_ = item + 6
		_ = item + 7
		_ = item + 8
		_ = item + 9
		_ = item + 10
		_ = item + 11
	}
}
