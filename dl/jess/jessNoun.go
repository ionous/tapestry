package jess

func (op *Noun) BuildNoun(traits, kinds []string) (ret DesiredNoun, err error) {
	a, flags := getOptionalArticle(op.Article)
	ret = DesiredNoun{
		Article: a,
		Flags:   flags,
		Noun:    op.ActualNoun,
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
	if ws := input.Words(); len(ws) > 0 {
		if m, width := q.FindNoun(ws, ""); width > 0 {
			op.ActualNoun = m
			op.Matched, *input, okay = input.Cut(width), input.Skip(width), true
		} else if ws[0].Hash() == keywords.You {
			// fix? it'd be nice if the mapping of "you" to "self" was handled by script;
			// or even not necessary at all.
			width := 1
			op.ActualNoun = PlayerSelf
			op.Matched, *input, okay = input.Cut(width), input.Skip(width), true
		}
	}
	return
}

// the noun that matched ( as opposed to the name that matched )
type ActualNoun = string
