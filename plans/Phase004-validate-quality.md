# Phase 004: 코드 품질 룰 검증 ✅ 완료

## 목표
코드 품질 룰(Q1~Q3) 검증 구현.

## 산출물

| 파일 | 룰 |
|---|---|
| `internal/validate/check_nesting_depth.go` | Q1: nesting depth ≤ 2 |
| `internal/validate/check_func_lines.go` | Q2/Q3: func line count (ERROR 1000, WARNING 100) |

## Q1 구현 방식
go/ast로 AST를 순회하며 depth를 카운트한다. if/for/switch/select 진입 시 depth++, 탈출 시 depth--. 최대 depth > 2이면 ERROR.

## Q2/Q3 구현 방식
go/ast에서 FuncDecl의 시작/끝 위치로 라인 수 계산. 1000 초과 ERROR, 100 초과 WARNING.

## 완료 기준
- depth 3 코드에 대해 Q1 ERROR 출력
- 1000줄 초과 func에 대해 Q2 ERROR 출력
- 100줄 초과 func에 대해 Q3 WARNING 출력
- 각 룰 테스트 통과
