# Phase 017: chain --prompt + --rate (reranker 관련도 필터링) ✅ 완료

## 목표

`filefunc chain` 출력의 chon=2 결과를 reranker 모델로 사용자 프롬프트와의 관련도를 평가하고, 임계값 이상만 필터링하여 정밀한 컨텍스트를 제공한다.

## 배경

chain의 현재 문제: chon=2 이상에서 co-called(형제) 등이 폭발하여 노이즈가 많다. 구조 기반 검색(chain)으로 후보를 확장한 뒤, reranker 리랭킹으로 정밀도를 올리는 RAG 방식.

```
1. chain chon=1     → 구조적으로 확실한 관련 함수 (항상 포함)
2. chain chon=2     → 후보군 확장
3. reranker 스코어링 → what 텍스트와 사용자 프롬프트의 관련도 0.0~1.0
4. --rate 필터      → 임계값 이상만 출력
```

벡터 DB 없이 호출 관계 + what 텍스트 매칭으로 동작. 후보가 20~100건.

## 기술 조사 결과

### ollama는 reranker에 부적합

Qwen3-Reranker(cross-encoder)는 "yes"/"no" 토큰의 logit 비율로 연속 점수를 산출한다. ollama API는 logit을 노출하지 않으므로 이진 분류(1.0 or 0.0)만 가능. 연속 점수 필터링에 부적합.

### vLLM이 reranker를 네이티브 지원

vLLM(Apache 2.0, UC Berkeley 출신, a16z 후원)은 `/v1/score` 엔드포인트로 cross-encoder의 logit 기반 연속 점수를 HTTP API로 제공한다.

```bash
# vLLM 서버 실행
vllm serve Qwen/Qwen3-Reranker-0.6B --task score \
  --hf_overrides '{"architectures":["Qwen3ForSequenceClassification"],"classifier_from_token":["no","yes"],"is_original_qwen3_reranker":true}'
```

```
POST http://localhost:8000/v1/score
{"model":"Qwen/Qwen3-Reranker-0.6B", "text_1":"사용자 프롬프트", "text_2":"함수명: what 텍스트"}
→ {
    "id": "score-xxx",
    "object": "list",
    "model": "Qwen/Qwen3-Reranker-0.6B",
    "data": [{"index": 0, "object": "score", "score": 0.95}],
    "usage": {"prompt_tokens": 42, "total_tokens": 42}
  }
```

점수 추출: `response.data[0].score`

Go에서는 `net/http`로 HTTP POST만 하면 됨. 기존 `filefunc llmc`의 ollama 호출 패턴과 동일한 구조.

### 모델 선택

| 모델 | 파라미터 | 타입 | 점수 방식 |
|---|---|---|---|
| **Qwen3-Reranker-0.6B** | 0.6B | cross-encoder 전용 | logit 기반 연속 점수 (정확) |
| qwen3:0.6b | 0.6B | 범용 LLM | 프롬프트로 텍스트 생성 (부정확) |

reranker 전용 모델이 동일 크기 범용 LLM보다 스코어링 정확도가 높고, 프롬프트 엔지니어링이 불필요하다.

## 설계

### 플래그

```bash
filefunc chain func RunAll --chon 2 --meta what \
  --prompt "nesting depth 검증 로직을 수정하려고 한다" \
  --rate 0.8
```

| 플래그 | 설명 | 기본값 |
|---|---|---|
| `--prompt` | 사용자 작업 의도 (관련도 평가 기준) | — |
| `--rate` | 관련도 임계값 (0.0~1.0). 이상만 출력 | 0.8 |
| `--model` | reranker 모델명 | `Qwen/Qwen3-Reranker-0.6B` |
| `--score-endpoint` | vLLM 엔드포인트 | `http://localhost:8000` |

- `--prompt` 없이 `--rate`만 지정하면 에러
- `--prompt`만 지정하면 `--rate 0.8` 기본 적용
- chon=1 결과는 항상 포함 (필터링 대상 아님). chon=2+ 결과만 reranker 평가

### 스코어링 로직

각 chon=2 함수에 대해 vLLM `/v1/score` 호출:

```json
// 요청
{"model":"Qwen/Qwen3-Reranker-0.6B",
 "text_1":"nesting depth 검증 로직을 수정하려고 한다",
 "text_2":"CheckNestingDepth: Q1: control과 dimension 기반으로 nesting depth 상한 검증"}

// 응답
{"id":"score-xxx","object":"list","model":"Qwen/Qwen3-Reranker-0.6B",
 "data":[{"index":0,"object":"score","score":0.95}],
 "usage":{"prompt_tokens":42,"total_tokens":42}}
```

