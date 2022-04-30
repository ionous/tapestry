package prim

func (op *Text_Unboxed_Slice) GetCompactValue() (ret interface{}) {
	strs := []string(*op)
	if len(strs) == 1 {
		ret = strs[0]
	} else {
		ret = strs
	}
	return
}
