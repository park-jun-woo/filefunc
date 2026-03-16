# test — 테스트 전략

## 테스트 분류

### 1. mutest (뮤테이션 테스트)

룰을 정확히 하나씩 위반 → validate가 감지하는가?

상세: `files/mutest.md`

현재 구현: F1~F4, Q1~Q3, A1~A3, A6, A13~A16, Q1 dimension

### 2. exception (예외 통과 테스트)

룰 예외 조건에서 위반이 발생하지 **않는가**?

| 예외 | 테스트 파일 | 검증 내용 |
|---|---|---|
| F5 | `test_file_test.go` | _test.go에 func 여러 개 → F1 미감지 |
| F6 | `func_with_param_type.go` | func + unexported type 함께 → F2 미감지 |
| F7 | `const_only.go` | const만 있는 파일 → F1 미감지 |
| F4 예외 | `var_with_init.go` | var + init() 조합 → F4 미감지 |

### 3. clean (정상 통과 테스트)

clean.go가 **모든** 룰을 통과하는지 일괄 검증.

```go
func TestClean_AllRules(t *testing.T) {
    gf := mustParse(t, "testdata/clean.go")
    violations := RunAll([]*model.GoFile{gf}, nil)
    expectNoViolation(t, violations)
}
```

### 4. //ff:type 경로 테스트

| 테스트 | 검증 내용 |
|---|---|
| A1 type 미감지 | type-only 파일에 //ff:type 없으면 A1 ERROR |
| A1 type 정상 | type-only 파일에 //ff:type 있으면 통과 |
| A2 type 코드북 | //ff:type의 값이 코드북에 없으면 A2 ERROR |

### 5. dimension (Q1 동적 상한 테스트)

| 테스트 | 검증 내용 |
|---|---|
| Q1 dimension 통과 | `dimension2_depth3.go` — dimension=2, depth 3 → Q1 통과 |
| A15 위반 | `iter_no_dimension.go` — iteration인데 dimension 없음 → ERROR |
| A16 위반 | `bad_dimension_value.go` — dimension=0 → ERROR |

## 테스트 위치

- 뮤턴트 파일: `internal/validate/testdata/`
- 뮤테이션 테스트: `internal/validate/mutest_test.go`
- 예외/clean/type 테스트: `internal/validate/exception_test.go`

## 원칙

- 뮤턴트 1개 = 룰 위반 1개
- 예외 테스트 1개 = 예외 조건 1개
- 새 룰 추가 시 mutest + exception 함께 작성
