package jess

func (op *AdditionalText) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.CommaAndOr.Match(q, &next) &&
		op.QuotedTexts.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

// its interesting that we dont have to store anything else
// all the trait info is in this... even additional traits.
func (op *QuotedTexts) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.QuotedText.Match(q, &next) {
		Optional(q, &next, &op.AdditionalText)
		*input, okay = next, true
	}
	return
}

// unwind the tree of additional traits
func (op *QuotedTexts) Iterate() TextListIterator {
	return TextListIterator{op}
}

func (op *QuotedTexts) Reduce() (ret []string) {
	for it := op.Iterate(); it.HasNext(); {
		ret = append(ret, it.GetNext())
	}
	return
}

// trait iterator
type TextListIterator struct {
	next *QuotedTexts
}

func (it TextListIterator) HasNext() bool {
	return it.next != nil
}

func (it *TextListIterator) GetNext() (ret string) {
	var next *QuotedTexts
	if more := it.next.AdditionalText; more != nil {
		next = &more.QuotedTexts
	}
	ret, it.next = it.next.QuotedText.String(), next
	return
}
