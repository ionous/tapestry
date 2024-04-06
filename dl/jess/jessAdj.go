package jess

import "git.sr.ht/~ionous/tapestry/rt/kindsOf"

func (op *Adjectives) Match(q Query, input *InputState) (okay bool) {
	next := *input
	traits := Optional(q, &next, &op.Traits)
	if traits {
		Optional(q, &next, &op.CommaAnd)
	}
	kind := Optional(q, &next, &op.Kind)
	if kind {
		Optional(q, &next, &op.AdditionalAdjectives)
	}
	if traits || kind {
		*input, okay = next, true
	}
	return
}

func (op *Adjectives) GetTraits() (ret Traitor) {
	if ts := op.Traits; ts != nil {
		ret = ts.GetTraits()
	}
	return
}

func (op Adjectives) Reduce() (ret NounProperties, err error) {
	for t := op; ; {
		ret.Traits = append(ret.Traits, ReduceTraits(op.GetTraits())...)
		if k := t.Kind; k != nil {
			// for something to have adjectives (ie. traits) it must be a noun of some sort
			if kn, e := k.Validate(kindsOf.Kind); e != nil {
				err = e
				break
			} else {
				ret.Kinds = append(ret.Kinds, kn)
			}
		}
		if next := t.AdditionalAdjectives; next == nil {
			break
		} else {
			t = next.Adjectives
		}
	}
	return
}

func (op *AdditionalAdjectives) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.CommaAnd.Match(q, &next) &&
		op.Adjectives.Match(q, &next) {
		*input, okay = next, true
	}
	return
}
