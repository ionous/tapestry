// Code generated by "stringer -type=Flags"; DO NOT EDIT.

package pattern

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Infix-1]
	_ = x[Prefix-2]
	_ = x[Postfix-4]
}

const (
	_Flags_name_0 = "InfixPrefix"
	_Flags_name_1 = "Postfix"
)

var (
	_Flags_index_0 = [...]uint8{0, 5, 11}
)

func (i Flags) String() string {
	switch {
	case 1 <= i && i <= 2:
		i -= 1
		return _Flags_name_0[_Flags_index_0[i]:_Flags_index_0[i+1]]
	case i == 4:
		return _Flags_name_1
	default:
		return "Flags(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
