# Phase 029: toulmin backing 마이그레이션 — 3-arg 시그니처 + Graph Builder API 변경

## 목표

toulmin 커밋 444f4d3 (`feat: backing as first-class Toulmin element`) 반영. (1) 룰 함수 시그니처를 2-arg에서 3-arg로 변경 (2) Graph Builder API 호출을 `(fn, backing, qualifier)` 형식으로 변경 (3) 테스트 헬퍼 시그니처 갱신 (4) go.mod 의존성 업데이트.

## 배경

### toulmin 변경 사항 (444f4d3)

| 항목 | Before | After |
|---|---|---|
| 룰 시그니처 | `func(claim any, ground any) (bool, any)` | `func(claim any, ground any, backing any) (bool, any)` |
| Graph Builder API | `Warrant(fn, qualifier)` | `Warrant(fn, backing, qualifier)` |
| Defeat 식별 | funcID 기반 | ruleID = funcID + "#" + backing |
| 신규 API | — | `DefeatWith(fromFn, fromBacking, toFn, toBacking)` |
| TraceEntry | Name, Role, Activated, Qualifier, Evidence | + Backing 필드 추가 |
| 레거시 지원 | — | 2-arg 함수는 `wrapLegacy`로 자동 래핑 (동작은 하지만 권장하지 않음) |

### filefunc 현재 상태

- 룰 함수 23개: 모두 2-arg 시그니처
- Graph Builder 호출: `Warrant(fn, qualifier)`, `Defeater(fn, qualifier)` 형식 (backing 파라미터 없음)
- 테스트 헬퍼 1개: `ruleViolations(fn func(any, any) (bool, any), ...)` 2-arg 타입
- filefunc은 backing을 사용하지 않으므로 모든 backing은 `nil`

## 이슈 목록

### 1. Graph Builder API 호출 변경 (CRITICAL)

| # | 파일 | 이슈 |
|---|---|---|
| G1 | `internal/validate/validate_graph.go` | `Warrant(fn, 1.0)` → `Warrant(fn, nil, 1.0)`. 20개 Warrant + 3개 Defeater = 23개 호출 |

### 2. 룰 함수 시그니처 변경 (MEDIUM)

레거시 2-arg도 동작하지만, 공식 시그니처에 맞추어 3-arg로 마이그레이션.

| # | 파일 | 함수 |
|---|---|---|
| R1 | `internal/validate/rule_f1.go` | `RuleF1(claim any, ground any) (bool, any)` |
| R2 | `internal/validate/rule_f2.go` | `RuleF2` |
| R3 | `internal/validate/rule_f3.go` | `RuleF3` |
| R4 | `internal/validate/rule_f4.go` | `RuleF4` |
| R5 | `internal/validate/rule_q1.go` | `RuleQ1` |
| R6 | `internal/validate/rule_q2q3.go` | `RuleQ2Q3` |
| R7 | `internal/validate/rule_a1.go` | `RuleA1` |
| R8 | `internal/validate/rule_a2.go` | `RuleA2` |
| R9 | `internal/validate/rule_a3.go` | `RuleA3` |
| R10 | `internal/validate/rule_a6.go` | `RuleA6` |
| R11 | `internal/validate/rule_a7.go` | `RuleA7` |
| R12 | `internal/validate/rule_a8.go` | `RuleA8` |
| R13 | `internal/validate/rule_a9.go` | `RuleA9` |
| R14 | `internal/validate/rule_a10.go` | `RuleA10` |
| R15 | `internal/validate/rule_a11.go` | `RuleA11` |
| R16 | `internal/validate/rule_a12.go` | `RuleA12` |
| R17 | `internal/validate/rule_a13.go` | `RuleA13` |
| R18 | `internal/validate/rule_a14.go` | `RuleA14` |
| R19 | `internal/validate/rule_a15.go` | `RuleA15` |
| R20 | `internal/validate/rule_a16.go` | `RuleA16` |
| R21 | `internal/validate/defeater_test_file.go` | `DefeaterTestFile` |
| R22 | `internal/validate/defeater_const_only.go` | `DefeaterConstOnly` |
| R23 | `internal/validate/defeater_no_func.go` | `DefeaterNoFunc` |

### 3. 테스트 헬퍼 시그니처 변경 (MEDIUM)

| # | 파일 | 이슈 |
|---|---|---|
| T1 | `internal/validate/mutest_test.go` | `ruleViolations(fn func(any, any) (bool, any), ...)` → `fn func(any, any, any) (bool, any)`. 함수 내 호출도 `fn(claim, ground)` → `fn(claim, ground, nil)` |

### 4. 의존성 업데이트

| # | 파일 | 이슈 |
|---|---|---|
| D1 | `go.mod` | `github.com/park-jun-woo/toulmin` 버전을 444f4d3 커밋으로 업데이트 |

## 설계

### G1: validate_graph.go — backing 파라미터 추가

