//ff:func feature=cli type=command control=sequence
//ff:what 단일 파일에 대해 LLM what-body 검증을 수행하고 결과를 반환
package cli

import (
	"fmt"
	"os"

	"github.com/park-jun-woo/filefunc/internal/annotate"
	"github.com/park-jun-woo/filefunc/internal/llm"
	"github.com/park-jun-woo/filefunc/internal/model"
)

// ProcessLlmcFile verifies a single file with LLM. Returns "pass", "fail", or "skip".
func ProcessLlmcFile(sf model.SourceFile, provider llm.Provider, modelName string, threshold float64) string {
	currentHash, err := CalcBodyHashForLang(sf)
	if err != nil {
		return "skip"
	}

	ann := sf.GetAnnotation()
	if len(ann.Checked) > 0 && ann.Checked["hash"] == currentHash {
		return "skip"
	}

	body := ExtractBodyForLlmc(sf)
	if body == "" {
		return "skip"
	}

	score, err := llm.VerifyWhat(provider, ann.What, body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s: LLM error: %v\n", sf.GetPath(), err)
		return "fail"
	}

	prefix := AnnotationPrefixForLang(sf)
	if score < threshold {
		fmt.Fprintf(os.Stderr, "[FAIL] %s: score=%.2f (threshold=%.2f) — update %swhat\n", sf.GetPath(), score, threshold, prefix)
		return "fail"
	}

	_, err = annotate.WriteAnnotationLine(sf.GetPath(), "checked", fmt.Sprintf("llm=%s hash=%s", modelName, currentHash))
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s: failed to write checked annotation: %v\n", sf.GetPath(), err)
		return "fail"
	}
	fmt.Printf("[PASS] %s: score=%.2f\n", sf.GetPath(), score)
	return "pass"
}
