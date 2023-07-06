package groktest

import "git.sr.ht/~ionous/tapestry/support/grok"

type SpanList [][]grok.Word

func PanicSpan(s string) grok.Span {
	out, e := grok.MakeSpan(s)
	if e != nil {
		panic(e)
	}
	return out
}

func PanicSpans(strs ...string) (out SpanList) {
	out = make(SpanList, len(strs))
	for i, str := range strs {
		out[i] = PanicSpan(str)
	}
	return
}

func (ws SpanList) FindMatch(words grok.Span) (ret grok.Match, none error) {
	if i, skip := ws.FindPrefix(words); skip > 0 {
		ret = grok.Span(ws[i])
	}
	return
}

// find the index and length of a prefix matching the passed words
func (ws SpanList) FindPrefix(words grok.Span) (retWhich int, retLen int) {
	if wordCount := len(words); wordCount > 0 {
		for prefixIndex, prefix := range ws {
			// every Word in el has to exist in words for it to be a prefix
			// and it has to be longer than any other previous match for it to be the best match
			// ( tbd? try a sort search? my first attempt failed miserably )
			if prefixLen := len(prefix); prefixLen <= wordCount && prefixLen > retLen {
				var failed bool
				for i, a := range prefix {
					if a.Hash() != words[i].Hash() {
						failed = true
						break
					}
				}
				if !failed {
					retWhich, retLen = prefixIndex, prefixLen
				}
			}
		}
	}
	return
}
