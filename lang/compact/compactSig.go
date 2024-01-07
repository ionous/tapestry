package compact

import (
	"strings"
)

type Signature struct {
	Key    string // full key
	Name   string
	Params []Param
}

type Param struct {
	Label  string
	Choice string // optional
}

// helper for debug printing
func (s *Signature) DebugString() string {
	var b strings.Builder
	b.WriteString(s.Name)
	for i, p := range s.Params {
		if i == 0 {
			b.WriteRune(':')
		} else {
			b.WriteRune(',')
		}
		b.WriteString(p.DebugString())
	}
	b.WriteRune(':')
	return b.String()
}

func (p *Param) Matches(s string) bool {
	return s == p.Label || s == "_" && len(p.Label) == 0
}

func (p *Param) DebugString() string {
	var b strings.Builder
	if l := p.Label; len(l) == 0 {
		b.WriteRune('_')
	} else {
		b.WriteString(l)
	}
	if len(p.Choice) > 0 {
		b.WriteRune(' ')
		b.WriteString(p.Choice)
	}
	return b.String()
}

func (p *Param) String() string {
	out := p.Label
	if len(p.Choice) > 0 {
		out = out + " " + p.Choice
	}
	return out
}
