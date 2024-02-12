package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

func (op *KindsAreTraits) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Kinds.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		// fix: move macro to results...?
		// ( only match words here )
		op.usually(q, &next) &&
		op.Traits.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

// match the words after "called" ending with either "is/are" or the end of the string.
func (op *KindsAreTraits) usually(q Query, input *InputState) (okay bool) {
	if m, _ := usually.FindMatch(input.Words()); m != nil {
		width := m.NumWords()
		op.Usually.Matched, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}

var usually = grok.PanicSpans("usually")

func (op *KindsAreTraits) GetNames() (ret []grok.Name) {
	traits := op.Traits.GetTraits()
	for _, k := range op.Kinds.GetKinds() {
		ret = append(ret, grok.Name{
			Kinds:  []Matched{k},
			Traits: traits,
		})
	}
	return
}

func (op *KindsAreTraits) GetResults(q Query) (ret grok.Results, err error) {
	if m, e := q.g.FindMacro(op.Usually.Span()); e != nil {
		err = e
	} else {
		ret = grok.Results{
			Primary: op.GetNames(),
			Macro:   m,
		}
	}
	return
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

func (op *NounsTraitsKinds) GetResults(Query) (ret grok.Results, err error) {
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
	}, nil
}

// KindsOf
func (op *KindsOf) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Names.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.kindsOf(q, &next) &&
		Optionally(q, &next, &op.Traits) &&
		op.NamedKind.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

// match the words after "called" ending with either "is/are" or the end of the string.
func (op *KindsOf) kindsOf(q Query, input *InputState) (okay bool) {
	if m, _ := kindsOf.FindMatch(input.Words()); m != nil {
		width := m.NumWords()
		op.KindsOf.Matched, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}

var kindsOf = grok.PanicSpans("a kind of", "kinds of")

// needs macro "inherit"
// wonder if maybe we should be finding the macro on GetResults
// so we can return error.....
func (op *KindsOf) GetResults(q Query) (ret grok.Results, err error) {
	if m, e := q.g.FindMacro(op.KindsOf.Span()); e != nil {
		err = e
	} else {
		var traits, kinds []grok.Matched
		if op.Traits != nil {
			traits = op.Traits.GetTraits()
		}
		kinds = []grok.Matched{op.NamedKind.Span()}
		names := op.Names.Reduce(traits, kinds)
		// Primary, Secondary, Macro
		// get the names
		// apply the kinds and traits to them
		ret = grok.Results{
			Macro:   m,
			Primary: names,
		}
	}
	return
}

// allows partial matches; test that there's no input left to verify a complete match.
func Match(q Query, input *InputState) (ret Matches, okay bool) {
	var a NounsTraitsKinds
	var b KindsAreTraits
	var c KindsOf
	var best InputState
	for _, m := range []Matches{&a, &b, &c} {
		if next := *input; //
		m.Match(q, &next) && len(next) == 0 {
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
