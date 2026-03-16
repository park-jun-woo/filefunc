# Q1 depth ≤ 2 — 다차원 데이터 순회 예외 보고

## 요약

iteration 대상이 2차원 이상의 데이터 구조일 때, depth는 데이터 차원을 반영할 뿐 논리 복잡도가 아니다. Q1 ≤ 2를 강제하면 순회 흐름이 여러 함수로 분산되어 오히려 인지 복잡도가 증가한다.

핵심 원칙: **다차원 데이터의 각 차원은 named type으로 모델링되어야 한다.** 타입 단언 체인의 named type 수 = 차원 수. 이 차원 수만큼 Q1을 완화한다.

## 대상 코드

fullend `internal/ssac/validator/`에서 `go/ast`를 import하는 8개 파일 중 Q1 위반 2건:

| 파일 | depth | 패턴 |
|---|---|---|
| `load_package_go_interfaces.go` | 8 | entries→decls→specs→struct→fields→names |
| `load_go_interfaces.go` | 6 | entries→decls→specs→iface→methods→names |

나머지 6개 파일(`contract/splice.go`, `contract/scan.go`, `is_context_type.go`, `expr_to_go_type.go`, `ssac/parser/parser.go`, `funcspec/parser.go`)은 Q1 위반 없음.

## Go AST의 고정 스키마

Go AST는 아래 트리 구조가 **언어 사양에 의해 고정**되어 있다:

```
File
 └─ Decls[]
     ├─ GenDecl
     │   └─ Specs[]
     │       └─ TypeSpec
     │           ├─ InterfaceType → Methods[] → FuncType → Params[] → Names[]
     │           └─ StructType → Fields[] → Names[]
     └─ FuncDecl
         └─ Doc → Comments[]
```

각 레벨은 **named type**이다. `ast.File`, `ast.GenDecl`, `ast.TypeSpec`, `ast.StructType`, `ast.Field` — 각 차원이 타입 시스템에 명시되어 있다. 이 구조는 변경할 수 없다.

## depth 발생 구조 분석

### load_package_go_interfaces.go (depth 8)

```go
for _, entry := range entries {                    // depth 1: 디렉토리 파일 순회
    for _, decl := range f.Decls {                 // depth 2: 선언 순회 — Q1 한계
        for _, spec := range gd.Specs {            // depth 3: 위반 시작
            for _, field := range st2.Fields.List { // depth 4
                for _, name := range field.Names {  // depth 5
                }
            }
        }
    }
}
```

각 루프의 타입 단언:

| depth | 타입 단언 | named type | 분기 |
|---|---|---|---|
| 1 | `entry` → `parser.ParseFile()` → `*ast.File` | `ast.File` | early-continue |
| 2 | `decl.(*ast.GenDecl)` | `ast.GenDecl` | early-continue |
| 3 | `spec.(*ast.TypeSpec)` → `ts.Type.(*ast.StructType)` | `ast.TypeSpec`, `ast.StructType` | early-continue |
| 4 | `st2.Fields.List` | `ast.Field` | 없음 |
| 5 | `field.Names` | `ast.Ident` | 없음 |

named type 단언 횟수: **5회**. 각 depth에 판단 로직 없음. 타입 단언 실패 시 continue, 성공 시 다음 레벨로 내려갈 뿐.

### load_go_interfaces.go (depth 6)

```go
for _, entry := range entries {                         // depth 1
    for _, decl := range f.Decls {                      // depth 2: Q1 한계
        for _, spec := range gd.Specs {                 // depth 3
            for _, method := range iface.Methods.List { // depth 4
                if len(method.Names) > 0 {              // depth 5
                    // 메서드 등록                        // depth 6
                }
            }
        }
    }
}
```

named type 단언: `ast.File` → `ast.GenDecl` → `ast.TypeSpec` → `ast.InterfaceType` → `ast.Field` = **5회**. 동일 패턴.

## 왜 extract-func가 복잡도를 올리는가

```go
// 원본: 한 함수에서 순회 전체가 보임
func loadGoInterfaces(dir string) {
    for _, entry := range entries {           // depth 1
        for _, decl := range f.Decls {        // depth 2
            for _, spec := range gd.Specs {   // depth 3 ← 위반
                processTypeSpec(ts, st)
            }
        }
    }
}

// 추출 후: 순회가 2개 함수로 분산
func loadGoInterfaces(dir string) {
    for _, entry := range entries {
        processFileDecls(f, st)               // 내부 순회가 숨겨짐
    }
}
func processFileDecls(f *ast.File, st *SymbolTable) {
    for _, decl := range f.Decls {
        processGenDeclSpecs(gd, st)           // 또 숨겨짐
    }
}
```

- 원본: 순회 경로가 **한 함수에서 연속적으로 보임**. entries→decls→specs→fields 흐름을 위에서 아래로 읽으면 끝.
- 분리 후: 순회 경로를 파악하려면 **3~4개 함수를 넘나들어야 함**. 각 함수는 단순하지만 전체 맥락이 흩어짐.

### 핵심 차이: 다차원 데이터 순회 vs 로직 중첩

