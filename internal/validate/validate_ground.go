//ff:type feature=validate type=model
//ff:what toulmin Evaluate에 전달되는 ground — GoFile + Codebook + 프로젝트 상태
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// ValidateGround is the ground passed to toulmin rule functions.
type ValidateGround struct {
	File       *model.GoFile
	Codebook   *model.Codebook
	HasChecked bool
}