```go
// Before
Warrant(RuleF1, 1.0).
Defeater(DefeaterTestFile, 1.0).

// After
Warrant(RuleF1, nil, 1.0).
Defeater(DefeaterTestFile, nil, 1.0).
```

filefunc은 backing을 사용하지 않으므로 모든 backing은 `nil`. Defeat 호출은 변경 없음 (같은 함수의 backing 구분이 불필요).

### R1-R23: 룰 함수 시그니처 변경

```go
// Before
func RuleF1(claim any, ground any) (bool, any) {

// After
func RuleF1(claim any, ground any, backing any) (bool, any) {
```

함수 본문은 backing을 사용하지 않으므로 변경 없음. 시그니처만 변경.

### T1: mutest_test.go 헬퍼 변경

```go
// Before
func ruleViolations(fn func(any, any) (bool, any), gf *model.GoFile, cb *model.Codebook) []model.Violation {
    // ...
    violated, evidence := fn(gf.Path, ground)

// After
func ruleViolations(fn func(any, any, any) (bool, any), gf *model.GoFile, cb *model.Codebook) []model.Violation {
    // ...
    violated, evidence := fn(gf.Path, ground, nil)
```

### D1: go.mod 업데이트

```bash
go get github.com/park-jun-woo/toulmin@444f4d3
go mod tidy
```

## 영향 범위

### 변경 파일

| 파일 | 변경 내용 |
|---|---|
| `go.mod` | toulmin 의존성 버전 업데이트 |
| `go.sum` | 자동 갱신 |
| `internal/validate/validate_graph.go` | Warrant/Defeater에 backing=nil 추가 (23개 호출) |
| `internal/validate/rule_f1.go` | 시그니처 3-arg |
| `internal/validate/rule_f2.go` | 시그니처 3-arg |
| `internal/validate/rule_f3.go` | 시그니처 3-arg |
| `internal/validate/rule_f4.go` | 시그니처 3-arg |
| `internal/validate/rule_q1.go` | 시그니처 3-arg |
| `internal/validate/rule_q2q3.go` | 시그니처 3-arg |
| `internal/validate/rule_a1.go` | 시그니처 3-arg |
| `internal/validate/rule_a2.go` | 시그니처 3-arg |
| `internal/validate/rule_a3.go` | 시그니처 3-arg |
| `internal/validate/rule_a6.go` | 시그니처 3-arg |
| `internal/validate/rule_a7.go` | 시그니처 3-arg |
| `internal/validate/rule_a8.go` | 시그니처 3-arg |
| `internal/validate/rule_a9.go` | 시그니처 3-arg |
| `internal/validate/rule_a10.go` | 시그니처 3-arg |
| `internal/validate/rule_a11.go` | 시그니처 3-arg |
| `internal/validate/rule_a12.go` | 시그니처 3-arg |
| `internal/validate/rule_a13.go` | 시그니처 3-arg |
| `internal/validate/rule_a14.go` | 시그니처 3-arg |
| `internal/validate/rule_a15.go` | 시그니처 3-arg |
| `internal/validate/rule_a16.go` | 시그니처 3-arg |
| `internal/validate/defeater_test_file.go` | 시그니처 3-arg |
| `internal/validate/defeater_const_only.go` | 시그니처 3-arg |
| `internal/validate/defeater_no_func.go` | 시그니처 3-arg |
| `internal/validate/mutest_test.go` | 헬퍼 시그니처 + 호출 변경 |

### 삭제/신규 파일

없음.

## 구현 순서

### Step 1: 의존성 업데이트 (D1)

1. `go get github.com/park-jun-woo/toulmin@444f4d3`
2. `go mod tidy`

### Step 2: Graph Builder API 변경 (G1)

3. `internal/validate/validate_graph.go` — 23개 Warrant/Defeater 호출에 backing=nil 추가

### Step 3: 룰 함수 시그니처 변경 (R1-R23)

4. 20개 룰 함수 시그니처에 `backing any` 파라미터 추가
5. 3개 defeater 함수 시그니처에 `backing any` 파라미터 추가

### Step 4: 테스트 헬퍼 변경 (T1)

6. `internal/validate/mutest_test.go` — 헬퍼 시그니처 + 호출 변경

### Step 5: 검증

7. `go build ./...`
8. `go vet ./...`
9. `go test ./...`
10. `filefunc validate`

## 완료 기준

- `go build ./...` 통과
- `go vet ./...` 통과
- `go test ./...` 전체 통과
- `filefunc validate` ERROR 0
- 모든 룰 함수가 3-arg 시그니처 사용
- Graph Builder 호출이 `(fn, nil, 1.0)` 형식

## 예상 규모

- 변경 파일: 27개 (go.mod, go.sum 포함)
- 삭제/신규 파일: 0개
- 예상 난이도: 하 (기계적 시그니처 변경)
- 핵심 난점: 없음. 모든 변경이 패턴화되어 있음
