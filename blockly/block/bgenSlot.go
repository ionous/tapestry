package block

import (
	"git.sr.ht/~ionous/tapestry/lang/inspect"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

func (m *bgen) newSlot(term string, blk *blockData) inspect.Callbacks {
	was := -1
	return inspect.Callbacks{
		OnFlow: func(w inspect.It) (_ error) {
			isFlow := true // previously, if there were strs implementing slots then
			was = blk.startInput(term)
			typeName := w.TypeInfo().(*typeinfo.Flow).Name
			return m.events.Push(m.newInnerBlock(w, &blk.inputs, typeName, isFlow))
		},
		OnEnd: func(w inspect.It) error {
			if was >= 0 {
				blk.endInput(was)
				was = -1
			}
			return m.events.Pop()
		},
	}
}
