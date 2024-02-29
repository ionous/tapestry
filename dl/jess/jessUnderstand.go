package jess

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/support/match"
)

func (op *Understand) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Understand.Match(q, &next, keywords.Understand) &&
		op.QuotedTexts.Match(q, &next) &&
		op.As.Match(q, &next, keywords.As) &&
		(Optional(q, &next, &op.Article) || true) &&
		(op.matchPluralOf(q, &next) || true) &&
		op.Names.Match(AddContext(q, IncludeExistingNouns), &next) {
		*input, okay = next, true
	}
	return
}

func (op *Understand) matchPluralOf(q Query, input *InputState) (okay bool) {
	if m, width := pluralOf.FindMatch(input.Words()); m != nil {
		op.PluralOf, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}

var pluralOf = match.PanicSpans("plural of")

func (op *Understand) Generate(rar Registrar) (err error) {
	if len(op.PluralOf) > 0 {
		err = op.applyPlurals(rar)
	} else {
		// check whether kind matches an action
		// ( although inform appears to eval the lhs first to see it matches any parser statement )
		// will need to understand (ahg. puns) that better.
		if actions, nouns, e := op.readRhs(rar); e != nil {
			err = e
		} else if len(actions) > 0 && len(nouns) > 0 {
			err = errors.New("jess doesn't support mixing noun and action understandings")
		} else if len(actions) > 0 {
			err = op.applyActions(rar, actions)
		} else if len(nouns) > 0 {
			err = op.applyAliases(rar, nouns)
		} else {
			err = errors.New("what's there to understand?")
		}
	}
	return
}

func (op *Understand) applyActions(rar Registrar, actions []string) (err error) {
Loop:
	for it := op.QuotedTexts.Iterate(); it.HasNext(); {
		phrase := it.GetNext()
		if m, e := BuildPhrase(phrase); e != nil {
			err = e
		} else {
			for _, act := range actions {
				if e := rar.AddGrammar(phrase, &grammar.Directive{
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

func (op *Understand) applyAliases(rar Registrar, nouns []string) (err error) {
	// for every noun on the rhs
	for _, noun := range nouns {
		if noun, e := rar.GetClosestNoun(inflect.Normalize(noun)); e != nil {
			err = e
			break
		} else {
			//  add the alias specified on the lhs
			for it := op.QuotedTexts.Iterate(); it.HasNext(); {
				alias := it.GetNext()
				if alias = inflect.Normalize(alias); len(alias) > 0 {
					// the -1 indicates that this is an alias; hrm.
					if e := rar.AddNounAlias(noun, alias, -1); e != nil {
						err = e
						break
					}
				}
			}
		}
	}
	return
}

func (op *Understand) readRhs(rar Registrar) (actions, nouns []string, err error) {
	for it := op.Names.Iterate(); it.HasNext(); {
		next := it.GetNext()
		if k := next.Kind; k != nil {
			if kt := k.ActualKind.base; kt != kindsOf.Action {
				err = fmt.Errorf("expected an action, but %q is a %s", k.String(), kt.String())
			} else {
				actions = append(actions, k.String())
			}
		} else {
			nouns = append(nouns, next.String())
		}
	}
	return
}

func (op *Understand) applyPlurals(rar Registrar) (err error) {
Loop:
	for as := op.Names.Iterate(); as.HasNext(); {
		// determine the "single" side of the plural request
		if n := as.GetNext(); n.Noun == nil {
			err = fmt.Errorf("unknown name %q, expected the name of an existing noun.", n)
		} else {
			name := n.Noun.ActualNoun
			for it := op.QuotedTexts.Iterate(); it.HasNext(); {
				plural := it.GetNext()
				if e := rar.AddPlural(plural, name); e != nil {
					err = e
					break Loop
				}
			}
		}
	}
	return
}
