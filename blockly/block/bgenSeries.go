package block

import (
	"git.sr.ht/~ionous/tapestry/lang/inspect"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// writes a list of inputs representing a repeating set of slots.
// unlike stacks, repeated inputs are all in the same block.
// "inputs": { "CONTAINS0": {"block":{...}}, "CONTAINS1": {"block":{...}}, ... }
func (m *bgen) newSeries(term string, inputs *js.Builder) inspect.Callbacks {
	open, close := js.Obj[0], js.Obj[1]
	var cnt int
	var writingSlot bool
	return inspect.Callbacks{
		OnSlot: func(w inspect.Iter) (_ error) {
			cnt++ // we count every slot, even if there is no block filling it.
			writingSlot = true
			return
		},
		OnFlow: func(w inspect.Iter) error {
			if inputs.Len() > 0 {
				inputs.R(js.Comma)
			}
			// writes: `"term#"`:{"block":{`
			inputs.
				Brace(js.Quotes, func(q *js.Builder) {
					q.Str(term).N(cnt - 1)
				}).
				R(js.Colon).R(open).
				Q("block").
				R(js.Colon).R(open)

			typeName := w.TypeInfo().(*typeinfo.Flow).Name
			return m.events.Push(inspect.OnEnd(m.newInnerFlow(w, inputs, typeName),
				// when a child ( the inner block ) has finished
				func(w inspect.Iter, err error) error {
					if err == nil {
						inputs.R(close, close)
					}
					return err
				}))
		},
		OnEnd: func(w inspect.Iter) (err error) {
			// note: we reuse the current state for each "OnSlot"
			// so we get ends for it and for the end of our own repeat.
			if writingSlot {
				writingSlot = false
			} else {
				err = m.events.Pop()
			}
			return
		},
	}
}
