# filefunc chain — 호출 관계 맵

## chain func 예시

### ParseGoFile (chon=2, --meta what)
```
ParseGoFile (what="Go 소스 파일을 파싱하여 GoFile 구조체로 변환")
  1촌 calls: CalcMaxDepth (what="파일 내 모든 함수의 최대 nesting depth 계산")
  1촌 calls: CollectFuncDecl (what="FuncDecl에서 func/method/init 정보를 GoFile에 수집")
  1촌 calls: CollectGenDecl (what="GenDecl에서 type/var 정보를 GoFile에 수집")
  1촌 calls: ParseAnnotation (what="Go 소스 파일 상단의 //ff: 어노테이션을 파싱")
  1촌 called-by: BuildGraph (what="프로젝트 루트에서 호출 그래프를 구성")
  2촌 calls-2depth: StmtDepth (what="문장 목록의 최대 nesting depth 계산")
  2촌 calls-2depth: ReceiverTypeName (what="메서드 리시버의 타입명을 추출")
  2촌 calls-2depth: CollectVarNames (what="ValueSpec에서 변수명 목록을 추출")
  2촌 calls-2depth: ApplyAnnotationLine (what="//ff: key-value 쌍을 Annotation 구조체에 적용")
  2촌 calls-2depth: IsSkippableLine (what="어노테이션 파싱 시 건너뛸 수 있는 라인인지 판별")
  2촌 co-called: BuildCallGraph (what="프로젝트 전체 호출 그래프를 양방향으로 구성")
  2촌 co-called: CollectProjectSymbols (what="프로젝트 전체 GoFile에서 func/type 이름을 수집하여 심볼 맵 반환")
  2촌 co-called: ParseFFIgnore (what=".ffignore 파일을 읽어 패턴 목록을 반환 (없으면 빈 목록)")
  2촌 co-called: ReadModulePath (what="go.mod에서 모듈 경로를 추출")
  2촌 co-called: WalkGoFiles (what="디렉토리를 재귀 순회하며 .go 파일 경로 목록 반환 (.ffignore 적용)")
```

### RunAll (chon=2, --meta what)
```
RunAll (what="모든 검증 룰을 실행하고 위반 목록을 반환")
  1촌 calls: CheckAnnotationPosition (what="A6: //ff: 어노테이션이 파일 최상단에 위치하는지 검증")
  1촌 calls: CheckAnnotationRequired (what="A1: func 파일은 //ff:func, type 파일은 //ff:type 필수 검증")
  1촌 calls: CheckCheckedHash (what="A7: //ff:checked 해시와 현재 body 해시를 대조하여 불일치 시 ERROR")
  1촌 calls: CheckCodebookValues (what="A2: 어노테이션 값이 코드북에 존재하는지 검증")
  1촌 calls: CheckControlIteration (what="A11: control=iteration인데 depth 1에 loop 없으면 ERROR")
  1촌 calls: CheckControlIterationNoSwitch (what="A14: control=iteration인데 depth 1에 switch 존재하면 ERROR")
  1촌 calls: CheckControlRequired (what="A9: func 파일은 control= 어노테이션 필수 (sequence/selection/iteration)")
  1촌 calls: CheckControlSelection (what="A10: control=selection인데 depth 1에 switch 없으면 ERROR")
  1촌 calls: CheckControlSelectionNoLoop (what="A13: control=selection인데 depth 1에 loop 존재하면 ERROR")
  1촌 calls: CheckControlSequence (what="A12: control=sequence인데 depth 1에 switch/loop 존재하면 ERROR")
  1촌 calls: CheckDimensionRequired (what="A15: control=iteration이면 dimension= 필수")
  1촌 calls: CheckDimensionValue (what="A16: dimension= 값은 양의 정수여야 함")
  1촌 calls: CheckFuncLines (what="Q2/Q3: func 라인 수 검증. Q3은 control별 기준 (sequence/iteration 100, selection 300)")
  1촌 calls: CheckInitStandalone (what="F4: init()만 단독으로 존재하는 파일 검증")
  1촌 calls: CheckNestingDepth (what="Q1: control과 dimension 기반으로 nesting depth 상한 검증")
  1촌 calls: CheckOneFileOneFunc (what="F1: 파일당 func 1개 검증")
  1촌 calls: CheckOneFileOneMethod (what="F3: 파일당 method 1개 검증")
  1촌 calls: CheckOneFileOneType (what="F2: 파일당 type 1개 검증")
  1촌 calls: CheckRequiredKeysInAnnotation (what="A8: 어노테이션에 codebook required 키가 모두 존재하는지 검증")
  1촌 calls: CheckWhatRequired (what="A3: func 또는 type 파일에 //ff:what 필수 검증")
  1촌 calls: HasAnyChecked (what="프로젝트 파일 목록에서 //ff:checked가 하나라도 있는지 판별")
  2촌 calls-2depth: CalcBodyHash (what="func signature+body의 SHA-256 해시 앞 8자리 계산")
  2촌 calls-2depth: AllowedValues (what="codebook의 required+optional에서 주어진 키의 허용 값 목록을 통합 반환")
  2촌 calls-2depth: Contains (what="문자열 슬라이스에 특정 항목이 포함되어 있는지 확인")
  2촌 calls-2depth: DetectControl (what="func body의 직계 자식(depth 1)을 순회하여 제어구조를 판별")
  2촌 calls-2depth: HasSwitchAtDepth1 (what="depth 1에 SwitchStmt/TypeSwitchStmt가 존재하는지 판별")
  2촌 calls-2depth: HasLoopAtDepth1 (what="depth 1에 ForStmt/RangeStmt가 존재하는지 판별")
  2촌 calls-2depth: HasBacktick (what="FuncDecl 내에 백틱 문자열 리터럴이 존재하는지 판별")
  2촌 calls-2depth: Q3Limit (what="control 값에 따라 Q3 줄 수 제한을 반환")
  2촌 calls-2depth: depthLimit (what="control과 dimension으로 Q1 depth 상한을 계산")
  2촌 calls-2depth: IsConstOnly (what="파일이 const/var만 포함하는지 판별")
```

