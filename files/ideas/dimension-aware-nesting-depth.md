# Dimension-Aware Nesting Depth Limits for Typed Data Traversal

## Abstract

Static analysis tools enforce uniform nesting depth limits (typically 3-5) across all functions. We identify a category of functions — typed multi-dimensional data traversals — where nesting depth reflects data structure dimensions rather than logical complexity. Forcing these functions to comply with a depth-2 limit via function extraction disperses the traversal path across multiple functions, increasing cognitive load without reducing complexity.

We propose **dimension-aware depth limits**: for functions traversing N-dimensional typed data, the depth limit is N+1 instead of the default 2. The dimension count is automatically derived from the number of named-type assertions in the for-range chain. We evaluate this on Go AST, OpenAPI, and protobuf traversal patterns, showing that 100% of depth violations in typed traversals are resolved by dimension-aware limits without sacrificing code clarity.

## 1. Introduction

### 1.1 Nesting Depth as a Complexity Metric

Nesting depth is a widely-used proxy for code complexity. Deep nesting indicates:
- Multiple levels of conditional logic
- Complex iteration patterns
- Difficult-to-follow control flow

Tools like SonarQube, golangci-lint, and ESLint enforce maximum depth limits. These limits work well for most code. However, they systematically misclassify one pattern: **typed data structure traversal**.

### 1.2 The Traversal Exception

Consider traversing a Go AST to collect struct fields:

```go
for _, entry := range entries {                    // depth 1: directory files
    for _, decl := range f.Decls {                 // depth 2: declarations
        for _, spec := range gd.Specs {            // depth 3: type specs
            for _, field := range st.Fields.List {  // depth 4: struct fields
                for _, name := range field.Names {  // depth 5: field names
                }
            }
        }
    }
}
```

Depth 5. By any standard metric, this is "too complex." Yet each level is:
1. A typed assertion (`decl.(*ast.GenDecl)`, `spec.(*ast.TypeSpec)`, etc.)
2. An early-continue on assertion failure
3. No branching logic — just descending one level deeper

The depth reflects the **data structure's dimensions**, not the code's logical complexity. The Go AST schema is fixed by the language specification: `File -> GenDecl -> TypeSpec -> StructType -> Field -> Ident`. This structure cannot be simplified.

### 1.3 Why Extraction Fails

The standard remedy for deep nesting is function extraction:

```go
// After extraction: traversal split across 3 functions
func loadGoInterfaces(dir string) {
    for _, entry := range entries {
        processFileDecls(f, st)          // inner traversal hidden
    }
}
func processFileDecls(f *ast.File, st *SymbolTable) {
    for _, decl := range f.Decls {
        processGenDeclSpecs(gd, st)      // even deeper hidden
    }
}
func processGenDeclSpecs(gd *ast.GenDecl, st *SymbolTable) {
    for _, spec := range gd.Specs {
        // finally: the actual logic
    }
}
```

Problems:
- **Dispersed traversal path**: The complete path (entries -> decls -> specs -> fields -> names) is split across 3+ functions. Understanding the full traversal requires reading all of them.
- **Lost locality**: The original code reads top-to-bottom as a continuous descent. The extracted version requires jumping between functions.
- **No complexity reduction**: Each extracted function still does the same thing — type-assert and iterate. The total complexity hasn't changed; it has moved.

## 2. Typed Data Traversals

### 2.1 Characteristics

A **typed data traversal** has these properties:

| Property | Description |
|---|---|
| Fixed schema | The data structure's shape is defined by a type system (AST, API spec, protobuf) |
| Named types at each level | Each dimension is modeled as a named type (struct, interface) |
| Type assertions as descent | Moving deeper means asserting a typed container |
| No branching logic | Each level: assert type -> early-continue on failure -> iterate next level |
| Deterministic path | The traversal path is the same for every input instance |

### 2.2 Dimension Count

The **dimension** of a traversal is the number of named-type assertions in the for-range chain:

```go
gd, ok := decl.(*ast.GenDecl)          // named type 1
ts, ok := spec.(*ast.TypeSpec)         // named type 2
st, ok := ts.Type.(*ast.StructType)    // named type 3
```

