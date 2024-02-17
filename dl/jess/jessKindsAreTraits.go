package jess

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

func (op *KindsAreTraits) GetResults() (ret grok.Results, err error) {
	panic("xxx")
}

func (op *KindsAreTraits) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Kinds.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.usually(q, &next) &&
		op.Traits.Match(q, &next) {
		q.note("matched KindsAreTraits")
		*input, okay = next, true
	}
	return
}

// match the words after "called" ending with either "is/are" or the end of the string.
func (op *KindsAreTraits) usually(q Query, input *InputState) (okay bool) {
	if m, _ := usually.FindMatch(input.Words()); m != nil {
		width := m.NumWords()
		op.Usually.Matched, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}

var usually = grok.PanicSpans("usually")

func (op *KindsAreTraits) Apply(rar Registrar) (err error) {
	traits := op.Traits.GetTraits()
	for kt := op.Kinds.Iterate(); kt.HasNext(); {
		k := kt.GetNext()
		name := inflect.Normalize(k.Kind.String())
		if lhs := k.GetTraits(); lhs.HasNext() {
			err = fmt.Errorf("unexpected traits before %s", name)
			break
		}
		if e := AddTraitsToKind(rar, name, traits); e != nil {
			err = e
			break
		}
	}
	return
}
