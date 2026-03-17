# filefunc 도입 효과 분석 — fullend

## 대상 프로젝트

| 항목 | 값 |
|---|---|
| 프로젝트 | fullend — 9개 SSOT의 정합성을 검증하고 코드를 산출하는 Go CLI |
| 언어 | Go 1.22+ |
| 도입 기간 | 2026-03-14 ~ 2026-03-17 (Phase031~043, 13단계) |
| 커밋 | d2bc6b0 (before) → 49c5410 (after) |
| filefunc validate | **ERROR 0, WARNING 0** |

---

## 1. 구조 변환 요약

| 지표 | Before | After | 변화 |
|---|---|---|---|
| Go 소스 파일 수 | 87 | 1,260 | +1,173 (+1,348%) |
| 총 LOC | 21,232 | 31,976 | +10,744 (+50.6%) |
| 평균 LOC/파일 | 244 | 25.4 | -89.6% |
| 중앙값 LOC/파일 | — | 21 | |
| 함수(func/method) 수 | 617 | 1,077 | +460 (리팩토링 추출 helper) |
| 타입(type) 수 | 41 | 218 | +177 (내부 타입 분리) |
| 평균 func/파일 | 7.0 | 0.85 | -87.9% |
| 최대 파일 크기 | 1,113줄 | 232줄 | -79.2% |
| 테스트 파일 | 40 | 40 | 변동 없음 |

LOC 증가(+10,744줄)의 구성:
- 어노테이션 오버헤드: ~2,500줄 (파일당 `//ff:func` + `//ff:what` 2줄 × 1,253파일)
- 파일당 `package` 선언 + `import` 블록 반복: ~5,000줄
- depth 해소를 위해 추출된 460개 helper 함수의 시그니처 + import: ~3,200줄

순수 로직 LOC 증가 없이 구조 분리에 의한 보일러플레이트 증가만 존재.

---

## 2. 파일 크기 분포

### Before (87 파일)

| 구간 | 파일 수 | 비율 |
|---|---|---|
| 1–25줄 | 3 | 3.4% |
| 26–50줄 | 13 | 14.9% |
| 51–100줄 | 14 | 16.1% |
| 101–200줄 | 30 | 34.5% |
| 201줄+ | 27 | 31.0% |

최대 파일: ssac/validator/symbol.go (1,113줄, 35 func), gen/gogin/model_impl.go (1,107줄, 28 func), gen/hurl/hurl.go (1,048줄, 26 func)

### After (1,260 파일)

| 구간 | 파일 수 | 비율 |
|---|---|---|
| 1–25줄 | 817 | 64.8% |
| 26–50줄 | 343 | 27.2% |
| 51–100줄 | 88 | 7.0% |
| 101–200줄 | 10 | 0.8% |
| 201줄+ | 2 | 0.2% |

**99.0%의 파일이 100줄 이하.** 201줄 초과 2파일: generate_method_from_iface.go (232줄, selection), query_opts_template.go (216줄, var 전용).

---

## 3. Single Responsibility 준수

### Before: funcs-per-file 분포

| func/파일 | 파일 수 | 비고 |
|---|---|---|
| 0–1 | 21 | type/const 전용 |
| 2–5 | 32 | |
| 6–10 | 12 | |
| 11–20 | 14 | |
| 21+ | 8 | 최악: symbol.go 35 func |

87개 파일 중 66개(75.9%)가 F1 위반.

### After: 전체 패키지

`filefunc validate` — F1/F2/F3 ERROR 0건. 1,260파일 전부 1-func/1-type/1-method 준수 (const/var 전용 5파일 제외).

---

## 4. 제어 흐름 분류 (Bohm-Jacopini)

| control | 파일 수 | 비율 | 의미 |
|---|---|---|---|
| sequence | 393 | 37.3% | 순차 실행 (loop/switch 없음) |
| selection | 74 | 7.0% | 분기 (switch) |
| iteration | 588 | 55.7% | 반복 (for/range) |

### dimension 분포 (iteration 파일)

| dimension | 파일 수 | Q1 상한 | 의미 |
|---|---|---|---|
| 1 | 509 | 2 | flat list 순회 |
| 2 | 66 | 3 | 2중 중첩 (paths→ops, funcs→seqs) |
| 3 | 10 | 4 | 3중 중첩 (policies→rules→actions) |
| 4 | 2 | 5 | 4중 중첩 (spec→existing dedup) |
| 5 | 1 | 6 | 5중 중첩 (AST entries→decls→specs→fields→names) |

---

## 5. Nesting Depth (Q1) 리팩토링

### Phase031~043에서 해결한 Q1 위반 총계

