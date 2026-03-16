# LLM 탐색 실험 보고

## 목표

사용자 프롬프트로부터 관련 코드 파일을 찾는 최적 방법을 탐색한다.

## 실험 1: reranker 방식 (실패)

### 방법

vLLM + Qwen3-Reranker (cross-encoder)로 chon=2 결과의 what 텍스트와 사용자 프롬프트의 관련도를 0.0~1.0으로 스코어링.

### 결과

| 모델 | 속도 (20건) | 정확도 |
|---|---|---|
| Qwen3-Reranker-0.6B | 0.4초 | 불량 — 관련/비관련 점수 차이 없음 |
| Qwen3-Reranker-4B | 0.6초 | 불량 — 동일 도메인 함수 전부 0.6~0.96 |

### 실패 원인

같은 패키지(validate) 안의 함수들이라 reranker가 전부 "관련 있다"고 판단. reranker는 "검색 쿼리 vs 문서"에 최적화되어 있고, "같은 도메인 안의 세밀한 구분"은 훈련 분포와 맞지 않음.

### 0.6B 점수 분포 (프롬프트: "nesting depth 검증 로직을 수정")

| 함수 | 실제 관련 | 점수 |
|---|---|---|
| HasAnyChecked (X) | X | 0.77 |
| CheckControlRequired (X) | △ | 0.17 |
| CheckDimensionRequired (O) | O | 0.17 |
| 나머지 17개 | 대부분 X | < 0.1 |

### 4B 점수 분포 (동일 프롬프트)

| 함수 | 실제 관련 | 점수 |
|---|---|---|
| CheckDimensionValue (O) | O | 0.96 |
| HasAnyChecked (X) | X | 0.96 |
| CheckCodebookValues (X) | X | 0.94 |
| CheckControlRequired (△) | △ | 0.93 |
| CheckAnnotationRequired (X) | X | 0.90 |
| CheckDimensionRequired (O) | O | 0.83 |

4B가 관련 있는 것에 높은 점수를 주긴 하지만, 관련 없는 것도 높게 줘서 **구분이 안 됨**.

---

## 실험 2: codebook 기반 grep 탐색어 생성 (성공)

### 방법

LLM(gpt-oss:20b, ollama)에게 codebook.yaml과 사용자 질문을 주고, grep용 `feature=X type=Y` 탐색어를 JSON으로 생성하게 함.

### 프롬프트

```json
{
  "task": "Given codebook, output grep search terms to find files related to the question. Use only codebook values.",
  "codebook": {"feature": [...], "type": [...]},
  "question": "사용자 질문",
  "format": "JSON array of {feature,type}, max 3"
}
```

### 결과

| 질문 | LLM 출력 | 정답 | 평가 |
|---|---|---|---|
| nesting depth 검증 수정 | `validate/rule`, `parse/parser`, `chain/walker` | validate/rule 핵심 | O |
| codebook feature 추가 | `codebook/model`, `annotate/util`, `cli/command` | codebook/model 핵심 | O |
| chain 출력 포맷 변경 | `chain/formatter`, `chain/command`, `chain/parser` | chain/formatter 핵심 | O |

**세 질문 모두 첫 번째 결과가 정답.**

### 장점

- LLM 호출 단 1회 (탐색어 생성)
- 이후는 전부 정적 분석 (grep → chain → read)
- vLLM 불필요 — ollama로 충분 (텍스트 생성 태스크)
- codebook 품질이 좋으면 정확도가 높음
- JSON 출력으로 파싱 안정적

---

## 결론: 최적 탐색 파이프라인

```
1. LLM 1회 호출  → codebook 기반 grep 탐색어 생성 (feature=X type=Y)
2. rg grep       → 파일 목록 확보
3. chain --meta  → 호출 관계 + what 텍스트
4. read          → 필요한 파일만 본문 읽기
```

- 1번만 LLM 의존 (ollama, gpt-oss:20b)
- 2~4번은 정적 분석 (AST + grep)
- reranker(건당 스코어링)보다 정확하고 빠름
- codebook이 탐색의 어휘 — codebook 품질 = 탐색 정확도

### reranker vs codebook grep 비교

| | reranker | codebook grep |
|---|---|---|
| LLM 호출 | 건당 1회 (20~100회) | 총 1회 |
| 인프라 | vLLM + GPU 필수 | ollama로 충분 |
| 정확도 | 같은 도메인 구분 불가 | codebook 어휘로 정확 구분 |
| 속도 | 0.4~1초 (웜업 후) | LLM 1회 + grep 즉시 |
| 의존성 | Qwen3-Reranker 모델 | codebook.yaml 품질 |
