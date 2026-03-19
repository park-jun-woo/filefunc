# Phase 019: filefunc context v2 — LLM feature 선택 기반 탐색 ✅ 완료

## 목표

대상 함수명 없이 프롬프트만으로 관련 코드를 찾는 4단계 파이프라인. Phase 018의 chain 의존을 제거하고 LLM feature 선택으로 대체.

```
1단계: LLM feature 선택  → codebook에서 feature 고르기 (LLM 1회)
2단계: feature 필터      → 프로젝트 전체에서 해당 feature 파일만 (정적)
3단계: what 스코어링     → LLM 1회 배치, rate≥0.2 필터
4단계: body 스코어링     → LLM 1회 배치, rate≥0.5 필터
```

LLM 호출: 3회. 대상 함수명 불필요.

## 배경

### Phase 018의 한계

- chain 의존: 대상 함수명을 알아야 함
- 고립 함수(ParseCodebook 등) 탐색 불가
- feature 필터가 대상 함수의 feature를 그대로 사용 → LLM 판단 아님

### 실험 결과

- LLM(gpt-oss:20b)이 codebook에서 feature를 고르는 것은 3/3 정확 (files/LLM-탐색-실험.md)
- feature 하나로 전체 파일의 70~99% 제거 가능

### 역할 분리

```
filefunc chain   → 정적 분석 (함수명 필요, LLM 불필요, 즉시)
filefunc context → LLM 탐색 (프롬프트만 필요, 함수명 불필요, ~15초)
```

## 설계

### 사용법

```bash
filefunc context "nesting depth 검증 로직을 수정하려고 한다"
filefunc context "codebook 파싱 로직을 변경하고 싶다"
filefunc context "chain 출력 포맷을 변경하고 싶다"
```

대상 함수명 불필요. 프롬프트만으로 동작.

진입점을 모르면 `whyso map`으로 키워드 검색하여 참고 가능.

### 파이프라인

```
입력: 사용자 프롬프트 + codebook.yaml + 프로젝트 전체 소스

1단계 — LLM feature 선택 (LLM 1회)
  codebook.yaml + 프롬프트를 LLM에 전달
  → feature 값 1~3개 선택 (JSON)
  → 프로젝트 전체 GoFile에서 해당 feature 매칭

2단계 — feature 필터 (정적)
  1단계 선택된 feature에 해당하는 파일만 유지
  → 후보 10~40개

3단계 — what 스코어링 (LLM 1회 배치)
  2단계 통과 파일들의 함수명 + what을 배치 프롬프트로 전달
  → rate≥0.2 필터 (후보 4~10개)

4단계 — 본문 스코어링 (LLM 1회 배치)
  3단계 통과 파일들의 본문을 배치 프롬프트로 전달
  → rate≥0.5 필터 (최종 2~5개)

출력: 최종 함수 목록 + 점수 + what
```

### 플래그

```bash
filefunc context "프롬프트"                          # 기본 (4단계 전부)
filefunc context "프롬프트" --depth 1                 # 1단계만 (feature 선택)
filefunc context "프롬프트" --depth 2                 # 2단계까지 (feature 필터)
filefunc context "프롬프트" --depth 3                 # 3단계까지 (what 스코어링)
filefunc context "프롬프트" --depth 4                 # 4단계까지 (본문 스코어링, 기본값)
filefunc context "프롬프트" --what-rate 0.2           # 3단계 임계값 (기본 0.2)
filefunc context "프롬프트" --body-rate 0.5           # 4단계 임계값 (기본 0.5)
filefunc context "프롬프트" --model gpt-oss:20b       # ollama 모델
filefunc context "프롬프트" --endpoint http://...     # ollama 엔드포인트
```

### 출력 포맷

```
[1/4] feature selection: validate (LLM)
[2/4] feature filter: 121 → 34
[3/4] what scoring: 34 → 6 (rate≥0.2)
[4/4] body scoring: 6 → 3 (rate≥0.5)

Results:
  CheckNestingDepth [0.90] (what="Q1: control과 dimension 기반으로 nesting depth 상한 검증")
  depthLimit [0.60] (what="control과 dimension으로 Q1 depth 상한을 계산")
  CheckDimensionValue [0.50] (what="dimension= 값은 양의 정수여야 함")
```

### LLM 프롬프트

#### 1단계: feature 선택

```json
{
  "task": "Given codebook, select the most relevant features for the user's question. Use only codebook feature values.",
  "codebook": {"feature": ["validate","annotate","chain","parse","codebook","report","cli","context"]},
  "question": "사용자 프롬프트",
  "format": "JSON array of feature strings, max 3"
}
```

