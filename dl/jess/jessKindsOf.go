package jess

import (
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/support/match"
)

// KindsOf
func (op *KindsOf) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Names.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.matchKindsOf(q, &next) &&
		(Optional(q, &next, &op.Traits) || true) &&
		op.Kind.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

// match "a kind of" or "kinds of"
func (op *KindsOf) matchKindsOf(q Query, input *InputState) (okay bool) {
	if m, _ := kindsSpan.FindMatch(input.Words()); m != nil {
		width := m.NumWords()
		op.KindsOf.Matched, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}

var kindsSpan = match.PanicSpans("a kind of", "kinds of")

func (op *KindsOf) GetTraits() (ret Traitor) {
	if op.Traits != nil {
		ret = op.Traits.GetTraits()
	}
	return
}

// The closed containers called safes are a kind of fixed in place thing.
func (op *KindsOf) Generate(rar Registrar) (err error) {
	if parent, e := op.Kind.Validate(); e != nil {
		err = e
	} else {
		traits := op.GetTraits()
		//
		for it := op.Names.Iterate(); it.HasNext(); {
			k := it.GetNext()
			// ugh. aspects are expected to be singular right now
			// ( see also AspectsAreTraits )
			var name string
			if op.Kind.ActualKind.base == kindsOf.Aspect {
				if n := k.String(); op.Are.IsPlural() {
					name = rar.GetSingular(n)
				} else {
					name = n
				}
			} else {
				if n := k.String(); op.Are.IsPlural() {
					name = n
				} else {
					name = rar.GetPlural(n)
				}
			}
			kind := inflect.Normalize(name)
			if e := rar.AddKind(kind, parent); e != nil {
				err = e
				break
			} else {
				// "x called" can have its own traits and kind
				if called := op.Names.KindCalled; called != nil {
					if calledKind, e := called.GetKind(); e != nil {
						err = e
					} else if e := rar.AddKind(kind, calledKind); e != nil {
						err = e // kind called already normalized because it matched the specific kind
						break
					} else if e := AddDefaultTraits(rar, kind, called.GetTraits()); e != nil {
						err = e
						break
					}
				}
				// add trailing traits.
				if e := AddDefaultTraits(rar, kind, traits); e != nil {
					err = e
					break
				}
			}
		}
	}
	return
}
