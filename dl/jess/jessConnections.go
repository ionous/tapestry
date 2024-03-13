package jess

import (
	"errors"
)

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
func (op *MapConnections) Generate(rar *Context) (err error) {
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
