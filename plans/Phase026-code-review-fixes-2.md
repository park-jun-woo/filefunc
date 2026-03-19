# Phase 026: 코드리뷰 조치 2차 — 버그 수정, 로직 보강, 테스트 확충 ✅

## 목표

2차 코드리뷰에서 발견된 이슈를 전부 조치한다: (1) 명확한 버그 2건 (2) 로직 결함 2건 (3) 안전성 보강 3건 (4) 테스트 커버리지 확충.

## 배경

### 현재 상태

- `go build`, `go vet`, `go test ./...`, `filefunc validate` 모두 통과
- Phase025에서 에러 처리, 중복 제거, else depth 버그 등 조치 완료
- 테스트 있는 패키지: 9개 (parse, validate, walk, chain, annotate, report, context, llm, cli)

### 리뷰 결과 요약

| 분류 | 건수 | 심각도 |
|---|---|---|
| 버그 (BUG) | 2 | HIGH |
| 로직 결함 (LOGIC) | 2 | MEDIUM |
| 안전성 보강 (SAFETY) | 3 | LOW |
| 테스트 확충 | 미테스트 핵심 로직 다수 | MEDIUM |

## 이슈 목록

### 1. 버그 (HIGH)

| # | 파일 | 이슈 |
|---|---|---|
| B1 | `internal/cli/context.go:43` | `ParseFFIgnore(root)` — 다른 3곳은 `root + "/.ffignore"` 전달. context 명령만 ignore 패턴 미적용 |
| B2 | `internal/cli/process_llmc_file.go:48` | `WriteAnnotationLine` 반환값 `(bool, error)` 무시. 쓰기 실패 시에도 `[PASS]` 출력 |

### 2. 로직 결함 (MEDIUM)

| # | 파일 | 이슈 |
|---|---|---|
| L1 | `internal/validate/rule_a12.go:19` | `control == ""` 일 때도 A12 발동. control 미지정은 A9가 담당하므로 A12는 `control == "sequence"`일 때만 발동해야 함 |
| L2 | `internal/validate/rule_a1.go:11` | `hasTypes := len(gf.Types) > 0 && !hasFuncs` — func+type 공존 파일에서 `//ff:type` 검사 누락. F1/F2 defeater 억제 시 A1이 type 검사를 건너뜀 |

### 3. 안전성 보강 (LOW)

| # | 파일 | 이슈 |
|---|---|---|
| S1 | `internal/validate/evaluate_file.go:14` | `r.Evidence.([]model.Violation)` 타입 단언 — 새 rule이 잘못된 타입 반환 시 패닉 |
| S2 | `internal/parse/read_module_path.go:27` | `go.mod`에 `module` 행 없으면 `("", nil)` 반환 → 모든 import가 프로젝트 내부로 판단됨 |
| S3 | `internal/chain/build_call_graph.go:28` | `ExtractCalls` 에러 시 `continue`로 무시 — 불완전 그래프 원인 디버깅 불가 |

### 4. 테스트 커버리지 확충

Phase025에서 순수 함수 단위 테스트를 추가했으나, 핵심 로직의 통합 테스트가 부족.

| 우선순위 | 패키지 | 미테스트 핵심 로직 |
|---|---|---|
| HIGH | chain | `BuildCallGraph`, `TraverseChon`, `TraverseDepth`, `FindSiblings`, `ExpandThrough` |
| HIGH | parse | `ExtractCalls`, `DetectControl`, `CalcMaxDepth`, `CallName` |
| MEDIUM | cli | `ResolveTarget`, `FindMatchingFuncs`, `FindGoMod` |
| MEDIUM | walk | `MatchFFIgnore`, `MatchPattern` |

## 설계

### B1: context.go ParseFFIgnore 경로 수정

```go
// Before
ignorePatterns := walk.ParseFFIgnore(root)

// After
ignorePatterns := walk.ParseFFIgnore(root + "/.ffignore")
```

### B2: process_llmc_file.go WriteAnnotationLine 에러 처리

