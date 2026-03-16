# filefunc 자기 분석 보고서

## 대상 프로젝트

| 항목 | 값 |
|---|---|
| 프로젝트 | filefunc — LLM 네이티브 Go 코드 구조 컨벤션 및 CLI 도구 |
| 언어 | Go (cobra) |
| 개발 기간 | 2026-03-15 ~ 2026-03-16 (2일, 21 커밋) |
| 특징 | 처음부터 filefunc 룰을 준수하며 개발 (대조군 없음) |
| filefunc validate | **No violations found** |

---

## 1. 구조 요약

| 지표 | 값 |
|---|---|
| Go 파일 수 | 121 (production) + 11 (test) + 24 (testdata) |
| 총 LOC | 3,037 (production, test/testdata 제외) |
| 평균 LOC/파일 | 25.1 |
| 함수 수 | 116 (exported) + 20 (unexported) = 136 |
| 타입 수 | 15 |
| 평균 func/파일 | 1.0 |
| 어노테이션 줄 수 | 325줄 (총 LOC의 10.7%) |

---

## 2. 파일 크기 분포

| 구간 | 파일 수 | 비율 |
|---|---|---|
| 1–10줄 | 9 | 7.4% |
| 11–25줄 | 65 | 53.7% |
| 26–50줄 | 40 | 33.1% |
| 51–100줄 | 7 | 5.8% |
| 101줄+ | 0 | 0.0% |

**94.2%의 파일이 50줄 이하. 101줄 초과 파일 0개.**

최대 파일: internal/cli/llmc.go (87줄).

---

## 3. 제어 흐름 분류 (Bohm-Jacopini)

| control | 파일 수 | 비율 | 의미 |
|---|---|---|---|
| sequence | 58 | 50.0% | 순차 실행 (loop/switch 없음) |
| iteration | 54 | 46.6% | 반복 (for/range) |
| selection | 4 | 3.4% | 분기 (switch) |

iteration 파일의 dimension 분포:

| dimension | 파일 수 | Q1 depth 상한 |
|---|---|---|
| 1 | 51 | 2 |
| 2 | 1 | 3 |

전부 dimension ≤ 2. filefunc 자체에는 다차원 순회가 거의 없음.

---

## 4. 어노테이션 커버리지

| 항목 | 수 |
|---|---|
| `//ff:func` | 111 |
| `//ff:type` | 14 |
| `//ff:what` | 120 |
| const/var 전용 파일 | 0 |
| **어노테이션 커버리지** | **120/121 = 99.2%** |

1개 미달: `cmd/filefunc/main.go` (cobra root 엔트리포인트).

---

## 5. 패키지별 상세

| 패키지 | 파일 수 | 최대 LOC | 역할 |
|---|---|---|---|
| cmd/filefunc | 1 | 15 | CLI 엔트리포인트 |
| internal/annotate | 3 | 57 | LLM 어노테이션 생성 |
| internal/chain | 15 | 36 | 호출 관계 추적 |
| internal/cli | 12 | 87 | cobra 명령 정의 |
| internal/llm | 13 | 39 | LLM 연동 (ollama) |
| internal/model | 4 | 18 | 데이터 구조체 |
| internal/parse | 32 | 56 | 소스 코드/어노테이션 파싱 |
| internal/report | 2 | 23 | 검증 결과 출력 |
| internal/validate | 34 | 63 | 검증 룰 구현 |
| internal/walk | 4 | 40 | 파일/디렉토리 순회 |

가장 큰 패키지: internal/validate (34파일) — 룰 1개 = 파일 1개 구조.

---

## 6. 자기검증 (dogfooding)

filefunc는 스스로 filefunc 룰을 통과해야 한다:

| 검증 | 결과 |
|---|---|
| `filefunc validate` | No violations found |
| `go build ./...` | OK |
| `go test ./...` | 32 tests PASS |
| F1/F2 (1 file 1 func/type) | 위반 0 |
| Q1 (depth ≤ control+dimension) | 위반 0 |
| A1~A16 (어노테이션 룰) | 위반 0 |

---

## 7. 비용 분석

filefunc은 처음부터 룰을 준수하며 만들어졌으므로 "전환 비용"이 아닌 "준수 비용"을 측정한다.

| 항목 | 수치 |
|---|---|
| 어노테이션 오버헤드 | 325줄 (10.7%) |
| 파일 수 (1 func 1 file로 인한) | 121개 |
| 평균 파일 크기 | 25.1줄 |

어노테이션 비율 10.7%는 whyso(7.9%)보다 높다. validate 룰 파일이 34개로 많고, 각각 짧은 함수에 2~3줄 어노테이션이 붙기 때문. 파일이 작을수록 어노테이션 비율은 올라간다.

---

## 8. 결론

filefunc은 자기 자신의 룰을 100% 준수하며 개발되었다:

1. **SRP 100% 준수** — 121개 파일 전부 1 func/type
2. **평균 25.1줄/파일** — 101줄 초과 파일 0개
3. **Q1 depth 위반 0건** — control+dimension 기반 동적 상한 적용
4. **어노테이션 커버리지 99.2%** — 120/121 파일
5. **32개 테스트 전체 통과** — mutest 16개 + exception 4개 + clean 1개 + type 3개 + 단위 8개
6. **dogfooding 완전 달성** — 도구가 자기 자신을 검증하고 통과
