//ff:func feature=cli type=loader control=iteration dimension=1
//ff:what ollama 엔드포인트에 모델이 존재하는지 API로 확인
package llm

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// ModelExists checks if the given model is available on the ollama endpoint.
func ModelExists(endpoint, model string) bool {
	resp, err := http.Get(endpoint + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	var tags ollamaTagsResponse
	if err := json.Unmarshal(body, &tags); err != nil {
		return false
	}

	for _, m := range tags.Models {
		if m.Name == model || strings.HasPrefix(m.Name, model+":") {
			return true
		}
	}
	return false
}
