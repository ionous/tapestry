package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

// called can have its own kind, its own specific article, and its name is flagged as "exact"
// ( where regular names are treated as potential aliases of existing names. )
func (op *KindCalled) GetName(traits, kinds []Matched) grok.Name {
	return grok.Name{
		Article: ReduceArticle(op.Kind.Article),
		Span:    op.Matched.(Span),
		Exact:   true,
		Traits:  traits,
		Kinds:   append(kinds, op.Kind.Span()),
	}
}

func (op *KindCalled) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Kind.Match(q, &next) &&
		op.Called.Match(q, &next, grok.Keyword.Called) &&
		Optionally(q, &next, &op.Article) &&
		op.matchName(q, &next) {
		*input, okay = next, true
	}
	return
}

// match the words after "called" ending with either "is/are" or the end of the string.
// this means that something like "a container called the bottle and a woman called the genie"
// generates a single object with a long, strange name.
func (op *KindCalled) matchName(q Query, input *InputState) (okay bool) {
	if width := beScan(input.Words()); width > 0 {
		op.Matched, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}
