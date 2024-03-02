package jess

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

func (op *KindsAreEither) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Kind.Match(q, &next) &&
		op.matchEither(q, &next) &&
		op.Traits.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

// match "can be", "are either", etc.
func (op *KindsAreEither) matchEither(q Query, input *InputState) (okay bool) {
	if m, width := canBeEither.FindMatch(input.Words()); m != nil {
		op.CanBe.Matched, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}

func (op *KindsAreEither) Generate(rar Registrar) (err error) {
	if k, e := op.Kind.Validate(kindsOf.Kind); e != nil {
		err = e
	} else {
		if op.Traits.NewTrait == nil {
			// mdl is smart enough to generate "not" aspects from bool fields
			name := inflect.Normalize(op.Traits.String())
			err = rar.AddFields(k, []mdl.FieldInfo{{
				Name:     name,
				Affinity: affine.Bool,
			}})
		} else {
			if name, e := op.generateAspect(rar); e != nil {
				err = e
			} else {
				err = rar.AddFields(k, []mdl.FieldInfo{{
					Name:     name,
					Affinity: affine.Text,
					Class:    name,
				}})
			}
		}
	}
	return
}

func (op *KindsAreEither) generateAspect(rar Registrar) (ret string, err error) {
	first := op.Traits
	aspect := inflect.Join([]string{first.String(), "status"})
	if e := rar.AddKind(aspect, kindsOf.Aspect.String()); e != nil {
		err = e
	} else {
		var traits []string
		for it := first.Iterate(); it.HasNext(); {
			traits = append(traits, it.GetNext())
		}
		if e := rar.AddTraits(aspect, traits); e != nil {
			err = e
		} else {
			ret = aspect
		}
	}
	return
}

var canBeEither = match.PanicSpans("can be", "are either", "is either", "can be either")

func (op *NewTrait) Match(q Query, input *InputState) (okay bool) {
	if next := *input; next.Len() > 0 {
		// look for 1) the end of the string, or 2) the separator "or"
		if firstSpan := scanUntil(next.Words(), keywords.Or); firstSpan < 0 {
			width := next.Len()
			op.Matched = next.Cut(width)          // eat everything
			*input, okay = next.Skip(width), true // all done.
		} else {
			// found the word "or":
			op.Matched = next.Cut(firstSpan)
			next = next.Skip(firstSpan + 1)
			// have to match after the "or" for the string to be valid.
			if Optional(q, &next, &op.NewTrait) {
				*input, okay = next, true
			}
		}
	}
	return
}

// unwind the tree of additional kinds
func (op *NewTrait) Iterate() NewTraitIterator {
	return NewTraitIterator{op}
}

func (op *NewTrait) String() string {
	return inflect.Normalize(op.Matched)
}

type NewTraitIterator struct {
	next *NewTrait
}

func (it NewTraitIterator) HasNext() bool {
	return it.next != nil
}

func (it *NewTraitIterator) GetNext() (ret string) {
	ret, it.next = it.next.String(), it.next.NewTrait
	return
}
