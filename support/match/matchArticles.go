package match

// return the name after removing leading articles
// eats any errors it encounters and returns the original name
func StripArticle(name string) (ret string) {
	if parts, e := MakeSpan(name); e != nil {
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
func FindCommonArticles(ws Span) (ret Span, width int) {
	if m, skip := determiners.FindPrefix(ws); skip > 0 {
		ret, width = m, skip
	}
	return
}

// fix? i feel like this should be part of package inflect instead
var determiners = PanicSpans("the", "a", "an", "some", "our")