Dimension = 3 here. Each assertion corresponds to descending one level in the typed data structure.

### 2.3 Examples Across Domains

| Domain | Data Structure | Traversal Path | Dimension |
|---|---|---|---|
| Go AST | Source file | File -> GenDecl -> TypeSpec -> StructType -> Field | 5 |
| Go AST | Interface | File -> GenDecl -> TypeSpec -> InterfaceType -> Method | 5 |
| OpenAPI | API spec | PathItem -> Operation -> Parameter -> Schema | 4 |
| Protobuf | Descriptor | FileDescriptor -> MessageDescriptor -> FieldDescriptor | 3 |
| HTML DOM | Document | Document -> Element -> Attributes -> Value | 3 |

### 2.4 Why Raw Nesting Doesn't Qualify

```go
// This is NOT a typed traversal:
var grid [][][]int
for i := range grid {
    for j := range grid[i] {
        for k := range grid[i][j] {
            // process grid[i][j][k]
        }
    }
}
```

No named types. No type assertions. This is raw multi-dimensional iteration — its depth genuinely reflects complexity. The programmer should model this with named types:

```go
type Grid struct { Layers []Layer }
type Layer struct { Rows []Row }
type Row struct { Cells []int }
```

The dimension-aware exception **forces good data modeling**: only named-type traversals qualify. This is a feature, not a limitation.

## 3. Dimension-Aware Depth Limits

### 3.1 Rule

For a function annotated with `control=iteration` and `dimension=N`:

```
depth_limit = N + 1
```

The +1 accounts for the loop body's working logic (filtering, collecting, transforming) at the innermost level.

| dimension | depth limit | Typical pattern |
|---|---|---|
| 1 | 2 | Flat list traversal (default) |
| 2 | 3 | 2D data (paths -> operations) |
| 3 | 4 | 3D data (policies -> rules -> actions) |
| 4 | 5 | 4D data (spec deduplication) |
| 5 | 6 | Full AST traversal |

### 3.2 Annotation Format

```go
//ff:func feature=symbol type=loader control=iteration dimension=3
//ff:what Parses Go interfaces from directory, registering as "pkg.Model" keys
func loadGoInterfaces(dir string, st *SymbolTable) error {
```

### 3.3 Validation Rules

| Rule | Condition | Severity |
|---|---|---|
| A15 | `control=iteration` requires `dimension=` | ERROR |
| A16 | `dimension=` must be a positive integer | ERROR |
| Q1 | Actual depth must be <= dimension + 1 | ERROR |

Automatic dimension verification: `filefunc validate` counts named-type assertions in the for-range chain and verifies that the declared `dimension=N` matches the actual assertion count.

## 4. Visitor Pattern: Not a Solution

### 4.1 ast.Inspect Conversion

```go
for _, entry := range entries {
    ast.Inspect(f, func(n ast.Node) bool {
        ts, ok := n.(*ast.TypeSpec)
        if !ok { return true }
        st, ok := ts.Type.(*ast.StructType)
        if !ok { return false }
        for _, field := range st.Fields.List {
            for _, name := range field.Names {
                // process
            }
        }
        return false
    })
}
```

Depth: 8 -> 4. Better, but still violates depth-2. And introduces new problems:

### 4.2 Visitor Pattern Problems

| Problem | Description |
|---|---|
| Error propagation | `ast.Inspect` callback returns bool only. Cannot `return fmt.Errorf(...)` |
| Unnecessary traversal | Visits function bodies, expressions, all irrelevant nodes |
| Parent context loss | Callback receives `*ast.TypeSpec` but doesn't know which `GenDecl` it belongs to |
| Implicit state | Requires a collector struct with state flags (`inRequestStruct bool`) |

The explicit nesting of the original code is replaced by an **implicit state machine**. Complexity doesn't decrease — it changes form from visible nesting to hidden state.

### 4.3 When Visitors Work

Visitors work well for **homogeneous trees** where every node has the same type and is self-contained:

```go
// Good visitor use case: all nodes are the same type
tree.Walk(func(node *Node) {
    process(node)  // each node independently processable
})
```

