//ff:type feature=validate type=model
//ff:what CountMax 룰의 backing — 수량 상한 판정 기준
package validate

// CountMaxBacking defines the judgment criteria for count-based rules.
type CountMaxBacking struct {
	Field   string // GoFile field: "Funcs", "Types", "Methods"
	Max     int
	Rule    string
	Message string
}
