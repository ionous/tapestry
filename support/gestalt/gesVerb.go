package gestalt

import "git.sr.ht/~ionous/tapestry/support/grok"

// matches a specified string and advances the input
type Verb struct {
	Str    string
	Action string
	matcher
}

func (op *Verb) Match(q Query, cs []InputState) (ret []InputState) {
	if match, e := op.makeSpan(op.Str); e == nil {
		for _, in := range cs {
			ws := in.Words()
			if len(match) == 0 {
				if m, width := q.FindMacro(ws); width > 0 {
					out := in.Next(width)
					out.AddResult(m)
					ret = append(ret, out)
				}

			} else if grok.HasPrefix(ws, match) {
				width := len(match)
				span := grok.Span(ws[:width])
				out := in.Next(width)
				out.AddResult(MatchedVerb{span, op.Action})
				ret = append(ret, out)
			}
		}

	}
	return
}

type MatchedVerb struct {
	grok.Span
	Action string
}
