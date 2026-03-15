//ff:func feature=codebook type=loader
//ff:what codebook.yaml 파일을 로드하여 Codebook 구조체로 파싱
//ff:uses Codebook
//ff:uses Codebook
//ff:checked llm=gpt-oss:20b hash=2c47479f
package parse

import (
	"os"

	"github.com/park-jun-woo/filefunc/internal/model"
	"gopkg.in/yaml.v3"
)

// ParseCodebook loads and parses a codebook.yaml file.
func ParseCodebook(path string) (*model.Codebook, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cb model.Codebook
	if err := yaml.Unmarshal(data, &cb); err != nil {
		return nil, err
	}
	return &cb, nil
}
