package jess

func (op *AdditionalNames) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.CommaAnd.Match(q, &next) &&
		op.Names.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

type Iterator struct {
	next *Names
}

func (it Iterator) HasNext() bool {
	return it.next != nil
}

func (it *Iterator) GetNext() (ret *Names) {
	var next *Names
	if more := it.next.AdditionalNames; more != nil {
		next = &more.Names
	}
	ret, it.next = it.next, next
	return
}
