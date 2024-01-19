package compact

import (
	"strings"
)

type Signature struct {
	Key    string // full key
	Name   string // lede
	Labels []string
}

// helper for debug printing
func (s *Signature) DebugString() string {
	var b strings.Builder
	b.WriteString(s.Name)
	for i, p := range s.Labels {
		if i == 0 {
			b.WriteRune(':')
		} else {
			b.WriteRune(',')
		}
		b.WriteString(p)
	}
	b.WriteRune(':')
	return b.String()
}
