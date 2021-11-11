package jsn

import "github.com/ionous/errutil"

func RepeatBlock(m Marshaler, slice SliceBlock) (err error) {
	if e := m.MarshalBlock(slice); e != nil {
		err = e
	} else {
		// note: marshals the block even if it lacks elements in order to record the empty "[]"
		for i, cnt := 0, slice.GetSize(); i < cnt; i++ {
			// fix: we exclude missing right now b/c of the way slot cin reading works :/
			if e := slice.MarshalEl(m, i); e != nil && e != Missing {
				m.Error(errutil.New(e, "in slice at", i))
			}
		}
		m.EndBlock()
	}
	return
}
