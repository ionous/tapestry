package jess

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// runs in the PropertyPhase; relies on the named kind existing.
func (op *KindsAreEither) Phase() weaver.Phase {
	return weaver.PropertyPhase
}

func (op *KindsAreEither) MatchLine(q JessContext, line InputState) (ret InputState, okay bool) {
	if next := line; //
	op.Kind.Match(q, &next) &&
		op.matchEither(&next) &&
		op.Traits.Match(q, &next) {
		ret, okay = next, true
	}
	return
}

// match "can be", "are either", etc.
func (op *KindsAreEither) matchEither(input *InputState) (okay bool) {
	if m, width := canBeEither.FindPrefix(input.Words()); m != nil {
		op.CanBe.Matched = input.Cut(width)
		*input, okay = input.Skip(width), true
	}
	return
}

func (op *KindsAreEither) Weave(w weaver.Weaves, run rt.Runtime) (err error) {
	if k, e := op.Kind.Validate(kindsOf.Kind); e != nil {
		err = e
	} else {
		if op.Traits.NewTrait == nil {
			// mdl is smart enough to generate "not" aspects from bool fields
			if name, e := match.NormalizeAll(op.Traits.Matched); e != nil {
				err = e
			} else {
				err = w.AddKindFields(k, []mdl.FieldInfo{{
					Name:     name,
					Affinity: affine.Bool,
				}})
			}
		} else {
			if name, e := op.generateAspect(w); e != nil {
				err = e
			} else {
				err = w.AddKindFields(k, []mdl.FieldInfo{{
					Name:     name,
					Affinity: affine.Text,
					Class:    name,
				}})
			}
		}
	}
	return
}

func (op *KindsAreEither) generateAspect(w weaver.Weaves) (ret string, err error) {
	first := op.Traits
	if aspect, e := match.NormalizeAll(append(first.Matched, statusToken)); e != nil {
		err = e
	} else if e := w.AddKind(aspect, kindsOf.Aspect.String()); e != nil {
		err = e
	} else if traits, e := normalizeTraits(first); e != nil {
		err = e
	} else if e := w.AddAspectTraits(aspect, traits); e != nil {
		err = e
	} else {
		ret = aspect
	}
	return
}

func normalizeTraits(ts NewTrait) (ret []string, err error) {
	for it := ts.Iterate(); it.HasNext(); {
		if str, e := match.NormalizeAll(it.GetNext().Matched); e != nil {
			err = e
			break
		} else {
			ret = append(ret, str)
		}
	}
	return
}

var statusToken = match.TokenValue{Token: match.String, Value: "status"}
var canBeEither = match.PanicSpans("can be", "are either", "is either", "can be either")

func (op *NewTrait) Match(q JessContext, input *InputState) (okay bool) {
	if next := *input; next.Len() > 0 {
		// look for 1) the end of the string, or 2) the separator "or"
		if firstSpan := scanUntil(next.words, keywords.Or); firstSpan < 0 {
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

type NewTraitIterator struct {
	next *NewTrait
}

func (it NewTraitIterator) HasNext() bool {
	return it.next != nil
}

func (it *NewTraitIterator) GetNext() (ret *NewTrait) {
	ret, it.next = it.next, it.next.NewTrait
	return
}
