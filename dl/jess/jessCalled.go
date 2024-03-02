package jess

import "git.sr.ht/~ionous/tapestry/rt/kindsOf"

// called can have its own kind, its own specific article, and its name is flagged as "exact"
// ( where regular names are treated as potential aliases of existing names. )
func (op *KindCalled) BuildNoun(traits, kinds []string) (ret DesiredNoun, err error) {
	if kind, e := op.GetKind(); e != nil {
		err = e
	} else {
		for ts := op.GetTraits(); ts.HasNext(); {
			t := ts.GetNext()
			traits = append(traits, t.String())
		}
		// ignores the article of the kind,
		// in favor of the article closest to the named noun
		ret, err = op.NamedNoun.BuildNoun(traits, append(kinds, kind))
	}
	return
}

func (op *KindCalled) GetNormalizedName() string {
	return op.NamedNoun.GetNormalizedName()
}

func (op *KindCalled) GetKind() (string, error) {
	return op.Kind.Validate(kindsOf.Kind)
}

func (op *KindCalled) GetTraits() (ret Traitor) {
	if op.Traits != nil {
		ret = op.Traits.GetTraits()
	}
	return
}

func (op *KindCalled) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Traits) || true) &&
		op.Kind.Match(q, &next) &&
		op.Called.Match(q, &next) &&
		op.NamedNoun.Match(AddContext(q, CheckIndefiniteArticles), &next) {
		*input, okay = next, true
	}
	return
}
