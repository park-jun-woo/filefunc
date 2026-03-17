# 타입 기반 데이터 순회를 위한 차원 인식 네스팅 깊이 상한

## 초록

정적 분석 도구는 모든 함수에 획일적 네스팅 깊이 상한(보통 3~5)을 적용한다. 우리는 타입화된 다차원 데이터 순회라는 함수 범주를 식별한다 — 이 범주에서 네스팅 깊이는 논리적 복잡도가 아닌 데이터 구조의 차원을 반영한다. 이 함수들에 깊이-2 상한을 강제하면 순회 경로가 여러 함수에 분산되어, 복잡도를 줄이지 않으면서 인지 부하를 증가시킨다.

**차원 인식 깊이 상한**을 제안한다: N차원 타입화된 데이터를 순회하는 함수의 깊이 상한은 기본 2가 아닌 N+1로 한다. 차원 수는 for-range 체인의 named type 단언 횟수에서 자동 도출된다. Go AST, OpenAPI, protobuf 순회 패턴에서 평가한 결과, 타입화된 순회에서의 깊이 위반 100%가 코드 명확성을 희생하지 않고 차원 인식 상한으로 해소됨을 보인다.

## 1. 서론

### 1.1 복잡도 지표로서의 네스팅 깊이

네스팅 깊이는 코드 복잡도의 널리 사용되는 대리 지표다. 깊은 중첩은 다음을 나타낸다:
- 다중 수준의 조건 로직
- 복잡한 반복 패턴
- 추적하기 어려운 제어 흐름

SonarQube, golangci-lint, ESLint 같은 도구가 최대 깊이 상한을 강제한다. 이 상한은 대부분의 코드에 잘 작동한다. 그러나 한 가지 패턴을 체계적으로 오분류한다: **타입화된 데이터 구조 순회**.

### 1.2 순회 예외

Go AST를 순회하여 struct 필드를 수집하는 경우를 보자:

```go
for _, entry := range entries {                    // 깊이 1: 디렉토리 파일
    for _, decl := range f.Decls {                 // 깊이 2: 선언
        for _, spec := range gd.Specs {            // 깊이 3: 타입 스펙
            for _, field := range st.Fields.List {  // 깊이 4: struct 필드
                for _, name := range field.Names {  // 깊이 5: 필드 이름
                }
            }
        }
    }
}
```

깊이 5. 어떤 표준 지표로든 "너무 복잡하다." 그러나 각 단계는:
1. 타입 단언 (`decl.(*ast.GenDecl)`, `spec.(*ast.TypeSpec)` 등)
2. 단언 실패 시 early-continue
3. 분기 로직 없음 — 한 단계 더 내려갈 뿐

깊이는 코드의 논리적 복잡도가 아닌 **데이터 구조의 차원**을 반영한다. Go AST 스키마는 언어 사양에 의해 고정되어 있다: `File -> GenDecl -> TypeSpec -> StructType -> Field -> Ident`. 이 구조는 단순화할 수 없다.

### 1.3 함수 추출이 실패하는 이유

깊은 중첩에 대한 표준 처방은 함수 추출이다:

```go
// 추출 후: 순회가 3개 함수로 분산
func loadGoInterfaces(dir string) {
    for _, entry := range entries {
        processFileDecls(f, st)          // 내부 순회가 숨겨짐
    }
}
func processFileDecls(f *ast.File, st *SymbolTable) {
    for _, decl := range f.Decls {
        processGenDeclSpecs(gd, st)      // 또 숨겨짐
    }
}
```

문제:
- **분산된 순회 경로**: 완전한 경로(entries -> decls -> specs -> fields -> names)가 3개 이상 함수에 분리됨. 전체 순회를 이해하려면 모두 읽어야 함.
- **지역성 상실**: 원본은 위에서 아래로 연속적인 하강으로 읽힘. 추출 버전은 함수 간 점프 필요.
- **복잡도 미감소**: 추출된 각 함수는 여전히 같은 일을 함 — 타입 단언 후 반복. 총 복잡도는 변하지 않고 이동함.

## 2. 타입화된 데이터 순회

### 2.1 특성

**타입화된 데이터 순회**는 다음 속성을 갖는다:

| 속성 | 설명 |
|---|---|
| 고정 스키마 | 데이터 구조의 형태가 타입 시스템(AST, API 스펙, protobuf)에 의해 정의 |
| 각 레벨에 named type | 각 차원이 named type(struct, interface)으로 모델링 |
| 하강 수단으로서의 타입 단언 | 더 깊이 내려가기 = 타입화된 컨테이너를 단언 |
| 분기 로직 없음 | 각 레벨: 타입 단언 -> 실패 시 early-continue -> 다음 레벨 반복 |
| 결정론적 경로 | 순회 경로가 모든 입력 인스턴스에 대해 동일 |

### 2.2 차원 수

순회의 **차원**은 for-range 체인의 named type 단언 횟수다:

```go
gd, ok := decl.(*ast.GenDecl)          // named type 1
ts, ok := spec.(*ast.TypeSpec)         // named type 2
st, ok := ts.Type.(*ast.StructType)    // named type 3
```

