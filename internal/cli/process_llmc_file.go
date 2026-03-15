//ff:func feature=cli type=command
//ff:what 단일 파일에 대해 LLM what-body 검증을 수행하고 결과를 반환
//ff:calls CalcBodyHash, ExtractFuncSource, VerifyWhat, WriteAnnotationLine
//ff:uses CalcBodyHash, ExtractFuncSource, GoFile, Provider, VerifyWhat, WriteAnnotationLine
//ff:checked llm=gpt-oss:20b hash=aab86087
package cli

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/filefunc/internal/annotate"
	"github.com/park-jun-woo/filefunc/internal/llm"
	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
)

// ProcessLlmcFile verifies a single file with LLM. Returns "pass", "fail", or "skip".
func ProcessLlmcFile(gf *model.GoFile, provider llm.Provider, modelName string, threshold float64) string {
	currentHash, err := parse.CalcBodyHash(gf.Path)
	if err != nil {
		return "skip"
	}

	if len(gf.Annotation.Checked) > 0 && gf.Annotation.Checked["hash"] == currentHash {
		return "skip"
	}

	src, err := os.ReadFile(gf.Path)
	if err != nil {
		return "skip"
	}

	body := parse.ExtractFuncSource(gf.Path, src)
	if body == "" {
		return "skip"
	}

	score, err := llm.VerifyWhat(provider, gf.Annotation.What, body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s: LLM error: %v\n", gf.Path, err)
		return "fail"
	}

	if score < threshold {
		fmt.Fprintf(os.Stderr, "[FAIL] %s: score=%.2f (threshold=%.2f) — update //ff:what\n", gf.Path, score, threshold)
		return "fail"
	}

	annotate.WriteAnnotationLine(gf.Path, "checked", fmt.Sprintf("llm=%s hash=%s", modelName, currentHash))
	fmt.Printf("[PASS] %s: score=%.2f\n", gf.Path, score)
	return "pass"
}
