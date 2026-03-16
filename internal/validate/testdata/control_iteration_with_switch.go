//ff:func feature=validate type=util control=iteration
//ff:what A14 위반 테스트: iteration인데 switch 존재
package testdata

// ControlIterationWithSwitch is a test func that violates A14.
func ControlIterationWithSwitch(kind string) int {
	switch kind {
	case "a":
		return 1
	case "b":
		return 2
	default:
		return 0
	}
}
