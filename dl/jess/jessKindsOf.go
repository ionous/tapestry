package jess

import (
	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

func (op *KindsOf) GetResults() (ret grok.Results, err error) {
	panic("xxx")
}

// KindsOf
func (op *KindsOf) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Names.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.kindsOf(q, &next) &&
		(Optional(q, &next, &op.Traits) || true) &&
		op.Kind.Match(q, &next) {
		q.note("matched KindsOf")
		*input, okay = next, true
	}
	return
}

// match the words after "called" ending with either "is/are" or the end of the string.
func (op *KindsOf) kindsOf(q Query, input *InputState) (okay bool) {
	if m, _ := kindsOf.FindMatch(input.Words()); m != nil {
		width := m.NumWords()
		op.KindsOf.Matched, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}

var kindsOf = grok.PanicSpans("a kind of", "kinds of")

func (op *KindsOf) GetTraits() (ret Traitor) {
	if op.Traits != nil {
		ret = op.Traits.GetTraits()
	}
	return
}

// The closed containers called safes are a kind of fixed in place thing.
func (op *KindsOf) Apply(rar Registrar) (err error) {
	parent := op.Kind.Matched
	traits := op.GetTraits()
	//
	for it := op.Names.Iterate(); it.HasNext(); {
		k := it.GetNext() // a 'Names'
		name := inflect.Normalize(k.String())
		if e := rar.AddKind(name, parent.String()); e != nil {
			err = e
			break
		} else {
			// "x called" can have its own traits and kind
			if called := op.Names.KindCalled; called != nil {
				if e := rar.AddKind(name, called.Kind.String()); e != nil {
					err = e // kind called already normalized because it matched the specific kind
					break
				} else if e := AddTraitsToKind(rar, name, called.GetTraits()); e != nil {
					err = e
					break
				}
			}
			// add trailing traits.
			if e := AddTraitsToKind(rar, name, traits); e != nil {
				err = e
				break
			}
		}
	}
	return
}
