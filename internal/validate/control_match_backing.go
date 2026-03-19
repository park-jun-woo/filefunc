//ff:type feature=validate type=model
//ff:what ControlMatch 룰의 backing — 제어 구조 일치 판정 기준
package validate

// ControlMatchBacking defines the judgment criteria for control structure consistency rules.
type ControlMatchBacking struct {
	Control     string // declared control: "selection", "iteration", "sequence"
	MustHave    string // required construct: "switch", "loop", ""
	MustNotHave string // forbidden construct: "loop", "switch", "switch|loop", ""
	Rule        string
	Message     string
}
