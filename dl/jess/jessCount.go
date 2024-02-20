package jess

func (op *CountedName) Match(q Query, input *InputState) (okay bool) {
	if start := *input; //
	Optional(q, &start, &op.Article) || true {
		if next := start; //
		op.MatchingNumber.Match(q, &next) &&
			op.Kind.Match(q, &next) {
			op.Matched = start.Cut(len(start) - len(next))
			*input, okay = next, true
		}
	}
	return
}

func (op *CountedName) String() string {
	return op.Matched.String()
}

func (op *CountedName) GetName(traits, kinds []Matched) (ret resultName) {
	return resultName{
		Article: articleResult{
			Count: int(op.MatchingNumber.Number),
		},
		Traits: traits,
		Kinds:  append(kinds, op.Kind.Matched),
		// no name, anonymous and counted.
	}
}
