package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

func (op *Counted) Match(q Query, input *InputState) (okay bool) {
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
