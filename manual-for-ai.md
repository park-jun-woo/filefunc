# filefunc — Manual for AI Agents

For Go application-layer projects: backend services, CLI tools, code generators, SSOT validators.

## How to Navigate

1. Read `codebook.yaml` — project vocabulary (required/optional keys and allowed values)
2. `filefunc chain func <target> --chon 2` — trace call relationships before modifying
3. `rg '//ff:func feature=validate'` — grep with codebook values to find files
4. Read `//ff:what` to narrow down — skip body if what is sufficient
5. Full read only the files you need, then work

---

## Rules

All rules defined in `rulebook.md` (SSOT). Categories: P (project), F (file structure), Q (code quality), A (annotation), C (codebook), N (naming).

---

## Annotations

Write at the **very top** of every func/type file (above package declaration):

```go
//ff:func feature=validate type=rule control=sequence
//ff:what F1: validates one func per file
//ff:checked llm=gpt-oss:20b hash=a3f8c1d2    (auto by llmc)
package validate
```

`control=` is required for all func files (A9). Values: `sequence`, `selection` (switch), `iteration` (loop). Böhm-Jacopini (1966). 1 func 1 control — no mixing.

`dimension=` is required for `control=iteration` (A15). Q1 depth limit = dimension + 1. dimension=1 for flat lists (depth ≤ 2). dimension ≥ 2 requires named type (struct/interface) nesting — raw `[][][]int` is not allowed.

| Annotation | Required | Description |
|---|---|---|
| `//ff:func` | func files | Metadata (feature, type, control). Values from codebook.yaml + control rule |
| `//ff:type` | type files | Metadata (feature, type). Values from codebook.yaml |
| `//ff:what` | func/type files | One-line description. What does this do? |
| `//ff:why` | optional | Why designed this way? User decisions only |
| `//ff:checked` | auto (llmc) | LLM verification signature. Do not write manually |

### Control-based read strategy

```
control=selection  → read entire body at once. Don't read cases partially.
control=iteration  → focus on loop body. Outside loop is initialization.
control=sequence   → read only the step you need. Other steps: what is enough.
```

### Naming

- Filenames: `snake_case`
- Variables/functions: `camelCase`
- Types: `PascalCase`
- gofmt compliance, early return pattern

---

## Codebook

`codebook.yaml` must exist in the project root (next to `go.mod`). `required` keys must be in every annotation (A8). `optional` keys are used when relevant.

```yaml
required:
  feature:
    validate: "code structure rule validation (F1,Q1,A1 etc.)"
    parse: "source code, annotation, codebook parsing"
  type:
    command: "cobra command entrypoint"
    rule: "individual validation rule"

optional:
  pattern:
    error-collection: "collect errors for batch reporting"
  level:
    error: ""
```

Each value has a description (`key: "description"`). Used by `filefunc context` for LLM feature selection.

Amend codebook.yaml when new values are needed.

---

## Commands

```bash
filefunc validate                                    # current dir as project root
filefunc validate /path/to/project                   # explicit project root
filefunc validate --format json
filefunc chain func RunAll --chon 2                  # call relationships
filefunc chain func RunAll --chon 2 --meta what      # with //ff:what annotations
filefunc chain func RunAll --chon 2 --meta all       # with all annotations
filefunc chain func RunAll --chon 2 --meta what \
  --prompt "nesting depth 수정" --rate 0.8            # reranker filtering
filefunc chain func ParseFile --package funcspec     # limit to specific package
filefunc chain feature validate                      # feature-wide chain
filefunc chain func RunAll --root /path/to/project   # explicit project root
filefunc context "nesting depth 수정"                   # LLM 4-stage context search
filefunc llmc                                        # LLM what-body verification
filefunc llmc /path/to/project
filefunc llmc --model qwen3:8b --threshold 0.9
```

Project root must contain `go.mod` and `codebook.yaml`. Omit to use current directory.

`--prompt` requires vLLM server: `pip install vllm && vllm serve Qwen/Qwen3-Reranker-0.6B --task score --hf_overrides '{"architectures":["Qwen3ForSequenceClassification"],"classifier_from_token":["no","yes"],"is_original_qwen3_reranker":true}'`

Exit code 1 on violations. Zero violations required before committing.

### .ffignore

Place in project root. Same syntax as `.gitignore`. Excludes paths from all commands.

```
vendor/
*.pb.go
*_gen.go
```

---

## Common Mistakes

| Mistake | Fix |
|---|---|
| Two funcs in one file | Extract helper functions into separate files |
| depth 3 (for→switch→if, for→if→if) | Type assertions + early continue, merge conditions, or extract func |
| Missing //ff:what | Write annotations first when creating a file |
| Value not in codebook | Check codebook.yaml first. Amend if absent |
| //ff:checked hash mismatch | Run `filefunc llmc` to re-verify |
| "codebook.yaml required" | Create codebook.yaml next to go.mod |
