//ff:func feature=validate type=util control=sequence
//ff:what test: ruleViolations
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

func ruleViolations(fn func(any, any, any) (bool, any), sf model.SourceFile, cb *model.Codebook, backing any) []model.Violation {
	g := &ValidateGround{File: sf, Codebook: cb, HasChecked: true}
	ok, ev := fn(sf.GetPath(), g, backing)
	if !ok || ev == nil {
		return nil
	}
	return ev.([]model.Violation)
}
