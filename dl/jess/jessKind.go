package jess

import "git.sr.ht/~ionous/tapestry/rt/kindsOf"

func (op *Kind) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.matchKind(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *Kind) matchKind(q Query, input *InputState) (okay bool) {
	var ks kindsOf.Kinds
	if m, width := q.FindKind(input.Words(), &ks); width > 0 && ks == kindsOf.Kind {
		m := matchedString{m, width}
		op.Matched, *input, okay = m, input.Skip(width), true
	}
	return
}

func (op *Kind) String() string {
	return op.Matched.String()
}

func (op *Kind) GetName(traits, kinds []Matched) (ret resultName) {
	return resultName{
		Traits: traits,
		// the order of kinds matters for "kinds of"
		// for: A container is a kind of thing.
		// the kinds should appear in that order in this list:
		Kinds: append([]Matched{op.Matched}, kinds...),
		// no name and no article because, the object itself is anonymous.
		// ( the article associated with the kind gets eaten )
	}
}
