package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// called can have its own kind, its own specific article, and its name is flagged as "exact"
// ( where regular names are treated as potential aliases of existing names. )
func (op *KindCalled) BuildNouns(q JessContext, w weaver.Weaves, run rt.Runtime, props NounProperties) (ret []DesiredNoun, err error) {
	if kind, e := op.GetKind(); e != nil {
		err = e
	} else {
		// ignores the article of the kind,
		// in favor of the article closest to the named noun
		ret, err = op.Name.BuildNouns(q, w, run, NounProperties{
			Traits: append(props.Traits, ReduceTraits(op.GetTraits())...),
			Kinds:  append(props.Kinds, kind),
		})
	}
	return
}

func (op *KindCalled) GetNormalizedName() (string, error) {
	return op.Name.GetNormalizedName()
}

func (op *KindCalled) GetKind() (string, error) {
	return op.Kind.Validate(kindsOf.Kind)
}

func (op *KindCalled) GetTraits() (ret *Traits) {
	if op.Traits != nil {
		ret = op.Traits.GetTraits()
	}
	return
}

// transparent container called (the) box.
func (op *KindCalled) Match(q JessContext, input *InputState) (okay bool) {
	if next := *input; // note: traits can eat up a leading article
	(Optional(q, &next, &op.Traits) || true) &&
		op.Kind.Match(q, &next) &&
		op.Called.Match(q, &next) &&
		op.Name.Match(AddContext(q, CheckIndefiniteArticles), &next) {
		*input, okay = next, true
	}
	return
}
