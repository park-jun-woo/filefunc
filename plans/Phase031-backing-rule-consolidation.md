# Phase 031: backing 기반 룰 통합 — 20개 룰을 9개 함수로 수렴

## 목표

toulmin backing을 활용하여 동일 개념의 룰 함수들을 통합한다. (1) 4개 범용 룰 함수 + backing 타입 신규 생성 (2) 5개 고유 룰 함수 이름 개선 (3) 16개 개별 룰 파일 삭제 (4) validate_graph.go를 선언적 룰북으로 재작성 (5) 테스트 갱신.

## 배경

### 현재 상태

- 20개 warrant 룰 + 3개 defeater = 23개 함수, 각각 별도 파일
- 모든 룰이 `backing=nil` — backing 미활용
- 동일 패턴의 함수가 복수 파일에 중복 존재

### 관찰된 개념 중복

| 개념 | 해당 룰 | 패턴 |
|---|---|---|
| 조건부 존재 검사 | A1a, A1b, A3, A9, A15, F4 | "조건 충족 시 대상이 존재해야 함" |
| 수량 상한 | F1, F2, F3 | "항목 수가 상한 초과하면 안 됨" |
| 제어 구조 일치 | A10, A11, A12, A13, A14 | "선언과 AST 실체가 일치해야 함" |
| 코드북 적합 | A2, A8 | "코드북 ↔ 어노테이션 교차 검증" |

### 설계 원칙

- **backing이 판정 기준을 담는다** — 함수는 개념, backing은 개별 룰의 구체적 기준
- **그래프 선언이 곧 룰북** — 위에서 아래로 읽으면 명세서
- **ruleID 안정성** — backing은 `fmt.Sprint` 가능한 값 타입만 사용 (함수 포인터 금지)

## 설계

### 범용 룰 함수 4개

#### 1. ExistsWhen — "조건 충족 시 대상이 존재해야 함"

```go
// ExistsWhenBacking은 조건부 존재 검사의 판정 기준이다.
type ExistsWhenBacking struct {
    When    string // 전제 조건: "HasFuncs", "HasTypes", "HasFuncOrType", "HasInit", "ControlIteration"
    Need    string // 존재해야 하는 것: "ff:func", "ff:type", "ff:what", "control", "dimension", "companion"
    Rule    string
    Level   string
    Message string
}
```

When/Need를 문자열로 표현하고 함수 내부에서 switch로 분기. 6개 backing:

| backing | When | Need | 원래 룰 |
|---|---|---|---|
| A1-func | `HasFuncs` | `ff:func` | A1 (func) |
| A1-type | `HasTypes` | `ff:type` | A1 (type) |
| A3 | `HasFuncOrType` | `ff:what` | A3 |
| A9 | `HasFuncs` | `control` | A9 |
| A15 | `ControlIteration` | `dimension` | A15 |
| F4 | `HasInit` | `companion` | F4 |

#### 2. CountMax — "항목 수가 상한 초과하면 안 됨"

```go
type CountMaxBacking struct {
    Field   string // "Funcs", "Types", "Methods"
    Max     int
    Rule    string
    Message string
}
```

3개 backing:

| backing | Field | Max | 원래 룰 |
|---|---|---|---|
| F1 | `Funcs` | 1 | F1 |
| F2 | `Types` | 1 | F2 |
| F3 | `Methods` | 1 | F3 |

#### 3. ControlMatch — "선언된 제어 구조와 AST 실체가 일치해야 함"

```go
type ControlMatchBacking struct {
    Control     string // 선언 control: "selection", "iteration", "sequence"
    MustHave    string // 있어야 하는 것: "switch", "loop", ""
    MustNotHave string // 없어야 하는 것: "loop", "switch", "switch|loop", ""
    Rule        string
    Message     string
}
```

5개 backing:

| backing | Control | MustHave | MustNotHave | 원래 룰 |
|---|---|---|---|---|
| A10 | selection | switch | | A10 |
| A11 | iteration | loop | | A11 |
| A12 | sequence | | switch\|loop | A12 |
| A13 | selection | | loop | A13 |
| A14 | iteration | | switch | A14 |

#### 4. InCodebook — "코드북 ↔ 어노테이션 교차 검증"

```go
type InCodebookBacking struct {
    Direction string // "value→codebook" (A2), "codebook→annotation" (A8)
    Rule      string
    Message   string
}
```

2개 backing:

