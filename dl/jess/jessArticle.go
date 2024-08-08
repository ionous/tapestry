package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// match articles ( such as: a, the, some )
// when the query is flagged with CheckIndefiniteArticles
// records *definite* article for use in indefinite contexts.
// ( re: append/applyArticle )
//
// note: capitalized articles in the middle of a sentence are treated as part of the name.
// ex. "a man called The Vampire" assumes The Vampire is proper named."
func (op *Article) Match(q JessContext, input *InputState) (okay bool) {
	ws := input.Words()
	// match: the, a, some, etc. ( case insensitive )
	if m, width := match.FindCommonArticles(ws); width > 0 {
		// ignore articles that seem to be part of a proper noun ( case filtering )
		if words := input.Cut(width); words[0].First || !startsUpper(words) {
			// build flags:
			article := ws[:width]
			if match.FindExactMatch(article, pluralNamed) >= 0 {
				op.flags.Plural = true
			} else if useIndefinite(q) && match.FindExactMatch(article, definiteArticles) >= 0 {
				op.flags.Indefinite = true
			}
			// return okay:
			op.Text, *input = m.String(), input.Skip(width)
			okay = true
		}
	}
	return
}

// extra info about specified articles determined during matching.
// stored as a private part of the Article generated type.
type ArticleFlags struct {
	Indefinite bool // use the specified article in indefinite contexts
	Plural     bool
}

var definiteArticles = match.PanicSpans("the", "our")
var pluralNamed = match.PanicSpans("some")

func TryArticle(q JessContext, in InputState,
	accept func(*Article, InputState),
	reject func(error),
) {
	q.Try(After(weaver.LanguagePhase), func(weaver.Weaves, rt.Runtime) {
		var art Article
		if !art.Match(q, &in) {
			accept(nil, in) // zero or more so no traits is okay.
		} else {
			accept(&art, in)
		}
	}, reject)
}
