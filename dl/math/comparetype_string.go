// Code generated by "stringer -type=CompareType"; DO NOT EDIT.

package math

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Compare_EqualTo-1]
	_ = x[Compare_GreaterThan-2]
	_ = x[Compare_LessThan-4]
}

const (
	_CompareType_name_0 = "Compare_EqualToCompare_GreaterThan"
	_CompareType_name_1 = "Compare_LessThan"
)

var (
	_CompareType_index_0 = [...]uint8{0, 15, 34}
)

func (i CompareType) String() string {
	switch {
	case 1 <= i && i <= 2:
		i -= 1
		return _CompareType_name_0[_CompareType_index_0[i]:_CompareType_index_0[i+1]]
	case i == 4:
		return _CompareType_name_1
	default:
		return "CompareType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
