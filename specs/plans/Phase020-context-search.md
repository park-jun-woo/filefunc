# Phase 020: filefunc context --search (직접 필터 지정) ✅ 완료

## 목표

`filefunc context`에 `--search` 플래그를 추가하여 LLM feature 선택을 건너뛰고 사용자가 직접 어노테이션 필터를 지정할 수 있게 한다.

## 배경

현재 context는 1단계에서 LLM이 codebook feature를 선택한다. 하지만 사용자가 이미 대상 영역을 아는 경우:
- LLM 호출 1회 절약 (~3초)
- 정확도 100% (LLM 편향 없음)
- fullend처럼 codebook이 복잡한 프로젝트에서 feature + ssot 조합 직접 지정 가능

## 설계

### 사용법

```bash
# LLM feature 선택 (기존)
filefunc context "SSaC 검증 룰 추가"

# 직접 필터 지정 (신규)
filefunc context "SSaC 검증 룰 추가" --search "feature=ssac-validate"
filefunc context "크로스체크 수정" --search "feature=crosscheck ssot=openapi"
filefunc context "DDL 코드젠 수정" --search "feature=gen-gogin ssot=ddl"
filefunc context "상태머신 검증" --search "feature=statemachine type=rule"
```

### 동작

```
--search 없음:
  1단계: LLM feature 선택 (LLM 1회)
  2단계: feature 필터
  3단계: what 스코어링 (LLM 1회)
  4단계: body 스코어링 (LLM 1회)

--search 있음:
  1단계: 스킵
  2단계: --search 값으로 어노테이션 매칭 필터
  3단계: what 스코어링 (LLM 1회)
  4단계: body 스코어링 (LLM 1회)
```

### --search 매칭 로직

`--search "feature=ssac-validate ssot=openapi"` → 공백 구분 key=value 쌍.

GoFile의 `Annotation.Func` (또는 `.Type`) 맵에서 **모든 key=value가 일치**하는 파일만 통과 (AND 조건).

- required 키 (feature, type): `Annotation.Func["feature"]` 매칭
- optional 키 (ssot, pattern): `Annotation.Func["ssot"]` 매칭
- 지정하지 않은 키는 무시 (부분 매칭)

### 출력 포맷

```
[1/4] search filter: feature=ssac-validate → 30 files
[2/4] (skipped — direct search)
[3/4] what scoring: 30 → 5 (rate≥0.2)
[4/4] body scoring: 5 → 3 (rate≥0.5)

Results:
  ...
```

## 산출물

| 파일 | 내용 | 구분 |
|---|---|---|
| `internal/context/parse_search.go` | --search 문자열을 key=value 맵으로 파싱 | 신규 |
| `internal/context/filter_search.go` | GoFile 어노테이션에서 key=value AND 매칭 필터 | 신규 |
| `internal/context/run_pipeline.go` | PipelineConfig에 Search 필드 추가, --search 있으면 1단계 스킵 | 수정 |
| `internal/cli/context.go` | --search 플래그 추가 | 수정 |

## 문서 업데이트

| 파일 | 변경 내용 |
|---|---|
| `CLAUDE.md` | context --search 사용법 추가 |
| `README.md` | context 섹션에 --search 추가 |
| `artifacts/manual-for-ai.md` | context --search 사용 예시 추가 |

## 구현 순서

1. `parse_search.go` 구현 — "feature=X ssot=Y" → map[string]string
2. `filter_search.go` 구현 — GoFile 어노테이션 AND 매칭
3. `run_pipeline.go` 수정 — --search 분기 추가
4. `context.go` 수정 — --search 플래그 추가
5. 문서 업데이트
6. `filefunc validate` 위반 0 확인

## 완료 기준

- `--search "feature=X"` 지정 시 LLM 1단계 스킵, 직접 필터링
- 복수 key=value AND 조건 매칭
- `--search` 없으면 기존 LLM feature 선택 동작 유지
- filefunc 자체 코드 validate 통과
- CLAUDE.md, README.md, manual-for-ai.md 업데이트
