//ff:type feature=validate type=model
//ff:what toulmin Evaluate에 전달되는 ground — SourceFile + Codebook + 프로젝트 상태
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// ValidateGround is the ground passed to toulmin rule functions.
type ValidateGround struct {
	File       model.SourceFile
	Codebook   *model.Codebook
	HasChecked bool
}
