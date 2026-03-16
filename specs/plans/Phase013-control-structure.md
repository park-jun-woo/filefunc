# Phase 013: control structure 어노테이션 및 검증 ✅ 완료

## 목표
`control=sequence|selection|iteration` 어노테이션을 도입하고, validate가 AST로 정합성을 검증한다. Böhm-Jacopini 정리(1966)의 3대 제어구조.

## 학술 근거
Böhm-Jacopini 정리(1966): 모든 프로그램은 sequence, selection, iteration 세 가지 제어구조의 조합으로 표현 가능. 분류가 수학적으로 완전하고 상호 배타적이므로 연구 없이 즉시 적용 가능.

## 설계 결정 (확정)

| 항목 | 결정 | 근거 |
|---|---|---|
| 분류 | sequence, selection, iteration 3개 | Böhm-Jacopini 정리. 완전하고 상호 배타적 |
| 어노테이션 | `control=` 키 사용. AI가 달고 AI가 읽음 | 제1시민이 AI. body 읽기 전 제어구조 예측 |
| 검증 | AST로 정합성 검증 | 거짓 분류 방지 |
| 1 func 1 control | sequence 함수에 switch/loop이 하나라도 있으면 ERROR | 제어구조 혼합 금지. 분리 강제 |
| Q1 | 전 control 동일: depth ≤ 2 | control별 Q1 분기 불필요 |
| Q3 sequence | 100줄 WARNING | fullend 데이터: sequence 100줄 초과는 대부분 분할 대상 |
| Q3 iteration | 100줄 WARNING | fullend 데이터: iteration 100줄 초과는 분할 대상 |
| Q3 selection | 300줄 WARNING | switch case 나열은 줄 수가 길어도 복잡도 낮음 |
| early return if | sequence에서 허용 | 분기가 아니라 탈출 |
| 대상 범위 | Go 응용 레이어 (백엔드, CLI, 코드 생성기, SSOT 검증기) | 알고리즘 라이브러리, 저수준 시스템은 비대상 |

## 어노테이션

```go
//ff:func feature=validate type=rule control=selection
//ff:what 타입별 필수 필드 누락 검증
func validateRequiredFields(...) ...
```

- `control=sequence` — 기본값. 생략 가능. depth 1에 switch/loop 없음.
- `control=selection` — depth 1에 switch 존재.
- `control=iteration` — depth 1에 loop (for, range) 존재.

## 검증 룰

| # | 룰 | 위반 시 | 검증 방법 |
|---|---|---|---|
| A9 | `control=selection`인데 depth 1에 switch 없음 | ERROR | AST: FuncDecl body 직계 자식에 SwitchStmt/TypeSwitchStmt 유무 |
| A10 | `control=iteration`인데 depth 1에 loop 없음 | ERROR | AST: FuncDecl body 직계 자식에 ForStmt/RangeStmt 유무 |
| A11 | `control=sequence`인데 depth 1에 switch 또는 loop이 하나라도 존재 | ERROR | AST: FuncDecl body 직계 자식에 SwitchStmt/TypeSwitchStmt/ForStmt/RangeStmt 유무 |

A9~A11은 `//ff:func` 파일(func이 있는 파일)에만 적용. `//ff:type`, var+init 파일은 대상 아님.

## Q3 control별 기준

| control | Q3 기준 |
|---|---|
| sequence | 100줄 WARNING |
| iteration | 100줄 WARNING |
| selection | 300줄 WARNING |

### 백틱 힌트 조건
Q3 WARNING 발생 시 (100줄 또는 300줄 초과), 해당 func에 백틱 문자열(`*ast.BasicLit` token.STRING, 값이 \`로 시작)이 **존재하면** 힌트 추가. 백틱 줄 수 계산 불필요 — 유무만 판별.

```
[WARNING] Q3: func GenerateQueryOpts is 216 lines; recommended maximum is 100
  hint: backtick string detected — consider extracting to a separate template file
```

## 산출물

| 파일 | 내용 | 구분 |
|---|---|---|
| `internal/validate/check_control_selection.go` | A9: control=selection인데 switch 없으면 ERROR | 신규 |
| `internal/validate/check_control_iteration.go` | A10: control=iteration인데 loop 없으면 ERROR | 신규 |
| `internal/validate/check_control_sequence.go` | A11: control=sequence인데 switch/loop 존재하면 ERROR | 신규 |
| `internal/parse/detect_control.go` | DetectControl — func body의 지배적 제어구조 판별 | 신규 |
| `internal/validate/check_func_lines.go` | control별 Q3 기준 분기 + 백틱 힌트 | 수정 |
| `internal/validate/run_all.go` | A9~A11 호출 추가 | 수정 |
| `codebook.yaml` | optional에 control 키 추가 | 수정 |

## DetectControl 판별 규칙

```
func body의 직계 자식(depth 1) statement를 순회:
  - SwitchStmt 또는 TypeSwitchStmt가 있음 → selection
  - ForStmt 또는 RangeStmt가 있음 → iteration
  - 둘 다 없음 → sequence
```

## AI 컨텍스트 전략

control로 body를 읽기 전에 "어떻게 읽을지"를 LLM 없이 기계적으로 결정:

```
control=selection  → 전체를 한 번에 read. case를 부분적으로 읽으면 안 됨.
control=iteration  → body 핵심은 루프 내부. 루프 외부는 초기화.
control=sequence   → 수정 대상 step만 read. 나머지 step은 what으로 충분.
```

## 문서 업데이트

| 파일 | 변경 내용 |
|---|---|
| `CLAUDE.md` | A9~A11 룰 추가. Q3 control별 기준 명시. 어노테이션 포맷에 control= 추가 |
| `README.md` | Annotation 룰 테이블에 A9~A11 추가. Annotations 섹션에 control= 예시. Q3 기준 control별 분기 설명 |
| `artifacts/manual-for-ai.md` | Rules 테이블에 A9~A11 추가. Annotations 테이블에 control 행 추가. 컨텍스트 전략에 control별 read 방법 추가 |

## 완료 기준
- A9: `control=selection`인데 switch 없으면 ERROR
- A10: `control=iteration`인데 loop 없으면 ERROR
- A11: `control=sequence`인데 switch/loop 존재하면 ERROR
- Q3: sequence/iteration 100줄, selection 300줄
- Q3 백틱 힌트 동작
- filefunc 자체 코드에 control 어노테이션 부착 (iteration ~42개, selection ~5개, sequence 나머지)
- filefunc 자체 코드 validate 통과
- codebook.yaml에 control 키 추가
- CLAUDE.md 업데이트
- README.md 업데이트
- artifacts/manual-for-ai.md 업데이트
