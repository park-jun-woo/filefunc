# Phase 024: Evidence 반환 — Rule 함수가 판정과 증거를 동시에 반환 ✅

## 목표

toulmin의 확장된 rule 시그니처 `func(any, any) (bool, any)`를 적용하여, Rule 함수가 판정과 증거(Violation)를 동시에 반환하게 한다. 이중 실행(Rule 판정 + Check 메시지)을 제거한다.

## 배경

### 현재 구조 (Phase023)

```
1. RuleF1(claim, ground) → true (판정)
2. verdict > 0 → violationsFor("RuleF1", gf, cb) → CheckOneFileOneFunc(gf) 재호출
3. CheckOneFileOneFunc → Violation{Rule:"F1", Message:"..."} (메시지)
```

동일 로직이 2번 실행됨. `violationsFor`가 20개 룰을 switch로 디스패치.

### Phase024 후

```
1. RuleF1(claim, ground) → (true, []Violation{{Rule:"F1", Message:"..."}})
2. verdict > 0 → EvalResult.Evidence에서 Violation 꺼냄
```

1번 실행으로 판정 + 증거. `violationsFor` 삭제. Check 함수 삭제 가능.

## 설계

### toulmin 의존성 업데이트

```bash
go get github.com/park-jun-woo/toulmin@84281b8
```

rule 시그니처: `func(any, any) (bool, any)` — toulmin 47ecf49에서 확장됨.

### Rule 함수 변경 패턴

기존 Check 함수의 로직을 Rule 함수에 통합. Evidence = `[]model.Violation`.

```go
// Before (rule_f1.go)
func RuleF1(claim any, ground any) bool {
    gf := ground.(*ValidateGround).File
    return len(gf.Funcs) > 1
}

// After (rule_f1.go)
func RuleF1(claim any, ground any) (bool, any) {
    gf := ground.(*ValidateGround).File
    if len(gf.Funcs) > 1 {
        return true, []model.Violation{{
            File: gf.Path, Rule: "F1", Level: "ERROR",
            Message: "file contains multiple funcs; expected 1 file 1 func",
        }}
    }
    return false, nil
}
```

Defeater 함수: evidence 불필요 → `return true, nil` / `return false, nil`

```go
func DefeaterTestFile(claim any, ground any) (bool, any) {
    return ground.(*ValidateGround).File.IsTest, nil
}
```

### evaluateFile 변경

```go
func evaluateFile(gf *model.GoFile, cb *model.Codebook, ground *ValidateGround) []model.Violation {
    results := ValidateGraph.Evaluate(gf.Path, ground)
    var violations []model.Violation
    for _, r := range results {
        if r.Verdict > 0 && r.Evidence != nil {
            violations = append(violations, r.Evidence.([]model.Violation)...)
        }
    }
    return violations
}
```

`violationsFor` 호출 제거. Evidence에서 직접 Violation 추출.

### 삭제 대상

| 파일 | 이유 |
|---|---|
| `internal/validate/violations_for.go` | evidence로 대체 — 디스패치 불필요 |
| `internal/validate/check_one_file_one_func.go` | RuleF1에 통합 |
| `internal/validate/check_one_file_one_type.go` | RuleF2에 통합 |
| `internal/validate/check_one_file_one_method.go` | RuleF3에 통합 |
| `internal/validate/check_init_standalone.go` | RuleF4에 통합 |
| `internal/validate/check_nesting_depth.go` | RuleQ1에 통합 |
| `internal/validate/check_func_lines.go` | RuleQ2Q3에 통합 |
| `internal/validate/check_annotation_required.go` | RuleA1에 통합 |
| `internal/validate/check_codebook_values.go` | RuleA2에 통합 |
| `internal/validate/check_what_required.go` | RuleA3에 통합 |
| `internal/validate/check_annotation_position.go` | RuleA6에 통합 |
| `internal/validate/check_checked_hash.go` | RuleA7에 통합 |
| `internal/validate/check_required_keys_in_annotation.go` | RuleA8에 통합 |
| `internal/validate/check_control_required.go` | RuleA9에 통합 |
| `internal/validate/check_control_selection.go` | RuleA10에 통합 |
| `internal/validate/check_control_iteration.go` | RuleA11에 통합 |
| `internal/validate/check_control_sequence.go` | RuleA12에 통합 |
| `internal/validate/check_control_selection_no_loop.go` | RuleA13에 통합 |
| `internal/validate/check_control_iteration_no_switch.go` | RuleA14에 통합 |
| `internal/validate/check_dimension_required.go` | RuleA15에 통합 |
| `internal/validate/check_dimension_value.go` | RuleA16에 통합 |

