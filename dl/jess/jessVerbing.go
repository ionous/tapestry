package jess

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// -------------------------------------------------------------------------
// Verb
// -------------------------------------------------------------------------

// the passed phrase is the macro to match
func (op *Verb) Match(q JessContext, input *InputState) (okay bool) {
	kind := Verbs
	if m, width := q.FindNoun(input.Words(), &kind); width > 0 {
		op.Text = m // holds the normalized name
		*input, okay = input.Skip(width), true
	}
	return
}

// -------------------------------------------------------------------------
// VerbPhrase
// -------------------------------------------------------------------------

func (op *VerbPhrase) Match(q JessContext, input *InputState) (okay bool) {
	if next := *input; //
	op.Verb.Match(q, &next) &&
		op.PlainNames.Match(AddContext(q, PlainNameMatching), &next) {
		*input, okay = next, true
	}
	return
}

// -------------------------------------------------------------------------
// VerbNamesAreNames
// *In* the lobby is a supporter.
// -------------------------------------------------------------------------

// runs in the NounPhase phase
func (op *VerbNamesAreNames) Phase() weaver.Phase {
	return weaver.NounPhase
}
func (op *VerbNamesAreNames) GetNouns() *MultipleNames {
	return &op.OtherNames // reverse left and right sides
}
func (op *VerbNamesAreNames) GetOtherNouns() *MultipleNames {
	return &op.Names
}
func (op *VerbNamesAreNames) GetAdjectives() (_ MultipleAdjectives) {
	return
}
func (op *VerbNamesAreNames) GetVerb() string {
	return op.Verb.Text
}

func (op *VerbNamesAreNames) Generate(ctx JessContext) error {
	return generateVerbPhrase(ctx, op)
}

func (op *VerbNamesAreNames) MatchLine(q JessContext, line InputState) (ret InputState, okay bool) {
	if next, q := line, //
		AddContext(q, MatchKindsOfKinds); //
	op.Verb.Match(q, &next) &&
		op.Names.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.OtherNames.Match(q, &next) {
		ret, okay = next, true
	}
	return
}

// -------------------------------------------------------------------------
// NamesVerbNames
// A container called the trunk is *in* the lobby.
// -------------------------------------------------------------------------
// matches in the NounPhase phase
func (op *NamesVerbNames) Phase() weaver.Phase {
	return weaver.NounPhase
}
func (op *NamesVerbNames) GetNouns() *MultipleNames {
	return &op.Names
}
func (op *NamesVerbNames) GetOtherNouns() *MultipleNames {
	return &op.OtherNames
}
func (op *NamesVerbNames) GetAdjectives() (_ MultipleAdjectives) {
	return
}
func (op *NamesVerbNames) GetVerb() string {
	return op.Verb.Text
}
func (op *NamesVerbNames) Generate(ctx JessContext) (err error) {
	if op.Names.HasAnonymousKind() {
		err = errors.New("can't start phrase with an anonymous leading kind")
	} else {
		err = generateVerbPhrase(ctx, op)
	}
	return
}
func (op *NamesVerbNames) MatchLine(q JessContext, line InputState) (ret InputState, okay bool) {
	if next, q := line, //
		AddContext(q, MatchKindsOfKinds|MatchPronouns);
	// like NamesAreLikeVerbs, this limits lhs matching to kinds which can be instanced
	op.Names.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.Verb.Match(q, &next) &&
		op.OtherNames.Match(q, &next) {
		ret, okay = next, true
	}
	return
}

// -------------------------------------------------------------------------
// NamesAreLikeVerbs
// The bottle is closed [*in* the library.]
// the adjectives ( traits and kinds ) are required;
// the verb and other noun are optional.
// -------------------------------------------------------------------------
// matches in the NounPhase phase
func (op *NamesAreLikeVerbs) Phase() weaver.Phase {
	return weaver.NounPhase
}
func (op *NamesAreLikeVerbs) GetNouns() *MultipleNames {
	return &op.Names
}
func (op *NamesAreLikeVerbs) GetOtherNouns() (ret *MultipleNames) {
	if v := op.VerbPhrase; v != nil {
		ret = &v.PlainNames
	}
	return
}
func (op *NamesAreLikeVerbs) GetAdjectives() MultipleAdjectives {
	return op.Adjectives
}
func (op *NamesAreLikeVerbs) GetVerb() (ret string) {
	if v := op.VerbPhrase; v != nil {
		ret = v.Verb.Text
	}
	return
}
func (op *NamesAreLikeVerbs) Generate(ctx JessContext) (err error) {
	if op.Names.HasAnonymousKind() {
		err = errors.New("can't start phrase with an anonymous leading kind")
	} else {
		err = generateVerbPhrase(ctx, op)
	}
	return
}

func (op *NamesAreLikeVerbs) MatchLine(q JessContext, in InputState) (ret InputState, okay bool) {
	if next, q := in, //
		AddContext(q, MatchKindsOfKinds); //
	op.Names.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.Adjectives.Match(q, &next) {
		Optional(q, &next, &op.VerbPhrase)
		ret, okay = next, true
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
	Phase() weaver.Phase
	GetNouns() *MultipleNames
	GetOtherNouns() *MultipleNames
	GetAdjectives() MultipleAdjectives
	GetVerb() string
}

func generateVerbPhrase(ctx JessContext, p jessVerbPhrase) error {
	return ctx.Schedule(p.Phase(), func(w weaver.Weaves, run rt.Runtime) (err error) {
		if props, e := p.GetAdjectives().Reduce(); e != nil {
			err = e
		} else if lhs, e := p.GetNouns().BuildNouns(ctx, w, run, props); e != nil {
			err = e
		} else if rhs, e := p.GetOtherNouns().BuildNouns(ctx, w, run, NounProperties{}); e != nil {
			err = e
		} else if e := genNounValues(ctx, lhs, nil); e != nil {
			err = e
		} else {
			if cnt := len(lhs); cnt > 0 {
				if cnt == 1 {
					n := lhs[0]
					ctx.SetTopic(ActualNoun{n.Noun, n.CreatedKind})
				} else {
					ctx.RejectTopic(errors.New("pronouns expect a single noun"))
				}
			}
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
	})
}

// note: some phrases "the box is open" dont have macros.
// in that case, genNounValues itself does all the work.
func applyVerb(ctx JessContext, verbName string, lhs, rhs []DesiredNoun) error {
	return ctx.Schedule(weaver.VerbPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if v, e := readVerb(run, verbName); e != nil {
			err = e
		} else {
			err = v.applyVerb(ctx, w, lhs, rhs)
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
