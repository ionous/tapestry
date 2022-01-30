package bgen

import (
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// writes a slice of repeating flows.
// unlike stacks, repeated inputs are all in the same block.
// ( ex. "inputs": { "CONTAINS0": {...}, "CONTAINS1": {...}, ... } )
func newSlice(m *chart.Machine, term string, inputs *js.Builder) *chart.StateMix {
	open, close, cnt := js.Obj[0], js.Obj[1], 0
	return &chart.StateMix{
		OnMap: func(typeName string, flow jsn.FlowBlock) bool {
			if inputs.Len() > 0 {
				inputs.R(js.Comma)
			}
			// write `"TERM#": {"block":{`
			inputs.Brace(js.Quotes, func(q *js.Builder) {
				q.S(term).N(cnt)
			}).R(js.Colon).R(open).
				Q("block").R(js.Colon).R(open)
			cnt++
			m.PushState(newInnerBlock(m, inputs, typeName))
			return true
		},
		// when a child state ( the inner block ) has finished
		OnCommit: func(interface{}) {
			inputs.R(close, close)
		},
		// the end of the repeat block which started us.
		OnEnd: func() {
			m.FinishState(nil)
		},
	}
}
