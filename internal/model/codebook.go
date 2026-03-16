//ff:type feature=codebook type=model
//ff:what codebook.yaml의 허용 값 + 설명을 담는 구조체
package model

// Codebook holds the allowed values and descriptions for annotations.
// Each key maps to a map of value → description.
type Codebook struct {
	Required map[string]map[string]string `yaml:"required"`
	Optional map[string]map[string]string `yaml:"optional"`
}
