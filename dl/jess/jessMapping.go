package jess

import "git.sr.ht/~ionous/tapestry/support/match"

// ----
func (op *MapLocations) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Links.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.DirectionFromLinks.Match(q, &next) {
		Optional(q, &next, &op.AdditionalLinks)
		*input, okay = next, true
	}
	return
}

func (op *MapLocations) Generate(Registrar) (err error) {
	panic("yyy")
}

// ----
func (op *MapDirections) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.DirectionFromLinks.Match(q, &next) &&
		op.Links.Match(q, &next) &&
		op.Redirect.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *MapDirections) Generate(Registrar) (err error) {
	panic("yyy")
}

// ----
func (op *MapConnections) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.matchThrough(q, &next) &&
		op.Doors.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.Links.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *MapConnections) matchThrough(q Query, input *InputState) (okay bool) {
	if width := input.MatchWord(keywords.Through); width > 0 {
		op.Through = true
		*input, okay = input.Skip(width), true
	}
	return
}

func (op *MapConnections) Generate(Registrar) (err error) {
	panic("yyy")
}

// ----
func (op *DirectionFromLinks) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Direction.Match(q, &next) &&
		op.matchFromOf(q, &next) &&
		op.Links.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *DirectionFromLinks) matchFromOf(q Query, input *InputState) (okay bool) {
	if m, width := fromOf.FindMatch(input.Words()); m != nil {
		op.FromOf.Matched, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}

var fromOf = match.PanicSpans("from", "of")

// ----
func (op *Direction) Match(q Query, input *InputState) (okay bool) {
	panic("xxx")
}

// ----
func (op *Links) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.matchNowhere(q, &next) ||
		// tbd. maybe this is better than context flags? i dunno.
		Optional(q, &next, &op.KindCalled) ||
		Optional(q, &next, &op.Noun) ||
		Optional(q, &next, &op.Name) {
		*input, okay = next, true
	}
	return
}

func (op *Links) matchNowhere(q Query, input *InputState) (okay bool) {
	if width := input.MatchWord(keywords.Nowhere); width > 0 {
		op.Nowhere = true
		*input, okay = input.Skip(width), true
	}
	return
}

// ----
func (op *AdditionalLinks) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.CommaAnd.Match(q, &next) &&
		op.Links.Match(q, &next) {
		Optional(q, &next, &op.AdditionalLinks)
		*input, okay = next, true
	}
	return
}

// links iterator
type Linkit struct {
	next *AdditionalLinks
}

func (it Linkit) HasNext() bool {
	return it.next != nil
}

func (it *Linkit) GetNext() (ret Links) {
	ret, it.next = it.next.Links, it.next.AdditionalLinks
	return
}
