package jsn

import "github.com/ionous/errutil"

func RepeatBlock(m Marshaler, slice SliceBlock) (err error) {
	if e := m.MarshalBlock(slice); e != nil {
		err = e
	} else {
		// note: marshals the block even if it lacks elements in order to record the empty "[]"
		for i, cnt := 0, slice.GetSize(); i < cnt; i++ {
			if e := slice.MarshalEl(m, i); e != nil {
				m.Error(errutil.New(e, "in slice at", i))
			}
		}
		m.EndBlock()
	}
	return
}
