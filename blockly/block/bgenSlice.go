package block

import (
	"git.sr.ht/~ionous/tapestry/lang/walk"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// writes a slice of repeating flows.
// unlike stacks, repeated inputs are all in the same block.
// ( ex. "inputs": { "CONTAINS0": {...}, "CONTAINS1": {...}, ... } )
func (m *bgen) newSlice(term string, inputs *js.Builder) walk.Callbacks {
	open, close, cnt := js.Obj[0], js.Obj[1], 0
	return walk.Callbacks{
		OnFlow: func(w walk.Walker) error {
			if inputs.Len() > 0 {
				inputs.R(js.Comma)
			}
			// write `"TERM#": {"block":{`
			inputs.Brace(js.Quotes, func(q *js.Builder) {
				q.Str(term).N(cnt)
			}).R(js.Colon).R(open).
				Q("block").R(js.Colon).R(open)
			cnt++
			return m.events.Push(
				walk.OnEnd(m.newInnerFlow(w, inputs, w.TypeName()),
					func(w walk.Walker, err error) error {
						if err == nil {
							inputs.R(close, close)
						}
						return err
					}))
		},
		// the end of the repeat block which started us.
		OnEnd: func(w walk.Walker) error {
			return m.events.Pop()
		},
	}
}
