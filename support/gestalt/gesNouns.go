package gestalt

import "git.sr.ht/~ionous/tapestry/support/grok"

// matches one or more names, their determiners, and separators.
// because this could generate a large number of possibilities,
// there are some assumptions for optimization:
// the words "is/are/comma/and" are never part of noun names.
// fix: quotes for "titles" ( which are then allowed to break the rules )
type Name struct{}

func (op *Name) Match(q Query, cs []InputState) (ret []InputState) {
	for _, in := range cs {
		ws := in.Words()
		if cnt := len(ws); cnt > 0 {
			if a, det := q.FindArticle(ws); det >= 0 {
				start := det
				width := keywordScan(ws[start:])
				out := in.Next(start + width)
				if det > 0 {
					out.AddResult(a)
				}
				n := MatchedName{grok.Span(ws[start : start+width])}
				out.AddResult(n)
				ret = append(ret, out)
			}
		}
	}
	return
}

type MatchedName struct{ grok.Span }

type Called struct{}

func (*Called) Match(q Query, cs []InputState) (ret []InputState) {
	return matchAll(cs, grok.Keyword.Called)
}

// returns a width
func keywordScan(ws []grok.Word) (ret int) {
	ret = len(ws) // provisionally the whole thing.
	for i, w := range ws {
		switch w.Hash() {
		case grok.Keyword.And,
			grok.Keyword.Are,
			grok.Keyword.Comma,
			grok.Keyword.Is:
			ret = i
			break
		}
	}
	return
}
