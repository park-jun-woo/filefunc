# Phase 021: codebook description 필드 도입 ✅ 완료

## 목표

codebook.yaml의 각 값에 description을 정식 필드로 추가하여 파서가 읽을 수 있게 하고, validate로 검증, context가 LLM에 정확히 전달할 수 있게 한다.

## 배경

### 현재 문제

- codebook 주석(`#`)은 YAML 파서가 버림 → `model.Codebook`에 description 없음
- `filefunc context`가 주석을 읽기 위해 원본 텍스트를 해킹 파싱 (`extractFeatureDescriptions`)
- description 유무/품질에 대한 검증 없음
- LLM feature 선택의 정확도가 description 품질에 직접 의존

### 해결

codebook 형식을 변경하여 description을 값으로 포함:

```yaml
# Before
feature:
  - validate    # 코드 구조 룰 검증

# After
feature:
  validate: "코드 구조 룰 검증 (F1,Q1,A1 등 정적 분석 룰)"
  annotate: "LLM what-body 검증 및 어노테이션 자동 생성 (llmc)"
```

## 설계

### codebook.yaml 새 형식

```yaml
required:
  feature:
    validate: "코드 구조 룰 검증 (F1,Q1,A1 등 정적 분석 룰)"
    annotate: "LLM what-body 검증 및 어노테이션 자동 생성 (llmc)"
    chain: "func/feature 데이터 흐름 추적"
    parse: "소스 코드, 어노테이션, codebook 파싱"
    codebook: "codebook 로드 및 관리"
    report: "검증 결과 출력"
    cli: "cobra 명령 정의 및 llmc/validate/chain/context 실행"
    context: "LLM 기반 컨텍스트 탐색"
  type:
    command: "cobra 명령 엔트리포인트"
    rule: "개별 검증 룰 구현"
    parser: "파싱 로직"
    walker: "파일/디렉토리 순회"
    model: "데이터 구조체"
    formatter: "출력 포매터"
    loader: "파일/설정 로드"
    util: "범용 유틸리티"

optional:
  pattern:
    error-collection: "에러 수집 후 일괄 보고"
    file-visitor: "파일 단위 순회 처리"
    rule-registry: "룰 등록/실행"
  level:
    error: ""
    warning: ""
    info: ""
```

- `[]string` → `map[string]string` (값 → 값:설명)
- description이 빈 문자열이어도 허용 (level 등 자명한 값)
- required 키의 description은 필수 검증 대상 (새 룰 C4)

### model.Codebook 변경

```go
// Before
type Codebook struct {
    Required map[string][]string `yaml:"required"`
    Optional map[string][]string `yaml:"optional"`
}

// After
type Codebook struct {
    Required map[string]map[string]string `yaml:"required"`
    Optional map[string]map[string]string `yaml:"optional"`
}
```

값 접근: `cb.Required["feature"]` → `map[string]string{"validate": "코드 구조 룰 검증 ..."}`
키 목록: `for name := range cb.Required["feature"]`

### 영향 범위

Codebook 구조체를 사용하는 모든 코드:

| 파일 | 변경 내용 |
|---|---|
| `internal/model/codebook.go` | `[]string` → `map[string]string` |
| `internal/parse/parse_codebook.go` | YAML 파싱 변경 |
| `internal/validate/check_codebook_values.go` | 값 매칭 로직 변경 |
| `internal/validate/check_codebook_duplicates.go` | 원본 텍스트에서 중복 키 검출로 변경 (YAML map은 조용히 덮어쓰므로) |
| `internal/validate/check_codebook_required_keys.go` | map 키 존재 확인으로 변경 |
| `internal/validate/check_codebook_value_format.go` | map 키 순회로 변경 |
| `internal/validate/check_required_keys_in_annotation.go` | 값 목록 접근 변경 |
| `internal/validate/allowed_values.go` | map 키를 슬라이스로 변환 |
| `internal/validate/find_duplicates.go` | 삭제 또는 수정 (map이면 중복 불가) |
| `internal/validate/check_values_format.go` | map 키 순회로 변경 |
| `internal/validate/validate_codebook.go` | C4 룰 추가 호출 |
| `internal/context/select_feature.go` | codebookRaw 해킹 제거, Codebook 직접 사용 |
| `internal/context/build_feature_prompt.go` | description을 Codebook에서 직접 읽기 |
| `internal/context/extract_feature_descriptions.go` | 삭제 |
| `internal/context/parse_feature_line.go` | 삭제 |
| `internal/cli/context.go` | codebookRaw → Codebook 전달로 변경 |

