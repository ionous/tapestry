package block

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/lang/inspect"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// writes a new block into what might be the topLevel array of blocks,
// or the value of a block or shadow key.
func Build(out *js.Builder, src typeinfo.Instance, _zeroPos bool) (err error) {
	zeroPos = _zeroPos
	// this setup is a little odd because we push to handle the visit
	// and then push to handle the src flow...
	var m bgen
	m.events.Push(inspect.Callbacks{
		OnFlow: func(w inspect.It) error {
			typeName := w.TypeInfo().(*typeinfo.Flow).Name
			return m.events.Push(m.newInnerFlow(w, out, typeName))
		},
		OnEnd: func(inspect.It) error {
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
func (m *bgen) newInnerFlow(w inspect.It, body *js.Builder, typeName string) inspect.Callbacks {
	return m.newInnerBlock(w, body, typeName, true)
}

func (m *bgen) newInnerBlock(w inspect.It, body *js.Builder, typeName string, allowExtraData bool) inspect.Callbacks {
	var term string // set per field, in CAP_UNDERSCORE format
	blk := blockData{
		id:             NewId(),
		typeName:       typeName,
		allowExtraData: allowExtraData,
		markup:         w.Markup(false), // for comments
		zeroPos:        zeroPos,
	}
	zeroPos = false
	return inspect.Callbacks{
		// one of every extant member of the flow ( the encoder skips optional elements lacking a value )
		// this might be a field or input
		// we might write to next when the block is *followed* by another in a repeat.
		// therefore we cant close the block in Commit --
		// but we might close child blocks
		OnField: func(w inspect.It) (_ error) {
			t := w.Term()
			term = strings.ToUpper(t.Name)
			return
		},

		// a member that is a flow.
		OnFlow: func(w inspect.It) error {
			was := blk.startInput(term)
			typeName := w.TypeInfo().(*typeinfo.Flow).Name
			return m.events.Push(
				inspect.OnEnd(m.newInnerFlow(w, &blk.inputs, typeName),
					func(w inspect.It, err error) error {
						blk.endInput(was)
						return err
					}))
		},

		// a value that fills a slot; this will be an input
		OnSlot: func(w inspect.It) (err error) {
			if cnt := w.Len(); cnt == 0 {
				m.skip()
			} else {
				err = m.events.Push(m.newSlot(term, &blk))
			}
			return
		},

		// a member that repeats
		OnRepeat: func(w inspect.It) (_ error) {
			if cnt := w.Len(); cnt == 0 {
				m.skip()
			} else {
				blk.writeCount(term, cnt)
				m.events.Push(m.newRepeat(term, &blk))
			}
			return
		},

		// a single value
		OnValue: func(w inspect.It) (err error) {
			if f := w.Term(); !f.Optional || !w.IsZero() {
				err = blk.writeValue(term, w)
			}
			return
		},

		// end of the inner block
		OnEnd: func(inspect.It) error {
			blk.writeTo(body)
			return m.events.Pop()
		},
	}
}

func (m *bgen) skip() {
	m.events.Push(inspect.Callbacks{
		OnEnd: func(inspect.It) (_ error) {
			return m.events.Pop()
		},
	})
}
