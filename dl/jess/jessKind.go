package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

type TraitsKind struct {
	Traits
	*KindName
}

// its interesting that we dont have to store anything else
// all the trait info is in this... even additional traits.
func (op *TraitsKind) Match(q Query, in InputState) (ret int) {
	if ts := op.Traits.Match(q, in); ts > 0 {
		kind := Optionally(q, in.Skip(ts), &op.KindName)
		ret = ts + kind
	}
	return
}

func (op *TraitsKind) GetTraitSet() grok.TraitSet {
	var k grok.Matched
	if op.KindName != nil {
		k = op.KindName.Matched
	}
	return grok.TraitSet{
		Kind:   k,
		Traits: op.Traits.GetTraits(),
	}
}
