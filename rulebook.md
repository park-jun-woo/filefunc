# filefunc rulebook

Single source of truth (SSOT) for all filefunc rules. All documents (CLAUDE.md, manual-for-ai.md) and code (validate_graph*.go) must stay in sync with this document.

---

## P: Project rules

Project-level structural rules. Checked before file-level validation.

| # | Rule | Severity | Verification |
|---|---|---|---|
| P1 | Single language per project. `.go`, `.py`, `.ts`/`.tsx` must not coexist | ERROR | WalkFiles for other language extensions (after .ffignore) |

---

## I: Import rules

Import structure rules. Python and TypeScript only (Go has no intra-package imports).

| # | Rule | Severity | Verification |
|---|---|---|---|
| I1 | No circular imports in project | ERROR | Build import graph from AST, detect cycles (DFS). Error message includes cycle path and fix advice: which import to move inside function body |

### I1 error format

```
[ERROR] I1: circular import — A → B → C → A
  fix: move "from .A import func" inside function body in C.py
```

The fix advice targets the last edge in the cycle — moving that one import to a lazy (function-body) import breaks the cycle.

---

## F: File structure rules

| # | Rule | Severity | Go | Python | TypeScript |
|---|---|---|---|---|---|
| F1 | One func per file (filename = func name). Includes test files | ERROR | includes `_test.go` | includes `test_*.py` | includes `*.test.ts` |
| F2 | One type per file (filename = type name) | ERROR | type | class | class/interface/type |
| F3 | One method per file | ERROR | `receiver_method.go` | **exempt** (F2 sufficient) | **exempt** (F2 sufficient) |
| F4 | `init()` must not exist alone (requires var or func) | ERROR | Go only | N/A | N/A |
| F5 | Semantically grouped constants allowed in one file | exception | const | module-level UPPER_CASE variables | export const |

### F3 Python Mixin pattern

```python
# server_start.py — 1 class, 1 method → # ff:func
class ServerStartMixin:
    def start(self): ...

# server.py — composition class, __init__ allowed → # ff:type
class Server(ServerStartMixin, ServerStopMixin):
    def __init__(self, config): ...
```

---

## Q: Code quality rules

| # | Rule | Severity | Verification |
|---|---|---|---|
| Q1 | Nesting depth: sequence=2, selection=2, iteration=dimension+1 | ERROR | MaxDepth vs limit |
| Q2 | Func body max 1000 lines (signature excluded) | ERROR | line count |
| Q3 | Sequence func body max 100 lines (signature excluded) | ERROR | line count + control=sequence |
| Q4 | Control body PURE > 10 lines → extract to sequence func | ERROR | body lines - inner control lines |

### Q1 depth target statements

| Category | Go | Python |
|---|---|---|
| Branching | if, switch, type switch | if, match |
| Looping | for, range | for, while |
| Excluded | select | with, try/except, comprehension |

- `elif` / `else if` does not increment depth

### Q4 PURE calculation

PURE = total control body lines - inner control statement lines (2-depth exemption).
For switch/match, Q4 applies per case clause individually.

---

## A: Annotation rules

| # | Rule | Severity | Verification |
|---|---|---|---|
| A1 | Func files require `//ff:func` (Go) / `# ff:func` (Python); type files require `//ff:type` / `# ff:type` | ERROR | annotation presence |
| A2 | Annotation values must exist in codebook | ERROR | codebook yaml lookup |
| A3 | Func/type files require `//ff:what` / `# ff:what` | ERROR | what presence |
| A6 | Annotations must be at the top of the file | ERROR | position check |
| A7 | `//ff:checked` / `# ff:checked` hash mismatch — signature broken | ERROR | hash comparison. LLM verification via `filefunc llmc` |
| A8 | All required codebook keys must be present in annotation | ERROR | codebook required lookup |
| A9 | Func files must have `control=` (sequence/selection/iteration) | ERROR | control presence + value check |
| A10 | `control=selection` but no switch(Go)/match(Python) at depth 1 | ERROR | AST verification |
| A11 | `control=iteration` but no loop at depth 1 | ERROR | AST verification |
| A12 | `control=sequence` but switch/match or loop exists at depth 1 | ERROR | AST verification |
| A13 | `control=selection` but loop exists at depth 1 | ERROR | AST verification |
| A14 | `control=iteration` but switch(Go)/match(Python) exists at depth 1 | ERROR | AST verification |
| A15 | `control=iteration` requires `dimension=` | ERROR | annotation presence |
| A16 | `dimension=` value must be a positive integer | ERROR | value parsing |

### Annotation format

```
Go:     //ff:func feature=validate type=rule control=sequence
Python: # ff:func feature=validate type=rule control=sequence
```

| Annotation | Required | Description |
|---|---|---|
| `//ff:func` / `# ff:func` | func files | Codebook values: feature, type, control, etc. |
| `//ff:type` / `# ff:type` | type files | Codebook values: feature, type, etc. |
| `//ff:what` / `# ff:what` | func/type files | What does this func/type do? |
| `//ff:why` / `# ff:why` | optional | Why designed this way? (user decision rationale) |
| `//ff:checked` / `# ff:checked` | auto (llmc) | LLM verification signature. Do not write manually |

### A6 top-of-file rule

| Go | Python |
|---|---|
| Before package declaration | shebang → encoding → `# ff:` → blank lines → import/code |

---

## C: Codebook rules

| # | Rule | Severity |
|---|---|---|
| C1 | `required` section must have at least one key with at least one value | ERROR |
| C2 | No duplicate keys within the same section | ERROR |
| C3 | Keys must be lowercase + hyphens only (`[a-z][a-z0-9-]*`) | ERROR |
| C4 | Required values should have non-empty descriptions | WARNING |

Codebook is validated before file-level rules. Codebook violations block code validation.

---

## N: Naming rules

| # | Rule |
|---|---|
| N1 | Filenames: snake_case |
| N2 | Variables/functions: camelCase (Go), snake_case (Python) |
| N3 | Types: PascalCase (Go), PascalCase (Python class) |
| N4 | gofmt compliance (Go) |
| N5 | Handle errors immediately (early return) |

---

## dimension (iteration only)

Dimensionality of the data being iterated. Q1 depth limit = dimension + 1.

| dimension | Meaning | Depth limit |
|---|---|---|
| 1 | Flat list traversal | 2 |
| 2 | 2D data traversal | 3 |
| N | N-dimensional | N + 1 |

dimension >= 2 requires named type nesting (Go: struct/interface, Python: class).

---

## Exceptions

- Const/constant-only files do not require annotations (F5 defeater)
- If no `//ff:checked` / `# ff:checked` exists in the project, A7 is skipped entirely
- Files without funcs are exempt from A9-A16 (HasNoFunc defeater)

---

## Validation order

```
1. P rules (project level) — mixed language check
2. C rules (codebook) — codebook.yaml integrity
3. F/Q/A rules (file level) — toulmin defeats graph
```

P or C violations block subsequent validation.
