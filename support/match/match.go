package match

import "strings"

// Span - implements Match for a chain of individual words.
type Span []Word

func Equals(ts []TokenValue, ws Span) (okay bool) {
	if okay = len(ts) == len(ws); okay { // provisionally set okay.
		for i, el := range ts {
			if el.Hash() != ws[i].hash {
				okay = false
				break
			}
		}
	}
	return
}

func JoinWords(ws []Word) string {
	var b strings.Builder
	for i, w := range ws {
		if i > 0 {
			b.WriteRune(' ')
		}
		b.WriteString(w.String())
	}
	return b.String()
}

func (s Span) String() string {
	return JoinWords(s)
}

func HasPrefix(ts []TokenValue, prefix []Word) (okay bool) {
	// a prefix must be the same as or shorter than us
	if len(prefix) <= len(ts) {
		okay = true // provisionally
		for i, a := range prefix {
			if a.Hash() != ts[i].Hash() {
				okay = false
				break
			}
		}
	}
	return
}

// search for a span in a list of spans;
// return the index of the span that matched.
func FindExactMatch(ts []TokenValue, spans []Span) (ret int) {
	ret = -1 // provisionally
	for i, span := range spans {
		if Equals(ts, span) {
			ret = i
			break
		}
	}
	return
}