| backing | Direction | 원래 룰 |
|---|---|---|
| A2 | `value→codebook` | A2 |
| A8 | `codebook→annotation` | A8 |

### 고유 룰 함수 5개 (이름 개선)

| 현재 파일 | 새 파일 | 새 함수명 | 룰 |
|---|---|---|---|
| `rule_q1.go` | `check_depth_limit.go` | `CheckDepthLimit` | Q1 |
| `rule_q2q3.go` | `check_func_lines.go` | `CheckFuncLines` | Q2, Q3 |
| `rule_a6.go` | `annotation_at_top.go` | `AnnotationAtTop` | A6 |
| `rule_a7.go` | `checked_hash_match.go` | `CheckedHashMatch` | A7 |
| `rule_a16.go` | `valid_dimension.go` | `ValidDimension` | A16 |

### Defeater 함수 3개 (이름 개선)

| 현재 파일 | 새 파일 | 새 함수명 |
|---|---|---|
| `defeater_test_file.go` | `is_test_file.go` | `IsTestFile` |
| `defeater_const_only.go` | `is_const_only_defeater.go` | `IsConstOnlyDefeater` |
| `defeater_no_func.go` | `has_no_func.go` | `HasNoFunc` |

### validate_graph.go — 선언적 룰북

```go
func newValidateGraph() *toulmin.Graph {
    g := toulmin.NewGraph("validate")

    // ── 파일 구조: 파일당 하나 ──
    wF1 := g.Warrant(CountMax, &CountMaxBacking{Field: "Funcs", Max: 1, Rule: "F1",
        Message: "file contains multiple funcs; expected 1 file 1 func"}, 1.0)
    wF2 := g.Warrant(CountMax, &CountMaxBacking{Field: "Types", Max: 1, Rule: "F2",
        Message: "file contains multiple types; expected 1 file 1 type"}, 1.0)
    wF3 := g.Warrant(CountMax, &CountMaxBacking{Field: "Methods", Max: 1, Rule: "F3",
        Message: "file contains multiple methods; expected 1 file 1 method"}, 1.0)

    // ── 파일 구조: init 단독 불허 ──
    _ = g.Warrant(ExistsWhen, &ExistsWhenBacking{When: "HasInit", Need: "companion", Rule: "F4",
        Level: "ERROR", Message: "init() must not exist alone; requires accompanying var or func"}, 1.0)

    // ── 코드 품질 ──
    _ = g.Warrant(CheckDepthLimit, nil, 1.0)
    _ = g.Warrant(CheckFuncLines, nil, 1.0)

    // ── 어노테이션: 존재 필수 ──
    wA1f := g.Warrant(ExistsWhen, &ExistsWhenBacking{When: "HasFuncs", Need: "ff:func", Rule: "A1",
        Level: "ERROR", Message: "file with func must have //ff:func annotation"}, 1.0)
    wA1t := g.Warrant(ExistsWhen, &ExistsWhenBacking{When: "HasTypes", Need: "ff:type", Rule: "A1",
        Level: "ERROR", Message: "file with type must have //ff:type annotation"}, 1.0)
    wA3 := g.Warrant(ExistsWhen, &ExistsWhenBacking{When: "HasFuncOrType", Need: "ff:what", Rule: "A3",
        Level: "ERROR", Message: "file with func or type must have //ff:what annotation"}, 1.0)
    wA9 := g.Warrant(ExistsWhen, &ExistsWhenBacking{When: "HasFuncs", Need: "control", Rule: "A9",
        Level: "ERROR", Message: "func file must have control= annotation"}, 1.0)
    wA15 := g.Warrant(ExistsWhen, &ExistsWhenBacking{When: "ControlIteration", Need: "dimension", Rule: "A15",
        Level: "ERROR", Message: "control=iteration requires dimension= annotation"}, 1.0)

    // ── 어노테이션: 제어 구조 일치 ──
    wA10 := g.Warrant(ControlMatch, &ControlMatchBacking{Control: "selection", MustHave: "switch", Rule: "A10",
        Message: "control=selection but no switch found at depth 1"}, 1.0)
    wA11 := g.Warrant(ControlMatch, &ControlMatchBacking{Control: "iteration", MustHave: "loop", Rule: "A11",
        Message: "control=iteration but no loop found at depth 1"}, 1.0)
    wA12 := g.Warrant(ControlMatch, &ControlMatchBacking{Control: "sequence", MustNotHave: "switch|loop", Rule: "A12",
        Message: "control=sequence but %s found at depth 1"}, 1.0)
    wA13 := g.Warrant(ControlMatch, &ControlMatchBacking{Control: "selection", MustNotHave: "loop", Rule: "A13",
        Message: "control=selection but loop found at depth 1"}, 1.0)
    wA14 := g.Warrant(ControlMatch, &ControlMatchBacking{Control: "iteration", MustNotHave: "switch", Rule: "A14",
        Message: "control=iteration but switch found at depth 1"}, 1.0)

    // ── 어노테이션: 코드북 적합 ──
    wA2 := g.Warrant(InCodebook, &InCodebookBacking{Direction: "value→codebook", Rule: "A2"}, 1.0)
    wA8 := g.Warrant(InCodebook, &InCodebookBacking{Direction: "codebook→annotation", Rule: "A8"}, 1.0)

    // ── 고유 검사 ──
    wA6 := g.Warrant(AnnotationAtTop, nil, 1.0)
    wA7 := g.Warrant(CheckedHashMatch, nil, 1.0)
    wA16 := g.Warrant(ValidDimension, nil, 1.0)

    // ══ 예외: test 파일은 구조/어노테이션 룰 면제 ══
    dTest := g.Defeater(IsTestFile, nil, 1.0)
    for _, w := range []*toulmin.Rule{wF1, wF2, wF3, wA1f, wA1t, wA2, wA3, wA6, wA7, wA8,
        wA9, wA10, wA11, wA12, wA13, wA14, wA15, wA16} {
        g.Defeat(dTest, w)
    }

    // ══ 예외: const 전용 파일은 F1 면제 ══
    dConst := g.Defeater(IsConstOnlyDefeater, nil, 1.0)
    g.Defeat(dConst, wF1)

    // ══ 예외: func 없는 파일은 제어 구조 룰 면제 ══
    dNoFunc := g.Defeater(HasNoFunc, nil, 1.0)
    for _, w := range []*toulmin.Rule{wA9, wA10, wA11, wA12, wA13, wA14, wA15, wA16} {
        g.Defeat(dNoFunc, w)
    }

    return g
}
```

