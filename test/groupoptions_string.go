// Code generated by "stringer -type=GroupOptions"; DO NOT EDIT.

package test

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[WithoutObjects-0]
	_ = x[WithoutArticles-1]
	_ = x[WithArticles-2]
}

const _GroupOptions_name = "WithoutObjectsWithoutArticlesWithArticles"

var _GroupOptions_index = [...]uint8{0, 14, 29, 41}

func (i GroupOptions) String() string {
	if i < 0 || i >= GroupOptions(len(_GroupOptions_index)-1) {
		return "GroupOptions(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _GroupOptions_name[_GroupOptions_index[i]:_GroupOptions_index[i+1]]
}