//ff:type feature=validate type=model
//ff:what InCodebook 룰의 backing — 코드북 적합 판정 기준
package validate

// InCodebookBacking defines the judgment criteria for codebook conformance rules.
type InCodebookBacking struct {
	Direction string // "value→codebook" (A2) or "codebook→annotation" (A8)
	Rule      string
}
