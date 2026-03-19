//ff:func feature=validate type=rule control=sequence
//ff:what func 없는 파일 defeater — A9~A16, Q1 등 func 전용 룰에서 예외
//ff:why annotation/control rules only apply to files with funcs
package validate

// DefeaterNoFunc returns true if the file has no funcs and no annotation.
// Used as a defeater against control/dimension/annotation rules.
func HasNoFunc(claim any, ground any, backing any) (bool, any) {
	gf := ground.(*ValidateGround).File
	return len(gf.Funcs) == 0 || gf.Annotation == nil, nil
}
