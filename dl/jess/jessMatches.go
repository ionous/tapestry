package jess

// allows partial matches; test that there's no input left to verify a complete match.
func (op *MatchingPhrases) Match(q Query, input *InputState) (ret Generator, okay bool) {
	// fix? could change to reflect ( or expand type info ) to walk generically
	var best InputState
	for _, m := range []interface {
		Generator
		Interpreter
	}{
		&op.KindsAreTraits,
		&op.KindsOf,
		&op.KindsHaveProperties,
		&op.KindsAreEither,
		&op.VerbNamesAreNames,
		&op.NamesVerbNames,
		&op.NamesAreLikeVerbs,
		&op.PropertyNounValue,
		&op.NounPropertyValue,
		&op.AspectsAreTraits,
	} {
		if next := *input; //
		m.Match(q, &next) /* && len(next) == 0 */ {
			if !okay || len(next) < len(best) {
				best = next
				ret, okay = m, true
				if len(best) == 0 {
					break
				}
			}
		}
	}
	if okay {
		*input = best
	}
	return
}

func (op *VerbNamesAreNames) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Verb.Match(q, &next) &&
		op.Names.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.OtherNames.Match(q, &next) {
		// q.note("matched VerbNamesAreNames")
		*input, okay = next, true
	}
	return
}

func (op *VerbNamesAreNames) Generate(rar Registrar) error {
	return applyResults(rar, op.compile())
}

func (op *VerbNamesAreNames) compile() localResults {
	return makeResult(
		op.Verb.Macro, !op.Verb.Macro.Reversed,
		op.Names.GetNames(nil, nil),
		op.OtherNames.GetNames(nil, nil),
	)
}

func (op *NamesVerbNames) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Names.Match(q, &next) &&
		!op.Names.HasAnonymousKind() &&
		op.Are.Match(q, &next) &&
		op.Verb.Match(q, &next) &&
		op.OtherNames.Match(q, &next) {
		// q.note("matched NamesVerbNames")
		*input, okay = next, true
	}
	return
}

func (op *NamesVerbNames) Generate(rar Registrar) error {
	return applyResults(rar, op.compile())
}

func (op *NamesVerbNames) compile() (ret localResults) {
	return makeResult(
		op.Verb.Macro,
		op.Verb.Macro.Reversed,
		op.Names.GetNames(nil, nil),
		op.OtherNames.GetNames(nil, nil),
	)
}

func (op *NamesAreLikeVerbs) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Names.Match(q, &next) &&
		!op.Names.HasAnonymousKind() &&
		op.Are.Match(q, &next) &&
		op.Adjectives.Match(q, &next) {
		// q.note("matched NamesAreLikeVerbs")
		Optional(q, &next, &op.VerbPhrase)
		*input, okay = next, true
	}
	return
}

func (op *NamesAreLikeVerbs) GetOtherNames() (ret Iterator) {
	if v := op.VerbPhrase; v != nil {
		ret = v.Names.Iterate()
	}
	return
}

func (op *NamesAreLikeVerbs) GetMacro() (ret Macro) {
	if v := op.VerbPhrase; v != nil {
		ret = v.Verb.Macro
	}
	return
}

func (op *NamesAreLikeVerbs) IsReversed() (okay bool) {
	if v := op.VerbPhrase; v != nil {
		okay = v.Verb.Macro.Reversed
	}
	return
}

func (op *NamesAreLikeVerbs) Generate(rar Registrar) error {
	return applyResults(rar, op.compile())
}

func (op *NamesAreLikeVerbs) compile() localResults {
	var b []resultName
	for it := op.GetOtherNames(); it.HasNext(); {
		n := it.GetNext()
		b = append(b, n.GetName(nil, nil))
	}
	a := op.Names.GetNames(op.Adjectives.Reduce())
	return makeResult(
		op.GetMacro(), op.IsReversed(),
		a, b)
}
