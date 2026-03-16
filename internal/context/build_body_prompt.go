//ff:func feature=context type=formatter control=iteration dimension=1
//ff:what 본문 스코어링용 배치 프롬프트를 생성
package context

import (
	"fmt"
	"os"
	"strings"

	"github.com/park-jun-woo/filefunc/internal/chain"
	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// BuildBodyPrompt creates a batch prompt for body-based relevance scoring.
// Only includes chon=2+ results (chon=1 is always 1.0).
func BuildBodyPrompt(prompt string, results []chain.ChonResult, fileMap map[string]*model.GoFile) (string, []int) {
	var lines []string
	var indices []int
	num := 1
	for i, r := range results {
		if r.Chon <= 1 {
			continue
		}
		gf := fileMap[r.Name]
		if gf == nil {
			continue
		}
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
