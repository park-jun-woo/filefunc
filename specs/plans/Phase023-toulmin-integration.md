# Phase 023: toulmin 엔진 도입 — validate를 논증 기반 룰 엔진으로 전환

## 목표

filefunc validate의 룰 실행을 toulmin Graph Builder API로 전환한다. 현재 하드코딩된 예외 분기(F5 test 예외, F6 const 예외)를 defeats 그래프로 선언적으로 표현하고, 향후 프로젝트별 예외 룰 확장을 가능하게 한다.

## 배경

### 현재 구조

```go
// CheckOneFileOneFunc — 예외가 if문으로 하드코딩
func CheckOneFileOneFunc(gf *model.GoFile) []model.Violation {
    if gf.IsTest { return nil }       // F5 예외
    if IsConstOnly(gf) { return nil } // F6 예외
    if len(gf.Funcs) > 1 { return violation }
    return nil
}
```

룰 22개, 예외 로직이 각 Check 함수 내부에 분산. 예외를 추가하려면 소스를 직접 수정해야 한다.

### toulmin 전환 후

```go
// 룰 함수 — 순수 조건만, 예외 분기 없음
//ff:func feature=validate type=rule control=sequence
//ff:what F1: 파일당 func 1개 검증
func CheckOneFileOneFunc(claim any, ground any) bool {
    return len(ground.(*ValidateGround).File.Funcs) > 1
}

// 예외 함수 — 독립 파일
//ff:func feature=validate type=rule control=sequence
//ff:what F5: test 파일은 복수 func 허용
func TestFileException(claim any, ground any) bool {
    return ground.(*ValidateGround).File.IsTest
}

// 그래프 — defeats 관계를 코드로 선언
var ValidateGraph = toulmin.NewGraph("validate").
    Warrant(CheckOneFileOneFunc, 1.0).
    Defeater(TestFileException, 1.0).
    Defeat(TestFileException, CheckOneFileOneFunc)
```

예외가 defeats 관계로 선언적 표현. warrant/defeater가 분리되어 각각 독립 파일. 그래프가 관계를 한눈에 보여줌.

### 도입 이유

1. **10개 오픈소스 실험 대비** — 프로젝트마다 다른 예외 패턴이 나올 수 있음. defeats 그래프에 defeater를 추가하면 기존 warrant 수정 없이 예외 확장 가능.
2. **"Constrain, Then Converge" 논문 강화** — 수렴적 제약이 이진 판정에서 연속 스케일([-1, +1])로 확장 가능함을 실증.
3. **논증 구조의 투명성** — 왜 이 파일이 위반인지, 왜 예외인지가 defeats 그래프에 명시적으로 보임.

## 설계

### 의존성 추가

```
go get github.com/park-jun-woo/toulmin/pkg/toulmin
```

### 룰 시그니처 변환

현재:
```go
func CheckXxx(gf *model.GoFile) []model.Violation
func CheckXxx(gf *model.GoFile, cb *model.Codebook) []model.Violation
```

toulmin:
```go
func(claim any, ground any) bool
```

**claim**: 파일 경로 `string`
**ground**: `*ValidateGround` 구조체

```go
type ValidateGround struct {
    File     *model.GoFile
    Codebook *model.Codebook
}
```

각 Check 함수는 `ground.(*ValidateGround).File`과 `.Codebook`으로 접근. 예외 분기는 제거하고 순수 조건만 남긴다.

### Graph Builder로 그래프 선언

YAML이나 코드 생성 불필요. Graph Builder API가 충분히 간결하다.

```go
// internal/validate/validate_graph.go
var ValidateGraph = toulmin.NewGraph("validate").
    // F 룰 (defeasible — 예외 가능)
    Warrant(CheckOneFileOneFunc, 1.0).
    Warrant(CheckOneFileOneType, 1.0).
    Warrant(CheckOneFileOneMethod, 1.0).
    Warrant(CheckInitStandalone, 1.0).
    // Q 룰
    Warrant(CheckNestingDepth, 1.0).
    Warrant(CheckFuncLines, 1.0).
    // A 룰 (strict — 예외 없음)
    Warrant(CheckAnnotationRequired, 1.0).
    Warrant(CheckCodebookValues, 1.0).
    Warrant(CheckRequiredKeysInAnnotation, 1.0).
    Warrant(CheckWhatRequired, 1.0).
    Warrant(CheckAnnotationPosition, 1.0).
    Warrant(CheckControlRequired, 1.0).
    Warrant(CheckControlSelection, 1.0).
    Warrant(CheckControlIteration, 1.0).
    Warrant(CheckControlSequence, 1.0).
    Warrant(CheckControlSelectionNoLoop, 1.0).
    Warrant(CheckControlIterationNoSwitch, 1.0).
    Warrant(CheckDimensionRequired, 1.0).
    Warrant(CheckDimensionValue, 1.0).
    Warrant(CheckCheckedHash, 1.0).
    // defeater (예외 룰)
    Defeater(TestFileException, 1.0).
    Defeater(ConstOnlyException, 1.0).
    // defeats 관계
    Defeat(TestFileException, CheckOneFileOneFunc).
    Defeat(TestFileException, CheckOneFileOneType).
    Defeat(TestFileException, CheckOneFileOneMethod).
    Defeat(ConstOnlyException, CheckOneFileOneFunc)
```

