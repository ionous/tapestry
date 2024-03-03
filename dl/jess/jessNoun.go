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
	if cnt := keywordScan(input.Words()); cnt > 0 {
		sub := input.CutSpan(cnt)
		// fix? it'd be nice if the mapping of "you" to "self" was handled by script;
		// or even not necessary at all.
		if width := 1; len(sub) == width && sub[0].Hash() == keywords.You {
			op.ActualNoun = PlayerSelf
			op.Matched, *input, okay = input.Cut(width), input.Skip(width), true
		} else {
			// match the subsection normally:
			if m, width := q.FindNoun(sub, ""); width > 0 {
				op.ActualNoun = m
				op.Matched, *input, okay = input.Cut(width), input.Skip(width), true
			}

		}
	}
	return
}

// the noun that matched ( as opposed to the name that matched )
type ActualNoun = string
