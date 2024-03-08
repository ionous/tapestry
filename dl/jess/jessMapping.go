package jess

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/support/match"
)

// ----
func (op *MapDirections) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.DirectionOfLinking.Match(q, &next) &&
		op.Linking.Match(q, &next) &&
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
	op.matchThrough(&next) &&
		op.Doors.Match(q, &next) &&
		(Optional(q, &next, &op.AdditionalLinks) || true) &&
		op.Are.Match(q, &next) &&
		op.Room.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *MapConnections) matchThrough(input *InputState) (okay bool) {
	if width := input.MatchWord(keywords.Through); width > 0 {
		op.Through = true
		*input, okay = input.Skip(width), true
	}
	return
}

// goals of generation:
// 1. ensure the rhs link is a room.
// 2. ensure all of the lhs links are doors
// 3. set the destination of those doors to the rhs room.
func (op *MapConnections) Generate(rar Registrar) (err error) {
	return rar.PostProcess(GenerateNouns, func(q Query) (err error) {
		if room, e := op.Room.GenerateNoun(q, rar, nil, []string{Rooms}); e != nil {
			err = e
		} else {
			for it := op.GetDoors(); it.HasNext(); {
				link := it.GetNext()
				// fix: rather than lists of things;
				// what about passing a "noun properties" instead
				// then the shared function could apply at the right time
				// rather than callers managing the timing.
				if door, e := link.BuildNoun(q, rar, nil, []string{Doors}); e != nil {
					err = e
					break
				} else if door == nil {
					// fix; we can go nowhere.
					err = errors.New("expected at least one door")
					break
				} else {
					if e := rar.PostProcess(GenerateValues, func(Query) (err error) {
						if e := door.generateValues(rar); e != nil {
							err = e
						} else {
							err = rar.AddNounValue(door.Noun, DoorDestination, text(room, Rooms))
						}
						return
					}); e != nil {
						err = e
						break
					}
				}
			}
		}
		return
	})
}

func (op *MapConnections) GetDoors() LinkIt {
	return IterateLinks(&op.Doors, op.AdditionalLinks)
}

// ----
func (op *DirectionOfLinking) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Direction.Match(q, &next) &&
		op.matchFromOf(&next) &&
		op.Linking.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *DirectionOfLinking) BuildNoun(q Query, rar Registrar, ts, ks []string) (*DesiredNoun, error) {
	return op.Linking.BuildNoun(q, rar, ts, ks)
}

func (op *DirectionOfLinking) matchFromOf(input *InputState) (okay bool) {
	if m, width := fromOf.FindPrefix(input.Words()); m != nil {
		op.FromOf.Matched, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}

var fromOf = match.PanicSpans("from", "of")

// ----
func (op *Direction) Match(q Query, input *InputState) (okay bool) {
	// options:
	// 1. look at the fields of the compass
	// 2. look at the noun instances of kind directions
	if m, width := q.FindNoun(input.Words(), Directions); width > 0 {
		op.Text = m
		*input, okay = input.Skip(width), true
	}
	return
}

// ----
func (op *Linking) Match(q Query, input *InputState) (okay bool) {
	if next := *input;        //
	op.matchNowhere(&next) || // tbd. maybe this is better than context flags? i dunno.
		Optional(q, &next, &op.KindCalled) ||
		Optional(q, &next, &op.Noun) ||
		Optional(q, &next, &op.Name) {
		*input, okay = next, true
	}
	return
}

func (op *Linking) matchNowhere(input *InputState) (okay bool) {
	if width := input.MatchWord(keywords.Nowhere); width > 0 {
		op.Nowhere = true
		*input, okay = input.Skip(width), true
	}
	return
}

// generate a room or door; an object if there's not enough information to know; or nil for nowhere.
func (op *Linking) BuildNoun(q Query, rar Registrar, ts, ks []string) (ret *DesiredNoun, err error) {
	if !op.Nowhere {
		if els, e := buildNounsFrom(q, rar, ts, ks, ref(op.KindCalled), ref(op.Noun), ref(op.Name)); e != nil {
			err = e
		} else {
			a := els[0]
			ret = &a
		}
	}
	return
}

// helper since we know there's linking doesnt support counted nouns, but does support nowhere;
// BuildNouns will always return a list of one or none.
func (op *Linking) GenerateNoun(q Query, rar Registrar, ts, ks []string) (ret string, err error) {
	if n, e := op.BuildNoun(q, rar, ts, ks); e != nil {
		err = e
	} else if n != nil {
		if e := rar.PostProcess(GenerateValues, func(Query) error {
			return n.generateValues(rar)
		}); e != nil {
			err = e
		} else {
			ret = n.Noun
		}
	}
	return
}

// ----
func (op *AdditionalLinks) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.CommaAnd.Match(q, &next) &&
		op.Linking.Match(q, &next) {
		Optional(q, &next, &op.AdditionalLinks)
		*input, okay = next, true
	}
	return
}

// links iterator
type LinkIt struct {
	next  *Linking
	queue *AdditionalLinks
}

func IterateLinks(first *Linking, queue *AdditionalLinks) LinkIt {
	return LinkIt{first, queue}
}

func (it LinkIt) HasNext() bool {
	return it.next != nil
}

func (it *LinkIt) GetNext() (ret Linking) {
	ret = *it.next
	if deq := it.queue; deq == nil {
		it.next = nil
	} else {
		it.next = &deq.Linking
		it.queue = deq.AdditionalLinks
	}
	return
}

// ----
func (op *AdditionalDirections) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.CommaAnd.Match(q, &next) &&
		op.DirectionOfLinking.Match(q, &next) {
		Optional(q, &next, &op.AdditionalDirections)
		*input, okay = next, true
	}
	return
}

// additional directions iterator
type DirectIt struct {
	next  *DirectionOfLinking
	queue *AdditionalDirections
}

func IterateDirections(first *DirectionOfLinking, queue *AdditionalDirections) DirectIt {
	return DirectIt{first, queue}
}

func (it DirectIt) HasNext() bool {
	return it.next != nil
}

func (it *DirectIt) GetNext() (ret DirectionOfLinking) {
	ret = *it.next
	if deq := it.queue; deq == nil {
		it.next = nil
	} else {
		it.next = &deq.DirectionOfLinking
		it.queue = deq.AdditionalDirections
	}
	return
}
