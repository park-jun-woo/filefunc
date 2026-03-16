# Phase 016: chain --meta 플래그 ✅ 완료

## 목표

`filefunc chain` 출력에 각 함수/타입의 어노테이션 메타데이터를 선택적으로 포함하여, 본문을 읽지 않고 호출 관계 + 함수 목적 + 설계 의도를 파악할 수 있게 한다.

## 배경

현재 `filefunc chain func X --chon 2` 출력은 함수명과 촌수만 표시:

```
RunAll
  1촌 child: CheckNestingDepth
  1촌 child: CheckDimensionRequired
  2촌 sibling: CheckControlSelection
```

`--meta` 플래그로 어노테이션을 함께 출력하면 LLM 컨텍스트로 직접 사용 가능:

```
RunAll (feature=validate type=command control=iteration dimension=1 what="모든 검증 룰을 실행하고 위반 목록을 반환")
  1촌 child: CheckNestingDepth (feature=validate type=rule control=sequence what="Q1: control과 dimension 기반으로 nesting depth 상한 검증")
  1촌 child: CheckDimensionRequired (feature=validate type=rule control=sequence what="A15: control=iteration이면 dimension= 필수")
```

## 설계

### --meta 플래그

```bash
filefunc chain func RunAll --chon 2 --meta what            # what만
filefunc chain func RunAll --chon 2 --meta meta,what       # //ff:func 키-값 + what
filefunc chain func RunAll --chon 2 --meta meta,what,why   # 키-값 + what + why
filefunc chain func RunAll --chon 2 --meta all             # 전부 (meta,what,why,checked)
```

`all`은 `meta,what,why,checked`의 축약.

| 값 | 출력 내용 |
|---|---|
| `meta` | `//ff:func` 또는 `//ff:type`의 키-값 쌍 (파일 종류에 따라 자동 선택) |
| `what` | `//ff:what` 내용 |
| `why` | `//ff:why` 내용 (없으면 생략) |
| `checked` | `//ff:checked` 내용 (없으면 생략) |

- 기본값: `--meta` 미지정 시 현재 동작 유지 (함수명 + 촌수만)
- 쉼표 구분으로 복수 지정
- 순서: meta → what → why → checked (지정 순서 무관, 항상 이 순서로 출력)

### 함수명 → 어노테이션 매핑

현재 chain은 함수명만 추적하고 파일 경로/어노테이션을 모른다. 매핑이 필요:

1. `BuildGraph` 시점에 이미 `WalkGoFiles` → `ParseGoFile`로 전체 파일을 순회
2. `BuildGraph` 호출부(`internal/cli/build_graph.go`)에서 `[]*model.GoFile`을 이미 보유
3. 함수명 → GoFile 매핑 (`map[string]*model.GoFile`)을 구축하여 FormatChain에 전달

### 출력 포맷

모든 메타데이터를 한 줄 괄호 안에 key=value 형태로 출력. 공백 포함 값은 큰따옴표로 감싼다.

```
RunAll (feature=validate type=command control=iteration dimension=1 what="모든 검증 룰을 실행하고 위반 목록을 반환")
  1촌 child: CheckNestingDepth (feature=validate type=rule control=sequence what="Q1: control과 dimension 기반으로 nesting depth 상한 검증")
  1촌 child: CheckDimensionRequired (feature=validate type=rule control=sequence what="A15: control=iteration이면 dimension= 필수")
```

why, checked 예시:

```
BuildPrompt (feature=llm type=util control=sequence what="LLM 검증용 프롬프트 생성" why="ollama 호환 단일 프롬프트 형식" checked="llm=gpt-oss:20b hash=5f7150eb")
```

- why, checked는 값이 없으면 생략
- 출력 순서: meta 키-값 → what → why → checked (항상 고정)

## 산출물

| 파일 | 내용 | 구분 |
|---|---|---|
| `internal/chain/format_chain.go` | meta 옵션에 따라 어노테이션 포함 출력 | 수정 |
| `internal/chain/func_file_map.go` | 함수명/메서드명/타입명 → GoFile 매핑 구축 | 신규 |
| `internal/cli/chain_func.go` | `--meta` 플래그 추가, `_` → `files` 변경, 매핑 전달 | 수정 |
| `internal/cli/chain_feature.go` | `--meta` 플래그 추가, 매핑 전달 | 수정 |

## 문서 업데이트

| 파일 | 변경 내용 |
|---|---|
| `CLAUDE.md` | chain 명령 설명에 `--meta` 플래그 추가 |
| `README.md` | chain 섹션에 `--meta` 플래그 및 사용 예시 추가 |
| `artifacts/manual-for-ai.md` | Commands 섹션에 `--meta` 플래그 및 사용 예시 추가 |

## 구현 순서

1. `func_file_map.go` 구현 — Funcs/Methods/Types → GoFile 매핑
2. `format_chain.go` 수정 — meta 옵션 파라미터 추가, 어노테이션 출력
3. `chain_func.go` 수정 — `--meta` 플래그 파싱, `_` → `files` 변경, 매핑 구축, FormatChain에 전달
4. `chain_feature.go` 수정 — `--meta` 플래그 파싱, 매핑 전달
5. 문서 업데이트
6. `filefunc validate` 위반 0 확인

## 완료 기준

- `--meta what` 지정 시 각 함수 옆에 what 출력
- `--meta meta,what` 지정 시 //ff:func 키-값 + what 출력
- `--meta all` 지정 시 meta,what,why,checked 전부 출력
- `--meta` 미지정 시 기존 동작 유지
- why, checked는 값이 있을 때만 출력
- func/type 파일 모두 대응 (meta가 //ff:func 또는 //ff:type 자동 선택)
- filefunc 자체 코드 validate 통과
- CLAUDE.md, README.md, manual-for-ai.md 업데이트
