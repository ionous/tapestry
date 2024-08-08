package jess

import (
	"git.sr.ht/~ionous/tapestry/support/match"
)

// iterator helper
func (op *MultipleKinds) Next() (ret *MultipleKinds) {
	if next := op.AdditionalKinds; next != nil {
		ret = &next.Kinds
	}
	return
}

func (op *AdditionalKinds) Match(q JessContext, input *InputState) (okay bool) {
	if next := *input; //
	op.CommaAnd.Match(q, &next) &&
		op.Kinds.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

// note: traits are matched, but prohibited by "kinds_are_traits"
func (op *MultipleKinds) Match(q JessContext, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Traits) || true) &&
		(Optional(q, &next, &op.Article) || true) &&
		op.matchName(&next) {
		Optional(q, &next, &op.AdditionalKinds)
		*input, okay = next, true
	}
	return
}

func (op *MultipleKinds) matchName(input *InputState) (okay bool) {
	if width := nameScan(input.Words()); width > 0 {
		op.Matched = input.Cut(width)
		*input, okay = input.Skip(width), true
	}
	return
}

func (op *MultipleKinds) GetNormalizedName() (string, error) {
	return match.NormalizeAll(op.Matched)
}

// unwind the tree of traits
func (op *MultipleKinds) GetTraits() (ret *Traits) {
	if ts := op.Traits; ts != nil {
		ret = ts.GetTraits()
	}
	return
}
