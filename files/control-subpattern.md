# control 세부 패턴 분석

## 배경

fullend 코드베이스(1,081개 함수)를 filefunc으로 리팩토링한 결과, 모든 함수가 3가지 control 타입(selection, sequence, iteration)으로 수렴함을 확인했다. 여기서 더 깊이 들어가, 같은 control 타입 내에서 반복되는 세부 패턴을 분석했다.

## 분포

| control | 파일 수 | 비율 |
|---|---|---|
| iteration | 603 | 56% |
| sequence | 403 | 37% |
| selection | 75 | 7% |

---

## selection (75개) — 입력 하나 → 분기 → 출력 하나

| 패턴 | 개수 | 설명 | 대표 예시 |
|---|---|---|---|
| 타입 매핑 | 13 | 타입 A → 타입 B 문자열 변환 | `pgTypeToGo`, `oaTypeToGo` |
| 문자열 패턴 분기 | 7 | prefix/suffix 매칭 → 변환 | `inputValueToCode`, `stripDirectivePrefix` |
| 열거형 → 위임 | 10 | enum/keyword → 하위 함수 호출 | `parseAnnotation`, `templateName` |
| 타입별 코드 생성 | 7 | Go 타입 → 코드 스니펫 | `generateExtractCode`, `generatePathParamCode` |
| AST 타입 분기 | 6 | `.(type)` switch | `exprToString`, `processDecl` |
| HTML 속성 디스패치 | 9 | `data-*` 유무 → 파서 분기 (STML 전용) | `dispatchFetchChild`, `dispatchActionChild` |
| 백엔드 디스패치 | 6 | "postgres"/"memory" → 구현체 분기 | `initQueue`, `publish` |
| 구조체 변이 | 5 | 조건별 필드 설정 | `applyAnnotation`, `buildErrTracking` |
| 기타 | 12 | CLI 분기, 포맷 분기, 복합 조건 등 | `main.go`, `formatHistory` |

순수 함수 비율: 63% (47/75).

가장 큰 군은 "타입 매핑" — A→B 변환 테이블로, 부작용 없이 switch/case만으로 구성된다.

---

## sequence (403개) — 순서대로 단계 실행

| 패턴 | 추정 수 | 설명 | 대표 예시 |
|---|---|---|---|
| A. 단순 위임 | ~65 | 1줄 래퍼, 함수 호출 전달 | `Generate` → `GenerateWith` |
| B. 파싱 | ~85 | file/string/AST → struct | `loadConfig`, `parseCallExpr`, `parseDir` |
| C. 검증 수집 | ~75 | struct → []error 누적 | `validateCallSpec`, `validateHurlEntry` |
| D. 구조체 조립 | ~55 | 빈 struct → 필드 채워서 반환 | `buildTemplateData`, `loadSymbolTable` |
| E. 버퍼 생성 | ~75 | Buffer/Builder에 코드 출력 | `assembleGoSource`, `buildHTTPFuncBody` |
| F. 파이프라인 | ~35 | step1→step2→...→stepN, 에러시 조기 복귀 | `genmodel/generate.go`, `runValidate` |
| G. 조회/해석 | ~45 | map 룩업 + 폴백 체인 | `lookupColumnType`, `resolveInputParamType` |

90%가 `if err != nil { return }` 패턴 사용.

feature별 지배적 패턴이 다르다:
- parser → B (파싱)
- validate/crosscheck → C (검증 수집)
- gen → E (버퍼 생성)
- orchestrator → F (파이프라인)

---

## iteration (603개) — 컬렉션 순회

| 패턴 | 추정 수 | 설명 | 대표 예시 |
|---|---|---|---|
| 1. collect-map | ~120 | 순회 → `map[K]V` 구축 | `buildOperationMap`, `collectDomainModels` |
| 2. filter | ~50 | 조건 불일치 제거 | `filterNonEmpty`, `filterAuthInputs` |
| 3. find-first | ~80 | 첫 매칭 즉시 return | `findFirstModelTable`, `hasDirectResponse` |
| 4. validate-all | ~110 | 각 항목 검증 → []error 누적 | `checkSsacOpenapi`, `validateWithSymbols` |
| 5. transform | ~60 | 각 요소 변환 → 새 슬라이스 | `toSnakeCase`, `renderChildNodes` |
| 6. reduce | ~40 | 컬렉션 → 스칼라 값 | `countDDLColumns`, `joinKeys` |
| 7. collect-slice | ~70 | 순회 → `[]T` 구축 | `extractMethods`, `detectSsots` |
| 8. recursive-walk | ~15 | 트리 재귀 순회 (DOM/AST) | `walkTopLevel`, `hasDescendantData` |
| 9. side-effect | ~45 | 항목마다 I/O 수행 | `printErrors`, `pollOnce` |
| 10. multi-stage | ~15 | 복수 루프 연쇄 (수집→검증) | `checkSsacOpenapi` (3단 루프) |

### dimension 분포

| dimension | 비율 | 의미 |
|---|---|---|
| 1 | 86% | 단일 `for range` 루프 |
| 2 | 12% | 2중 중첩 (files→lines, funcs→sequences) |
| 3+ | 2% | 3중 이상 (AST 깊은 순회 등) |

단일 루프가 압도적 — 함수 분해가 잘 되어 있다는 증거.

---

## 교차 분석: control × role

세부 패턴을 다시 묶으면, control(어떻게)과 직교하는 **role(무엇을)** 축이 드러난다.

| role | selection | sequence | iteration |
|---|---|---|---|
| **변환** | 타입 매핑, 코드생성 | 파싱(B), 조회(G) | transform, collect-map |
| **검증** | 검증 규칙 분기 | 검증 수집(C) | validate-all |
| **조립** | — | 버퍼 생성(E), 구조체 조립(D) | collect-slice, side-effect |
| **탐색** | — | — | find-first, filter |
| **조율** | 위임 디스패치 | 파이프라인(F), 단순 위임(A) | multi-stage |
| **집약** | — | — | reduce |

- **control** = 제어구조. 함수가 **어떻게** 흐르는가 (분기/순차/반복)
- **role** = 의미역. 함수가 **무엇을** 하는가 (변환/검증/조립/탐색/조율/집약)

이 두 축이 교차하면 함수를 더 정밀하게 분류할 수 있다.
