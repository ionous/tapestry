// Code generated by "stringer -type=Type"; DO NOT EDIT.

package walk

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Flow-0]
	_ = x[Slot-1]
	_ = x[Swap-2]
	_ = x[Str-3]
	_ = x[Num-4]
	_ = x[Value-5]
}

const _Type_name = "FlowSlotSwapStrNumValue"

var _Type_index = [...]uint8{0, 4, 8, 12, 15, 18, 23}

func (i Type) String() string {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
