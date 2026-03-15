# filefunc

## 핵심 원칙

**파일 하나에 개념 하나. 파일명 = 개념명.**

func이든 type이든 interface든 const 묶음이든 동일하다. 이 원칙 하나에서 모든 룰이 파생된다.

**1 file 1 concept. 이것만으로 read 한 번에 불필요한 코드가 딸려오지 않는다. 컨텍스트 오염 차단은 파일 구조 자체가 해결한다.**

### 왜 이 원칙인가

Claude Code는 grep → read로 코드를 탐색한다. read의 단위가 파일이므로, 파일이 개념의 단위와 일치해야 한다.

- 1 file 1 func → read 한 번 = func 하나
- 1 file 1 type → read 한 번 = type 하나

파일 하나에 func이 20개면, CrossError 하나가 필요해서 read했는데 19개가 딸려온다. 이게 컨텍스트 오염의 본질이다.

**필요한 5-10개를 집는 것보다, 불필요한 290개를 안 여는 게 더 중요하다.**

### 제1시민은 AI 에이전트다

filefunc의 제1시민은 사람이 아니라 AI 에이전트(Claude Code)다.

Claude Code는 `ls`가 아니라 `grep`으로 탐색한다. 파일 500개든 1000개든 `rg '//ff:func feature=validate'` 한 번이면 끝. 파일이 많을수록 각 파일이 작고, read 한 번에 딸려오는 노이즈가 줄어들어 오히려 유리하다.

파일 수 폭발은 약점이 아니다. 사람의 불편은 뷰 레이어(VSCode 확장 등)에서 해결한다. filefunc의 구조를 사람에게 맞춰 타협하지 않는다.

### Go 한정

Go가 아니면 filefunc 구조화가 쉽지 않다. gofmt가 코드 포맷을 강제하고, early return이 관례이고, 예외가 없고, 패키지 = 디렉토리. 다른 언어로 확장하려면 주석 문법이 아니라 gofmt 수준의 구조 강제 전략이 필요하다. 이는 filefunc의 범위 밖이며, 필요한 자들이 각자 해결할 영역이다.

---

## Claude Code 작업 동선

### 기존
```
사용자 요청
→ 뭐가 있는지 몰라서 ls, find
→ 파일 열어보고 구조 파악
→ 관련 파일 찾으러 또 grep
→ 열었더니 func 20개, 대부분 불필요
→ 탐색 비용 > 실제 작업 시간
```

### filefunc
```
사용자 요청 + 코드북 제공
→ 코드북 보고 grep 쿼리 즉시 구성
→ 파일 20-30개 read (각각 1개념, 전부 유효 컨텍스트)
→ 작업
```

코드북이 Claude Code의 프로젝트 지도다. 코드북이 없으면 어휘를 모르는 상태로 탐색을 시작한다. 코드북이 있으면 `feature=crosscheck`, `type=rule` 같은 정확한 쿼리를 탐색 없이 바로 던진다.

파일 수가 많아도 괜찮다. 30개를 read해도 전부 유효한 컨텍스트면 30개가 문제가 아니다. 1개를 read했는데 30개 분량이 딸려오는 게 문제다.

---

## 룰

### 파일 구조 룰

| # | 룰 | 위반 시 |
|---|---|---|
| 1 | 파일 하나에 func 하나 (파일명 = 함수명) | ERROR |
| 2 | 파일 하나에 type 하나 (파일명 = 타입명) | ERROR |
| 3 | 각 파일은 최대 1개의 `init()`을 선택적으로 가질 수 있다. `init()`만 단독 불허 (var 또는 func과 함께) | ERROR |
| 4 | 메서드: 1 file 1 method (`server_start.go`, `server_stop.go`) | ERROR |
| 5 | `_test.go`는 복수 func 허용 | 예외 |
| 6 | 함수 전용 파라미터 타입은 해당 func 파일에 함께 허용 | 예외 |
| 7 | 의미적으로 한 묶음인 const는 같은 파일 허용 | 예외 |

### 코드 품질 룰

