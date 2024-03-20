// Code generated by "stringer -type=Phase"; DO NOT EDIT.

package mdl

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DependencyPhase-1]
	_ = x[LanguagePhase-2]
	_ = x[AncestryPhase-3]
	_ = x[PropertyPhase-4]
	_ = x[MappingPhase-5]
	_ = x[NounPhase-6]
	_ = x[MacroPhase-7]
	_ = x[ConnectionPhase-8]
	_ = x[FallbackPhase-9]
	_ = x[ValuePhase-10]
	_ = x[RulePhase-11]
	_ = x[FinalPhase-12]
	_ = x[NumPhases-13]
}

const _Phase_name = "DependencyPhaseLanguagePhaseAncestryPhasePropertyPhaseMappingPhaseNounPhaseMacroPhaseConnectionPhaseFallbackPhaseValuePhaseRulePhaseFinalPhaseNumPhases"

var _Phase_index = [...]uint8{0, 15, 28, 41, 54, 66, 75, 85, 100, 113, 123, 132, 142, 151}

func (i Phase) String() string {
	i -= 1
	if i < 0 || i >= Phase(len(_Phase_index)-1) {
		return "Phase(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Phase_name[_Phase_index[i]:_Phase_index[i+1]]
}
