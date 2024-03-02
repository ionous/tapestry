package jess

import (
	"fmt"
)

func (op *KindsAreTraits) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Kinds.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.Usually.Match(q, &next, keywords.Usually) &&
		op.Traits.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *KindsAreTraits) Generate(rar Registrar) (err error) {
	traits := op.Traits.GetTraits()
	for kt := op.Kinds.Iterate(); kt.HasNext(); {
		k := kt.GetNext()
		name := k.String()
		if lhs := k.GetTraits(); lhs.HasNext() {
			err = fmt.Errorf("unexpected traits before %s", name)
			break
		}
		if e := AddKindTraits(rar, name, traits); e != nil {
			err = e
			break
		}
	}
	return
}
