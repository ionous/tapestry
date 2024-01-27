package block

import (
	"git.sr.ht/~ionous/tapestry/lang/inspect"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// writes a fields ( in dummy inputs ) representing a repeating set of primitives.
// "fields": { "VALUES0": "a", "VALUES1": "b"... }
func (m *bgen) newList(term string, fields *js.Builder) inspect.Callbacks {
	var cnt int
	return inspect.Callbacks{
		OnValue: func(w inspect.It) (err error) {
			if b, e := valueToBytes(w); e != nil {
				err = e
			} else {
				if fields.Len() > 0 {
					fields.R(js.Comma)
				}
				fields.Brace(js.Quotes, func(q *js.Builder) { q.Str(term).N(cnt) }).
					R(js.Colon).
					Write(b)
				cnt++
			}
			return
		},
		// we dont enter a new state for "OnValue".. but values dont have a matching End.
		// we only get the end of our own repeat.
		OnEnd: func(w inspect.It) error {
			return m.events.Pop()
		},
	}
}
