package testdata

func DepthNested(x int, items []string) {
	for _, item := range items {
		if item == "a" {
			if x > 0 {
				_ = x
			}
		}
	}
}
