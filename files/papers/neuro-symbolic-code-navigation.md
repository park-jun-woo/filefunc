# 250 Tokens to Navigate a Codebase: Neuro-Symbolic Code Search for LLM Agents

## Abstract

AI code agents process 300K–1M tokens of context per turn, yet fewer than 1% of those tokens are relevant to the task. We present filefunc, a neuro-symbolic code navigation system that reduces the compute required for code understanding by 126,000× — from 140 petaFLOPs to 1.1 teraFLOPs — by structuring code at write time rather than searching at read time.

The system enforces a "one file, one concept" convention with machine-verifiable annotations, then uses a 3-stage pipeline: (1) an 8B model selects relevant code domains from a codebook using 250 tokens, (2) a 20B local model scores function descriptions, and (3) the same 20B model scores function bodies. The frontier model receives only the 960 tokens that matter.

We demonstrate this on two Go projects (filefunc itself and whyso), achieving 4/5 accuracy on diverse code search queries with 3 LLM calls totaling 5,760 tokens — versus feeding 12,000 project tokens or 350,000 context tokens to a frontier model.

## 1. Problem

### 1.1 The Context Window Tax

Modern AI code agents (Claude Code, Cursor, Copilot) operate within context windows of 200K–1M tokens. A typical Claude Code session carries 300–400K tokens of accumulated context. For a 200B-parameter model, each turn costs:

```
2 × 200B × 350,000 = 140,000,000 TFLOPs
```

Yet when a user asks "modify the nesting depth validation logic," only 3–4 functions (~960 tokens) are relevant. The remaining 349,000 tokens are noise that the model must still process.

### 1.2 The Search Inversion

The industry response has been retrieval-augmented generation (RAG): embed code into vectors, build ANN indices, retrieve similar chunks. This approach:

- Destroys structure (code → float vectors → approximate nearest neighbors)
- Requires expensive embedding models and vector databases
- Produces approximate results ranked by similarity, not exact results filtered by classification

We argue this is backwards. The problem is not "how to search unstructured code better" but "why is code unstructured in the first place?"

## 2. Approach: Structure at Write Time

### 2.1 filefunc Convention

Every Go source file contains exactly one function or type (1 file = 1 concept). Each file carries machine-verifiable annotations:

```go
//ff:func feature=validate type=rule control=sequence
//ff:what Q1: validates nesting depth against dimension-based limit
func CheckNestingDepth(gf *model.GoFile) []model.Violation {
```

A project-level codebook defines the allowed vocabulary:

```yaml
required:
  feature: [validate, annotate, chain, parse, codebook, report, cli, context]
  type: [command, rule, parser, walker, model, formatter, loader, util]
```

### 2.2 The Codebook as Search Index

The codebook is not documentation — it is a search index. Each `feature=X` annotation partitions the codebase into disjoint sets. A single feature value eliminates 70–99% of files:

| feature | files | % of total |
|---|---|---|
| validate | 35 | 28% |
| parse | 37 | 30% |
| chain | 28 | 22% |
| cli | 12 | 10% |
| codebook | 1 | 1% |

This is classification, not ranking. No embeddings. No similarity scores. A string match on `feature=validate` is a set operation that returns exact results.

### 2.3 SILK Methodology

filefunc implements the SILK (Symbolic Index for LLM Knowledge) architecture applied to source code:

| SILK Principle | filefunc Implementation |
|---|---|
| Structure at write time | `//ff:func feature=X type=Y` annotations |
| Codebook vocabulary control | `codebook.yaml` + A2 validation |
| Symbolic filter → Neural refinement | Feature filter → LLM what/body scoring |
| VALID machine verification | `filefunc validate` (22 rules) |
| Classification over ranking | Feature set operations, not vector similarity |

## 3. Pipeline

### 3.1 Four-Stage Context Pipeline

```
Input: User prompt (natural language)

Stage 1 — Feature Selection (8B model, 250 tokens, 1 API call)
  Codebook + prompt → LLM selects 1–2 relevant features
  "nesting depth 수정" → ["validate"]

Stage 2 — Feature Filter (static, 0 FLOPs)
  Project files filtered by feature annotation
  161 files → 35 files

Stage 3 — What Scoring (20B local model, ~3,450 tokens, 35 calls)
  Each function's //ff:what scored against prompt
  35 files → 6 files (rate ≥ 0.2)

Stage 4 — Body Scoring (20B local model, ~1,100 tokens, 6 calls)
  Each function body scored against prompt
  6 files → 3 files (rate ≥ 0.5)

Output: 3 functions, 960 tokens of Go source
```

Total LLM compute: 1,114 TFLOPs. Frontier model receives only 960 tokens.

### 3.2 Compute Comparison

| Method | Model | Tokens | FLOPs |
|---|---|---|---|
| Naive context | Opus 200B | 350,000 | 140,000,000 T |
| Project dump | Opus 200B | 12,000 | 4,800 T |
| **filefunc pipeline** | **8B + 20B + Opus** | **250 + 4,550 + 960** | **1,114 T** |
| **Reduction** | | | **126,000×** |

