package jess

func (op *NamedNoun) GetNormalizedName() (ret string) {
	if n := op.Noun; n != nil {
		ret = n.ActualNoun // the actual name is already normalized
	} else if n := op.Name; n != nil {
		ret = n.GetNormalizedName()
	} else {
		panic("NamedNoun was unmatched")
	}
	return
}

// panics if not matched
func (op *NamedNoun) BuildNouns(ctx *Context, ts, ks []string) ([]DesiredNoun, error) {
	return buildNounsFrom(ctx, ts, ks, ref(op.Noun), ref(op.Name))
}

func (op *NamedNoun) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	Optional(q, &next, &op.Noun) ||
		Optional(q, &next, &op.Name) {
		*input, okay = next, true
	}
	return
}
