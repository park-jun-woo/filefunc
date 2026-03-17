# Phase 022: CallGraph 패키지 인식 키 전환

## 목표

CallGraph 키를 `FuncName` → `pkg.FuncName`으로 변경하여 동명이패키지 함수 충돌(BUG001)을 근본 해결한다.

## 배경

### 현재 문제

`CallGraph`와 `BuildFuncFileMap`이 함수명만 키로 사용하여 동명 함수가 덮어써진다.

```bash
filefunc chain func ParseFile --chon 1 --meta what
# → stml/parser.ParseFile이 출력됨 (의도: funcspec.ParseFile)

filefunc chain func ParseFile --package funcspec --chon 1 --meta what
# → Error: func "ParseFile" not found in package "funcspec" (found in "parser")
```

fullend 기준 동명 함수:

| 함수명 | 존재 패키지 |
|---|---|
| ParseFile | funcspec, stml/parser, ssac/parser, policy, statemachine |
| ParseDir | funcspec, stml/parser, ssac/parser, statemachine, policy, scenario |

Phase021에서 추가한 `--package` 필터는 출력 필터링이므로, 그래프 자체가 덮어써진 상태에서는 동작하지 않는다.

### 근본 원인

```go
// call_graph.go — 키가 함수명 단독
type CallGraph struct {
    Children map[string][]string  // "ParseFile" → callees
    Parents  map[string][]string  // "ParseFile" → callers
}

// func_file_map.go — 동명 함수 덮어쓰기
m[name] = gf  // 마지막 등록이 승리
```

### 해결

키를 `pkg.FuncName` 형식으로 변경:

```go
Children["funcspec.ParseFile"] = ["funcspec.parseCommentGroup", "funcspec.processDecl"]
Children["parser.ParseFile"]   = ["parser.ParseReader"]
```

## 설계

### 키 형식

`qualifiedName(pkg, name) string` → `"pkg.FuncName"`

- `funcspec.ParseFile`, `parser.ParseFile` — 충돌 없음
- 같은 패키지 내 함수는 같은 prefix → `--package` 필터가 정확히 동작

### CallName 변경 (핵심)

현재 `CallName`은 callee 이름만 반환:

```go
// Before: 외부 패키지 호출 → sel.Sel.Name (함수명만)
if ok && projImports[ident.Name] {
    return sel.Sel.Name
}
// Before: 같은 패키지 호출 → ident.Name (함수명만)
if ok && projFuncs[ident.Name] {
    return ident.Name
}
```

변경: callee에 패키지명 포함

```go
// After: 외부 패키지 호출 → "alias.FuncName" (alias는 import에서 추출)
// After: 같은 패키지 호출 → "callerPkg.FuncName"
```

**핵심 난점**: `CallName`이 현재 `projImports map[string]bool`을 받는데, 이를 `map[string]string` (alias → 실제 패키지명)으로 변경해야 한다. `BuildImportMap`이 import 경로의 마지막 세그먼트를 alias로 사용하므로, 이 정보를 보존하면 된다.

### BuildImportMap 변경

```go
// Before: map[string]bool (alias → 프로젝트 내부 여부)
func BuildImportMap(f *ast.File, modulePath string) map[string]bool

// After: map[string]string (alias → 패키지명)
func BuildImportMap(f *ast.File, modulePath string) map[string]string
```

import 경로 `github.com/user/proj/internal/funcspec`에서:
- alias: `funcspec` (또는 `import foo "..."` → `foo`)
- 패키지명: 경로의 마지막 세그먼트 = `funcspec`

alias import 시: `import fp "internal/funcspec"` → alias=`fp`, 패키지명=`funcspec`

### CollectProjectSymbols 변경

```go
// Before: funcs map[string]bool
// After: funcs map[string]string (funcName → pkg)
```

`ExtractCalls`가 같은 패키지 호출 시 caller의 패키지를 알아야 하므로, `projFuncs`에 패키지 정보를 포함시킨다.

## 영향 범위

### 변경 파일

