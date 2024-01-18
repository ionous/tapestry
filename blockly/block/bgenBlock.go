package block

import (
	r "reflect"
	"strings"

	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/lang/walk"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// writes a new block into what might be the topLevel array of blocks,
// or the value of a block or shadow key.
func Build(out *js.Builder, src jsn.Marshalee, types bconst.Types, _zeroPos bool) (err error) {
	w := walk.Walk(r.ValueOf(src).Elem())
	zeroPos = _zeroPos

	m := bgen{types: types}
	m.events.Push(walk.Callbacks{
		OnFlow: func(w walk.Walker) error {
			typeName := w.TypeName()
			return m.events.Push(m.newInnerFlow(w, out, typeName))
		},
		OnEnd: func(walk.Walker) error {
			return m.events.Pop()
		},
	})

	out.R(js.Obj[0])
	if e := walk.VisitFlow(w, &m.events); e != nil {
		err = e
	} else {
		out.R(js.Obj[1])
	}
	return
}

type bgen struct {
	events walk.Stack
	types  bconst.Types
}

var zeroPos = false

// writes most of the contents of a block, without its surrounding {}
// ( to support the nested linked lists of blocks used for stacks )
func (m *bgen) newInnerFlow(w walk.Walker, body *js.Builder, typeName string) walk.Callbacks {
	return m.newInnerBlock(w, body, typeName, true)
}

func (m *bgen) newInnerBlock(w walk.Walker, body *js.Builder, typeName string, allowExtraData bool) walk.Callbacks {
	var term string // set per field, in CAP_UNDERSCORE format
	blk := blockData{
		id:             NewId(),
		typeName:       typeName,
		allowExtraData: allowExtraData,
		markup:         w.Markup().Interface().(map[string]any),
		zeroPos:        zeroPos,
	}
	zeroPos = false
	return walk.Callbacks{
		// one of every extant member of the flow ( the encoder skips optional elements lacking a value )
		// this might be a field or input
		// we might write to next when the block is *followed* by another in a repeat.
		// therefore we cant close the block in Commit --
		// but we might close child blocks
		OnField: func(w walk.Walker) (_ error) {
			f := w.Field()
			term = strings.ToUpper(f.FieldName())
			return
		},

		// a member that is a flow.
		OnFlow: func(w walk.Walker) error {
			was := blk.startInput(term)
			return m.events.Push(
				walk.OnEnd(m.newInnerFlow(w, &blk.inputs, w.TypeName()),
					func(w walk.Walker, err error) error {
						blk.endInput(was)
						return err
					}))
		},

		// a value that fills a slot; this will be an input
		OnSlot: func(w walk.Walker) error {
			return m.events.Push(m.newSlot(term, &blk))
		},

		// a member that repeats
		OnRepeat: func(w walk.Walker) (_ error) {
			if cnt := w.Len(); cnt > 0 {
				blk.writeCount(term, cnt)
				m.events.Push(m.newRepeat(term, &blk))
			}
			return
		},

		// a single value
		OnValue: func(w walk.Walker) (err error) {

			// fix: when going through marshal and the autognerated things
			// this will use things like Bool_Unboxed_Marshal
			//  n m.MarshalValue(Bool_Type, jsn.BoxBool(val))
			// which has GetValue() $TRUE $FALSE
			// and compact value true, value
			if v := w.RawValue(); v.IsValid() {
				if f := w.Field(); !f.Optional() || !v.IsZero() {
					// see valueToBytes which tries to use GetCompactValue if it eixsts
					// Text_Unboxed_Repeats_Marshal for repeating text ( collapsing an array of lines to one )
					// and Enum
					// i think this has to do all that maually --
					// except that type info should provide access to choices for strings
					// and bool in type info should have str choices too,
					err = blk.writeValue(term, v.Interface())
				}
			}
			return
		},

		// end of the inner block
		OnEnd: func(walk.Walker) error {
			blk.writeTo(body)
			return m.events.Pop()
		},
	}
}
