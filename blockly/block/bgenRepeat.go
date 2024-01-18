package block

import (
	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/lang/walk"
)

// a repeating member of a flow
// could be a stack, a series of inputs, or a series of fields
//
// fix: doesnt handle repeating swaps there are none that i know of currently
//
//	and not quite sure what they'd look like off hand....
func (m *bgen) newRepeat(term string, blk *blockData) walk.Callbacks {
	return walk.Callbacks{
		// ex. a series of a specific flow
		OnFlow: func(w walk.Walker) error {
			next := m.newSlice(term, &blk.inputs)
			m.events.Replace(next)
			return next.OnFlow(w)
		},
		// possibly a single stack, or a series of inputs
		OnSlot: func(w walk.Walker) error {
			var next walk.Callbacks
			slotType := w.TypeName()
			if slot := bconst.FindSlotRule(m.types, slotType); slot.Stack {
				next = m.newStack(term, blk)
			} else {
				next = m.newSeries(term, &blk.inputs)
			}
			m.events.Replace(next)
			return next.OnSlot(w)
		},
		// a series of fields
		OnValue: func(w walk.Walker) error {
			next := m.newList(term, &blk.fields)
			m.events.Replace(next)
			return next.OnValue(w)
		},
	}
}
