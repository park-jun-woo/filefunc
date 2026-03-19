# Phase 011: chain — func/feature 데이터 흐름 추적 ✅ 완료

## 목표
`filefunc chain`으로 func의 호출 관계와 데이터 흐름을 추적한다. 실시간 AST 분석으로 그래프를 구성한다. 코드북의 feature가 줌 레벨.

## 명령

```bash
filefunc chain func CheckOneFileOneFunc    # 이 함수의 호출 체인 추적
filefunc chain feature validate            # validate feature 전체 체인
```

## 설계 원칙
- 어노테이션 캐시(calls/uses) 없이 실시간 AST 분석
- ExtractCalls, ExtractUses 기존 코드 재사용
- callgraph와의 차이: 프로젝트 내부만, feature 단위로 스코프 제한

## chain func 출력

대상 func에서 시작하여 호출하는 func, 그 func이 호출하는 func... 재귀적으로 추적.

```
CheckOneFileOneFunc
  ├─ calls: IsConstOnly
  └─ called by: RunAll

IsConstOnly
  └─ (leaf — no project-internal calls)
```

## chain feature 출력

해당 feature의 모든 func를 수집하고, 호출 관계를 그래프로 표시.

```
feature=validate

RunAll
  ├─ CheckOneFileOneFunc → IsConstOnly
  ├─ CheckOneFileOneType
  ├─ CheckOneFileOneMethod
  ├─ CheckInitStandalone
  ├─ CheckNestingDepth
  ├─ CheckFuncLines
  ├─ CheckAnnotationRequired
  ├─ CheckCodebookValues → AllowedValues, Contains
  ├─ CheckWhatRequired
  ├─ CheckAnnotationPosition
  ├─ CheckRequiredKeysInAnnotation
  ├─ CheckCheckedHash → CalcBodyHash
  └─ HasAnyChecked
```

## 산출물

| 파일 | 내용 |
|---|---|
| `internal/chain/build_call_graph.go` | BuildCallGraph — 프로젝트 전체 호출 그래프 구성 (caller→callee, callee→caller 양방향) |
| `internal/chain/find_callers.go` | FindCallers — 특정 func을 호출하는 func 목록 (역방향) |
| `internal/chain/find_siblings.go` | FindSiblings — 같은 부모에게 호출되는 func 목록 (2촌 형제) |
| `internal/chain/traverse_chon.go` | TraverseChon — 촌수 기반 그래프 탐색 |
| `internal/chain/traverse_depth.go` | TraverseDepth — 단방향 깊이 탐색 (child-depth, parent-depth) |
| `internal/chain/filter_by_feature.go` | FilterByFeature — feature별 func 필터링 |
| `internal/chain/format_chain.go` | FormatChain — 체인 결과를 텍스트로 출력 |
| `internal/cli/chain.go` | chainCmd — cobra 부모 서브커맨드 |
| `internal/cli/chain_func.go` | chainFuncCmd — chain func 서브커맨드 |
| `internal/cli/chain_feature.go` | chainFeatureCmd — chain feature 서브커맨드 |

## 실행 흐름

### chain func

```
filefunc chain func CheckOneFileOneFunc
  │
  ├─ ReadModulePath → 모듈 경로
  ├─ WalkGoFiles → ParseGoFile 전체
  ├─ CollectProjectSymbols → 심볼 맵
  │
  ├─ BuildCallGraph(파일들, 모듈경로, 심볼) → map[func][]callee
  │
  ├─ 대상 func에서 시작:
  │    ├─ calls: ExtractCalls로 호출 대상 추출
  │    ├─ called by: FindCallers로 역방향 추적
  │    └─ 재귀적으로 체인 순회 (depth 제한)
  │
  └─ FormatChain → stdout
```

### chain feature

```
filefunc chain feature validate
  │
  ├─ (동일 초기화)
  │
  ├─ FilterByFeature(파일들, "validate") → feature에 속하는 func 목록
  ├─ BuildCallGraph (feature 범위 내)
  │
  └─ FormatChain → stdout
```

## CLI 플래그

| 플래그 | 설명 | 기본값 | 범위 |
|---|---|---|---|
| `--chon` | 촌수 기반 탐색 (수직+수평). 1촌=부모+자식, 2촌=형제+조부모+손자, 3촌=삼촌+조카 | 1 | 1~3 |
| `--child-depth` | 아래로만 추적 (의존 체인) | — | 제한 없음 |
| `--parent-depth` | 위로만 추적 (영향 범위) | — | 제한 없음 |
| `--format` | 출력 형식 (text / json) | text | — |

`--chon`이 기본 모드. `--child-depth` 또는 `--parent-depth`를 지정하면 chon 대신 단방향 탐색.

### 촌수 정의

```
1촌: 직계 부모(caller) + 직계 자식(callee)
2촌: 조부모 + 손자 + 형제(같은 부모에게 호출되는 func)
3촌: 삼촌(부모의 형제) + 조카(형제의 자식) — 최대값
```

## 기존 코드 재사용

| 파일 | 용도 |
|---|---|
| `parse/extract_calls.go` | func body에서 호출 함수명 추출 |
| `parse/extract_uses.go` | func에서 사용 타입명 추출 |
| `parse/collect_project_symbols.go` | 프로젝트 심볼 수집 |
| `parse/build_import_map.go` | import 맵 생성 |
| `parse/call_name.go` | CallExpr에서 함수명 추출 |
| `cli/find_go_mod.go` | go.mod 탐색 |

## 완료 기준
- `filefunc chain func <name>` 동작
- `filefunc chain feature <name>` 동작
- 재귀 추적 depth 제한 동작
- filefunc 자체 코드에서 chain 실행하여 검증
- filefunc 자체 코드가 전체 validate 통과
- README.md, manual-for-ai.md 업데이트
