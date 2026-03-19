# Phase 028: 코드리뷰 조치 4차 — 런타임 패닉, 인덱스 버그, AST 안전성, 에러 처리

## 목표

4차 코드리뷰에서 발견된 10건을 전부 조치한다: (1) 런타임 패닉 2건 수정 (2) 인덱스 매핑 버그 수정 (3) AST Body nil 체크 3건 추가 (4) 에러 처리 일관성 3건 수정 (5) Scanner 에러 체크 추가.

## 배경

### 현재 상태

- `go build`, `go vet`, `go test ./...`, `filefunc validate` 모두 통과
- Phase027까지 3차례 코드리뷰 조치 완료
- 6개 에이전트 병렬 리뷰: cli, validate, chain, parse, context+llm+annotate+walk+model+report, 테스트 커버리지

### 리뷰 결과 요약

| 분류 | 건수 | 심각도 |
|---|---|---|
| 런타임 패닉 (nil 인덱싱) | 1 | CRITICAL |
| 인덱스 매핑 버그 | 1 | CRITICAL |
| AST Body nil 체크 누락 | 3 | CRITICAL |
| 에러 처리 누락/무시 | 3 | MEDIUM |
| Scanner 에러 미확인 | 1 | MEDIUM |
| 경로 매칭 부정확 | 1 | MEDIUM |

## 이슈 목록

### 1. 런타임 패닉: nil scores 인덱싱 (CRITICAL)

| # | 파일 | 이슈 |
|---|---|---|
| P1 | `internal/context/format_result.go:17` | `scores[i]` 접근 시 scores=nil이면 패닉. `run_pipeline.go:40`에서 `FormatResult(w, filtered, nil)` 호출 |

### 2. 인덱스 매핑 버그: filter_by_rate (CRITICAL)

| # | 파일 | 이슈 |
|---|---|---|
| P2 | `internal/chain/filter_by_rate.go:22` | `keptScores[len(kept)]`로 새 인덱스 기준 저장하지만, `format_chain.go:23`에서 `scores[i]`는 순차 인덱스 0부터 조회. 필터링 후 첫 항목(chon<2)은 scores에 없으므로 인덱스가 1씩 밀림 |

### 3. AST Body nil 체크 누락 (CRITICAL)

| # | 파일 | 이슈 |
|---|---|---|
| A1 | `internal/parse/calc_body_hash.go:29` | `fd.Body == nil` 체크 없음. 같은 패키지 `CalcMaxDepth`, `DetectControl`, `ExtractCalls`는 모두 체크 |
| A2 | `internal/parse/extract_func_source.go:21` | 동일 패턴 누락 |
| A3 | `internal/parse/extract_uses.go:23` | 동일 패턴 누락. `CollectTypeRefs(fd, ...)` 호출 시 Body nil 전달 가능 |

### 4. 에러 처리 누락 (MEDIUM)

| # | 파일 | 이슈 |
|---|---|---|
| E1 | `internal/cli/llmc.go:51-56` | 파싱 실패 시 `continue`로 무시. `validate.go`는 `fmt.Fprintf(os.Stderr, "warning: ...")` 출력 |
| E2 | `internal/context/build_feature_prompt.go:20` | `b, _ := json.Marshal(req)` 에러 무시 |
| E3 | `internal/context/score_body.go:20-22` | 파일 읽기 실패 시 `removed++`로 건너뛰지만 경고 없음 |

### 5. Scanner 에러 미확인 (MEDIUM)

| # | 파일 | 이슈 |
|---|---|---|
| S1 | `internal/validate/rule_a6.go:24-40` | `bufio.Scanner` 루프 후 `scanner.Err()` 미확인. I/O 에러 시 `false, nil` 반환 → 잘못된 결과 |

### 6. 경로 매칭 부정확 (MEDIUM)

| # | 파일 | 이슈 |
|---|---|---|
| M1 | `internal/walk/match_pattern.go:19` | `strings.Contains(path+"/", pattern)` — "node" 패턴이 "my_nodes/" 경로도 매칭 |

## 설계

### P1: format_result.go nil scores 방어

```go
// Before
score := scores[i]

// After
var score float64
if scores != nil {
    score = scores[i]
}
```

scores=nil이면 score=0.0으로 처리. 점수 미산출 상태를 0으로 표시.

### P2: filter_by_rate.go 인덱스 수정

```go
// Before
keptScores[len(kept)] = s

// After
keptScores[len(kept)-1] = s
```

`kept = append(kept, results[i])` 직후이므로 `len(kept)-1`이 방금 추가된 항목의 인덱스. `format_chain.go`에서 `scores[i]`로 순차 조회하므로 0-based 인덱스와 일치.

### A1: calc_body_hash.go Body nil 체크

```go
// Before
if !ok || fd.Name.Name == "init" {
    continue
}

// After
if !ok || fd.Body == nil || fd.Name.Name == "init" {
    continue
}
```

### A2: extract_func_source.go Body nil 체크

```go
// Before
if !ok || fd.Name.Name == "init" {
    continue
}

// After
if !ok || fd.Body == nil || fd.Name.Name == "init" {
    continue
}
```

### A3: extract_uses.go Body nil 체크

```go
// Before
if !ok || fd.Name.Name == "init" {
    continue
}

// After
if !ok || fd.Body == nil || fd.Name.Name == "init" {
    continue
}
```

