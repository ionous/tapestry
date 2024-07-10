package jess

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// runs in the PropertyPhase
func (op *KindsAreTraits) Phase() weaver.Phase {
	return weaver.PropertyPhase
}

func (op *KindsAreTraits) MatchLine(q Query, line InputState) (ret InputState, okay bool) {
	if next := line; //
	op.Kinds.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.Usually.Match(q, &next, keywords.Usually) &&
		op.Traits.Match(q, &next) {
		ret, okay = next, true
	}
	return
}

func (op *KindsAreTraits) Weave(w weaver.Weaves, run rt.Runtime) (err error) {
	traits := op.Traits.GetTraits()
	for k := &op.Kinds; k != nil; k = k.Next() {
		if name, e := k.GetNormalizedName(); e != nil {
			err = e
			break
		} else {
			if lhs := k.GetTraits(); lhs != nil {
				err = fmt.Errorf("unexpected traits before %s", name)
				break
			}
			if e := AddKindTraits(w, name, traits); e != nil {
				err = e
				break
			}
		}
	}
	return
}
