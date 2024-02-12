package jess

import (
	"git.sr.ht/~ionous/tapestry/support/grok"
)

func (op *Are) Match(q Query, input *InputState) (okay bool) {
	if width := input.MatchWord(grok.Keyword.Are, grok.Keyword.Is); width > 0 {
		op.Matched = input.Cut(width)
		*input = input.Skip(width)
		okay = true
	}
	return
}

func (op *CommaAnd) Match(q Query, input *InputState) (okay bool) {
	if sep, e := grok.CommaAnd(input.Words()); e != nil {
		q.error("comma and", e)
	} else {
		width := sep.Len()
		op.Matched = input.Cut(width)
		*input = input.Skip(width)
		okay = true
	}
	return
}

func (op *Words) Match(q Query, input *InputState, hashes ...uint64) (okay bool) {
	if width := input.MatchWord(hashes...); width > 0 {
		op.Matched = input.Cut(width)
		*input = input.Skip(width)
		okay = true
	}
	return
}

func (op *Words) Span() grok.Span {
	return op.Matched.(grok.Span)
}

func (op *MacroName) Match(q Query, input *InputState, phrase grok.Span) (okay bool) {
	if m, e := q.g.FindMacro(phrase); e != nil {
		q.error("find macro", e)
	} else if grok.HasPrefix(input.Words(), phrase) {
		width := len(phrase)
		op.Macro = m
		op.Matched = input.Cut(width)
		*input = input.Skip(width)
		okay = true
	}
	return
}
