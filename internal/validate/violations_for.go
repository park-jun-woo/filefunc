//ff:func feature=validate type=util control=selection
//ff:what toulmin verdict가 양수인 룰에 대해 원본 Check 함수를 호출하여 Violation 반환
package validate

import "github.com/park-jun-woo/filefunc/internal/model"

// violationsFor dispatches to the original Check function by rule name
// to get the detailed violation with dynamic messages.
func violationsFor(ruleName string, gf *model.GoFile, cb *model.Codebook) []model.Violation {
	switch ruleName {
	case "RuleF1":
		return CheckOneFileOneFunc(gf)
	case "RuleF2":
		return CheckOneFileOneType(gf)
	case "RuleF3":
		return CheckOneFileOneMethod(gf)
	case "RuleF4":
		return CheckInitStandalone(gf)
	case "RuleQ1":
		return CheckNestingDepth(gf)
	case "RuleQ2Q3":
		return CheckFuncLines(gf)
	case "RuleA1":
		return CheckAnnotationRequired(gf)
	case "RuleA2":
		return CheckCodebookValues(gf, cb)
	case "RuleA3":
		return CheckWhatRequired(gf)
	case "RuleA6":
		return CheckAnnotationPosition(gf)
	case "RuleA7":
		return CheckCheckedHash(gf)
	case "RuleA8":
		return CheckRequiredKeysInAnnotation(gf, cb)
	case "RuleA9":
		return CheckControlRequired(gf)
	case "RuleA10":
		return CheckControlSelection(gf)
	case "RuleA11":
		return CheckControlIteration(gf)
	case "RuleA12":
		return CheckControlSequence(gf)
	case "RuleA13":
		return CheckControlSelectionNoLoop(gf)
	case "RuleA14":
		return CheckControlIterationNoSwitch(gf)
	case "RuleA15":
		return CheckDimensionRequired(gf)
	case "RuleA16":
		return CheckDimensionValue(gf)
	default:
		return nil
	}
}
