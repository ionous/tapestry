package jess

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

// runs in the PropertyPhase phase
func (op *KindsAreTraits) Phase() Phase {
	return mdl.PropertyPhase
}

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

func (op *KindsAreTraits) Generate(rar *Context) (err error) {
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
