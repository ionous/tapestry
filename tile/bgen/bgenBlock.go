package bgen

import (
	"encoding/json"

	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// writes a new block into what might be the topLevel array of blocks,
// or the value of a block or shadow key.
// FIX: handle user comments
func newTopBlock(m *chart.Machine, blk *js.Builder) chart.State {
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

// note: without looking up the definitions and checking whether a field was optional
// we don't know whether extraState is required -- so, we just write it for everything.
type blockData struct {
	typeName                         string
	fields, inputs, next, extraState js.Builder
}

// writes most of the contents of a block, without its surrounding {}
// ( to support the nested linked lists of blocks used for stacks )
func newInnerBlock(m *chart.Machine, blk *js.Builder, typeName string) chart.State {
	var term string // set per key
	data := blockData{typeName: typeName}
	return &chart.StateMix{
		// a member that is a flow.
		OnMap: func(string, jsn.FlowBlock) (okay bool) {
			return
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
		OnSlot: func(string, jsn.SlotBlock) bool {
			data.writeCount(term, 1)
			data.inputs.Q(term).R(js.Colon) // "TEXT": ....
			m.PushState(newSlot(m, &data.inputs))
			return true
		},

		// a member that's a swap
		OnSwap: func(string, jsn.SwapBlock) (okay bool) {
			// 			data.writeCount(term, 1)
			return
		},

		// a member that repeats
		OnRepeat: func(_ string, slice jsn.SliceBlock) (okay bool) {
			// FIX: test some zero sized arrays
			if cnt := slice.GetSize(); cnt > 0 {
				m.PushState(newRepeat(m, term, cnt, &data))
			}
			return true
		},

		// a single value
		OnValue: func(_ string, pv interface{}) (err error) {
			if b, e := json.Marshal(unpackValue(pv)); e != nil {
				err = e
			} else {
				// fields are named the same as the input
				// see the tapestry_generic_mixin, createInput javascript.
				data.writeCount(term, 1)
				data.fields.Q(term).R(js.Colon).Write(b)
			}
			return
		},

		// end of the inner block
		OnEnd: func() {
			data.write(blk)
			m.FinishState(nil)
		},
	}
}

func unpackValue(pv interface{}) (ret interface{}) {
	switch pv := pv.(type) {
	// case interface{ GetCompactValue() interface{} }:
	// 	ret = pv.GetCompactValue()
	case interface{ GetValue() interface{} }:
		ret = pv.GetValue()
	default:
		ret = pv
	}
	return
}

func (b *blockData) writeCount(term string, cnt int) {
	if b.extraState.Len() > 0 {
		b.extraState.R(js.Comma)
	}
	b.extraState.Q(term).R(js.Colon).N(cnt)
}

func (b *blockData) write(out *js.Builder) {
	out.Kv("type", b.typeName)
	// note: always have to write extraState or blockly gets unhappy...
	writeContents(out, "extraState", &b.extraState)
	if els := &b.fields; els.Len() > 0 {
		writeContents(out, "fields", els)
	}
	if els := &b.inputs; els.Len() > 0 {
		writeContents(out, "inputs", els)
	}
	if b.next.Len() > 0 {
		out.S(b.next.String())
	}
}

func writeContents(out *js.Builder, key string, contents *js.Builder) {
	out.R(js.Comma).Q(key).R(js.Colon).Brace(js.Obj, func(out *js.Builder) {
		out.S(contents.String())
	})
}
