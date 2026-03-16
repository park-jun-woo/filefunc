# Phase 015: dimension 어노테이션 도입 ✅ 완료

## 목표

순회/분기 대상 데이터의 차원 수를 어노테이션으로 명시하고, Q1 depth 상한을 차원 수 기반으로 동적 적용한다.

## 배경

Q1 ≤ 2는 로직 복잡도 제한에는 효과적이지만, Go AST처럼 다차원 데이터를 순회하는 코드에서 depth가 데이터 구조의 차원을 반영할 뿐 논리 복잡도가 아닌 경우가 있다. 차원 수를 명시하면 Q1 상한을 데이터 구조에 맞게 조정할 수 있다.

상세 분석: `files/GO-AST-예외-보고.md`

## 설계

### Q1 depth 상한 공식

| control | Q1 허용 depth | dimension |
|---|---|---|
| sequence | 2 (고정) | 불필요 |
| selection | 2 (고정) | 불필요 |
| iteration | dimension + 1 | 필수 |

- sequence: 분기/순회 없음. depth 2 (에러 핸들링 등 early return 허용).
- selection: switch + 내부 분기까지 depth 2. 그 이상은 함수 분리 대상.
- iteration: 순회 대상 데이터의 차원 수에 따라 동적. dimension=1이면 depth 2 (현재 기본값과 동일).

### dimension 값 예시 (iteration 전용)

| dimension | 의미 | Q1 허용 depth |
|---|---|---|
| 1 | flat list 순회 | 2 |
| 2 | 2D 데이터 순회 | 3 |
| 3 | 3D 데이터 순회 | 4 |
| 5 | Go AST leaf 순회 | 6 |

### 어노테이션 형식

```go
//ff:func feature=parse type=parser control=iteration dimension=1
//ff:what statement 목록에서 switch/loop 존재 여부로 제어구조 판별

//ff:func feature=symbol type=loader control=iteration dimension=5
//ff:what 디렉토리에서 Go interface를 파싱하여 "pkg.Model" 키로 등록한다

//ff:func feature=parse type=parser control=selection
//ff:what 어노테이션 라인을 키별로 분류하여 Annotation 구조체에 적용
```

- `control=sequence`: dimension 불필요
- `control=selection`: dimension 불필요 (depth 2 고정)
- `control=iteration`: `dimension=` 필수

### 기존 코드 영향

- `control=selection` → 변경 없음 (depth 2 고정, 기존과 동일)
- `control=iteration` → `dimension=1` 추가 (depth ≤ 2, 기존과 동일)
- `control=sequence` → 변경 없음

기존 Q1 동작과 동일하게 유지된다.

## 추가/수정 룰

| # | 룰 | 위반 시 | 검증 방법 |
|---|---|---|---|
| A15 | `control=iteration`이면 `dimension=` 필수 | ERROR | 어노테이션 유무 |
| A16 | `dimension=` 값은 양의 정수여야 함 | ERROR | 값 파싱 |
| Q1 수정 | depth 상한을 동적 적용 (sequence=2, selection=2, iteration=dimension+1) | ERROR | MaxDepth vs 상한 비교 |

### dimension별 입력 데이터 요건 (iteration 전용)

| dimension | 입력 요건 | 비고 |
|---|---|---|
| 1 | Go가 허용하는 모든 iterable (slice, map, string, channel, int, struct field, receiver field, 함수 반환값, 패키지 변수) | 제한 없음 |
| ≥ 2 | named type (struct 또는 interface) 중첩 필수 | raw 중첩 불허 |

```go
// dimension=1: 모든 iterable 허용
func processList(items []string) { ... }
func handleEvents(ch chan Event) { ... }
func (s *Server) run() { for _, h := range s.handlers { ... } }

// dimension=2: named type 중첩 필수
type Grid struct { Layers []Layer }
type Layer struct { Rows []Row }
func processGrid(g Grid) { ... }

// dimension ≥ 2 자격 없음: raw 중첩
var data [][][]int                    // named type 아님
var lookup map[string]map[string]int  // named type 아님

// dimension=5: Go AST (interface + struct 혼합 named type)
// ast.File → ast.GenDecl → ast.TypeSpec → ast.StructType → ast.Field
func loadInterfaces(entries []fs.DirEntry) { ... }
```

