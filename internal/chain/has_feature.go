//ff:func feature=chain type=util control=sequence
//ff:what GoFile의 어노테이션이 지정된 feature 값을 가지는지 판별
package chain

import "github.com/park-jun-woo/filefunc/internal/model"

func hasFeature(gf *model.GoFile, feature string) bool {
	if gf.Annotation == nil {
		return false
	}
	return gf.Annotation.Func["feature"] == feature || gf.Annotation.Type["feature"] == feature
}
