# Phase 010: codebook required/optional 체계 ✅ 완료

## 목표
codebook.yaml을 required/optional로 구분한다. required 키는 모든 `//ff:func`, `//ff:type` 어노테이션에 필수. optional은 해당할 때만. validate가 required 키 누락을 ERROR로 잡아 grep 신뢰성을 보장한다.

## codebook.yaml 구조 변경

### Before
```yaml
feature: [validate, parse, cli]
type: [command, rule, parser]
pattern: [error-collection]
level: [error, warning]
```

### After
```yaml
required:
  feature: [validate, annotate, chain, parse, codebook, report, cli]
  type: [command, rule, parser, walker, model, formatter, loader, util]

optional:
  pattern: [error-collection, file-visitor, rule-registry]
  level: [error, warning, info]
```

## 검증 룰 변경

| 룰 | 변경 |
|---|---|
| A2 | required + optional 전체에서 값 존재 확인 (기존과 동일, 범위만 확장) |
| A8 (신규) | required 키가 `//ff:func` 또는 `//ff:type`에 모두 존재하는지 검증. 누락 시 ERROR |
| C1 | required 섹션에 최소 1개 키 필수. required 내 각 키는 최소 1개 값 필수 |
| C2, C3 | required + optional 전체 대상으로 동일하게 적용 |

## 산출물

| 파일 | 내용 |
|---|---|
| `internal/model/codebook.go` | Codebook 구조체 변경 (Required/Optional 중첩) |
| `internal/parse/parse_codebook.go` | 새 구조 파싱 대응 |
| `internal/validate/check_required_keys_in_annotation.go` | A8: 어노테이션에 required 키 누락 검증 |
| `internal/validate/allowed_values.go` | required + optional 통합 조회로 변경 |
| `internal/validate/check_codebook_required_keys.go` | C1: required 섹션 검증으로 변경 |
| `internal/validate/check_codebook_duplicates.go` | required + optional 전체 대상 |
| `internal/validate/check_codebook_value_format.go` | required + optional 전체 대상 |
| `internal/validate/run_all.go` | A8 호출 추가 |

## 실행 흐름

```
filefunc validate ./internal/
  │
  ├─ codebook.yaml 로드 (새 구조)
  ├─ ValidateCodebook → C1~C3 (required 섹션 필수, 전체 중복/형식 검사)
  │
  ├─ 각 파일에 대해:
  │    ├─ A2: 어노테이션 값이 required + optional에 존재하는지
  │    └─ A8: required 키가 어노테이션에 모두 있는지
  │
  └─ ... 기존 룰
```

## AI 에이전트 grep 신뢰성

required 키는 모든 어노테이션에 있으므로 grep이 완전하다.
```bash
rg '//ff:func feature=validate'       # required → 빠짐 없음
rg '//ff:func pattern=error-collection'  # optional → 일부만 매칭 (정상)
```

## 완료 기준
- codebook.yaml이 required/optional 구조로 변경
- A8: required 키 누락 시 ERROR
- A2: required + optional 통합 조회
- C1~C3: 새 구조 대응
- 기존 mutest + 새 뮤테이션 테스트 통과
- filefunc 자체 코드가 전체 validate 통과
- README.md, manual-for-ai.md 업데이트
