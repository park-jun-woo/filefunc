# Phase 008: LLM 검증 및 //ff:checked 서명 ✅ 완료

## 목표
명령을 분리하여 관심사를 분리한다.
- `filefunc validate`: `//ff:checked`가 있으면 해시 대조. 불일치 시 ERROR. 파일 수정 안 함. checked가 프로젝트에 하나도 없으면 A7 스킵.
- `filefunc llmc`: 소형 LLM(로컬 ollama)으로 what-body 일치 판정 후 `//ff:checked` 서명 기록.

## 설계 원칙
- validate에 LLM 의존성을 넣지 않는다. validate는 어디서든 실행 가능해야 한다.
- validate는 읽기 전용. 파일을 수정하지 않는다.
- checked가 있다는 건 llmc를 돌렸다는 뜻. 서명이 깨졌으면 신뢰성이 깨진 것이므로 ERROR.
- checked가 프로젝트에 하나도 없으면 llmc를 도입하지 않은 것이므로 A7 전체 스킵.

## 확정 사항

| 항목 | 결정 |
|---|---|
| LLM API | 로컬 ollama (OpenAI 호환 엔드포인트) |
| 프롬프트 | 0.0~1.0 사이 점수 반환 |
| 판정 기준 | 0.8 이상이면 일치 (추후 조정 가능) |
| validate 동작 | checked가 프로젝트에 하나도 없으면 A7 스킵 |

## 산출물

| 파일 | 내용 |
|---|---|
| `internal/parse/calc_body_hash.go` | CalcBodyHash — func body SHA-256 해시 (앞 8자리) |
| `internal/validate/check_checked_hash.go` | A7: checked 해시 대조. 불일치 시 ERROR |
| `internal/validate/has_any_checked.go` | HasAnyChecked — 프로젝트에 checked가 하나라도 있는지 판별 |
| `internal/llm/provider.go` | Provider 인터페이스 정의 |
| `internal/llm/ollama.go` | OllamaProvider — ollama /api/generate 호출 |
| `internal/llm/new_provider.go` | NewProvider — provider 문자열로 Provider 생성 |
| `internal/llm/verify_what.go` | VerifyWhat — Provider로 what-body 일치 점수 판정 |
| `internal/llm/build_prompt.go` | BuildPrompt — 검증 프롬프트 생성 |
| `internal/llm/parse_score.go` | ParseScore — LLM 응답에서 0.0~1.0 점수 추출 |
| `internal/cli/llmc.go` | llmcCmd — cobra 서브커맨드 |

## 프롬프트

```
You are a code reviewer. Rate how accurately the description matches the Go function.

Description: {what}

Function:
{body}

Respond with a single number between 0.0 and 1.0 only.
0.0 = completely wrong, 1.0 = perfectly accurate.
```

## validate 흐름 (읽기 전용)

```
filefunc validate ./internal/
  │
  ├─ 기존 룰 (F1~F4, Q1~Q3, A1~A6) 실행
  │
  └─ A7:
       ├─ 프로젝트 전체에 //ff:checked가 하나도 없음 → A7 전체 SKIP
       ├─ //ff:checked가 1개 이상 존재:
       │    ├─ 해당 파일에 //ff:checked 없음 → PASS
       │    ├─ //ff:checked 있고 해시 일치 → PASS
       │    └─ //ff:checked 있고 해시 불일치 → ERROR
```

## llmc 흐름 (파일 수정)

```
filefunc llmc ./internal/
  │
  ├─ 모델 존재 확인 (ollama list)
  │    ├─ 있음 → 진행
  │    └─ 없음 → "Model not found. Pull gpt-oss:20b? [y/N]" 질의
  │         ├─ y → ollama pull 실행 후 진행
  │         └─ N → 중단
  │
  ├─ 각 func 파일에 대해:
  │    ├─ //ff:checked 있고 해시 일치 → SKIP (이미 검증됨)
  │    ├─ //ff:checked 없음 또는 해시 불일치:
  │    │    ├─ LLM에 what + body 전달
  │    │    ├─ 점수 ≥ 0.8 → //ff:checked llm=모델명 hash=body해시 기록
  │    │    └─ 점수 < 0.8 → ERROR 출력 (점수 표시, what 수정 필요)
  │
  └─ 결과 요약 출력 (통과/실패/스킵 수)
```

## CLI 플래그 (llmc)

| 플래그 | 설명 | 기본값 |
|---|---|---|
| `--provider` | LLM 제공자 | `ollama` |
| `--model` | 모델명 | `gpt-oss:20b` |
| `--endpoint` | API 엔드포인트 | `http://localhost:11434` (ollama 기본) |
| `--threshold` | 일치 판정 기준 점수 | `0.8` |

환경변수로도 설정 가능: `FILEFUNC_LLM_PROVIDER`, `FILEFUNC_LLM_MODEL`, `FILEFUNC_LLM_ENDPOINT`

## LLM Provider 인터페이스

```go
// Provider는 LLM API 제공자의 인터페이스.
type Provider interface {
    Generate(prompt string) (string, error)
}
```

현재 구현: `ollama` (POST /api/generate)
향후 확장: `openai`, `anthropic`, `gemini` 등 — Provider 인터페이스 구현만 추가하면 됨.

## 산출물 (Provider 반영)

| 파일 | 내용 |
|---|---|
| `internal/llm/provider.go` | Provider 인터페이스 정의 |
| `internal/llm/ollama.go` | OllamaProvider — ollama /api/generate 호출 |
| `internal/llm/new_provider.go` | NewProvider — provider 문자열로 Provider 생성 |

## 해시 범위
- func signature + body 해시 (어노테이션/import/package 제외)
- SHA-256의 앞 8자리

## 기존 파일 수정

| 파일 | 변경 |
|---|---|
| `internal/model/annotation.go` | `Checked map[string]string` 필드 추가 |
| `internal/parse/apply_annotation_line.go` | `"checked"` case 추가 |
| `internal/parse/parse_annotation.go` | Checked 맵 초기화 |
| `internal/validate/run_all.go` | A7 룰 호출 추가 |

## 완료 기준
- `filefunc validate`: checked 해시 불일치 시 ERROR, checked 0개면 A7 스킵
- `filefunc llmc`: ollama 호출 후 점수 ≥ 0.8이면 서명 기록
- checked 있고 해시 일치 시 llmc SKIP
- 뮤테이션 테스트 추가
- filefunc 자체 코드가 전체 validate 통과
