# filefunc 평가와 설계 결정

## 강점

**1. 문제 정의가 정확하다**

"불필요한 290개를 안 여는 게 더 중요하다" — LLM 컨텍스트 오염 문제를 정확히 짚었다. 학술 근거(Lost in the Middle, Context Rot)도 주장을 뒷받침한다.

**2. 설계가 일관적이다**

"1 file 1 concept" 원칙 하나에서 파일 구조 룰, 어노테이션, 코드북, chain이 모두 파생된다. 임의 룰의 집합이 아니라 하나의 원칙에서 나온 체계.

**3. 기존 도구 없이도 즉시 적용 가능하다**

CLI 도구가 없어도 컨벤션 자체만으로 효과가 있다. `//ff:` 주석과 파일 분리는 오늘 당장 손으로 할 수 있다. 도구는 자동화일 뿐.

**4. whyso 연동이 자연스럽다**

func = file이면 git log = 함수 변경 이력. 별도 분석 없이 파일 이력이 곧 함수 이력. 암묵적 커플링 검출도 파일 동시 변경 통계만으로 가능.

---

## 설계 결정: 제1시민은 AI 에이전트다

filefunc의 제1시민은 사람이 아니라 AI 에이전트(Claude Code)다.

Claude Code는 `ls`가 아니라 `grep`으로 탐색한다. 파일 500개든 1000개든 `rg '//ff:func feature=validate'` 한 번이면 끝. 파일이 많을수록 각 파일이 작고, read 한 번에 딸려오는 노이즈가 줄어들어 오히려 유리하다.

메서드 20개 = 파일 20개도 문제없다. 사람은 불편하게 느끼지만, Claude Code는 필요한 메서드 파일 3개만 정확히 read한다. 나머지 17개를 안 여는 게 핵심이다.

### 파일 수 폭발은 약점이 아니다

| 우려 | 판정 | 이유 |
|---|---|---|
| 파일 수 폭발 | 약점 아님 | Claude Code에게는 장점 |
| 메서드 파일 분리 응집도 | 약점 아님 | grep으로 리시버 타입 기준 검색 |
| PR 리뷰 부담 | 2시민 문제 | AI 코드리뷰로 해소 가능 |
| IDE 탐색 불편 | 2시민 문제 | VSCode 확장으로 해결 |

### 사람의 불편은 뷰 레이어에서 해결한다

파일은 Claude Code에 최적화된 물리 구조로 유지하고, 사람은 VSCode 확장에서 타입별 메서드 그룹핑, feature별 func 목록 등 원하는 뷰로 모아본다. filefunc의 구조를 사람에게 맞춰 타협하지 않는다.

---

## 설계 결정: //ff:what과 //ff:why

기존 `//ff:desc`와 `//ff:ssot`를 `//ff:what`과 `//ff:why`로 대체한다. 직관적이고 기억하기 쉽다.

```go
//ff:func feature=validate type=rule
//ff:what 파일당 func 개수를 검증한다
//ff:why 제1시민은 AI 에이전트. 파일 수 폭발은 약점이 아니라 장점이다
//ff:calls check_file_funcs
//ff:uses FileResult, Violation
func CheckOneFileOneFunc(...) ...
```

- `//ff:what` — 필수. 이 함수가 뭘 하는가. **소형 LLM이 func body와 대조하여 일치 여부를 검증하는 기준**
- `//ff:why` — 선택. 왜 이렇게 만들었는가. 반드시 사용자의 요구/결정이 근거여야 한다. AI의 판단이나 추측은 why가 아니다. 검증 불가, 이력으로 남기는 것

---

## 남은 약점 (2개)

1. **어노테이션 drift 방지** — `//ff:checked` 서명으로 해결 확정
   - `//ff:calls`, `//ff:uses`: 코드와 기계적 대조 가능
   - `//ff:what`: 소형 LLM이 func body와 대조 → 일치 시 `//ff:checked llm=모델명 hash=body해시` 자동 기록. body 수정 시 해시 불일치로 자동 무효화 → ERROR. `--no-llm`으로 비활성화 가능 (환경 제약 대응)
2. **코드북 설계 품질** — 코드북이 나쁘면 grep이 부정확. 의도된 트레이드오프. 코드북으로 어휘를 정규화하면 빠진 feature, 중복된 type, 애매한 분류가 목록에서 드러난다. 구멍이 보여야 관리가 된다. 대응 전략: 도메인을 날카롭게 좁히고, 프로젝트마다 코드북을 맞춤 작성한다

---

## 해결된 미결 사항

| 미결 사항 | 결론 |
|---|---|
| Go 멀티리턴 output 표현 | 미결 아님. 어노테이션은 Go 코드 위에 붙고, func signature 자체가 input/output을 명시하므로 별도 표현 불필요 |
| cyclomatic complexity 상한 | nesting depth 2 룰(Q1)로 충분. "반복 1depth, 분기 1depth"로 세분화하는 안도 검토했으나, 이중 반복 등 실무 패턴을 과도하게 제한함. 룰은 단순할수록 지켜진다. "depth 2" 하나로 확정 |
| 어노테이션 desc/ssot 네이밍 | `//ff:desc` → `//ff:what`, `//ff:ssot` → `//ff:why`로 확정. what은 필수(LLM 검증 대상), why는 선택(사용자 결정 이력) |
| F4 init() 룰 개정 | 기존: `init()`은 `main.go`에만 허용. Phase 001에서 cobra 패턴과 충돌 확인. 개정: 각 파일은 최대 1개의 `init()`을 선택적으로 가질 수 있다. `init()`만 단독으로 존재하는 파일은 불허. 1 file 1 func이므로 init()과 func은 자연스럽게 1:1 |
| A1/A3 어노테이션 적용 범위 개정 | `//ff:type` 추가. func 파일은 `//ff:func`, type 파일은 `//ff:type` 필수. `//ff:what`은 func 또는 type이 있는 파일에 필수. const-only 파일은 대상 아님 |

## 미결 사항

| 미결 사항 | 우선순위 | 이유 |
|---|---|---|
| 파라미터 개수 상한 | 중간 | 유용하지만 없어도 동작함 |
| 언어별 구조 강제 정책 | 범위 밖 | Go가 아니면 filefunc 구조화가 쉽지 않다. 주석 문법이 아니라 gofmt 수준의 구조 강제 전략이 본질. 필요한 자들이 각자 해결할 영역 |
| 레지스트리 호스팅 | 낮음 | MVP 범위 밖 |
| LLM 어노테이션 품질 보장 | 낮음 | MVP 범위 밖 |

---

## 향후 구상

- **VSCode 확장**: func 모아보기 (타입별 메서드 그룹핑, feature별 func 목록)
- **그래프 GUI 시각화**: feature 선택 → 관련 func 노드를 input/output 타입으로 연결한 그래프 표시. 코드북의 feature가 곧 줌 레벨.

---

## 결론

핵심 아이디어는 강하다. MVP는 validate(파일 구조 룰 검증)부터 만들어서 filefunc 자체에 적용하고, annotate → chain 순으로 확장.
