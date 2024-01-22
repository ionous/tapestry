package block

import (
	"git.sr.ht/~ionous/tapestry/lang/inspect"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// writes a slice of repeating flows.
// unlike stacks, repeated inputs are all in the same block.
// ( ex. "inputs": { "CONTAINS0": {...}, "CONTAINS1": {...}, ... } )
func (m *bgen) newSlice(term string, inputs *js.Builder) inspect.Callbacks {
	open, close, cnt := js.Obj[0], js.Obj[1], 0
	return inspect.Callbacks{
		OnFlow: func(w inspect.Iter) error {
			if inputs.Len() > 0 {
				inputs.R(js.Comma)
			}
			// write `"TERM#": {"block":{`
			inputs.Brace(js.Quotes, func(q *js.Builder) {
				q.Str(term).N(cnt)
			}).R(js.Colon).R(open).
				Q("block").R(js.Colon).R(open)
			cnt++
			typeName := w.TypeInfo().(*typeinfo.Flow).Name
			return m.events.Push(
				inspect.OnEnd(m.newInnerFlow(w, inputs, typeName),
					func(w inspect.Iter, err error) error {
						if err == nil {
							inputs.R(close, close)
						}
						return err
					}))
		},
		// the end of the repeat block which started us.
		OnEnd: func(w inspect.Iter) error {
			return m.events.Pop()
		},
	}
}
