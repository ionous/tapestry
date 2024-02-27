package jess

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/support/match"
)

func (op *Understandings) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Understand.Match(q, &next, keywords.Understand) &&
		op.QuotedTexts.Match(q, &next) &&
		op.As.Match(q, &next, keywords.As) &&
		(Optional(q, &next, &op.Article) || true) &&
		(op.matchPluralOf(q, &next) || true) &&
		op.Names.Match(AddContext(q, ExcludeNounCreation), &next) {
		*input, okay = next, true
	}
	return
}

func (op *Understandings) matchPluralOf(q Query, input *InputState) (okay bool) {
	if m, _ := pluralOf.FindMatch(input.Words()); m != nil {
		width := m.NumWords()
		op.PluralOf, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}

var pluralOf = match.PanicSpans("plural of")

func (op *Understandings) Generate(rar Registrar) (err error) {
	if op.PluralOf != nil {
		err = op.generatePluralOf(rar)
	} else {
		// check whether kind matches an action
		// ( although inform actually looks at the left hand side to see it matches any parser statement )
		// will need to understand that better.
	}
	return
}

func (op *Understandings) generatePluralOf(rar Registrar) (err error) {
Loop:
	for as := op.Names.Iterate(); as.HasNext(); {
		// determine the "single" side of the plural request
		var name string
		if n := as.GetNext(); n.Kind != nil {
			err = errors.New("plural kinds not supported") // yet?
			break Loop
		} else if n.Name != nil {
			// in inform these have to be real objects
			// we attempt to do the same thing here....
			// but the parser doesn't really handle plurals all that well yet.
			shortName := n.Name.String()
			if _, e := rar.GetClosestNoun(shortName); e != nil {
				err = e
				break
			} else {
				// for now, after establishing the noun exists;
				// we use the name originally specified for the plurals table.
				name = shortName
			}
		} else {
			panic("unexpected match of name")
		}

		for it := op.QuotedTexts.Iterate(); it.HasNext(); {
			plural := it.GetNext()
			if e := rar.AddPlural(plural, name); e != nil {
				err = e
				break Loop
			}
		}
	}
	return
}
