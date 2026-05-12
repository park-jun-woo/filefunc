//ff:func feature=validate type=rule control=sequence
//ff:what F7 defeater — const/var 전용 파일은 F1에서 예외
//ff:why const-only files have no func to validate
package validate

// DefeaterConstOnly returns true if the file contains only const/var declarations.
// Used as a defeater against F1.
func IsConstOnlyDefeater(claim any, ground any, backing any) (bool, any) {
	sf := ground.(*ValidateGround).File
	return len(sf.GetFuncs()) == 0 && len(sf.GetTypes()) == 0 && len(sf.GetMethods()) == 0 && !sf.GetHasInit(), nil
}
