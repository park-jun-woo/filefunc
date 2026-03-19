# Phase 027: 코드리뷰 조치 3차 — JSON injection, 에러 처리, 문서 정합, 안전성 ✅

## 목표

3차 코드리뷰에서 발견된 8건을 전부 조치한다: (1) JSON injection 제거 (2) 에러 무시 수정 (3) 문서/주석 불일치 수정 (4) filepath.Join 통일 (5) 입력 검증 추가.

## 배경

### 현재 상태

- `go build`, `go vet`, `go test ./...`, `filefunc validate` 모두 통과
- Phase026에서 버그, 로직 결함, 안전성, 테스트 확충 조치 완료
- 4개 에이전트 병렬 리뷰: validate, chain, parse, cli+context+llm+annotate+report+walk+model

### 리뷰 결과 요약

| 분류 | 건수 | 심각도 |
|---|---|---|
| 보안 (JSON injection) | 1 | HIGH |
| 에러 무시 | 2 | HIGH / LOW |
| defeat 그래프 누락 | 1 | MEDIUM |
| 문서/주석 불일치 | 2 | MEDIUM / LOW |
| 경로 처리 미흡 | 1 (5파일) | LOW |
| 입력 검증 누락 | 1 | LOW |

## 이슈 목록

### 1. JSON injection (HIGH)

| # | 파일 | 이슈 |
|---|---|---|
| S1 | `internal/llm/pull_model.go:16` | `fmt.Sprintf`로 JSON 문자열 직접 조립. model 이름에 `"`가 포함되면 JSON 구조 파괴. 같은 패키지 `ollama_generate.go`는 `json.Marshal` 사용 |

### 2. 에러 무시 (HIGH / LOW)

| # | 파일 | 이슈 |
|---|---|---|
| E1 | `internal/context/score_body.go:19` | `os.ReadFile` 에러 무시 (`src, _ :=`). 실패 시 nil이 `ExtractFuncSource`에 전달됨 |
| E2 | `internal/llm/pull_model.go:24` | `io.ReadAll` 에러 무시. 에러 응답 본문 읽기 실패 시 빈 에러 메시지 반환 |

### 3. defeat 그래프 누락 (MEDIUM)

| # | 파일 | 이슈 |
|---|---|---|
| D1 | `internal/validate/validate_graph.go:43-44` | 테스트 파일에 대해 A1~A3, A6, A8~A16 모두 defeat하는데 A7만 누락. 테스트 파일의 `//ff:checked` 해시 불일치가 ERROR로 잡힐 수 있음 |

### 4. 문서/주석 불일치 (MEDIUM / LOW)

| # | 파일 | 이슈 |
|---|---|---|
| C1 | `internal/chain/chon_result.go:9` | Rel 필드 주석이 완전히 틀림. 주석: `"child", "parent", "sibling", ...` → 실제: `"calls", "called-by", "co-called", ...` |
| C2 | `internal/chain/traverse_depth.go:6` | direction 파라미터 주석 `"child" or "parent"` → 실제: `"calls" or "called-by"` |

### 5. 경로 처리 미흡 (LOW)

| # | 파일 | 이슈 |
|---|---|---|
| P1 | `internal/cli/check_project_root.go:19,22` | `root + "/go.mod"` — 문자열 직접 연결. `filepath.Join` 미사용 |
| P1 | `internal/cli/validate.go:35,49` | 동일 패턴 |
| P1 | `internal/cli/llmc.go:43` | 동일 패턴 |
| P1 | `internal/cli/context.go:38,43` | 동일 패턴 |
| P1 | `internal/cli/build_graph.go:16,21` | 동일 패턴 |

### 6. 입력 검증 누락 (LOW)

| # | 파일 | 이슈 |
|---|---|---|
| V1 | `internal/chain/traverse_depth_recur.go:11-15` | `direction == "calls"`만 확인하고 나머지 전부 parent로 처리. 잘못된 direction 값이 무시됨 |

## 설계

### S1: pull_model.go JSON injection 제거

```go
// Before
reqBody := fmt.Sprintf(`{"name":"%s","stream":false}`, model)
resp, err := http.Post(endpoint+"/api/pull", "application/json", strings.NewReader(reqBody))

// After
reqBody, err := json.Marshal(map[string]interface{}{"name": model, "stream": false})
if err != nil {
    return fmt.Errorf("marshal request: %w", err)
}
resp, err := http.Post(endpoint+"/api/pull", "application/json", bytes.NewReader(reqBody))
```

import: `strings` → `bytes`, `encoding/json` 추가.

### E1: score_body.go 에러 처리

```go
// Before
src, _ := os.ReadFile(gf.Path)

// After
src, err := os.ReadFile(gf.Path)
if err != nil {
    removed++
    continue
}
```

### E2: pull_model.go 응답 본문 에러 처리

```go
// Before
body, _ := io.ReadAll(resp.Body)
return fmt.Errorf("pull failed (%d): %s", resp.StatusCode, string(body))

// After
body, err := io.ReadAll(resp.Body)
if err != nil {
    return fmt.Errorf("pull failed (%d): read body: %w", resp.StatusCode, err)
}
return fmt.Errorf("pull failed (%d): %s", resp.StatusCode, string(body))
```

### D1: validate_graph.go A7 defeater 추가

