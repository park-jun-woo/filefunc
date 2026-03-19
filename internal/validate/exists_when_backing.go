//ff:type feature=validate type=model
//ff:what ExistsWhen 룰의 backing — 조건부 존재 검사 판정 기준
package validate

// ExistsWhenBacking defines the judgment criteria for conditional existence rules.
type ExistsWhenBacking struct {
	When    string // precondition: "HasFuncs", "HasTypes", "HasFuncOrType", "HasInit", "ControlIteration"
	Need    string // must exist: "ff:func", "ff:type", "ff:what", "control", "dimension", "companion"
	Rule    string
	Level   string
	Message string
}
