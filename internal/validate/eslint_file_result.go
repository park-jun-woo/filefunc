//ff:type feature=validate type=model
//ff:what ESLint JSON 출력의 파일 단위 결과를 담는 구조체
package validate

type eslintFileResult struct {
	FilePath   string          `json:"filePath"`
	Messages   []eslintMessage `json:"messages"`
	ErrorCount int             `json:"errorCount"`
}