### 새 검증 룰

| # | 룰 | 위반 시 | 검증 방법 |
|---|---|---|---|
| C4 | required 키의 각 값에 description이 비어있지 않아야 함 | WARNING | map value 검사 |

WARNING으로 — description이 없어도 동작은 하지만 context 정확도가 떨어짐.

## 산출물

| 파일 | 내용 | 구분 |
|---|---|---|
| `codebook.yaml` | 새 형식으로 변환 (filefunc 자체) | 수정 |
| `internal/model/codebook.go` | `map[string]map[string]string` 타입 변경 | 수정 |
| `internal/parse/parse_codebook.go` | YAML 파싱 변경 | 수정 |
| `internal/validate/allowed_values.go` | map 키 → 슬라이스 변환 | 수정 |
| `internal/validate/check_codebook_values.go` | map 키 매칭 | 수정 |
| `internal/validate/check_codebook_duplicates.go` | 원본 텍스트에서 중복 키 검출로 변경 | 수정 |
| `internal/validate/check_codebook_required_keys.go` | map 키 검사 | 수정 |
| `internal/validate/check_codebook_value_format.go` | map 키 순회 | 수정 |
| `internal/validate/check_values_format.go` | map 키 순회 | 수정 |
| `internal/validate/find_duplicates.go` | 원본 텍스트 기반 중복 키 검출로 변경 | 수정 |
| `internal/validate/mutest_test.go` | Codebook 생성 코드 `map[string]map[string]string`으로 변경 | 수정 |
| `internal/validate/check_required_keys_in_annotation.go` | map 키 접근 | 수정 |
| `internal/validate/validate_codebook.go` | C4 추가 | 수정 |
| `internal/validate/check_codebook_description.go` | C4: required description 검증 | 신규 |
| `internal/context/build_feature_prompt.go` | Codebook.Required에서 description 직접 읽기 | 수정 |
| `internal/context/select_feature.go` | codebookRaw → Codebook 사용 | 수정 |
| `internal/context/run_pipeline.go` | PipelineConfig.CodebookRaw 제거, Codebook 추가 | 수정 |
| `internal/context/extract_feature_descriptions.go` | 삭제 | 삭제 |
| `internal/context/parse_feature_line.go` | 삭제 | 삭제 |
| `internal/cli/context.go` | Codebook 전달로 변경 | 수정 |

## 문서 업데이트

| 파일 | 변경 내용 |
|---|---|
| `CLAUDE.md` | codebook 형식 변경 설명, C4 룰 추가 |
| `README.md` | Codebook 섹션 형식 변경, C4 추가 |
| `artifacts/manual-for-ai.md` | Codebook 형식 변경, C4 추가 |

## 구현 순서

1. `model/codebook.go` 타입 변경
2. `parse/parse_codebook.go` YAML 파싱 변경
3. `codebook.yaml` 새 형식으로 변환 (filefunc 자체)
4. validate 패키지 전체 수정 (Codebook 접근 방식 변경)
5. `find_duplicates.go` + `check_codebook_duplicates.go` 수정 — 원본 텍스트 기반 중복 키 검출
6. `check_codebook_description.go` 신규 (C4)
7. `mutest_test.go` 수정 — Codebook 생성 코드 변경
8. context 패키지 수정 (codebookRaw → Codebook)
9. `extract_feature_descriptions.go`, `parse_feature_line.go` 삭제
10. `cli/context.go` 수정
11. 문서 업데이트
12. `filefunc validate` 위반 0 확인
13. whyso, fullend codebook.yaml도 새 형식으로 변환

## 완료 기준

- codebook.yaml이 `map[string]string` 형식
- `model.Codebook`에서 description 직접 접근 가능
- C4: required 키 description 비어있으면 WARNING
- context가 Codebook에서 description을 직접 읽어 LLM에 전달
- `extractFeatureDescriptions` 해킹 코드 삭제
- filefunc 자체 + whyso + fullend codebook 변환 완료
- filefunc validate 위반 0
- CLAUDE.md, README.md, manual-for-ai.md 업데이트

## 하위 호환성

codebook 형식이 바뀌므로 **기존 codebook.yaml이 파싱 실패**합니다. 마이그레이션 필요:
- filefunc, whyso, fullend 전부 변환
- README에 마이그레이션 안내 추가
