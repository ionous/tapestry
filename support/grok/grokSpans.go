package grok

type SpanList [][]Word

func (ws SpanList) FindMatch(words Span) (ret Match, none error) {
	if i, skip := ws.FindPrefix(words); skip > 0 {
		ret = Span(ws[i])
	}
	return
}

// find the index and length of the longest prefix matching the passed words
func (ws SpanList) FindPrefix(words Span) (retWhich int, retLen int) {
	if wordCount := len(words); wordCount > 0 {
		for prefixIndex, prefix := range ws {
			// every Word in el has to exist in words for it to be a prefix
			// and it has to be longer than any other previous match for it to be the best match
			// ( tbd? try a sort search? my first attempt failed miserably )
			if prefixLen := len(prefix); prefixLen <= wordCount && prefixLen > retLen {
				if HasPrefix(words, prefix) {
					retWhich, retLen = prefixIndex, prefixLen
				}
			}
		}
	}
	return
}

func PanicSpan(s string) Span {
	out, e := MakeSpan(s)
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