### 테스트 전략

`ruleViolations` 헬퍼에 backing 파라미터 추가:

```go
// Before
func ruleViolations(fn func(any, any, any) (bool, any), gf *model.GoFile, cb *model.Codebook) []model.Violation {
    ok, ev := fn(gf.Path, g, nil)

// After
func ruleViolations(fn func(any, any, any) (bool, any), gf *model.GoFile, cb *model.Codebook, backing any) []model.Violation {
    ok, ev := fn(gf.Path, g, backing)
```

개별 테스트는 backing을 명시:

```go
// Before
ruleViolations(RuleF1, mustParse(t, "testdata/multi_func.go"), nil)

// After
ruleViolations(CountMax, mustParse(t, "testdata/multi_func.go"), nil,
    &CountMaxBacking{Field: "Funcs", Max: 1, Rule: "F1", Message: "..."})
```

## 영향 범위

### 삭제 파일 (16개)

| 파일 | 이유 |
|---|---|
| `rule_f1.go` | → CountMax |
| `rule_f2.go` | → CountMax |
| `rule_f3.go` | → CountMax |
| `rule_f4.go` | → ExistsWhen |
| `rule_a1.go` | → ExistsWhen |
| `rule_a2.go` | → InCodebook |
| `rule_a3.go` | → ExistsWhen |
| `rule_a8.go` | → InCodebook |
| `rule_a9.go` | → ExistsWhen |
| `rule_a10.go` | → ControlMatch |
| `rule_a11.go` | → ControlMatch |
| `rule_a12.go` | → ControlMatch |
| `rule_a13.go` | → ControlMatch |
| `rule_a14.go` | → ControlMatch |
| `rule_a15.go` | → ExistsWhen |
| `defeater_const_only.go` | → is_const_only_defeater.go |

### 신규 파일 (8개)

| 파일 | 내용 |
|---|---|
| `exists_when.go` | `ExistsWhen` 함수 |
| `exists_when_backing.go` | `ExistsWhenBacking` 타입 |
| `count_max.go` | `CountMax` 함수 |
| `count_max_backing.go` | `CountMaxBacking` 타입 |
| `control_match.go` | `ControlMatch` 함수 |
| `control_match_backing.go` | `ControlMatchBacking` 타입 |
| `in_codebook.go` | `InCodebook` 함수 |
| `in_codebook_backing.go` | `InCodebookBacking` 타입 |

