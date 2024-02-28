package jess

func (op *Property) String() string {
	return op.Matched
}

func (op *Property) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.matchProperty(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *Property) matchProperty(q Query, input *InputState) (okay bool) {
	if m, width := q.FindField(input.Words()); width > 0 {
		op.Matched, *input, okay = m, input.Skip(width), true
	}
	return
}

func (op *PropertyNounValue) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.Property.Match(q, &next) &&
		op.Of.Match(q, &next, keywords.Of) &&
		op.Noun.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.SingleValue.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *PropertyNounValue) Generate(rar Registrar) (err error) {
	if n, e := rar.GetClosestNoun(op.Noun.String()); e != nil {
		err = e
	} else {
		err = rar.AddNounValue(n, op.Property.String(), op.SingleValue.Assignment())
	}
	return
}

func (op *NounPropertyValue) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Noun.Match(q, &next) &&
		op.Has.Match(q, &next, keywords.Has) &&
		(Optional(q, &next, &op.Article) || true) &&
		op.Property.Match(q, &next) &&
		(op.matchOf(q, &next) || true) &&
		op.SingleValue.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *NounPropertyValue) matchOf(q Query, input *InputState) (okay bool) {
	var w Words
	if w.Match(q, input, keywords.Of) {
		op.Of, okay = &w, true
	}
	return
}

func (op *NounPropertyValue) Generate(rar Registrar) (err error) {
	if n, e := rar.GetClosestNoun(op.Noun.String()); e != nil {
		err = e
	} else {
		err = rar.AddNounValue(n, op.Property.String(), op.SingleValue.Assignment())
	}
	return
}
