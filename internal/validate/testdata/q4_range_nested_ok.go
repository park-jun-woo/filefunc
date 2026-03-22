package testdata

// Q4 pass: range body has inner control, but PURE is within limit
func Q4RangeNestedOK() {
	items := []int{1, 2, 3}
	for _, item := range items {
		_ = item
		_ = item + 1
		_ = item + 2
		if item > 0 {
			_ = item + 3
			_ = item + 4
			_ = item + 5
			_ = item + 6
			_ = item + 7
			_ = item + 8
			_ = item + 9
		}
	}
}