```go
// Before
annotate.WriteAnnotationLine(gf.Path, "checked", fmt.Sprintf("llm=%s hash=%s", modelName, currentHash))
fmt.Printf("[PASS] %s: score=%.2f\n", gf.Path, score)
return "pass"

// After
_, err = annotate.WriteAnnotationLine(gf.Path, "checked", fmt.Sprintf("llm=%s hash=%s", modelName, currentHash))
if err != nil {
    fmt.Fprintf(os.Stderr, "[ERROR] %s: failed to write checked annotation: %v\n", gf.Path, err)
    return "fail"
}
fmt.Printf("[PASS] %s: score=%.2f\n", gf.Path, score)
return "pass"
```

### L1: rule_a12.go control 빈값 처리

```go
// Before
control := gf.Annotation.Func["control"]
if control != "" && control != "sequence" {
    return false, nil
}

// After
control := gf.Annotation.Func["control"]
if control != "sequence" {
    return false, nil
}
```

`control == ""` 이면 A9이 검출. A12는 `control == "sequence"`로 선언된 파일만 검증.

### L2: rule_a1.go func+type 공존 처리

```go
// Before
hasFuncs := len(gf.Funcs) > 0
hasTypes := len(gf.Types) > 0 && !hasFuncs

// After
hasFuncs := len(gf.Funcs) > 0
hasTypes := len(gf.Types) > 0
```

func이 있든 없든 type이 있으면 `//ff:type` 검사 수행. 단, F1/F2가 둘 다 있는 파일을 ERROR로 잡으므로, 실제로 func+type 공존은 사전에 걸러짐. A1은 방어적으로 양쪽 모두 검사.

### S1: evaluate_file.go 안전한 타입 단언

```go
// Before
violations = append(violations, r.Evidence.([]model.Violation)...)

// After
if vs, ok := r.Evidence.([]model.Violation); ok {
    violations = append(violations, vs...)
}
```

### S2: read_module_path.go 빈 모듈 경로 에러

```go
// Before
return "", scanner.Err()

// After
if err := scanner.Err(); err != nil {
    return "", err
}
return "", fmt.Errorf("module directive not found in %s", goModPath)
```

### S3: build_call_graph.go 에러 로깅

```go
// Before
calls, err := parse.ExtractCalls(gf.Path, modulePath, projFuncs, gf.Package)
if err != nil {
    continue
}

// After
calls, err := parse.ExtractCalls(gf.Path, modulePath, projFuncs, gf.Package)
if err != nil {
    fmt.Fprintf(os.Stderr, "[WARN] %s: extract calls failed: %v\n", gf.Path, err)
    continue
}
```

### 테스트 전략

testdata 디렉토리에 최소한의 Go 소스 파일을 배치하여 통합 테스트 수행.

#### chain 패키지 테스트 추가

- `BuildCallGraph`: 3개 파일(A→B→C) testdata로 그래프 구성 → Children/Parents 검증
- `TraverseChon`: 구성된 그래프에서 chon=1, chon=2 결과 검증
- `TraverseDepth`: child-depth, parent-depth 방향별 결과 검증
- `FindSiblings`: 같은 부모 호출 형제 검증
- `ExpandThrough`: 중간 노드 경유 확장 검증

#### parse 패키지 테스트 추가

- `CallName`: SelectorExpr(pkg.Func), Ident(samePackage), 미등록 함수, 비-Ident receiver 케이스
- `DetectControl`: sequence/selection/iteration 각각의 testdata 파일
- `CalcMaxDepth`: 중첩 깊이별 testdata 파일
- `ExtractCalls`: 프로젝트 내부/외부 호출 구분 검증

#### walk 패키지 테스트 추가

- `MatchFFIgnore`: 패턴 매칭 케이스 (glob, directory, negation)
- `MatchPattern`: 단일 패턴 매칭 검증

#### cli 패키지 테스트 추가

- `FindGoMod`: go.mod 탐색 (존재/부재)
- `ResolveTarget`: qualified/unqualified name 해석

## 영향 범위

### 변경 파일