점수 추출: `response.data[0].score` → 0.95

- cross-encoder가 query-document 쌍의 관련도를 logit 기반으로 평가
- 프롬프트 엔지니어링 불필요
- 건당 수십 ms (0.6B, 분류 태스크). 20~100건 → 1~3초
- Go에서 goroutine 병렬 호출로 추가 최적화 가능

### 출력 포맷

기존 chain 출력과 동일하되, 필터링된 결과만 표시. 점수를 접미사로:

```
CheckNestingDepth (what="Q1: control과 dimension 기반으로 nesting depth 상한 검증")
  1촌 calls: depthLimit (what="control과 dimension으로 Q1 depth 상한을 계산")
  1촌 called-by: RunAll (what="모든 검증 룰을 실행하고 위반 목록을 반환")
  2촌 co-called: CheckDimensionValue (what="A16: dimension= 값은 양의 정수여야 함") [0.85]
  2촌 co-called: CheckDimensionRequired (what="A15: control=iteration이면 dimension= 필수") [0.82]
  -- 18 results filtered by rate>0.8 (3 shown) --
```

- chon=1은 점수 없이 항상 표시
- chon=2+는 점수 표시 + 임계값 미만 제거
- 마지막 줄에 필터링 요약

## 산출물

| 파일 | 내용 | 구분 |
|---|---|---|
| `internal/chain/score_relevance.go` | vLLM `/v1/score` 호출 + 관련도 스코어링 | 신규 |
| `internal/chain/build_score_input.go` | 함수명 + what → document 텍스트 생성 | 신규 |
| `internal/chain/filter_by_rate.go` | 임계값 이상 결과만 필터링 | 신규 |
| `internal/chain/format_chain.go` | 점수 접미사 + 필터링 요약 출력 | 수정 |
| `internal/cli/chain_func.go` | `--prompt`, `--rate`, `--model`, `--score-endpoint` 플래그 추가 | 수정 |
| `internal/cli/chain_feature.go` | 동일 | 수정 |

## 문서 업데이트

| 파일 | 변경 내용 |
|---|---|
| `CLAUDE.md` | chain 명령에 --prompt, --rate 플래그 설명 추가 |
| `README.md` | chain 섹션에 reranker 필터링 사용법 + vLLM 서버 실행 방법 추가 |
| `artifacts/manual-for-ai.md` | Commands 섹션에 --prompt, --rate 사용 예시 추가 |

## 구현 순서

1. `build_score_input.go` 구현 — 함수명 + what → document 텍스트 생성
2. `score_relevance.go` 구현 — vLLM `/v1/score` HTTP POST + 점수 수집
3. `filter_by_rate.go` 구현 — 임계값 필터링
4. `format_chain.go` 수정 — 점수 접미사 + 필터링 요약
5. `chain_func.go` 수정 — 플래그 추가 + 스코어링 파이프라인 연결
6. `chain_feature.go` 수정 — 동일
7. Qwen3-Reranker-0.6B 정확도 + 속도 검증
8. 문서 업데이트
9. `filefunc validate` 위반 0 확인

## 완료 기준

- `--prompt` + `--rate` 지정 시 chon=2 결과가 reranker 관련도로 필터링됨
- chon=1은 항상 포함 (필터링 대상 아님)
- vLLM `/v1/score` 엔드포인트로 건당 스코어링 (20~100건 → 1~3초 이내)
- `--prompt` 없이 `--rate`만 지정하면 에러
- Qwen3-Reranker-0.6B에서 동작 확인
- filefunc 자체 코드 validate 통과
- CLAUDE.md, README.md, manual-for-ai.md 업데이트

## 사전 요구사항

vLLM 서버가 Qwen3-Reranker-0.6B를 서빙 중이어야 한다:

```bash
pip install vllm
vllm serve Qwen/Qwen3-Reranker-0.6B --task score \
  --hf_overrides '{"architectures":["Qwen3ForSequenceClassification"],"classifier_from_token":["no","yes"],"is_original_qwen3_reranker":true}'
```

GPU 필요 (CUDA). CPU 전용 환경에서는 `--device cpu` 옵션으로 동작하나 속도 저하.
