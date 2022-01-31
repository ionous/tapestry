package bgen

import "git.sr.ht/~ionous/tapestry/web/js"

// helper for writing the inner bits of a blockly block definition
// ( ex. anything inside of "block" or "shadow" keys )
type blockData struct {
	typeName                         string
	fields, inputs, next, extraState js.Builder
}

// start a new input connection from this block to some new block
func (blk *blockData) startInput(term string) int {
	open := js.Obj[0]
	oldPos := blk.inputs.Len()
	blk.inputs.
		Q(term).R(js.Colon).
		R(open).Q("block").R(js.Colon).R(open)
	blk.writeCount(term, 1)
	return oldPos
}

// end a previous startInput connection.
func (blk *blockData) endInput(was int) {
	close := js.Obj[1]
	if now := blk.inputs.Len(); now > was {
		blk.inputs.R(close, close)
	}
}

// fields are named the same as the input
// see the tapestry_generic_mixin, createInput javascript.
func (b *blockData) writeValue(term string, pv interface{}) (err error) {
	if v, e := valueToBytes(pv); e != nil {
		err = e
	} else {
		b.writeCount(term, 1)
		b.fields.Q(term).R(js.Colon).Write(v)
	}
	return
}

// note: without looking up the definitions and checking whether a field was optional
// we don't know whether extraState is required -- so, we just write it for everything.
func (b *blockData) writeCount(term string, cnt int) {
	if cnt > 0 { // zero's the default for the tapestry_generic_mixin getExtraState()
		if b.extraState.Len() > 0 {
			b.extraState.R(js.Comma)
		}
		b.extraState.Q(term).R(js.Colon).N(cnt)
	}
}

func (b *blockData) writeTo(out *js.Builder) {
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

// helper to write a key:object where the object {} contains some arbitrary contents.
func writeContents(out *js.Builder, key string, contents *js.Builder) {
	out.R(js.Comma).Q(key).R(js.Colon).Brace(js.Obj, func(out *js.Builder) {
		out.S(contents.String())
	})
}
