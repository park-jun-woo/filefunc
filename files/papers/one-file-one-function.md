# One File, One Function: How File Granularity Affects AI Code Agent Performance

## Abstract

AI code agents read code file-by-file. When a file contains 20 functions, reading one function forces the agent to ingest 19 irrelevant functions into its context window — we call this **context contamination**. We present empirical evidence from refactoring two Go projects (21,232 LOC and 1,773 LOC) to a strict "one file, one function" convention, measuring the impact on code structure metrics.

After conversion, average file size dropped from 244 to 25.4 lines (-89.6%), single-responsibility violations dropped from 75.9% to 0%, and nesting depth violations dropped from 148+ to 0. We argue that file granularity is the primary lever for AI code agent efficiency: smaller files mean less noise per read, and file-level metadata (name, annotations) becomes function-level metadata for free.

## 1. Introduction

### 1.1 The Read Unit Problem

AI code agents (Claude Code, Cursor, Copilot) navigate codebases through `read file` operations. The atomic unit of code retrieval is the file. This creates a fundamental mismatch:

- **Developer intent**: "I need to understand function X"
- **Agent action**: "I must read the entire file containing function X"
- **Result**: Function X is 30 lines, but the file is 500 lines

The 470 irrelevant lines are not free. Research demonstrates that:

- **Lost in the Middle** (Liu et al., 2024): Relevant information in the middle of long contexts drops LLM performance by 30%+
- **Context Length Alone Hurts** (Amazon, 2025): Even blank padding tokens degrade performance by 13.9-85%
- **Context Rot** (Chroma Research): Focused prompts outperform full prompts across all models tested

### 1.2 The File Granularity Thesis

If the read unit is the file, then **file granularity determines context quality**. We propose:

> Minimize the gap between "what the agent needs" and "what the file contains" by making each file contain exactly one concept.

This is the 1 file = 1 function (1F1F) convention.

## 2. The 1F1F Convention

### 2.1 Rules

| Rule | Description | Severity |
|---|---|---|
| F1 | One file, one function (filename = function name) | ERROR |
| F2 | One file, one type (filename = type name) | ERROR |
| F3 | One file, one method (receiver_method.go) | ERROR |
| F4 | init() allowed only alongside a func or var | ERROR |
| F5 | _test.go files may contain multiple functions | Exception |
| F6 | Semantically grouped const blocks allowed | Exception |

### 2.2 Naming Convention

```
check_nesting_depth.go    → func CheckNestingDepth(...)
server.go                 → type Server struct{...}
server_start.go           → func (s *Server) Start(...)
```

File name = concept name. No grep needed to find a function — the file system is the index.

### 2.3 Annotations

Each file carries machine-readable metadata:

```go
//ff:func feature=validate type=rule control=iteration dimension=1
//ff:what Q1: validates nesting depth against dimension-based limit
func CheckNestingDepth(gf *model.GoFile) []model.Violation {
```

Since each file is one function, file-level annotations are automatically function-level annotations.

## 3. Case Study 1: fullend (21K LOC)

### 3.1 Project Description

fullend is a Go CLI that validates consistency across 9 single sources of truth (SSOT) and generates code. Converted over 13 phases (Phase031-043).

### 3.2 Structural Metrics

| Metric | Before | After | Change |
|---|---|---|---|
| Go source files | 87 | 1,260 | +1,348% |
| Total LOC | 21,232 | 31,976 | +50.6% |
| Avg LOC/file | 244 | 25.4 | -89.6% |
| Median LOC/file | — | 21 | |
| Functions | 617 | 1,077 | +460 (extracted helpers) |
| Types | 41 | 218 | +177 (split internal types) |
| Avg func/file | 7.0 | 0.85 | -87.9% |
| Max file size | 1,113 lines | 232 lines | -79.2% |

### 3.3 LOC Increase Breakdown

The +10,744 LOC increase is entirely structural overhead:

