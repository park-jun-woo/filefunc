//ff:func feature=validate type=util control=iteration dimension=1
//ff:what A11 위반 테스트: iteration인데 loop 없음
package testdata

// IterationNoLoop is a test func that violates A11.
func IterationNoLoop(x int) int {
	return x + 1
}
