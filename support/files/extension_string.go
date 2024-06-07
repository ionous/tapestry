// Code generated by "stringer -type=Extension -linecomment"; DO NOT EDIT.

package files

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[BlockExt-1]
	_ = x[CompactExt-2]
	_ = x[SpecExt-3]
	_ = x[TellSpec-4]
	_ = x[TellStory-5]
}

const _Extension_name = ".block.if.ifspecs.idl.tell"

var _Extension_index = [...]uint8{0, 6, 9, 17, 21, 26}

func (i Extension) String() string {
	i -= 1
	if i < 0 || i >= Extension(len(_Extension_index)-1) {
		return "Extension(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Extension_name[_Extension_index[i]:_Extension_index[i+1]]
}
