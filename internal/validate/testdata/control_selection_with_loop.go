//ff:func feature=validate type=util control=selection
//ff:what A13 위반 테스트: selection인데 loop 존재
package testdata

// ControlSelectionWithLoop is a test func that violates A13.
func ControlSelectionWithLoop(items []string) string {
	for _, item := range items {
		if item == "found" {
			return item
		}
	}
	return ""
}
