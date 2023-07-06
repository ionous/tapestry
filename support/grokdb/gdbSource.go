package grokdb

import (
	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/tables"
)

type Source struct {
	inner dbSource
}

func NewSource(db *tables.Cache) Source {
	var x Source
	x.inner.db = db
	return x
}

func (x *Source) Grok(domain, phrase string) (grok.Results, error) {
	x.inner.domain = domain
	return grok.Grok(&x.inner, phrase)
}

func (x *Source) MatchArticle(ws []string) (ret int, err error) {
	if len(ws) > 0 {
		if s, e := grok.MakeSpan(ws[0]); e != nil {
			err = e
		} else if m, e := det.FindMatch(s); e != nil {
			err = e
		} else if m != nil {
			ret = m.NumWords()
		}
	}
	return
}
