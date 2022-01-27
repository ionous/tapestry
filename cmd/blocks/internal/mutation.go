package blocks

import "git.sr.ht/~ionous/tapestry/dl/spec"

// write the mutation block used by the 'tapestry_generic_mutation' mutator
// see also: writeMuiData
func writeMutator(out *Js, blockType *spec.TypeSpec, flow *spec.FlowSpec) *Js {
	// 1. write header
	//  "type": "_text_value_mutator",
	//  "style": "logic_blocks",
	//  "inputsInline": false,
	// 2. write args and message.
	return out.Brace(obj, func(block *Js) {
		block.
			Q("type").R(colon, quote, score).S(blockType.Name).S("_mutator").R(quote, comma).
			Kv("style", "logic_blocks").R(comma).
			Q("inputsInline").R(colon).S("false").R(comma).
			If(true, func(args *Js) {
				writeMuiMsgArgs(args, blockType, flow)
			})
	})
}
