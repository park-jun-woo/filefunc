//ff:func feature=validate type=rule control=sequence
//ff:what F4: init()만 단독으로 존재하는 파일 검증
//ff:checked llm=gpt-oss:20b hash=bbfc055b
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// CheckInitStandalone checks F4: init() must not exist alone in a file.
// It must accompany a var or func.
func CheckInitStandalone(gf *model.GoFile) []model.Violation {
	if !gf.HasInit {
		return nil
	}
	if len(gf.Funcs) == 0 && len(gf.Vars) == 0 && len(gf.Methods) == 0 && len(gf.Types) == 0 {
		return []model.Violation{{
			File:    gf.Path,
			Rule:    "F4",
			Level:   "ERROR",
			Message: "init() must not exist alone; requires accompanying var or func",
		}}
	}
	return nil
}