### CheckOneFileOneFunc (chon=3, --meta what)
```
CheckOneFileOneFunc (what="F1: 파일당 func 1개 검증")
  1촌 calls: IsConstOnly (what="파일이 const/var만 포함하는지 판별")
  1촌 called-by: RunAll (what="모든 검증 룰을 실행하고 위반 목록을 반환")
  2촌 co-called: CheckAnnotationPosition (what="A6: //ff: 어노테이션이 파일 최상단에 위치하는지 검증")
  2촌 co-called: CheckAnnotationRequired (what="A1: func 파일은 //ff:func, type 파일은 //ff:type 필수 검증")
  2촌 co-called: CheckCheckedHash (what="A7: //ff:checked 해시와 현재 body 해시를 대조하여 불일치 시 ERROR")
  2촌 co-called: CheckCodebookValues (what="A2: 어노테이션 값이 코드북에 존재하는지 검증")
  2촌 co-called: CheckControlIteration (what="A11: control=iteration인데 depth 1에 loop 없으면 ERROR")
  2촌 co-called: CheckControlIterationNoSwitch (what="A14: control=iteration인데 depth 1에 switch 존재하면 ERROR")
  2촌 co-called: CheckControlRequired (what="A9: func 파일은 control= 어노테이션 필수 (sequence/selection/iteration)")
  2촌 co-called: CheckControlSelection (what="A10: control=selection인데 depth 1에 switch 없으면 ERROR")
  2촌 co-called: CheckControlSelectionNoLoop (what="A13: control=selection인데 depth 1에 loop 존재하면 ERROR")
  2촌 co-called: CheckControlSequence (what="A12: control=sequence인데 depth 1에 switch/loop 존재하면 ERROR")
  2촌 co-called: CheckDimensionRequired (what="A15: control=iteration이면 dimension= 필수")
  2촌 co-called: CheckDimensionValue (what="A16: dimension= 값은 양의 정수여야 함")
  2촌 co-called: CheckFuncLines (what="Q2/Q3: func 라인 수 검증. Q3은 control별 기준 (sequence/iteration 100, selection 300)")
  2촌 co-called: CheckInitStandalone (what="F4: init()만 단독으로 존재하는 파일 검증")
  2촌 co-called: CheckNestingDepth (what="Q1: control과 dimension 기반으로 nesting depth 상한 검증")
  2촌 co-called: CheckOneFileOneMethod (what="F3: 파일당 method 1개 검증")
  2촌 co-called: CheckOneFileOneType (what="F2: 파일당 type 1개 검증")
  2촌 co-called: CheckRequiredKeysInAnnotation (what="A8: 어노테이션에 codebook required 키가 모두 존재하는지 검증")
  2촌 co-called: CheckWhatRequired (what="A3: func 또는 type 파일에 //ff:what 필수 검증")
  2촌 co-called: HasAnyChecked (what="프로젝트 파일 목록에서 //ff:checked가 하나라도 있는지 판별")
  3촌 peer-calls: CalcBodyHash (what="func signature+body의 SHA-256 해시 앞 8자리 계산")
  3촌 peer-calls: AllowedValues (what="codebook의 required+optional에서 주어진 키의 허용 값 목록을 통합 반환")
  3촌 peer-calls: Contains (what="문자열 슬라이스에 특정 항목이 포함되어 있는지 확인")
  3촌 peer-calls: DetectControl (what="func body의 직계 자식(depth 1)을 순회하여 제어구조를 판별")
  3촌 peer-calls: HasSwitchAtDepth1 (what="depth 1에 SwitchStmt/TypeSwitchStmt가 존재하는지 판별")
  3촌 peer-calls: HasLoopAtDepth1 (what="depth 1에 ForStmt/RangeStmt가 존재하는지 판별")
  3촌 peer-calls: HasBacktick (what="FuncDecl 내에 백틱 문자열 리터럴이 존재하는지 판별")
  3촌 peer-calls: Q3Limit (what="control 값에 따라 Q3 줄 수 제한을 반환")
  3촌 peer-calls: depthLimit (what="control과 dimension으로 Q1 depth 상한을 계산")
```

