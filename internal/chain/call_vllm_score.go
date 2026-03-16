//ff:func feature=chain type=loader control=sequence
//ff:what vLLM /v1/score 엔드포인트에 단일 query-document 쌍을 전송하여 관련도 점수 반환
package chain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func callVLLMScore(endpoint string, modelName string, text1 string, text2 string) (float64, error) {
	reqBody, err := json.Marshal(scoreRequest{Model: modelName, Text1: text1, Text2: text2})
	if err != nil {
		return 0, err
	}
	resp, err := http.Post(endpoint+"/v1/score", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return 0, fmt.Errorf("vLLM server not available at %s\nInstall: pip install vllm\nRun:     vllm serve Qwen/Qwen3-Reranker-0.6B --task score --hf_overrides '{\"architectures\":[\"Qwen3ForSequenceClassification\"],\"classifier_from_token\":[\"no\",\"yes\"],\"is_original_qwen3_reranker\":true}'", endpoint)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("vllm returned %d: %s", resp.StatusCode, string(body))
	}
	var result scoreResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}
	if len(result.Data) == 0 {
		return 0, fmt.Errorf("vllm returned empty data")
	}
	return result.Data[0].Score, nil
}
