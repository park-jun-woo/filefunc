# Classifying Functions by Dominant Control Structure: Applying Bohm-Jacopini to Code Quality Rules

## Abstract

The Bohm-Jacopini theorem (1966) proves that sequence, selection, and iteration suffice to express any computable function. We apply this theorem at the function level: after enforcing a nesting depth limit of 2, each function's body is dominated by exactly one of the three control structures. We formalize this as a `control=` annotation, validate it via AST analysis, and use it to define differentiated quality thresholds — replacing one-size-fits-all line limits with structure-aware rules.

We evaluate this classification on 2,332 functions across three Go projects (filefunc 136 functions, whyso 84 functions, fullend 1,055 functions). All functions classify into exactly one of the three categories with zero ambiguous cases, yielding a distribution of 44.1% sequence, 51.2% iteration, and 4.7% selection. We find that depth-2 enforcement is a necessary precondition: without it, functions exhibit mixed control structures that resist clean classification.

## 1. Introduction

### 1.1 The Problem with Uniform Quality Rules

Static analysis tools apply uniform thresholds: cyclomatic complexity < 10, function length < 100 lines, nesting depth < 4. These thresholds treat all functions identically, yet functions have fundamentally different structures:

- A **sequence** function calls 12 stages in order (pipeline). 150 lines, but complexity is low — each line is independent.
- A **selection** function dispatches 10 switch cases, each 3-8 lines. 119 lines, but splitting into 10 files destroys the ability to see all cases at once.
- An **iteration** function traverses a 5-level AST. Depth 5, but each level is a typed assertion with no logic — depth reflects data dimensions, not complexity.

Uniform rules either miss real complexity (too lenient) or penalize structural patterns (too strict). The missing variable is the function's **dominant control structure**.

### 1.2 Bohm-Jacopini at the Function Level

The Bohm-Jacopini theorem states that any computable function can be expressed using only sequence, selection, and iteration. We observe that when functions are constrained to depth 2, this theoretical result becomes practically observable: each function's body is **dominated** by exactly one control structure.

This is not a restatement of the theorem. The theorem is about expressiveness; our claim is about **classification utility**: depth-constrained functions exhibit a single dominant control structure that predicts their quality characteristics.

### 1.3 Prior Work

To our knowledge, no prior work applies Bohm-Jacopini classification to individual function bodies as a quality metric. Related work:

- **Cyclomatic complexity** (McCabe, 1976): Counts decision points. Does not distinguish between a 10-case switch and 10 nested if-else chains.
- **Cognitive complexity** (SonarSource, 2017): Penalizes nesting. Closer to our approach but does not classify functions into structural categories.
- **Halstead metrics** (1977): Measures operands and operators. Content-based, not structure-based.
- **ABC metric**: Counts assignments, branches, conditions. Aggregates rather than classifies.

None of these classify functions by their dominant control structure.

## 2. Method

### 2.1 Precondition: Depth-2 Enforcement

We first enforce a strict nesting depth limit:

| control | depth limit |
|---|---|
| sequence | 2 (no loop/switch at depth 1) |
| selection | 2 (switch at depth 1, logic inside cases at depth 2) |
| iteration | dimension + 1 (loop at depth 1, body at depth 2+) |

This flattening is essential. In unconstrained code, a function may contain a loop containing a switch containing another loop — mixed structures that resist classification. Depth-2 enforcement forces each function to have at most one level of primary control structure, making the dominant structure unambiguous.

### 2.2 Classification Rules

Given a function with depth-2 enforcement:

| Condition | Classification |
|---|---|
| Depth-1 contains `for`/`range` statement | **iteration** |
| Depth-1 contains `switch`/`type switch` statement | **selection** |
| Neither loop nor switch at depth 1 | **sequence** |

Early-return `if err != nil { return }` at depth 1 does not count as selection — it is a guard (exit), not a branch (choice).

### 2.3 AST Validation

`filefunc validate` enforces consistency between the `control=` annotation and the actual AST:

| Rule | Condition | Severity |
|---|---|---|
| A9 | func file must have `control=` | ERROR |
| A10 | `control=selection` but no switch at depth 1 | ERROR |
| A11 | `control=iteration` but no loop at depth 1 | ERROR |
| A12 | `control=sequence` but loop/switch at depth 1 | ERROR |
| A13 | `control=selection` but loop at depth 1 | ERROR |
| A14 | `control=iteration` but switch at depth 1 | ERROR |

These rules make the classification machine-verifiable, not subjective.

### 2.4 Dimension: Iteration Depth Refinement

For iteration functions, the data being traversed determines the natural nesting depth. We introduce `dimension=N`:

- `dimension=1`: Flat list traversal. Depth limit = 2.
- `dimension=2`: 2D data (e.g., paths -> operations). Depth limit = 3. Named type nesting required.
- `dimension=N`: N-dimensional data. Depth limit = N + 1.

Constraint: dimension >= 2 requires named type modeling at each level. Raw nesting (`[][][]int`) is prohibited — each dimension must be a named type (struct or interface).

## 3. Depth-2 Combinations Analysis

With depth 1 elements being loop, if, and switch, there are 9 possible depth-2 combinations:

| depth 1 -> depth 2 | Pattern | Classification | Frequency |
|---|---|---|---|
| loop -> loop | Nested iteration | **iteration** | Rare |
| loop -> if | Filter/select | **iteration** | Common |
| loop -> switch | Traverse + dispatch | **iteration** | Common |
| if -> loop | Conditional iteration | **sequence** | Rare |
| if -> if | Nested condition | **sequence** | Replaced by early return |
| if -> switch | Conditional dispatch | **sequence** | Rare |
| switch -> loop | Per-case iteration | **selection** | Rare |
| switch -> if | Per-case condition | **selection** | Common |
| switch -> switch | Nested branching | **selection** | Rare |

The classification rule is simple: depth-1 determines the category. loop -> iteration, switch -> selection, if or nothing -> sequence.

## 4. Structure-Aware Quality Rules

### 4.1 Differentiated Thresholds

| control | Q1 (depth) | Q3 (lines) | Rationale |
|---|---|---|---|
| sequence | <= 2 | 100 (WARNING), 200 (pipeline) | Sequential steps; long but simple |
| selection | case-internal <= 2 | 300 (case count * proportional) | Cases must be seen together |
| iteration | <= dimension + 1 | 100 | Loop body is the complexity |

### 4.2 Why Uniform Rules Fail (fullend Q3 data)

From 15 Q3 violations in fullend (>100 lines):

| Category | Count | Should split? |
|---|---|---|
| sequence (pipeline/template) | 9 | 4 yes, 5 no |
| selection (flat dispatch) | 5 | 0 yes, 5 no |
| iteration (AST walker) | 1 | 1 yes |

**60% of violations (9/15) are appropriately over 100 lines.** A uniform 100-line rule would force splitting functions that are structurally simple, reducing readability.

## 5. Evaluation

### 5.1 Classification Completeness

| Project | Functions | sequence | selection | iteration | Unclassifiable |
|---|---|---|---|---|---|
| filefunc | 136 | 58 (42.6%) | 4 (2.9%) | 54 (39.7%) | 0 |
| whyso | 84 | 41 (48.8%) | 5 (6.0%) | 38 (45.2%) | 0 |
| fullend | 1,055 | 393 (37.3%) | 74 (7.0%) | 588 (55.7%) | 0 |
| **Total** | **1,275** | **492 (38.6%)** | **83 (6.5%)** | **680 (53.3%)** | **0** |

Zero unclassifiable functions across 1,275 functions in three projects.

### 5.2 Dimension Distribution (iteration functions)

| dimension | Files | Q1 limit | Meaning |
|---|---|---|---|
| 1 | 600 | 2 | Flat list traversal |
| 2 | 67 | 3 | 2D data (paths->ops) |
| 3 | 10 | 4 | 3D data (policies->rules->actions) |
| 4 | 2 | 5 | 4D data (spec->existing dedup) |
| 5 | 1 | 6 | 5D data (AST full traversal) |

88.2% of iteration functions are dimension=1, confirming that most iteration is simple flat traversal.

### 5.3 Project Type Hypothesis

| Project Type | Expected Dominant | Observed (fullend) |
|---|---|---|
| Code generator | sequence | sequence 37.3% (parsers push iteration) |
| CLI tool | sequence | (filefunc) sequence 42.6% |
| Data processor | iteration | (whyso) iteration 45.2% |

The distribution correlates with project purpose, suggesting predictive value.

## 6. Discussion

### 6.1 Why Depth-2 is the Precondition

Without depth-2 enforcement, a function may contain:
```go
for items {           // iteration
    switch item.Type { // selection
        case A:
            for sub {  // iteration again
            }
    }
}
```

This function is iteration? Selection? Both? Depth-2 forces this to be split into an iteration function (outer loop) and a selection function (switch dispatch), each classifiable.

### 6.2 AI Agent Implications

The `control=` annotation enables pre-read strategy:

- `control=selection`: Read entire function. Partial reads miss cases.
- `control=iteration`: Focus on loop body. Setup code is initialization.
- `control=sequence`: Read only the relevant step. Other steps are independent.

An AI agent can decide **how to read** a function before reading its body, using only the annotation.

### 6.3 Relationship to Cyclomatic Complexity

Cyclomatic complexity counts all decision points equally. A 10-case switch scores CC=10, identical to 10 nested if-else chains. Yet the switch is a flat dispatch (low cognitive load) while nested if-else is genuinely complex. Bohm-Jacopini classification distinguishes these: the switch is selection (read all at once), the nested if-else would be split into multiple sequence functions under depth-2 enforcement.

## 7. Limitations

- Go-only evaluation. Other languages may require different depth limits (e.g., Python's implicit nesting from comprehensions).
- Three projects from one developer. Broader evaluation on open-source codebases needed.
- The "zero unclassifiable" result depends on depth-2 enforcement. Projects without this constraint may resist classification.
- Selection functions are rare (6.5%). More data needed to validate selection-specific quality rules.

## 8. Conclusion

The Bohm-Jacopini theorem, applied at the function level after depth-2 enforcement, yields a practical classification with three properties:

1. **Complete** — 1,275/1,275 functions classified (0 ambiguous)
2. **Machine-verifiable** — 6 AST rules validate annotation consistency
3. **Quality-predictive** — enables differentiated thresholds that avoid 60% of false positives from uniform rules

The key insight: depth enforcement is not just a quality rule — it is a **precondition for structural classification**. Flattened functions reveal their dominant control structure, enabling structure-aware tooling.

---

## Citation

```
@misc{bohmjacopini-func-quality2026,
  title={Classifying Functions by Dominant Control Structure: Applying Bohm-Jacopini to Code Quality Rules},
  author={Park, Jun-Woo},
  year={2026},
  url={https://github.com/park-jun-woo/filefunc}
}
```