한 파일에서 전체 defeats 관계가 보인다. 타입 세이프. IDE 자동완성.

### 어노테이션 — `//ff:` 단일 체계

`//rule:` 어노테이션을 사용하지 않는다. `//ff:` 어노테이션만 사용한다.

- `//ff:what` = toulmin의 what
- `//ff:why` = toulmin의 backing

role, qualifier, defeats, strength는 Graph Builder 코드에서 선언하므로 어노테이션에 넣지 않는다. 하나의 약속 체계.

### RunAll 전환

```go
func RunAll(files []*model.GoFile, cb *model.Codebook) []model.Violation {
    var violations []model.Violation
    hasChecked := HasAnyChecked(files)
    for _, gf := range files {
        ground := &ValidateGround{File: gf, Codebook: cb, HasChecked: hasChecked}
        results := ValidateGraph.Evaluate(gf.Path, ground)
        for _, r := range results {
            if r.Verdict > 0 {
                violations = append(violations, toViolation(gf, r))
            }
        }
    }
    return violations
}
```

- `NewEngine()` + `Register()` 대신 `ValidateGraph.Evaluate()` 직접 호출
- verdict > 0 → ERROR, verdict == 0 → 예외 상쇄, verdict < 0 → rebutted

### ValidateGround 구조체

```go
type ValidateGround struct {
    File       *model.GoFile
    Codebook   *model.Codebook
    HasChecked bool
}
```

`HasChecked`는 현재 `HasAnyChecked(files)` 결과. A7 룰에서 사용.

### toViolation 매핑

```go
func toViolation(gf *model.GoFile, r toulmin.EvalResult) model.Violation {
    rule, level, msg := ruleInfo(r.Name)
    if r.Verdict <= 0.5 {
        level = "WARNING"
    }
    return model.Violation{
        File:    gf.Path,
        Rule:    rule,
        Level:   level,
        Message: msg,
    }
}
```

`ruleInfo`는 함수명 → (룰 ID, 기본 level, 메시지) 매핑. 현재 각 Check 함수 내에 하드코딩된 것을 중앙화.

### 룰 분류

| 현재 룰 | Graph Builder 역할 | defeats |
|---|---|---|
| F1 CheckOneFileOneFunc | Warrant | ← TestFileException, ConstOnlyException |
| F2 CheckOneFileOneType | Warrant | ← TestFileException |
| F3 CheckOneFileOneMethod | Warrant | ← TestFileException |
| F4 CheckInitStandalone | Warrant | — |
| Q1 CheckNestingDepth | Warrant | — |
| Q2/Q3 CheckFuncLines | Warrant | — |
| A1~A16 (11개) | Warrant | — |
| C1~C4 | 별도 (ValidateCodebook 유지) | — |
| F5 TestFileException | Defeater | → F1, F2, F3 |
| F6/F7 ConstOnlyException | Defeater | → F1 |

C1~C4(codebook 룰)은 파일 단위가 아닌 codebook 단위 검증이므로 `ValidateCodebook`에 유지. toulmin 그래프에 포함하지 않는다.

## 영향 범위

### 변경 파일

