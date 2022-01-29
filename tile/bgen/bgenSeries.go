package bgen

import (
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// writes a list of inputs representing a repeating set of slots.
// unlike stacks, repeated inputs are all in the same block.
// "inputs": { "CONTAINS0": {...}, "CONTAINS1": {...}, ... }
func newSeries(m *chart.Machine, term string, inputs *js.Builder) *chart.StateMix {
	open, close := js.Obj[0], js.Obj[1]
	var cnt int
	var writingBlock bool
	flushBlock := func() {
		if writingBlock {
			inputs.R(close)
			writingBlock = false
		}
	}
	return &chart.StateMix{
		OnMap: func(typeName string, _ jsn.FlowBlock) bool {
			m.PushState(newInnerBlock(m, inputs, typeName))
			flushBlock()
			return true
		},
		OnSlot: func(string, jsn.SlotBlock) (okay bool) {
			if inputs.Len() > 0 {
				inputs.R(js.Comma)
			}
			inputs.Brace(js.Quotes, func(q *js.Builder) {
				inputs.S(term).N(cnt)
			}).R(js.Colon).R(open).
				Q("block").R(js.Colon).R(open)
			cnt++
			writingBlock = true
			return true
		},
		// child block
		OnCommit: func(interface{}) {
			flushBlock()
		},
		OnEnd: func() {
			flushBlock()
			m.FinishState(nil)
		},
	}
}
