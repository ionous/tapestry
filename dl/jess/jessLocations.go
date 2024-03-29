package jess

import (
	"errors"
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// runs in the MappingPhase phase
func (op *MapLocations) Phase() weaver.Phase {
	return weaver.MappingPhase
}

func (op *MapLocations) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Linking.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.DirectionOfLinking.Match(q, &next) {
		Optional(q, &next, &op.AdditionalDirections)
		*input, okay = next, true
	}
	return
}

// return an iterator that is capable of walking over the right hand side of the mapping.
func (op *MapLocations) GetOtherLocations() DirectIt {
	return IterateDirections(&op.DirectionOfLinking, op.AdditionalDirections)
}

func (op *MapLocations) Generate(ctx Context) error {
	return ctx.Schedule(op.Phase(), func(w weaver.Weaves, run rt.Runtime) (err error) {
		if links, e := op.generateLinks(ctx, w, run); e != nil {
			err = e
		} else {
			err = ctx.Schedule(weaver.ConnectionPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
				if e := writeLinkTypes(w, links); e != nil {
					err = e
				} else {
					err = connectPlaceToPlaces(w, run, links[0], links[1:])
				}
				return
			})
		}
		return
	})
}

func (op *MapLocations) generateLinks(q Query, w weaver.Weaves, run rt.Runtime) (ret []*jessLink, err error) {
	if room, e := op.Linking.BuildNoun(q, w, run, NounProperties{}); e != nil {
		err = e
	} else {
		ret = append(ret, makeLink(*room, ""))
		for it := op.GetOtherLocations(); it.HasNext(); {
			link := it.GetNext()
			if el, e := link.buildLink(q, w, run); e != nil {
				err = e
				break
			} else {
				ret = append(ret, el)
			}
		}
	}
	return
}

// connection behavior depends on what's room door and what's room room...
// so the generation of locations and connections needs to be after applying the default kinds
func connectPlaceToPlaces(w weaver.Weaves, run rt.Runtime, src *jessLink, dst []*jessLink) (err error) {
	if !src.roomLike {
		err = connectDoorToRooms(w, run, src, dst)
	} else {
		err = connectRoomToPlaces(w, run, src, dst)
	}
	return
}

// assuming the lhs of the phrase was a door, try the rhs as rooms.
func connectDoorToRooms(w weaver.Weaves, run rt.Runtime, door *jessLink, places []*jessLink) (err error) {
	for i, cnt := 0, len(places); i < cnt && err == nil; i++ {
		if p := places[i]; !p.roomLike {
			err = errors.New("both sides cant be doors")
		} else {
			err = connectDoorToRoom(w, run, door, p, p.direction)
		}
	}
	return
}

// assuming the lhs of the phrase was a room, try the rhs as rooms or doors.
func connectRoomToPlaces(w weaver.Weaves, run rt.Runtime, room *jessLink, places []*jessLink) (err error) {
	for i, cnt := 0, len(places); i < cnt && err == nil; i++ {
		if p := places[i]; !p.roomLike {
			err = connectRoomToDoor(w, run, room, p, p.direction)
		} else {
			err = connectRoomToRoom(w, run, room, p, p.direction)
		}
	}
	return
}

// ex. door R is direction from room P.
// this is limited to the door being in the room
// ( i dont think there's any other interpretations
// - because directions always lead to doors, and doors are always exists. )
func connectDoorToRoom(w weaver.Weaves, run rt.Runtime, door, room *jessLink, direction string) (err error) {
	if e := room.writeDoor(w, door.Noun); e != nil {
		err = e
	} else {
		// doesn't set the room to room direction because we don't have a known destination for the door.
		err = room.writeCompass(w, direction, door.Noun)
	}
	return
}

// room is direction from door:
// there are two different interpretations:
// 1. the door is inside the room on its opposite side.
// 2. the door is in some other room. in that room, moving in the specified direction
// leads to the door, and the door exits into room R.
// inform only handles the first; but it seems either are valid for tapestry.
func connectRoomToDoor(w weaver.Weaves, run rt.Runtime, room, door *jessLink, direction string) (err error) {
	if back, e := readReverse(run, direction); e != nil {
		err = e
	} else if parent, e := door.readParent(run); e != nil {
		err = e // ^ pairs must already have been written
	} else if len(parent) == 0 || parent == room.Noun {
		// put the door in the room; traveling reverse gets us there.
		err = connectDoorToRoom(w, run, door, room, back)
	} else {
		otherRoom := makeRoom(parent)
		if res, e := writeDirection(w, direction, otherRoom, room); e != nil {
			err = e
		} else if res == setDirectionConflict || res == setDirectionDupe {
			err = errors.New("direction already set") // fix? might be to handle dupe in some cases
		} else if e := door.writeDestination(w, room.Noun); e != nil {
			err = e
		} else if e := otherRoom.writeCompass(w, direction, door.Noun); e != nil {
			err = e
		} else {
			_, err = createPrivateDoor(w, back, room, otherRoom)
		}
	}
	return
}

// room R is direction from room P.
// the primary goal is to ensure there's door in room P that leads to room R;
// secondarily, try to put a door in room R leading to P.
func connectRoomToRoom(w weaver.Weaves, run rt.Runtime, room, otherRoom *jessLink, direction string) (err error) {
	if door, e := createPrivateDoor(w, direction, otherRoom, room); e != nil {
		err = e
	} else if len(door) == 0 {
		err = errors.New("room already has a door")
	} else if back, e := readReverse(run, direction); e != nil {
		err = e
	} else {
		// try to create the reverse door; dont worry if it creates nothing.
		_, err = createPrivateDoor(w, back, room, otherRoom)
	}
	return
}

func readReverse(run rt.Runtime, direction string) (ret string, err error) {
	if rev := run.OppositeOf(direction); ret == direction {
		err = fmt.Errorf("couldnt determine the opposite of %q", direction)
	} else {
		ret = rev
	}
	return
}

// generate a door in the first room so that going in the specified direction leads into the other room.
// ex. SOUTH from Lhs is Rhs.
// can return the empty string if there was already room door on that side of the room.
func createPrivateDoor(w weaver.Weaves, direction string, room, otherRoom *jessLink) (ret string, err error) {
	if res, e := writeDirection(w, direction, room, otherRoom); e != nil {
		err = e
	} else if res == setDirectionConflict || res == setDirectionDupe {
		err = errors.New("direction already set") // fix? might be to handle dupe in some cases
	} else {
		// create room magic name: room-direction-door
		door := strings.Replace(room.Noun, " ", "-", -1) + "-" + strings.Replace(direction, " ", "-", -1) + "-door"
		e := room.writeCompass(w, direction, door)
		if res, e := translateError(e); e != nil {
			err = e
		} else if res == setDirectionOkay {
			if e := w.AddNounKind(door, Doors); e != nil {
				err = e
			} else if e := w.AddNounName(door, door, 0); e != nil {
				err = e
			} else if e := w.AddNounTrait(door, Scenery); e != nil {
				err = e
			} else if e := w.AddNounTrait(door, Private); e != nil {
				err = e
			} else if e := w.AddNounValue(door, DoorDestination, text(otherRoom.Noun, Rooms)); e != nil {
				err = e
			} else if e := room.writeDoor(w, door); e != nil {
				err = e
			} else {
				ret = door
			}
		}
	}
	return
}
