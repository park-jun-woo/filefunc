# Phase 009: codebook.yaml 형식 검증 ✅ 완료

## 목표
codebook.yaml 로드 시 형식과 내용을 검증한다. 잘못된 codebook은 모든 후속 검증의 신뢰성을 해치므로 조기에 차단.

## 검증 룰

| # | 룰 | 위반 시 |
|---|---|---|
| C1 | feature, type 키 필수 (최소 1개 값) | ERROR |
| C2 | 동일 키 내 중복 값 불허 | ERROR |
| C3 | 모든 값은 소문자 + 하이픈만 허용 (`[a-z][a-z0-9-]*`) | ERROR |

빈 배열이면 아예 키를 제거하는 게 올바르다. C1에서 빈 배열과 키 누락을 동일하게 처리.

## 산출물

| 파일 | 내용 |
|---|---|
| `internal/validate/check_codebook_required_keys.go` | C1: feature, type 최소 1개 값 필수 |
| `internal/validate/check_codebook_duplicates.go` | C2: 동일 키 내 중복 값 검출 |
| `internal/validate/check_codebook_value_format.go` | C3: 값 형식 검증 (소문자 + 하이픈) |
| `internal/validate/validate_codebook.go` | ValidateCodebook — C1~C3 실행 오케스트레이터 |

## 실행 흐름

```
filefunc validate ./internal/
  │
  ├─ codebook.yaml 로드 (ParseCodebook)
  │
  ├─ ValidateCodebook(cb) → []Violation
  │    ├─ C1: feature, type 최소 1개 값?
  │    ├─ C2: 각 키 내 중복 값 없는가?
  │    └─ C3: 각 값이 [a-z][a-z0-9-]* 형식인가?
  │
  ├─ codebook 위반 있으면 즉시 출력 + exit 1 (후속 룰 실행 안 함)
  │
  └─ codebook 정상이면 기존 룰 (F1~A7) 실행
```

## CLI 연결
- validateCmd에서 ParseCodebook 직후 ValidateCodebook 호출
- codebook 위반이 있으면 파일 검증을 진행하지 않고 즉시 종료

## 기존 파일 수정

| 파일 | 변경 |
|---|---|
| `internal/cli/validate.go` | ParseCodebook 직후 ValidateCodebook 호출 추가 |
| `codebook.yaml` | level 값 소문자로 변경 (error, warning, info) |

## 완료 기준
- C1: feature 또는 type 누락/빈배열 시 ERROR
- C2: 중복 값 시 ERROR
- C3: 대문자, 공백, 특수문자 포함 시 ERROR
- codebook 위반 시 후속 파일 검증 미실행
- 뮤테이션 테스트 추가
- filefunc 자체 codebook.yaml이 검증 통과
- README.md 업데이트 (C1~C3 룰 추가)
- artifacts/manual-for-ai.md 업데이트 (codebook 형식 규칙 추가)
