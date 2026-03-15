# mutest — 뮤테이션 테스트

## 개념

각 룰을 정확히 하나씩 위반하는 파일을 만들고, validate가 해당 위반을 감지하는지 확인한다. 룰당 1개의 뮤턴트 파일, 1개의 테스트.

## 위치

- 뮤턴트 파일: `internal/validate/testdata/`
- 테스트: `internal/validate/mutest_test.go`
- WalkGoFiles는 testdata를 건너뛰므로 자기검증에 영향 없음

## 전체 현황

| 룰 | 뮤턴트 파일 | 위반 내용 | 상태 |
|---|---|---|---|
| F1 | `multi_func.go` | func 2개 | ✅ |
| F2 | `multi_type.go` | exported type 2개 | ✅ |
| F3 | `multi_method.go` | method 2개 | ✅ |
| F4 | `init_alone.go` | init()만 단독 존재 | ✅ |
| Q1 | `deep_nesting.go` | nesting depth 3 | ✅ |
| Q2 | `long_func.go` | func 1010줄 | ✅ |
| Q3 | `medium_func.go` | func 110줄 (WARNING) | ✅ |
| A1 | `no_annotation.go` | //ff:func 없음 | ✅ |
| A2 | `bad_codebook_value.go` | 코드북에 없는 feature값 | ✅ |
| A3 | `no_what.go` | //ff:what 없음 | ✅ |
| A6 | `annotation_after_func.go` | 어노테이션이 func 뒤에 위치 | ✅ |

정상 파일: `clean.go` — 모든 룰 통과 확인용.

## 원칙

- 뮤턴트 1개 = 룰 위반 1개. 복합 위반 금지.
- 새 룰 추가 시 반드시 뮤턴트와 테스트를 함께 작성.
