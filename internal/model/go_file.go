//ff:type feature=parse type=model
//ff:what Go 소스 파일의 파싱된 구조를 담는 구조체
package model

// GoFile holds the parsed structure of a single Go source file.
type GoFile struct {
	Path       string      // file path
	Package    string      // package name
	Funcs      []string    // top-level func names (excluding methods and init)
	Types      []string    // type names
	Methods    []string    // method names (receiver.Method)
	HasInit    bool        // whether the file contains func init()
	Vars       []string    // top-level var names
	Annotation *Annotation // parsed //ff: annotation (nil if absent)
	Lines      int         // total line count
	MaxDepth   int         // maximum nesting depth
	IsTest     bool        // whether the file is a _test.go file
}
