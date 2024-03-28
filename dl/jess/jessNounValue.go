package jess

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

func (op *Property) String() string {
	return op.Matched
}

func (op *Property) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.matchProperty(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *Property) matchProperty(q Query, input *InputState) (okay bool) {
	if m, width := q.FindField(input.Words()); width > 0 {
		op.Matched, *input, okay = m, input.Skip(width), true
	}
	return
}

// unexpectedly, runs in the fallback phase;
// because, to match inform, if the named noun doesn't exist yet
// this creates a thing with that name.
func (op *PropertyNounValue) Phase() Phase {
	return weave.FallbackPhase
}

func (op *PropertyNounValue) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.Property.Match(q, &next) &&
		op.Of.Match(q, &next, keywords.Of) &&
		op.NamedNoun.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.SingleValue.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *PropertyNounValue) Generate(ctx *Context) (err error) {
	return generatePropertyPhrase(ctx, op.NamedNoun, op.Property, op.SingleValue)
}

// like PropertyNounValue, runs in FallbackPhase; see notes there.
func (op *NounPropertyValue) Phase() Phase {
	return weave.FallbackPhase
}

func (op *NounPropertyValue) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.NamedNoun.Match(q, &next) &&
		op.Has.Match(q, &next, keywords.Has) &&
		(Optional(q, &next, &op.Article) || true) &&
		op.Property.Match(q, &next) &&
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

func (op *NounPropertyValue) Generate(ctx *Context) error {
	return generatePropertyPhrase(ctx, op.NamedNoun, op.Property, op.SingleValue)
}

func generatePropertyPhrase(ctx *Context, n NamedNoun, p Property, v SingleValue) (err error) {
	if ns, e := n.BuildNouns(ctx, NounProperties{}); e != nil {
		err = e
	} else if e := tryAsThings(ctx, ns); e != nil {
		err = e
	} else {
		err = genNounValues(ctx, ns, func(n string) error {
			// fix: can i add this to "desired noun" instead of as a callback
			return ctx.AddNounValue(n, p.String(), v.Assignment())
		})
	}
	return
}

// try to apply one of the passed kinds to each of the desired nouns
// the first one not to generate a conflict succeeds.
func generateFallbacks(ctx *Context, ns []DesiredNoun, kinds ...string) error {
	return ctx.PostProcess(weave.FallbackPhase, func() (err error) {
	Loop:
		for _, n := range ns {
			for _, k := range kinds {
				if e := ctx.AddNounKind(n.Noun, k); e == nil || errors.Is(e, mdl.Duplicate) {
					err = nil // applying a duplicate kind is considered a success
					break Loop
				} else {
					err = e // keep one of the conflicts; only cleared on success
					if !errors.Is(e, mdl.Conflict) {
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
func tryAsThings(ctx *Context, ns []DesiredNoun) (err error) {
	return ctx.PostProcess(weave.FallbackPhase, func() (err error) {
		for _, n := range ns {
			e := ctx.AddNounKind(n.Noun, Things)
			if e != nil && !errors.Is(e, mdl.Conflict) && !errors.Is(e, mdl.Duplicate) {
				err = e
				break
			}
		}
		return
	})
}
