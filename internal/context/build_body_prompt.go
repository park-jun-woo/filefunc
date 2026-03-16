//ff:func feature=context type=formatter control=iteration dimension=1
//ff:what 본문 스코어링용 배치 프롬프트를 GoFile 목록으로 생성
package context

import (
	"fmt"
	"os"
	"strings"

	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// BuildBodyPrompt creates a batch prompt for body-based relevance scoring.
// Returns the prompt and indices of files that have extractable bodies.
func BuildBodyPrompt(prompt string, files []*model.GoFile) (string, []int) {
	var lines []string
	var indices []int
	num := 1
	for i, gf := range files {
		src, _ := os.ReadFile(gf.Path)
		body := parse.ExtractFuncSource(gf.Path, src)
		if body == "" {
			continue
		}
		lines = append(lines, fmt.Sprintf("%d. %s", num, body))
		indices = append(indices, i)
		num++
	}
	header := fmt.Sprintf(`사용자가 수정하려는 작업과 각 함수의 관련도를 평가하시오.
관련도: 0.0(무관) ~ 0.5(간접 관련) ~ 1.0(직접 관련)
직접 수정 대상이면 0.8 이상, 영향을 받는 함수면 0.4~0.7, 무관하면 0.2 이하.

작업: "%s"

%s

각 번호에 대해 점수만. 형식: 번호. 점수`, prompt, strings.Join(lines, "\n\n"))
	return header, indices
}