여기서 차원 = 3. 각 단언은 타입화된 데이터 구조에서 한 단계 내려가는 것에 대응.

### 2.3 도메인별 예시

| 도메인 | 데이터 구조 | 순회 경로 | 차원 |
|---|---|---|---|
| Go AST | 소스 파일 | File -> GenDecl -> TypeSpec -> StructType -> Field | 5 |
| Go AST | 인터페이스 | File -> GenDecl -> TypeSpec -> InterfaceType -> Method | 5 |
| OpenAPI | API 스펙 | PathItem -> Operation -> Parameter -> Schema | 4 |
| Protobuf | 디스크립터 | FileDescriptor -> MessageDescriptor -> FieldDescriptor | 3 |
| HTML DOM | 문서 | Document -> Element -> Attributes -> Value | 3 |

### 2.4 Raw 중첩은 해당하지 않음

```go
// 타입화된 순회가 아님:
var grid [][][]int
for i := range grid {
    for j := range grid[i] {
        for k := range grid[i][j] {
            // grid[i][j][k] 처리
        }
    }
}
```

Named type 없음. 타입 단언 없음. 이것은 raw 다차원 반복이다 — 깊이가 진짜 복잡도를 반영한다. 프로그래머는 named type으로 모델링해야 한다:

```go
type Grid struct { Layers []Layer }
type Layer struct { Rows []Row }
type Row struct { Cells []int }
```

차원 인식 예외는 **좋은 데이터 모델링을 강제한다**: named type 순회만 자격을 얻는다. 이것은 한계가 아니라 기능이다.

## 3. 차원 인식 깊이 상한

### 3.1 룰

`control=iteration`과 `dimension=N`으로 어노테이션된 함수에 대해:

```
깊이_상한 = N + 1
```

+1은 가장 안쪽 레벨에서의 작업 로직(필터링, 수집, 변환)을 위한 것.

| dimension | 깊이 상한 | 전형적 패턴 |
|---|---|---|
| 1 | 2 | 단층 리스트 순회 (기본값) |
| 2 | 3 | 2차원 데이터 (paths -> operations) |
| 3 | 4 | 3차원 데이터 (policies -> rules -> actions) |
| 4 | 5 | 4차원 데이터 (spec 중복 제거) |
| 5 | 6 | 전체 AST 순회 |

### 3.2 어노테이션 형식

```go
//ff:func feature=symbol type=loader control=iteration dimension=3
//ff:what 디렉토리에서 Go 인터페이스를 파싱하여 "pkg.Model" 키로 등록
func loadGoInterfaces(dir string, st *SymbolTable) error {
```

### 3.3 검증 룰

| 룰 | 조건 | 심각도 |
|---|---|---|
| A15 | `control=iteration`이면 `dimension=` 필수 | ERROR |
| A16 | `dimension=` 값은 양의 정수여야 함 | ERROR |
| Q1 | 실제 깊이 <= dimension + 1이어야 함 | ERROR |

자동 차원 검증: `filefunc validate`가 for-range 체인의 named type 단언 횟수를 세고, 선언된 `dimension=N`이 실제 단언 횟수와 일치하는지 검증.

## 4. Visitor 패턴: 해결책이 아님

### 4.1 ast.Inspect 변환

```go
for _, entry := range entries {
    ast.Inspect(f, func(n ast.Node) bool {
        ts, ok := n.(*ast.TypeSpec)
        if !ok { return true }
        st, ok := ts.Type.(*ast.StructType)
        if !ok { return false }
        for _, field := range st.Fields.List {
            for _, name := range field.Names {
                // 처리
            }
        }
        return false
    })
}
```

깊이: 8 -> 4. 개선이지만 여전히 깊이-2 위반. 그리고 새로운 문제 도입:

### 4.2 Visitor 패턴의 문제

| 문제 | 설명 |
|---|---|
| 에러 전파 불가 | `ast.Inspect` 콜백은 bool만 반환. `return fmt.Errorf(...)` 불가 |
| 불필요한 순회 | 함수 본문, 표현식 등 무관한 노드까지 전부 방문 |
| 부모 컨텍스트 소실 | 콜백이 `*ast.TypeSpec`을 받을 때 어떤 GenDecl 소속인지 모름 |
| 암묵적 상태 | 상태 플래그를 가진 collector struct 필요 (`inRequestStruct bool`) |

원본의 명시적 중첩이 **암묵적 상태 머신**(visitor)으로 바뀔 뿐, 복잡도가 감소하지 않고 형태만 바뀐다.

### 4.3 Visitor가 적합한 경우

Visitor는 모든 노드가 같은 타입이고 자기 완결적인 **동질적 트리**에 잘 작동한다:

```go
// 좋은 visitor 사용 예: 모든 노드가 같은 타입
tree.Walk(func(node *Node) {
    process(node)  // 각 노드를 독립적으로 처리 가능
})
```

Go AST는 **이질적**(각 레벨마다 다른 타입)이고 **컨텍스트 의존적**(처리가 부모 타입에 따라 달라짐). Visitor는 잘못된 추상화다.

## 5. 평가

### 5.1 데이터

