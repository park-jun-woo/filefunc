# Phase 018: filefunc context — 4단계 컨텍스트 파이프라인

## 목표

사용자 프롬프트로부터 관련 코드를 최소 토큰으로 정밀하게 찾는 4단계 파이프라인을 구현한다.

```
1단계: chain chon=2          → 구조적 후보 확보 (AST, 정적)
2단계: same-feature 필터     → chon=2 중 같은 feature만 유지, chon=1은 무조건 유지 (정적)
3단계: what 스코어링 rate≥0.2 → LLM 1회 배치
4단계: body 스코어링 rate≥0.5 → LLM 1회 배치
```

LLM 호출: 2회만. 1~2단계는 순수 정적 분석.

## 배경

### 실험 결과 (files/LLM-탐색-실험.md)

- reranker(vLLM + Qwen3-Reranker): 같은 도메인 구분 불가. 실패.
- codebook grep LLM 탐색어 생성: feature 단독은 정확, type 조합은 누락 발생.
- gpt-oss:20b what 스코어링: 관련/비관련 구분 가능 (0.1 vs 0.3~1.0).
- gpt-oss:20b 본문 스코어링: what보다 정밀 (0.1 vs 0.5~0.9).

### 핵심 발견

- feature 하나로 전체 파일의 70~99% 제거 가능 (filefunc 기준 최대 validate 34/121=28%)
- chain chon=1은 구조적으로 확실한 관련 → 무조건 유지
- chain chon=2는 노이즈 포함 → same-feature 필터 + LLM 스코어링으로 정밀화

### SILK 방법론 대응

```
SILK 구조적 쿼리 (80%) → 1단계 chain + 2단계 feature 필터
SILK 의미적 쿼리 (15%) → 3단계 what + 4단계 본문 스코어링
SILK VALID 검증        → filefunc validate
```

## 설계

### 사용법

```bash
filefunc context <func-name> "프롬프트"

filefunc context CheckNestingDepth "nesting depth 검증 로직을 수정하려고 한다"
```

대상 함수명을 지정하면 해당 함수의 chon=2 체인에서 파이프라인 시작.
대상 함수를 모르면 `whyso map`으로 키워드 검색하여 진입점 확보.

### 파이프라인 상세

```
입력: 대상 함수명 + 사용자 프롬프트

1단계 — chain chon=2 (정적)
  대상 함수의 호출 관계를 chon=2까지 탐색
  → 후보 20~40개

2단계 — same-feature 필터 (정적)
  대상 함수의 feature 값 확인 (예: feature=validate)
  chon=1 결과: 무조건 유지 (점수 1.0)
  chon=2 결과: 같은 feature인 것만 유지, 나머지 제거
  → 후보 10~20개

3단계 — what 스코어링 (LLM 1회 배치)
  2단계 통과 함수들의 함수명 + what을 배치 프롬프트로 전달
  chon=1: 스코어링 생략 (점수 1.0 유지)
  chon=2: 각 함수 0.0~1.0 점수
  → rate≥0.2 필터 (후보 4~8개)

4단계 — 본문 스코어링 (LLM 1회 배치)
  3단계 통과 함수들의 본문을 배치 프롬프트로 전달 (parse.ExtractFuncSource 사용)
  chon=1: 스코어링 생략 (점수 1.0 유지)
  chon=2: 각 함수 0.0~1.0 점수
  → rate≥0.5 필터 (최종 2~5개)

출력: 최종 함수 목록 + 점수 + 촌수
```

### 플래그

```bash
filefunc context <func> "프롬프트"                    # 기본 (4단계 전부)
filefunc context <func> "프롬프트" --depth 1           # 1단계만 (chain)
filefunc context <func> "프롬프트" --depth 2           # 2단계까지 (feature 필터)
filefunc context <func> "프롬프트" --depth 3           # 3단계까지 (what 스코어링)
filefunc context <func> "프롬프트" --depth 4           # 4단계까지 (본문 스코어링, 기본값)
filefunc context <func> "프롬프트" --what-rate 0.2     # 3단계 임계값 (기본 0.2)
filefunc context <func> "프롬프트" --body-rate 0.5     # 4단계 임계값 (기본 0.5)
filefunc context <func> "프롬프트" --model gpt-oss:20b # ollama 모델
filefunc context <func> "프롬프트" --endpoint http://... # ollama 엔드포인트
filefunc context <func> "프롬프트" --meta what         # 출력에 meta 포함
```

### 출력 포맷