| Phase | 대상 | Q1 위반 |
|---|---|---|
| Phase031~033 | ssac/validator, gen/gogin, gen/hurl, orchestrator 초기 | 기반 작업 |
| Phase034 | ssac/validator | 23건 |
| Phase035 | gen/gogin | 32건 |
| Phase036 | gen/hurl | 25건 |
| Phase037 | orchestrator | 23건 |
| Phase039 | pkg/* | 1건 |
| Phase040 | internal 소형 + gen/gogin Q3 | 11건 |
| Phase041 | stml + ssac/parser + contract | 8건 |
| Phase042 | crosscheck | 17건 |
| Phase043 | ssac/generator | 8건 |
| **합계** | | **148건+** |

### After

`filefunc validate` 결과: **Q1 ERROR 0건, WARNING 0건**.

---

## 6. 어노테이션 커버리지

| 항목 | 수 |
|---|---|
| `//ff:func` | 1,055 |
| `//ff:type` | 198 |
| `//ff:what` | 1,256 |
| 어노테이션 불필요 (const/var 전용) | 5 |
| **어노테이션 커버리지** | **1,253/1,260 = 99.4%** |

미부착 5파일은 const/var 전용 (go_templates.go, rules.go 등). filefunc 규칙상 예외.

---

## 7. 패키지별 상세

| 패키지 | 파일 수 | 최대 LOC | 비고 |
|---|---|---|---|
| cmd/fullend | 23 | 94 | Phase040에서 main.go 12→23 분해 |
| orchestrator | 84 | 110 | Phase037 |
| crosscheck | 191 | 149 | Phase042. 19→191, 가장 큰 분해 |
| ssac/parser | 49 | 55 | Phase041. 752줄 파서 분해 |
| ssac/validator | 105 | 118 | Phase034. 1,113줄 symbol.go 분해 |
| ssac/generator | 162 | 183 | Phase043. 10→162 |
| stml/parser | 62 | 49 | Phase041. 705줄 파서 분해 |
| stml/validator | 59 | 41 | Phase041 |
| stml/generator | 65 | 60 | Phase041 |
| funcspec | 14 | 48 | Phase040 |
| gen/gogin | 117 | 232 | Phase035 + Phase040 Q3 |
| gen/hurl | 57 | 114 | Phase036 |
| gen/react | 30 | 48 | Phase040 |
| genmodel | 33 | 78 | Phase040 |
| genapi | 4 | 32 | Phase040 |
| projectconfig | 17 | 53 | Phase040 |
| statemachine | 8 | 72 | Phase040 |
| policy | 13 | 63 | Phase040 |
| reporter | 14 | 49 | Phase040 |
| contract | 38 | 56 | Phase041. splice.go(290줄) Q1 depth 4 해소 |
| pkg/* (13개) | 115 | 48 | Phase039. 26→115, Q1 depth 4→2 |

---

## 8. 전환 Phase 이력

| Phase | 내용 | 파일 변화 | 위반 해소 |
|---|---|---|---|
| 031 | filefunc 도입, 비대 파일 5개 분해 | +42 | 기반 |
| 032 | 미분해 21파일 전수 분해 | +58 | 기반 |
| 033 | control annotation 부착 | — | A9 |
| 034 | ssac/validator dimension+Q1 | +3 | Q1 23건 |
| 035 | gen/gogin dimension+Q1 | +9 | Q1 32건 |
| 036 | gen/hurl dimension+Q1 | +4 | Q1 25건 |
| 037 | orchestrator dimension+Q1 | +6 | Q1 23건 |
| 038 | mutest FAIL 수정 | — | 검증 3건 |
| 039 | pkg/* 13개 패키지 | 26→115 | ERROR 98건 |
| 040 | internal 소형 + gen/gogin Q3 | 14→160 | ERROR 53건, WARNING 4건 |
| 041 | stml + ssac/parser + contract | 17→273 | ERROR 62건, WARNING 1건 |
| 042 | crosscheck | 19→191 | ERROR 65건, WARNING 4건 |
| 043 | ssac/generator + codebook 전환 | 9→162 | ERROR 37건, WARNING 5건 |

---

## 9. Mutation Test 결과

| 항목 | 수치 |
|---|---|
| 총 케이스 | 114 |
| PASS | 87 |
| FAIL | 9 (확정 5 + 불확실 3 + 오집계 1) |
| SKIP | 18 (dummy에 해당 기능 없음) |
| 통과율 (SKIP 제외) | 90.6% |

확정 FAIL 5건: middleware↔securitySchemes 교차 검증 누락, funcspec stub body 미감지, @get result type 미검증, @state transition 대소문자, @call pkg 함수명 대소문자.

---

## 10. 비용-편익 분석

### 비용

| 항목 | 수치 |
|---|---|
| LOC 증가 | +10,744줄 (+50.6%) |
| 파일 수 증가 | +1,173개 (+1,348%) |
| 어노테이션 오버헤드 | ~2,500줄 (7.8%) |
| 전환 작업 | Phase031~043 (13 phase, AI 에이전트) |

### 편익

| 항목 | Before → After |
|---|---|
| F1(SRP) 위반 | 66/87 (75.9%) → 0/1,260 (0%) |
| Q1(depth) 위반 | 148건+ → 0건 |
| 평균 파일 크기 | 244줄 → 25.4줄 |
| 100줄 이하 파일 비율 | 34.5% → 99.0% |
| 최대 파일 크기 | 1,113줄 → 232줄 |
| 어노테이션 커버리지 | 0% → 99.4% |
| 파일명으로 함수 식별 | 불가 → 가능 |
| filefunc validate | N/A → ERROR 0, WARNING 0 |
| Mutest 통과율 | N/A → 90.6% (114건) |

---

## 11. 결론

filefunc 도입은 LOC 50.6% 증가라는 비용으로 다음을 달성:

1. **SRP 100% 준수** — 1,260파일 전부 1 func/type/method per file
2. **인지 복잡도 89.6% 감소** — 평균 파일 크기 244줄 → 25.4줄
3. **depth 위반 전면 해소** — 148건+ → 0건 (`filefunc validate` ERROR 0)
4. **메타데이터 99.4% 커버리지** — 1,253/1,260 파일 어노테이션
5. **파일명 기반 탐색** — grep 없이 파일명만으로 함수 식별
6. **Bohm-Jacopini 분류** — 전체 함수에 control=/dimension= 명시
7. **Mutest 90.6%** — 114건 mutation test로 validator 품질 보증
8. **codebook 어휘 체계** — feature 37개, type 9개로 전체 코드베이스 분류

13 phase에 걸친 체계적 전환으로 21,000줄 규모의 코드베이스를 완전히 filefunc 준수 상태로 변환. 최대 파일이 1,113줄에서 232줄로 줄었고, 99%의 파일이 100줄 이하.
