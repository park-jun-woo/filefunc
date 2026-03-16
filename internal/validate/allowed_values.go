//ff:func feature=validate type=util control=sequence
//ff:what codebook의 required+optional에서 주어진 키의 허용 값 목록을 통합 반환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// AllowedValues returns the combined allowed values from required and optional
// for a given annotation key. Returns nil if the key is not in the codebook.
func AllowedValues(cb *model.Codebook, key string) []string {
	req := cb.Required[key]
	opt := cb.Optional[key]
	if req == nil && opt == nil {
		return nil
	}
	var result []string
	result = append(result, req...)
	result = append(result, opt...)
	return result
}
