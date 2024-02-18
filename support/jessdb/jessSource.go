package jessdb

import (
	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/support/match"
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

func (x *Source) MatchSpan(domain string, span match.Span) (jess.Applicant, error) {
	x.inner.domain = domain
	return jess.Match(&x.inner, span)
}