세 Go 프로젝트(filefunc, whyso, fullend)에서 리팩토링 전 깊이 > 2인 모든 함수를 식별:

| 프로젝트 | 총 깊이 위반 | 타입화된 순회 위반 | 비순회 위반 |
|---|---|---|---|
| filefunc | 0 (룰 준수 개발) | 0 | 0 |
| whyso | 23 | 2 | 21 |
| fullend | 148+ | 13 | 135+ |

### 5.2 타입화된 순회 위반

두 프로젝트의 타입화된 순회 위반 15건 전부:

| 함수 | 원래 깊이 | 차원 | 새 상한 | 해소? |
|---|---|---|---|---|
| loadPackageGoInterfaces | 8 | 5 | 6 | 예 |
| loadGoInterfaces | 6 | 5 | 6 | 예 |
| (fullend AST 워커) | 4~6 | 2~4 | 3~5 | 전부 예 |

차원 인식 상한으로 타입화된 순회 위반의 100%가 해소.

### 5.3 비순회 위반

나머지 156건 이상의 위반(비타입화된 순회)은 다음으로 해소:

| 기법 | 건수 | 설명 |
|---|---|---|
| early return/continue | ~40 | 중첩 if를 가드 절로 대체 |
| 조건 병합 | ~20 | `if a { if b {`를 `if a && b {`로 결합 |
| 헬퍼 추출 | ~96 | 내부 로직을 별도 함수로 추출 |

이것들은 깊이-2 강제의 혜택을 받는 진짜 복잡도다. 차원 인식 예외가 이들을 정확히 제외한다.

## 6. 설계 결정

### 6.1 Named Type 요구의 이유

각 차원 레벨에서 named type을 요구하는 것은 이중 목적:

1. **검증**: 도구가 타입 단언을 세어 선언된 차원을 검증 가능
2. **설계 압력**: 개발자가 raw 중첩(`[][][]int`) 대신 named type으로 데이터를 모델링하도록 강제

이것은 코드 품질을 개선하면서 동시에 깊이 예외를 가능케 하는 의도적 제약이다.

### 6.2 dimension + 1인 이유, dimension이 아닌 이유

+1은 가장 안쪽 레벨에서의 실제 작업을 위한 것:

```go
for _, field := range st.Fields.List {  // 깊이 = dimension
    for _, name := range field.Names {  // 깊이 = dimension + 1
        result[name.Name] = fieldType   // 가장 깊은 레벨에서의 실제 로직
    }
}
```

+1 없이는 순회 자체가 허용 깊이를 모두 소비하여, 각 요소를 처리하는 로직을 위한 여유가 없다.

### 6.3 재귀 깊이를 다루지 않는 이유

일부 순회는 중첩 루프 대신 재귀를 사용한다. 재귀 깊이는 정적으로 검증하기 어려우며 별도의 문제다. 본 논문은 정적으로 알려진 깊이를 가진 반복적 다차원 순회만 다룬다.

## 7. 관련 연구

- **순환 복잡도** (McCabe, 1976): 모든 경로를 센다. 데이터 주도 중첩과 로직 주도 중첩을 구분하지 않음.
- **인지 복잡도** (SonarSource, 2017): 중첩에 점증적 페널티. 타입화된 순회도 여전히 높은 점수.
- **SonarQube 깊이 상한**: 3~5로 고정. 차원 인식 없음.
- **golangci-lint nestif**: 설정 가능한 깊이. 타입화된 순회 예외 없음.

이 도구들 중 어느 것도 "데이터 차원에 의한 깊이"와 "로직 복잡도에 의한 깊이"를 구분하지 않는다.

## 8. 한계

- Go 전용 평가. 다른 언어는 다른 순회 패턴을 가짐 (예: Python 제너레이터, Rust 이터레이터).
- 차원 검증은 타입 단언의 Go AST 분석을 필요로 함. 다른 언어는 다른 감지 전략이 필요할 수 있음.
- 반복적 순회만 다룸. 재귀 순회(트리 워커, 그래프 DFS)는 범위 밖.
- 소규모 데이터셋 (타입화된 순회 위반 15건). 다양한 Go 프로젝트에서의 대규모 평가 필요.

## 9. 결론

네스팅 깊이 상한은 데이터에 의한 깊이와 로직에 의한 깊이를 구분해야 한다. 타입화된 다차원 데이터 순회에서 깊이는 데이터 구조의 차원을 반영한다 — 스키마의 고정 속성이지 코드 복잡도의 가변 속성이 아니다.

차원 인식 상한(`깊이 <= dimension + 1`)은 다른 모든 코드에 대한 엄격한 깊이 상한의 혜택을 유지하면서 타입화된 순회 위반의 100%를 해소한다. Named type 요구는 잘 모델링된 데이터 구조만 예외 자격을 얻도록 보장하여, 긍정적인 설계 압력을 만든다.

---

## 인용

```
@misc{dimension-nesting2026,
  title={Dimension-Aware Nesting Depth Limits for Typed Data Traversal},
  author={Park, Jun-Woo},
  year={2026},
  url={https://github.com/park-jun-woo/filefunc}
}
```
