// Code generated by "stringer -type=eventSuffix -linecomment"; DO NOT EDIT.

package rules

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[continues-0]
	_ = x[stops-1]
	_ = x[jumps-2]
}

const _eventSuffix_name = "continuestopjump"

var _eventSuffix_index = [...]uint8{0, 8, 12, 16}

func (i eventSuffix) String() string {
	if i < 0 || i >= eventSuffix(len(_eventSuffix_index)-1) {
		return "eventSuffix(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _eventSuffix_name[_eventSuffix_index[i]:_eventSuffix_index[i+1]]
}