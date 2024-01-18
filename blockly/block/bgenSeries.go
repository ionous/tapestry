package block

import (
	"git.sr.ht/~ionous/tapestry/lang/walk"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// writes a list of inputs representing a repeating set of slots.
// unlike stacks, repeated inputs are all in the same block.
// "inputs": { "CONTAINS0": {"block":{...}}, "CONTAINS1": {"block":{...}}, ... }
func (m *bgen) newSeries(term string, inputs *js.Builder) walk.Callbacks {
	open, close := js.Obj[0], js.Obj[1]
	var cnt int
	var writingSlot bool
	return walk.Callbacks{
		OnSlot: func(w walk.Walker) (_ error) {
			cnt++ // we count every slot, even if there is no block filling it.
			writingSlot = true
			return
		},
		OnFlow: func(w walk.Walker) error {
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
			return m.events.Push(walk.OnEnd(m.newInnerFlow(w, inputs, w.TypeName()),
				// when a child ( the inner block ) has finished
				func(w walk.Walker, err error) error {
					if err == nil {
						inputs.R(close, close)
					}
					return err
				}))
		},
		OnEnd: func(w walk.Walker) (err error) {
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
