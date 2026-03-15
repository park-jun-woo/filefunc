//ff:func feature=validate type=util
//ff:what 코드북에서 주어진 키의 허용 값 목록을 반환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// AllowedValues returns the codebook's allowed values for a given annotation key.
// Returns nil if the key is not a codebook-controlled field.
func AllowedValues(cb *model.Codebook, key string) []string {
	switch key {
	case "feature":
		return cb.Feature
	case "type":
		return cb.Type
	case "pattern":
		return cb.Pattern
	case "level":
		return cb.Level
	}
	return nil
}
