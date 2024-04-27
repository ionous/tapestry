// Code generated by "stringer -type=Signal -trimprefix=Signal"; DO NOT EDIT.

package game

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SignalUnknown-0]
	_ = x[SignalQuit-1]
	_ = x[SignalSave-2]
	_ = x[SignalLoad-3]
}

const _Signal_name = "UnknownQuitSaveLoad"

var _Signal_index = [...]uint8{0, 7, 11, 15, 19}

func (i Signal) String() string {
	if i < 0 || i >= Signal(len(_Signal_index)-1) {
		return "Signal(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Signal_name[_Signal_index[i]:_Signal_index[i+1]]
}
