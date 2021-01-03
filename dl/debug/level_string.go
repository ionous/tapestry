// Code generated by "stringer -type=Level"; DO NOT EDIT.

package debug

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Note-0]
	_ = x[ToDo-1]
	_ = x[Warning-2]
	_ = x[Fix-3]
}

const _Level_name = "NoteToDoWarningFix"

var _Level_index = [...]uint8{0, 4, 8, 15, 18}

func (i Level) String() string {
	if i < 0 || i >= Level(len(_Level_index)-1) {
		return "Level(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Level_name[_Level_index[i]:_Level_index[i+1]]
}
