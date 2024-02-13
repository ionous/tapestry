package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

// called can have its own kind, its own specific article, and its name is flagged as "exact"
// ( where regular names are treated as potential aliases of existing names. )
func (op *Called) GetName(traits, kinds []Matched) grok.Name {
	return grok.Name{
		Article: ReduceArticle(op.TheKind.Article),
		Span:    op.Matched.(Span),
		Exact:   true,
		Traits:  traits,
		Kinds:   append(kinds, op.TheKind.Span()),
	}
}

func (op *Called) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.TheKind.Match(q, &next) &&
		op.Called.Match(q, &next, grok.Keyword.Called) &&
		Optionally(q, &next, &op.Article) &&
		op.matchName(q, &next) {
		*input, okay = next, true
	}
	return
}

// match the words after "called" ending with either "is/are" or the end of the string.
func (op *Called) matchName(q Query, input *InputState) (okay bool) {
	if width := beScan(input.Words()); width > 0 {
		op.Matched, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}

func TryNameCalled(q Query, input *InputState, out *NameCalled) (okay bool) {
	var b Called // called is more specific than name, so has to be checked first.
	var a TheName
	if ok := b.Match(q, input); ok {
		(*out), okay = &b, true
	} else if ok := a.Match(q, input); ok {
		(*out), okay = &a, true
	}
	return
}
