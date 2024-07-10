package jess

import (
	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/rt"
)

// --------------------------------------------------------------
// QuotedTexts
// --------------------------------------------------------------
func (op *QuotedTexts) Next() (ret *QuotedTexts) {
	if next := op.AdditionalText; next != nil {
		ret = &next.QuotedTexts
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

func (op *QuotedTexts) Assignment() (ret rt.Assignment) {
	var els []rt.TextEval
	for it := op; it != nil; it = it.Next() {
		els = append(els, it.QuotedText.TextEval())
	}
	return &call.FromTextList{
		Value: &list.MakeTextList{
			Values: els,
		},
	}
}

func (op *QuotedTexts) Reduce() (ret []string) {
	for it := op; it != nil; it = it.Next() {
		str := it.QuotedText.String()
		ret = append(ret, str)
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
