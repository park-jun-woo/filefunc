# filefunc

**One file, one concept.** The filename is the concept name.

A code structure convention and CLI toolchain for LLM-native Go application development — backend services, CLI tools, code generators, and SSOT validators. Not intended for algorithm libraries, low-level systems programming, or performance-critical hot paths.

## Why

AI code agents (Claude Code, etc.) navigate code via `grep → read`. The unit of `read` is a file. If file = concept, then every read returns exactly one relevant thing.

```
# Without filefunc
read utils.go → 20 funcs, 19 irrelevant. Context pollution.

# With filefunc
read check_one_file_one_func.go → 1 func. Exactly what you needed.
```

**It's more important to NOT open 290 irrelevant functions than to pick the right 5-10.**

The primary citizen of filefunc is the AI agent, not the human. File count explosion is a feature, not a bug — more files means smaller files, less noise per read. Human convenience is solved at the view layer (VSCode extensions, etc.).

## Install

```bash
go install github.com/park-jun-woo/filefunc/cmd/filefunc@latest
```

Or build from source:

```bash
git clone https://github.com/park-jun-woo/filefunc.git
cd filefunc
go build ./cmd/filefunc/
```

Requires Go 1.22+.

## Commands

### validate — Check code structure rules

```bash
filefunc validate                    # current directory as project root
filefunc validate /path/to/project   # explicit project root
filefunc validate --format json
```

Project root must contain `go.mod` and `codebook.yaml`. Read-only. Exit code 1 on violations. Respects `.ffignore`.

### chain — Trace call relationships

```bash
filefunc chain func RunAll              # 1촌 (default, current dir)
filefunc chain func RunAll --chon 2     # 2촌 (co-called included)
filefunc chain func RunAll --chon 3     # 3촌 (max)
filefunc chain func RunAll --child-depth 3   # calls only
filefunc chain func RunAll --parent-depth 3  # callers only
filefunc chain feature validate         # all funcs in feature
filefunc chain func RunAll --root /path/to/project  # explicit project root
filefunc chain func RunAll --chon 2 --meta what     # with //ff:what annotations
filefunc chain func RunAll --chon 2 --meta all      # with all annotations
filefunc chain func RunAll --chon 2 --meta what \
  --prompt "nesting depth 수정" --rate 0.8           # reranker filtering
```

Real-time AST analysis. Respects `.ffignore`.

| Flag | Description | Default |
|---|---|---|
| `--root` | Project root | `.` |
| `--chon` | Relationship distance (1~3) | 1 |
| `--child-depth` | Trace calls only to this depth | — |
| `--parent-depth` | Trace callers only to this depth | — |
| `--meta` | Include annotation metadata (meta,what,why,checked,all) | — |
| `--prompt` | User task intent for relevance scoring (requires vLLM) | — |
| `--rate` | Relevance score threshold (0.0~1.0) | 0.8 |
| `--model` | Reranker model name | `Qwen/Qwen3-Reranker-0.6B` |
| `--score-endpoint` | vLLM endpoint for reranker | `http://localhost:8000` |

`--prompt` requires a vLLM server running Qwen3-Reranker-0.6B:

```bash
pip install vllm
vllm serve Qwen/Qwen3-Reranker-0.6B --task score \
  --hf_overrides '{"architectures":["Qwen3ForSequenceClassification"],"classifier_from_token":["no","yes"],"is_original_qwen3_reranker":true}'
```

### llmc — LLM verification

```bash
filefunc llmc                           # current directory
filefunc llmc /path/to/project          # explicit project root
filefunc llmc --model qwen3:8b
filefunc llmc --threshold 0.9
```

Verifies `//ff:what` matches func body using local LLM (ollama). Scores 0.0~1.0, threshold 0.8. On pass, writes `//ff:checked` signature. Respects `.ffignore`.

| Flag | Description | Default |
|---|---|---|
| `--provider` | LLM provider | `ollama` |
| `--model` | Model name | `gpt-oss:20b` |
| `--endpoint` | API endpoint | `http://localhost:11434` |
| `--threshold` | Minimum passing score | `0.8` |

## Rules

### File structure

| Rule | Description | Severity |
|---|---|---|
| F1 | One func per file (filename = func name) | ERROR |
| F2 | One type per file (filename = type name) | ERROR |
| F3 | One method per file | ERROR |
| F4 | init() must not exist alone (requires var or func) | ERROR |
| F5 | _test.go files may have multiple funcs | exception |
| F6 | Semantically grouped consts allowed in one file | exception |

### Code quality

