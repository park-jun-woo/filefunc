# Phase 012: .ffignore 지원 ✅ 완료

## 목표
`.ffignore` 파일로 filefunc의 모든 명령(validate, chain, llmc)에서 제외할 경로를 지정한다. `.gitignore`와 동일한 패턴 문법.

## .ffignore 예시

```
# 서드파티
vendor/

# 생성 코드
*.pb.go
*_gen.go

# 특정 디렉토리
internal/legacy/
```

## 적용 범위
- `WalkGoFiles`에서 .ffignore 패턴에 매칭되는 파일/디렉토리를 건너뜀
- validate, chain, llmc 모두 WalkGoFiles를 사용하므로 한 곳만 수정하면 전체 적용

## chain 탐색 결과

```
WalkGoFiles
  1촌 called-by: BuildGraph
  2촌 co-called: BuildCallGraph, CollectProjectSymbols, FindGoMod, ParseGoFile, ReadModulePath, WalkGoFiles
```

## 산출물

| 파일 | 내용 | 구분 |
|---|---|---|
| `internal/walk/parse_ffignore.go` | ParseFFIgnore — .ffignore 파일을 읽어 패턴 목록 반환 | 신규 |
| `internal/walk/match_ffignore.go` | MatchFFIgnore — 경로가 패턴에 매칭되는지 판별 | 신규 |
| `internal/walk/walk_go_files.go` | 시그니처 변경: `WalkGoFiles(root string, ignorePatterns []string)` | 수정 |
| `internal/cli/build_graph.go` | WalkGoFiles 호출부 수정 (ignorePatterns 전달) | 수정 |
| `internal/cli/validate.go` | WalkGoFiles 호출부 수정 (ignorePatterns 전달) | 수정 |

## 시그니처 변경 영향

`WalkGoFiles(root string)` → `WalkGoFiles(root string, ignorePatterns []string)`

호출하는 곳:
- `cli/build_graph.go` — chain, llmc가 사용
- `cli/validate.go` — validate가 사용

두 곳 모두 .ffignore를 로드하여 전달.

## 기존 코드 재사용

| 파일 | 용도 |
|---|---|
| `cli/find_go_mod.go` | go.mod 위치로 .ffignore 위치 추정 |
| `cli/find_go_mod_dir.go` | 디렉토리 경로 반환 |

## 패턴 문법
- `.gitignore`와 동일: glob 패턴, `/`로 디렉토리 지정, `#`으로 주석, 빈 줄 무시
- `filepath.Match` 기반 + 디렉토리 매칭은 `strings.HasPrefix`

## 실행 흐름

```
validate/chain/llmc
  │
  ├─ FindGoModDir → 프로젝트 루트
  ├─ ParseFFIgnore(루트/.ffignore) → []string (없으면 빈 목록)
  │
  └─ WalkGoFiles(target, ignorePatterns)
       └─ filepath.Walk:
            ├─ 디렉토리: testdata 또는 MatchFFIgnore 매칭 → SkipDir
            └─ 파일: MatchFFIgnore 매칭 → skip
```

## .ffignore 위치
- go.mod와 같은 디렉토리에서 탐색
- 없으면 무시 (에러 아님, 빈 패턴 목록)

## 완료 기준
- .ffignore에 지정된 경로가 validate/chain/llmc에서 제외
- .ffignore 없으면 기존과 동일하게 동작
- 디렉토리 패턴 (`vendor/`) 동작
- 파일 패턴 (`*.pb.go`) 동작
- 테스트 추가
- filefunc 자체 코드가 전체 validate 통과
- README.md, manual-for-ai.md 업데이트
