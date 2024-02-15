package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

func (op *Kind) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.matchKind(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *Kind) matchKind(q Query, input *InputState) (okay bool) {
	if m, width := q.FindKind(*input); width > 0 {
		// we want to return the matched kind, not the span because
		// it might have additional info about the match ( ex. a db key )
		op.Matched, *input, okay = m, input.Skip(width), true
	}
	return
}

func (op *Kind) Span() grok.Span {
	return op.Matched.(Span)
}

func (op *Kind) GetName(traits, kinds []Matched) (ret grok.Name) {
	return grok.Name{
		Traits: traits,
		Kinds:  append(kinds, op.Matched),
		// no name and no article because, the object itself is anonymous.
		// ( the article associated with the kind gets eaten )
	}
}

func (op *Kinds) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Kind.Match(q, &next) {
		Optional(q, &next, &op.AdditionalKinds)
		*input, okay = next, true
	}
	return
}

// for backwards compatibility, kinds are sometimes understood as names
func (op *Kinds) GetKinds() []grok.Matched {
	var out []grok.Matched
	for t := *op; ; {
		out = append(out, t.Kind.Matched)
		if next := t.AdditionalKinds; next == nil {
			break
		} else {
			t = next.Kinds
		}
	}
	return out
}

func (op *AdditionalKinds) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.CommaAnd.Match(q, &next) &&
		op.Kinds.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

type TraitsKind struct {
	Traits
	*Kind
}

// its interesting that we dont have to store anything else
// all the trait info is in this... even additional traits.
func (op *TraitsKind) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Traits.Match(q, &next) {
		Optional(q, &next, &op.Kind)
		*input, okay = next, true
	}
	return
}

func (op *TraitsKind) GetTraitSet() grok.TraitSet {
	var k grok.Matched
	if op.Kind != nil {
		k = op.Kind.Matched
	}
	return grok.TraitSet{
		Kind:   k,
		Traits: op.Traits.GetTraits(),
	}
}
