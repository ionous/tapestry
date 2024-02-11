package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

func (op *NamedKind) Match(q Query, input *InputState) (okay bool) {
	if m, width := q.FindKind(*input); width > 0 {
		// we want to return the matched kind, not the span because
		// it might have additional info about the match ( ex. a db key )
		op.Matched, *input, okay = m, input.Skip(width), true
	}
	return
}

func (op *Kinds) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	Optionally(q, &next, &op.Article) &&
		op.NamedKind.Match(q, &next) &&
		Optionally(q, &next, &op.AdditionalKinds) {
		*input, okay = next, true
	}
	return
}

// for backwards compatibility, kinds are sometimes understood as names
func (op *Kinds) GetKinds() []grok.Matched {
	var out []grok.Matched
	for t := *op; ; {
		out = append(out, t.NamedKind.Matched)
		if next := t.AdditionalKinds; next == nil {
			break
		} else {
			t = next.Kinds
		}
	}
	return out
}

// // for backwards compatibility, kinds are sometimes understood as names
// func (op *Kinds) GetNames() []grok.Name {
// 	var out []grok.Name
// 	for t := *op; ; {
// 		var art grok.Article
// 		if a := t.Article; a != nil {
// 			art.Matched = a.Matched
// 		}
// 		out = append(out, grok.Name{
// 			Article: art,
// 			Span:    t.NamedKind.Matched.(Span),
// 		})
// 		if next := t.AdditionalKinds; next == nil {
// 			break
// 		} else {
// 			t = next.Kinds
// 		}
// 	}
// 	return out
// }

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
	*NamedKind
}

// its interesting that we dont have to store anything else
// all the trait info is in this... even additional traits.
func (op *TraitsKind) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Traits.Match(q, &next) &&
		Optionally(q, &next, &op.NamedKind) {
		*input, okay = next, true
	}
	return
}

func (op *TraitsKind) GetTraitSet() grok.TraitSet {
	var k grok.Matched
	if op.NamedKind != nil {
		k = op.NamedKind.Matched
	}
	return grok.TraitSet{
		Kind:   k,
		Traits: op.Traits.GetTraits(),
	}
}
