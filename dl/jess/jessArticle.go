package jess

import "git.sr.ht/~ionous/tapestry/support/match"

func (op *Article) Match(q Query, input *InputState) (okay bool) {
	ws := input.Words()
	if m, width := match.FindCommonArticles(ws); width > 0 {
		// ignore articles that seem to be part of a proper noun
		// by only allowing capitalized articles at the start of a sentence
		// ex. "a man called The Vampire" assumes The Vampire is proper named.
		if words := input.Cut(width); words[0].First || !startsUpper(words) {
			// build flags:
			article := ws[:width]
			if match.FindExactMatch(article, pluralNamed) >= 0 {
				op.Flags.Plural = true
			} else if useIndefinite(q) && match.FindExactMatch(article, indefinite) >= 0 {
				op.Flags.Indefinite = true
			}
			// return okay:
			op.Text, *input = m.String(), input.Skip(width)
			okay = true
		}
	}
	return
}

type ArticleFlags struct {
	Indefinite bool
	Plural     bool
}

var indefinite = match.PanicSpans("the", "our")
var pluralNamed = match.PanicSpans("some")