| # | 룰 | 위반 시 | 검증 방법 |
|---|---|---|---|
| 1 | nesting depth 2까지만 허용 | ERROR | tree-sitter depth 측정 |
| 2 | func max 1000 lines | ERROR | line count |
| 3 | func 권고 100 lines | WARNING | line count |
| 4 | 정형 구조 강제 (CLI는 cobra 등) | ERROR | import 검사 |

> nesting depth 2는 Go의 early return 패턴으로 해결한다. if err != nil 중첩이 쌓이면 룰 위반 신호다. "반복 1depth, 분기 1depth"로 세분화하는 안도 검토했으나, 이중 반복 등 실무 패턴을 과도하게 제한하므로 "depth 2" 하나로 확정.

### 어노테이션 룰

| # | 룰 | 위반 시 | 검증 방법 |
|---|---|---|---|
| 1 | func이 있는 파일은 `//ff:func`, type이 있는 파일은 `//ff:type` 필수 | ERROR | `//ff:func` 또는 `//ff:type` 유무 |
| 2 | 어노테이션 값은 코드북에 존재해야 함 | ERROR | codebook yaml 대조 |
| 3 | func 또는 type이 있는 파일은 `//ff:what` 필수 | ERROR | `//ff:what` 유무 |
| 4 | input/output 타입 명시 | ERROR | Go AST / tree-sitter |
| 5 | calls/uses는 실제 코드와 정합성 검증 | ERROR | tree-sitter 교차검증 |
| 6 | 어노테이션은 파일 최상단에 위치 | ERROR | 위치 검사 |
| 7 | `//ff:checked` 해시 불일치 시 서명 깨짐 | ERROR | `validate`는 해시 대조만 (읽기 전용). LLM 검증은 `filefunc llmc`로 별도 실행 |

---

## 메타데이터 어노테이션

```go
//ff:func feature=crosscheck type=rule source=SSaC target=OpenAPI
//ff:what SSaC 함수명↔OpenAPI operationId 양방향 정합성 검증
//ff:why 제1시민은 AI 에이전트. 파일 수 폭발은 약점이 아니라 장점이다
//ff:calls check_response_fields, check_err_status
//ff:uses CrossError, ServiceFunc
func CheckSSaCOpenAPI(funcs []ServiceFunc, st *SymbolTable, doc *openapi3.T, specs []FuncSpec) []CrossError {
```

어노테이션은 파일 최상단에 위치해야 한다. body 전체를 read하지 않아도 상단 5줄로 메타 파악이 가능하도록.

| 어노테이션 | 내용 | 필수 | 이유 |
|---|---|---|---|
| `//ff:func` | func 파일의 feature, type 등 메타 | func 파일 필수 | 불변, 짧음, grep 대상 |
| `//ff:type` | type 파일의 feature, type 등 메타 | type 파일 필수 | 불변, 짧음, grep 대상 |
| `//ff:what` | 1줄 설명 (이 함수/타입이 뭘 하는가) | func/type 파일 필수 | body 안 읽어도 파악. **소형 LLM이 body와 대조하여 일치 여부를 검증하는 기준** |
| `//ff:why` | 왜 이렇게 만들었는가 | 선택 | 사용자의 요구/결정이 근거. AI의 판단이나 추측은 why가 아니다. 검증 불가, 이력으로 남기는 것 |
| `//ff:calls` | 호출하는 함수 목록 | 자동 생성 | Go AST에서 추출. 도구가 생성/관리 |
| `//ff:uses` | 사용하는 타입 목록 | 자동 생성 | Go AST에서 추출. 도구가 생성/관리 |
| `//ff:checked` | LLM 검증 서명 | 자동 생성 | `llm=모델명 hash=body해시`. body 변경 시 자동 무효화 |

형식: `//ff:key key1=value1 key2=value2`

- grep/ripgrep으로 즉시 검색 가능 (`rg '//ff:'`)
- 정형화된 key-value로 도구가 파싱 가능
- Go의 `//go:generate`, `//go:embed` 관례와 동일한 패턴

