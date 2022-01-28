package bgen

import (
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// some tapestry slices of slots are represented as a blockly statement stack.
// ( for example: story or execute statements )
// stacks in blockly are.... interesting.
// they are a nested linked list of values.
// this writes the inner halves of the list
func newStack(m *chart.Machine, blk *js.Builder) chart.State {
	// the whole chain is going to be encapsulated by object braces {}
	// we try to keep the same state going for as long as we can...
	open, close := js.Obj[0], js.Obj[1]
	// var waitingOnBlock bool
	var cnt, end int // alt: might read the .Len() of the slice
	return &chart.StateMix{
		OnMap: func(typeName string, _ jsn.FlowBlock) bool {
			m.PushState(newInnerBlock(m, blk, typeName))
			// waitingOnBlock = false
			return true
		},
		OnSlot: func(string, jsn.SlotBlock) (okay bool) {
			// this could (theoretically) happen with nil slots
			// if waitingOnBlock {
			// 	blk.S("{}") // how to correctly handle this?
			// 	log.Println("skipped nil slot in stack")
			// }
			blk.R(js.Comma).Q("next").R(js.Colon).R(open).
				Q("block").R(js.Colon).R(open)
			// the contents will (usually) be written by OnMap
			// waitingOnBlock = true
			cnt++
			return true
		},
		OnEnd: func() {
			if end = end + 1; end > cnt {
				for i := 0; i < cnt*2; i++ {
					blk.R(close)
				}
				m.FinishState(nil)
			}
		},
	}
}
