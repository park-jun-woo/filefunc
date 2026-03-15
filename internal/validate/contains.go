//ff:func feature=validate type=util
//ff:what 문자열 슬라이스에 특정 항목이 포함되어 있는지 확인
package validate

// Contains returns true if the slice contains the given string.
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