---

## 코드북

코드북은 filefunc 설계에서 가장 중요한 위치를 차지한다. 어노테이션 룰보다 코드북이 먼저다. 코드북이 잘 설계되어야 grep 쿼리가 정밀해지고, grep이 정밀해야 read 목록이 깨끗해진다.

```yaml
# codebook.yaml
feature: [crosscheck, validate, gen, parse, report, contract, orchestrate]
type: [rule, parser, validator, generator, handler, middleware, loader, util]
pattern: [rulebook, target-interface, symbol-table, error-collection]
level: [ERROR, WARNING, INFO]
```

- 코드북에 없는 값을 어노테이션에 쓰면 validate ERROR
- 프로젝트 시작 시 Claude Code에 코드북을 제공하면 탐색 없이 바로 grep 가능
- 코드북 설계 품질이 grep 정밀도를 좌우한다 — 의도된 트레이드오프. 대응 전략: 도메인을 날카롭게 좁히고, 프로젝트마다 코드북을 맞춤 작성한다. 코드북으로 어휘를 정규화하면 빠진 feature, 중복된 type, 애매한 분류가 목록에서 드러난다. 구멍이 보여야 관리가 된다.

---

## func 노드 그래프

func의 본질만 남긴 계약 그래프:

```yaml
node:
  name: CheckSSaCOpenAPI
  input: [ServiceFuncs, SymbolTable, OpenAPIDoc, FuncSpecs]
  output: [CrossError[]]
  what: "SSaC 함수명↔OpenAPI operationId 양방향 정합성 검증"
```

- body 안 읽어도 함수가 뭘 하는지 앎
- 그래프 순회로 의존성 체인 자동 추적
- 입출력 타입만으로 연결 가능한 함수 조합 자동 제안

### func chain

```
ParseAll() → ParsedSSOTs
                ↓
        Run(CrossValidateInput) → []CrossError
                                      ↓
                              Print(Report)
```

```bash
filefunc chain func CheckSSaCOpenAPI   # 이 함수의 데이터 흐름 추적
filefunc chain feature crosscheck      # crosscheck feature 전체 체인
```

기존 `go callgraph`와의 차이: callgraph는 모든 호출을 정적 분석해 수천 노드가 나온다. func chain은 같은 feature 안에서만 input/output 타입 매칭으로 연결한다. 코드북의 feature가 곧 줌 레벨이다.

---

## LLM 탐색 파이프라인

어노테이션이 SIDX와 같은 역할을 한다. 벡터 임베딩 같은 무거운 인프라 없이 동작한다.

```
1. 구조적 축소 (LLM 불필요, grep)
   코드북 기반으로 grep 쿼리 구성
   → 후보 파일 20-30개 추출

2. 메타 판정 (LLM 불필요 또는 초소형)
   각 파일 상단 어노테이션만 read
   → name/input/output/what으로 실제 필요한 파일 5-10개로 좁힘

3. 정밀 작업 (대형 LLM, 최소 컨텍스트)
   5-10개 파일만 full read
   → 코드 수정/생성
```

---

## func = file의 부수 효과

### whyso 연동

func = file이므로 함수 단위 변경 이력이 파일 단위로 정확히 떨어진다.

```bash
whyso history check_ssac_openapi.go   # CheckSSaCOpenAPI 함수의 변경 이력
```

지금은 한 파일에 func 여러 개 있으면 어느 함수가 바뀐 건지 diff를 뒤져야 한다. filefunc면 파일 변경 = 함수 변경. 추적 비용 제로.

### 암묵적 커플링 검출

```bash
whyso coupling check_ssac_openapi.go

같은 요청에 함께 수정된 함수:
  check_response_fields.go  8회
  check_err_status.go       5회
  types.go                  4회
```

calls/uses에 명시적 관계가 없는데 coupling 통계에서 자꾸 나오면 숨은 의존성 신호다.

- 같은 비즈니스 규칙을 다른 각도에서 구현한 함수들
- interface 없이 암묵적으로 format을 맞추고 있는 것들
- 버그가 항상 같이 터지는 것들

