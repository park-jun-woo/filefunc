# Phase 030: toulmin Rule reference API 마이그레이션 — 체이닝 제거, Defeat(*Rule, *Rule)

## 목표

toulmin 커밋 b1b78c3 (`refactor: Rule reference API, remove chaining, GraphBuilder → Graph`) 반영. (1) 체이닝 구문을 개별 문장으로 변환 (2) `Defeat(fn, fn)` → `Defeat(*Rule, *Rule)` (3) go.mod 의존성 업데이트.

## 배경

### toulmin 변경 사항 (b1b78c3)

| 항목 | Before | After |
|---|---|---|
| 타입명 | `GraphBuilder` | `Graph` |
| Warrant/Rebuttal/Defeater 반환 | `*GraphBuilder` (체이닝) | `*Rule` (참조) |
| Defeat 파라미터 | `Defeat(fromFn, toFn)` — 함수 참조 | `Defeat(*Rule, *Rule)` — Rule 참조 |
| DefeatWith | 존재 (`fromFn, fromBacking, toFn, toBacking`) | 제거 — `*Rule` 참조로 대체 |
| 그래프 선언 | 체이닝: `NewGraph().Warrant().Warrant()...` | 개별 문장: `g.Warrant(...)`, `g.Defeat(r, w)` |
| 패키지 구조 | `graph_builder.go`, `graph_builder_evaluate.go`, `graph_builder_defeat_with.go` | `graph.go`, `graph_warrant.go`, `graph_rebuttal.go`, `graph_defeater.go`, `graph_defeat.go`, `graph_evaluate.go`, `rule.go` |

### filefunc 현재 상태

- `validate_graph.go`: 체이닝 방식으로 20 Warrant + 3 Defeater + 37 Defeat 선언
- `evaluate_file.go`: `ValidateGraph.Evaluate(claim, ground)` — 변경 불필요 (메서드 시그니처 동일)
- 룰 함수 시그니처: Phase029에서 3-arg로 이미 마이그레이션 완료

## 이슈 목록

### 1. validate_graph.go 체이닝 → 개별 문장 (CRITICAL)

| # | 이슈 |
|---|---|
| G1 | `var ValidateGraph = toulmin.NewGraph("validate").Warrant(...).Warrant(...)...` 체이닝을 개별 문장으로 분리 |
| G2 | `Defeat(DefeaterTestFile, RuleF1)` — 함수 참조를 `*Rule` 변수 참조로 변경 |

### 2. 의존성 업데이트

| # | 파일 | 이슈 |
|---|---|---|
| D1 | `go.mod` | `github.com/park-jun-woo/toulmin` 버전을 b1b78c3 커밋으로 업데이트 |

## 설계

### G1+G2: validate_graph.go — 전면 재작성

