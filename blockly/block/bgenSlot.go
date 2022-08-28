package block

import (
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
)

// see also: newSwap
func (m *bgen) newSlot(term string, blk *blockData) *chart.StateMix {
	was := -1
	return &chart.StateMix{
		OnBlock: func(block jsn.Block) (err error) {
			was = blk.startInput(term)
			_, ok := block.(jsn.FlowBlock)
			m.PushState(m.newInnerBlock(&blk.inputs, block.GetType(), ok))
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
