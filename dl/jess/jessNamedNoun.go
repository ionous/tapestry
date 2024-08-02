package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

func (op *NamedNoun) GetNormalizedName() (ret string, err error) {
	if n := op.Noun; n != nil {
		ret = n.actualNoun.Name // the actual name is already normalized
	} else if n := op.Name; n != nil {
		ret, err = match.NormalizeAll(n.Matched)
	} else {
		panic("NamedNoun was unmatched")
	}
	return
}

// requires that BuildNouns have been called already
func (op *NamedNoun) GetDesiredNouns() []DesiredNoun {
	return op.desiredNouns
}

// panics if not matched
func (op *NamedNoun) BuildNouns(q Query, w weaver.Weaves, run rt.Runtime, props NounProperties) (ret []DesiredNoun, err error) {
	if nouns, e := buildNounsFrom(q, w, run, props,
		nillable(op.Pronoun),
		nillable(op.KindCalled),
		nillable(op.Noun),
		nillable(op.Name),
	); e != nil {
		err = e
	} else {
		op.desiredNouns = nouns
		ret = nouns
	}
	return
}

func (op *NamedNoun) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	Optional(q, &next, &op.Pronoun) ||
		Optional(q, &next, &op.KindCalled) ||
		Optional(q, &next, &op.Noun) ||
		Optional(q, &next, &op.Name) {
		*input, okay = next, true
	}
	return
}
