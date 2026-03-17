# 30K Tokens for 10,000 Functions: Compact Symbol Maps as LLM Code Navigation Entry Points

## Abstract

LLM code agents need a global view of a codebase before navigating to specific files. Current approaches either dump entire project structures (expensive) or rely on vector search (lossy). We present **symbol maps** — a compact representation format that encodes all function names in a codebase using `[package]func1,func2,...` notation, achieving ~3 tokens per function.

A 10,000-function codebase fits in ~30K tokens — small enough to include in every LLM prompt. Combined with a call-graph traversal tool, this creates a two-stage navigation pattern: **map for entry point, chain for impact scope**. We evaluate this on three Go projects, measuring token efficiency and comparing against codebook-based search and full project dumps.

## 1. Introduction

### 1.1 The Cold Start Problem

When an AI code agent begins a task, it faces a cold start: which files are relevant? Current approaches:

| Approach | Token cost | Accuracy | Latency |
|---|---|---|---|
| Full project dump | O(total LOC) | Perfect | Expensive |
| RAG / vector search | O(query + k results) | Approximate | Index required |
| User-specified files | O(specified) | Depends on user | Manual |
| grep + read | O(matches * file_size) | Good for exact terms | Multiple rounds |

None provide a cheap, complete, accurate global view.

### 1.2 The Symbol Map Approach

A symbol map lists every function/type in the project, grouped by package, using minimal syntax:

```
# whyso/v1

## go
[codemap]BuildMap,FormatMap,NeedsUpdate,buildSections,...
[history]BuildHistories,BuildHistoriesIncremental,BuildIndex,...
[main]clearCache,formatChangeRow,getSessionsDir,main,...
[model]ContentBlocks,IsUserMessage,UserContent
[output]FormatJSON,FormatYAML,OutputPath,ReadYAML,...
[parser]DetectSessionsDir,ExtractChanges,ListSessions,...
```

This is whyso's complete symbol map: ~80 functions in ~200 tokens. The key trade-off: **no metadata per function** (no descriptions, no types, no relationships) — just names, grouped by package.

### 1.3 Why No Metadata

Adding metadata (description, input/output types) would increase tokens linearly:

| Format | Tokens per function | 10K functions |
|---|---|---|
| Name only (map) | ~3 | ~30K |
| Name + 1-line description | ~15 | ~150K |
| Name + types + description | ~30 | ~300K |
| Full function signature | ~50 | ~500K |

At ~30K tokens, the symbol map fits comfortably in any LLM context. At 150K+, it competes with actual code for context space. The design choice is deliberate: **function names are sufficient for LLM navigation because LLMs infer purpose from naming conventions**.

## 2. Design

### 2.1 Format Specification

```
# <tool>/v1

## <language>
[<package>]<func1>,<func2>,...,<funcN>
```

Rules:
- Functions within a package are comma-separated on one line
- Package names serve as group headers
- Language sections allow multi-language projects
- Version header enables format evolution

### 2.2 Generation

whyso uses tree-sitter to extract symbols from supported languages:

| Language | Extracted symbols |
|---|---|
| Go | Functions, methods, types |
| TypeScript/JavaScript | Functions, classes, exports |
| Python | Functions, classes |
| Rust | Functions, structs, traits |
| OpenAPI | Operation IDs, paths |
| SQL | Table names, function names |
| Rego | Rule names |
| SSaC | Service functions |

Generation is incremental: only re-parses files modified since last run (`mtime` comparison).

### 2.3 Token Efficiency

The compact format achieves high information density:

```
[codemap]BuildMap,FormatMap,NeedsUpdate,buildSections,collectCaptures,dedupe,...
```

Tokenization (cl100k_base):
- `[codemap]` → 3 tokens
- Each function name → 1-3 tokens (camelCase splits predictably)
- Commas → shared with adjacent tokens
- Average: ~3 tokens per function

### 2.4 Staleness Management

The map uses `mtime`-based caching:
- `whyso map` checks source file modification times against `.whyso/_map.md`
- Only re-parses changed files
- Force regeneration with `whyso map -f`
- Prints "up to date" when no changes detected

## 3. Two-Stage Navigation: Map + Chain

### 3.1 The Pattern

```
Stage 1: whyso map → "Where is everything?"
  LLM reads symbol map (~30K tokens)
  Identifies entry point function(s) by name

Stage 2: filefunc chain <func> → "What does it affect?"
  Call graph traversal from entry point
  Returns callers, callees, co-called functions with //ff:what descriptions
```

### 3.2 Complementary Roles

| Tool | Provides | Token cost | Scope |
|---|---|---|---|
| whyso map | Global index (all functions) | ~30K fixed | Entire project |
| filefunc chain | Local graph (call relationships) | ~500 per query | Function neighborhood |
| codebook | Domain vocabulary (features, types) | ~200 fixed | Classification |

The map answers "what exists?" The chain answers "what's connected?" The codebook answers "how is it organized?"

### 3.3 Example Workflow

User request: "Fix the bug in history merging"

```
1. LLM reads symbol map
   → Spots [history]BuildHistories,BuildHistoriesIncremental,Merge,...
   → Identifies Merge as likely target

2. filefunc chain func Merge --chon 2 --meta what
   → Merge (what="Merges overlapping history entries for the same file")
     1촌 calls: entryKey, lastEntry
     1촌 called-by: BuildHistories, BuildHistoriesIncremental
     2촌 co-called: processChanges, skipBySince

3. LLM reads: merge.go, entry_key.go, last_entry.go
   → 3 files, ~75 lines total, zero noise
```

