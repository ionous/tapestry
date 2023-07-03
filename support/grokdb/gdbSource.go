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