| 파일 | 변경 내용 |
|---|---|
| `internal/cli/context.go` | B1: `.ffignore` 경로 수정 |
| `internal/cli/process_llmc_file.go` | B2: WriteAnnotationLine 에러 처리 |
| `internal/validate/rule_a12.go` | L1: control 빈값 제외 |
| `internal/validate/rule_a1.go` | L2: func+type 독립 검사 |
| `internal/validate/evaluate_file.go` | S1: 안전한 타입 단언 |
| `internal/parse/read_module_path.go` | S2: 빈 모듈 경로 에러 |
| `internal/chain/build_call_graph.go` | S3: 에러 로깅 추가 |

### 신규/변경 테스트 파일

| 파일 | 내용 |
|---|---|
| `internal/chain/chain_test.go` | BuildCallGraph, TraverseChon, TraverseDepth 등 통합 테스트 추가 |
| `internal/parse/parse_go_file_test.go` | CallName, DetectControl, CalcMaxDepth, ExtractCalls 테스트 추가 |
| `internal/walk/walk_go_files_test.go` | MatchFFIgnore, MatchPattern 테스트 추가 |
| `internal/cli/cli_test.go` | FindGoMod, ResolveTarget 테스트 추가 |
| `internal/validate/mutest_test.go` | L1, L2 수정 검증 케이스 추가 |

### 신규 testdata 파일

| 파일 | 용도 |
|---|---|
| `internal/chain/testdata/*.go` | call graph 구성용 샘플 소스 |
| `internal/parse/testdata/call_name_*.go` | CallName 테스트용 샘플 소스 |
| `internal/parse/testdata/detect_control_*.go` | DetectControl 테스트용 샘플 소스 |

## 구현 순서

### Step 1: 버그 수정 (B1, B2)

1. `internal/cli/context.go` — ParseFFIgnore 경로 수정
2. `internal/cli/process_llmc_file.go` — WriteAnnotationLine 에러 처리

### Step 2: 로직 수정 (L1, L2)

3. `internal/validate/rule_a12.go` — control 빈값 조건 제거
4. `internal/validate/rule_a1.go` — hasTypes 독립 조건으로 변경
5. `internal/validate/mutest_test.go` — L1, L2 수정 검증 테스트 추가

### Step 3: 안전성 보강 (S1, S2, S3)

6. `internal/validate/evaluate_file.go` — 안전한 타입 단언
7. `internal/parse/read_module_path.go` — 빈 모듈 경로 에러
8. `internal/chain/build_call_graph.go` — 에러 로깅

### Step 4: 테스트 확충

9. chain testdata 작성 + chain_test.go 통합 테스트 추가
10. parse testdata 작성 + parse_go_file_test.go 단위 테스트 추가
11. walk_go_files_test.go 패턴 매칭 테스트 추가
12. cli_test.go 유틸리티 테스트 추가

### Step 5: 검증

13. `go build ./...`
14. `go vet ./...`
15. `go test ./...`
16. `filefunc validate`

## 완료 기준

- `go build ./...` 통과
- `go vet ./...` 통과
- `go test ./...` 전체 통과
- `filefunc validate` ERROR 0
- B1: context 명령에서 `.ffignore` 패턴 정상 적용
- B2: WriteAnnotationLine 실패 시 `[ERROR]` 출력 + `"fail"` 반환
- L1: `control=""` 일 때 A12 미발동 (A9만 발동)
- L2: func+type 공존 파일에서 `//ff:type` 누락 시 A1 발동
- S1: evaluate_file에서 잘못된 evidence 타입에 패닉하지 않음
- S2: go.mod에 module 행 없으면 에러 반환
- S3: ExtractCalls 실패 시 stderr에 경고 출력
- 테스트: BuildCallGraph, TraverseChon, CallName, DetectControl, CalcMaxDepth, MatchFFIgnore 테스트 존재

## 예상 규모

- 변경 파일: ~12개
- 신규 파일: ~8개 (testdata + 테스트 확장)
- 예상 난이도: 중 (버그/로직 수정은 기계적, 테스트 작성이 주요 작업)
- 핵심 난점: chain 통합 테스트를 위한 testdata 설계