| Source | Lines |
|---|---|
| Annotations (//ff:func + //ff:what per file) | ~2,500 |
| Package declaration + import blocks per file | ~5,000 |
| Extracted helper function signatures + imports | ~3,200 |

Zero net logic LOC increase. Every additional line is boilerplate from file separation.

### 3.4 File Size Distribution

| Range | Before | After |
|---|---|---|
| 1-25 lines | 3.4% | 64.8% |
| 26-50 lines | 14.9% | 27.2% |
| 51-100 lines | 16.1% | 7.0% |
| 101-200 lines | 34.5% | 0.8% |
| 201+ lines | 31.0% | 0.2% |

**99.0% of files are under 100 lines after conversion.**

### 3.5 SRP Violations

Before: 66/87 files (75.9%) violated F1 (multiple functions per file). Worst offender: symbol.go with 35 functions in one file.

After: 0/1,260 files violate F1. Every file contains exactly one function, one type, or one method.

### 3.6 Nesting Depth

148+ Q1 (depth) violations resolved to 0 across 13 phases. The depth-2 enforcement forced extraction of helper functions, which in turn created more single-function files.

## 4. Case Study 2: whyso (1.7K LOC)

### 4.1 Project Description

whyso is a Go CLI that extracts file change histories from Claude Code sessions and builds keyword maps via tree-sitter parsing.

### 4.2 Structural Metrics

| Metric | Before | After | Change |
|---|---|---|---|
| Go source files | 12 | 99 | +725% |
| Total LOC | 1,773 | 2,417 | +36.3% |
| Avg LOC/file | 147.8 | 24.4 | -83.5% |
| Functions | 65 | 84 | +19 helpers |
| Avg func/file | 5.4 | 1.0 | -81.5% |
| Max file size | 410 lines | 109 lines | -73.4% |

### 4.3 Conversion Cost

| Metric | Value |
|---|---|
| Conversion time | ~5 minutes (6 AI agents in parallel) |
| LOC increase | +644 lines (+36.3%) |
| Annotation overhead | 192 lines (7.9% of total) |

### 4.4 Depth Resolution

23 Q1 violations resolved via:

| Technique | Count | New files needed |
|---|---|---|
| Condition merge | 5 | 0 |
| Early continue/return | 5 | 0 |
| Helper func extraction | 13 | +13 files |

## 5. Context Contamination Analysis

### 5.1 Theoretical Model

For a project with N functions across F files, the average context contamination per read is:

```
contamination = (N/F - 1) / (N/F) = 1 - F/N
```

| Project | Before (F/N) | After (F/N) | Before contamination | After contamination |
|---|---|---|---|---|
| fullend | 87/617 = 0.14 | 1,260/1,077 = 1.17 | 86% | 0% |
| whyso | 12/65 = 0.18 | 99/84 = 1.18 | 82% | 0% |

Before conversion, reading one function brings in 5-7 irrelevant functions on average. After conversion, reading one function brings in exactly zero irrelevant functions.

### 5.2 Token Waste Estimate

Assuming average function size of 30 tokens (after conversion):

| Project | Before: tokens wasted per read | After: tokens wasted per read |
|---|---|---|
| fullend | 30 * (7.0 - 1) = 180 tokens | 0 tokens |
| whyso | 30 * (5.4 - 1) = 132 tokens | 0 tokens |

For a typical Claude Code session reading 50 files, this saves 6,600-9,000 tokens of noise — tokens that would degrade model performance per the "Lost in the Middle" findings.

## 6. The File System as Index

### 6.1 Before 1F1F

To find function `CheckNestingDepth`:
```
1. grep -r "func CheckNestingDepth" .     → find file
2. read validate.go                        → 500 lines, 20 functions
3. locate the function within the file     → scroll/search
```

### 6.2 After 1F1F

```
1. read check_nesting_depth.go            → 25 lines, 1 function
```

The file system itself becomes a searchable index. File name = function name. No grep needed for exact matches. `ls internal/validate/` shows all validation rules.

### 6.3 Git Integration

File-level git operations become function-level operations:

- `git log check_nesting_depth.go` = function change history
- `git blame check_nesting_depth.go` = function authorship
- `git diff check_nesting_depth.go` = function diff

No hunk analysis needed. The entire file diff is the function diff.

## 7. Costs and Trade-offs

### 7.1 File Count Explosion

87 -> 1,260 files (+1,348%) is the primary cost. This affects:

- **IDE navigation**: Mitigated by filename = function name convention
- **PR review**: More files changed per PR. Mitigated by smaller diffs per file.
- **Build system**: Go compiles by package, not file. No build impact.

### 7.2 LOC Overhead

+50.6% LOC increase (fullend) from structural boilerplate. This is the price of structure. Each file needs its own package declaration, imports, and annotations.

The overhead percentage decreases as function complexity increases: a 5-line function in its own file has 300% overhead (15 lines of boilerplate), but a 50-line function has 30% overhead.

### 7.3 Import Duplication

The same import may appear in 50 files within a package. This is textual duplication, not logical duplication. Go's build system deduplicates at compile time.

## 8. Comparison with Existing Practices

| Approach | File granularity | Index method | AI agent cost |
|---|---|---|---|
| Standard Go | 200-500 LOC/file | grep + read | High contamination |
| Java (1 class/file) | 100-300 LOC/file | Class name | Medium contamination |
| **1F1F** | **20-30 LOC/file** | **Filename** | **Zero contamination** |
| Microservices | N/A (service level) | API discovery | Different scope |

1F1F is the logical extreme of the "one concept per file" principle that Java partially implements with its one-class-per-file rule.

## 9. Limitations

- Two Go projects from one developer. Cross-language and cross-team evaluation needed.
- No direct measurement of AI agent task completion accuracy (before/after). We measure structural proxy metrics.
- The "Lost in the Middle" effect is model-dependent. Future models may be less sensitive to context noise.
- File count explosion may have practical limits at very large scales (100K+ files).

## 10. Conclusion

File granularity is the primary lever for AI code agent context quality. The 1F1F convention eliminates context contamination — the irrelevant code that enters an agent's context window when reading a file to find one function.

Across two projects totaling 23,005 LOC:

1. **Context contamination: 82-86% -> 0%** — every read returns exactly what's needed
2. **Average file size: 147-244 lines -> 24-25 lines** — 84-90% reduction
3. **SRP violations: 76-100% -> 0%** — structural quality guaranteed
4. **Depth violations: 23-148+ -> 0** — complexity constrained

The cost — 36-51% LOC increase and 725-1,348% file count increase — is structural overhead that buys zero-contamination reads. In AI-assisted development where reads vastly outnumber writes, this trade-off favors structure.

---

## Citation

```
@misc{one-file-one-func2026,
  title={One File, One Function: How File Granularity Affects AI Code Agent Performance},
  author={Park, Jun-Woo},
  year={2026},
  url={https://github.com/park-jun-woo/filefunc}
}
```
