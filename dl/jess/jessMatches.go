package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

func (op *KindsAreTraits) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Kinds.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.Usually.Match(q, &next, usually) &&
		op.Traits.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

var usually = grok.PanicSpan("usually")

func (op *KindsAreTraits) GetMatch() ([]grok.Name, grok.Macro) {
	var out []grok.Name
	traits := op.Traits.GetTraits()
	for _, k := range op.Kinds.GetKinds() {
		out = append(out, grok.Name{
			Kinds:  []Matched{k},
			Traits: traits,
		})
	}
	return out, op.Usually.Macro
}

func (op *KindsAreTraits) GetResults() (ret grok.Results) {
	ns, m := op.GetMatch()
	return grok.Results{Primary: ns, Macro: m}
}

// NounsTraitsKinds
func (op *NounsTraitsKinds) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Names.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		len(next) > 0 && // something needs to be after "is/are"
		Optionally(q, &next, &op.Traits) &&
		Optionally(q, &next, &op.CommaAnd) &&
		Optionally(q, &next, &op.Kinds) {
		*input, okay = next, true
	}
	return
}

func (op *NounsTraitsKinds) GetResults() (ret grok.Results) {
	var traits, kinds []grok.Matched
	if op.Traits != nil {
		traits = op.Traits.GetTraits()
	}
	if op.Kinds != nil {
		kinds = op.Kinds.GetKinds()
	}
	names := op.Names.Reduce(traits, kinds)
	// Primary, Secondary, Macro
	// get the names
	// apply the kinds and traits to them
	return grok.Results{
		Primary: names,
	}
}
