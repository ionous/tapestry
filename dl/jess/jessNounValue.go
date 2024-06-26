package jess

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// --------------------------------------------------------------
// Property
// --------------------------------------------------------------

func (op *Property) String() string {
	return op.Matched
}

func (op *Property) Match(q Query, kind string, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.matchProperty(q, kind, &next) {
		*input, okay = next, true
	}
	return
}

func (op *Property) matchProperty(q Query, kind string, input *InputState) (okay bool) {
	if m, width := q.FindField(kind, *input); width > 0 {
		op.Matched, *input, okay = m, input.Skip(width), true
	}
	return
}

// --------------------------------------------------------------
// PropertyNounValue
// `The description of the pen is "mightier than the sword.`
// --------------------------------------------------------------

// matches in the noun phase, but mostly runs in the fallback and value phase
func (op *PropertyNounValue) Phase() weaver.Phase {
	return weaver.NounPhase
}
func (op *PropertyNounValue) GetNamedNoun() NamedNoun {
	return op.NamedNoun
}
func (op *PropertyNounValue) GetProperty() Property {
	return op.Property
}
func (op *PropertyNounValue) GetValue() (ret rt.Assignment) {
	if v := op.SingleValue; v != nil {
		ret = v.Assignment()
	} else if v := op.QuotedTexts; v != nil {
		ret = v.Assignment()
	} else {
		panic("not implemented")
	}
	return
}

func (op *PropertyNounValue) Match(q Query, input *InputState) (okay bool) {
	next := *input
	Optional(q, &next, &op.Article)
	if index := scanUntil(next, keywords.Of); index > 0 {
		rest := next.Skip(index + 1) // everything after "of"
		if op.NamedNoun.Match(q, &rest) &&
			op.Are.Match(q, &rest) &&
			// either: single value or quoted texts
			((op.Are.IsPlural() && Optional(q, &rest, &op.QuotedTexts)) ||
				(!op.Are.IsPlural() && Optional(q, &rest, &op.SingleValue))) {

			// try the phrase before the word "of"
			// the whole string must be consumed
			property := InputState(next.Cut(index))
			if op.Property.Match(q, namedKind(op.NamedNoun), &property) && //
				property.Len() == 0 {
				*input, okay = rest, true
			}
		}
	}
	return
}

func (op *PropertyNounValue) Generate(ctx Context) error {
	return genNounValuePhrase(ctx, op)
}

// --------------------------------------------------------------
// NounPropertyValue
// `The pen has (a) description (of) "mightier than the sword.`
// --------------------------------------------------------------

// like PropertyNounValue, matches in the noun phase,
// but mostly runs in the fallback and value phase
func (op *NounPropertyValue) Phase() weaver.Phase {
	return weaver.NounPhase
}
func (op *NounPropertyValue) GetNamedNoun() NamedNoun {
	return op.NamedNoun
}
func (op *NounPropertyValue) GetProperty() Property {
	return op.Property
}
func (op *NounPropertyValue) GetValue() rt.Assignment {
	return op.SingleValue.Assignment()
}

func (op *NounPropertyValue) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.NamedNoun.Match(q, &next) &&
		op.Has.Match(q, &next, keywords.Has) &&
		(Optional(q, &next, &op.Article) || true) &&
		op.Property.Match(q, namedKind(op.NamedNoun), &next) &&
		(op.matchOf(q, &next) || true) &&
		op.SingleValue.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *NounPropertyValue) matchOf(q Query, input *InputState) (okay bool) {
	var w Words
	if w.Match(q, input, keywords.Of) {
		op.Of, okay = &w, true
	}
	return
}

func (op *NounPropertyValue) Generate(ctx Context) error {
	return genNounValuePhrase(ctx, op)
}

// --------------------------------------------------------------
// support for noun property handling
// --------------------------------------------------------------

func namedKind(named NamedNoun) (ret string) {
	if n := named.Noun; n != nil {
		ret = n.actualNoun.Kind
	} else if n := named.Name; n != nil {
		ret = Things // if it hasn't matched this is the default that will be generated
	} else {
		panic("unexpected namedKind")
	}
	return
}

type nounValuePhrase interface {
	Phase() weaver.Phase
	GetNamedNoun() NamedNoun
	GetProperty() Property
	GetValue() rt.Assignment
}

func genNounValuePhrase(ctx Context, phrase nounValuePhrase) (err error) {
	return ctx.Schedule(phrase.Phase(), func(w weaver.Weaves, run rt.Runtime) (err error) {
		n, p, v := phrase.GetNamedNoun(), phrase.GetProperty(), phrase.GetValue()
		if ns, e := n.BuildNouns(ctx, w, run, NounProperties{}); e != nil {
			err = e
		} else if e := tryAsThings(ctx, ns); e != nil {
			err = e
		} else {
			err = genNounValues(ctx, ns, func(n string) error {
				// fix: can i add this to "desired noun" instead of as a callback
				return w.AddNounValue(n, p.String(), v)
			})
		}
		return
	})
}

// try to apply one of the passed kinds to each of the desired nouns
// the first one not to generate a conflict succeeds.
func generateFallbacks(u Scheduler, ns []DesiredNoun, kinds ...string) error {
	return u.Schedule(weaver.FallbackPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
	Loop:
		for _, n := range ns {
			for _, k := range kinds {
				if e := w.AddNounKind(n.Noun, k); e == nil || errors.Is(e, weaver.Duplicate) {
					err = nil // applying a duplicate kind is considered a success
					break Loop
				} else {
					err = e // keep one of the conflicts; only cleared on success
					if !errors.Is(e, weaver.Conflict) {
						break Loop // some other error is an immediate problem
					}
				}
			}
		}
		return
	})
}

// here, we don't care if we aren't able to set "Things"
// this is really and truly a "if nothing else applied" situation.
func tryAsThings(u Scheduler, ns []DesiredNoun) (err error) {
	return u.Schedule(weaver.FallbackPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		for _, n := range ns {
			e := w.AddNounKind(n.Noun, Things)
			if e != nil && !errors.Is(e, weaver.Conflict) && !errors.Is(e, weaver.Duplicate) {
				err = e
				break
			}
		}
		return
	})
}