### VerifyWhat (chon=2, --meta what)
```
VerifyWhat (what="Provider로 what-body 일치 점수를 판정")
  1촌 calls: BuildPrompt (what="what과 func body로 LLM 검증 프롬프트를 생성")
  1촌 calls: ParseScore (what="LLM 응답 문자열에서 0.0~1.0 점수를 추출")
  1촌 called-by: ProcessLlmcFile (what="단일 파일에 대해 LLM what-body 검증을 수행하고 결과를 반환")
  2촌 co-called: CalcBodyHash (what="func signature+body의 SHA-256 해시 앞 8자리 계산")
  2촌 co-called: ExtractFuncSource (what="Go AST로 파일의 첫 번째 func(init 제외)의 signature+body를 소스 텍스트로 추출하여 반환")
  2촌 co-called: WriteAnnotationLine (what="파일의 //ff: 라인을 추가/갱신/제거")
```

## chain feature 예시

### feature=validate (34 funcs, --meta what)

핵심 구조:
- `ValidateCodebook` → C1~C3 (codebook 형식 검증)
- `RunAll` → F1~F4, Q1~Q3, A1~A16 (코드 검증 전체)
- 각 Check 함수는 독립적, RunAll에서만 호출됨

### feature=parse (35 funcs, --meta what)

핵심 구조:
- `ParseGoFile` → CalcMaxDepth, CollectFuncDecl, CollectGenDecl, ParseAnnotation
- `BuildCallGraph` → ExtractCalls → BuildImportMap, CallName, SortedKeys
- `DetectControl` → detectFromBody (control 구조 판별)
- `HasLoopAtDepth1`/`HasSwitchAtDepth1` → firstFuncBodyStmts + containsLoop/containsSwitch

### feature=chain (18 funcs, --meta what)

핵심 구조:
- `TraverseChon` → CollectChon, ExpandThrough, FindSiblings (촌수 탐색)
- `TraverseDepth` → traverseDepthRecur (단방향 깊이 탐색)
- `FormatChain` → formatMeta → metaPairs, checkedString (출력)
- `BuildFuncFileMap` (함수명 → GoFile 매핑)
- `ParseMetaFlags` (--meta 플래그 파싱)

### feature=cli (15 funcs, --meta what)

핵심 구조:
- `main` → `Execute` (엔트리포인트)
- `BuildGraph` → WalkGoFiles, ParseGoFile, BuildCallGraph, CollectProjectSymbols
- `ProcessLlmcFile` → CalcBodyHash, ExtractFuncSource, VerifyWhat, WriteAnnotationLine
- `CheckModel` → ModelExists, PullModel (ollama 모델 관리)
