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
	if m, width := q.FindField(kind, input.words); width > 0 {
		op.Matched, *input, okay = m, input.Skip(width), true
	}
	return
}

// --------------------------------------------------------------
// PropertyNounValue
// `The description of the pen is "mightier than the sword.`
// --------------------------------------------------------------

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

// tbd: should this match pronouns: `The age of it is 42.`?
func (op *PropertyNounValue) MatchLine(q Query, line InputState) (ret InputState, okay bool) {
	next := line
	Optional(q, &next, &op.Article)
	//
	if index := scanUntil(next.words, keywords.Of); index > 0 {
		rest := next.Skip(index + 1) // everything after "of"
		if op.NamedNoun.Match(q, &rest) &&
			op.Are.Match(q, &rest) &&
			// either: single value or quoted texts
			((op.Are.IsPlural() && Optional(q, &rest, &op.QuotedTexts)) ||
				(!op.Are.IsPlural() && Optional(q, &rest, &op.SingleValue))) {

			// try the phrase before the word "of"
			// the whole string must be consumed
			property := next
			property.words = next.Cut(index)
			if op.Property.Match(q, matchedKind(op.NamedNoun), &property) && //
				property.Len() == 0 {
				ret, okay = rest, true
			}
		}
	}
	return
}

func (op *PropertyNounValue) Generate(ctx Context) error {
	return ctx.Schedule(weaver.NounPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if ns, e := op.NamedNoun.BuildNouns(ctx, w, run, NounProperties{}); e != nil {
			err = e
		} else if e := tryAsThings(ctx, ns); e != nil {
			err = e
		} else {
			err = genNounValues(ctx, ns, func(id string) error {
				return w.AddNounValue(id, op.Property.String(), op.GetValue())
			})
		}
		return
	})
}

// --------------------------------------------------------------
// NounPropertyValue
// `The pen has (a) description (of) "mightier than the sword.`
// --------------------------------------------------------------

func (op *NounPropertyValue) MatchLine(q Query, line InputState) (ret InputState, okay bool) {
	if next := line; //
	op.NamedNoun.Match(AddContext(q, MatchPronouns), &next) &&
		op.Has.Match(q, &next, keywords.Has) &&
		op.PropertyValues.Match(q, matchedKind(op.NamedNoun), &next) {
		//
		next.pronouns.setPronounSource(&op.NamedNoun)
		ret, okay = next, true
	}
	return
}

func (op *NounPropertyValue) Generate(ctx Context) error {
	return ctx.Schedule(weaver.NounPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if ns, e := op.NamedNoun.BuildNouns(ctx, w, run, NounProperties{}); e != nil {
			err = e
		} else if e := tryAsThings(ctx, ns); e != nil {
			err = e
		} else {
			err = genNounValues(ctx, ns, func(id string) (err error) {
				for pv := &op.PropertyValues; pv != nil; pv = pv.Next() {
					if e := w.AddNounValue(id,
						pv.Property.String(),
						pv.Value.Assignment()); e != nil {
						err = e
						break
					}
				}
				return
			})
		}
		return
	})
}

// --------------------------------------------------------------
// support for noun property handling
// --------------------------------------------------------------

func matchedKind(named NamedNoun) (ret string) {
	if n := named.Pronoun; n != nil {
		//  pronoun might have to hold a namedNoun so it can do the right thing.
		// or even share constructed noun memory
		ret = Things
	} else if n := named.KindCalled; n != nil {
		ret = n.Kind.actualKind.Name
	} else if n := named.Noun; n != nil {
		ret = n.actualNoun.Kind
	} else if n := named.Name; n != nil {
		ret = Things // if it hasn't matched, this is the default that will be generated
	} else {
		panic("unexpected matchedKind")
	}
	return
}

// try to apply one of the passed kinds to each of the desired nouns
// the first one not to generate a conflict succeeds.
func generateFallbacks(u Scheduler, ns []DesiredNoun, kinds ...string) error {
	return u.Schedule(weaver.FallbackPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
	Loop:
		for _, n := range ns {
			for _, k := range kinds {
				if e := w.AddNounKind(n.Noun, k); e == nil || errors.Is(e, weaver.ErrDuplicate) {
					err = nil // applying a duplicate kind is considered a success
					break Loop
				} else {
					err = e // keep one of the conflicts; only cleared on success
					if !errors.Is(e, weaver.ErrConflict) {
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
			if e != nil && !errors.Is(e, weaver.ErrConflict) && !errors.Is(e, weaver.ErrDuplicate) {
				err = e
				break
			}
		}
		return
	})
}
