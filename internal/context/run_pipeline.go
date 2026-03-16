//ff:func feature=context type=command control=sequence
//ff:what 4단계 컨텍스트 파이프라인 오케스트레이터
package context

import (
	"fmt"
	"io"

	"github.com/park-jun-woo/filefunc/internal/chain"
	"github.com/park-jun-woo/filefunc/internal/model"
)

// PipelineConfig holds context pipeline parameters.
type PipelineConfig struct {
	Prompt   string
	Depth    int
	WhatRate float64
	BodyRate float64
	Generate func(string) (string, error)
}

// RunPipeline executes the 4-stage context pipeline.
func RunPipeline(w io.Writer, target string, g *chain.CallGraph, files []*model.GoFile, cfg PipelineConfig) error {
	fileMap := chain.BuildFuncFileMap(files)

	// 1단계: chain chon=2
	results := chain.TraverseChon(g, target, 2)
	fmt.Fprintf(w, "[1/4] chain chon=2: %s → %d funcs\n", target, len(results))
	if cfg.Depth <= 1 {
		FormatResult(w, results, chon1Scores(results), fileMap)
		return nil
	}

	// 2단계: same-feature 필터
	targetFeature := getFeature(target, fileMap)
	before := len(results)
	results = FilterFeature(results, targetFeature, fileMap)
	fmt.Fprintf(w, "[2/4] feature filter (%s): %d → %d\n", targetFeature, before, len(results))
	if cfg.Depth <= 2 {
		FormatResult(w, results, chon1Scores(results), fileMap)
		return nil
	}

	// 3단계: what 스코어링
	before = len(results)
	kept, scores, removed, err := ScoreWhat(results, cfg.Prompt, cfg.WhatRate, fileMap, cfg.Generate)
	if err != nil {
		return fmt.Errorf("what scoring failed: %w", err)
	}
	fmt.Fprintf(w, "[3/4] what scoring: %d → %d (rate≥%.1f, %d removed)\n", before, len(kept), cfg.WhatRate, removed)
	results = kept
	if cfg.Depth <= 3 {
		FormatResult(w, results, scores, fileMap)
		return nil
	}

	// 4단계: 본문 스코어링
	before = len(results)
	kept, scores, removed, err = ScoreBody(results, cfg.Prompt, cfg.BodyRate, fileMap, cfg.Generate)
	if err != nil {
		return fmt.Errorf("body scoring failed: %w", err)
	}
	fmt.Fprintf(w, "[4/4] body scoring: %d → %d (rate≥%.1f, %d removed)\n", before, len(kept), cfg.BodyRate, removed)

	FormatResult(w, kept, scores, fileMap)
	return nil
}
