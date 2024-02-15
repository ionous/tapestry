package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

func (op *KindsAreTraits) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Kinds.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.usually(q, &next) &&
		op.Traits.Match(q, &next) {
		q.note("matched KindsAreTraits")
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
		q.note("matched KindsOf")
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
		ret = grok.Results{
			Macro:   m,
			Primary: names,
		}
	}
	return
}

func (op *VerbLinks) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Verb.Match(q, &next) &&
		op.Nouns.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.OtherNouns.Match(q, &next) {
		q.note("matched VerbLinks")
		*input, okay = next, true
	}
	return
}

func (op *VerbLinks) GetResults(Query) (ret grok.Results, _ error) {
	ret = makeResult(
		op.Verb.Macro,
		op.Nouns.Reduce(nil, nil),      // fix?
		op.OtherNouns.Reduce(nil, nil), // fix?
	)
	return
}

func (op *LinksVerb) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Nouns.LimitedMatch(q, &next) &&
		op.Are.Match(q, &next) &&
		op.Verb.Match(q, &next) &&
		op.OtherNouns.Match(q, &next) {
		q.note("matched LinksVerb")
		*input, okay = next, true
	}
	return
}

func (op *LinksVerb) GetResults(Query) (ret grok.Results, _ error) {
	ret = makeResult(
		op.Verb.Macro,
		op.Nouns.Reduce(nil, nil),      // fix?
		op.OtherNouns.Reduce(nil, nil), // fix?
	)
	return
}

func (op *LinksAdjectives) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Nouns.LimitedMatch(q, &next) &&
		op.Are.Match(q, &next) &&
		op.Adjectives.Match(q, &next) {
		q.note("matched LinksAdjectives")
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
	ret = makeResult(m, ones, twos)
	return
}
