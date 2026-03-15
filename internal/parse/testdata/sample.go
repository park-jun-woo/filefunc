//ff:func feature=validate type=rule
//ff:what 파일당 func 개수를 검증한다
//ff:why 제1시민은 AI 에이전트
//ff:calls WalkGoFiles, ParseGoFile
//ff:uses GoFile, Violation
package testdata

type SampleParam struct {
	Name string
}

func CheckSample(p SampleParam) error {
	for _, item := range []string{"a", "b"} {
		if item == "a" {
			_ = item
		}
	}
	return nil
}