21파일 삭제. Rule 함수가 Check 함수를 대체.

### 유지 대상

| 파일 | 이유 |
|---|---|
| `is_const_only.go` | DefeaterConstOnly에서 호출. 유틸리티. |
| `allowed_values.go` | RuleA2에서 호출. 유틸리티. |
| `contains.go` | RuleA2에서 호출. 유틸리티. |
| `q3_limit.go` | RuleQ2Q3에서 호출. 유틸리티. |
| `depth_limit.go` | RuleQ1에서 호출. 유틸리티. |
| `has_backtick.go` | RuleQ2Q3에서 호출. 유틸리티. |
| `detect_control.go` | RuleA10~A14에서 호출. 유틸리티. |
| `has_loop_at_depth1.go` | RuleA11, A13에서 호출. 유틸리티. |
| `has_switch_at_depth1.go` | RuleA10, A14에서 호출. 유틸리티. |
| `has_any_checked.go` | RunAll에서 호출. 유틸리티. |
| `validate_codebook.go` | C1~C4는 별도 경로. |
| `find_duplicates.go` | C2에서 호출. |
| `check_codebook_*.go` (C1~C4) | codebook 검증은 toulmin 밖. |
| `mutest_test.go` | 테스트 — 시그니처 변경 대응 |
| `exception_test.go` | 테스트 — 이미 RunAll 경유 |

## 영향 범위

### 변경 파일

| 파일 | 변경 내용 |
|---|---|
| `go.mod`, `go.sum` | toulmin 최신 버전 |
| `internal/validate/rule_f1.go` ~ `rule_a16.go` (20개) | 시그니처 `(bool, any)` + Check 로직 통합 |
| `internal/validate/defeater_test_file.go` | 시그니처 `(bool, any)` |
| `internal/validate/defeater_const_only.go` | 시그니처 `(bool, any)` |
| `internal/validate/defeater_no_func.go` | 시그니처 `(bool, any)` |
| `internal/validate/evaluate_file.go` | Evidence에서 Violation 추출 |
| `internal/validate/mutest_test.go` | Check → Rule 함수 호출로 변경 |

### 삭제 파일

21개 (violations_for.go + check_*.go 20개)

### 신규 파일

0개

## 구현 순서

1. `go get github.com/park-jun-woo/toulmin@84281b8`
2. Defeater 3개: 시그니처 `(bool, any)` 변경 (evidence nil)
3. Rule 20개: 시그니처 `(bool, any)` + Check 로직 통합 (evidence = []Violation)
4. `evaluate_file.go`: Evidence에서 Violation 추출
5. `mutest_test.go`: Check → Rule 함수 호출로 변경
6. Check 함수 20개 + violations_for.go 삭제 (21파일)
7. `go build ./...` + `go test ./...` + `filefunc validate`

## 완료 기준

- `go build ./...` 통과
- `go test ./...` 통과 (mutest 27개 전체)
- `filefunc validate` ERROR 0
- Check 함수 0개 (전부 Rule 함수로 대체)
- `violations_for.go` 삭제
- 이중 실행 제거 확인

## 예상 규모

- 변경 파일: ~26개
- 삭제 파일: 21개
- 신규 파일: 0개
- 예상 난이도: 중 (기계적 시그니처 변환 + 로직 복사)
- 핵심 난점: CheckFuncLines(Q2/Q3)가 다중 Violation을 반환 — Rule 함수에서 동일하게 처리