```
[1/4] chain chon=2: CheckNestingDepth → 22 funcs
[2/4] feature filter (validate): 22 → 14
[3/4] what scoring: 14 → 5 (rate≥0.2)
[4/4] body scoring: 5 → 3 (rate≥0.5)

Results:
  CheckNestingDepth [1.00] 1촌 (what="Q1: control과 dimension 기반으로 nesting depth 상한 검증")
  depthLimit [1.00] 1촌 (what="control과 dimension으로 Q1 depth 상한을 계산")
  CheckDimensionValue [0.50] 2촌 (what="dimension= 값은 양의 정수여야 함")
```

### LLM 프롬프트

#### 3단계: what 스코어링

```
사용자가 수정하려는 작업과 각 함수의 관련도를 평가하시오.
관련도: 0.0(무관) ~ 0.5(간접 관련) ~ 1.0(직접 관련)
직접 수정 대상이면 0.8 이상, 영향을 받는 함수면 0.4~0.7, 무관하면 0.2 이하.

작업: "사용자 프롬프트"

1. FuncName: "what 텍스트"
2. FuncName: "what 텍스트"
...

각 번호에 대해 점수만. 형식: 번호. 점수
```

#### 4단계: 본문 스코어링

```
(동일 가이드라인)

작업: "사용자 프롬프트"

1. func FuncName(args) ReturnType {
    본문
}
...

각 번호에 대해 점수만. 형식: 번호. 점수
```

## 산출물

| 파일 | 내용 | 구분 |
|---|---|---|
| `internal/context/run_pipeline.go` | 4단계 파이프라인 오케스트레이터 | 신규 |
| `internal/context/filter_feature.go` | 2단계: same-feature 필터 | 신규 |
| `internal/context/score_what.go` | 3단계: what 배치 스코어링 | 신규 |
| `internal/context/score_body.go` | 4단계: 본문 배치 스코어링 (parse.ExtractFuncSource 사용) | 신규 |
| `internal/context/build_what_prompt.go` | what 스코어링용 프롬프트 생성 | 신규 |
| `internal/context/build_body_prompt.go` | 본문 스코어링용 프롬프트 생성 | 신규 |
| `internal/context/parse_scores.go` | LLM 응답에서 점수 파싱 | 신규 |
| `internal/context/format_result.go` | 결과 출력 | 신규 |
| `internal/cli/context.go` | cobra 명령 정의 | 신규 |

## 문서 업데이트

| 파일 | 변경 내용 |
|---|---|
| `CLAUDE.md` | context 명령 추가, 작업 흐름에 context 단계 추가 |
| `README.md` | context 섹션 추가 |
| `artifacts/manual-for-ai.md` | context 사용법 추가 |
| `codebook.yaml` | feature에 `context` 추가 |

## 구현 순서

1. codebook.yaml에 `context` feature 추가
2. `internal/context/` 패키지 생성
3. `filter_feature.go` 구현 — same-feature 필터
4. `build_what_prompt.go` + `score_what.go` 구현 — what 배치 스코어링
5. `build_body_prompt.go` + `score_body.go` 구현 — 본문 배치 스코어링
6. `parse_scores.go` 구현 — 응답 파싱
7. `run_pipeline.go` 구현 — 4단계 오케스트레이터
8. `format_result.go` 구현 — 출력
9. `internal/cli/context.go` 구현 — cobra 명령
10. gpt-oss:20b로 정확도 + 속도 검증
11. 문서 업데이트
12. `filefunc validate` 위반 0 확인

## 완료 기준

- `filefunc context <func> "프롬프트"` 실행 시 4단계 파이프라인 동작
- 1단계: chain chon=2 구조적 후보 (정적)
- 2단계: same-feature 필터, chon=1 무조건 유지 (정적)
- 3단계: what 스코어링 rate≥0.2 (LLM 1회)
- 4단계: 본문 스코어링 rate≥0.5 (LLM 1회)
- `--depth` 플래그로 단계 조절 가능
- gpt-oss:20b (ollama)에서 동작 확인
- filefunc 자체 코드 validate 통과
- CLAUDE.md, README.md, manual-for-ai.md 업데이트

## Phase 017과의 관계

Phase 017(chain --prompt/--rate)은 vLLM reranker 방식. reranker가 같은 도메인 구분에 부적합함이 실험으로 확인됨. Phase 018은 chain + feature 필터 + 범용 LLM 스코어링으로 대체. Phase 017 코드는 유지하되, 실제 사용은 `filefunc context`를 권장.
