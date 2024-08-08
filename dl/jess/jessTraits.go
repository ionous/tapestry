package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// iterator helper
func (op *Traits) Next() (ret *Traits) {
	if next := op.AdditionalTraits; next != nil {
		ret = &next.Traits
	}
	return
}

func (op *AdditionalTraits) Match(q JessContext, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.CommaAnd) || true) &&
		op.Traits.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *Trait) String() string {
	return inflect.Normalize(op.Matched)
}

func (op *Trait) Match(q JessContext, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.matchTrait(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *Trait) matchTrait(q JessContext, input *InputState) (okay bool) {
	if str, width := q.FindTrait(input.Words()); width > 0 {
		op.Matched, *input, okay = str, input.Skip(width), true
	}
	return
}

// its interesting that we dont have to store anything else
// all the trait info is in this... even additional traits.
func (op *Traits) Match(q JessContext, input *InputState) (okay bool) {
	if next := *input; //
	op.Trait.Match(q, &next) {
		Optional(q, &next, &op.AdditionalTraits)
		*input, okay = next, true
	}
	return
}

// unwind the tree of additional traits
func (op *Traits) GetTraits() *Traits {
	return op
}

// fix: is this all trait iteration is being used for?
func ReduceTraits(it *Traits) (ret []string) {
	for ; it != nil; it = it.Next() {
		ret = append(ret, it.Trait.String())
	}
	return
}

// search the input left to right for zero or matching traits;
// regardless of the desired kind.
func TryTraits(q JessContext, in InputState,
	accept func(*Traits, InputState),
	reject func(error),
) {
	q.Try(After(weaver.PropertyPhase), func(weaver.Weaves, rt.Runtime) {
		var traits Traits
		if !traits.Match(q, &in) {
			accept(nil, in) // zero or more so no traits is okay.
		} else {
			accept(&traits, in)
		}
		return
	}, reject)
}