### 이름 변경 파일 (8개)

| 현재 | 새 이름 |
|---|---|
| `rule_q1.go` | `check_depth_limit.go` |
| `rule_q2q3.go` | `check_func_lines.go` |
| `rule_a6.go` | `annotation_at_top.go` |
| `rule_a7.go` | `checked_hash_match.go` |
| `rule_a16.go` | `valid_dimension.go` |
| `defeater_test_file.go` | `is_test_file.go` |
| `defeater_no_func.go` | `has_no_func.go` |
| `defeater_const_only.go` | `is_const_only_defeater.go` |

### 수정 파일 (2개)

| 파일 | 변경 |
|---|---|
| `validate_graph.go` | 선언적 룰북으로 전면 재작성 |
| `mutest_test.go` | 헬퍼 + 테스트 갱신 |
| `exception_test.go` | 함수 참조 갱신 |

## 구현 순서

### Step 1: backing 타입 생성

1. `exists_when_backing.go`, `count_max_backing.go`, `control_match_backing.go`, `in_codebook_backing.go`

### Step 2: 범용 룰 함수 구현

2. `exists_when.go` — When/Need 문자열 기반 switch 분기
3. `count_max.go` — Field 문자열로 GoFile 필드 참조
4. `control_match.go` — MustHave/MustNotHave 기반 AST 검사 위임
5. `in_codebook.go` — Direction 기반 교차 검증

### Step 3: 고유 룰 이름 변경

6. `rule_q1.go` → `check_depth_limit.go` (함수명 `CheckDepthLimit`)
7. `rule_q2q3.go` → `check_func_lines.go` (함수명 `CheckFuncLines`)
8. `rule_a6.go` → `annotation_at_top.go` (함수명 `AnnotationAtTop`)
9. `rule_a7.go` → `checked_hash_match.go` (함수명 `CheckedHashMatch`)
10. `rule_a16.go` → `valid_dimension.go` (함수명 `ValidDimension`)

### Step 4: defeater 이름 변경

11. `defeater_test_file.go` → `is_test_file.go` (함수명 `IsTestFile`)
12. `defeater_const_only.go` → `is_const_only_defeater.go` (함수명 `IsConstOnlyDefeater`)
13. `defeater_no_func.go` → `has_no_func.go` (함수명 `HasNoFunc`)

### Step 5: validate_graph.go 재작성

14. 선언적 룰북 구조로 전면 재작성

### Step 6: 구 룰 파일 삭제

15. 16개 파일 삭제

### Step 7: 테스트 갱신

16. `mutest_test.go` — ruleViolations 헬퍼에 backing 추가, 테스트 갱신
17. `exception_test.go` — 함수 참조 갱신

### Step 8: 검증

18. `go build ./...`
19. `go vet ./...`
20. `go test ./...`
21. `filefunc validate`

## 파일 수 변화

| 구분 | Before | After |
|---|---|---|
| 룰 함수 파일 | 20 | 9 (4 범용 + 5 고유) |
| backing 타입 파일 | 0 | 4 |
| defeater 파일 | 3 | 3 |
| 기타 (graph, ground, helper 등) | 변경 없음 | 변경 없음 |
| **validate 패키지 총 .go 파일** | **약 40개** | **약 32개** |

룰 함수: 20 → 9 (55% 감소). 파일 총 수: 8개 감소.

## 완료 기준

- `go build ./...` 통과
- `go vet ./...` 통과
- `go test ./...` 전체 통과
- `filefunc validate` ERROR 0
- 구 룰 파일 16개 삭제 완료
- 4개 범용 룰 함수가 backing으로 16개 룰 커버
- validate_graph.go가 위에서 아래로 읽히는 룰북 구조
- 기존 테스트 커버리지 유지 (동일 testdata, 동일 기대 결과)

## 예상 규모

- 삭제 파일: 16개
- 신규 파일: 8개
- 이름 변경: 8개
- 수정 파일: 3개
- 예상 난이도: 중 (범용 함수 내 switch 분기 설계가 핵심, 테스트 갱신 범위 넓음)
- 핵심 난점: ExistsWhen의 When/Need 조합이 6가지 — switch 분기가 깔끔해야 한다