named type으로 설계하면:
- 컴파일러 타입 안전성 보장
- gofmt 구조 정리
- json.Marshal 중첩 직렬화 자동 처리
- 연속 메모리 배치로 캐시 효율 최적화
- AI가 struct 생성 비용을 제거하므로 타이핑 부담 없음

### 남용 방지 자동 검증

dimension ≥ 2의 자동 검증(입력 타입의 named type 중첩 수와 dimension 일치 여부)은 구현 복잡도가 높으므로 별도 Phase로 분리한다. Phase 015에서는 어노테이션 필수화와 Q1 동적 상한까지만 구현한다.

## 산출물

| 파일 | 내용 | 구분 |
|---|---|---|
| `internal/validate/check_dimension_required.go` | A15: iteration이면 dimension= 필수 | 신규 |
| `internal/validate/check_dimension_value.go` | A16: dimension 값이 양의 정수인지 검증 | 신규 |
| `internal/validate/check_nesting_depth.go` | Q1: depth 상한을 dimension 기반으로 동적 적용 | 수정 |
| `internal/validate/run_all.go` | A15, A16 호출 추가 | 수정 |
| `internal/validate/check_dimension_required_test.go` | A15 테스트 | 신규 |
| `internal/validate/check_dimension_value_test.go` | A16 테스트 | 신규 |
| `internal/validate/check_nesting_depth_test.go` | Q1 dimension 연동 테스트 | 신규 |
| `internal/validate/testdata/iter_no_dimension.go` | A15 위반: iteration인데 dimension 없음 | 신규 |
| `internal/validate/testdata/bad_dimension_value.go` | A16 위반: dimension=0 | 신규 |
| `internal/validate/testdata/dimension2_depth3.go` | Q1 통과: dimension=2, iteration, depth 3 | 신규 |

### 기존 파일 dimension 추가

기존 `control=iteration` 파일 전체에 `dimension=1` 추가 (약 51개). selection, sequence 파일은 변경 없음.

## 문서 업데이트

| 파일 | 변경 내용 |
|---|---|
| `CLAUDE.md` | A15, A16 룰 추가, Q1 공식 변경, dimension 설명, type struct 설계 안내 |
| `README.md` | Annotation 룰 및 Annotations 섹션에 dimension 추가, type struct 설계 안내 |
| `artifacts/manual-for-ai.md` | Rules 및 Annotations 섹션에 dimension 추가, type struct 설계 안내 |

### 문서에 포함할 설계 안내

dimension은 iteration 전용. selection은 depth 2 고정 (그 이상은 함수 분리).

dimension=1은 Go의 모든 iterable을 허용한다. dimension ≥ 2는 순회 대상 데이터를 named type (struct 또는 interface) 중첩으로 설계해야 한다:

```go
// Bad: raw 중첩 — dimension ≥ 2 자격 없음
var data [][][]int
var lookup map[string]map[string]int

// Good: named type — dimension=2
type Grid struct { Layers []Layer }
type Layer struct { Rows []Row }
type Row struct { Cells []int }

// Good: 외부 named type — dimension=5 (Go AST)
// ast.File → ast.GenDecl → ast.TypeSpec → ast.StructType → ast.Field
```

named type의 이점:
- 컴파일러 타입 검사, gofmt 정리, json.Marshal 자동 직렬화
- 연속 메모리 배치, GC 부담 최소화
- AI가 struct 생성 비용을 제거하므로 타이핑 부담 없음

## 구현 순서

1. sequence 파일 depth 전수 조사 → depth 3 이상인 파일은 control 재분류 또는 함수 분리
2. 기존 iteration 파일에 dimension=1 추가 (약 51개)
3. A15 구현 + 테스트
4. A16 구현 + 테스트
5. Q1 수정 (dimension 연동) + 테스트
6. run_all.go에 A15, A16 등록
7. 문서 업데이트
8. `filefunc validate` 위반 0 확인

## 완료 기준

- A15: iteration인데 dimension 없으면 ERROR
- A16: dimension 값이 양의 정수가 아니면 ERROR
- Q1: sequence → depth ≤ 2, selection → depth ≤ 2, iteration → depth ≤ dimension + 1
- 기존 iteration 파일 전체에 dimension=1 추가 완료
- filefunc 자체 코드 validate 통과
- CLAUDE.md, README.md, manual-for-ai.md 업데이트
