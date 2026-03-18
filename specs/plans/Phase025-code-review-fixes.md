# Phase 025: 코드리뷰 조치 — 에러 처리, 중복 제거, 로직 수정, 테스트 보강

## 목표

코드리뷰에서 발견된 5개 카테고리 이슈를 전부 조치한다: (1) 무시된 에러 처리 (2) 코드 중복 제거 (3) 로직 버그 수정 (4) 미사용 코드 삭제 (5) 테스트 커버리지 보강.

## 배경

### 현재 상태

- `go build`, `go vet`, `filefunc validate` 모두 통과
- 테스트 있는 패키지: parse, validate, walk (3개)
- 테스트 없는 패키지: annotate, chain, cli, context, llm, model, report (8개)
- 에러 무시 5건, 코드 중복 2쌍, 로직 버그 2건, 미사용 코드 1건

## 이슈 목록

### 1. 에러 처리 누락 (HIGH)

| # | 파일 | 위치 | 이슈 |
|---|---|---|---|
| E1 | `internal/cli/context.go` | :75 | `json.Marshal()` 에러 무시 |
| E2 | `internal/cli/context.go` | :84 | `io.ReadAll()` 에러 무시 |
| E3 | `internal/cli/context.go` | :89 | `json.Unmarshal()` 에러 완전 무시 |
| E4 | `internal/llm/check_model.go` | :22 | `reader.ReadString()` 에러 무시 |

### 2. 코드 중복 (HIGH)

| # | 파일 A | 파일 B | 함수 |
|---|---|---|---|
| D1 | `cli/name_from_qualified.go` | `chain/name_from_qualified.go` | `nameFromQualified` |
| D2 | `cli/pkg_from_qualified.go` | `chain/pkg_from_qualified.go` | `pkgFromQualified` |

두 패키지 모두 unexported(소문자)로 동일 로직. chain 패키지를 정본으로 하고 exported 함수로 전환. cli는 chain 패키지를 호출.

### 3. 로직 버그 (MEDIUM)

| # | 파일 | 이슈 |
|---|---|---|
| L1 | `internal/parse/if_else_depth.go:15` | `else` 블록의 depth 과소 계산 |
| L2 | `internal/chain/find_callers.go` | exported이나 프로젝트 내 미사용 |

#### L1 상세: else 블록 depth 버그

```go
// 현재 (버그)
func IfElseDepth(s *ast.IfStmt, current int) int {
    d := StmtDepth(s.Body.List, current+1)  // if body: current+1 ✓
    if s.Else == nil { return d }
    ed := NodeDepth(s.Else, current)         // else body: current (BlockStmt → StmtDepth(list, current))
    ...
}
```

- `if` body: `current+1`에서 계산 (정상)
- `else` body: `NodeDepth(s.Else, current)` → `BlockStmt` case → `StmtDepth(s.List, current)` → `current`에서 계산 (1 부족)
- `else if`: `NodeDepth(s.Else, current)` → `IfStmt` case → `IfElseDepth(s, current)` → body가 `current+1` (정상)

**수정**: else 분기를 직접 처리하여 `BlockStmt`일 때 `current+1`로 계산:

```go
func IfElseDepth(s *ast.IfStmt, current int) int {
    d := StmtDepth(s.Body.List, current+1)
    if s.Else == nil { return d }
    var ed int
    switch e := s.Else.(type) {
    case *ast.IfStmt:
        ed = IfElseDepth(e, current)
    case *ast.BlockStmt:
        ed = StmtDepth(e.List, current+1)
    default:
        ed = NodeDepth(s.Else, current)
    }
    if ed > d { return ed }
    return d
}
```

### 4. 테스트 보강

| 패키지 | 우선순위 | 테스트 대상 | 비고 |
|---|---|---|---|
| chain | HIGH | `nameFromQualified`, `pkgFromQualified`, `qualifiedName`, `funcName`, `addUnique`, `filterByFeature`, `hasFeature`, `collectChon`, `countChon2Plus` | 순수 함수, 단위 테스트 용이 |
| annotate | HIGH | `WriteAnnotationLine`, `ReplaceAnnotationLine`, `InsertAfterAnnotations` | 순수 함수 |
| report | MEDIUM | `FormatText`, `FormatJSON` | 출력 포매터 |
| context | MEDIUM | `ParseFeatures`, `ParseScores`, `ParseSingleScore`, `ParseSearch`, `FilterFeature`, `MatchAnnotation` | 파싱/필터 함수 |
| llm | LOW | `BuildPrompt`, `ParseScore` | 외부 의존 없는 함수만 |
| model | LOW | 구조체만 — 테스트 불필요 | — |
| cli | LOW | `FindGoMod`, `CheckProjectRoot`, `EnvOrDefault` | 유틸리티만 |