### E1: llmc.go 파싱 실패 경고 출력

```go
// Before
if err != nil {
    continue
}

// After
if err != nil {
    fmt.Fprintf(os.Stderr, "warning: skipping %s: %v\n", p, err)
    continue
}
```

`validate.go:60`과 동일 패턴.

### E2: build_feature_prompt.go json.Marshal 에러 처리

```go
// Before
b, _ := json.Marshal(req)
return string(b)

// After
b, err := json.Marshal(req)
if err != nil {
    return ""
}
return string(b)
```

### E3: score_body.go 경고 출력 추가

파일 읽기 실패 시 stderr에 경고 출력.

```go
// Before
removed++
continue

// After
fmt.Fprintf(os.Stderr, "warning: skipping %s: %v\n", gf.Path, err)
removed++
continue
```

import `"fmt"`, `"os"` 추가 필요 여부 확인.

### S1: rule_a6.go scanner.Err() 체크

```go
// Before
for scanner.Scan() {
    // ...
}
return false, nil

// After
for scanner.Scan() {
    // ...
}
if err := scanner.Err(); err != nil {
    return false, err
}
return false, nil
```

### M1: match_pattern.go 경로 매칭 개선

```go
// Before
return strings.Contains(path+"/", pattern)

// After
return strings.Contains("/"+path+"/", "/"+pattern+"/") || strings.Contains("/"+path+"/", "/"+pattern)
```

패턴 앞뒤에 "/" 경계를 추가하여 부분 매칭 방지. 단, 기존 동작과의 호환성을 확인해야 하므로 `match_pattern.go`의 전체 로직을 읽고 정확한 수정안을 결정한다.

## 영향 범위

### 변경 파일

| 파일 | 변경 내용 |
|---|---|
| `internal/context/format_result.go` | P1: scores nil 방어 |
| `internal/chain/filter_by_rate.go` | P2: 인덱스 매핑 수정 |
| `internal/parse/calc_body_hash.go` | A1: fd.Body nil 체크 추가 |
| `internal/parse/extract_func_source.go` | A2: fd.Body nil 체크 추가 |
| `internal/parse/extract_uses.go` | A3: fd.Body nil 체크 추가 |
| `internal/cli/llmc.go` | E1: 파싱 실패 경고 출력 |
| `internal/context/build_feature_prompt.go` | E2: json.Marshal 에러 처리 |
| `internal/context/score_body.go` | E3: 파일 읽기 실패 경고 출력 |
| `internal/validate/rule_a6.go` | S1: scanner.Err() 체크 |
| `internal/walk/match_pattern.go` | M1: 경로 매칭 정확도 개선 |

### 삭제/신규 파일

없음.

## 구현 순서

### Step 1: 런타임 패닉 수정 (P1, P2)

1. `internal/context/format_result.go` — scores nil 방어
2. `internal/chain/filter_by_rate.go` — 인덱스 매핑 수정

### Step 2: AST 안전성 (A1, A2, A3)

3. `internal/parse/calc_body_hash.go` — fd.Body nil 체크
4. `internal/parse/extract_func_source.go` — fd.Body nil 체크
5. `internal/parse/extract_uses.go` — fd.Body nil 체크

### Step 3: 에러 처리 (E1, E2, E3)

6. `internal/cli/llmc.go` — 파싱 실패 경고 출력
7. `internal/context/build_feature_prompt.go` — json.Marshal 에러 처리
8. `internal/context/score_body.go` — 파일 읽기 실패 경고 출력

### Step 4: Scanner 에러 (S1)

9. `internal/validate/rule_a6.go` — scanner.Err() 체크

### Step 5: 경로 매칭 (M1)

10. `internal/walk/match_pattern.go` — 매칭 로직 개선 (기존 테스트 확인 후 수정)

### Step 6: 검증

11. `go build ./...`
12. `go vet ./...`
13. `go test ./...`
14. `filefunc validate`

## 완료 기준

- `go build ./...` 통과
- `go vet ./...` 통과
- `go test ./...` 전체 통과
- `filefunc validate` ERROR 0
- P1: `FormatResult`에 scores=nil 전달 시 패닉 없음
- P2: `FilterByRate` 후 `FormatChain`에서 점수 정상 표시
- A1/A2/A3: 인터페이스 메서드(Body=nil) 포함 파일 파싱 시 패닉 없음
- E1: `llmc` 실행 시 파싱 실패 파일이 stderr에 경고 출력
- E2: `BuildFeaturePrompt`에서 json.Marshal 실패 시 빈 문자열 반환 (패닉 없음)
- E3: `ScoreBody`에서 파일 읽기 실패 시 stderr 경고 출력
- S1: `RuleA6`에서 I/O 에러 시 에러 반환 (잘못된 false 반환 방지)
- M1: "node" 패턴이 "my_nodes" 경로를 매칭하지 않음

## 예상 규모

- 변경 파일: 10개
- 삭제/신규 파일: 0개
- 예상 난이도: 하 (대부분 기계적 수정, M1만 기존 테스트 확인 필요)
- 핵심 난점: P2 인덱스 수정 시 `format_chain.go`와의 인덱스 계약 확인, M1 기존 .ffignore 패턴 호환성
