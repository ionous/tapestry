package jess

import (
	"slices"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
)

func (op *Kind) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.matchKind(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *Kind) matchKind(q Query, input *InputState) (okay bool) {
	var k kindsOf.Kinds
	if m, width := q.FindKind(input.Words(), &k); width > 0 && filterKinds(q, k) {
		op.DeclaredKind = DeclaredKind{m, k}
		op.Matched, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}

// returns the real ( generally plural ) name of the kind
func (op *Kind) String() string {
	return op.DeclaredKind.actual
}

func (op *Kind) GetName(traits, kinds []string) (ret resultName) {
	return resultName{
		Traits: traits,
		// the order of kinds matters for "kinds of"
		// for: A container is a kind of thing.
		// the kinds should appear in that order in this list:
		Kinds: append([]string{op.String()}, kinds...),
		// no name and no article because, the object itself is anonymous.
		// ( the article associated with the kind gets eaten )
	}
}

func filterKinds(q Query, k kindsOf.Kinds) (okay bool) {
	flags := q.GetContext()
	if (flags & PropertyKinds) != 0 {
		okay = slices.Contains([]kindsOf.Kinds{kindsOf.Kind, kindsOf.Aspect, kindsOf.Record}, k)
	} else {
		okay = k == kindsOf.Kind
	}
	return
}

type DeclaredKind struct {
	actual string // as opposed to just what matched
	base   kindsOf.Kinds
}
