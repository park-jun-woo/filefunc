//ff:type feature=context type=model
//ff:what 컨텍스트 파이프라인 설정을 담는 구조체
package context

import "github.com/park-jun-woo/filefunc/internal/model"

// PipelineConfig holds context pipeline parameters.
type PipelineConfig struct {
	Prompt   string
	Search   string
	Depth    int
	WhatRate float64
	BodyRate float64
	Codebook *model.Codebook
	Generate func(string) (string, error)
}
