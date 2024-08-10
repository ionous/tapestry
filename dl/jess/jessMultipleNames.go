package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

func (op *MultipleNames) Next() (ret *MultipleNames) {
	if next := op.AdditionalNames; next != nil {
		ret = &next.Names
	}
	return
}

// some callers want to fail matching on anonymous leading kinds
// tbd: would it be better to match, and error on generation?
// ( ie. to produce a message )
func (op *MultipleNames) HasAnonymousKind() bool {
	return op.Kind != nil
}

// unwind the traits ( if any ) of the names
func (op *MultipleNames) GetTraits() (ret *Traits) {
	if c := op.KindCalled; c != nil {
		ret = c.GetTraits()
	}
	return
}

// checks Query flags to control matching
func (op *MultipleNames) Match(q JessContext, input *InputState) (okay bool) {
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

func (op *MultipleNames) GetActualNoun() ActualNoun {
	return op.actualNoun
}

// implements NounMaker by calling BuildNouns on all matched names
func (op *MultipleNames) BuildNouns(q JessContext, w weaver.Weaves, run rt.Runtime, props NounProperties) (ret []DesiredNoun, err error) {
	for at := op; at != nil; at = at.Next() {
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
			out := append(ret, ns...)
			var an ActualNoun
			if len(out) == 1 {
				n := out[0]
				an = ActualNoun{
					n.Noun,
					n.CreatedKind,
				}
			}
			op.actualNoun = an
			ret = out
		}
	}
	return
}

func (op *AdditionalNames) Match(q JessContext, input *InputState) (okay bool) {
	if next := *input; //
	op.CommaAnd.Match(q, &next) &&
		op.Names.Match(q, &next) {
		*input, okay = next, true
	}
	return
}
