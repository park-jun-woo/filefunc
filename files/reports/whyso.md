# filefunc 도입 효과 분석 — whyso

## 대상 프로젝트

| 항목 | 값 |
|---|---|
| 프로젝트 | whyso — Claude Code 세션에서 파일 변경의 "why"를 추출하는 CLI |
| 언어 | Go |
| 도입 시점 | 2026-03-16 |
| 커밋 | efb0093 (before) → a4066bb (after, Q1 수정 포함) |
| filefunc validate | **No violations found** |

---

## 1. 구조 변환 요약

| 지표 | Before | After | 변화 |
|---|---|---|---|
| Go 파일 수 | 12 | 99 | +87 (+725%) |
| 총 LOC | 1,773 | 2,417 | +644 (+36.3%) |
| 평균 LOC/파일 | 147.8 | 24.4 | -83.5% |
| 함수(func/method) 수 | 65 | 84 | +19 (depth 리팩토링으로 추출된 helper) |
| 타입(struct) 수 | 12 | 12 | 변동 없음 |
| 평균 func/파일 | 5.4 | 1.0 | -81.5% (SRP 달성) |

LOC 증가(+644줄)의 구성:
- 어노테이션 오버헤드: 192줄 (파일당 `//ff:func` + `//ff:what` 2줄)
- 파일당 `package` 선언 + `import` 블록 반복
- depth 해소를 위해 추출된 19개 helper 함수의 시그니처 + import

순수 로직 LOC 증가 없이 구조 분리에 의한 보일러플레이트 증가만 존재.

---

## 2. 파일 크기 분포

### Before (12 파일)

| 구간 | 파일 수 | 비율 |
|---|---|---|
| 1–50줄 | 0 | 0% |
| 51–100줄 | 5 | 41.7% |
| 101줄+ | 7 | 58.3% |

최대 파일: cmd/whyso/main.go (410줄, 9 func)

### After (99 파일)

| 구간 | 파일 수 | 비율 |
|---|---|---|
| 1–10줄 | 12 | 12.1% |
| 11–25줄 | 57 | 57.6% |
| 26–50줄 | 24 | 24.2% |
| 51–100줄 | 5 | 5.1% |
| 101줄+ | 1 | 1.0% |

93.9%의 파일이 50줄 이하. 101줄 초과 파일은 1개(read_yaml.go 109줄)로 Q3 WARNING 수준.

---

## 3. Single Responsibility 준수

### Before: funcs-per-file 분포

| func/파일 | 파일 수 |
|---|---|
| 2 | 1 |
| 3 | 4 |
| 4 | 3 |
| 5 | 1 |
| 9 | 1 |
| 10 | 1 |
| 15 | 1 |

12개 파일 전부 F1(1 func per file) 위반. 최악: treesitter.go에 15개 함수 혼재.

### After: funcs-per-file 분포

- func 파일: 84개 (전부 1 func)
- type 파일: 12개 (전부 1 type)
- const/var 전용: 3개

**F1/F2 위반: 0건. SRP 100% 달성.**

---

## 4. 제어 흐름 분류 (Bohm-Jacopini)

filefunc는 모든 func에 `control=` 속성을 요구하여 제어 흐름을 명시적으로 분류한다.

| control | 파일 수 | 비율 | 의미 |
|---|---|---|---|
| sequence | 41 | 48.8% | 순차 실행 (loop/switch 없음) |
| iteration | 38 | 45.2% | 반복 (for/range) |
| selection | 5 | 6.0% | 분기 (switch) |

iteration 파일은 전부 `dimension=1` (flat list 순회, depth ≤ 2).

---

## 5. Nesting Depth (Q1) 리팩토링

### Before: depth 위반 23건

초기 분리 시 발견된 18건 + `dimension=1` 도입 후 추가 발견 5건.

| depth 경로 패턴 | 건수 |
|---|---|
| for→for→if | 6 |
| for→if→if | 7 |
| for→switch→if | 2 |
| for→else→for | 1 |
| for→for→switch | 2 |
| for→for→for | 1 |
| Walk 콜백 if→if | 2 |
| closure 내 if→if→if | 1 |
| if→for→switch (출력부) | 1 |

