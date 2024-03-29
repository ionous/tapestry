package jess

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// runs in the PropertyPhase phase
func (op *KindsAreTraits) Phase() weaver.Phase {
	return weaver.PropertyPhase
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

func (op *KindsAreTraits) Weave(w weaver.Weaves, run rt.Runtime) (err error) {
	traits := op.Traits.GetTraits()
	for kt := op.Kinds.Iterate(); kt.HasNext(); {
		k := kt.GetNext()
		name := k.String()
		if lhs := k.GetTraits(); lhs.HasNext() {
			err = fmt.Errorf("unexpected traits before %s", name)
			break
		}
		if e := AddKindTraits(w, name, traits); e != nil {
			err = e
			break
		}
	}
	return
}
