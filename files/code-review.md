# filefunc 소스 품질 평가

Phase 001~010 완료 후 전체 소스 코드 리뷰.

## 수치

| 지표 | 값 |
|---|---|
| 프로덕션 파일 수 | 76개 |
| 총 라인 수 | 2,348줄 |
| 파일당 평균 | 31줄 |
| 최대 라인 (추정) | ~70줄 (validate.go, llmc.go의 RunE 클로저) |
| nesting depth 위반 | 0 |
| F1 위반 | 0 |

## 긍정적

**1. 가독성이 높다.** 각 파일을 열면 어노테이션 3줄 + func 1개. 뭘 하는지 즉시 파악된다. `SplitTrim`, `SortedKeys`, `IsSkippableLine` 같은 10줄짜리 파일이 과하다고 느낄 수 있지만, read 한 번에 정확히 하나만 들어온다.

**2. early return이 일관적이다.** 모든 func에서 에러/nil 체크 후 즉시 반환. 중첩이 없다.

**3. 의존 관계가 명시적이다.** `//ff:calls`, `//ff:uses`로 파일 간 관계가 어노테이션에 보인다. body를 안 읽어도 `ParseGoFile`이 `CollectFuncDecl`, `CollectGenDecl`, `CalcMaxDepth`, `ParseAnnotation`을 호출한다는 걸 안다.

**4. 네이밍이 균일하다.** `check_one_file_one_func.go` → `CheckOneFileOneFunc`. 파일명으로 함수를 예측할 수 있다.

## 개선 필요

**1. `//ff:uses` 중복 라인**

여러 파일에서 `//ff:uses` 라인이 중복 출력되어 있다:
```
//ff:uses GoFile, Violation
//ff:uses GoFile, Violation
```
`annotate`의 WriteAnnotationLine 버그가 완전히 해결되지 않았다. 기존 라인과 새 라인의 값이 같을 때 중복이 남는 경로가 있다.

**2. ollama.go에 unexported type 3개**

`ollamaRequest`, `ollamaResponse`, `OllamaProvider`가 같은 파일에 있다. F2(파일당 type 1개) 위반은 아닌가? `OllamaProvider`만 exported이고 나머지는 unexported이므로 F6 예외에 해당하긴 하지만, 엄밀히는 F6는 "함수 전용 파라미터 타입"에 대한 예외다. request/response 구조체는 파라미터 타입이 아니라 내부 구현 타입이다.

**3. validate.go, llmc.go의 RunE 클로저가 길다**

cobra 명령의 RunE 클로저가 50줄 이상이다. func으로 분리되어 있지만(ProcessLlmcFile 등), 오케스트레이션 로직 자체가 긴 편이다. Q3 위반(100줄 권고)은 아니지만 filefunc 철학에서 보면 아슬아슬하다.

**4. model_exists.go에도 unexported type**

`ollamaTagsResponse`가 `ModelExists` func과 같은 파일에 있다. ollama.go와 같은 이슈.

## 결론: 코드 품질이 향상되었는가?

**그렇다.** 일반적인 Go 프로젝트와 비교하면:

- 파일당 31줄 평균은 극도로 작다. 보통 Go 프로젝트는 파일당 200~500줄.
- depth 2 제한이 모든 func에서 지켜지고 있다. 일반 Go 코드에서는 depth 4~5가 흔하다.
- 함수가 모두 짧고 단일 책임이다. "이 func이 뭘 하는가"가 명확하다.

**filefunc 룰이 코드 품질 향상을 강제한다**는 것이 자체 코드로 증명되었다. 다만 annotate 중복 버그와 unexported type 처리 기준은 정리가 필요하다.

## 발견된 버그/이슈 목록

| # | 이슈 | 심각도 | 상태 |
|---|---|---|---|
| 1 | `//ff:uses` 중복 라인 (annotate 버그) | 중 | 미수정 |
| 2 | ollama.go unexported type 3개 (F2/F6 경계) | 낮 | 룰 명확화 필요 |
| 3 | model_exists.go unexported type (같은 이슈) | 낮 | 룰 명확화 필요 |
| 4 | RunE 클로저 길이 (50줄+) | 낮 | 모니터링 |
