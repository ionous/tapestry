package jess

// allows partial matches; test that there's no input left to verify a complete match.
func (op *MatchingPhrases) Match(q Query, input *InputState) (ret Generator, okay bool) {
	// fix? could change to reflect ( or expand type info ) to walk generically
	var best InputState
	for _, m := range []interface {
		Generator
		Interpreter
	}{
		// understand "..." as .....
		&op.Understand,
		// names are "a kind of"/"kinds of" [traits] kind:any.
		&op.KindsOf,
		// kind:objects are "usually" traits.
		&op.KindsAreTraits,
		// kinds:records|objects "have" a ["list of"] number|text|records|objects|aspects ["called a" ...]
		&op.KindsHaveProperties,
		// kinds:objects ("can be"|"are either") new_trait [or new_trait...]
		&op.KindsAreEither,
		// kind:aspects are names
		&op.AspectsAreTraits,
		&op.VerbNamesAreNames,
		&op.NamesVerbNames,
		&op.NamesAreLikeVerbs,
		&op.PropertyNounValue,
		&op.NounPropertyValue,
		&op.MapLocations,
		&op.MapDirections,
		&op.MapConnections,
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

func (op *VerbNamesAreNames) Generate(rar Registrar) (err error) {
	if res, e := op.compile(); e != nil {
		err = e
	} else {
		err = applyResults(rar, res)
	}
	return
}

func (op *VerbNamesAreNames) compile() (ret localResults, err error) {
	if lhs, e := op.Names.GetNames(nil, nil); e != nil {
		err = e
	} else if rhs, e := op.OtherNames.GetNames(nil, nil); e != nil {
		err = e
	} else {
		ret = makeResult(op.Verb.Macro, !op.Verb.Macro.Reversed, lhs, rhs)
	}
	return
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

func (op *NamesVerbNames) Generate(rar Registrar) (err error) {
	if res, e := op.compile(); e != nil {
		err = e
	} else {
		err = applyResults(rar, res)
	}
	return
}

func (op *NamesVerbNames) compile() (ret localResults, err error) {
	if lhs, e := op.Names.GetNames(nil, nil); e != nil {
		err = e
	} else if rhs, e := op.OtherNames.GetNames(nil, nil); e != nil {
		err = e
	} else {
		ret = makeResult(op.Verb.Macro, op.Verb.Macro.Reversed, lhs, rhs)
	}
	return
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

func (op *NamesAreLikeVerbs) GetOtherNames() (ret []resultName, err error) {
	if v := op.VerbPhrase; v != nil {
		ret, err = op.VerbPhrase.Names.GetNames(nil, nil)
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

func (op *NamesAreLikeVerbs) Generate(rar Registrar) (err error) {
	if res, e := op.compile(); e != nil {
		err = e
	} else {
		err = applyResults(rar, res)
	}
	return
}

func (op *NamesAreLikeVerbs) compile() (ret localResults, err error) {
	if rhs, e := op.GetOtherNames(); e != nil {
		err = e
	} else if ts, ks, e := op.Adjectives.Reduce(); e != nil {
		err = e
	} else if lhs, e := op.Names.GetNames(ts, ks); e != nil {
		err = e
	} else {
		ret = makeResult(op.GetMacro(), op.IsReversed(), lhs, rhs)
	}
	return
}
