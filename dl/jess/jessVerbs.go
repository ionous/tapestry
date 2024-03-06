package jess

func (op *VerbPhrase) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Verb.Match(q, &next) &&
		op.PlainNames.Match(AddContext(q, PlainNameMatching), &next) {
		*input, okay = next, true
	}
	return
}

// the passed phrase is the macro to match
func (op *Verb) Match(q Query, input *InputState) (okay bool) {
	if m, width := q.FindMacro(input.Words()); width > 0 {
		op.Macro = m
		op.Matched = input.Cut(width)
		*input, okay = input.Skip(width), true
	}
	return
}
