# filefunc rulebook

filefunc 룰의 단일 진실 공급원(SSOT). 모든 문서(CLAUDE.md, manual-for-ai.md)와 코드(validate_graph*.go)는 이 문서를 기준으로 동기화한다.

---

## P: 프로젝트 룰

프로젝트 레벨 구조 룰. validate 실행 시 파일 검증 이전에 먼저 체크.

| # | 룰 | 위반 시 | 검증 방법 |
|---|---|---|---|
| P1 | 프로젝트는 단일 언어. `.go`와 `.py`가 공존하면 ERROR | ERROR | WalkFiles로 양쪽 확장자 존재 여부 (.ffignore 적용 후) |

---

## F: 파일 구조 룰

| # | 룰 | 위반 시 | Go | Python |
|---|---|---|---|---|
| F1 | 파일 하나에 func 하나 (파일명 = 함수명). 테스트 파일 포함 | ERROR | `_test.go` 포함 | `test_*.py` 포함 |
| F2 | 파일 하나에 type 하나 (파일명 = 타입명) | ERROR | type | class |
| F3 | 1 file 1 method | ERROR | `receiver_method.go` | Mixin class (1 class 1 method). `__init__` 면제 |
| F4 | `init()` 단독 불허 (var 또는 func과 함께) | ERROR | Go 전용 | 해당 없음 |
| F5 | 의미적으로 한 묶음인 const/상수는 같은 파일 허용 | 예외 | const | 모듈 레벨 대문자 변수 |

### F3 Python Mixin 패턴

```python
# server_start.py — 1 class, 1 method → # ff:func
class ServerStartMixin:
    def start(self): ...

# server.py — 합성 class, __init__ 허용 → # ff:type
class Server(ServerStartMixin, ServerStopMixin):
    def __init__(self, config): ...
```

---

## Q: 코드 품질 룰

| # | 룰 | 위반 시 | 검증 방법 |
|---|---|---|---|
| Q1 | nesting depth: sequence=2, selection=2, iteration=dimension+1 | ERROR | MaxDepth vs 상한 비교 |
| Q2 | func max 1000 lines | ERROR | line count |
| Q3 | sequence func max 100 lines | ERROR | line count + control=sequence |
| Q4 | 제어문 body PURE 10줄 초과 시 sequence func으로 추출 | ERROR | body lines - inner control lines |

### Q1 depth 대상 제어문

| 구분 | Go | Python |
|---|---|---|
| 분기 | if, switch, type switch | if, match |
| 반복 | for, range | for, while |
| 비대상 | select | with, try/except, comprehension |

- `elif` / `else if`는 depth를 증가시키지 않음

### Q4 PURE 계산

PURE = 제어문 body 총 줄 수 - 내부 제어문 줄 수 (2-depth 면제).
switch/match는 case 절별로 개별 적용.

---

## A: 어노테이션 룰

| # | 룰 | 위반 시 | 검증 방법 |
|---|---|---|---|
| A1 | func 파일은 `//ff:func` (Go) / `# ff:func` (Python), type 파일은 `//ff:type` / `# ff:type` 필수 | ERROR | 어노테이션 유무 |
| A2 | 어노테이션 값은 코드북에 존재해야 함 | ERROR | codebook yaml 대조 |
| A3 | func 또는 type이 있는 파일은 `//ff:what` / `# ff:what` 필수 | ERROR | what 유무 |
| A6 | 어노테이션은 파일 최상단에 위치 | ERROR | 위치 검사 |
| A7 | `//ff:checked` / `# ff:checked` 해시 불일치 시 서명 깨짐 | ERROR | 해시 대조. LLM 검증은 `filefunc llmc` |
| A8 | required 코드북 키가 어노테이션에 모두 존재해야 함 | ERROR | codebook required 대조 |
| A9 | func 파일은 `control=` 필수 (sequence/selection/iteration) | ERROR | control 유무 + 값 검증 |
| A10 | `control=selection`인데 depth 1에 switch(Go)/match(Python) 없음 | ERROR | AST 검증 |
| A11 | `control=iteration`인데 depth 1에 loop 없음 | ERROR | AST 검증 |
| A12 | `control=sequence`인데 depth 1에 switch/match 또는 loop 존재 | ERROR | AST 검증 |
| A13 | `control=selection`인데 depth 1에 loop 존재 | ERROR | AST 검증 |
| A14 | `control=iteration`인데 depth 1에 switch(Go)/match(Python) 존재 | ERROR | AST 검증 |
| A15 | `control=iteration`이면 `dimension=` 필수 | ERROR | 어노테이션 유무 |
| A16 | `dimension=` 값은 양의 정수여야 함 | ERROR | 값 파싱 |

### 어노테이션 포맷

```
Go:     //ff:func feature=validate type=rule control=sequence
Python: # ff:func feature=validate type=rule control=sequence
```

| 어노테이션 | 필수 | 설명 |
|---|---|---|
| `//ff:func` / `# ff:func` | func 파일 | feature, type, control 등 코드북 값 |
| `//ff:type` / `# ff:type` | type 파일 | feature, type 등 코드북 값 |
| `//ff:what` / `# ff:what` | func/type 파일 | 이 함수/타입이 뭘 하는가 |
| `//ff:why` / `# ff:why` | 선택 | 왜 이렇게 만들었는가 (사용자 결정 근거) |
| `//ff:checked` / `# ff:checked` | 자동 (llmc) | LLM 검증 서명. 수동 작성 금지 |

### A6 최상단 규칙

| Go | Python |
|---|---|
| package 선언 이전 | shebang → encoding → `# ff:` → 빈 줄 → import/code |

---

## C: 코드북 룰

| # | 룰 | 위반 시 |
|---|---|---|
| C1 | `required` 섹션에 최소 1개 키와 1개 값 | ERROR |
| C2 | 같은 섹션 내 중복 키 없음 | ERROR |
| C3 | 키는 소문자 + 하이픈만 (`[a-z][a-z0-9-]*`) | ERROR |
| C4 | required 값의 description이 비어있음 | WARNING |

코드북은 파일 검증 이전에 먼저 체크. 코드북 위반 시 코드 검증 미실행.

---

## N: 네이밍 룰

| # | 룰 |
|---|---|
| N1 | 파일명: snake_case |
| N2 | 변수/함수: camelCase (Go), snake_case (Python) |
| N3 | 타입: PascalCase (Go), PascalCase (Python class) |
| N4 | gofmt 준수 (Go) |
| N5 | 에러 즉시 처리 (early return) |

---

## dimension (iteration 전용)

순회 대상 데이터의 차원 수. Q1 depth 상한 = dimension + 1.

| dimension | 의미 | depth 상한 |
|---|---|---|
| 1 | flat list 순회 | 2 |
| 2 | 2D 데이터 순회 | 3 |
| N | N차원 | N + 1 |

dimension >= 2는 named type (Go: struct/interface, Python: class) 중첩으로 설계해야 한다.

---

## 예외

- const/상수 전용 파일은 어노테이션 불필요 (F5 defeater)
- 프로젝트에 `//ff:checked` / `# ff:checked`가 하나도 없으면 A7 전체 건너뜀
- func 없는 파일은 A9-A16 면제 (HasNoFunc defeater)

---

## 검증 순서

```
1. P 룰 (프로젝트 레벨) — 혼합 언어 체크
2. C 룰 (코드북) — codebook.yaml 정합성
3. F/Q/A 룰 (파일 레벨) — toulmin defeats graph
```

P 또는 C 위반 시 후속 검증 미실행.
