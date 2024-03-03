package jess

import "git.sr.ht/~ionous/tapestry/support/match"

func (op *Article) Match(q Query, input *InputState) (okay bool) {
	if m, width := FindCommonArticles(input.Words()); width > 0 {
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

func getOptionalArticle(art *Article) (retText string, retFlags ArticleFlags) {
	if art != nil {
		retText = art.Text
		retFlags = art.Flags
	}
	return
}

// return the name after removing leading articles
// eats any errors it encounters and returns the original name
func StripArticle(name string) (ret string) {
	if parts, e := match.MakeSpan(name); e != nil {
		ret = name
	} else if len(parts) <= 1 {
		ret = name
	} else if _, width := FindCommonArticles(parts); width > 0 {
		words := parts[width:]
		ret = words.String()
	} else {
		ret = name
	}
	return
}

// for now, the common articles are a fixed set.
// when the author specifies some particular indefinite article for a noun
// that article only gets used for printing the noun;
// it doesn't enhance the parsing of the story.
// ( it would take some work to lightly hold the relation between a name and an article
// then parse a sentence matching names to nouns in the
// fwiw: the articles in inform also seems to be predetermined in this way.  )
func FindCommonArticles(ws match.Span) (ret match.Span, width int) {
	if m, skip := determiners.FindPrefix(ws); skip > 0 {
		ret, width = m, skip
	}
	return
}

var determiners = match.PanicSpans("the", "a", "an", "some", "our")
var indefinite = match.PanicSpans("the", "our")
var pluralNamed = match.PanicSpans("some")
