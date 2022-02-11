package block

import (
	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
)

// a repeating member of a flow
// could be a stack, a series of inputs, or a series of fields
//
// fix: doesnt handle repeating swaps there are none that i know of currently
//      and not quite sure what they'd look like off hand....
func newRepeat(m *chart.Machine, term string, blk *blockData) chart.State {
	return &chart.StateMix{
		// ex. a series of a specific flow
		OnMap: func(typeName string, flow jsn.FlowBlock) bool {
			next := newSlice(m, term, &blk.inputs)
			m.PushState(next)
			return next.OnMap(typeName, flow)
		},
		// possibly a single stack, or a series of inputs
		OnSlot: func(slotType string, slotBlock jsn.SlotBlock) bool {
			var next *chart.StateMix
			if slot := bconst.FindSlotRule(slotType); slot.Stack {
				next = newStack(m, term, blk)
			} else {
				next = newSeries(m, term, &blk.inputs)
			}
			m.PushState(next)
			return next.OnSlot(slotType, slotBlock)
		},
		// a series of fields
		OnValue: func(n string, pv interface{}) error {
			next := newList(m, term, &blk.fields)
			m.PushState(next)
			return next.OnValue(n, pv)
		},
		OnCommit: func(interface{}) {
			m.FinishState(nil)
		},
	}
}
