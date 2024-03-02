package jess

// implements MatchedName
func (op *NamedNoun) BuildNoun(traits, kinds []string) (DesiredNoun, error) {
	return op.GetMatchedName().BuildNoun(traits, kinds)
}

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
func (op *NamedNoun) GetMatchedName() (ret MatchedName) {
	if n := op.Noun; n != nil {
		ret = n
	} else if n := op.Name; n != nil {
		ret = n
	} else {
		panic("NamedNoun was unmatched")
	}
	return
}

func (op *NamedNoun) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	Optional(q, &next, &op.Noun) ||
		Optional(q, &next, &op.Name) {
		*input, okay = next, true
	}
	return
}
