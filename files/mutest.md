# mutest — 뮤테이션 테스트

## 개념

각 룰을 정확히 하나씩 위반하는 파일을 만들고, validate가 해당 위반을 감지하는지 확인한다. 룰당 1개의 뮤턴트 파일, 1개의 테스트.

## 위치

- 뮤턴트 파일: `internal/validate/testdata/`
- 테스트: `internal/validate/mutest_test.go`
- WalkGoFiles는 testdata를 건너뛰므로 자기검증에 영향 없음

## 전체 현황

### 파일 구조 룰

| 룰 | 뮤턴트 파일 | 위반 내용 | 상태 |
|---|---|---|---|
| F1 | `multi_func.go` | func 2개 | ✅ |
| F2 | `multi_type.go` | exported type 2개 | ✅ |
| F3 | `multi_method.go` | method 2개 | ✅ |
| F4 | `init_alone.go` | init()만 단독 존재 | ✅ |

### 코드 품질 룰

| 룰 | 뮤턴트 파일 | 위반 내용 | 상태 |
|---|---|---|---|
| Q1 | `deep_nesting.go` | nesting depth 3 | ✅ |
| Q1 (dimension) | `dimension2_depth3.go` | dimension=2, depth 3 → 통과 확인 | ✅ |
| Q2 | `long_func.go` | func 1010줄 | ✅ |
| Q3 | `medium_func.go` | func 110줄 (WARNING) | ✅ |
| Q3 (backtick) | `q3_backtick.go` | 100줄+ backtick 포함 → hint 메시지 확인 | ✅ |

### 어노테이션 룰

| 룰 | 뮤턴트 파일 | 위반 내용 | 상태 |
|---|---|---|---|
| A1 | `no_annotation.go` | //ff:func 없음 | ✅ |
| A2 | `bad_codebook_value.go` | 코드북에 없는 feature값 | ✅ |
| A3 | `no_what.go` | //ff:what 없음 | ✅ |
| A6 | `annotation_after_func.go` | 어노테이션이 func 뒤에 위치 | ✅ |
| A9 | `no_control.go` | control= 없음 | ✅ |
| A10 | `selection_no_switch.go` | selection인데 switch 없음 | ✅ |
| A11 | `iteration_no_loop.go` | iteration인데 loop 없음 | ✅ |
| A12 | `sequence_with_loop.go` | sequence인데 loop 존재 | ✅ |
| A13 | `control_selection_with_loop.go` | selection인데 loop 존재 | ✅ |
| A14 | `control_iteration_with_switch.go` | iteration인데 switch 존재 | ✅ |
| A15 | `iter_no_dimension.go` | iteration인데 dimension 없음 | ✅ |
| A16 | `bad_dimension_value.go` | dimension=0 (양의 정수 아님) | ✅ |

### 코드북 룰

| 룰 | 위반 내용 | 상태 |
|---|---|---|
| C1 | required 섹션 비어있음 | — |
| C2 | 동일 섹션 내 중복 키 (원본 텍스트 기반) | — |
| C3 | 키가 소문자+하이픈 아님 | — |
| C4 | required description 비어있음 (WARNING) | — |

### 미구현 뮤턴트

| 룰 | 위반 내용 | 비고 |
|---|---|---|
| A7 | //ff:checked 해시 불일치 | 프로젝트에 checked 없으면 스킵 |
| A8 | required 키 누락 | codebook 의존 |

정상 파일: `clean.go` — 모든 룰 통과 확인용.

## 원칙

- 뮤턴트 1개 = 룰 위반 1개. 복합 위반 금지.
- 새 룰 추가 시 반드시 뮤턴트와 테스트를 함께 작성.