#### 3단계: what 스코어링

```
사용자가 수정하려는 작업과 각 함수의 관련도를 평가하시오.
관련도: 0.0(무관) ~ 0.5(간접 관련) ~ 1.0(직접 관련)
직접 수정 대상이면 0.8 이상, 영향을 받는 함수면 0.4~0.7, 무관하면 0.2 이하.

작업: "사용자 프롬프트"

1. FuncName: "what 텍스트"
...

각 번호에 대해 점수만. 형식: 번호. 점수
```

#### 4단계: 본문 스코어링 (동일 가이드라인)

## 산출물

Phase 018의 `internal/context/` 패키지를 수정.

| 파일 | 내용 | 구분 |
|---|---|---|
| `internal/context/select_feature.go` | 1단계: LLM으로 codebook feature 선택 | 신규 |
| `internal/context/build_feature_prompt.go` | feature 선택용 프롬프트 생성 | 신규 |
| `internal/context/parse_features.go` | LLM 응답에서 feature JSON 파싱 | 신규 |
| `internal/context/filter_feature.go` | 2단계: GoFile 목록에서 feature 필터링 (ChonResult 제거) | 수정 |
| `internal/context/build_what_prompt.go` | ChonResult → GoFile 기반으로 변경 | 수정 |
| `internal/context/build_body_prompt.go` | ChonResult → GoFile 기반으로 변경 | 수정 |
| `internal/context/score_what.go` | ChonResult → GoFile 기반으로 변경 | 수정 |
| `internal/context/score_body.go` | ChonResult → GoFile 기반으로 변경 | 수정 |
| `internal/context/filter_by_score.go` | ChonResult → GoFile 기반으로 변경 | 수정 |
| `internal/context/format_result.go` | ChonResult → GoFile 기반으로 변경 | 수정 |
| `internal/context/run_pipeline.go` | chain 의존 제거, feature 선택 기반으로 변경 | 수정 |
| `internal/context/chon1_scores.go` | 삭제 (ChonResult 의존) | 삭제 |
| `internal/context/get_feature.go` | 삭제 (fileMap 의존 → GoFile 직접 접근) | 삭제 |
| `internal/cli/context.go` | 인자 변경 (func-name 제거, prompt만) | 수정 |

## 문서 업데이트

| 파일 | 변경 내용 |
|---|---|
| `CLAUDE.md` | context 명령 사용법 업데이트 |
| `README.md` | context 섹션 업데이트 |
| `artifacts/manual-for-ai.md` | context 사용법 업데이트 |

## 구현 순서

1. `build_feature_prompt.go` 구현 — feature 선택 프롬프트
2. `parse_features.go` 구현 — JSON 파싱
3. `select_feature.go` 구현 — LLM 호출 + feature 선택
4. context 패키지 전체를 GoFile 기반으로 리팩토링 (ChonResult 의존 제거)
   - `filter_feature.go` — `[]*model.GoFile` 입력/출력
   - `build_what_prompt.go` — GoFile에서 함수명+what 추출
   - `build_body_prompt.go` — GoFile에서 본문 추출
   - `score_what.go`, `score_body.go` — GoFile 기반
   - `filter_by_score.go` — GoFile 기반 필터링
   - `format_result.go` — GoFile 기반 출력
   - `chon1_scores.go`, `get_feature.go` 삭제
5. `run_pipeline.go` 수정 — chain 제거, feature 선택 → 필터 → what → body
6. `context.go` 수정 — cobra 인자 변경 (func-name 제거)
7. gpt-oss:20b로 정확도 + 속도 검증
8. 문서 업데이트
9. `filefunc validate` 위반 0 확인

## 완료 기준

- `filefunc context "프롬프트"` 실행 시 4단계 파이프라인 동작
- 대상 함수명 불필요
- 1단계: LLM이 codebook feature 선택 (1회)
- 2단계: feature 필터 (정적)
- 3단계: what 스코어링 rate≥0.2 (LLM 1회)
- 4단계: body 스코어링 rate≥0.5 (LLM 1회)
- 고립 함수(ParseCodebook 등)도 탐색 가능
- `--depth` 플래그로 단계 조절 가능
- gpt-oss:20b (ollama)에서 동작 확인
- filefunc 자체 코드 validate 통과
- CLAUDE.md, README.md, manual-for-ai.md 업데이트

## Phase 018과의 관계

Phase 018은 chain 기반 context. Phase 019는 LLM feature 선택 기반 context. Phase 019가 Phase 018을 대체. `internal/context/` 패키지의 what/body 스코어링 코드는 재사용.