| 파일 | 변경 내용 |
|---|---|
| `internal/parse/build_import_map.go` | 반환 `map[string]bool` → `map[string]string` |
| `internal/parse/call_name.go` | callee에 패키지명 포함 반환, `projImports` 타입 변경 |
| `internal/parse/extract_calls.go` | caller 패키지 전달, `projFuncs` 타입 변경 |
| `internal/parse/collect_project_symbols.go` | `map[string]bool` → `map[string]string` (name → pkg) |
| `internal/chain/build_call_graph.go` | caller 키를 qualified name으로, ExtractCalls 인자 변경 |
| `internal/chain/func_file_map.go` | 키를 qualified name으로 |
| `internal/chain/func_name.go` | `qualifiedName` 함수 추가 |
| `internal/chain/format_chain.go` | 출력 시 패키지명 표시/생략 정책 |
| `internal/cli/chain_func.go` | target을 qualified name으로 변환, 동명 함수 안내 |
| `internal/cli/chain_feature.go` | FilterByFeature 결과가 qualified name |
| `internal/chain/filter_by_feature.go` | qualified name 반환 |
| `internal/chain/filter_by_package.go` | qualified key에서 패키지 추출 |
| `internal/chain/score_relevance.go` | fileMap 조회 키 변경 |
| `internal/chain/build_score_input.go` | fileMap 조회 키를 qualified name으로 대응 |
| `internal/parse/extract_uses.go` | `BuildImportMap` 반환 타입 변경 대응 (`map[string]bool` → `map[string]string`) |
| `internal/parse/collect_type_refs.go` | `projImports` 타입 `map[string]bool` → `map[string]string` 대응 (`_, ok :=` 패턴으로 변경) |
| `internal/cli/build_graph.go` | CollectProjectSymbols 반환 타입 대응 |

### 영향 없는 파일

- `traverse_chon.go`, `traverse_depth.go`, `traverse_depth_recur.go` — string 키 그대로, 형식만 변경
- `collect_chon.go`, `find_siblings.go`, `expand_through.go` — string 연산, 변경 없음

## 출력 정책

### `--package` 미지정 + 동명 함수

```
Multiple funcs named "ParseFile" found:
  funcspec.ParseFile
  parser.ParseFile (stml)
  parser.ParseFile (ssac)
Use --package to specify.
```

### `--package` 지정 시

```
ParseFile (what="단일 func spec .go 파일을 파싱하여 FuncSpec을 반환한다")
  1촌 calls: parseCommentGroup (what="...")
  1촌 called-by: ParseDir (what="...")
```

출력에서 같은 패키지 함수는 패키지명 생략, 다른 패키지 함수는 `pkg.FuncName` 형식.

### `--package` 미지정 + 동명 없음 (기존과 동일)

```
RunAll (what="모든 검증 룰을 실행하고 위반 목록을 반환")
  1촌 calls: CheckOneFileOneFunc (what="...")
```

## 구현 순서

1. `parse/build_import_map.go` — `map[string]string` 반환으로 변경
2. `parse/collect_project_symbols.go` — `map[string]string` (name → pkg) 반환
3. `parse/call_name.go` — callee에 패키지명 포함, 인자 타입 변경
4. `parse/extract_calls.go` — caller 패키지 전달, qualified callee 반환
5. `parse/collect_type_refs.go` — `projImports` 타입 대응 (`_, ok :=` 패턴)
6. `parse/extract_uses.go` — `BuildImportMap` 반환 타입 대응
7. `chain/func_name.go` — `qualifiedName` 추가
8. `chain/build_call_graph.go` — qualified key로 그래프 구성
9. `chain/func_file_map.go` — qualified key로 매핑
10. `chain/filter_by_feature.go` — qualified name 반환
11. `chain/filter_by_package.go` — qualified key 대응
12. `chain/format_chain.go` — 출력 정책 구현
13. `chain/score_relevance.go` — fileMap 조회 키 대응
14. `chain/build_score_input.go` — fileMap 조회 키 대응
15. `cli/chain_func.go` — 동명 함수 안내 + target 변환
16. `cli/chain_feature.go` — qualified name 대응
17. `cli/build_graph.go` — CollectProjectSymbols 반환 타입 대응
18. 빌드 + 테스트 + validate

## 완료 기준

- `go build ./...` 통과
- `go test ./...` 통과
- `filefunc validate` ERROR 0
- filefunc 자체에서 `chain func` 기존과 동일 동작
- fullend에서 재현 테스트:
  ```bash
  filefunc chain func ParseFile --chon 1 --meta what
  # → 동명 함수 후보 목록 출력
  filefunc chain func ParseFile --package funcspec --chon 1 --meta what
  # → funcspec.ParseFile의 정확한 호출 관계 출력
  ```

## 하위 호환성

- `chain func FuncName` (패키지 미지정) — 동명 함수 없으면 기존과 동일
- `chain func FuncName` (동명 함수 존재) — 후보 목록 안내 (기존: 임의 선택)
- `chain feature FeatureName` — 기존과 동일 (feature 필터는 이름 충돌 없음)
- `--package` 플래그 — 기존 필터링에서 그래프 수준 정확 매칭으로 승격

## 예상 규모

- 변경 파일: 17개
- 신규 파일: 0개
- 예상 난이도: 중상
- 핵심 난점: `CallName`에서 alias import → 실제 패키지명 매핑
