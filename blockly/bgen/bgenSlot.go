package bgen

import (
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
)

// see also: newSwap
func newSlot(m *chart.Machine, term string, blk *blockData) *chart.StateMix {
	var was int
	return &chart.StateMix{
		OnBlock: func(block jsn.Block) (err error) {
			was = blk.startInput(term)
			m.PushState(newInnerBlock(m, &blk.inputs, block.GetType()))
			return
		},
		OnEnd: func() {
			blk.endInput(was)
			m.FinishState(nil)
		},
	}
}
