# func control structure 분류

## 배경

filefunc의 F룰(1 file 1 func)과 Q1(depth ≤ 2)이 함수를 순수하게 만들면서, 함수의 제어구조(control structure)가 드러났다.

## 학술 근거

1966년 Böhm-Jacopini 정리(구조적 프로그램 정리): **sequence, selection, iteration 세 가지 제어구조의 조합으로 모든 계산 가능한 함수를 표현할 수 있다.**

## control = 제어구조 3분류

| control | 정통 용어 | 판별 기준 | 본질 |
|---|---|---|---|
| **sequence** | Sequence | 지배적 제어구조 없음 (순차 호출 + early return) | 모든 경로를 순서대로 실행 |
| **selection** | Selection | depth 1이 switch | 여러 경로 중 하나를 선택 |
| **iteration** | Iteration | depth 1이 loop | 같은 경로를 반복 실행 |

코드는 순차 진행하거나, 분기하거나, 반복한다. 이 밖의 제어구조는 없다.

### 어노테이션 사용

```go
//ff:func feature=validate type=rule control=selection
//ff:what 타입별 필수 필드 누락 검증
func validateRequiredFields(...) ...
```

```go
//ff:func feature=parse type=walker control=iteration
//ff:what 디렉토리를 재귀 순회하며 .go 파일 경로 목록 반환
func WalkGoFiles(...) ...
```

```go
//ff:func feature=cli type=command control=sequence
//ff:what validate 서브커맨드 정의 및 코드 구조 룰 검증 실행
func RunAll(...) ...
```

`control=sequence`는 기본값이므로 생략 가능.

### 제어구조로만 분류한다

template(문자열 리터럴), computation(계산) 등은 **본문의 내용 유형**이지 제어구조가 아니다. 내용 유형을 넣으면 끝없이 추가해야 한다. control은 제어구조 축만 담당.

template 함수의 제어구조는 sequence다(순차적으로 문자열 조립 → 파일 쓰기). 내용이 문자열일 뿐.

### "호출하는가 말단인가"는 control이 아니라 chain이 답한다

오케스트레이터(calls 13개)인지 말단(calls 0개)인지는 `filefunc chain`이 이미 제공. control에 중복으로 넣을 필요 없다.

### early return if는 selection이 아니라 탈출

서비스 함수의 `if err != nil { return nil, err }`는 depth 1이지만 실질 depth 0. 분기가 아니라 에러 탈출. sequence에 속한다.

## depth 2 제어구조 조합 분석

depth 1 요소: loop, if, switch (3개)
반복문(for, range)은 구조적으로 동일하므로 loop으로 통합.

### 9개 조합

| depth 1 → depth 2 | 패턴 | control | 빈도 |
|---|---|---|---|
| loop → loop | 중첩 반복 | **iteration** | 드뭄 |
| loop → if | 필터/선별 | **iteration** | 흔함 |
| loop → switch | 순회+분기 | **iteration** | 흔함 |
| if → loop | 조건부 반복 | **sequence** | 드뭄 |
| if → if | 중첩 조건 | **sequence** | early return으로 대체 |
| if → switch | 조건부 dispatch | **sequence** | 드뭄 |
| switch → loop | case별 반복 | **selection** | 드뭄 |
| switch → if | case별 조건 | **selection** | 흔함 |
| switch → switch | 중첩 분기 | **selection** | 드뭄 |

판별 규칙: depth 1이 **switch → selection**, **loop → iteration**, **if 또는 없음 → sequence**.

## AI 에이전트 컨텍스트 전략

control로 **body를 읽기 전에 "어떻게 읽을지"를 LLM 없이 기계적으로 결정**한다.

```
control=selection  → 전체를 한 번에 read. case를 부분적으로 읽으면 안 됨.
control=iteration  → body 핵심은 루프 내부. 루프 외부는 초기화.
control=sequence   → 수정 대상 step만 read. 나머지 step은 what으로 충분.
```

what이 "뭘 하는가", control이 "어떤 제어구조인가". read 전에 함수의 구조를 예측 가능.

## fullend Q3 위반 15건 재분류

| control | 해당 함수 | 건수 |
|---|---|---|
| selection | validateRequiredFields, ValidateWith, parseDDLTables, buildScenarioOrder, statusCmd | 5 |
| sequence | GenWith, generateMainWithDomains, generateServerStruct, generateCentralServer, transformSource, generateModelFile, generateQueryOpts, generateMain | 9 |
| iteration | loadPackageGoInterfaces | 1 |

15건 전부 3분류로 분류 가능. 분류 불가 0건.

## control별 Q 기준 (초안)

| control | Q1 기준 | Q3 기준 | 근거 |
|---|---|---|---|
| selection | switch 내부 depth ≤ 2 | 완화 (case 수에 비례) | switch 자체는 복잡도가 아닌 분기 나열 |
| iteration | depth ≤ 2 | 200줄 | 순회 내부 로직이 핵심 |
| sequence | depth ≤ 2 | 200줄 | 순차 나열이므로 길어도 복잡도 낮음 |

## 상태

**연구 단계**. fullend + filefunc 두 프로젝트 데이터 기반. filefunc 룰 개정(Q1, Q3, codebook control 키)의 기반 연구. 오픈소스 라이브러리 적용으로 데이터 수집 후 분류 체계 확정 및 룰 개정.

## 미결

- control별 Q3 기준의 구체적 숫자 확정
- AST 자동 분류 가능 여부 및 구현 시점
- codebook optional 키로 control을 추가하는 시점

## 연구 방향: 오픈소스 라이브러리 적용

유명 Go 라이브러리에 filefunc을 적용하여 control structure 데이터를 수집한다.

### 방법
1. 대상 라이브러리에 `filefunc validate` 실행 — Q1/Q3 위반 수집
2. 위반 함수를 filefunc 룰에 맞게 리팩토링 (1 file 1 func, depth ≤ 2)
3. 리팩토링 후 각 함수의 control 분류
4. 프로젝트 성격별 control 분포 분석

### 가설: 프로젝트 성격별 control 분포

| 프로젝트 유형 | 예상 지배 control |
|---|---|
| CLI 도구 (cobra 기반) | sequence |
| 파서/컴파일러 | selection + iteration |
| 웹 프레임워크 | sequence (handler) + selection (router) |
| 코드 생성기 | sequence |
| 데이터 처리 | iteration |

### 후보 라이브러리

| 라이브러리 | 유형 | 규모 |
|---|---|---|
| cobra | CLI 프레임워크 | 중 |
| gin | 웹 프레임워크 | 대 |
| esbuild | 컴파일러/번들러 (Go) | 대 |
| sqlc | 코드 생성기 | 중 |
| goldmark | Markdown 파서 | 중 |

### 기대 성과
- control 3분류의 범용성 검증 (3개로 충분한가, 4번째가 필요한가)
- 프로젝트 유형별 control 분포 데이터
- control별 Q3 기준의 실증적 근거
- filefunc 리팩토링 전후 코드 품질 비교 데이터

### 선행 연구
- Böhm-Jacopini 정리(1966): sequence, selection, iteration으로 모든 계산 가능 함수를 표현 가능
- 이 정리를 함수 본문의 지배적 제어구조로 분류하여 품질 기준에 적용하는 연구는 발견되지 않음
- filefunc의 depth 2 강제가 함수를 평탄화한 뒤에야 control 분류가 실용적으로 가능해짐
