# filefunc

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![skills.sh](https://skills.sh/b/park-jun-woo/filefunc)](https://skills.sh/park-jun-woo/filefunc)

**One file, one concept.** 파일명이 곧 개념명.

**AI 에이전트가 `check_nesting_depth.go`를 열면 함수 하나만 나온다. `utils.go`에서 무관한 헬퍼 19개를 스크롤할 일이 없다.**

Go, Python, TypeScript를 위한 LLM 네이티브 코드 구조 컨벤션 및 CLI 도구. 백엔드 서비스, CLI 도구, 코드 생성기, SSOT 검증기가 대상이다. 알고리즘 라이브러리, 저수준 시스템 프로그래밍은 대상이 아니다.

---

## 적용 사례

### typer — Python CLI 프레임워크 (테스트 1155개, 실패 0)

[fastapi/typer](https://github.com/fastapi/typer)를 `filefunc validate --lang python` 위반 0으로 리팩토링한 [포크](https://github.com/park-jun-woo/typer).

| 지표 | 원본 | 리팩토링 |
|---|---|---|
| 소스 파일 | 16 | 197 |
| filefunc 위반 | 69 | 0 |
| pytest 통과 | 1155 | 1155 |
| pytest 실패 | 0 | 0 |

모든 public API, import 경로, 런타임 동작이 원본과 동일. 성능 저하 없음 (import +2% 오차 범위, 나머지 벤치마크 동일). pytest 전체 통과 및 전수 비교로 검증.

### hono — TypeScript 웹 프레임워크 (테스트 4419개, 신규 실패 0)

[honojs/hono](https://github.com/honojs/hono)를 `filefunc validate --lang typescript` 위반 0으로 리팩토링한 [포크](https://github.com/park-jun-woo/hono).

| 지표 | 원본 | 리팩토링 |
|---|---|---|
| 소스 파일 | 186 | 626 |
| filefunc 위반 | 397 | 0 |
| vitest 통과 | 4419 | 4419 |
| vitest 실패 | 4 | 4 (기존 결함) |

모든 import 경로, 런타임 동작이 원본과 동일. vitest 전체 통과로 검증.

---

## 빠른 시작

```bash
npx skills add park-jun-woo/filefunc
```

```bash
go install github.com/park-jun-woo/filefunc/cmd/filefunc@latest
```

소스에서 빌드:

```bash
git clone https://github.com/park-jun-woo/filefunc.git
cd filefunc
go build ./cmd/filefunc/
```

Go 1.18+ 필요.

---

## AI와 함께 사용하기

위 적용 사례는 AI 에이전트가 리팩토링하고 filefunc이 검증하는 방식으로 측정했다. Claude Code, Codex, Copilot, Cursor — 어떤 에이전트든 된다.

에이전트를 시작하고 프롬프트 하나를 준다:

```
SKILL.md를 읽고 이 프로젝트를 filefunc validate를 통과하도록 리팩토링해.
```

AI가 파일을 분리한다. `filefunc validate`가 구조 위반을 즉시 잡는다. 에이전트는 레일 안에서 자유롭고, 레일을 벗어나면 검증이 바로 실패한다.

---

## 왜

AI 코드 에이전트(Claude Code 등)는 `grep → read`로 코드를 탐색한다. `read`의 단위는 파일이다. 파일 = 개념이면, 한 번의 read가 정확히 하나의 관련된 것을 반환한다.

```
# filefunc 없이
read utils.go → 함수 20개, 19개는 무관. 컨텍스트 오염.

# filefunc 적용
read check_one_file_one_func.go → 함수 1개. 정확히 필요한 것.
```

**무관한 함수 290개를 안 여는 것이, 필요한 5~10개를 고르는 것보다 중요하다.**

filefunc의 1등 시민은 AI 에이전트이지 사람이 아니다. 파일 수 폭발은 기능이지 버그가 아니다 — 파일이 많을수록 파일이 작아지고, read당 노이즈가 줄어든다. 사람의 편의는 뷰 레이어(VSCode 확장 등)에서 해결한다.

---

## 명령

### validate — 코드 구조 룰 검증

```bash
filefunc validate                    # 현재 디렉토리를 프로젝트 루트로
filefunc validate /path/to/project   # 프로젝트 루트 명시
filefunc validate --format json
filefunc validate --lang python      # 언어 지정
```

```
[PASS] 127 files checked, 0 violations
```

프로젝트 루트에 `go.mod`(Go), `pyproject.toml`(Python), `package.json`(TypeScript) 중 하나와 `codebook.yaml` 필수. 읽기 전용. 위반 시 종료 코드 1. `.ffignore` 적용. 언어는 자동 감지되며, `--lang go|python|typescript`로 지정 가능. [toulmin](https://github.com/park-jun-woo/toulmin) 논증 엔진 기반 — 룰은 backing 기반 범용 함수, 예외는 그래프의 defeat로 선언.

### chain — 호출 관계 추적

```bash
filefunc chain func RunAll --chon 2 --meta what
```

```
── Func Chain: RunAll (chon=2) ──

  [self]    RunAll                  //ff:what 전체 검증 실행
  [child]   RunProjectRules         //ff:what P 룰 실행
  [child]   RunCodebookRules        //ff:what C 룰 실행
  [child]   RunFileRules            //ff:what F/Q/A 룰 실행
  [co]      FormatResult            //ff:what 검증 결과 포매팅
  [co]      PrintReport             //ff:what 검증 결과 출력
```

```bash
filefunc chain func RunAll              # 1촌 (기본, 현재 디렉토리)
filefunc chain func RunAll --chon 2     # 2촌 (공동 호출 포함)
filefunc chain func RunAll --chon 3     # 3촌 (최대)
filefunc chain func RunAll --child-depth 3   # 호출만
filefunc chain func RunAll --parent-depth 3  # 호출자만
filefunc chain feature validate         # feature 내 전체 함수
filefunc chain func RunAll --chon 2 --meta what \
  --prompt "nesting depth 수정" --rate 0.8           # 리랭커 필터링
filefunc chain func ParseFile --package funcspec     # 특정 패키지 한정
```

실시간 AST 분석. `.ffignore` 적용.

| 플래그 | 설명 | 기본값 |
|---|---|---|
| `--root` | 프로젝트 루트 | `.` |
| `--chon` | 관계 거리 (1~3) | 1 |
| `--child-depth` | 호출 방향만 추적 (깊이) | — |
| `--parent-depth` | 호출자 방향만 추적 (깊이) | — |
| `--meta` | 어노테이션 메타데이터 포함 (meta,what,why,checked,all) | — |
| `--package` | Go 패키지 한정 (chain func 전용) | — |
| `--prompt` | 관련도 점수를 위한 사용자 작업 의도 (vLLM 필요) | — |
| `--rate` | 관련도 점수 임계값 (0.0~1.0) | 0.8 |
| `--model` | 리랭커 모델명 | `Qwen/Qwen3-Reranker-0.6B` |
| `--score-endpoint` | 리랭커용 vLLM 엔드포인트 | `http://localhost:8000` |

`--prompt`는 Qwen3-Reranker-0.6B를 실행하는 vLLM 서버 필요:

```bash
pip install vllm
vllm serve Qwen/Qwen3-Reranker-0.6B --task score \
  --hf_overrides '{"architectures":["Qwen3ForSequenceClassification"],"classifier_from_token":["no","yes"],"is_original_qwen3_reranker":true}'
```

### context — LLM 컨텍스트 탐색

```bash
filefunc context "nesting depth 검증 수정"                        # 4단계 파이프라인
filefunc context "modify depth logic" --depth 2                    # feature 필터만
filefunc context "depth 수정" --what-rate 0.3                      # what 임계값 조정
filefunc context "depth 수정" --body-rate 0.5                      # body 임계값 조정
filefunc context "depth 수정" --search "feature=validate"          # LLM 건너뛰고 직접 필터
filefunc context "cross 수정" --search "feature=crosscheck ssot=openapi"  # 다중 키 AND
```

4단계 파이프라인: LLM feature 선택 → feature 필터 → what 점수 산출(LLM) → body 점수 산출(LLM). 함수명 불필요. ollama + gpt-oss:20b 필요.

| 플래그 | 설명 | 기본값 |
|---|---|---|
| `--depth` | 파이프라인 깊이 (1-4) | 4 |
| `--what-rate` | what 점수 임계값 | 0.2 |
| `--body-rate` | body 점수 임계값 | 0.5 |
| `--model` | ollama 모델 | `gpt-oss:20b` |
| `--endpoint` | ollama 엔드포인트 | `http://localhost:11434` |
| `--search` | 직접 어노테이션 필터 (LLM feature 선택 건너뜀) | — |

### llmc — LLM 검증

```bash
filefunc llmc                           # 현재 디렉토리
filefunc llmc /path/to/project          # 프로젝트 루트 명시
filefunc llmc --model qwen3:8b
filefunc llmc --threshold 0.9
```

`//ff:what`이 함수 본문과 일치하는지 로컬 LLM(ollama)으로 검증. 점수 0.0~1.0, 임계값 0.8. 통과 시 `//ff:checked` 서명 기록. `.ffignore` 적용.

| 플래그 | 설명 | 기본값 |
|---|---|---|
| `--provider` | LLM 제공자 | `ollama` |
| `--model` | 모델명 | `gpt-oss:20b` |
| `--endpoint` | API 엔드포인트 | `http://localhost:11434` |
| `--threshold` | 최소 통과 점수 | `0.8` |

---

## 룰

전체 룰은 [`rulebook.md`](rulebook.md) 참조 (SSOT). 카테고리: P (프로젝트), I (임포트), F (파일 구조), Q (코드 품질), A (어노테이션), C (코드북), N (네이밍).

검증 순서:

```
1. P 룰 (프로젝트 레벨) — 단일 언어 검사
2. C 룰 (코드북) — codebook.yaml 무결성
3. I 룰 (임포트) — 순환 임포트 탐지 (Python, TypeScript)
4. F/Q/A 룰 (파일 레벨) — toulmin defeats graph
```

P 또는 C 위반 시 후속 검증이 차단된다.

## 어노테이션

```go
//ff:func feature=validate type=rule control=sequence
//ff:what F1: 파일당 func 1개 검증
//ff:why AI 에이전트가 1등 시민. 1 파일 1 개념이 컨텍스트 오염을 방지.
//ff:checked llm=gpt-oss:20b hash=a3f8c1d2     (llmc가 자동 생성)
func CheckOneFileOneFunc(gf *model.GoFile) []model.Violation {
```

`control=`은 모든 func 파일에 필수 (A9). 값: `sequence`, `selection` (switch), `iteration` (loop). Bohm-Jacopini 정리(1966) 기반. 1 func 1 control.

`dimension=`은 `control=iteration` 파일에 필수 (A15). 순회 대상 데이터의 차원 수. Q1 깊이 상한 = dimension + 1. dimension=1이면 flat list (depth <= 2), dimension >= 2이면 named type(struct/interface) 중첩 필수.

| 어노테이션 | 용도 | 필수 |
|---|---|---|
| `//ff:func` | func 메타데이터 (feature, type 등) | 예 (func 파일) |
| `//ff:type` | type 메타데이터 (feature, type 등) | 예 (type 파일) |
| `//ff:what` | 한 줄 설명 — 뭘 하는가 | 예 |
| `//ff:why` | 설계 결정 — 왜 이렇게 만들었는가 | 아니오 |
| `//ff:checked` | LLM 검증 서명 | 자동 (`filefunc llmc`) |

Python은 `# ff:func`, `# ff:what` 등을 사용. TypeScript는 `// ff:func`, `// ff:what` 등을 사용.

## 코드북

코드북은 어노테이션에 사용 가능한 값의 허용 목록이다. 프로젝트의 어휘 — `grep`을 정확하게 만드는 지도.

```yaml
# codebook.yaml
required:
  feature:
    validate: "코드 구조 룰 검증 (F1,Q1,A1 등 정적 분석 룰)"
    parse: "소스 코드, 어노테이션, codebook 파싱"
  type:
    command: "cobra 명령 엔트리포인트"
    rule: "개별 검증 룰 구현"

optional:
  pattern:
    error-collection: "에러 수집 후 일괄 보고"
  level:
    error: ""
    warning: ""
```

각 값에 description 포함 (`key: "설명"`). description은 `filefunc context`의 LLM feature 선택에 사용. `required` 키는 모든 어노테이션에 필수 (A8). `optional` 키는 해당할 때만 사용.

코드북에 없는 값을 쓰면 `A2 ERROR`. 프로젝트마다 고유한 코드북. `codebook.yaml` 필수.

코드북 포맷 룰 (C1-C4)은 [`rulebook.md`](rulebook.md) 참조. 코드북이 먼저 검증된다. 코드북이 실패하면 코드 검증은 실행되지 않는다.

## .ffignore

모든 filefunc 명령에서 경로를 제외. 프로젝트 루트(`go.mod` 옆)에 `.ffignore` 배치. `.gitignore`와 동일 문법. 경로 기반 패턴 지원.

```
# .ffignore 예시
vendor/
*.pb.go
*_gen.go
internal/legacy/
```

선택 사항. 없으면 제외 없음.

---

## 학술 근거

- **"Lost in the Middle" (Stanford, 2024)** — 컨텍스트 중간의 관련 정보는 성능을 30%+ 저하시킨다.
- **"Context Length Alone Hurts LLM Performance" (Amazon, 2025)** — 빈 토큰조차 성능을 저하시킨다 (13.9~85%). 짧고 집중된 컨텍스트가 이긴다.
- **"Context Rot" (Chroma Research)** — 집중된 프롬프트 > 전체 프롬프트, 모든 모델에서.

연구가 "짧은 컨텍스트가 낫다"는 것을 증명했다. filefunc은 코드를 구조적으로 분리하여 관련 부분만 컨텍스트에 들어가게 하는, 그 빠진 도구다.

## 라이선스

MIT — [LICENSE](LICENSE) 참조.
