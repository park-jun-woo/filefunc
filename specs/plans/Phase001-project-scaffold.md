# Phase 001: 프로젝트 스캐폴드 ✅ 완료

## 목표
Go 모듈 초기화, cobra CLI 뼈대, 디렉토리 구조 생성. `filefunc validate` 명령이 빈 껍데기로 실행 가능한 상태.

## 산출물

| 파일 | 내용 |
|---|---|
| `go.mod` | 모듈 초기화 |
| `cmd/filefunc/main.go` | 엔트리포인트 |
| `internal/cli/root.go` | rootCmd (cobra) |
| `internal/cli/validate.go` | validateCmd (빈 껍데기) |

## 완료 기준
- `go build ./cmd/filefunc/` 성공
- `filefunc validate` 실행 시 "not implemented" 메시지 출력
- filefunc 룰 자기준수 확인 (F1, F4, N1~N5)
