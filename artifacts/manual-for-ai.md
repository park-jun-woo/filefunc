# filefunc — Manual for AI Agents

This document is the reference manual for AI code agents (Claude Code, etc.) when writing code in filefunc projects.

---

## How to Navigate

### 1. Read codebook.yaml first

The `codebook.yaml` at the project root is the project's vocabulary. Check feature and type values, then compose grep queries.

### 2. Find target files with grep

```bash
rg '//ff:func feature=validate'            # all funcs in validate feature
rg '//ff:func feature=validate type=rule'  # validate rules only
rg '//ff:type feature=parse'               # all types in parse feature
```

### 3. Read top annotations to narrow down

Each file's `//ff:what` tells you what it does without reading the body.

### 4. Full read only the files you need, then work

---

## Code Writing Rules

### File Creation

- **func file**: filename = func name (snake_case). `check_one_file_one_func.go` → `CheckOneFileOneFunc`
- **type file**: filename = type name (snake_case). `violation.go` → `Violation`
- **method file**: `{receiver}_{method}.go`. `server_start.go` → `(Server).Start`
- **One func/type/method per file.** Helper functions go in separate files.
- init() is only allowed alongside a var or func. Standalone init() files are forbidden.

### Annotations

Write annotations at the **very top** of every func/type file (above the package declaration).

```go
//ff:func feature=validate type=rule
//ff:what F1: validates one func per file
package validate
```

```go
//ff:type feature=validate type=model
//ff:what holds validation violation results
package model
```

| Annotation | Required | Description |
|---|---|---|
| `//ff:func` | Yes (func files) | Metadata: feature, type, etc. Values must exist in codebook.yaml |
| `//ff:type` | Yes (type files) | Metadata: feature, type, etc. Values must exist in codebook.yaml |
| `//ff:what` | Yes (func/type files) | One-line description. What does this do? |
| `//ff:why` | No | Why was it designed this way? Only record user decisions/requirements |
| `//ff:calls` | No | List of functions this func calls |
| `//ff:uses` | No | List of types this func uses |

### Code Quality

- **Nesting depth ≤ 2.** One loop + one branch is the maximum. Beyond that, extract a function or merge conditions.
- **Func recommended ≤ 100 lines, hard limit 1000 lines.**
- Early return pattern. Handle errors immediately and return.
- gofmt compliance.

### Naming

- Filenames: `snake_case` (`check_one_file_one_func.go`)
- Variables/functions: `camelCase` (`maxDepth`, `walkGoFiles`)
- Types: `PascalCase` (`GoFile`, `Violation`)

---

## Codebook Values (Current)

Allowed values defined in codebook.yaml. Only these values may be used in annotations.

```
feature: validate, annotate, chain, parse, codebook, report, cli
type:    command, rule, parser, walker, model, formatter, loader, util
pattern: error-collection, file-visitor, rule-registry
level:   ERROR, WARNING, INFO
```

If you need a new value, amend codebook.yaml.

---

## Validation

```bash
# Validate code
filefunc validate ./internal/

# Validate with codebook
filefunc validate --codebook codebook.yaml ./internal/

# JSON output
filefunc validate --format json ./internal/
```

Exit code 1 on violations. Zero violations required before committing.

---

## Exceptions (Not Violations)

- `_test.go` files may have multiple funcs (F5)
- Unexported types alongside a func in the same file are allowed (F6)
- Semantically grouped consts in one file are allowed (F7)
- const-only and var-only files do not require annotations

---

## Common Mistakes and Fixes

| Mistake | Cause | Fix |
|---|---|---|
| Two funcs in one file | Habit | Extract helper functions into separate files |
| for → switch → if (depth 3) | Nested branching | Use type assertions + early continue, or extract inner logic to a func |
| for → if → if (depth 3) | Nested conditions | Merge conditions `if a && b`, or use early break |
| Missing //ff:what | Forgot annotation | Write annotations first when creating a file |
| Value not in codebook | Didn't check vocabulary | Check codebook.yaml first. If absent, amend codebook |
