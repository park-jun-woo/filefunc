//ff:func feature=validate type=util control=sequence
//ff:what A12 위반 테스트: sequence인데 loop 존재
package testdata

// SequenceWithLoop is a test func that violates A12.
func SequenceWithLoop(items []string) {
	for range items {
	}
}
