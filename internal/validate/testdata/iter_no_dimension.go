//ff:func feature=validate type=util control=iteration
//ff:what A15 위반 테스트: iteration인데 dimension 없음
package testdata

// IterNoDimension is a test func that violates A15.
func IterNoDimension(items []string) int {
	for range items {
	}
	return 0
}
