package match

import "strings"

type SpanList []Span

func (ws SpanList) FindExactMatch(ts []TokenValue) (ret Span, width int) {
	if idx := FindExactMatch(ts, ws); idx >= 0 {
		ret, width = ws[idx], len(ts)
	}
	return
}

// this is the same as FindPrefixIndex only it returns a Span instead of an index
func (ws SpanList) FindPrefix(words []TokenValue) (ret Span, width int) {
	// idx is the index of the span in the list
	if idx, skip := ws.FindPrefixIndex(words); skip > 0 {
		ret, width = ws[idx], skip
	}
	return
}

// see anything in our span list starts the passed words.
// for instance, if the span list contains the span "oh hello"
// then the words "oh hello world" will match
// returns the index of the  index and length of the longest prefix
func (ws SpanList) FindPrefixIndex(words []TokenValue) (retWhich int, retWidth int) {
	if wordCount := len(words); wordCount > 0 {
		for prefixIndex, prefix := range ws {
			// every Word in el has to exist in words for it to be a prefix
			// and it has to be longer than any other previous match for it to be the best match
			// ( tbd? try a sort search? my first attempt failed miserably )
			if prefixLen := len(prefix); prefixLen <= wordCount && prefixLen > retWidth {
				if HasPrefix(words, prefix) {
					retWhich, retWidth = prefixIndex, prefixLen
				}
			}
		}
	}
	return
}

// transform a space separated string into a slice of hashes
func PanicSpan(s string) (ret Span) {
	for _, f := range strings.Fields(s) {
		ret = append(ret, MakeWord(f))
	}
	return
}

func PanicSpans(strs ...string) (out SpanList) {
	out = make(SpanList, len(strs))
	for i, str := range strs {
		out[i] = PanicSpan(str)
	}
	return
}
