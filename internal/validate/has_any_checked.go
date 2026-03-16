//ff:func feature=validate type=util control=iteration
//ff:what 프로젝트 파일 목록에서 //ff:checked가 하나라도 있는지 판별
//ff:checked llm=gpt-oss:20b hash=70d865a6
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// HasAnyChecked returns true if any GoFile in the list has a //ff:checked annotation.
func HasAnyChecked(files []*model.GoFile) bool {
	for _, gf := range files {
		if gf.Annotation != nil && len(gf.Annotation.Checked) > 0 {
			return true
		}
	}
	return false
}
