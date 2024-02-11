package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

func (op *KindName) Match(q Query, input *InputState) (okay bool) {
	if _, width := q.FindKind(*input); width > 0 {
		op.Matched, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}

// its interesting that we dont have to store anything else
// all the trait info is in this... even additional traits.
func (op *Kinds) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	Optionally(q, &next, &op.Article) &&
		op.KindName.Match(q, &next) &&
		Optionally(q, &next, &op.AdditionalKinds) {
		*input, okay = next, true
	}
	return
}

// for backwards compatibility, kinds are sometimes understood as names
func (op *Kinds) GetKinds() []grok.Matched {
	var out []grok.Matched
	for t := *op; ; {
		out = append(out, t.KindName.Matched)
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
// 			Span:    t.KindName.Matched.(Span),
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
	*KindName
}

// its interesting that we dont have to store anything else
// all the trait info is in this... even additional traits.
func (op *TraitsKind) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Traits.Match(q, &next) &&
		Optionally(q, &next, &op.KindName) {
		*input, okay = next, true
	}
	return
}

func (op *TraitsKind) GetTraitSet() grok.TraitSet {
	var k grok.Matched
	if op.KindName != nil {
		k = op.KindName.Matched
	}
	return grok.TraitSet{
		Kind:   k,
		Traits: op.Traits.GetTraits(),
	}
}

func (op *KindsAreTraits) Match(q Query, input *InputState) (okay bool) {
	t := Zt_KindsAreTraits
	if i := t.TermIndex("usually"); i < 0 {
		panic("missing typeinfo")
	} else if next := *input; //
	op.Kinds.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.Usually.Match(q, &next, t.Terms[i].Markup) &&
		op.Traits.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

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
