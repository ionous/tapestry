package prim

import "strings"

func (op *Text_Unboxed_Slice) GetCompactValue() (ret interface{}) {
	strs := []string(*op)
	if len(strs) == 1 {
		ret = strs[0]
	} else {
		ret = strs
	}
	return
}

func (op *Lines) GetLines() []string {
	return strings.FieldsFunc(op.Str, func(r rune) bool { return r == newline })
}

const newline = '\n'
