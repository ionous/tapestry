package block

import (
	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/lang/walk"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// writes blockly statement stacks ( for example for: story or execute statements )
// stacks in blockly are.... interesting.
// they are a nested linked list of values.
// this writes the inner halves of the list
func (m *bgen) newStack(term string, blk *blockData) walk.Callbacks {
	// the whole chain is going to be encapsulated by object braces {}
	// we try to keep the same state going for as long as we can...
	var cnt int
	var writingSlot bool
	open, close := js.Obj[0], js.Obj[1]
	return walk.Callbacks{
		// happens before each slat which is received in OnFlow.
		OnSlot: func(w walk.Walker) (_ error) {
			writingSlot = true
			return
		},
		// happens after OnSlot, if and only if the slot is filled.
		OnFlow: func(w walk.Walker) (_ error) {
			typeName := w.TypeName()
			if cnt == 0 {
				_ = blk.startInputWithoutCount(term) // the repeat already wrote the count
			} else {
				blk.inputs.R(js.Comma).Q("next").R(js.Colon).R(open).
					Q("block").R(js.Colon).R(open)
			}
			cnt++ // increment here (rather than OnSlot) to skip any empty slots.
			return m.events.Push(m.newInnerFlow(w, &blk.inputs, bconst.StackedName(typeName)))
		},
		// called after each slot and slot.
		OnEnd: func(w walk.Walker) (err error) {
			// we dont enter a new state for "OnSlot"
			// so we get ends for it and for the end of our own repeat.
			if writingSlot {
				writingSlot = false
			} else {
				for i := 0; i < cnt*2; i++ {
					blk.inputs.R(close)
				}
				err = m.events.Pop()
			}
			return
		},
	}
}
