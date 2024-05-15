package jess

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// runs in the LanguagePhase phase
func (op *Understand) Phase() weaver.Phase {
	return weaver.LanguagePhase
}

func (op *Understand) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Understand.Match(q, &next, keywords.Understand) &&
		op.QuotedTexts.Match(q, &next) &&
		op.As.Match(q, &next, keywords.As) &&
		(Optional(q, &next, &op.Article) || true) &&
		(op.matchPluralOf(&next) || true) &&
		op.Names.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *Understand) matchPluralOf(input *InputState) (okay bool) {
	if m, width := pluralOf.FindPrefix(input.Words()); m != nil {
		op.PluralOf = m.String()
		*input, okay = input.Skip(width), true
	}
	return
}

var pluralOf = match.PanicSpans("plural of")

func (op *Understand) Generate(ctx Context) error {
	return ctx.Schedule(op.Phase(), func(w weaver.Weaves, run rt.Runtime) (err error) {
		if len(op.PluralOf) > 0 {
			err = op.applyPlurals(w)
		} else {
			// fix: parse lhs first, into a map keyed by its string
			// then we can error better when strings or grammars appear on the wrong side.
			// (and probably simplify some)
			err = ctx.Schedule(weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
				// check whether kind matches an action or noun
				if actions, nouns, e := op.readRhs(ctx); e != nil {
					err = e
				} else if len(actions) > 0 && len(nouns) > 0 {
					err = errors.New("jess doesn't support mixing noun and action understandings")
				} else if len(actions) > 0 {
					err = op.applyActions(w, actions)
				} else if len(nouns) > 0 {
					err = op.applyAliases(w, nouns)
				} else {
					err = errors.New("what's there to understand?")
				}
				return
			})
		}
		return
	})
}

func (op *Understand) applyActions(w weaver.Weaves, actions []string) (err error) {
Loop:
	for it := op.QuotedTexts.Iterate(); it.HasNext(); {
		phrase := it.GetNext()
		if m, e := BuildPhrase(phrase); e != nil {
			err = e
		} else {
			for _, act := range actions {
				if e := w.AddGrammar(phrase, &grammar.Directive{
					Name: phrase,
					Series: []grammar.ScannerMaker{
						m, &grammar.Action{Action: act},
					}}); e != nil {
					err = e
					break Loop
				}
			}
		}
	}
	return
}

// fix: should this work through desired noun instead?
func (op *Understand) applyAliases(w weaver.Weaves, rhsNouns []string) (err error) {
	// for every noun on the rhs
	for _, noun := range rhsNouns {
		//  add the alias specified on the lhs
		for it := op.QuotedTexts.Iterate(); it.HasNext(); {
			alias := it.GetNext()
			if alias = inflect.Normalize(alias); len(alias) > 0 {
				// the -1 indicates that this is an alias; hrm.
				if e := w.AddNounName(noun, alias, -1); e != nil {
					err = e
					break
				}
			}
		}
	}
	return
}

func (op *Understand) readRhs(q Query) (actions, nouns []string, err error) {
	for it := op.Names.GetNames(); it.HasNext(); {
		next := it.GetNext()
		if n := next.Noun; n != nil {
			nouns = append(nouns, n.ActualNoun.Name)
		} else if k := next.Kind; k != nil && k.ActualKind.BaseKind == kindsOf.Action {
			actions = append(actions, k.ActualKind.Name)
		} else if n := next.Name; n == nil {
			err = errors.New("understandings can only match existing nouns or existing actions")
			break
		} else {
			// fix? if we're going to check at the end; maybe shift to plain names instead.
			// fix? pass a filter to FindNoun so you can't understand things that arent objects.
			if noun, width := q.FindNoun(n.Matched, nil); width < 0 {
				err = fmt.Errorf("no noun found called %q", n.Matched.DebugString())
				break
			} else {
				nouns = append(nouns, noun)
			}
		}
	}
	return
}

func (op *Understand) applyPlurals(w weaver.Weaves) (err error) {
Loop:
	for as := op.Names.GetNames(); as.HasNext(); {
		// determine the "single" side of the plural request
		if n := as.GetNext(); n.Noun == nil {
			err = errors.New("unknown name, expected the name of an existing noun")
		} else {
			name := n.Noun.ActualNoun.Name
			for it := op.QuotedTexts.Iterate(); it.HasNext(); {
				plural := it.GetNext()
				if e := w.AddPlural(plural, name); e != nil {
					err = e
					break Loop
				}
			}
		}
	}
	return
}
