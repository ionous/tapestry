package block

import (
	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/lang/inspect"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// a repeating member of a flow
// could be a stack, a series of inputs, or a series of fields
func (m *bgen) newRepeat(term string, blk *blockData) inspect.Callbacks {
	return inspect.Callbacks{
		// ex. a series of a specific flow
		OnFlow: func(w inspect.It) error {
			next := m.newSlice(term, &blk.inputs)
			m.events.Replace(next)
			return next.OnFlow(w)
		},
		// possibly a single stack, or a series of inputs
		OnSlot: func(w inspect.It) error {
			var next inspect.Callbacks
			slotType := w.TypeInfo().(*typeinfo.Slot)
			if slot := bconst.MakeSlotRule(slotType); slot.Stack {
				next = m.newStack(term, blk)
			} else {
				next = m.newSeries(term, &blk.inputs)
			}
			m.events.Replace(next)
			return next.OnSlot(w)
		},
		// a series of fields
		OnValue: func(w inspect.It) error {
			next := m.newList(term, &blk.fields)
			m.events.Replace(next)
			return next.OnValue(w)
		},
	}
}
