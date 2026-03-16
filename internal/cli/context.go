//ff:func feature=cli type=command control=sequence
//ff:what context 서브커맨드 — LLM 기반 4단계 컨텍스트 탐색
package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	ffcontext "github.com/park-jun-woo/filefunc/internal/context"
	"github.com/spf13/cobra"
)

var contextCmd = &cobra.Command{
	Use:   "context <func-name> <prompt>",
	Short: "Find relevant code context using LLM scoring",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		prompt := args[1]
		root, _ := cmd.Flags().GetString("root")
		depth, _ := cmd.Flags().GetInt("depth")
		whatRate, _ := cmd.Flags().GetFloat64("what-rate")
		bodyRate, _ := cmd.Flags().GetFloat64("body-rate")
		modelName, _ := cmd.Flags().GetString("model")
		endpoint, _ := cmd.Flags().GetString("endpoint")

		if err := CheckProjectRoot(root); err != nil {
			return err
		}

		g, files, err := BuildGraph(root)
		if err != nil {
			return err
		}

		generate := func(p string) (string, error) {
			return ollamaGenerate(endpoint, modelName, p)
		}

		return ffcontext.RunPipeline(os.Stdout, target, g, files, ffcontext.PipelineConfig{
			Prompt:   prompt,
			Depth:    depth,
			WhatRate: whatRate,
			BodyRate: bodyRate,
			Generate: generate,
		})
	},
}

func ollamaGenerate(endpoint, model, prompt string) (string, error) {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"model": model, "prompt": prompt, "stream": false,
	})
	resp, err := http.Post(endpoint+"/api/generate", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("ollama not available at %s", endpoint)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned %d: %s", resp.StatusCode, string(body))
	}
	var result struct{ Response string }
	json.Unmarshal(body, &result)
	return result.Response, nil
}

func init() {
	contextCmd.Flags().String("root", ".", "project root")
	contextCmd.Flags().Int("depth", 4, "pipeline depth (1-4)")
	contextCmd.Flags().Float64("what-rate", 0.2, "what scoring threshold")
	contextCmd.Flags().Float64("body-rate", 0.5, "body scoring threshold")
	contextCmd.Flags().String("model", "gpt-oss:20b", "ollama model name")
	contextCmd.Flags().String("endpoint", "http://localhost:11434", "ollama endpoint")
	rootCmd.AddCommand(contextCmd)
}
