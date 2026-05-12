//ff:func feature=validate type=util control=iteration dimension=1
//ff:what 프로젝트 파일 목록에서 //ff:checked가 하나라도 있는지 판별
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// HasAnyChecked returns true if any SourceFile in the list has a //ff:checked annotation.
func HasAnyChecked(files []model.SourceFile) bool {
	for _, sf := range files {
		ann := sf.GetAnnotation()
		if ann != nil && len(ann.Checked) > 0 {
			return true
		}
	}
	return false
}
