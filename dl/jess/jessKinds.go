package jess

import "git.sr.ht/~ionous/tapestry/support/inflect"

func (op *AdditionalKinds) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.CommaAnd.Match(q, &next) &&
		op.Kinds.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

// note: traits are matched, but prohibited by "kinds_are_traits"
func (op *Kinds) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Traits) || true) &&
		(Optional(q, &next, &op.Article) || true) &&
		op.matchName(q, &next) {
		Optional(q, &next, &op.AdditionalKinds)
		*input, okay = next, true
	}
	return
}

func (op *Kinds) matchName(q Query, input *InputState) (okay bool) {
	if width := keywordScan(input.Words()); width > 0 {
		op.Matched, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}

func (op *Kinds) String() string {
	return inflect.Normalize(op.Matched)
}

// unwind the tree of traits
func (op *Kinds) GetTraits() (ret Traitor) {
	if ts := op.Traits; ts != nil {
		ret = ts.GetTraits()
	}
	return
}

// unwind the tree of additional kinds
func (op *Kinds) Iterate() Kinder {
	return Kinder{op}
}

type Kinder struct {
	next *Kinds
}

func (it Kinder) HasNext() bool {
	return it.next != nil
}

func (it *Kinder) GetNext() (ret *Kinds) {
	var next *Kinds
	if more := it.next.AdditionalKinds; more != nil {
		next = &more.Kinds
	}
	ret, it.next = it.next, next
	return
}
