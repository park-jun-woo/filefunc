# filefunc

**One file, one concept.** The filename is the concept name.

A code structure convention and CLI toolchain for LLM-native Go development.

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

## Usage

```bash
# Validate current directory
filefunc validate

# Validate specific path
filefunc validate ./internal/

# With codebook validation
filefunc validate --codebook codebook.yaml ./internal/

# JSON output
filefunc validate --format json ./internal/
```

## Rules

### File structure

| Rule | Description | Severity |
|---|---|---|
| F1 | One func per file (filename = func name) | ERROR |
| F2 | One type per file (filename = type name) | ERROR |
| F3 | One method per file | ERROR |
| F4 | init() must not exist alone (requires var or func) | ERROR |
| F5 | _test.go files may have multiple funcs | exception |
| F6 | Unexported param types allowed alongside func | exception |
| F7 | Semantically grouped consts allowed in one file | exception |

### Code quality

| Rule | Description | Severity |
|---|---|---|
| Q1 | Nesting depth ≤ 2 | ERROR |
| Q2 | Func max 1000 lines | ERROR |
| Q3 | Func recommended max 100 lines | WARNING |

### Annotation

| Rule | Description | Severity |
|---|---|---|
| A1 | Func files require `//ff:func`, type files require `//ff:type` | ERROR |
| A2 | Annotation values must exist in codebook | ERROR |
| A3 | Func/type files require `//ff:what` | ERROR |
| A6 | Annotations must be at the top of the file | ERROR |

## Annotations

```go
//ff:func feature=validate type=rule
//ff:what F1: validates one func per file
//ff:why Primary citizen is AI agent. 1 file 1 concept prevents context pollution.
//ff:calls IsConstOnly
//ff:uses GoFile, Violation
func CheckOneFileOneFunc(gf *model.GoFile) []model.Violation {
```

| Annotation | Purpose | Required |
|---|---|---|
| `//ff:func` | Func metadata (feature, type, etc.) | Yes (func files) |
| `//ff:type` | Type metadata (feature, type, etc.) | Yes (type files) |
| `//ff:what` | One-line description — what it does | Yes |
| `//ff:why` | Design decision — why it's this way | No |
| `//ff:calls` | Functions this func calls | No |
| `//ff:uses` | Types this func uses | No |

## Codebook

The codebook defines allowed values for annotations. It's the project's vocabulary — the map that makes `grep` precise.

```yaml
# codebook.yaml
feature: [validate, parse, codebook, report, cli]
type: [command, rule, parser, walker, model, formatter, loader, util]
pattern: [error-collection, file-visitor, rule-registry]
level: [ERROR, WARNING, INFO]
```

Values not in the codebook trigger `A2 ERROR`.

## Options

| Flag | Description | Default |
|---|---|---|
| `--codebook` | Path to codebook.yaml | (none) |
| `--format` | Output format (text / json) | text |

## Academic basis

- **"Lost in the Middle" (Stanford, 2024)** — Relevant info in the middle of context drops performance 30%+.
- **"Context Length Alone Hurts LLM Performance" (Amazon, 2025)** — Even blank tokens degrade performance (13.9~85%). Short focused context wins.
- **"Context Rot" (Chroma Research)** — Focused prompt > full prompt across all models.

Research proved "shorter context is better." filefunc is the missing tool that structurally splits code so only the relevant parts enter context.

## License

MIT License — see [LICENSE](LICENSE).
