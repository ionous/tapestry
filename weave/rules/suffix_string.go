// Code generated by "stringer -type=Suffix -linecomment"; DO NOT EDIT.

package rules

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Continues-0]
	_ = x[Stops-1]
	_ = x[Jumps-2]
	_ = x[Begins-3]
	_ = x[Ends-4]
}

const _Suffix_name = "then continuethen stopthen jumpbeginsend"

var _Suffix_index = [...]uint8{0, 13, 22, 31, 37, 40}

func (i Suffix) String() string {
	if i < 0 || i >= Suffix(len(_Suffix_index)-1) {
		return "Suffix(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Suffix_name[_Suffix_index[i]:_Suffix_index[i+1]]
}