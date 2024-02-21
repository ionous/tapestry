package jess

func (op *AdditionalTraits) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.CommaAnd) || true) &&
		op.Traits.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *Trait) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.matchTrait(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *Trait) matchTrait(q Query, input *InputState) (okay bool) {
	if m, width := q.FindTrait(input.Words()); width > 0 {
		m := matchedString{m, width}
		op.Matched, *input, okay = m, input.Skip(width), true
	}
	return
}

// its interesting that we dont have to store anything else
// all the trait info is in this... even additional traits.
func (op *Traits) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Trait.Match(q, &next) {
		Optional(q, &next, &op.AdditionalTraits)
		*input, okay = next, true
	}
	return
}

// unwind the tree of additional traits
func (op *Traits) GetTraits() Traitor {
	return Traitor{op}
}

// trait iterator
type Traitor struct {
	next *Traits
}

func (it Traitor) HasNext() bool {
	return it.next != nil
}

func (it *Traitor) GetNext() (ret Trait) {
	var next *Traits
	if more := it.next.AdditionalTraits; more != nil {
		next = &more.Traits
	}
	ret, it.next = it.next.Trait, next
	return
}