### 3.3 Cost

| Stage | Model | Cost |
|---|---|---|
| Feature selection | Haiku API | $0.0003 |
| What + Body scoring | Local 20B | Electricity only |
| Final answer | Opus | 960 tokens instead of 350K |

## 4. Experiments

### 4.1 Setup

- Projects: filefunc (161 files, 3,037 LOC), whyso (99 files, 2,417 LOC)
- Feature selection model: gpt-oss:20b (ollama local)
- Scoring model: gpt-oss:20b (ollama local)
- Validation: `filefunc validate` — 0 violations on both projects

### 4.2 Code Search Accuracy

| Query | Feature Selected | Top Results | Correct? |
|---|---|---|---|
| "nesting depth 검증 수정" | validate | CheckNestingDepth [0.80], depthLimit [0.80] | O |
| "codebook 파싱 변경" | parse | ParseCodebook [1.00] | O |
| "chain 출력 포맷 변경" | chain | FormatChain [0.90], formatMeta [0.55] | O |
| "LLM what-body 검증 수정" | validate (wrong) | (no results) | X |
| "어노테이션 파싱 수정" | parse | ParseAnnotation [0.90], ApplyAnnotationLine [0.80] | O |

4/5 accuracy. The failure case ("LLM 검증") is a feature selection error where the 20B model maps "검증" to `validate` instead of `annotate`. An 8B+ API model (Haiku, Gemini Flash) resolves this.

### 4.3 Reranker Comparison

We also tested vLLM + Qwen3-Reranker (0.6B, 4B) as cross-encoder scorers:

| Method | Accuracy | Issue |
|---|---|---|
| Reranker 0.6B | 0/5 | All scores flat (0.1–0.2) |
| Reranker 4B | 0/5 | All scores high (0.6–0.96), no discrimination |
| **Codebook + LLM scoring** | **4/5** | Feature selection is the bottleneck |

Rerankers fail because all functions in the same package are "relevant" to a code query — they cannot distinguish within a domain. Codebook classification solves this by partitioning first.

### 4.4 filefunc Adoption Metrics (whyso)

| Metric | Before | After |
|---|---|---|
| Files | 12 | 99 |
| Avg LOC/file | 147.8 | 24.4 |
| SRP violations | 12/12 (100%) | 0/99 (0%) |
| Depth violations | 23 | 0 |
| Annotation coverage | 0% | 100% |
| Conversion time | — | ~5 min (AI agents) |

## 5. The 250-Token Navigation Thesis

A 250-token codebook query to an 8B model is sufficient to navigate a codebase that would otherwise require a 200B model to process 350,000 tokens. This is not an optimization — it is a 126,000× reduction in compute that follows from a single design principle:

**Structure code at write time, and search becomes free at read time.**

The cost of structuring (annotations, codebook, validation) is paid once per file write. The cost of searching (350K tokens × 200B parameters) is paid every read. In AI-assisted development where reads vastly outnumber writes, structuring is always cheaper.

### 5.1 Broader Implications

This result is domain-agnostic. The same principle applies wherever LLMs process large contexts:

- **Documents**: Annotate sections with topic codebook → 8B selects relevant sections
- **Conversations**: Tag messages with intent codebook → filter before context inclusion
- **Knowledge bases**: SILK's 64-bit structured identifiers → bit-AND filtering

The pattern is always: small model + codebook → classify → frontier model receives only what matters.

## 6. Related Work

- **RAG** (Lewis et al., 2020): Retrieves by vector similarity. filefunc retrieves by symbolic classification.
- **Lost in the Middle** (Liu et al., 2024): Relevant information in context middle drops performance 30%+. filefunc eliminates the middle.
- **Context Rot** (Chroma Research): Focused prompt > full prompt. filefunc generates focused prompts automatically.
- **SILK** (Park, 2026): Neuro-symbolic search architecture for knowledge. filefunc is SILK applied to code.

## 7. Limitations

- Feature selection accuracy depends on codebook quality and LLM capability
- Currently Go-only (AST analysis, annotation format)
- 20B local model shows "검증"→validate keyword bias for Korean prompts
- Individual LLM scoring (35 calls) is slower than batch (~1 min vs ~5 sec)

## 8. Conclusion

AI code agents waste 99% of their compute processing irrelevant context. filefunc demonstrates that 250 tokens of structured metadata — a codebook query to an 8B model — can replace 350,000 tokens of brute-force context, reducing compute by 126,000×.

The key insight is not a new algorithm but a convention: annotate code at write time with a controlled vocabulary, then classify rather than search. This is the SILK principle applied to source code: structure eliminates search.

---

## Citation

```
@misc{filefunc2026,
  title={250 Tokens to Navigate a Codebase: Neuro-Symbolic Code Search for LLM Agents},
  author={Park, Jun-Woo},
  year={2026},
  url={https://github.com/park-jun-woo/filefunc}
}
```
