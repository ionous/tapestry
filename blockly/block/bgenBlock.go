package block

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/lang/inspect"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// writes a new block into what might be the topLevel array of blocks,
// or the value of a block or shadow key.
func Build(out *js.Builder, src typeinfo.Inspector, _zeroPos bool) (err error) {
	zeroPos = _zeroPos

	var m bgen
	m.events.Push(inspect.Callbacks{
		OnFlow: func(w inspect.Iter) error {
			typeName := w.TypeInfo().(*typeinfo.Flow).Name
			return m.events.Push(m.newInnerFlow(w, out, typeName))
		},
		OnEnd: func(inspect.Iter) error {
			return m.events.Pop()
		},
	})

	out.R(js.Obj[0])
	if e := inspect.Visit(src, &m.events); e != nil {
		err = e
	} else {
		out.R(js.Obj[1])
	}
	return
}

type bgen struct {
	events inspect.Stack
}

var zeroPos = false

// writes most of the contents of a block, without its surrounding {}
// ( to support the nested linked lists of blocks used for stacks )
func (m *bgen) newInnerFlow(w inspect.Iter, body *js.Builder, typeName string) inspect.Callbacks {
	return m.newInnerBlock(w, body, typeName, true)
}

func (m *bgen) newInnerBlock(w inspect.Iter, body *js.Builder, typeName string, allowExtraData bool) inspect.Callbacks {
	var term string // set per field, in CAP_UNDERSCORE format
	blk := blockData{
		id:             NewId(),
		typeName:       typeName,
		allowExtraData: allowExtraData,
		markup:         w.Markup(),
		zeroPos:        zeroPos,
	}
	zeroPos = false
	return inspect.Callbacks{
		// one of every extant member of the flow ( the encoder skips optional elements lacking a value )
		// this might be a field or input
		// we might write to next when the block is *followed* by another in a repeat.
		// therefore we cant close the block in Commit --
		// but we might close child blocks
		OnField: func(w inspect.Iter) (_ error) {
			t := w.Term()
			term = strings.ToUpper(t.Name)
			return
		},

		// a member that is a flow.
		OnFlow: func(w inspect.Iter) error {
			was := blk.startInput(term)
			typeName := w.TypeInfo().(*typeinfo.Flow).Name
			return m.events.Push(
				inspect.OnEnd(m.newInnerFlow(w, &blk.inputs, typeName),
					func(w inspect.Iter, err error) error {
						blk.endInput(was)
						return err
					}))
		},

		// a value that fills a slot; this will be an input
		OnSlot: func(w inspect.Iter) error {
			return m.events.Push(m.newSlot(term, &blk))
		},

		// a member that repeats
		OnRepeat: func(w inspect.Iter) (_ error) {
			if cnt := w.Len(); cnt > 0 {
				blk.writeCount(term, cnt)
				m.events.Push(m.newRepeat(term, &blk))
			}
			return
		},

		// a single value
		OnValue: func(w inspect.Iter) (err error) {
			if f := w.Term(); !f.Optional || !w.ZeroValue() {
				err = blk.writeValue(term, w)
			}
			return
		},

		// end of the inner block
		OnEnd: func(inspect.Iter) error {
			blk.writeTo(body)
			return m.events.Pop()
		},
	}
}
