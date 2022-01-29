package bgen

import (
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/tile/bc"
)

// a repeating member of a flow
// could be a stack, a series of inputs, or a series of fields
//
// fix: doesnt handle repeating swaps there are none that i know of currently
//      and not quite sure what they'd look like off hand....
func newRepeat(m *chart.Machine, term string, data *blockData) chart.State {
	return &chart.StateMix{
		// a series of inputs
		// OnMap: func(typeName string, _ jsn.FlowBlock) bool {
		// 	return false
		// },
		// possibly a single stack, or a series of inputs
		OnSlot: func(slotType string, slotBlock jsn.SlotBlock) bool {
			var next *chart.StateMix
			if slot := bc.FindSlotRule(slotType); slot.Stack {
				next = newStack(m, &data.next)
			} else {
				next = newSeries(m, term, &data.inputs)
			}
			m.PushState(next)
			return next.OnSlot(slotType, slotBlock)
		},
		// a series of fields
		OnValue: func(n string, pv interface{}) error {
			next := newList(m, term, &data.fields)
			m.PushState(next)
			return next.OnValue(n, pv)
		},
		OnCommit: func(interface{}) {
			m.FinishState(nil)
		},
	}
}
