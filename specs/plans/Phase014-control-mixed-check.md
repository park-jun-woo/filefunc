# Phase 014: control 혼합 금지 검증 (A13, A14) ✅ 완료

## 목표
"1 func 1 control" 원칙을 룰 레벨에서 완전히 강제한다. 기존 A10~A12는 "있어야 할 것이 있는가"만 검사. A13/A14는 "있으면 안 되는 것이 있는가"를 검사.

## 빈틈 분석

| 현재 룰 | 검사 내용 | 빈틈 |
|---|---|---|
| A10 | selection인데 switch 없음 | selection인데 loop 있음 미검사 |
| A11 | iteration인데 loop 없음 | iteration인데 switch 있음 미검사 |
| A12 | sequence인데 switch/loop 존재 | (빈틈 없음) |

## 추가 룰

| # | 룰 | 위반 시 | 검증 방법 |
|---|---|---|---|
| A13 | `control=selection`인데 depth 1에 loop 존재 | ERROR | AST: FuncDecl body 직계 자식에 ForStmt/RangeStmt 유무 |
| A14 | `control=iteration`인데 depth 1에 switch 존재 | ERROR | AST: FuncDecl body 직계 자식에 SwitchStmt/TypeSwitchStmt 유무 |

## 완전성 검증표

A13/A14 추가 후 모든 조합이 빈틈 없이 검증:

| 파일 내용 | control=sequence | control=selection | control=iteration |
|---|---|---|---|
| switch + loop | A12 ERROR | A13 ERROR | A14 ERROR |
| switch만 | A12 ERROR | ✅ | A14 ERROR |
| loop만 | A12 ERROR | A13 ERROR | ✅ |
| 둘 다 없음 | ✅ | A10 ERROR | A11 ERROR |

mixed 파일은 어떤 control 값을 붙여도 통과 불가. 분리하지 않으면 검증 통과 불가능.

## 산출물

| 파일 | 내용 | 구분 |
|---|---|---|
| `internal/validate/check_control_selection_no_loop.go` | A13: selection인데 loop 있으면 ERROR | 신규 |
| `internal/validate/check_control_iteration_no_switch.go` | A14: iteration인데 switch 있으면 ERROR | 신규 |
| `internal/parse/has_loop_at_depth1.go` | HasLoopAtDepth1 — depth 1에 loop 존재 여부 | 신규 |
| `internal/parse/has_switch_at_depth1.go` | HasSwitchAtDepth1 — depth 1에 switch 존재 여부 | 신규 |
| `internal/validate/run_all.go` | A13~A14 호출 추가 | 수정 |
| `internal/validate/check_control_selection_no_loop_test.go` | A13 테스트 | 신규 |
| `internal/validate/check_control_iteration_no_switch_test.go` | A14 테스트 | 신규 |
| `internal/validate/testdata/control_selection_with_loop.go` | A13 위반 테스트 데이터 (selection인데 loop) | 신규 |
| `internal/validate/testdata/control_iteration_with_switch.go` | A14 위반 테스트 데이터 (iteration인데 switch) | 신규 |

## 구현 참고

- `internal/parse/detect_from_body.go`의 `detectFromBody()`가 이미 depth 1에서 switch/loop을 판별한다. `HasLoopAtDepth1`/`HasSwitchAtDepth1`는 동일 로직의 bool 반환 변형이므로, 파일 파싱 부분은 `detect_control.go`의 패턴을 따른다.

## 문서 업데이트

| 파일 | 변경 내용 |
|---|---|
| `CLAUDE.md` | A13~A14 룰 추가 |
| `README.md` | Annotation 룰 테이블에 A13~A14 추가 |
| `artifacts/manual-for-ai.md` | Rules 테이블에 A13~A14 추가 |

## 완료 기준
- A13: control=selection인데 loop 존재하면 ERROR
- A14: control=iteration인데 switch 존재하면 ERROR
- 완전성 검증표의 모든 조합 통과
- filefunc 자체 코드 validate 통과
- CLAUDE.md, README.md, manual-for-ai.md 업데이트
