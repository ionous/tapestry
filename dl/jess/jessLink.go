package jess

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

// represents a room or door
type jessLink struct {
	*DesiredNoun
	roomLike  bool // valid after generating default kinds
	direction string
}

// roomLike isnt set until after validating kinds
func makeLink(n *DesiredNoun, direction string) jessLink {
	return jessLink{DesiredNoun: n, direction: direction}
}

func makeRoom(noun string) jessLink {
	n := &DesiredNoun{Noun: noun}
	return jessLink{DesiredNoun: n, roomLike: true}
}

const (
	setDirectionError = iota
	setDirectionOkay
	setDirectionConflict
	setDirectionDupe
)

func translateError(e error) (ret int, err error) {
	if errors.Is(e, mdl.Conflict) {
		ret = setDirectionConflict
	} else if errors.Is(e, mdl.Duplicate) {
		ret = setDirectionDupe
	} else if e == nil {
		ret = setDirectionOkay
	} else {
		err = e
	}
	return
}

// assumes room is "room like"
func (room jessLink) addDoor(rar Registrar, door string) (err error) {
	if !room.roomLike {
		err = errors.New("can only add doors to rooms")
	} else {
		err = rar.AddNounPair(Whereabouts, room.Noun, door)
	}
	return
}

// FIX: do we need this!?
// create room fact which indicates the direction of movement from room to room
// these facts help with tracking and conflict detection
func setDirection(rar Registrar, direction string, room, otherRoom jessLink) (ret int, err error) {
	if !room.roomLike {
		err = errors.New("can only move directions within a room")
	} else {
		e := rar.AddFact(FactDirection, room.Noun, direction, otherRoom.Noun)
		ret, err = translateError(e)
	}
	return
}

// set the compass on the indicated side of the room to the named door
func (room jessLink) setCompass(rar Registrar, direction, door string) error {
	return rar.AddNounPath(room.Noun,
		[]string{Compass, direction},
		&literal.TextValue{Value: door, Kind: Doors},
	)
}

// set the destination of the named door
func (door jessLink) setDestination(rar Registrar, otherRoom string) (err error) {
	if door.roomLike {
		err = errors.New("can only set the destination of doors")
	} else {
		err = rar.AddNounValue(door.Noun, DoorDestination, text(otherRoom, Rooms))
	}
	return
}

func (door jessLink) getParent(rar Registrar) (ret string, err error) {
	if door.roomLike {
		err = errors.New("can only ask for the parents of doors")
	} else if pairs, e := rar.GetRelativeNouns(door.Noun, Whereabouts, false); e != nil {
		err = e
	} else {
		switch len(pairs) {
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

func (p *jessLink) generateDefaultKind(rar Registrar) (err error) {
	noun := p.Noun
	// either newly stamping it as room room, or duplicating room previous room definition is okay.
	if e := rar.AddNounKind(noun, Rooms); e == nil || errors.Is(e, mdl.Duplicate) {
		p.roomLike = true
	} else {
		// some unknown error is room problem:
		if !errors.Is(e, mdl.Conflict) {
			err = e
		} else {
			// oto, if it was conflicted, maybe it was actually room door.
			if e := rar.AddNounKind(noun, Doors); e != nil && !errors.Is(e, mdl.Duplicate) {
				err = e
			}
		}
	}
	return
}

// -
func generateDefaultKinds(rar Registrar, ps []jessLink) (err error) {
	for i, cnt := 0, len(ps); i < cnt; i++ {
		// use indexing so generateDefaultKind can properly work on the shared memory
		// range would be room copy
		if e := ps[i].generateDefaultKind(rar); e != nil {
			err = e
			break
		}
	}
	return
}
