package bgen

import (
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/web/js"
)

func newSlot(m *chart.Machine, blk *js.Builder) chart.State {
	open, close := js.Obj[0], js.Obj[1]
	blk.R(open).Q("block").R(js.Colon).R(open)
	return &chart.StateMix{
		OnMap: func(typeName string, _ jsn.FlowBlock) bool {
			m.PushState(newInnerBlock(m, blk, typeName))
			// waitingOnBlock = false
			return true
		},
		OnEnd: func() {
			blk.R(close, close)
			m.FinishState(nil)
		},
	}
}
