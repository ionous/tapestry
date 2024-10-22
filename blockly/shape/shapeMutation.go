package shape

import (
	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// write the mutation block used by the 'tapestry_generic_mutation' mutator
// see also: writeMuiData
func (w *ShapeWriter) writeMutator(out *js.Builder, flow *typeinfo.Flow) *js.Builder {
	// 1. write header
	//  "type": "_text_value_mutator",
	//  "style": "logic_blocks",
	//  "inputsInline": false,
	// 2. write args and message.
	return out.Brace(js.Obj, func(block *js.Builder) {
		block.
			Q("type").
			R(js.Colon).
			Q(bconst.MutatorName(flow.TypeName())).
			R(js.Comma).
			Kv("style", "logic_blocks").R(js.Comma).
			Q("inputsInline").R(js.Colon).Raw("false").R(js.Comma).
			If(true, func(args *js.Builder) {
				w.writeMuiMsgArgs(args, flow)
			})
	})
}