| 파일 | 변경 내용 |
|---|---|
| `go.mod`, `go.sum` | toulmin 의존성 추가 |
| `internal/validate/run_all.go` | Graph Builder 기반 실행으로 전환 |
| `internal/validate/check_one_file_one_func.go` | 시그니처 변환 + 예외 분기 제거 |
| `internal/validate/check_one_file_one_type.go` | 동일 |
| `internal/validate/check_one_file_one_method.go` | 동일 |
| `internal/validate/check_init_standalone.go` | 시그니처 변환 |
| `internal/validate/check_nesting_depth.go` | 시그니처 변환 |
| `internal/validate/check_func_lines.go` | 시그니처 변환 |
| `internal/validate/check_annotation_required.go` | 시그니처 변환 |
| `internal/validate/check_codebook_values.go` | 시그니처 변환 |
| `internal/validate/check_required_keys_in_annotation.go` | 시그니처 변환 |
| `internal/validate/check_what_required.go` | 시그니처 변환 |
| `internal/validate/check_annotation_position.go` | 시그니처 변환 |
| `internal/validate/check_control_required.go` | 시그니처 변환 |
| `internal/validate/check_control_selection.go` | 시그니처 변환 |
| `internal/validate/check_control_iteration.go` | 시그니처 변환 |
| `internal/validate/check_control_sequence.go` | 시그니처 변환 |
| `internal/validate/check_control_selection_no_loop.go` | 시그니처 변환 |
| `internal/validate/check_control_iteration_no_switch.go` | 시그니처 변환 |
| `internal/validate/check_dimension_required.go` | 시그니처 변환 |
| `internal/validate/check_dimension_value.go` | 시그니처 변환 |
| `internal/validate/check_checked_hash.go` | 시그니처 변환 |
| `internal/validate/has_any_checked.go` | 삭제 또는 ValidateGround로 이동 |
| `internal/validate/mutest_test.go` | 테스트 적응 (시그니처 변경) |
| `internal/validate/exception_test.go` | 테스트 적응 |

### 신규 파일

| 파일 | 내용 |
|---|---|
| `internal/validate/validate_ground.go` | ValidateGround 구조체 |
| `internal/validate/validate_graph.go` | Graph Builder 코드 (전체 defeats 관계 선언) |
| `internal/validate/test_file_exception.go` | F5 defeater 함수 |
| `internal/validate/const_only_exception.go` | F6/F7 defeater 함수 |
| `internal/validate/to_violation.go` | EvalResult → Violation 변환 |
| `internal/validate/rule_info.go` | 함수명 → (룰 ID, level, message) 매핑 |

### 변경 없음

- `internal/validate/validate_codebook.go` — C1~C4는 파일 단위가 아니므로 기존 유지
- `internal/validate/is_const_only.go` — defeater에서 호출, 로직 유지
- `internal/validate/contains.go` — 유틸리티, 변경 없음
- `internal/validate/allowed_values.go` — codebook 유틸리티, 변경 없음
- `internal/validate/q3_limit.go` — Q3 유틸리티, 변경 없음
- `internal/validate/depth_limit.go` — Q1 유틸리티, 변경 없음
- `internal/validate/has_backtick.go` — Q3 유틸리티, 변경 없음
- `internal/validate/detect_control.go` — 제어구조 판별, 변경 없음

## 구현 순서

1. `go get github.com/park-jun-woo/toulmin/pkg/toulmin`
2. `validate_ground.go` — ValidateGround 구조체 정의
3. `to_violation.go` + `rule_info.go` — EvalResult → Violation 변환
4. warrant 룰 변환 (19개): 시그니처 `func(claim any, ground any) bool`로 변경, 예외 분기 제거
5. defeater 룰 신규 작성: `test_file_exception.go`, `const_only_exception.go`
6. `validate_graph.go` — Graph Builder로 전체 그래프 선언
7. `run_all.go` — ValidateGraph.Evaluate 기반으로 전환
8. `mutest_test.go`, `exception_test.go` — 테스트 적응
9. `go build ./...` + `go test ./...`
10. `filefunc validate` ERROR 0 확인
11. README, CLAUDE.md 업데이트

## 완료 기준

- `go build ./...` 통과
- `go test ./...` 통과 (mutest 전체)
- `filefunc validate` ERROR 0 (filefunc 자체)
- ValidateGraph에서 전체 defeats 관계가 보임
- 기존 22개 룰 동일 동작 (이진 판정 호환)
- F5/F7 예외가 defeats 그래프로 표현

## 하위 호환성

- `filefunc validate` 출력 형식 변경 없음 (ERROR/WARNING 동일)
- verdict > 0 → ERROR, verdict == 0 → 예외 상쇄 (기존과 동일 동작)
- 사용자 관점에서 변경 없음. 내부 엔진만 교체.

## 예상 규모

- 변경 파일: ~27개 (validate 룰 19개 + run_all + 테스트 2개 + go.mod 등)
- 신규 파일: 6개
- 예상 난이도: 중 (시그니처 변환은 기계적, 핵심은 ground 설계와 테스트 적응)
- 핵심 난점: mutest 테스트 적응 — 기존 테스트가 직접 Check 함수를 호출하므로 시그니처 변경에 맞춰 수정 필요
