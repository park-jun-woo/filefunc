//ff:func feature=cli type=command control=sequence
//ff:what validate 서브커맨드 정의 및 코드 구조 룰 검증 실행
package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/filefunc/internal/model"
	"github.com/park-jun-woo/filefunc/internal/parse"
	"github.com/park-jun-woo/filefunc/internal/report"
	"github.com/park-jun-woo/filefunc/internal/validate"
	"github.com/park-jun-woo/filefunc/internal/walk"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [project-root]",
	Short: "Validate code structure rules",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root := "."
		if len(args) > 0 {
			root = args[0]
		}

		if err := CheckProjectRoot(root); err != nil {
			return err
		}

		codebookPath, _ := cmd.Flags().GetString("codebook")
		format, _ := cmd.Flags().GetString("format")

		if codebookPath == "" {
			codebookPath = filepath.Join(root, "codebook.yaml")
		}

		cb, err := parse.ParseCodebook(codebookPath)
		if err != nil {
			return fmt.Errorf("codebook.yaml required: %w", err)
		}

		cbViolations := validate.ValidateCodebook(cb, codebookPath)
		if len(cbViolations) > 0 {
			report.FormatText(os.Stdout, cbViolations)
			return fmt.Errorf("codebook.yaml has %d violation(s) — fix before validating code", len(cbViolations))
		}

		ignorePatterns := walk.ParseFFIgnore(filepath.Join(root, ".ffignore"))
		paths, err := walk.WalkGoFiles(root, ignorePatterns)
		if err != nil {
			return fmt.Errorf("walking files: %w", err)
		}

		var files []*model.GoFile
		for _, p := range paths {
			gf, err := parse.ParseGoFile(p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: skipping %s: %v\n", p, err)
				continue
			}
			files = append(files, gf)
		}

		violations := validate.RunAll(files, cb)

		if format == "json" {
			report.FormatJSON(os.Stdout, violations)
		} else {
			report.FormatText(os.Stdout, violations)
		}

		if len(violations) > 0 {
			return fmt.Errorf("%d violation(s) found", len(violations))
		}
		return nil
	},
}

func init() {
	validateCmd.Flags().String("codebook", "", "path to codebook.yaml (default: <project-root>/codebook.yaml)")
	validateCmd.Flags().String("format", "text", "output format (text or json)")
	rootCmd.AddCommand(validateCmd)
}
