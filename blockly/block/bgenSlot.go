package block

import (
	"git.sr.ht/~ionous/tapestry/lang/walk"
)

// see also: newSwap
func (m *bgen) newSlot(term string, blk *blockData) walk.Callbacks {
	was := -1
	return walk.Callbacks{
		OnFlow: func(w walk.Walker) (_ error) {
			isFlow := true // previously, if there were strs implementing slots then
			was = blk.startInput(term)
			return m.events.Push(m.newInnerBlock(w, &blk.inputs, w.TypeName(), isFlow))
		},
		OnEnd: func(w walk.Walker) error {
			if was >= 0 {
				blk.endInput(was)
				was = -1
			}
			return m.events.Pop()
		},
	}
}
