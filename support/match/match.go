package match

import "strings"

// Span - implements Match for a chain of individual words.
type Span []Word

func (s Span) Equals(ws Span) (okay bool) {
	if okay = len(s) == len(ws); okay { // provisionally set okay.
		for i, el := range s {
			if el.hash != ws[i].hash {
				okay = false
				break
			}
		}
	}
	return
}

func (s Span) String() string {
	var b strings.Builder
	for i, w := range s {
		if i > 0 {
			b.WriteRune(' ')
		}
		b.WriteString(w.String())
	}
	return b.String()
}

func (s Span) NumWords() int {
	return len(s)
}

func HasPrefix(s, prefix []Word) (okay bool) {
	// a prefix must be the same as or shorter than us
	if len(prefix) <= len(s) {
		okay = true // provisionally
		for i, a := range prefix {
			if a.Hash() != s[i].Hash() {
				okay = false
				break
			}
		}
	}
	return
}

func FindMatch(s Span, spans []Span) (ret int) {
	ret = -1 // provisionally
	for i, el := range spans {
		if s.Equals(el) {
			ret = i
			break
		}
	}
	return
}
