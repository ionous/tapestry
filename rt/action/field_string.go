// Code generated by "stringer -type=Field -linecomment"; DO NOT EDIT.

package action

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Noun-0]
	_ = x[OtherNoun-1]
	_ = x[Actor-2]
	_ = x[Target-3]
	_ = x[CurrentTarget-4]
	_ = x[Interupt-5]
	_ = x[Cancel-6]
}

const _Field_name = "nounother nounactortargetcurrent targetinterrupt eventcancel event"

var _Field_index = [...]uint8{0, 4, 14, 19, 25, 39, 54, 66}

func (i Field) String() string {
	if i < 0 || i >= Field(len(_Field_index)-1) {
		return "Field(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Field_name[_Field_index[i]:_Field_index[i+1]]
}
