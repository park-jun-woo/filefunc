//ff:type feature=validate type=model
//ff:what ESLint JSON 출력의 개별 린트 메시지를 담는 구조체
package validate

type eslintMessage struct {
	RuleID   string `json:"ruleId"`
	Message  string `json:"message"`
	Line     int    `json:"line"`
	Severity int    `json:"severity"`
}