```go
//ff:func feature=validate type=command control=sequence
//ff:what toulmin defeats graph — 전체 validate 룰과 예외 관계를 선언
package validate

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// ValidateGraph declares all validation rules and their defeat relationships.
var ValidateGraph = newValidateGraph()

func newValidateGraph() *toulmin.Graph {
	g := toulmin.NewGraph("validate")
	// F rules (file structure)
	wF1 := g.Warrant(RuleF1, nil, 1.0)
	wF2 := g.Warrant(RuleF2, nil, 1.0)
	wF3 := g.Warrant(RuleF3, nil, 1.0)
	_ = g.Warrant(RuleF4, nil, 1.0)
	// Q rules (code quality)
	_ = g.Warrant(RuleQ1, nil, 1.0)
	_ = g.Warrant(RuleQ2Q3, nil, 1.0)
	// A rules (annotation)
	wA1 := g.Warrant(RuleA1, nil, 1.0)
	wA2 := g.Warrant(RuleA2, nil, 1.0)
	wA3 := g.Warrant(RuleA3, nil, 1.0)
	wA6 := g.Warrant(RuleA6, nil, 1.0)
	wA7 := g.Warrant(RuleA7, nil, 1.0)
	wA8 := g.Warrant(RuleA8, nil, 1.0)
	wA9 := g.Warrant(RuleA9, nil, 1.0)
	wA10 := g.Warrant(RuleA10, nil, 1.0)
	wA11 := g.Warrant(RuleA11, nil, 1.0)
	wA12 := g.Warrant(RuleA12, nil, 1.0)
	wA13 := g.Warrant(RuleA13, nil, 1.0)
	wA14 := g.Warrant(RuleA14, nil, 1.0)
	wA15 := g.Warrant(RuleA15, nil, 1.0)
	wA16 := g.Warrant(RuleA16, nil, 1.0)
	// Defeaters
	dTestFile := g.Defeater(DefeaterTestFile, nil, 1.0)
	dConstOnly := g.Defeater(DefeaterConstOnly, nil, 1.0)
	dNoFunc := g.Defeater(DefeaterNoFunc, nil, 1.0)
	// Defeat edges: test files defeat F/A rules
	g.Defeat(dTestFile, wF1)
	g.Defeat(dTestFile, wF2)
	g.Defeat(dTestFile, wF3)
	g.Defeat(dTestFile, wA1)
	g.Defeat(dTestFile, wA2)
	g.Defeat(dTestFile, wA3)
	g.Defeat(dTestFile, wA6)
	g.Defeat(dTestFile, wA7)
	g.Defeat(dTestFile, wA8)
	g.Defeat(dTestFile, wA9)
	g.Defeat(dTestFile, wA10)
	g.Defeat(dTestFile, wA11)
	g.Defeat(dTestFile, wA12)
	g.Defeat(dTestFile, wA13)
	g.Defeat(dTestFile, wA14)
	g.Defeat(dTestFile, wA15)
	g.Defeat(dTestFile, wA16)
	// Defeat edges: const-only files defeat F1
	g.Defeat(dConstOnly, wF1)
	// Defeat edges: no-func files defeat control/dimension/annotation rules
	g.Defeat(dNoFunc, wA9)
	g.Defeat(dNoFunc, wA10)
	g.Defeat(dNoFunc, wA11)
	g.Defeat(dNoFunc, wA12)
	g.Defeat(dNoFunc, wA13)
	g.Defeat(dNoFunc, wA14)
	g.Defeat(dNoFunc, wA15)
	g.Defeat(dNoFunc, wA16)
	return g
}
```

핵심 변경:
- `var ValidateGraph = ...` 체이닝 → `newValidateGraph()` 팩토리 함수
- 각 `Warrant`/`Defeater` 반환값을 `*Rule` 변수로 캡처
- `Defeat(fn, fn)` → `Defeat(*Rule, *Rule)`
- defeat 대상이 아닌 warrant는 `_ =`로 반환값 무시 (RuleF4, RuleQ1, RuleQ2Q3)

### evaluate_file.go — 변경 없음

`ValidateGraph.Evaluate(claim, ground)` 호출은 그대로 유효. `Graph` 타입에 동일한 `Evaluate` 메서드 존재.

## 영향 범위

### 변경 파일

| 파일 | 변경 내용 |
|---|---|
| `go.mod` | toulmin 의존성 버전 업데이트 |
| `go.sum` | 자동 갱신 |
| `internal/validate/validate_graph.go` | 체이닝 → 개별 문장, Defeat(*Rule, *Rule) |

### 삭제/신규 파일

없음.

## 구현 순서

### Step 1: 의존성 업데이트 (D1)

1. `go get github.com/park-jun-woo/toulmin@b1b78c3`
2. `go mod tidy`

### Step 2: validate_graph.go 재작성 (G1+G2)

3. 체이닝 구문을 팩토리 함수로 변환

### Step 3: 검증

4. `go build ./...`
5. `go vet ./...`
6. `go test ./...`
7. `filefunc validate`

## 완료 기준

- `go build ./...` 통과
- `go vet ./...` 통과
- `go test ./...` 전체 통과
- `filefunc validate` ERROR 0
- 체이닝 구문 완전 제거
- 모든 Defeat 호출이 `*Rule` 참조 사용

## 예상 규모

- 변경 파일: 3개 (go.mod, go.sum, validate_graph.go)
- 삭제/신규 파일: 0개
- 예상 난이도: 하 (구조적 변환, 로직 변경 없음)
- 핵심 난점: 없음
