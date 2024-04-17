package jess

import (
	"git.sr.ht/~ionous/tapestry/support/match"
)

// only valid after matching
func (op *Are) IsPlural() bool {
	return len(op.Matched) == 1 && op.Matched[0].Hash() == keywords.Are
}

func (op *Are) Match(_ Query, input *InputState) (okay bool) {
	if width := input.MatchWord(keywords.Are, keywords.Is); width > 0 {
		op.Matched = input.Cut(width)
		*input = input.Skip(width)
		okay = true
	}
	return
}

func (op *Called) Match(_ Query, input *InputState) (okay bool) {
	if width := input.MatchWord(keywords.Called); width > 0 {
		op.Matched = input.Cut(width)
		*input = input.Skip(width)
		okay = true
	}
	return
}

func (op *CommaAnd) Match(_ Query, input *InputState) (okay bool) {
	if sep, e := ReadCommaAnd(input.Words()); e == nil && sep != 0 {
		width := sep.Len()
		op.Matched = input.Cut(width)
		*input = input.Skip(width)
		okay = true
	}
	return
}

func (op *CommaAndOr) Match(_ Query, input *InputState) (okay bool) {
	if sep, e := ReadCommaAndOr(input.Words()); e == nil && sep != 0 {
		width := sep.Len()
		op.Matched = input.Cut(width)
		*input = input.Skip(width)
		okay = true
	}
	return
}

func (op *Words) Match(_ Query, input *InputState, hashes ...uint64) (okay bool) {
	if width := input.MatchWord(hashes...); width > 0 {
		op.Matched = input.Cut(width)
		*input = input.Skip(width)
		okay = true
	}
	return
}

// make customizable?
var keywords = struct {
	And, Are, As, Called,
	Colon, Comma,
	Has, Have,
	Is, Nowhere, Of, Or, Quote,
	Rule, Someone,
	Through, Understand, Usually, You uint64
}{
	And:    match.Hash("and"),
	Are:    match.Hash("are"),
	As:     match.Hash("as"),
	Called: match.Hash("called"),
	Colon:  match.Hash(`:`),
	Comma:  match.Hash(","),
	Has:    match.Hash("has"),
	Have:   match.Hash("have"),
	//
	Is:      match.Hash("is"),
	Nowhere: match.Hash("nowhere"),
	Of:      match.Hash("of"),
	Or:      match.Hash("or"),
	Quote:   match.Hash(`"`),
	Rule:    match.Hash("rule"),
	Someone: match.Hash("someone"),
	//
	Through:    match.Hash("through"),
	Understand: match.Hash("understand"),
	Usually:    match.Hash("usually"),
	You:        match.Hash("you"),
}
