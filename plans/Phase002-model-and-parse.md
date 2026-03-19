# Phase 002: 모델 및 파싱 ✅ 완료

## 목표
핵심 타입 정의와 파싱 기능 구현. Go 소스 파일에서 func/type/method 추출, 어노테이션 파싱, codebook.yaml 로드.

## 산출물

| 파일 | 내용 |
|---|---|
| `internal/model/annotation.go` | Annotation 타입 |
| `internal/model/codebook.go` | Codebook 타입 |
| `internal/model/violation.go` | Violation 타입 |
| `internal/model/go_file.go` | GoFile 타입 |
| `internal/parse/parse_go_file.go` | ParseGoFile — go/ast로 func/type/method 추출 |
| `internal/parse/parse_annotation.go` | ParseAnnotation — //ff: 주석 파싱 |
| `internal/parse/parse_codebook.go` | ParseCodebook — codebook.yaml 로드 |
| `internal/walk/walk_go_files.go` | WalkGoFiles — 디렉토리 내 .go 파일 순회 |

## 의존성
- `go/parser`, `go/ast`, `go/token` (표준 라이브러리)
- `gopkg.in/yaml.v3`

## 완료 기준
- ParseGoFile: .go 파일에서 func 이름, type 이름, method 이름, 라인 수 추출
- ParseAnnotation: //ff:func, //ff:what, //ff:why, //ff:calls, //ff:uses 파싱
- ParseCodebook: codebook.yaml → Codebook 구조체 로드
- WalkGoFiles: 디렉토리 재귀 순회, _test.go 구분
- 각 파일 테스트 통과
- filefunc 룰 자기준수 확인 (F1, F2, N1~N5)
