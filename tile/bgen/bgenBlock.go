package bgen

import (
	"encoding/json"
	"strings"

	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// about to start a new flow; a fresh block that has no influence on parent blocks
// "blocks" might be the topLevel array of blocks, or the value of a block or shadow key.
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

func newInnerBlock(m *chart.Machine, blk *js.Builder, typeName string) chart.State {
	var fields, inputs, next js.Builder
	var nextKey string
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
		OnKey: func(key string, field string) (err error) {
			field = field[1:] // strip off $
			nextKey = field   // this is the name we want
			return
		},
		// a value that fills a slot
		OnSlot: func(string, jsn.SlotBlock) (okay bool) { return },
		// a member that's a swap
		OnSwap: func(string, jsn.SwapBlock) (okay bool) { return },
		// a member that repeats
		OnRepeat: func(string, jsn.SliceBlock) (okay bool) {
			// fix: determine if this is a stack or a series of inputs within the same block
			// note: slice actually gives us the depth if that's helpful
			m.PushState(newStack(m, &next))
			return true
		},
		// a single value
		OnValue: func(_ string, pv interface{}) (err error) {
			if b, e := json.Marshal(unpackValue(pv)); e != nil {
				err = e
			} else {
				// fields are lower case ( inputs are upper case )
				// fix: probably a distinction without meaning.
				name := strings.ToLower(nextKey)
				fields.Q(name).R(js.Colon).Write(b)
			}
			return
		},

		OnEnd: func() {
			blk.Kv("type", typeName).R(js.Comma)
			blk.Q("extraState").R(js.Colon).Brace(js.Obj, func(out *js.Builder) {
				// it seems the extra state blk is required for loading to work
				// even if there's nothing to put in it.
			})
			writeContents(blk, "fields", &fields)
			writeContents(blk, "inputs", &inputs)
			blk.If(next.Len() > 0, func(out *js.Builder) {
				out.S(next.String())
			})
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

func writeContents(out *js.Builder, key string, contents *js.Builder) {
	if contents.Len() > 0 {
		out.R(js.Comma).Q(key).R(js.Colon).Brace(js.Obj, func(out *js.Builder) {
			out.S(contents.String())
		})
	}
}
