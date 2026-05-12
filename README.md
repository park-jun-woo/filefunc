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

Requires Go 1.18+.

## Commands

### validate — Check code structure rules

```bash
filefunc validate                    # current directory as project root
filefunc validate /path/to/project   # explicit project root
filefunc validate --format json
```

Project root must contain `go.mod` and `codebook.yaml`. Read-only. Exit code 1 on violations. Respects `.ffignore`. Powered by [toulmin](https://github.com/park-jun-woo/toulmin) argumentation engine — rules are generic functions with backing-based judgment criteria, exceptions are defeats in a graph.

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
filefunc chain func ParseFile --package funcspec     # limit to specific package
```

Real-time AST analysis. Respects `.ffignore`.

| Flag | Description | Default |
|---|---|---|
| `--root` | Project root | `.` |
| `--chon` | Relationship distance (1~3) | 1 |
| `--child-depth` | Trace calls only to this depth | — |
| `--parent-depth` | Trace callers only to this depth | — |
| `--meta` | Include annotation metadata (meta,what,why,checked,all) | — |
| `--package` | Limit to funcs in this Go package (chain func only) | — |
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

### context — LLM context search

```bash
filefunc context "nesting depth 검증 수정"                        # 4-stage pipeline
filefunc context "modify depth logic" --depth 2                    # feature filter only
filefunc context "depth 수정" --what-rate 0.3                      # adjust what threshold
filefunc context "depth 수정" --body-rate 0.5                      # adjust body threshold
filefunc context "depth 수정" --search "feature=validate"          # skip LLM, direct filter
filefunc context "cross 수정" --search "feature=crosscheck ssot=openapi"  # multi-key AND
```

4-stage pipeline: LLM feature selection → feature filter → what scoring (LLM) → body scoring (LLM). No func name needed. Requires ollama with gpt-oss:20b.

| Flag | Description | Default |
|---|---|---|
| `--depth` | Pipeline depth (1-4) | 4 |
| `--what-rate` | What scoring threshold | 0.2 |
| `--body-rate` | Body scoring threshold | 0.5 |
| `--model` | ollama model | `gpt-oss:20b` |
| `--endpoint` | ollama endpoint | `http://localhost:11434` |
| `--search` | Direct annotation filter (skip LLM feature selection) | — |

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

All rules defined in [`rulebook.md`](rulebook.md) (SSOT). Categories: P (project), F (file structure), Q (code quality), A (annotation), C (codebook), N (naming).

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
  feature:
    validate: "코드 구조 룰 검증 (F1,Q1,A1 등 정적 분석 룰)"
    parse: "소스 코드, 어노테이션, codebook 파싱"
  type:
    command: "cobra 명령 엔트리포인트"
    rule: "개별 검증 룰 구현"

optional:
  pattern:
    error-collection: "에러 수집 후 일괄 보고"
  level:
    error: ""
    warning: ""
```

Each value has a description (`key: "description"`). Descriptions are used by `filefunc context` for LLM feature selection. `required` keys must be present in every annotation (A8). `optional` keys are used when relevant.

Values not in the codebook trigger `A2 ERROR`. Each project has its own codebook. `codebook.yaml` is required.

Codebook format rules (C1-C4) in [`rulebook.md`](rulebook.md). Codebook is validated first. If codebook fails, code validation does not run.

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

## Case study

### typer — Python CLI framework (1155 tests, 0 failures)

[Refactored fork](https://github.com/park-jun-woo/typer) of [fastapi/typer](https://github.com/fastapi/typer), restructured to pass `filefunc validate --lang python` with zero violations.

| Metric | Original | Refactored |
|---|---|---|
| Source files | 16 | 197 |
| filefunc violations | 69 | 0 |
| pytest passed | 1155 | 1155 |
| pytest failed | 0 | 0 |

All public APIs, import paths, and runtime behavior are identical to the original. No performance regression (import +2% within noise, all other benchmarks identical). Verified by full pytest suite and exhaustive comparison.

### hono — TypeScript web framework (4419 tests, 0 new failures)

[Refactored fork](https://github.com/park-jun-woo/hono) of [honojs/hono](https://github.com/honojs/hono), restructured to pass `filefunc validate --lang typescript` with zero violations.

| Metric | Original | Refactored |
|---|---|---|
| Source files | 186 | 626 |
| filefunc violations | 397 | 0 |
| vitest passed | 4419 | 4419 |
| vitest failed | 4 | 4 (pre-existing) |

All import paths and runtime behavior identical to original. Verified by full vitest suite.

## Academic basis

- **"Lost in the Middle" (Stanford, 2024)** — Relevant info in the middle of context drops performance 30%+.
- **"Context Length Alone Hurts LLM Performance" (Amazon, 2025)** — Even blank tokens degrade performance (13.9~85%). Short focused context wins.
- **"Context Rot" (Chroma Research)** — Focused prompt > full prompt across all models.

Research proved "shorter context is better." filefunc is the missing tool that structurally splits code so only the relevant parts enter context.

## License

MIT License — see [LICENSE](LICENSE).
