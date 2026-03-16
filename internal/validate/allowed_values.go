//ff:func feature=validate type=util control=iteration dimension=1
//ff:what codebook의 required+optional에서 주어진 키의 허용 값 목록을 통합 반환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// AllowedValues returns the combined allowed values from required and optional
// for a given annotation key. Returns nil if the key is not in the codebook.
func AllowedValues(cb *model.Codebook, key string) []string {
	var result []string
	for name := range cb.Required[key] {
		result = append(result, name)
	}
	for name := range cb.Optional[key] {
		result = append(result, name)
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
