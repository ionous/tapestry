package block

import (
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// writes a list of inputs representing a repeating set of slots.
// unlike stacks, repeated inputs are all in the same block.
// "inputs": { "CONTAINS0": {"block":{...}}, "CONTAINS1": {"block":{...}}, ... }
func newSeries(m *chart.Machine, term string, inputs *js.Builder) *chart.StateMix {
	open, close := js.Obj[0], js.Obj[1]
	var cnt int
	var writingSlot bool
	return &chart.StateMix{
		OnMap: func(typeName string, _ jsn.FlowBlock) bool {
			if inputs.Len() > 0 {
				inputs.R(js.Comma)
			}
			// writes: `"term#"`:{"block":{`
			inputs.
				Brace(js.Quotes, func(q *js.Builder) {
					q.X(term).N(cnt - 1)
				}).
				R(js.Colon).R(open).
				Q("block").
				R(js.Colon).R(open)
			m.PushState(newInnerFlow(m, inputs, typeName))
			return true
		},
		// when a child ( the inner block ) has finished
		OnCommit: func(interface{}) {
			inputs.R(close, close)
		},
		OnSlot: func(string, jsn.SlotBlock) (okay bool) {
			cnt++ // we count every slot, even if there is no block filling it.
			writingSlot = true
			return true
		},
		OnEnd: func() {
			// note: we reuse the current state for each "OnSlot"
			// so we get ends for it and for the end of our own repeat.
			if writingSlot {
				writingSlot = false
			} else {
				m.FinishState(nil)
			}
		},
	}
}
