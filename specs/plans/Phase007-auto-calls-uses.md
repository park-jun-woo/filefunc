# Phase 007: calls/uses 자동 산출 ✅ 완료

## 목표
Go AST에서 `//ff:calls`와 `//ff:uses`를 자동 추출하여 파일에 기록한다. `filefunc annotate` 명령으로 실행.

## 범위
- **calls**: 같은 프로젝트 내 func 호출만. 표준 라이브러리, built-in(`make`, `len`, `append`) 제외.
- **uses**: 같은 프로젝트 내 type 사용만. built-in 타입(`string`, `int`, `error` 등) 제외.

## 판별 기준: 프로젝트 내부 여부
- 같은 패키지 내 호출: 프로젝트 전체 func/type 이름 목록과 대조하여 built-in과 구분
- 다른 패키지 호출: import 경로가 `go.mod` 모듈 경로로 시작하면 프로젝트 내부

## 산출물

| 파일 | 내용 |
|---|---|
| `internal/parse/collect_project_symbols.go` | CollectProjectSymbols — 프로젝트 전체 func/type 이름 수집 |
| `internal/parse/extract_calls.go` | ExtractCalls — func body에서 프로젝트 내 호출 함수명 추출 |
| `internal/parse/extract_uses.go` | ExtractUses — func에서 프로젝트 내 사용 타입명 추출 |
| `internal/parse/read_module_path.go` | ReadModulePath — go.mod에서 모듈 경로 추출 |
| `internal/annotate/write_annotation_line.go` | WriteAnnotationLine — 파일의 //ff: 라인을 추가/갱신 |
| `internal/cli/annotate.go` | annotateCmd — cobra 서브커맨드 |

## 실행 흐름

```
filefunc annotate ./internal/
  │
  ├─ ReadModulePath(go.mod) → 모듈 경로
  │
  ├─ WalkGoFiles → ParseGoFile 전체 → CollectProjectSymbols → 프로젝트 func/type 목록
  │
  ├─ 각 파일에 대해:
  │    ├─ ExtractCalls(파일, 모듈경로, 프로젝트심볼) → []string
  │    ├─ ExtractUses(파일, 모듈경로, 프로젝트심볼) → []string
  │    └─ WriteAnnotationLine(파일, "calls", 값)
  │       WriteAnnotationLine(파일, "uses", 값)
  │
  └─ 변경된 파일 수 출력
```

## WriteAnnotationLine 동작
- 기존 `//ff:calls` 라인이 있으면 교체
- 없으면 `//ff:what` 또는 `//ff:why` 다음, `package` 선언 전에 삽입
- calls/uses가 비어있으면 해당 라인 생략 (기존 라인 있으면 제거)
- 파일 전체 읽기 → 라인 수정 → 전체 다시 쓰기

## 완료 기준
- `filefunc annotate`로 calls/uses 자동 생성
- 기존 //ff:calls, //ff:uses가 있으면 갱신
- 프로젝트 외부 호출/타입은 제외
- 빈 calls/uses는 라인 생략
- filefunc 자체 코드에 annotate 실행하여 검증
- filefunc 자체 코드가 전체 validate 통과
- 테스트 추가
