package bgen

import (
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
)

// see also: newSwap
func newSlot(m *chart.Machine, term string, blk *blockData) *chart.StateMix {
	was := -1
	return &chart.StateMix{
		OnBlock: func(block jsn.Block) (err error) {
			was = blk.startInput(term)
			m.PushState(newInnerBlock(m, &blk.inputs, block.GetType()))
			return
		},
		OnEnd: func() {
			if was >= 0 {
				blk.endInput(was)
				was = -1
			}
			m.FinishState(nil)
		},
	}
}
