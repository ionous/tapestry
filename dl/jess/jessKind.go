package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

func (op *KindName) Match(q Query, input *InputState) (okay bool) {
	if _, width := q.FindKind(*input); width > 0 {
		op.Matched, *input, okay = input.Cut(width), input.Skip(width), true
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
	if next := *input; op.Traits.Match(q, &next) &&
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
	} else if next := *input; Optionally(q, &next, &op.Article) &&
		op.KindName.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.Usually.Match(q, &next, t.Terms[i].Markup) &&
		op.Traits.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

// FIX: multiple kinds
func (op *KindsAreTraits) GetMatch() (retNames []grok.Name, retMacro grok.Macro) {
	n := grok.Name{
		// fix: the current grok tests ( and code ) assume that this phrase produces names
		// not kinds... even though they *are* and have to be kinds
		Span:   op.KindName.Matched.(Span),
		Traits: op.Traits.GetTraits(),
	}
	return []grok.Name{n}, op.Usually.Macro
}

func (op *KindsAreTraits) GetResults() (ret grok.Results) {
	ns, m := op.GetMatch()
	return grok.Results{Primary: ns, Macro: m}
}
