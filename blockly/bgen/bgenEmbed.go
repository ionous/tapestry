package bgen

// // see also: newSwap
// func newEmbed(m *chart.Machine, term string, blk *blockData) *chart.StateMix {
// 	open, close := js.Obj[0], js.Obj[1]
// 	was := blk.inputs.Len()
// 	return &chart.StateMix{
// 		OnBlock: func(block jsn.Block) (err error) {
// 			blk.inputs.Q(term).R(js.Colon).
// 				R(open).Q("block").R(js.Colon).R(open)
// 			blk.writeCount(term, 1)
// 			m.PushState(newInnerBlock(m, &blk.inputs, block.GetType()))
// 			// return true
// 			return
// 		},
// 		OnEnd: func() {
// 			if now := blk.inputs.Len(); was != now {
// 				blk.inputs.R(close, close)
// 			}
// 			m.FinishState(nil)
// 		},
// 	}
// }
