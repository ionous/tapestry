package jess

import (
	"git.sr.ht/~ionous/tapestry/support/match"
)

// only valid after matching
func (op *Are) IsPlural() bool {
	return keywordEquals(op.Matched, keywords.Are)
}

func (op *Are) Match(_ Query, input *InputState) (okay bool) {
	if width := input.MatchWord(keywords.Are, keywords.Is); width > 0 {
		op.Matched = input.Cut(width)
		*input = input.Skip(width)
		okay = true
	}
	return
}

// given a phrase, split around the word is or are.
func (op *Are) Split(input InputState) (lhs, rhs InputState, okay bool) {
	if split, ok := keywordSplit(input, keywords.Are, keywords.Is); ok {
		lhs, rhs = split.lhs, split.rhs
		op.Matched = split.matched
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

func (op *Called) Split(in InputState) (lhs, rhs InputState, okay bool) {
	if split, ok := keywordSplit(in, keywords.Called); ok {
		lhs, rhs = split.lhs, split.rhs
		op.Matched = split.matched
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
func (op *Words) Split(in InputState, hashes ...uint64) (lhs, rhs InputState, okay bool) {
	if split, ok := keywordSplit(in, hashes...); ok {
		lhs, rhs = split.lhs, split.rhs
		op.Matched = split.matched
		okay = true
	}
	return
}

func keywordEquals(m Matched, hash uint64) bool {
	return len(m) == 1 && m[0].Hash() == hash
}

// TODO: check boundaries
// expects words before and after the split
func keywordSplit(in InputState, hashes ...uint64) (ret split, okay bool) {
	if start := scanUntil(in.words, hashes...); start > 0 {
		if end := start + 1; end < in.Len() {
			ret, okay = split{
				in.Slice(start),
				in.Skip(end),
				in.words[start:end],
			}, true
		}
	}
	return
}

type split struct {
	lhs, rhs InputState
	matched  Matched
}

// make customizable?
var keywords = struct {
	And, Are, As, Called,
	Colon, Comma,
	Has, Have,
	Is, It, Nowhere, Of, Or, Quote,
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
	It:      match.Hash("it"),
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
