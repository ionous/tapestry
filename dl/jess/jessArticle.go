package jess

import "git.sr.ht/~ionous/tapestry/support/match"

func (op *Article) Match(q Query, input *InputState) (okay bool) {
	if m, width := match.FindCommonArticles(input.Words()); width > 0 {
		if words := input.Cut(width); input.Offset() == 0 || !startsUpper(words) {
			// build flags:
			if match.FindExactMatch(m, pluralNamed) >= 0 {
				op.Flags.Plural = true
			} else if useIndefinite(q) && match.FindExactMatch(m, indefinite) >= 0 {
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
