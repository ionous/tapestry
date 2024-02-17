package jess

import (
	"git.sr.ht/~ionous/tapestry/support/grok"
)

func (op *Article) Match(q Query, input *InputState) (okay bool) {
	if m, width := q.FindArticle(*input); width > 0 {
		op.Matched, *input, okay = m, input.Skip(width), true
	}
	return
}

func ReduceArticle(op *Article) (ret grok.Article) {
	if op != nil {
		ret = grok.Article{
			Matched: op.Matched,
		}
	}
	return
}
