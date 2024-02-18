package jess

func (op *VerbLinks) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Verb.Match(q, &next) &&
		op.Names.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.OtherNames.Match(q, &next) {
		// q.note("matched VerbLinks")
		*input, okay = next, true
	}
	return
}

func (op *VerbLinks) Apply(rar Registrar) error {
	return grokNounPhrase(rar, op.GetResults())
}

func (op *VerbLinks) GetResults() localResults {
	return makeResult(
		op.Verb.Macro, !op.Verb.Macro.Reversed,
		op.Names.GetNames(nil, nil),
		op.OtherNames.GetNames(nil, nil),
	)
}

func (op *LinksVerb) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Names.Match(q, &next) &&
		!op.Names.HasAnonymousKind() &&
		op.Are.Match(q, &next) &&
		op.Verb.Match(q, &next) &&
		op.OtherNames.Match(q, &next) {
		// q.note("matched LinksVerb")
		*input, okay = next, true
	}
	return
}

func (op *LinksVerb) Apply(rar Registrar) error {
	return grokNounPhrase(rar, op.GetResults())
}

func (op *LinksVerb) GetResults() (ret localResults) {
	return makeResult(
		op.Verb.Macro,
		op.Verb.Macro.Reversed,
		op.Names.GetNames(nil, nil),
		op.OtherNames.GetNames(nil, nil),
	)
}

func (op *LinksAdjectives) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Names.Match(q, &next) &&
		!op.Names.HasAnonymousKind() &&
		op.Are.Match(q, &next) &&
		op.Adjectives.Match(q, &next) {
		// q.note("matched LinksAdjectives")
		Optional(q, &next, &op.VerbPhrase)
		*input, okay = next, true
	}
	return
}

func (op *LinksAdjectives) GetOtherNames() (ret Iterator) {
	if v := op.VerbPhrase; v != nil {
		ret = v.Names.Iterate()
	}
	return
}

func (op *LinksAdjectives) GetMacro() (ret Macro) {
	if v := op.VerbPhrase; v != nil {
		ret = v.Verb.Macro
	}
	return
}

func (op *LinksAdjectives) IsReversed() (okay bool) {
	if v := op.VerbPhrase; v != nil {
		okay = v.Verb.Macro.Reversed
	}
	return
}

func (op *LinksAdjectives) Apply(rar Registrar) error {
	return grokNounPhrase(rar, op.GetResults())
}

func (op *LinksAdjectives) GetResults() localResults {
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
