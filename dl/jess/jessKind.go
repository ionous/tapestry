package jess

func (op *Kind) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.matchKind(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *Kind) matchKind(q Query, input *InputState) (okay bool) {
	if m, width := q.FindKind(input.Words()); width > 0 {
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
