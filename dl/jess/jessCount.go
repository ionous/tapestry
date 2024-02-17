package jess

import (
	"git.sr.ht/~ionous/tapestry/support/grok"
)

func (op *CountedName) Match(q Query, input *InputState) (okay bool) {
	if start := *input; //
	Optional(q, &start, &op.Article) || true {
		if next := start; //
		op.matchNumber(q, &next) &&
			op.Kind.Match(q, &next) {
			op.Matched = start.Cut(len(start) - len(next))
			*input, okay = next, true
		}
	}
	return
}

// try a word as a number...
func (op *CountedName) matchNumber(q Query, input *InputState) (okay bool) {
	if ws := input.Words(); len(ws) > 0 {
		word := ws[0].String()
		if v, ok := grok.WordsToNum(word); ok && v > 0 {
			const width = 1
			op.Number = float64(v)
			*input, okay = input.Skip(width), true
		}
	}
	return
}

func (op *CountedName) String() string {
	return op.Matched.String()
}

func (op *CountedName) GetName(traits, kinds []Matched) (ret grok.Name) {
	return grok.Name{
		Article: grok.Article{
			Count: int(op.Number),
		},
		Traits: traits,
		Kinds:  append(kinds, op.Kind.Matched),
		// no name, anonymous and counted.
	}
}
