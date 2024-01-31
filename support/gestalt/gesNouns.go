package gestalt

import "git.sr.ht/~ionous/tapestry/support/grok"

// matches one or more *potential* nouns, their determiners, and separators.

// matches names that might be nouns.
// because this could generate a large number of possibilities,
// there are some assumptions for optimization:
// the words "is/are/comma/and" are never part of noun names.
// fix: quotes for "titles" ( which are then allowed to break the rules )
type NounExactly struct{}

func (op *NounExactly) Match(q Query, cs []InputState) (ret []InputState) {
	for _, in := range cs {
		ws := in.Words()
		if cnt := len(ws); cnt > 0 {
			if a, det := q.FindArticle(ws); det >= 0 {
				next := keywordScan(ws)
				out := in.Next(next)
				out.AddResult(TypeArticle, a)
				out.AddResult(TypeExactName, grok.Span(ws[det:next]))
				ret = append(ret, out)
			}
		}
	}
	return
}

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
