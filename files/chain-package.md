# `--package` 옵션 제안 — chain func 패키지 한정 필터

## 문제

`filefunc chain func ParseFile`을 실행하면 동명 함수가 여러 패키지에 존재할 때 어느 것이 잡힐지 예측할 수 없다.

fullend 프로젝트 실례:
```
filefunc chain func ParseFile --chon 1 --meta what

ParseFile (what="단일 HTML 파일을 파싱하여 PageSpec 반환")   ← stml
```

의도는 `internal/funcspec.ParseFile`이었으나 `internal/stml.ParseFile`이 반환됨. CallGraph의 키가 함수명 단독(`string`)이므로 패키지 구분이 불가능하다.

### 영향받는 함수 (fullend 기준)

| 함수명 | 존재 패키지 |
|---|---|
| ParseFile | funcspec, stml, ssac |
| ParseDir | funcspec, stml, ssac, statemachine, policy, scenario |
| ParseReader | stml |

## 제안: `--package` 플래그

```bash
filefunc chain func ParseFile --package funcspec --chon 1 --meta what
filefunc chain func ParseDir --package stml --chon 2 --meta what
```

### 동작

1. `--package` 미지정: 현재 동작 유지 (첫 번째 매칭)
2. `--package` 지정: 해당 Go 패키지에 속하는 함수만 대상으로 chain 수행
3. 매칭 없으면 에러: `Error: func "ParseFile" not found in package "funcspec"`

### 변경 지점

**1. `internal/cli/chain_func.go`** — 플래그 등록 + 값 전달

```go
// init()에 추가
chainFuncCmd.Flags().String("package", "", "limit to funcs in this Go package")

// RunE에서 읽기
pkg, _ := cmd.Flags().GetString("package")
```

**2. `internal/chain/call_graph.go`** — CallGraph 키 체계는 변경하지 않음

CallGraph는 `map[string][]string`으로 함수명만 키로 사용한다. 키를 `pkg.FuncName`으로 변경하면 전체 그래프 빌드/탐색/출력에 영향이 크다. 대신 필터링 방식으로 해결한다.

**3. `internal/chain/filter_by_package.go`** (신규)

`BuildFuncFileMap`의 결과(`map[string]*GoFile`)를 활용하여, `--package`에 매칭되지 않는 함수를 chain 결과에서 제거한다.

```go
func FilterByPackage(results []ChonResult, pkg string, fileMap map[string]*model.GoFile) []ChonResult {
    var filtered []ChonResult
    for _, r := range results {
        if gf, ok := fileMap[r.Name]; ok && gf.Package == pkg {
            filtered = append(filtered, r)
        }
    }
    return filtered
}
```

**4. `internal/cli/chain_func.go` RunE** — 시작 노드 검증 + 결과 필터

```go
// 시작 노드 패키지 검증
if pkg != "" {
    if gf, ok := fileMap[target]; ok {
        if gf.Package != pkg {
            return fmt.Errorf("func %q not found in package %q (found in %q)", target, pkg, gf.Package)
        }
    }
}

// chain 결과 필터링 (--package 지정 시 시작 노드의 실제 호출 관계만 남김)
```

### 한계

- CallGraph 키가 함수명 단독이므로, 동명 함수가 서로를 호출하는 경우 구분 불가
- 이 제안은 **출력 필터링**이지 그래프 자체의 패키지 인식이 아님
- 근본 해결은 CallGraph 키를 `pkg.FuncName`으로 변경하는 것이나, 파급 범위가 큼

### chain feature에는 불필요

`filefunc chain feature funcspec`은 `//ff:func feature=funcspec` 어노테이션으로 필터하므로 이름 충돌이 발생하지 않는다. `--package`는 `chain func` 전용이다.
