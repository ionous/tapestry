package match

type SpanList []Span

// this is the same as FindPrefix only it returns a Span instead of an index
func (ws SpanList) FindMatch(words Span) (ret Span, retLen int) {
	// idx is the index of the span in the list
	if idx, skip := ws.FindPrefix(words); skip > 0 {
		ret, retLen = ws[idx], skip
	}
	return
}

// see anything in our span list starts the passed words.
// for instance, if the span list contains the span "oh hello"
// then the words "oh hello world" will match
// returns the index of the  index and length of the longest prefix
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
