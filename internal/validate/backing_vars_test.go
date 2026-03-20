package validate

var (
	backingF1 = &CountMaxBacking{Field: "Funcs", Max: 1, Rule: "F1",
		Message: "file contains multiple funcs; expected 1 file 1 func"}
	backingF2 = &CountMaxBacking{Field: "Types", Max: 1, Rule: "F2",
		Message: "file contains multiple types; expected 1 file 1 type"}
	backingF3 = &CountMaxBacking{Field: "Methods", Max: 1, Rule: "F3",
		Message: "file contains multiple methods; expected 1 file 1 method"}
	backingF4 = &ExistsWhenBacking{When: "HasInit", Need: "companion", Rule: "F4",
		Level: "ERROR", Message: "init() must not exist alone; requires accompanying var or func"}
	backingA1f = &ExistsWhenBacking{When: "HasFuncs", Need: "ff:func", Rule: "A1",
		Level: "ERROR", Message: "file with func must have //ff:func annotation"}
	backingA1t = &ExistsWhenBacking{When: "HasTypes", Need: "ff:type", Rule: "A1",
		Level: "ERROR", Message: "file with type must have //ff:type annotation"}
	backingA3 = &ExistsWhenBacking{When: "HasFuncOrType", Need: "ff:what", Rule: "A3",
		Level: "ERROR", Message: "file with func or type must have //ff:what annotation"}
	backingA9 = &ExistsWhenBacking{When: "HasFuncs", Need: "control", Rule: "A9",
		Level: "ERROR", Message: "func file must have control= annotation (sequence, selection, or iteration)"}
	backingA15 = &ExistsWhenBacking{When: "ControlIteration", Need: "dimension", Rule: "A15",
		Level: "ERROR", Message: "control=iteration requires dimension= annotation"}
	backingA2  = &InCodebookBacking{Direction: "value→codebook", Rule: "A2"}
	backingA8  = &InCodebookBacking{Direction: "codebook→annotation", Rule: "A8"}
	backingA10 = &ControlMatchBacking{Control: "selection", MustHave: "switch", Rule: "A10",
		Message: "control=selection but no switch found at depth 1"}
	backingA11 = &ControlMatchBacking{Control: "iteration", MustHave: "loop", Rule: "A11",
		Message: "control=iteration but no loop found at depth 1"}
	backingA12 = &ControlMatchBacking{Control: "sequence", MustNotHave: "switch|loop", Rule: "A12",
		Message: "control=sequence but %s found at depth 1; add control=%s or extract to separate func"}
	backingA13 = &ControlMatchBacking{Control: "selection", MustNotHave: "loop", Rule: "A13",
		Message: "control=selection but loop found at depth 1; extract loop to separate func"}
	backingA14 = &ControlMatchBacking{Control: "iteration", MustNotHave: "switch", Rule: "A14",
		Message: "control=iteration but switch found at depth 1; extract switch to separate func"}
)
