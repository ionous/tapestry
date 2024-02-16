package jess

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/support/grok"
)

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

func (op *Kind) GetName(traits, kinds []Matched) (ret grok.Name) {
	return grok.Name{
		Traits: traits,
		// the order of kinds matters for "kinds of"
		// for: A container is a kind of thing.
		// the kinds should appear in that order in this list:
		Kinds: append([]Matched{op.Matched}, kinds...),
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

// exists for testing the matching of trait sets
// not used directly during normal matching.
func MatchTraits(q Query, in InputState) (ret grok.TraitSet, err error) {
	var traits Traits
	if next := in; //
	!traits.Match(q, &next) {
		err = errors.New("failed to match traits")
	} else {
		// after the traits, there might be a kind
		var kind *Kind
		Optional(q, &next, &kind)
		if cnt := len(next); cnt != 0 {
			err = fmt.Errorf("partially matched %d traits", len(in)-cnt)
		} else {
			var k grok.Matched
			if kind != nil {
				k = kind.Matched
			}
			ret = grok.TraitSet{
				Kind:   k,
				Traits: traits.GetTraits(),
			}
		}
	}
	return
}
