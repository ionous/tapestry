package bgen

import (
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// writes a fields ( in dummy inputs ) representing a repeating set of primitives.
// "fields": { "VALUES0": "a", "VALUES1": "b"... }
func newList(m *chart.Machine, term string, fields *js.Builder) *chart.StateMix {
	var cnt int
	return &chart.StateMix{
		OnValue: func(_ string, pv interface{}) (err error) {
			if b, e := valueToBytes(pv); e != nil {
				err = e
			} else {
				if fields.Len() > 0 {
					fields.R(js.Comma)
				}
				fields.Brace(js.Quotes, func(q *js.Builder) { q.S(term).N(cnt) }).
					R(js.Colon).
					Write(b)
				cnt++
			}
			return
		},
		OnEnd: func() {
			m.FinishState(nil)
		},
	}
}
