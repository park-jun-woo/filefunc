//ff:func feature=validate type=rule control=selection
//ff:what 코드북과 어노테이션의 교차 적합성을 검증하여 violation 반환
package validate

// InCodebook returns (true, []model.Violation) if codebook and annotation values are inconsistent.
func InCodebook(claim any, ground any, backing any) (bool, any) {
	b := backing.(*InCodebookBacking)
	g := ground.(*ValidateGround)
	gf := g.File
	cb := g.Codebook
	if cb == nil || gf.Annotation == nil {
		return false, nil
	}

	switch b.Direction {
	case "value→codebook":
		return checkValuesInCodebook(gf, cb, b.Rule)
	case "codebook→annotation":
		return checkCodebookInAnnotation(gf, cb, b.Rule)
	}
	return false, nil
}
