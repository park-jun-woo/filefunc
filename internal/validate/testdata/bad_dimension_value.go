//ff:func feature=validate type=util control=iteration dimension=0
//ff:what A16 위반 테스트: dimension=0
package testdata

// BadDimensionValue is a test func that violates A16.
func BadDimensionValue(items []string) int {
	for range items {
	}
	return 0
}
