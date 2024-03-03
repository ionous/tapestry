package jess

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/match"
)

// KindsOf
func (op *KindsOf) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Names.Match(AddContext(q, ExcludeNounMatching), &next) &&
		op.Are.Match(q, &next) &&
		op.matchKindsOf(&next) &&
		(Optional(q, &next, &op.Traits) || true) &&
		op.Kind.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

// match "a kind of" or "kinds of"
func (op *KindsOf) matchKindsOf(input *InputState) (okay bool) {
	if m, width := kindsSpan.FindPrefix(input.Words()); m != nil {
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

// ex. "... are a kind of aspect."
// panics if not matched.
func (op *KindsOf) IsAspect() bool {
	return op.Kind.ActualKind.base == kindsOf.Aspect
}

// The closed containers called safes are a kind of fixed in place thing.
func (op *KindsOf) Generate(rar Registrar) (err error) {
	if parent, e := op.Kind.Validate(); e != nil {
		err = e
	} else {
		traits := op.GetTraits()
		isPlural, isAspect := op.Are.IsPlural(), op.IsAspect()
		// the names are kinds we have not yet created
		for it := op.Names.Iterate(); it.HasNext(); {
			at := it.GetNext()
			if at.CountedKind != nil {
				err = errors.New(countedKindMsg)
				break
			} else {
				// determine the name of the desired kind
				// ex. the rhs of "the k called desired kind"
				var kind string
				if k := at.Kind; k != nil {
					// if it was a known kind, then that's easy.
					kind = k.ActualKind.name
				} else {
					// otherwise, get the specified name
					if n := getKindOfName(at); isAspect && isPlural {
						// ick. aspects are expected to be singular
						kind = rar.GetSingular(n)
					} else if !isAspect && !isPlural {
						// all other kinds are are expected to be plural
						kind = rar.GetPlural(n)
					} else {
						kind = n
					}
				}
				// register our new kind ( or new kind of hierarchy )
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
						} else if e := AddKindTraits(rar, kind, called.GetTraits()); e != nil {
							err = e
							break
						}
					}
					// add trailing traits.
					if e := AddKindTraits(rar, kind, traits); e != nil {
						err = e
						break
					}
				}
			}
		}
	}
	return
}

func getKindOfName(at *Names) string {
	var name *Name
	if n := at.Name; n != nil {
		name = n
	} else if kc := at.KindCalled; kc != nil {
		// we excluded existing nouns; so only names must exist
		name = kc.NamedNoun.Name
	}
	if name == nil {
		panic("unexpected match")
	}
	return name.GetNormalizedName()
}

const countedKindMsg = `Defining a new kind using a leading number is prohibited. 
If you're sure you'd like a number to be part of the name, use "called the" instead.`
