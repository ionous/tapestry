package jess

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/match"
)

// panics if unmatched
func (op *SingleValue) Assignment() (ret rt.Assignment) {
	if n := op.QuotedText; n != nil {
		ret = n.Assignment()
	} else if n := op.MatchingNumber; n != nil {
		ret = n.Assignment()
	} else {
		panic("unmatched assignment")
	}
	return
}
func (op *SingleValue) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	Optional(q, &next, &op.QuotedText) ||
		Optional(q, &next, &op.MatchingNumber) {
		*input, okay = next, true
	}
	return
}

func (op *QuotedText) Assignment() rt.Assignment {
	return text(op.Matched.String(), "")
}

// match combines double quoted and backtick text:
// generating a leading "QuotedText" indicator
// and a single "word" containing the entire quoted text.
func (op *QuotedText) Match(q Query, input *InputState) (okay bool) {
	if width := input.MatchWord(match.Keywords.QuotedText); width > 0 {
		next := input.Skip(width) // skip over the quote indicator (1 word)
		op.Matched, *input, okay = next.Cut(1), next.Skip(1), true
	}
	return
}

func (op *MatchingNumber) Assignment() rt.Assignment {
	return number(op.Number, "")
}

func (op *MatchingNumber) Match(q Query, input *InputState) (okay bool) {
	if ws := input.Words(); len(ws) > 0 {
		word := ws[0].String()
		if v, ok := WordsToNum(word); ok && v > 0 {
			const width = 1
			op.Number = float64(v)
			*input, okay = input.Skip(width), true
		}
	}
	return
}

// tbd: i'm not sold on the idea that registar takes assignments
// maybe it'd make more sense to pass in generic "any" values,
// to have add factory functions to Registrar,
// or to have individual methods for the necessary types
// ( maybe just three: trait, text, number )
func text(value, kind string) rt.Assignment {
	return &assign.FromText{
		Value: &literal.TextValue{Value: value, Kind: kind},
	}
}

func number(value float64, kind string) rt.Assignment {
	return &assign.FromNumber{
		Value: &literal.NumValue{Value: value, Kind: kind},
	}
}
