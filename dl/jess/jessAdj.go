package jess

func (op *Adjectives) Match(q Query, input *InputState) (okay bool) {
	next := *input
	traits := Optional(q, &next, &op.Traits)
	if traits {
		Optional(q, &next, &op.CommaAnd)
	}
	kind := Optional(q, &next, &op.Kind)
	if kind {
		Optional(q, &next, &op.AdditionalAdjectives)
	}
	if traits || kind {
		*input, okay = next, true
	}
	return
}
func (op *Adjectives) GetTraits() (ret Traitor) {
	if ts := op.Traits; ts != nil {
		ret = ts.GetTraits()
	}
	return
}

func (op *Adjectives) Reduce() (retTraits, retKinds []Matched) {
	for t := *op; ; {
		for ts := op.GetTraits(); ts.HasNext(); {
			t := ts.GetNext()
			retTraits = append(retTraits, t.Matched)
		}
		if k := t.Kind; k != nil {
			retKinds = append(retKinds, k.Matched)
		}
		if next := t.AdditionalAdjectives; next == nil {
			break
		} else {
			t = next.Adjectives
		}
	}
	return
}

func (op *AdditionalAdjectives) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.CommaAnd.Match(q, &next) &&
		op.Adjectives.Match(q, &next) {
		*input, okay = next, true
	}
	return
}