## 설계

### D1/D2 중복 제거 전략

1. `chain/name_from_qualified.go` → `NameFromQualified` (exported)
2. `chain/pkg_from_qualified.go` → `PkgFromQualified` (exported)
3. `cli/name_from_qualified.go` 삭제 → `chain.NameFromQualified()` 호출로 대체
4. `cli/pkg_from_qualified.go` 삭제 → `chain.PkgFromQualified()` 호출로 대체
5. cli 패키지 내 호출부 전부 변경

### L2 미사용 코드 처리

`FindCallers`는 `g.Parents[name]` 한 줄 래퍼. 프로젝트 내 미사용 — 삭제.

## 영향 범위

### 변경 파일

| 파일 | 변경 내용 |
|---|---|
| `internal/cli/context.go` | E1~E3: json.Marshal, io.ReadAll, json.Unmarshal 에러 처리 추가 |
| `internal/llm/check_model.go` | E4: reader.ReadString 에러 처리 추가 |
| `internal/chain/name_from_qualified.go` | D1: `nameFromQualified` → `NameFromQualified` (exported) |
| `internal/chain/pkg_from_qualified.go` | D2: `pkgFromQualified` → `PkgFromQualified` (exported) |
| `internal/chain/*.go` (호출부) | `nameFromQualified` → `NameFromQualified` 등 호출부 변경 |
| `internal/cli/resolve_qualified_target.go` 등 | cli 내 호출부 → `chain.NameFromQualified()` |
| `internal/parse/if_else_depth.go` | L1: else 블록 depth 수정 |

### 삭제 파일

| 파일 | 이유 |
|---|---|
| `internal/cli/name_from_qualified.go` | D1: chain 패키지로 통합 |
| `internal/cli/pkg_from_qualified.go` | D2: chain 패키지로 통합 |
| `internal/chain/find_callers.go` | L2: 미사용 코드 |

### 신규 파일

| 파일 | 내용 |
|---|---|
| `internal/chain/chain_test.go` | chain 패키지 단위 테스트 |
| `internal/annotate/annotate_test.go` | annotate 패키지 단위 테스트 |
| `internal/report/report_test.go` | report 패키지 단위 테스트 |
| `internal/context/context_test.go` | context 파싱/필터 테스트 |
| `internal/llm/llm_test.go` | BuildPrompt, ParseScore 테스트 |
| `internal/cli/cli_test.go` | 유틸리티 함수 테스트 |
| `internal/parse/if_else_depth_test.go` | L1 수정 검증 테스트 |

## 구현 순서

### Step 1: 에러 처리 수정 (E1~E4)

1. `internal/cli/context.go` — `ollamaGenerate` 함수 에러 처리 추가
2. `internal/llm/check_model.go` — `reader.ReadString` 에러 처리 추가

### Step 2: 로직 버그 수정 (L1)

3. `internal/parse/if_else_depth.go` — else 블록 depth 수정
4. `internal/parse/if_else_depth_test.go` — 수정 검증 테스트 작성

### Step 3: 코드 중복 제거 (D1, D2)

5. `chain/name_from_qualified.go` — exported로 변경
6. `chain/pkg_from_qualified.go` — exported로 변경
7. chain 패키지 내 호출부 전부 변경 (unexported → exported)
8. cli 패키지 호출부 변경 → `chain.NameFromQualified()` 사용
9. `cli/name_from_qualified.go`, `cli/pkg_from_qualified.go` 삭제

### Step 4: 미사용 코드 삭제 (L2)

10. `chain/find_callers.go` 삭제

### Step 5: 테스트 작성

11. `internal/chain/chain_test.go`
12. `internal/annotate/annotate_test.go`
13. `internal/report/report_test.go`
14. `internal/context/context_test.go`
15. `internal/llm/llm_test.go`
16. `internal/cli/cli_test.go`

### Step 6: 검증

17. `go build ./...`
18. `go vet ./...`
19. `go test ./...`
20. `filefunc validate`

## 완료 기준

- `go build ./...` 통과
- `go vet ./...` 통과
- `go test ./...` 전체 통과 (기존 + 신규)
- `filefunc validate` ERROR 0
- 에러 무시 0건 (`_, _` 패턴 제거)
- 중복 함수 0쌍 (cli 측 삭제, chain 측 통합)
- `FindCallers` 삭제
- `IfElseDepth` else 블록 depth 정상
- 테스트 있는 패키지: parse, validate, walk + chain, annotate, report, context, llm, cli (9개)

## 예상 규모

- 변경 파일: ~12개
- 삭제 파일: 3개
- 신규 파일: 7개 (테스트)
- 예상 난이도: 중하 (대부분 기계적 수정 + 테스트 작성)
- 핵심 난점: chain 패키지 exported 전환 시 호출부 전수 조사
