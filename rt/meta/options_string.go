// Code generated by "stringer -type=Options"; DO NOT EDIT.

package meta

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[PrintResponseNames-0]
	_ = x[PrintPatternNames-1]
	_ = x[CacheErrors-2]
	_ = x[SaveDir-3]
	_ = x[NumOptions-4]
}

const _Options_name = "PrintResponseNamesPrintPatternNamesCacheErrorsSaveDirNumOptions"

var _Options_index = [...]uint8{0, 18, 35, 46, 53, 63}

func (i Options) String() string {
	if i < 0 || i >= Options(len(_Options_index)-1) {
		return "Options(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Options_name[_Options_index[i]:_Options_index[i+1]]
}
