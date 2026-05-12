//ff:type feature=parse type=model
//ff:what ts_ast.js JSON 출력을 디시리얼라이즈하는 구조체
package parse

// TsAstResult is the JSON structure output by scripts/ts_ast.js.
type TsAstResult struct {
	Path            string         `json:"path"`
	Functions       []string       `json:"functions"`
	Classes         []string       `json:"classes"`
	Interfaces      []string       `json:"interfaces"`
	TypeAliases     []string       `json:"type_aliases"`
	Methods         []string       `json:"methods"`
	HasConstructor  bool           `json:"has_constructor"`
	Vars            []string       `json:"vars"`
	Lines           int            `json:"lines"`
	MaxDepth        int            `json:"max_depth"`
	Control         string         `json:"control"`
	HasLoopAtDepth1 bool           `json:"has_loop_at_depth1"`
	HasSwitchAtD1   bool           `json:"has_switch_at_depth1"`
	FuncLines       map[string]int `json:"func_lines"`
	Q4Results       []TsQ4JSON     `json:"q4_results"`
	Calls           []string       `json:"calls"`
	Imports         []TsImportJSON `json:"imports"`
	BodyHash        string         `json:"body_hash"`
	Error           string         `json:"error,omitempty"`
}
