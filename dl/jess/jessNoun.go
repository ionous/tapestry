package jess

func (op *Noun) String() string {
	return op.ActualNoun
}

func (op *Noun) GetName(traits, kinds []string) (ret resultName, err error) {
	ret = resultName{
		Article: reduceArticle(op.Article),
		Matched: op.ActualNoun,
		Traits:  traits,
		Kinds:   kinds,
	}
	return
}

func (op *Noun) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.matchNoun(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *Noun) matchNoun(q Query, input *InputState) (okay bool) {
	if m, width := q.FindNoun(input.Words()); width > 0 {
		op.ActualNoun = m
		op.Matched, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}

// the noun that matched ( as opposed to the name that matched )
type ActualNoun = string