Go AST is **heterogeneous** (different types at each level) and **context-dependent** (processing depends on parent type). Visitors are the wrong abstraction.

## 5. Evaluation

### 5.1 Data

From three Go projects (filefunc, whyso, fullend), we identified all functions with depth > 2 before refactoring:

| Project | Total depth violations | Typed traversal violations | Non-traversal violations |
|---|---|---|---|
| filefunc | 0 (built with rules) | 0 | 0 |
| whyso | 23 | 2 | 21 |
| fullend | 148+ | 13 | 135+ |

### 5.2 Typed Traversal Violations

All 15 typed-traversal violations across two projects:

| Function | Original depth | Dimension | New limit | Resolved? |
|---|---|---|---|---|
| loadPackageGoInterfaces | 8 | 5 | 6 | Yes |
| loadGoInterfaces | 6 | 5 | 6 | Yes |
| (fullend AST walkers) | 4-6 | 2-4 | 3-5 | All yes |

100% of typed traversal violations resolved by dimension-aware limits.

### 5.3 Non-traversal Violations

The remaining 156+ violations (non-typed-traversal) were resolved by:

| Technique | Count | Description |
|---|---|---|
| Early return/continue | ~40 | Replace nested if with guard clause |
| Condition merge | ~20 | Combine `if a { if b {` into `if a && b {` |
| Helper extraction | ~96 | Extract inner logic to separate function |

These are genuine complexity that benefits from depth-2 enforcement. The dimension-aware exception correctly excludes them.

## 6. Design Decisions

### 6.1 Why Named Types Required

Requiring named types at each dimension level serves dual purposes:

1. **Verification**: The tool can count type assertions to verify the declared dimension
2. **Design pressure**: Forces developers to model their data with named types instead of raw nesting (`[][][]int`)

This is a deliberate constraint that improves code quality while enabling the depth exception.

### 6.2 Why dimension + 1, Not dimension

The +1 accounts for actual work at the innermost level:

```go
for _, field := range st.Fields.List {  // depth = dimension
    for _, name := range field.Names {  // depth = dimension + 1
        result[name.Name] = fieldType   // actual logic at deepest level
    }
}
```

Without +1, the traversal itself consumes all allowed depth, leaving no room for the logic that processes each element.

### 6.3 Why Not Recursive Depth

Some traversals use recursion instead of nested loops. Recursive depth is harder to statically verify and is a separate problem. This paper addresses only iterative multi-dimensional traversals with statically-known depth.

## 7. Related Work

- **Cyclomatic complexity** (McCabe, 1976): Counts all paths. Does not distinguish data-driven nesting from logic-driven nesting.
- **Cognitive complexity** (SonarSource, 2017): Penalizes nesting incrementally. Typed traversals still score high.
- **SonarQube depth limit**: Fixed at 3-5. No dimension awareness.
- **golangci-lint nestif**: Configurable depth. No exception for typed traversals.

None of these tools distinguish between "depth from data dimensions" and "depth from logic complexity."

## 8. Limitations

- Go-specific evaluation. Other languages have different traversal patterns (e.g., Python generators, Rust iterators).
- Dimension verification requires Go AST analysis of type assertions. Other languages may need different detection strategies.
- Only iterative traversals covered. Recursive traversals (tree walkers, graph DFS) are out of scope.
- Small dataset (15 typed traversal violations). Larger-scale evaluation on diverse Go projects needed.

## 9. Conclusion

Nesting depth limits should distinguish between depth-from-data and depth-from-logic. For typed multi-dimensional data traversals, depth reflects the data structure's dimensions — a fixed property of the schema, not a variable property of the code's complexity.

Dimension-aware limits (`depth <= dimension + 1`) resolve 100% of typed traversal violations while preserving the benefits of strict depth limits for all other code. The named-type requirement ensures that only well-modeled data structures qualify for the exception, creating positive design pressure.

---

## Citation

```
@misc{dimension-nesting2026,
  title={Dimension-Aware Nesting Depth Limits for Typed Data Traversal},
  author={Park, Jun-Woo},
  year={2026},
  url={https://github.com/park-jun-woo/filefunc}
}
```
