# Q3 100줄 제한 평가 — fullend 실데이터 분석

## 분석 일시

2026-03-16

## 대상

fullend 프로젝트 `filefunc validate` Q3 WARNING 15건 전수 분석.

## 데이터

| 함수 | 파일 | 총줄 | 로직줄 | 템플릿줄 | 패키지 |
|---|---|---|---|---|---|
| generateMethodFromIface | gen/gogin/generate_method_from_iface.go | 285 | 205 | 0 | gen/gogin |
| generateQueryOpts | gen/gogin/generate_query_opts.go | 216 | 133 | 212 | gen/gogin |
| generateMainWithDomains | gen/gogin/generate_main_with_domains.go | 192 | 135 | 58 | gen/gogin |
| loadPackageGoInterfaces | ssac/validator/load_package_go_interfaces.go | 157 | 111 | 0 | ssac/validator |
| GenWith | orchestrator/gen_with.go | 151 | 98 | 0 | orchestrator |
| generateServerStruct | gen/gogin/generate_server_struct.go | 137 | 97 | 0 | gen/gogin |
| generateCentralServer | gen/gogin/generate_central_server.go | 137 | 96 | 0 | gen/gogin |
| buildScenarioOrder | gen/hurl/build_scenario_order.go | 130 | 88 | 0 | gen/hurl |
| generateMain | gen/gogin/generate_main.go | 126 | 95 | 54 | gen/gogin |
| validateRequiredFields | ssac/validator/validate_required_fields.go | 119 | 79 | 0 | ssac/validator |
| TestParseIdempotency | orchestrator/parse_idempotent_test.go | 117 | 135 | 0 | orchestrator |
| parseDDLTables | ssac/validator/parse_ddl_tables.go | 108 | 62 | 0 | ssac/validator |
| generateModelFile | gen/gogin/generate_model_file.go | 104 | 77 | 0 | gen/gogin |
| transformSource | gen/gogin/transform_source.go | 103 | 63 | 0 | gen/gogin |
| ValidateWith | orchestrator/validate_with.go | 102 | 84 | 0 | orchestrator |

> **로직줄**: 빈줄, 주석줄, 닫는 괄호만 있는 줄 제외.
> **템플릿줄**: 백틱 문자열 내부 줄 수.

## 분류 결과

### 1. 쪼개야 하는 함수 (진짜 복잡) — 4개

| 함수 | 줄 | 로직 | 판단 근거 |
|---|---|---|---|
| generateMethodFromIface | 285 | 205 | switch 7개 case, 각 30~50줄. 독립적인 코드 생성 패턴이 하나의 함수에 밀집. case별 분리가 가독성과 테스트 모두 개선. |
| loadPackageGoInterfaces | 157 | 111 | 3단계 AST 순회(Request struct 수집 → interface 파싱 → standalone func 파싱)가 한 함수에 밀집. 각 단계가 독립적 변환. |
| GenWith | 151 | 98 | 12개 gen 단계 순차 호출 + preserve 로직. pipeline 패턴이지만 로직줄 98로 경계선. 블록 추출 시 인자 폭발 우려. |
| generateMainWithDomains | 192 | 135 | domain init + import + queue 3블록 혼재. 템플릿 58줄 포함이지만 로직 자체가 135줄로 복잡. |

### 2. 쪼개면 오히려 나빠지는 함수 — 5개

| 함수 | 줄 | 로직 | 판단 근거 |
|---|---|---|---|
| validateRequiredFields | 119 | 79 | switch 10개 case, 각 3~8줄. flat dispatch — 전체 검증 규칙을 한눈에 보는 게 핵심 가치. 10개 파일로 분리하면 "SeqGet 필수 필드가 뭐지?"에 10번 읽기 필요. |
| generateQueryOpts | 216 | 133 | **212줄이 백틱 템플릿**. 실제 Go 로직 = `src := \`...\`; return os.WriteFile(...)` 4줄. 함수가 긴 게 아니라 생성할 코드가 긴 것. |
| ValidateWith | 102 | 84 | switch 11개 kind 디스패처. validateRequiredFields와 동일 패턴. 겨우 2줄 초과. |
| parseDDLTables | 108 | 62 | 라인별 파서. for→switch 관용 패턴. 실질 로직 62줄은 100줄 이내. |
| buildScenarioOrder | 130 | 88 | 정렬+분류 로직. 3단계가 공유 변수(authSteps, midSteps, readSteps, deleteSteps)에 의존. 분리하면 인자 4~5개 + 반환값 3~4개 필요. |

### 3. 경계선 (분리 가능하나 필수는 아님) — 6개

| 함수 | 줄 | 로직 | 분리 가능 블록 | 분리 효과 |
|---|---|---|---|---|
| generateServerStruct | 137 | 97 | pathParam 수집 (20줄) | 중 — depth 6 해결에도 기여 |
| generateCentralServer | 137 | 96 | route 빌더 (55줄) | 중 |
| generateMain | 126 | 95 | queue 블록 (30줄) | 낮 — 템플릿 54줄 제외 시 72줄 |
| generateModelFile | 104 | 77 | import 감지 (15줄) | 낮 — 15줄 추출로 89줄 |
| transformSource | 103 | 63 | type-assertion 블록 (15줄) | 낮 — 이미 63줄 |
| TestParseIdempotency | 117 | 135 | 테스트 헬퍼 | 낮 — F5 예외 파일 |

## 결론

### 100줄 제한이 적절한 경우 (9/15 = 60%)

순수 로직 함수 + 경계선 함수. 분리가 가독성과 유지보수를 개선하거나 최소한 해치지 않음.

### 100줄 제한이 과도한 경우 (6/15 = 40%)

| 패턴 | 해당 함수 | 이유 |
|---|---|---|
| **flat dispatch** | validateRequiredFields, ValidateWith, parseDDLTables | switch case 테이블형. 한 파일에서 전체를 보는 게 핵심. |
| **template** | generateQueryOpts | 백틱 문자열이 본문 대부분. 실제 로직은 4줄. |
| **shared-state pipeline** | buildScenarioOrder | 공유 변수 의존으로 분리 시 인자 폭발. |

### 제안

1. `logic=dispatch` — switch가 주 구조이고 case별 ≤ 10줄이면 Q3 면제. AST로 자동 검증.
2. `logic=template` — 백틱 줄이 50% 이상이면 템플릿 줄 제외 후 Q3 적용. AST로 자동 검증.
3. `logic=pipeline` — 순차 호출 비율 ≥ 70%이면 Q3 상한 200줄. AST로 자동 검증.
4. `logic=walker` — AST 순회 패턴이면 Q1 depth ≤ 3 허용.
5. `logic=pure` (기본값) — 현행 Q1 ≤ 2, Q3 ≤ 100 유지.

→ `func분류.md`에 상세 설계.
