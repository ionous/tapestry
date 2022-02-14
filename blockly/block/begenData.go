package block

import "git.sr.ht/~ionous/tapestry/web/js"

// helper for writing the inner bits of a blockly block definition
// ( ex. anything inside of "block" or "shadow" keys )
type blockData struct {
	id, typeName               string
	fields, inputs, extraState js.Builder
	allowExtraData             bool
	comment                    string
}

// start a new input connection from this block to some new block
func (b *blockData) startInput(term string) int {
	ret := b.startInputWithoutCount(term)
	b.writeCount(term, 1)
	return ret
}

func (b *blockData) startInputWithoutCount(term string) int {
	open := js.Obj[0]
	oldPos := b.inputs.Len()
	if oldPos > 0 {
		b.inputs.R(js.Comma)
	}
	b.inputs.
		Q(term).R(js.Colon).
		R(open).Q("block").R(js.Colon).R(open)
	return oldPos
}

// end a previous startInput connection.
func (b *blockData) endInput(was int) {
	close := js.Obj[1]
	if now := b.inputs.Len(); now > was {
		b.inputs.R(close, close)
	}
}

// fields are named the same as the input
// see the tapestry_generic_mixin, createInput javascript.
func (b *blockData) writeValue(term string, pv interface{}) (err error) {
	if v, e := valueToBytes(pv); e != nil {
		err = e
	} else {
		if b.fields.Len() > 0 {
			b.fields.R(js.Comma)
		}
		b.fields.Q(term).R(js.Colon).Write(v)
		b.writeCount(term, 1)
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
	out.
		Kv("id", b.id).R(js.Comma).
		Kv("type", b.typeName)
	// comment, if any:
	if cmt := b.comment; len(cmt) > 0 {
		out.R(js.Comma).
			Q("icons").R(js.Colon).Brace(js.Obj, func(icons *js.Builder) {
			icons.Q("comment").R(js.Colon).Brace(js.Obj, func(com *js.Builder) {
				com.
					Kv("text", cmt).R(js.Comma).
					Q("pinned").R(js.Colon).S(js.False).R(js.Comma).
					Q("height").R(js.Colon).N(80).R(js.Comma).
					Q("width").R(js.Colon).N(280)
			})
		})
	}
	// note: blockly throws exceptions if you have extraData but not a loadExtraState function on the shape.
	// and it only allows you to have loadExtraState on shapes with mutations.
	// ( in loadConnection, appendPrivate, it tries to treat extraData as xml data when that function is missing )
	if b.allowExtraData {
		writeContents(out, "extraState", &b.extraState)
	}
	if els := &b.fields; els.Len() > 0 {
		writeContents(out, "fields", els)
	}
	if els := &b.inputs; els.Len() > 0 {
		writeContents(out, "inputs", els)
	}
}

// helper to write a key:object where the object {} contains some arbitrary contents.
func writeContents(out *js.Builder, key string, contents *js.Builder) {
	out.R(js.Comma).Q(key).R(js.Colon).Brace(js.Obj, func(out *js.Builder) {
		out.S(contents.String())
	})
}
