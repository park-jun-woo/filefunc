//ff:func feature=validate type=util control=iteration dimension=2
//ff:what Q1 통과 테스트: dimension=2이면 depth 3 허용
package testdata

// Dimension2Depth3 is a test func that passes Q1 with dimension=2.
func Dimension2Depth3(layers [][]string) {
	for _, layer := range layers {
		for _, item := range layer {
			_ = item
		}
	}
}
