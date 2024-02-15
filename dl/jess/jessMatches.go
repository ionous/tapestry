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

// KindsOf
func (op *KindsOf) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Nouns.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.kindsOf(q, &next) &&
		(Optional(q, &next, &op.Traits) || true) &&
		op.Kind.Match(q, &next) {
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
		kinds = []grok.Matched{op.Kind.Span()}
		names := op.Nouns.Reduce(traits, kinds)
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

// func (op *NounsTraitsKinds) GetResults(Query) (ret grok.Results, err error) {
// 	var traits, kinds []grok.Matched
// 	if op.Traits != nil {
// 		traits = op.Traits.GetTraits()
// 	}
// 	if op.Kinds != nil {
// 		kinds = op.Kinds.GetKinds()
// 	}
// 	names := op.Names.Reduce(traits, kinds)
// 	// Primary, Secondary, Macro
// 	// get the names
// 	// apply the kinds and traits to them
// 	return grok.Results{
// 		Primary: names,
// 	}, nil
// }

func (op *VerbLinks) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Verb.Match(q, &next) &&
		op.Nouns.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.OtherNouns.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *VerbLinks) GetResults(Query) (ret grok.Results, _ error) {
	// fix: swapping
	ret = grok.Results{
		Primary:   op.Nouns.Reduce(nil, nil),      // ASLKLA+
		Secondary: op.OtherNouns.Reduce(nil, nil), // ASLKLA+
		Macro:     op.Verb.Macro,
	}
	return
}

func (op *LinksVerb) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Nouns.LimitedMatch(q, &next) &&
		op.Are.Match(q, &next) &&
		op.Verb.Match(q, &next) &&
		op.OtherNouns.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *LinksVerb) GetResults(Query) (ret grok.Results, _ error) {
	ret = grok.Results{
		Primary:   op.Nouns.Reduce(nil, nil),      // ASLKLA+
		Secondary: op.OtherNouns.Reduce(nil, nil), // ASLKLA+
		Macro:     op.Verb.Macro,
	}
	return
}

func (op *LinksAdjectives) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Nouns.LimitedMatch(q, &next) &&
		op.Are.Match(q, &next) &&
		op.Adjectives.Match(q, &next) {
		Optional(q, &next, &op.VerbPhrase)
		*input, okay = next, true
	}
	return
}

func (op *LinksAdjectives) GetResults(Query) (ret grok.Results, _ error) {
	var m grok.Macro
	var twos []grok.Name
	if c := op.VerbPhrase; c != nil {
		m, twos = c.Verb.Macro, c.Names.Reduce(nil, nil) // traits and kinds?
	}
	ones := op.Nouns.Reduce(op.Adjectives.Reduce())
	ret = grok.Results{
		Primary:   ones,
		Secondary: twos,
		Macro:     m,
	}
	return
}
