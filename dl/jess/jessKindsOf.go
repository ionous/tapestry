package jess

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// runs in the AncestryPhase phase
func (op *KindsAreKind) Phase() weaver.Phase {
	return weaver.AncestryPhase
}

func (op *KindsAreKind) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Names.Match(AddContext(q, ExcludeNounMatching), &next) &&
		op.Are.Match(q, &next) &&
		op.matchKindsOf(&next) &&
		(Optional(q, &next, &op.Traits) || true) &&
		op.Name.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

// match "a kind of" or "kinds of"
func (op *KindsAreKind) matchKindsOf(input *InputState) (okay bool) {
	if m, width := kindsSpan.FindPrefix(input.Words()); m != nil {
		op.KindsAreKind.Matched = input.Cut(width)
		*input, okay = input.Skip(width), true
	}
	return
}

var kindsSpan = match.PanicSpans("a kind of", "kinds of")

func (op *KindsAreKind) GetTraits() (ret Traitor) {
	if op.Traits != nil {
		ret = op.Traits.GetTraits()
	}
	return
}

// The closed containers called safes are a kind of fixed in place thing.
func (op *KindsAreKind) Generate(ctx Context) error {
	// manually schedule, so we can query FindKind()
	return ctx.Schedule(op.Phase(), func(w weaver.Weaves, run rt.Runtime) (err error) {
		var base kindsOf.Kinds
		span := op.Name.Matched
		if parent, width := ctx.FindKind(span, &base); width != len(span) {
			err = fmt.Errorf("%w kind %s", weaver.Missing, span.DebugString())
		} else {
			traits := op.GetTraits()
			isPlural, isAspect := op.Are.IsPlural(), base == kindsOf.Aspect
			// the names are kinds we have not yet created
			for it := op.Names.GetNames(); it.HasNext(); {
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
						kind = k.actualKind.Name
					} else {
						// otherwise, get the specified name
						if n := getKindOfName(at); isAspect && isPlural {
							// ick. aspects are expected to be singular
							kind = run.SingularOf(n)
						} else if !isAspect && !isPlural {
							// all other kinds are are expected to be plural
							kind = run.PluralOf(n)
						} else {
							kind = n
						}
					}
					// register our new kind ( or new kind of hierarchy )
					if e := w.AddKind(kind, parent); e != nil {
						err = e
						break
					} else {
						// "x called" can have its own traits and kind
						if called := op.Names.KindCalled; called != nil {
							if calledKind, e := called.GetKind(); e != nil {
								err = e
							} else if e := w.AddKind(kind, calledKind); e != nil {
								err = e // kind called already normalized because it matched the specific kind
								break
							} else if e := AddKindTraits(w, kind, called.GetTraits()); e != nil {
								err = e
								break
							}
						}
						// add trailing traits.
						if e := AddKindTraits(w, kind, traits); e != nil {
							err = e
							break
						}
					}
				}
			}
		}
		return
	})
}

func getKindOfName(at *Names) (ret string) {
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
	if n, e := name.GetNormalizedName(); e != nil {
		panic(e)
	} else {
		return n
	}
}

const countedKindMsg = `defining a new kind using a leading number is prohibited. 
If you're sure you'd like a number to be part of the name, use "called the" instead`
