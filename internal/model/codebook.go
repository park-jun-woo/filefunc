//ff:type feature=codebook type=model
//ff:what codebook.yaml의 허용 값 목록을 담는 구조체
package model

// Codebook holds the allowed values for //ff:func annotations.
type Codebook struct {
	Feature []string `yaml:"feature"`
	Type    []string `yaml:"type"`
	Pattern []string `yaml:"pattern"`
	Level   []string `yaml:"level"`
}
