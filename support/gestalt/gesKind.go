package gestalt

import "git.sr.ht/~ionous/tapestry/support/grok"

// matches one pre-existing kind.
type Kind struct{}

func (*Kind) Match(q Query, cs []InputState) (ret []InputState) {
	for _, in := range cs {
		ws := in.Words()
		if _, width := q.FindKind(ws); width > 0 {
			out := in.Next(width)
			out.AddResult(MatchedKind{grok.Span(ws[:width])})
			ret = append(ret, out)
		}
	}
	return
}

type MatchedKind struct {
	grok.Span
}
