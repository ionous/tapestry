// Code generated by "stringer -type=Suffix -linecomment"; DO NOT EDIT.

package rules

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UnspecfiedSuffix-0]
	_ = x[Continues-1]
	_ = x[Stops-2]
	_ = x[Skips-3]
}

const _Suffix_name = "UnspecfiedSuffixthen continuethen stopthen skip phase"

var _Suffix_index = [...]uint8{0, 16, 29, 38, 53}

func (i Suffix) String() string {
	if i < 0 || i >= Suffix(len(_Suffix_index)-1) {
		return "Suffix(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Suffix_name[_Suffix_index[i]:_Suffix_index[i+1]]
}
