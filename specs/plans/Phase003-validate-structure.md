# Phase 003: 파일 구조 룰 검증 ✅ 완료

## 목표
파일 구조 룰(F1~F4) 검증 구현. `filefunc validate`가 실제로 위반을 탐지하고 보고한다.

## 산출물

| 파일 | 룰 |
|---|---|
| `internal/validate/check_one_file_one_func.go` | F1: 파일당 func 1개 |
| `internal/validate/check_one_file_one_type.go` | F2: 파일당 type 1개 |
| `internal/validate/check_one_file_one_method.go` | F3: 파일당 method 1개 |
| `internal/validate/check_init_standalone.go` | F4: init()만 단독 불허 (var 또는 func과 함께) |
| `internal/validate/run_all.go` | RunAll — 전체 룰 실행 오케스트레이터 |
| `internal/report/format_text.go` | FormatText — 텍스트 출력 |

## 예외 처리
- F5: _test.go는 복수 func 허용 → F1 검사에서 _test.go 제외
- F6: 함수 전용 파라미터 타입은 해당 func 파일에 허용 → F2 검사에서 unexported type 예외
- F7: 의미적 const 묶음은 같은 파일 허용 → const 파일은 F1 검사 제외

## CLI 연결
- validateCmd가 RunAll을 호출하고 FormatText로 결과 출력
- 위반 시 exit code 1

## 완료 기준
- 위반 코드에 대해 정확한 ERROR 출력
- 정상 코드에 대해 위반 없음 출력
- filefunc 자체 코드가 validate 통과
- 각 룰 테스트 통과
