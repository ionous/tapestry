package jess

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/support/match"
)

func (op *Understandings) Generate(rar Registrar) (err error) {
	return errors.New("not implemented")
}

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
