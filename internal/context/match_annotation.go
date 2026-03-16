//ff:func feature=context type=util control=iteration dimension=1
//ff:what Annotation의 Func 또는 Type 맵이 모든 query key=value를 만족하는지 판별
package context

import "github.com/park-jun-woo/filefunc/internal/model"

// matchAnnotation returns true if all query key=value pairs match the annotation.
func matchAnnotation(ann *model.Annotation, query map[string]string) bool {
	for k, v := range query {
		funcVal := ann.Func[k]
		typeVal := ann.Type[k]
		if funcVal != v && typeVal != v {
			return false
		}
	}
	return true
}
