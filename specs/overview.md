# filefunc 구조 설계

## Go 모듈

```
module github.com/park-jun-woo/filefunc
go 1.22
```

## 패키지 레이아웃

```
cmd/filefunc/
  main.go                         # 엔트리포인트

internal/
  cli/                            # cobra 명령 정의
    root.go                       # rootCmd
    validate.go                   # validateCmd
    annotate.go                   # annotateCmd (Phase 2)
    chain.go                      # chainCmd (Phase 3)

  parse/                          # 파싱
    parse_annotation.go           # ParseAnnotation — //ff: 주석 파싱
    parse_codebook.go             # ParseCodebook — codebook.yaml 로드
    parse_go_file.go              # ParseGoFile — Go 소스 파일 구조 파싱

  validate/                       # 검증 룰 (각 파일 = 각 룰)
    check_one_file_one_func.go    # F1: 파일당 func 1개
    check_one_file_one_type.go    # F2: 파일당 type 1개
    check_one_file_one_method.go  # F3: 파일당 method 1개
    check_init_only_main.go       # F4: init()은 main.go에만
    check_nesting_depth.go        # Q1: nesting depth ≤ 2
    check_func_lines.go           # Q2/Q3: func line count
    check_annotation_required.go  # A1: //ff:func 필수
    check_codebook_values.go      # A2: 코드북 값 검증
    check_what_required.go        # A3: //ff:what 필수
    check_annotation_position.go  # A6: 어노테이션 최상단 위치
    run_all.go                    # RunAll — 전체 룰 실행 오케스트레이터

  model/                          # 데이터 구조체
    annotation.go                 # Annotation
    codebook.go                   # Codebook
    violation.go                  # Violation
    go_file.go                    # GoFile

  report/                         # 검증 결과 출력
    format_text.go                # FormatText
    format_json.go                # FormatJSON

  walk/                           # 파일 순회
    walk_go_files.go              # WalkGoFiles — 디렉토리 내 .go 파일 순회
```

## 핵심 타입

```go
// annotation.go
type Annotation struct {
    Func    map[string]string // feature=validate type=rule ...
    What    string            // 1줄 설명
    Why     string            // 설계 의도 (선택)
    Calls   []string          // 호출하는 함수 목록
    Uses    []string          // 사용하는 타입 목록
}

// codebook.go
type Codebook struct {
    Feature []string
    Type    []string
    Pattern []string
    Level   []string
}

// violation.go
type Violation struct {
    File    string // 위반 파일 경로
    Rule    string // 룰 ID (F1, Q1, A1 등)
    Level   string // ERROR / WARNING
    Message string // 설명
}

// go_file.go
type GoFile struct {
    Path        string
    Package     string
    Funcs       []string     // 파일 내 func 이름 목록
    Types       []string     // 파일 내 type 이름 목록
    Methods     []string     // 파일 내 method 이름 목록
    Annotation  *Annotation  // 파싱된 어노테이션 (없으면 nil)
    Lines       int          // 총 라인 수
    MaxDepth    int          // 최대 nesting depth
}
```

## 실행 흐름 (validate)

```
filefunc validate ./target/
  │
  ├─ ParseCodebook(codebook.yaml)     → Codebook
  │
  ├─ WalkGoFiles(./target/)           → []string (파일 경로 목록)
  │
  ├─ 각 파일에 대해:
  │    ├─ ParseGoFile(path)           → GoFile
  │    ├─ ParseAnnotation(path)       → Annotation
  │    └─ 각 룰 실행                   → []Violation
  │
  ├─ RunAll(files, codebook)          → []Violation
  │
  └─ FormatText(violations)           → stdout 출력
```

## 구현 단계

### Phase 1: validate (MVP)

파일 구조 룰(F1~F4)과 코드 품질 룰(Q1~Q3)을 검증한다. Go AST(`go/parser`, `go/ast`)로 구현. tree-sitter는 Phase 1에서 불필요 — Go 표준 라이브러리만으로 충분하다.

| 우선순위 | 대상 |
|---|---|
| 1 | cli (root, validate), walk, model |
| 2 | parse (go_file, annotation, codebook) |
| 3 | validate 룰 (F1~F4, Q1~Q3) |
| 4 | validate 룰 (A1~A3, A6) |
| 5 | report (text, json) |

### Phase 2: annotate

LLM이 Go 소스를 읽고 `//ff:func`, `//ff:what`, `//ff:calls`, `//ff:uses`를 자동 생성한다.

### Phase 3: chain

`//ff:calls`와 func signature의 input/output 타입을 기반으로 데이터 흐름 그래프를 구성한다.

## 의존성

Phase 1 기준:

| 의존성 | 용도 |
|---|---|
| `github.com/spf13/cobra` | CLI 프레임워크 |
| `go/parser`, `go/ast`, `go/token` | Go 소스 파싱 (표준 라이브러리) |
| `gopkg.in/yaml.v3` | codebook.yaml 파싱 |

## codebook.yaml 위치

- 프로젝트 루트의 `codebook.yaml`을 기본으로 탐색
- `--codebook` 플래그로 경로 지정 가능
