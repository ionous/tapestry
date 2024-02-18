package jess

import "git.sr.ht/~ionous/tapestry/support/match"

func (op *Article) Match(q Query, input *InputState) (okay bool) {
	if m, width := FindCommonArticles(input.Words()); width > 0 {
		op.Matched, *input, okay = m, input.Skip(width), true
	}
	return
}

func reduceArticle(op *Article) (ret articleResult) {
	if op != nil {
		ret = articleResult{
			Matched: op.Matched,
		}
	}
	return
}

func StripArticle(name string) (ret string, err error) {
	if parts, e := match.MakeSpan(name); e != nil {
		err = e
	} else if len(parts) <= 1 {
		ret = name
	} else if _, width := FindCommonArticles(parts); width > 0 {
		words := parts[width:]
		ret = words.String()
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
	if i, skip := determiners.FindPrefix(ws); skip > 0 {
		ret, width = match.Span(determiners[i]), skip
	}
	return
}

var determiners = match.PanicSpans("the", "a", "an", "some", "our")
