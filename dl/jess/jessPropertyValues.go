package jess

func (op *PropertyValues) Next() (ret *PropertyValues) {
	if next := op.AdditionalPropertyValues; next != nil {
		ret = &next.Values
	}
	return
}

func (op *AdditionalPropertyValues) Match(q Query, kind string, input *InputState) (okay bool) {
	if next := *input; //
	op.CommaAnd.Match(q, &next) &&
		op.Values.Match(q, kind, &next) {
		*input, okay = next, true
	}
	return
}

// valid properties are filtered by kind.
// ( which allows the values on the rhs of the property to be unquoted nouns and kinds )
func (op *PropertyValues) Match(q Query, kind string, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.Property.Match(q, kind, &next) &&
		(op.matchOf(q, &next) || true) && // optionally the wod "of"
		op.Value.Match(q, &next) {
		//
		var optional AdditionalPropertyValues
		if optional.Match(q, kind, &next) {
			op.AdditionalPropertyValues = &optional
		}
		*input, okay = next, true
	}
	return
}

func (op *PropertyValues) matchOf(q Query, input *InputState) (okay bool) {
	var w Words
	if w.Match(q, input, keywords.Of) {
		op.Of, okay = &w, true
	}
	return
}
