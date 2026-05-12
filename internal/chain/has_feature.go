//ff:func feature=chain type=util control=sequence
//ff:what SourceFile의 어노테이션이 지정된 feature 값을 가지는지 판별
package chain

import "github.com/park-jun-woo/filefunc/internal/model"

func hasFeature(sf model.SourceFile, feature string) bool {
	ann := sf.GetAnnotation()
	if ann == nil {
		return false
	}
	return ann.Func["feature"] == feature || ann.Type["feature"] == feature
}
