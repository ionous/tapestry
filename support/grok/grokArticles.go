package grok

func FindCommonArticles(ws Span) (ret Article, err error) {
	if ws.NumWords() > 0 {
		if cnt, ok := WordsToNum(ws[0].String()); ok {
			ret = Article{
				Match: Span(ws[:1]),
				Count: cnt,
			}
		} else if m, e := determiners.FindMatch(ws); e != nil {
			err = e
		} else {
			ret = Article{
				Match: m,
			}
		}
	}
	return
}

// for now, these are fixed.
// when the author specifies some particular indefinite article for a noun
// that article only gets used for printing the noun;
// it doesn't enhance the parsing of the story.
// it would take some work to lightly hold the relation between a name and an article
// then parse a sentence matching names to nouns in the grok
var determiners = PanicSpans("the", "a", "an", "some", "our")
