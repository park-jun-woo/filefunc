# Phase 005: 어노테이션 룰 검증 ✅ 완료

## 목표
어노테이션 룰(A1~A3, A6) 검증 구현. 코드북과 어노테이션의 정합성을 검증한다.

## 산출물

| 파일 | 룰 |
|---|---|
| `internal/validate/check_annotation_required.go` | A1: func이 있는 파일은 //ff:func 필수 |
| `internal/validate/check_codebook_values.go` | A2: 어노테이션 값이 코드북에 존재하는지 |
| `internal/validate/check_what_required.go` | A3: func이 있는 파일은 //ff:what 필수 |
| `internal/validate/check_annotation_position.go` | A6: 어노테이션 파일 최상단 위치 |

## A2 구현 방식
ParseAnnotation 결과의 Func 맵에서 feature, type 등 키-값을 추출하고, ParseCodebook 결과의 허용 목록과 대조. 코드북에 없는 값이면 ERROR.

## A6 구현 방식
//ff: 주석이 package 선언과 import 사이 또는 package 선언 전에 위치해야 한다. func/type 선언보다 뒤에 있으면 ERROR.

## 완료 기준
- func이 있는데 어노테이션 없는 .go 파일에 대해 A1 ERROR
- 코드북에 없는 값 사용 시 A2 ERROR
- //ff:what 없는 파일에 대해 A3 ERROR
- 어노테이션이 func 선언 아래에 있으면 A6 ERROR
- filefunc 자체 코드가 전체 validate 통과
- 각 룰 테스트 통과