### 해결 기법 분포

| 기법 | 적용 건수 | 새 파일 필요 |
|---|---|---|
| 조건 병합 (condition merge) | 5 | 0 |
| early continue/return | 5 | 0 |
| helper func 추출 | 13 | +13 파일 |

### After

`filefunc validate` 결과: **No violations found** (Q1 ERROR 0건).

---

## 6. 어노테이션 커버리지

| 항목 | 수 |
|---|---|
| `//ff:func` | 84 |
| `//ff:type` | 12 |
| `//ff:what` | 96 |
| 어노테이션 불필요 (const/var 전용) | 3 |
| **어노테이션 커버리지** | **96/96 = 100%** |
| 어노테이션 총 줄 수 | 192줄 (총 LOC의 7.9%) |

---

## 7. 탐색성 개선 효과

### Before

- 함수를 찾으려면 12개 파일을 열어 grep 필요
- 파일명이 도메인을 반영하지 않음 (codemap.go에 10개 함수 혼재)
- 코드 의도를 파악하려면 전체 함수 본문을 읽어야 함

### After

- **파일명 = 함수명**: `parse_go.go`에는 `parseGo` 함수만 존재
- **`//ff:what` 한 줄 요약**: 본문을 읽지 않고 함수 목적 파악 가능
- **codebook 기반 grep**: `rg '//ff:func feature=session'`으로 도메인 단위 탐색
- **control= 기반 읽기 전략**: iteration은 loop body에 집중, sequence는 필요한 step만

---

## 8. 비용-편익 분석

### 비용

| 항목 | 수치 |
|---|---|
| LOC 증가 | +644줄 (+36.3%) |
| 파일 수 증가 | +87개 (+725%) |
| 어노테이션 오버헤드 | 192줄 (7.9%) |
| 전환 작업 시간 | ~5분 (AI 에이전트 6개 병렬) |

### 편익

| 항목 | Before → After |
|---|---|
| SRP 위반 | 12/12 (100%) → 0/99 (0%) |
| depth 위반 | 23건 → 0건 |
| 평균 파일 크기 | 147.8줄 → 24.4줄 |
| 50줄 이하 파일 비율 | 0% → 93.9% |
| 어노테이션 커버리지 | 0% → 100% |
| 파일명으로 함수 식별 | 불가 → 가능 |
| filefunc validate | N/A → No violations |

---

## 9. 패키지별 상세

| 패키지 | Before 파일 | After 파일 | Before 최대 LOC | After 최대 LOC |
|---|---|---|---|---|
| cmd/whyso | 1 | 19 | 410 | 79 |
| internal/output | 2 | 9 | 139 | 109 |
| pkg/codemap | 3 | 34 | 252 | 59 |
| pkg/history | 3 | 16 | 141 | 61 |
| pkg/model | 1 | 7 | 78 | 19 |
| pkg/parser | 2 | 14 | 155 | 46 |

가장 큰 변환: pkg/codemap (3→34파일). treesitter.go의 15개 파서 함수가 각각 독립 파일로 분리.
가장 큰 depth 개선: cmd/whyso (410줄 1파일 → 19파일, 최대 79줄, depth 4→2).

---

## 10. 결론

filefunc 도입은 LOC 36.3% 증가라는 비용으로 다음을 달성:

1. **SRP 100% 준수** — 1 file = 1 func/type
2. **인지 복잡도 83.5% 감소** — 평균 파일 크기 147.8줄 → 24.4줄
3. **depth 위반 전면 해소** — 23건 → 0건 (`filefunc validate` 통과)
4. **메타데이터 완전 커버리지** — 96/96 파일 어노테이션 (100%)
5. **파일명 기반 탐색** — grep 없이 파일명만으로 함수 식별
6. **Bohm-Jacopini 분류** — 모든 함수에 control=/dimension= 명시

전환 비용(작업 시간 ~5분, LOC +36.3%)은 코드베이스 규모가 커질수록 구조적 이점에 의해 상쇄된다.
