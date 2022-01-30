package bgen

import (
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/web/js"
)

func newSlot(m *chart.Machine, term string, blk *blockData) chart.State {
	open, close := js.Obj[0], js.Obj[1]
	var slotExists bool
	return &chart.StateMix{
		OnMap: func(typeName string, _ jsn.FlowBlock) bool {
			if !slotExists {
				// "TEXT": { "block": { ....
				blk.inputs.Q(term).R(js.Colon).
					R(open).Q("block").R(js.Colon).R(open)
				blk.writeCount(term, 1)
				slotExists = true
			}
			m.PushState(newInnerBlock(m, &blk.inputs, typeName))
			// waitingOnBlock = false
			return true
		},
		OnEnd: func() {
			if slotExists {
				blk.inputs.R(close, close)
			}
			m.FinishState(nil)
		},
	}
}
