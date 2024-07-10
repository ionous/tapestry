package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// some callers want to fail matching on anonymous leading kinds
// tbd: would it be better to match, and error on generation?
// ( ie. to produce a message )
func (op *Names) HasAnonymousKind() bool {
	return op.Kind != nil
}

// unwind the tree of additional names
func (op *Names) GetNames() Iterator {
	return Iterator{op}
}

// unwind the traits ( if any ) of the names
func (op *Names) GetTraits() (ret Traitor) {
	if c := op.KindCalled; c != nil {
		ret = c.GetTraits()
	}
	return

}

// checks Query flags to control matching
func (op *Names) Match(q Query, input *InputState) (okay bool) {
	if next := *input;                   //
	Optional(q, &next, &op.Pronoun) || ( //
	//
	matchKinds(q) &&
		// "5 containers",
		Optional(q, &next, &op.CountedKind) ||
		// "the container called the bottle"
		Optional(q, &next, &op.KindCalled) ||
		// "the container"
		Optional(q, &next, &op.Kind)) || ( //
	// "the bottle", or "the unknown name"
	Optional(q, &next, &op.Name)) {
		// can only match in the first name
		q = ClearContext(q, MatchPronouns)
		// as long as one succeeded, try matching additional names too...
		// inform seems to only allow "kind called" at the front of a list of names...
		// fix? but maybe it'd be better to match and reject when used incorrectly.
		// [ an advantage is -- we could register a whole side ( lhs/rhs ) to nouns at once ]
		if !Optional(q, &next, &op.AdditionalNames) || op.AdditionalNames.Names.KindCalled == nil {
			*input, okay = next, true
		}
	}
	return
}

// implements NounBuilder by calling BuildNouns on all matched names
func (op *Names) BuildNouns(q Query, w weaver.Weaves, run rt.Runtime, props NounProperties) (ret []DesiredNoun, err error) {
	for n := op.GetNames(); n.HasNext(); {
		at := n.GetNext()
		if ns, e := buildNounsFrom(q, w, run, props,
			nillable(at.Pronoun),
			nillable(at.CountedKind),
			nillable(at.KindCalled),
			nillable(at.Kind),
			nillable(at.Name),
		); e != nil {
			err = e
			break
		} else {
			ret = append(ret, ns...)
		}
	}
	return
}
