# scripts

Go AST 기반 코드 통계 스크립트. `go run scripts/<name>.go [dir]`로 실행.

| 스크립트 | 측정 대상 | 용도 |
|---------|----------|------|
| `bodystat.go` | 함수 body 라인 수 (TOTAL/PURE) | Q3 규칙 검증 — 제어문 제외 순수 시퀀스 길이 |
| `blockstat.go` | 제어문 body 라인 수 (전체) | 제어문 body 제한 규칙 근거 |
| `blockstat_by_kind.go` | 종류별 분포 (if/for/range/switch) | 종류별 제한 기준 차등화 근거 |
| `patternstat.go` | 함수 depth-1 제어문 패턴 | control= 분류 체계 검증 (mixed 비율 등) |
| `casestat.go` | case 절 내부 라인 수 | case body 제한 규칙 근거 |
