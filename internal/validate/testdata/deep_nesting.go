package testdata

func DeepNesting() {
	for i := 0; i < 10; i++ {
		if i > 5 {
			for j := 0; j < 10; j++ {
				_ = j
			}
		}
	}
}
