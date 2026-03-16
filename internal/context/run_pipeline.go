//ff:func feature=context type=command control=sequence
//ff:what 4단계 컨텍스트 파이프라인 오케스트레이터
package context

import (
	"fmt"
	"io"
	"strings"

	"github.com/park-jun-woo/filefunc/internal/model"
)

// PipelineConfig holds context pipeline parameters.
type PipelineConfig struct {
	Prompt      string
	Search      string
	Depth       int
	WhatRate    float64
	BodyRate    float64
	CodebookRaw string
	Generate    func(string) (string, error)
}

// RunPipeline executes the 4-stage context pipeline.
func RunPipeline(w io.Writer, files []*model.GoFile, cfg PipelineConfig) error {
	var filtered []*model.GoFile

	if cfg.Search != "" {
		// --search: LLM 스킵, 직접 필터
		query := ParseSearch(cfg.Search)
		filtered = FilterSearch(files, query)
		fmt.Fprintf(w, "[1/4] search filter: %s → %d files\n", cfg.Search, len(filtered))
		fmt.Fprintln(w, "[2/4] (skipped — direct search)")
	} else {
		// LLM feature 선택
		features, err := SelectFeature(cfg.Prompt, cfg.CodebookRaw, cfg.Generate)
		if err != nil {
			return fmt.Errorf("feature selection failed: %w", err)
		}
		if len(features) == 0 {
			fmt.Fprintln(w, "[1/4] feature selection: (none)")
			fmt.Fprintln(w, "\nResults:\n  (no results)")
			return nil
		}
		fmt.Fprintf(w, "[1/4] feature selection: %s (LLM)\n", strings.Join(features, ", "))
		filtered = FilterFeature(files, features)
		fmt.Fprintf(w, "[2/4] feature filter: %d → %d\n", len(files), len(filtered))
	}

	if cfg.Depth <= 2 {
		FormatResult(w, filtered, nil)
		return nil
	}

	// 3단계: what 스코어링
	before := len(filtered)
	kept, scores, removed, err := ScoreWhat(filtered, cfg.Prompt, cfg.WhatRate, cfg.Generate)
	if err != nil {
		return fmt.Errorf("what scoring failed: %w", err)
	}
	fmt.Fprintf(w, "[3/4] what scoring: %d → %d (rate≥%.1f, %d removed)\n", before, len(kept), cfg.WhatRate, removed)
	if cfg.Depth <= 3 {
		FormatResult(w, kept, scores)
		return nil
	}

	// 4단계: 본문 스코어링
	before = len(kept)
	kept, scores, removed, err = ScoreBody(kept, cfg.Prompt, cfg.BodyRate, cfg.Generate)
	if err != nil {
		return fmt.Errorf("body scoring failed: %w", err)
	}
	fmt.Fprintf(w, "[4/4] body scoring: %d → %d (rate≥%.1f, %d removed)\n", before, len(kept), cfg.BodyRate, removed)

	FormatResult(w, kept, scores)
	return nil
}
