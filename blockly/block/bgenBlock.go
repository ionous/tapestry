package block

import (
	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// writes a new block into what might be the topLevel array of blocks,
// or the value of a block or shadow key.
func NewTopBlock(cm *chart.Machine, types bconst.Types, blk *js.Builder, _zeroPos bool) chart.State {
	m := bgen{cm, types}
	open, close := js.Obj[0], js.Obj[1]
	zeroPos = _zeroPos
	return &chart.StateMix{
		// initially our own block name
		// later, a member of our flow
		OnMap: func(typeName string, flow jsn.FlowBlock) (okay bool) {
			blk.R(open)
			m.PushState(m.newInnerFlow(blk, typeName))
			return true
		},
		// listen to the end of the inner block
		OnCommit: func(interface{}) {
			blk.R(close)
			m.FinishState(nil)
		},
	}
}

type bgen struct {
	*chart.Machine
	types bconst.Types
}

var zeroPos = false

// writes most of the contents of a block, without its surrounding {}
// ( to support the nested linked lists of blocks used for stacks )
func (m *bgen) newInnerFlow(body *js.Builder, typeName string) *chart.StateMix {
	return m.newInnerBlock(body, typeName, true)
}

func (m *bgen) newInnerBlock(body *js.Builder, typeName string, allowExtraData bool) *chart.StateMix {
	var markup map[string]any
	if ptr := m.Markout; ptr != nil {
		markup, m.Markout = *ptr, nil
	}
	//
	var term string // set per key
	blk := blockData{id: NewId(), typeName: typeName, allowExtraData: allowExtraData, markup: markup, zeroPos: zeroPos}
	zeroPos = false
	return &chart.StateMix{
		// one of every extant member of the flow ( the encoder skips optional elements lacking a value )
		// this might be a field or input
		// we might write to next when the block is *followed* by another in a repeat.
		// therefore we cant close the block in Commit --
		// but we might close child blocks
		OnKey: func(_ string, field string) (noerr error) {
			term = field[1:] // strip off the $
			return
		},

		// a member that is a flow.
		OnMap: func(_ string, flow jsn.FlowBlock) (alwaysTrue bool) {
			was := blk.startInput(term)
			inner := m.newInnerFlow(&blk.inputs, flow.GetType())
			m.PushState(inner)
			defaultEnding := inner.OnEnd
			inner.OnEnd = func() {
				defaultEnding() // flushes its block, so call first
				blk.endInput(was)
			}
			return true
		},

		// a value that fills a slot; this will be an input
		OnSlot: func(string, jsn.SlotBlock) (alwaysTrue bool) {
			m.PushState(m.newSlot(term, &blk))
			return true
		},

		// a member that's a swap
		OnSwap: func(_ string, swap jsn.SwapBlock) (alwaysTrue bool) {
			m.PushState(m.newSwap(term, swap, &blk))
			return true
		},

		// a member that repeats
		OnRepeat: func(_ string, slice jsn.SliceBlock) (okay bool) {
			if cnt := slice.GetSize(); cnt > 0 {
				blk.writeCount(term, cnt)
				m.PushState(m.newRepeat(term, &blk))
				okay = true
			}
			return
		},

		// a single value
		OnValue: func(_ string, pv interface{}) error {
			return blk.writeValue(term, pv)
		},

		// end of the inner block
		OnEnd: func() {
			blk.writeTo(body)
			m.FinishState(nil)
		},
	}
}