| 구분 | 다차원 데이터 순회 | 동적 트리/그래프 순회 | 분기 중첩 |
|---|---|---|---|
| 경로 예측 | 항상 동일 (스키마 고정) | 실행 시 결정 | 조건에 따라 분기 |
| 각 depth의 로직 | named type 단언 + continue | 재귀/큐 + 방문 체크 | if/switch 판단 |
| depth 제한 효과 | 맥락 분산 (복잡도 ↑) | 관심사 분리 (복잡도 ↓) | 평탄화 (복잡도 ↓) |
| 예시 | Go AST, OpenAPI, protobuf | 재귀 DFS, BFS, visitor | for→if→if, for→switch→if |

## visitor 패턴 전환 검토

### ast.Inspect 적용 시

```go
for _, entry := range entries {                          // depth 1
    ast.Inspect(f, func(n ast.Node) bool {               // depth 2 (func literal)
        ts, ok := n.(*ast.TypeSpec)                      // depth 2
        if !ok { return true }                           // depth 3
        st2, ok := ts.Type.(*ast.StructType)
        if !ok { return false }
        for _, field := range st2.Fields.List {          // depth 3
            for _, name := range field.Names {           // depth 4
            }
        }
        return false
    })
}
```

depth 8 → 4. 절반으로 줄지만 Q1 ≤ 2는 여전히 위반.

### visitor 패턴의 근본 한계

일반 트리 재귀가 depth 1로 동작하는 이유: 모든 노드가 **동일 타입**이고 **자기 완결적**. `handle(node)` 하나로 처리 가능.

Go AST는 **이종 타입 + 부모 컨텍스트 의존**:
- `*ast.Field`를 처리할 때 부모가 Request struct인지 interface method인지 알아야 처리가 달라짐
- `*ast.StructType`을 처리할 때 부모 TypeSpec.Name이 "XxxRequest"인지 확인 필요

handle을 타입별로 분리하면 handle 간 상태 공유가 필요:

```go
type collector struct {
    inRequestStruct bool    // 부모 컨텍스트 추적
    currentName     string  // 상위 레벨 정보
    fields          map[string]string
}
```

명시적 중첩(원본)이 **암묵적 상태 머신**(visitor)으로 바뀔 뿐, 복잡도가 감소하지 않고 이동한다.

추가 문제:
- **에러 전파 불가**: `ast.Inspect` 콜백은 bool만 반환. `return fmt.Errorf(...)` 불가
- **불필요 순회**: Inspect는 함수 본문, 표현식 등 무관한 노드까지 전부 방문
- **부모 컨텍스트 소실**: 콜백이 `*ast.TypeSpec`을 받을 때 어떤 GenDecl 소속인지 모름

### 결론

visitor 패턴은 이 케이스에서 해결책이 아니다. depth를 줄여주지만 Q1 ≤ 2는 못 맞추고, 에러 전파·부모 컨텍스트·불필요 순회라는 새로운 복잡도가 발생한다.

## 제안: `traverse` 어노테이션

### 설계 원칙

`control=`은 Böhm-Jacopini 3분류(sequence/selection/iteration)를 유지한다. `traverse`는 별도 어노테이션으로 Q1 예외를 명시한다.

**다차원 데이터의 각 차원은 named type으로 모델링되어야 한다.** `[][][]int` 같은 raw 중첩은 타입 단언이 없으므로 `traverse` 대상이 아니다. 각 차원을 named type으로 모델링해야만 자격이 생긴다. 룰이 좋은 데이터 설계를 강제한다.

### 어노테이션 형식

```go
//ff:func feature=symbol type=loader control=iteration traverse
//ff:what 디렉토리에서 Go interface를 파싱하여 "pkg.Model" 키로 등록한다
```

### 차원 수 자동 검출

filefunc는 for-range 체인 안의 **named type 타입 단언 횟수**로 차원 수 N을 자동 산출한다:

```go
gd, ok := decl.(*ast.GenDecl)          // named type 1
ts, ok := spec.(*ast.TypeSpec)         // named type 2
st, ok := ts.Type.(*ast.StructType)    // named type 3
```

외부 패키지 타입 해석 불필요. AST의 타입 단언 표현식에서 named type을 세면 된다.

### 검증 룰

| 룰 | 설명 | severity |
|---|---|---|
| T1 | `traverse` 있지만 for-range 체인에 named type 타입 단언이 2회 미만 | ERROR |
| T2 | `traverse` 있으면 Q1 depth 상한 = 2 + N (N = 타입 단언 체인의 named type 수) | exception |

### Q3 상한

`traverse`는 `control=iteration`의 부분집합이므로 Q3 상한은 iteration과 동일(100줄). 다만 차원이 깊으면 줄 수도 자연히 늘어나므로, 별도 상한(예: 200줄) 검토 여지 있음.

### 적용 대상

Go AST에 한정되지 않는다. **각 차원이 named type으로 모델링된 다차원 데이터 구조** 전체가 대상:
- Go AST: `File → GenDecl → TypeSpec → StructType → Field`
- OpenAPI spec: `PathItem → Operation → Parameter → Schema`
- protobuf descriptor: `FileDescriptor → MessageDescriptor → FieldDescriptor`

`[][][]int` 같은 raw 중첩 배열은 named type 단언이 없으므로 대상 아님. named type으로 래핑해야 자격 획득:

```go
// 대상 아님: raw 중첩
var grid [][][]int

// 대상: named type 모델링
type Grid struct { Layers []Layer }
type Layer struct { Rows []Row }
type Row struct { Cells []int }
```
