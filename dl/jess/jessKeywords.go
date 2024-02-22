package jess

import (
	"git.sr.ht/~ionous/tapestry/support/match"
)

// only valid after matching
func (op *Are) IsPlural() bool {
	span := op.Matched.(match.Span) //yikes.
	return span[0].Hash() == keywords.Are
}

func (op *Are) Match(q Query, input *InputState) (okay bool) {
	if width := input.MatchWord(keywords.Are, keywords.Is); width > 0 {
		op.Matched = input.Cut(width)
		*input = input.Skip(width)
		okay = true
	}
	return
}

func (op *CommaAnd) Match(q Query, input *InputState) (okay bool) {
	if sep, e := ReadCommaAnd(input.Words()); e == nil && sep != 0 {
		width := sep.Len()
		op.Matched = input.Cut(width)
		*input = input.Skip(width)
		okay = true
	}
	return
}

func (op *Words) Match(q Query, input *InputState, hashes ...uint64) (okay bool) {
	if width := input.MatchWord(hashes...); width > 0 {
		op.Matched = input.Cut(width)
		*input = input.Skip(width)
		okay = true
	}
	return
}

func (op *Words) Span() match.Span {
	return op.Matched.(match.Span)
}

// make customizable?
var keywords = struct {
	And, Are, Called, Comma, Has, Is, Of, Quote uint64
}{
	And:    match.Hash("and"),
	Are:    match.Hash("are"),
	Called: match.Hash("called"),
	Comma:  match.Hash(","),
	Has:    match.Hash("has"),
	Is:     match.Hash("is"),
	Of:     match.Hash("of"),
	Quote:  match.Hash(`"`),
}
