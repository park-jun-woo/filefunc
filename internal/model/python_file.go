//ff:type feature=parse type=model
//ff:what Python 소스 파일의 파싱된 구조를 담는 구조체
package model

// PythonFile holds the parsed structure of a single Python source file.
type PythonFile struct {
	Path             string
	Module           string
	Funcs            []string
	Classes          []string
	Methods          []string
	HasInitMethod    bool
	Vars             []string
	Annotation       *Annotation
	Lines            int
	MaxDepth         int
	IsTest           bool
	Control          string
	HasLoopAtDepth1  bool
	HasMatchAtDepth1 bool
	FuncLines        map[string]int
	Q4Violations     []Q4Result
	Calls            []string
	BodyHash         string
}
