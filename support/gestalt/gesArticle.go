package gestalt

import "git.sr.ht/~ionous/tapestry/support/grok"

// optional article
type Article struct {
	Optional bool
	Ignore   bool
}

func (op *Article) Match(q Query, cs []InputState) (ret []InputState) {
	for _, in := range cs {
		ws := in.Words()
		if m, width := q.FindArticle(ws); width == 0 {
			if op.Optional {
				ret = append(ret, in)
			}
		} else {
			out := in.Next(width)
			if !op.Ignore {
				out.AddResult(m)
			}
			ret = append(ret, out)
		}
	}
	return
}

type Count struct{}

func (op *Count) Match(q Query, cs []InputState) (ret []InputState) {
	for _, in := range cs {
		ws := in.Words()
		if v, ok := grok.WordsToNum(ws[0].String()); ok && v > 0 {
			const width = 1
			out := in.Next(width)
			out.AddResult(grok.Article{
				// fix? grok excludes the original text in its match
				// Match: grok.Span(ws[:width]),
				Count: v,
			})
			ret = append(ret, out)
		}
	}
	return
}
