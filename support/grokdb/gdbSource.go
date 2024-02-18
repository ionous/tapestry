package grokdb

import (
	"git.sr.ht/~ionous/tapestry/dl/jess"
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

func (x *Source) MatchSpan(domain string, span grok.Span) (jess.Applicant, error) {
	x.inner.domain = domain
	return jess.Match(&x.inner, span)
}

func (x *Source) MatchArticle(ws []string) (ret int, err error) {
	if len(ws) > 0 {
		// assumes all articles are one word.
		if s, e := grok.MakeSpan(ws[0]); e != nil {
			err = e
		} else if m, e := grok.FindCommonArticles(s); e != nil {
			err = e
		} else if m != nil {
			ret = m.NumWords()
		}
	}
	return
}