Without the map, the agent would grep for "merge", potentially hitting false positives across packages. The map provides immediate, unambiguous function location.

## 4. Comparison with Alternatives

### 4.1 vs. Codebook-Based Search

filefunc's codebook + `filefunc context` pipeline:
```
Codebook (250 tokens) → LLM selects feature → grep feature=X → read what → read body
```

Symbol map:
```
Map (30K tokens) → LLM identifies function → chain → read body
```

| Aspect | Codebook | Symbol map |
|---|---|---|
| Token cost | 250 (codebook) + LLM calls | 30K (map) + 0 LLM calls for selection |
| Requires LLM for selection | Yes (feature selection) | No (LLM reads map directly) |
| Granularity | Feature level (10-50 files) | Function level (1 file) |
| Setup cost | Manual codebook design | Automatic (tree-sitter) |

The codebook is more token-efficient for large projects (250 vs 30K tokens) but requires manual vocabulary design. The symbol map is automatic but costs more tokens.

**Best combined**: Codebook for domain classification, symbol map for function-level identification. They operate at different granularity levels and don't compete.

### 4.2 vs. Full Project Dump

A full project dump includes all source code:

| Project | Full dump tokens | Symbol map tokens | Ratio |
|---|---|---|---|
| filefunc (3K LOC) | ~12,000 | ~500 | 24× |
| whyso (2.4K LOC) | ~10,000 | ~300 | 33× |
| fullend (32K LOC) | ~130,000 | ~4,000 | 32× |

The symbol map is ~30× more compact than a full dump, while retaining the complete function inventory.

### 4.3 vs. RAG / Vector Search

| Aspect | RAG | Symbol map |
|---|---|---|
| Infrastructure | Embedding model + vector DB | File on disk |
| Index size | O(chunks * dim) floats | O(functions) tokens |
| Query | Approximate nearest neighbor | Exact string lookup |
| Result type | Ranked list (may miss) | Complete inventory |
| Update cost | Re-embed changed chunks | Re-parse changed files |

RAG retrieves relevant chunks but may miss functions with non-obvious names. The symbol map provides the complete list — nothing is hidden.

## 5. Evaluation

### 5.1 Token Measurements

| Project | Functions | Types | Map tokens | Tokens/symbol |
|---|---|---|---|---|
| filefunc | 136 | 15 | ~500 | 3.3 |
| whyso | 84 | 12 | ~300 | 3.1 |
| fullend | 1,077 | 218 | ~4,000 | 3.1 |

Consistent ~3 tokens per symbol across projects of different sizes.

### 5.2 Projected Scaling

| Project size | Functions | Estimated map tokens | % of 200K context |
|---|---|---|---|
| Small (5K LOC) | 200 | 600 | 0.3% |
| Medium (30K LOC) | 1,200 | 3,600 | 1.8% |
| Large (100K LOC) | 5,000 | 15,000 | 7.5% |
| Very large (500K LOC) | 10,000 | 30,000 | 15% |
| Extreme (1M LOC) | 20,000 | 60,000 | 30% |

Up to 500K LOC (10K functions), the map consumes <15% of a 200K context window. Beyond that, hierarchical grouping or on-demand loading may be needed.

### 5.3 Name Informativeness

The symbol map relies on function names being informative. In projects following naming conventions (filefunc's snake_case files = camelCase functions):

```
[validate]CheckNestingDepth,CheckOneFileOneFunc,CheckCodebookValues,...
```

An LLM reading these names can:
- Identify `CheckNestingDepth` as relevant to "nesting depth" queries
- Infer that `validate` package contains validation rules
- Distinguish `CheckOneFileOneFunc` from `CheckOneFileOneType` without descriptions

In projects with poor naming (`doStuff`, `helper1`, `process`), the symbol map degrades. **Map quality = naming quality**.

## 6. Limitations

- Token efficiency depends on naming conventions. Poorly-named functions reduce navigation accuracy.
- No metadata means the LLM must infer function purpose from names alone. For ambiguous names, chain's `--meta what` provides descriptions as a second step.
- Scaling beyond 20K functions (60K tokens) may require hierarchical maps or on-demand loading.
- Currently measured on Go projects only. Other languages may have different tokens-per-symbol ratios.
- No formal accuracy evaluation of LLM navigation using symbol maps (planned future work).

## 7. Future Work

- **Accuracy benchmarks**: Measure LLM task completion accuracy with/without symbol map in context
- **Hierarchical maps**: For very large projects, package-level summaries with drill-down
- **Cross-project maps**: Symbol maps for dependency libraries (similar to filefunc's registry concept)
- **Incremental context**: Start with map, progressively add chain results and file contents as needed

## 8. Conclusion

A symbol map encoding `[package]func1,func2,...` achieves ~3 tokens per function, fitting 10,000 functions in 30K tokens. This provides LLM code agents with a complete, cheap, always-available global index of the codebase.

Combined with call-graph traversal (`filefunc chain`), the two-stage pattern — **map for entry, chain for scope** — enables precise navigation without vector search infrastructure, manual codebook design, or expensive project dumps.

The key insight: **function names are metadata**. In well-named codebases, a list of names is sufficient for an LLM to navigate. The symbol map makes this list compact enough to include in every prompt.

---

## Citation

```
@misc{compact-symbol-maps2026,
  title={30K Tokens for 10,000 Functions: Compact Symbol Maps as LLM Code Navigation Entry Points},
  author={Park, Jun-Woo},
  year={2026},
  url={https://github.com/park-jun-woo/whyso}
}
```
