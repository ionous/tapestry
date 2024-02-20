package jess

func (op *NounValue) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.matchProperty(q, &next) && // ends with "of"
		op.Noun.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.SingleValue.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *NounValue) Generate(rar Registrar) (err error) {
	return rar.AddNounValue(op.Noun.String(), op.Property.String(), op.SingleValue.Assignment())
}

// read until "of" is found
func (op *NounValue) matchProperty(q Query, input *InputState) (okay bool) {
	if width := scanUntil(input.Words(), keywords.Of); width > 0 {
		// skip one more than the width to account for the size of the keyword (of)
		op.Property, *input, okay = input.Cut(width), input.Skip(width+1), true
	}
	return
}
