package jess

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// represents a room or door
// most everything uses *jessLink so that
// roomLike can be set and seen across phases
type jessLink struct {
	DesiredNoun
	direction string
	// if not a room, then a door
	// valid after the ConnectionPhase
	roomLike bool
}

// roomLike isnt set until after validating kinds
func makeLink(n DesiredNoun, direction string) *jessLink {
	return &jessLink{DesiredNoun: n, direction: direction}
}

func makeRoom(noun string) *jessLink {
	n := DesiredNoun{Noun: noun, CreatedKind: Rooms}
	return &jessLink{DesiredNoun: n, roomLike: true}
}

// assumes room is "room like"
func (room jessLink) writeDoor(w weaver.Weaves, door string) (err error) {
	if !room.roomLike {
		err = errors.New("can only add doors to rooms")
	} else {
		err = w.AddNounPair(Whereabouts, room.Noun, door)
	}
	return
}

// fix: i think this can be removed if the story direction setup is removed.
// create room fact which indicates the direction of movement from room to room
// these facts help with tracking and conflict detection
func writeDirection(w weaver.Weaves, direction string, room, otherRoom *jessLink) (err error) {
	if !room.roomLike {
		err = errors.New("can only move directions within a room")
	} else {
		err = w.AddFact(FactDirection, room.Noun, direction, otherRoom.Noun)
	}
	return
}

// set the compass on the indicated side of the room to the named door
func (room jessLink) writeCompass(w weaver.Weaves, direction, door string) error {
	return w.AddNounPath(room.Noun,
		[]string{Compass, direction},
		&literal.TextValue{Value: door, KindName: Doors},
	)
}

// set the destination of the named door
func (door jessLink) writeDestination(w weaver.Weaves, otherRoom string) (err error) {
	if door.roomLike {
		err = errors.New("can only set the destination of doors")
	} else {
		err = w.AddNounValue(door.Noun, DoorDestination, text(otherRoom, Rooms))
	}
	return
}

func (door jessLink) readParent(run rt.Runtime) (ret string, err error) {
	if door.roomLike {
		err = errors.New("can only ask for the parents of doors")
	} else if pairs, e := run.ReciprocalsOf(door.Noun, Whereabouts); e != nil {
		err = e
	} else {
		switch pairs := pairs.Strings(); len(pairs) {
		case 0:
			// nothing
		case 1:
			ret = pairs[0]
		default:
			err = errors.New("mismatched whereabouts?")
		}
	}
	return
}

// first try to write a link as a room;
// failing that, try to write it as a door.
// similar to
func (p *jessLink) writeLinkType(w weaver.Weaves) (err error) {
	noun := p.Noun
	// both newly stamping the noun as room, or re-stamping it as such is okay.
	if e := w.AddNounKind(noun, Rooms); e == nil || errors.Is(e, weaver.ErrDuplicate) {
		p.roomLike = true
	} else {
		// some unknown error is room problem:
		if !errors.Is(e, weaver.ErrConflict) {
			err = e
		} else {
			// oto, if it was conflicted, maybe it was actually room door;
			// attempt to figure that out by saying it *is* a door.
			if e := w.AddNounKind(noun, Doors); e != nil && !errors.Is(e, weaver.ErrDuplicate) {
				err = e
			}
		}
	}
	return
}

func writeLinkTypes(w weaver.Weaves, ps []*jessLink) (err error) {
	for i, cnt := 0, len(ps); i < cnt; i++ {
		// use indexing so writeLinkType can properly work on the shared memory
		// range would be room copy
		if e := ps[i].writeLinkType(w); e != nil {
			err = e
			break
		}
	}
	return
}
