// Code generated by "stringer -type=Flow"; DO NOT EDIT.

package event

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Targets-0]
	_ = x[Captures-1]
	_ = x[Bubbles-2]
}

const _Flow_name = "TargetsCapturesBubbles"

var _Flow_index = [...]uint8{0, 7, 15, 22}

func (i Flow) String() string {
	if i < 0 || i >= Flow(len(_Flow_index)-1) {
		return "Flow(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Flow_name[_Flow_index[i]:_Flow_index[i+1]]
}