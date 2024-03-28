package jess

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/weave"
)

// -------------------------------------------------------------------------
// Verb
// -------------------------------------------------------------------------

// the passed phrase is the macro to match
func (op *Verb) Match(q Query, input *InputState) (okay bool) {
	if m, width := q.FindNoun(input.Words(), Verbs); width > 0 {
		op.Text = m // holds the normalized name
		*input, okay = input.Skip(width), true
	}
	return
}

// -------------------------------------------------------------------------
// VerbPhrase
// -------------------------------------------------------------------------

func (op *VerbPhrase) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Verb.Match(q, &next) &&
		op.PlainNames.Match(AddContext(q, PlainNameMatching), &next) {
		*input, okay = next, true
	}
	return
}

// -------------------------------------------------------------------------
// VerbNamesAreNames
// -------------------------------------------------------------------------

// runs in the NounPhase phase
func (op *VerbNamesAreNames) Phase() Phase {
	return weave.NounPhase
}
func (op *VerbNamesAreNames) GetNouns() Names {
	return op.OtherNames // reverse left and right sides
}
func (op *VerbNamesAreNames) GetOtherNouns() Names {
	return op.Names
}
func (op *VerbNamesAreNames) GetAdjectives() (_ Adjectives) {
	return
}
func (op *VerbNamesAreNames) GetVerb() string {
	return op.Verb.Text
}
func (op *VerbNamesAreNames) Generate(rar *Context) error {
	return generateVerbPhrase(rar, op)
}
func (op *VerbNamesAreNames) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Verb.Match(q, &next) &&
		op.Names.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.OtherNames.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

// -------------------------------------------------------------------------
// NamesVerbNames
// -------------------------------------------------------------------------
// runs in the NounPhase phase
func (op *NamesVerbNames) Phase() Phase {
	return weave.NounPhase
}
func (op *NamesVerbNames) GetNouns() Names {
	return op.Names
}
func (op *NamesVerbNames) GetOtherNouns() Names {
	return op.OtherNames
}
func (op *NamesVerbNames) GetAdjectives() (_ Adjectives) {
	return
}
func (op *NamesVerbNames) GetVerb() string {
	return op.Verb.Text
}
func (op *NamesVerbNames) Generate(rar *Context) error {
	return generateVerbPhrase(rar, op)
}
func (op *NamesVerbNames) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Names.Match(q, &next) &&
		!op.Names.HasAnonymousKind() &&
		op.Are.Match(q, &next) &&
		op.Verb.Match(q, &next) &&
		op.OtherNames.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

// -------------------------------------------------------------------------
// NamesAreLikeVerbs
// -------------------------------------------------------------------------
// runs in the NounPhase phase
func (op *NamesAreLikeVerbs) Phase() Phase {
	return weave.NounPhase
}
func (op *NamesAreLikeVerbs) GetNouns() Names {
	return op.Names
}
func (op *NamesAreLikeVerbs) GetOtherNouns() (ret Names) {
	if v := op.VerbPhrase; v != nil {
		ret = v.PlainNames
	}
	return
}
func (op *NamesAreLikeVerbs) GetAdjectives() Adjectives {
	return op.Adjectives
}
func (op *NamesAreLikeVerbs) GetVerb() (ret string) {
	if v := op.VerbPhrase; v != nil {
		ret = v.Verb.Text
	}
	return
}
func (op *NamesAreLikeVerbs) Generate(rar *Context) error {
	return generateVerbPhrase(rar, op)
}
func (op *NamesAreLikeVerbs) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Names.Match(q, &next) &&
		!op.Names.HasAnonymousKind() &&
		op.Are.Match(q, &next) &&
		op.Adjectives.Match(q, &next) {
		Optional(q, &next, &op.VerbPhrase)
		*input, okay = next, true
	}
	return
}

// -------------------------------------------------------------------------
// jessVerbPhrase
// helper to generate nouns;
// could be a struct, or even multiple parameters instead of a pull/getter
// not sure what's best: can other phrases use the same patter... whatever it is?
// ( esp. seeing as they all need to schedule to the same phase )
// -------------------------------------------------------------------------

// fix? this interface means that Names can contain zero matches.
type jessVerbPhrase interface {
	GetNouns() Names
	GetOtherNouns() Names
	GetAdjectives() Adjectives
	GetVerb() string
}

func generateVerbPhrase(ctx *Context, p jessVerbPhrase) (err error) {
	if props, e := p.GetAdjectives().Reduce(); e != nil {
		err = e
	} else if lhs, e := p.GetNouns().BuildNouns(ctx, props); e != nil {
		err = e
	} else if rhs, e := p.GetOtherNouns().BuildNouns(ctx, NounProperties{}); e != nil {
		err = e
	} else if e := genNounValues(ctx, lhs, nil); e != nil {
		err = e
	} else {
		if verbName := p.GetVerb(); len(verbName) > 0 {
			if e := genNounValues(ctx, rhs, nil); e != nil {
				err = e
			} else {
				err = applyVerb(ctx, verbName, lhs, rhs)
			}
		} else if len(rhs) > 0 {
			err = errors.New("missing verb")
		} else {
			err = tryAsThings(ctx, lhs)
		}
	}
	return
}

// note: some phrases "the box is open" dont have macros.
// in that case, genNounValues itself does all the work.
func applyVerb(ctx *Context, verbName string, lhs, rhs []DesiredNoun) (err error) {
	return ctx.PostProcess(weave.VerbPhrase, func() (err error) {
		if v, e := readVerb(ctx, verbName); e != nil {
			err = e
		} else {
			err = v.applyVerb(ctx, lhs, rhs)
		}
		return
	})
}
