package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

func (op *Called) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.NamedKind.Match(q, &next) &&
		op.Called.Match(q, &next, grok.Keyword.Called) &&
		op.matchName(q, &next) {
		*input, okay = next, true
	}
	return
}

// match the words after "called" ending with either "is/are" or the end of the string.
func (op *Called) matchName(q Query, input *InputState) (okay bool) {
	if width := beScan(input.Words()); width > 0 {
		op.Name, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}
