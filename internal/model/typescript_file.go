//ff:type feature=parse type=model
//ff:what TypeScript 소스 파일의 파싱된 구조를 담는 구조체
package model

// TypeScriptFile holds the parsed structure of a single TypeScript source file.
type TypeScriptFile struct {
	Path              string
	Module            string
	Funcs             []string
	Classes           []string
	Interfaces        []string
	TypeAliases       []string
	Methods           []string
	HasConstructor    bool
	Vars              []string
	Annotation        *Annotation
	Lines             int
	MaxDepth          int
	IsTest            bool
	Control           string
	HasLoopAtDepth1   bool
	HasSwitchAtDepth1 bool
	FuncLines         map[string]int
	Q4Violations      []Q4Result
	Calls             []string
	BodyHash          string
}