```go
// Before
Defeat(DefeaterTestFile, RuleA6).
Defeat(DefeaterTestFile, RuleA8).

// After
Defeat(DefeaterTestFile, RuleA6).
Defeat(DefeaterTestFile, RuleA7).
Defeat(DefeaterTestFile, RuleA8).
```

### C1: chon_result.go Rel 주석 수정

```go
// Before
Rel  string // "child", "parent", "sibling", "grandchild", "grandparent", "uncle", "nephew"

// After
Rel  string // "calls", "called-by", "calls-2depth", "called-by-2depth", "co-called", "caller-peer", "peer-calls"
```

### C2: traverse_depth.go direction 주석 수정

```go
// Before
// direction is "child" or "parent".

// After
// direction is "calls" or "called-by".
```

### P1: cli 패키지 filepath.Join 통일

5개 파일의 `root + "/파일명"` 패턴을 전부 `filepath.Join(root, "파일명")`으로 변경. 각 파일에 `"path/filepath"` import 추가.

대상 파일:
- `check_project_root.go`: `root + "/go.mod"`, `root + "/codebook.yaml"`
- `validate.go`: `root + "/codebook.yaml"`, `root + "/.ffignore"`
- `llmc.go`: `root + "/.ffignore"`
- `context.go`: `root + "/codebook.yaml"`, `root + "/.ffignore"`
- `build_graph.go`: `root + "/go.mod"`, `root + "/.ffignore"`

### V1: traverse_depth_recur.go direction 검증

```go
// Before
var nexts []string
if direction == "calls" {
    nexts = g.Children[current]
} else {
    nexts = g.Parents[current]
}

// After
if direction != "calls" && direction != "called-by" {
    return
}

var nexts []string
if direction == "calls" {
    nexts = g.Children[current]
} else {
    nexts = g.Parents[current]
}
```

early return으로 처리. switch 문 사용 시 `control=iteration` 파일에서 A14 (switch 금지) 위반이므로 if 문 사용.

## 영향 범위

### 변경 파일

| 파일 | 변경 내용 |
|---|---|
| `internal/llm/pull_model.go` | S1: json.Marshal 사용, E2: io.ReadAll 에러 처리 |
| `internal/context/score_body.go` | E1: os.ReadFile 에러 처리 |
| `internal/validate/validate_graph.go` | D1: A7 defeater 추가 |
| `internal/chain/chon_result.go` | C1: Rel 주석 수정 |
| `internal/chain/traverse_depth.go` | C2: direction 주석 수정 |
| `internal/chain/traverse_depth_recur.go` | V1: direction 검증 추가 |
| `internal/cli/check_project_root.go` | P1: filepath.Join |
| `internal/cli/validate.go` | P1: filepath.Join |
| `internal/cli/llmc.go` | P1: filepath.Join |
| `internal/cli/context.go` | P1: filepath.Join |
| `internal/cli/build_graph.go` | P1: filepath.Join |

### 삭제/신규 파일

없음.

## 구현 순서

### Step 1: 보안 수정 (S1)

1. `internal/llm/pull_model.go` — `fmt.Sprintf` → `json.Marshal`, import 변경

### Step 2: 에러 처리 수정 (E1, E2)

2. `internal/context/score_body.go` — `os.ReadFile` 에러 처리
3. `internal/llm/pull_model.go` — `io.ReadAll` 에러 처리 (Step 1과 함께)

### Step 3: defeat 그래프 수정 (D1)

4. `internal/validate/validate_graph.go` — A7 defeater 추가

### Step 4: 문서 수정 (C1, C2)

5. `internal/chain/chon_result.go` — Rel 주석 수정
6. `internal/chain/traverse_depth.go` — direction 주석 수정

### Step 5: 경로 처리 통일 (P1)

7. `internal/cli/check_project_root.go` — filepath.Join
8. `internal/cli/validate.go` — filepath.Join
9. `internal/cli/llmc.go` — filepath.Join
10. `internal/cli/context.go` — filepath.Join
11. `internal/cli/build_graph.go` — filepath.Join

### Step 6: 입력 검증 (V1)

12. `internal/chain/traverse_depth_recur.go` — direction 검증

### Step 7: 검증

13. `go build ./...`
14. `go vet ./...`
15. `go test ./...`
16. `filefunc validate`

## 완료 기준

- `go build ./...` 통과
- `go vet ./...` 통과
- `go test ./...` 전체 통과
- `filefunc validate` ERROR 0
- S1: `PullModel`에서 model 이름에 특수문자 포함 시 JSON 구조 유지
- E1: `ScoreBody`에서 파일 읽기 실패 시 `removed++`로 건너뜀
- E2: `PullModel`에서 응답 본문 읽기 실패 시 에러 메시지 포함
- D1: 테스트 파일의 `//ff:checked` 해시 불일치가 defeat됨
- C1/C2: 주석이 실제 코드 동작과 일치
- P1: cli 패키지 내 `root + "/"` 패턴 0건
- V1: 잘못된 direction 값 입력 시 silent wrong behavior 대신 early return

## 예상 규모

- 변경 파일: 11개
- 삭제/신규 파일: 0개
- 예상 난이도: 하 (전부 기계적 수정)
- 핵심 난점: `traverse_depth_recur.go`에서 switch 사용 시 A14 위반 — if 문으로 우회
