# func 분류 — logic 키 설계

## 동기

filefunc Q1(네스팅 ≤ 2), Q3(함수 ≤ 100줄)은 일반 로직에는 적절하지만, 특정 구조 패턴의 함수에는 과도하다.

- **dispatch**: switch 10개 case, 각 3~8줄. 119줄이지만 실질 복잡도는 낮음. 10개 파일로 쪼개면 전체 규칙을 한눈에 볼 수 없어 가독성 저하.
- **template**: 백틱 코드젠 216줄 중 212줄이 문자열 리터럴. 실제 Go 로직 4줄. 쪼갤 게 없음.
- **pipeline**: GenWith 151줄. 12개 gen 단계를 순차 호출. 각 호출이 독립적이지만 공유 변수(report, parsed)에 의존. 분리하면 인자 폭발.
- **walker**: AST 순회 157줄. for→for→if→for 패턴이 AST 구조에서 불가피. depth 3~4가 자연스러움.

현재 해결 방법: 없음. Q3는 WARNING이라 무시하거나, 억지로 쪼개서 오히려 코드 품질을 해침.

## 설계

### codebook.yaml optional 키

```yaml
optional:
  logic:
    - pure        # 순수 로직 (기본값, 생략 가능)
    - dispatch    # switch/case 테이블형 분기
    - template    # 백틱 코드젠 (문자열 리터럴이 주 본문)
    - pipeline    # 순차 호출 체인 (오케스트레이션)
    - walker      # AST/파일/트리 순회
```

### 어노테이션 사용

```go
//ff:func feature=ssac-validate type=rule logic=dispatch
//ff:what 타입별 필수 필드 누락 검증
func validateRequiredFields(sf parser.ServiceFunc) []ValidationError {
```

```go
//ff:func feature=gen-gogin type=generator logic=template
//ff:what creates model/queryopts.go with parseQueryOpts
func generateQueryOpts(modelDir string) error {
```

`logic=pure`는 기본값이므로 생략 가능. 명시하지 않으면 pure로 간주.

### logic별 품질 기준

| logic | Q1 기준 | Q3 기준 | 근거 |
|---|---|---|---|
| `pure` | depth ≤ 2 | ≤ 100줄 | 일반 로직. 엄격 적용 |
| `dispatch` | case 내부 depth ≤ 2 | case 수 × 10줄까지 허용 | switch 자체는 네스팅이 아닌 분기. 전체를 한눈에 보는 게 중요 |
| `template` | 템플릿 외부만 depth ≤ 2 | 템플릿 줄 제외 후 ≤ 100줄 | 문자열 리터럴은 로직이 아님. 줄 수에 포함하면 불합리 |
| `pipeline` | depth ≤ 2 | ≤ 200줄 | 순차 호출 나열이므로 줄 수가 길어도 복잡도는 낮음 |
| `walker` | depth ≤ 3 | ≤ 100줄 | 순회 구조상 depth 3이 자연스러움 (for→type-assert→for) |

### AST 검증 (거짓 분류 방지)

filefunc validate가 `logic=` 값의 합당성을 AST로 검증:

| logic | AST 검증 조건 | 위반 시 |
|---|---|---|
| `dispatch` | 함수 본문의 주 구조가 `*ast.SwitchStmt` 또는 `*ast.TypeSwitchStmt`이고, case별 줄 수 ≤ 10 | ERROR: logic=dispatch이지만 switch 구조가 아닙니다 |
| `template` | 백틱 `*ast.BasicLit` (token.STRING)의 줄 수가 함수 본문의 50% 이상 | ERROR: logic=template이지만 템플릿 비율이 50% 미만입니다 |
| `pipeline` | 함수 본문에서 함수 호출(`*ast.CallExpr`) 비율이 70% 이상 (줄 기준) | ERROR: logic=pipeline이지만 순차 호출 비율이 70% 미만입니다 |
| `walker` | 함수 본문에 `*ast.RangeStmt` 또는 `*ast.ForStmt`가 존재하고, 내부에 type assertion이 있음 | ERROR: logic=walker이지만 순회 패턴이 감지되지 않습니다 |

### AI 에이전트 활용

1. **탐색 시**: `grep 'logic=dispatch'` → 테이블형 함수, 전체를 한눈에 읽어야 함
2. **수정 시**: `logic=template` 함수는 템플릿 내용만 수정, Go 로직은 건드리지 않음
3. **리팩토링 시**: `logic=pure`이면서 Q3 위반 → 분리 대상. `logic=pipeline`이면서 Q3 위반 → 200줄 이내면 허용
4. **코드 리뷰 시**: `logic=` 값이 없으면 → 반드시 100줄 이내여야 함. 초과하면 logic 분류부터 결정

## 적용 예시 (fullend 기준)

### dispatch (5개)

```
validateRequiredFields     119줄  case 10개 × 3~8줄
ValidateWith               102줄  switch 11개 kind 디스패처
statusCmd (Status)         128줄  switch kind별 summary
parseDDLTables             108줄  라인별 파서 (for→switch 패턴)
buildScenarioOrder         130줄  분류+정렬 (여러 switch/if 분기)
```

### template (3개)

```
generateQueryOpts          216줄  백틱 212줄 / 로직 4줄
generateMain               126줄  백틱 54줄 / 로직 72줄
generateMainWithDomains    192줄  백틱 58줄 / 로직 134줄 (pipeline과 혼합)
```

### pipeline (2개)

```
GenWith                    151줄  12개 gen 단계 순차 호출
generateMainWithDomains    192줄  template+pipeline 혼합
```

### walker (1개)

```
loadPackageGoInterfaces    157줄  AST 3단계 순회
```

### pure (나머지)

```
generateMethodFromIface    285줄  → 분리 필수 (switch case별 추출)
generateServerStruct       137줄  → 분리 가능 (pathParam 수집)
generateCentralServer      137줄  → 분리 가능 (route 빌더)
transformSource            103줄  → 분리 가능 (2블록 추출)
generateModelFile          104줄  → 분리 가능 (import 감지)
```

## 미결정 사항

1. **혼합 분류**: `generateMainWithDomains`는 template+pipeline. 복수 logic 허용? (`logic=template,pipeline`) 아니면 주 패턴 하나만?
2. **case별 줄 수 상한**: dispatch에서 case당 10줄? 15줄? 20줄?
3. **pipeline 호출 비율**: 70%가 적절한지, 검증 방법이 실용적인지
4. **walker depth 상한**: depth 3이면 대부분 커버되지만 AST 깊은 순회는 depth 4~5 필요할 수 있음
