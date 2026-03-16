//ff:func feature=validate type=util control=selection
//ff:what A10 위반 테스트: selection인데 switch 없음
package testdata

// SelectionNoSwitch is a test func that violates A10.
func SelectionNoSwitch(x int) int {
	return x + 1
}
