package jess

func (op *AdditionalKinds) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.CommaAnd.Match(q, &next) &&
		op.Kinds.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *Kinds) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Kind.Match(q, &next) {
		Optional(q, &next, &op.AdditionalKinds)
		*input, okay = next, true
	}
	return
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
