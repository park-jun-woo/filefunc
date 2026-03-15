# Phase 006: 리포트 및 마무리 ✅ 완료

## 목표
JSON 출력 포맷 추가, CLI 플래그 정리, 자기검증 완료.

## 산출물

| 파일 | 내용 |
|---|---|
| `internal/report/format_json.go` | FormatJSON — JSON 출력 |

## CLI 플래그

| 플래그 | 설명 | 기본값 |
|---|---|---|
| `--codebook` | codebook.yaml 경로 | 프로젝트 루트 자동 탐색 |
| `--format` | 출력 포맷 (text / json) | text |

## 완료 기준
- `filefunc validate --format json` 동작
- `filefunc validate --codebook path/to/codebook.yaml` 동작
- filefunc 자체 코드가 모든 룰(F1~F4, Q1~Q3, A1~A3, A6) 통과
- 전체 테스트 통과
- `go vet ./...` 통과
- README.md 작성 (사용법, 설치법)
