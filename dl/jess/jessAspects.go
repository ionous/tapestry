package jess

import "git.sr.ht/~ionous/tapestry/support/inflect"

func (op *AspectsAreTraits) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	// there's nothing in the current query interface to limit /filter the kind
	// so we let the generate do that
	op.Aspect.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.Names.Match(AddContext(q, PlainNameMatching), &next) {
		*input, okay = next, true
	}
	return
}

func (op *AspectsAreTraits) Generate(rar Registrar) (err error) {
	var names []string
	for it := op.Names.Iterate(); it.HasNext(); {
		n := it.GetNext()
		names = append(names, inflect.Normalize(n.String()))
	}
	return rar.AddTraits(op.Aspect.String(), names)
}