자동 WARNING 가능: "이 두 함수는 명시적 관계 없이 8회 함께 수정됨. 의존성을 명시하세요."

---

## LLM 자동 어노테이션

기존 코드를 안 건드리고 메타데이터만 씌운다.

```bash
filefunc annotate ./internal/       # LLM이 func 읽고 어노테이션 + what 자동 생성
filefunc validate ./internal/       # 정합성 검증
```

비침투적. 기존 Go 라이브러리든 뭐든 코드 그대로, 메타데이터 레이어만 추가. 도입 저항 제로.

---

## 어노테이션 drift 방지

- `//ff:calls`, `//ff:uses`: 코드와 기계적 대조 가능 (tree-sitter 교차검증)
- `//ff:what`: 자연어라 기계적 검증 불가 → **소형 LLM 검수 + `//ff:checked` 서명으로 해결**

### //ff:checked 메커니즘

```go
//ff:func feature=validate type=rule
//ff:what F1: 파일당 func 1개 검증
//ff:checked llm=gpt-oss-20b hash=a3f8c1d2
func CheckOneFileOneFunc(...) ...
```

**validate와 llmc의 역할 분리:**

- `filefunc validate`: checked 해시 대조만. 불일치 시 WARNING + checked 삭제. LLM 없음. 어디서든 실행 가능.
- `filefunc llmc`: 소형 LLM으로 what-body 일치 판정. 통과 시 `//ff:checked llm=모델명 hash=body해시` 기록. 불일치 시 ERROR.

body가 수정되면 해시가 달라지므로 validate가 checked를 삭제한다. llmc를 다시 돌려야 서명이 복원된다.

---

## 레지스트리

GitHub 스타 10k+ 유명 라이브러리의 func 메타를 미리 정의한다. 사용자가 `import "github.com/gin-gonic/gin"` 하면 해당 func 메타가 자동으로 그래프에 편입.

| 대상 | 전략 | filefunc 룰 강제 |
|---|---|---|
| 외부 라이브러리 | public API 메타만 추출. 내부 구조 무관 | X |
| 자체 코드 | filefunc 룰 전면 강제. validate 통과 필수 | O |

초기 구축 파이프라인: pkg.go.dev API로 public API signature 수집 → LLM 1차 draft → 사람 검증 → 레지스트리 merge.

`filefunc validate` 통과 못하면 레지스트리 등록 불가. 생태계에 올라온 라이브러리는 구조가 보장됨.

---

## filefunc CLI

```bash
filefunc validate ./internal/           # 코드 구조 룰 검증 (LLM 없음)
filefunc annotate ./internal/           # calls/uses 자동 산출 (AST)
filefunc llmc ./internal/               # LLM으로 what-body 일치 검증 + checked 서명
filefunc chain func CheckSSaCOpenAPI    # 함수 데이터 흐름 추적
filefunc chain feature crosscheck       # feature 전체 체인
```

---

## 학술 근거

- **"Lost in the Middle" (Stanford, 2024)** — 관련 정보가 컨텍스트 중간에 있으면 성능 30% 이상 하락. 긴 컨텍스트 전용 모델에서도 동일.
- **"Context Length Alone Hurts LLM Performance" (Amazon, 2025)** — 불필요한 토큰이 공백이어도 성능 하락 (13.9~85%). 관련 정보만 추출한 짧은 컨텍스트가 압도적 우수.
- **"Context Rot" (Chroma Research)** — 모든 모델에서 focused prompt > full prompt 성능 확인.

연구는 "컨텍스트 짧을수록 좋다"고 증명했지만, 코드를 구조적으로 짧게 쪼개서 필요한 것만 넣는 도구가 없었다. filefunc가 그 빈자리다.

---

## 미결 사항

- 파라미터 개수 상한 (N개 초과 시 struct로 묶어라)
- 레지스트리 호스팅/배포 방식
- LLM 어노테이션 품질 보장 (사람 리뷰 필수? 자동 승인?)
