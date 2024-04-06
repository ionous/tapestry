package jess

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/rt"
)

// --------------------------------------------------------------
// QuotedTexts
// --------------------------------------------------------------

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

func (op *QuotedTexts) Assignment() (ret rt.Assignment) {
	var els []rt.TextEval
	for it := op.Iterate(); it.HasNext(); {
		next := it.GetNextText()
		els = append(els, next.TextEval())
	}
	return &assign.FromTextList{
		Value: &list.MakeTextList{
			Values: els,
		},
	}
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

// --------------------------------------------------------------
// AdditionalText
// --------------------------------------------------------------

func (op *AdditionalText) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.CommaAndOr.Match(q, &next) &&
		op.QuotedTexts.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

// --------------------------------------------------------------
// TextListIterator
// --------------------------------------------------------------

type TextListIterator struct {
	next *QuotedTexts
}

func (it TextListIterator) HasNext() bool {
	return it.next != nil
}

func (it *TextListIterator) GetNext() string {
	return it.GetNextText().String()
}

func (it *TextListIterator) GetNextText() (ret *QuotedText) {
	var next *QuotedTexts
	if more := it.next.AdditionalText; more != nil {
		next = &more.QuotedTexts
	}
	ret, it.next = &it.next.QuotedText, next
	return
}
