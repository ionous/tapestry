package bgen

import (
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// writes blockly statement stacks ( for example for: story or execute statements )
// stacks in blockly are.... interesting.
// they are a nested linked list of values.
// this writes the inner halves of the list
func newStack(m *chart.Machine, blk *js.Builder) *chart.StateMix {
	// the whole chain is going to be encapsulated by object braces {}
	// we try to keep the same state going for as long as we can...
	var cnt int
	var writingSlot bool
	open, close := js.Obj[0], js.Obj[1]
	return &chart.StateMix{
		OnMap: func(typeName string, _ jsn.FlowBlock) bool {
			cnt++ // we increment count here, eliding any empty slots.
			blk.R(js.Comma).Q("next").R(js.Colon).R(open).
				Q("block").R(js.Colon).R(open)
			m.PushState(newInnerBlock(m, blk, typeName))
			return true
		},
		OnSlot: func(string, jsn.SlotBlock) (okay bool) {
			writingSlot = true
			return true
		},
		OnEnd: func() {
			// we dont enter a new state for "OnSlot"
			// so we get ends for it and for the end of our own repeat.
			if writingSlot {
				writingSlot = false
			} else {
				for i := 0; i < cnt*2; i++ {
					blk.R(close)
				}
				m.FinishState(nil)
			}
		},
	}
}
