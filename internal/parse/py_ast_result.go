//ff:type feature=parse type=model
//ff:what py_ast.py JSON 출력을 디시리얼라이즈하는 구조체
package parse

// PyAstResult is the JSON structure output by scripts/py_ast.py.
type PyAstResult struct {
	Path            string         `json:"path"`
	Functions       []string       `json:"functions"`
	Classes         []string       `json:"classes"`
	Methods         []string       `json:"methods"`
	HasInitMethod   bool           `json:"has_init_method"`
	Vars            []string       `json:"vars"`
	Lines           int            `json:"lines"`
	MaxDepth        int            `json:"max_depth"`
	Control         string         `json:"control"`
	HasLoopAtDepth1 bool           `json:"has_loop_at_depth1"`
	HasMatchAtD1    bool           `json:"has_match_at_depth1"`
	FuncLines       map[string]int `json:"func_lines"`
	Q4Results       []PyQ4JSON     `json:"q4_results"`
	Calls           []string       `json:"calls"`
	Imports         []PyImportJSON `json:"imports"`
	ModuleImports   []PyImportJSON `json:"module_imports"`
	BodyHash        string         `json:"body_hash"`
	Error           string         `json:"error,omitempty"`
}
