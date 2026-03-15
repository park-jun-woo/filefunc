//ff:type feature=codebook type=model
//ff:what codebook.yaml의 허용 값 목록을 담는 구조체
package model

// Codebook holds the allowed values for //ff:func and //ff:type annotations.
type Codebook struct {
	Required map[string][]string `yaml:"required"`
	Optional map[string][]string `yaml:"optional"`
}
