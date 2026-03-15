//ff:func feature=cli type=formatter
//ff:what what과 func body로 LLM 검증 프롬프트를 생성
//ff:checked llm=gpt-oss:20b hash=f3a78f60
package llm

import "fmt"

// BuildPrompt creates a prompt for LLM verification of what-body match.
func BuildPrompt(what string, body string) string {
	return fmt.Sprintf(`You are a code reviewer. Rate how accurately the description matches the Go function.

Description: %s

Function:
%s

Respond with a single number between 0.0 and 1.0 only.
0.0 = completely wrong, 1.0 = perfectly accurate.`, what, body)
}