| Rule | Description | Severity |
|---|---|---|
| Q1 | Nesting depth: sequence=2, selection=2, iteration=dimension+1 | ERROR |
| Q2 | Func max 1000 lines | ERROR |
| Q3 | Func recommended max: sequence/iteration 100, selection 300 | WARNING |

### Annotation

| Rule | Description | Severity |
|---|---|---|
| A1 | Func files require `//ff:func`, type files require `//ff:type` | ERROR |
| A2 | Annotation values must exist in codebook | ERROR |
| A3 | Func/type files require `//ff:what` | ERROR |
| A6 | Annotations must be at the top of the file | ERROR |
| A7 | `//ff:checked` hash mismatch (body changed after LLM verification) | ERROR |
| A8 | Required codebook keys must be present in annotation | ERROR |
| A9 | Func files must have `control=` (sequence/selection/iteration) | ERROR |
| A10 | `control=selection` but no switch at depth 1 | ERROR |
| A11 | `control=iteration` but no loop at depth 1 | ERROR |
| A12 | `control=sequence` but switch/loop exists at depth 1 | ERROR |
| A13 | `control=selection` but loop exists at depth 1 | ERROR |
| A14 | `control=iteration` but switch exists at depth 1 | ERROR |
| A15 | `control=iteration` requires `dimension=` | ERROR |
| A16 | `dimension=` value must be a positive integer | ERROR |

## Annotations

```go
//ff:func feature=validate type=rule control=sequence
//ff:what F1: validates one func per file
//ff:why Primary citizen is AI agent. 1 file 1 concept prevents context pollution.
//ff:checked llm=gpt-oss:20b hash=a3f8c1d2     (auto-generated by llmc)
func CheckOneFileOneFunc(gf *model.GoFile) []model.Violation {
```

`control=` is required for all func files (A9). Values: `sequence`, `selection` (switch), `iteration` (loop). Based on Böhm-Jacopini theorem (1966). 1 func 1 control.

`dimension=` is required for `control=iteration` files (A15). Specifies the dimensionality of the data being iterated. Q1 depth limit = dimension + 1. dimension=1 for flat lists (depth ≤ 2), dimension ≥ 2 requires named type (struct/interface) nesting.

| Annotation | Purpose | Required |
|---|---|---|
| `//ff:func` | Func metadata (feature, type, etc.) | Yes (func files) |
| `//ff:type` | Type metadata (feature, type, etc.) | Yes (type files) |
| `//ff:what` | One-line description — what it does | Yes |
| `//ff:why` | Design decision — why it's this way | No |
| `//ff:checked` | LLM verification signature | Auto (`filefunc llmc`) |

## Codebook

The codebook defines allowed values for annotations. It's the project's vocabulary — the map that makes `grep` precise.

```yaml
# codebook.yaml
required:
  feature: [validate, annotate, chain, parse, codebook, report, cli]
  type: [command, rule, parser, walker, model, formatter, loader, util]

optional:
  pattern: [error-collection, file-visitor, rule-registry]
  level: [error, warning, info]
```

`required` keys must be present in every `//ff:func` and `//ff:type` annotation (A8). This guarantees grep reliability — required keys never have gaps. `optional` keys are used only when relevant.

Values not in the codebook trigger `A2 ERROR`. Each project has its own codebook. `codebook.yaml` is required — validate will error without it.

### Codebook format rules

| Rule | Description | Severity |
|---|---|---|
| C1 | `required` section must have at least one key with at least one value | ERROR |
| C2 | No duplicate values within the same key | ERROR |
| C3 | All values must be lowercase + hyphens only (`[a-z][a-z0-9-]*`) | ERROR |

Codebook is validated first. If codebook fails, code validation does not run.

## .ffignore

Exclude paths from all filefunc commands. Place `.ffignore` in the project root (next to `go.mod`). Same syntax as `.gitignore`.

```
# Example .ffignore
vendor/
*.pb.go
*_gen.go
internal/legacy/
```

Optional. If absent, nothing is excluded.

## Academic basis

- **"Lost in the Middle" (Stanford, 2024)** — Relevant info in the middle of context drops performance 30%+.
- **"Context Length Alone Hurts LLM Performance" (Amazon, 2025)** — Even blank tokens degrade performance (13.9~85%). Short focused context wins.
- **"Context Rot" (Chroma Research)** — Focused prompt > full prompt across all models.

Research proved "shorter context is better." filefunc is the missing tool that structurally splits code so only the relevant parts enter context.

## License

MIT License — see [LICENSE](LICENSE).
