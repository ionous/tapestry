package block

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/web/js"
	"github.com/ionous/errutil"
)

// we are basically the same as slot
// we just write a choice of the swap to the block's fields
func (m *bgen) newSwap(term string, swap jsn.SwapBlock, blk *blockData) chart.State {
	was := -1
	return &chart.StateMix{
		OnBlock: func(block jsn.Block) (err error) {
			if was, err = m.openSwap(blk, term, swap); err == nil {
				_, ok := block.(jsn.FlowBlock)
				m.PushState(m.newInnerBlock(&blk.inputs, block.GetType(), ok))
			}
			return
		},
		// for values, we have to manufacture a standalone block
		// this matches the btypes swap/standalone setup:
		// an input for every swap.
		OnValue: func(typeName string, pv interface{}) (err error) {
			if was, err = m.openSwap(blk, term, swap); err == nil {
				field := strings.ToUpper(typeName) // see: shape.writeStandalone
				faux := blockData{id: NewId(), typeName: typeName}
				if e := faux.writeValue(field, pv); e != nil {
					err = e
				} else {
					faux.writeTo(&blk.inputs)
				}
			}
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

func (m *bgen) openSwap(blk *blockData, term string, swap jsn.SwapBlock) (ret int, err error) {
	if choice, _ := swap.GetSwap(); len(choice) == 0 {
		err = errutil.New("expected valid choice")
	} else {
		// write the block field:
		// note: choice is $NAME style.
		if blk.fields.Len() > 0 {
			blk.fields.R(js.Comma)
		}
		blk.fields.Kv(term, choice)
		// write the opening of the connecting block
		ret = blk.startInput(term)
	}
	return
}
