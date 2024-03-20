package jess

import "git.sr.ht/~ionous/tapestry/weave"

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
	if ns, e := op.NamedNoun.BuildNouns(ctx, nil, nil); e != nil {
		err = e
	} else {
		err = genNounValues(ctx, ns, func(n string) error {
			return ctx.AddNounValue(n, op.Property.String(), op.SingleValue.Assignment())
		})
	}
	return
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

func (op *NounPropertyValue) Generate(ctx *Context) (err error) {
	if ns, e := op.NamedNoun.BuildNouns(ctx, nil, nil); e != nil {
		err = e
	} else {
		err = genNounValues(ctx, ns, func(n string) error {
			return ctx.AddNounValue(n, op.Property.String(), op.SingleValue.Assignment())
		})
	}
	return
}
