package bgen

import (
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// writes a new block into what might be the topLevel array of blocks,
// or the value of a block or shadow key.
func NewTopBlock(m *chart.Machine, blk *js.Builder) chart.State {
	open, close := js.Obj[0], js.Obj[1]
	return &chart.StateMix{
		// initially our own block name
		// later, a member of our flow
		OnMap: func(typeName string, flow jsn.FlowBlock) (okay bool) {
			blk.R(open)
			m.PushState(newInnerBlock(m, blk, typeName))
			return true
		},
		// listen to the end of the inner block
		OnCommit: func(interface{}) {
			blk.R(close)
			m.FinishState(nil)
		},
	}
}

// writes most of the contents of a block, without its surrounding {}
// ( to support the nested linked lists of blocks used for stacks )
func newInnerBlock(m *chart.Machine, body *js.Builder, typeName string) *chart.StateMix {
	var term string // set per key
	blk := blockData{typeName: typeName}
	return &chart.StateMix{
		// a member that is a flow.
		OnMap: func(_ string, flow jsn.FlowBlock) bool {
			was := blk.startInput(term)
			next := newInnerBlock(m, &blk.inputs, flow.GetType())
			m.PushState(next)
			prev := next.OnEnd
			next.OnEnd = func() {
				prev() // flushes its block, so call first
				blk.endInput(was)
			}
			return true
		},

		// one of every extant member of the flow ( skipping optional elements lacking a value )
		// this might be a field or input
		// we might write to next when the block is *followed* by another in a repeat.
		// therefore we cant close the block in Commit --
		// but we might close child blocks
		OnKey: func(_ string, field string) (err error) {
			term = field[1:] // strip off the $
			return
		},

		// a value that fills a slot; this will be an input
		OnSlot: func(string, jsn.SlotBlock) (alwaysTrue bool) {
			m.PushState(newSlot(m, term, &blk))
			return true
		},

		// a member that's a swap
		OnSwap: func(_ string, swap jsn.SwapBlock) (alwaysTrue bool) {
			m.PushState(newSwap(m, term, swap, &blk))
			return true
		},

		// a member that repeats
		OnRepeat: func(_ string, slice jsn.SliceBlock) (alwaysTrue bool) {
			if cnt := slice.GetSize(); cnt > 0 {
				blk.writeCount(term, cnt)
				m.PushState(newRepeat(m, term, &blk))
			}
			return true
		},

		// a single value
		OnValue: func(_ string, pv interface{}) (err error) {
			return blk.writeValue(term, pv)
		},

		// end of the inner block
		OnEnd: func() {
			blk.writeTo(body)
			m.FinishState(nil)
		},
	}
}
