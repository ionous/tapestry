package jess

import (
	"errors"
	"strings"

	"git.sr.ht/~ionous/tapestry/weave"
)

// runs in the MappingPhase phase
func (op *MapLocations) Phase() Phase {
	return weave.MappingPhase
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

func (op *MapLocations) Generate(ctx *Context) (err error) {
	if links, e := op.generateLinks(ctx); e != nil {
		err = e
	} else {
		err = ctx.PostProcess(weave.ConnectionPhase, func() (err error) {
			if e := assignDefaultKinds(ctx, links); e != nil {
				err = e
			} else {
				err = connectPlaceToPlaces(ctx, links[0], links[1:])
			}
			return
		})
	}
	return
}
func (op *MapLocations) generateLinks(ctx *Context) (ret []jessLink, err error) {
	if room, e := op.Linking.BuildNoun(ctx, nil, nil); e != nil {
		err = e
	} else {
		ret = append(ret, makeLink(room, ""))
		for it := op.GetOtherLocations(); it.HasNext(); {
			link := it.GetNext()
			if el, e := link.buildLink(ctx); e != nil {
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
func connectPlaceToPlaces(ctx *Context, src jessLink, dst []jessLink) (err error) {
	if !src.roomLike {
		err = connectDoorToRooms(ctx, src, dst)
	} else {
		err = connectRoomToPlaces(ctx, src, dst)
	}
	return
}

// assuming the lhs of the phrase was a door, try the rhs as rooms.
func connectDoorToRooms(ctx *Context, door jessLink, places []jessLink) (err error) {
	for i, cnt := 0, len(places); i < cnt && err == nil; i++ {
		if p := places[i]; !p.roomLike {
			err = errors.New("both sides cant be doors")
		} else {
			err = connectDoorToRoom(ctx, door, p, p.direction)
		}
	}
	return
}

// assuming the lhs of the phrase was a room, try the rhs as rooms or doors.
func connectRoomToPlaces(ctx *Context, room jessLink, places []jessLink) (err error) {
	for i, cnt := 0, len(places); i < cnt && err == nil; i++ {
		if p := places[i]; !p.roomLike {
			err = connectRoomToDoor(ctx, room, p, p.direction)
		} else {
			err = connectRoomToRoom(ctx, room, p, p.direction)
		}
	}
	return
}

// ex. door R is direction from room P.
// this is limited to the door being in the room
// ( i dont think there's any other interpretations
// - because directions always lead to doors, and doors are always exists. )
func connectDoorToRoom(ctx *Context, door, room jessLink, direction string) (err error) {
	if e := room.addDoor(ctx, door.Noun); e != nil {
		err = e
	} else {
		// doesn't set the room to room direction because we don't have a known destination for the door.
		err = room.setCompass(ctx, direction, door.Noun)
	}
	return
}

// room is direction from door:
// there are two different interpretations:
// 1. the door is inside the room on its opposite side.
// 2. the door is in some other room. in that room, moving in the specified direction
// leads to the door, and the door exits into room R.
// inform only handles the first; but it seems either are valid for tapestry.
func connectRoomToDoor(ctx *Context, room, door jessLink, direction string) (err error) {
	if back, e := reverse(ctx, direction); e != nil {
		err = e
	} else if parent, e := door.getParent(ctx); e != nil {
		err = e
	} else if len(parent) == 0 || parent == room.Noun {
		// put the door in the room; traveling reverse gets us there.
		err = connectDoorToRoom(ctx, door, room, back)
	} else {
		otherRoom := makeRoom(parent)
		if res, e := setDirection(ctx, direction, otherRoom, room); e != nil {
			err = e
		} else if res == setDirectionConflict || res == setDirectionDupe {
			err = errors.New("direction already set") // fix? might be to handle dupe in some cases
		} else if e := door.setDestination(ctx, room.Noun); e != nil {
			err = e
		} else if e := otherRoom.setCompass(ctx, direction, door.Noun); e != nil {
			err = e
		} else {
			_, err = createPrivateDoor(ctx, back, room, otherRoom)
		}
	}
	return
}

// room R is direction from room P.
// the primary goal is to ensure there's door in room P that leads to room R;
// secondarily, try to put a door in room R leading to P.
func connectRoomToRoom(ctx *Context, room, otherRoom jessLink, direction string) (err error) {
	if door, e := createPrivateDoor(ctx, direction, otherRoom, room); e != nil {
		err = e
	} else if len(door) == 0 {
		err = errors.New("room already has a door")
	} else if back, e := reverse(ctx, direction); e != nil {
		err = e
	} else {
		// try to create the reverse door; dont worry if it creates nothing.
		_, err = createPrivateDoor(ctx, back, room, otherRoom)
	}
	return
}

func reverse(ctx *Context, direction string) (string, error) {
	return ctx.GetOpposite(direction)
}

// generate a door in the first room so that going in the specified direction leads into the other room.
// ex. SOUTH from Lhs is Rhs.
// can return the empty string if there was already room door on that side of the room.
func createPrivateDoor(ctx *Context, direction string, room, otherRoom jessLink) (ret string, err error) {
	if res, e := setDirection(ctx, direction, room, otherRoom); e != nil {
		err = e
	} else if res == setDirectionConflict || res == setDirectionDupe {
		err = errors.New("direction already set") // fix? might be to handle dupe in some cases
	} else {
		// create room magic name: room-direction-door
		door := strings.Replace(room.Noun, " ", "-", -1) + "-" + strings.Replace(direction, " ", "-", -1) + "-door"
		e := room.setCompass(ctx, direction, door)
		if res, e := translateError(e); e != nil {
			err = e
		} else if res == setDirectionOkay {
			if e := ctx.AddNounKind(door, Doors); e != nil {
				err = e
			} else if e := ctx.AddNounName(door, door, 0); e != nil {
				err = e
			} else if e := ctx.AddNounTrait(door, Scenery); e != nil {
				err = e
			} else if e := ctx.AddNounTrait(door, Private); e != nil {
				err = e
			} else if e := ctx.AddNounValue(door, DoorDestination, text(otherRoom.Noun, Rooms)); e != nil {
				err = e
			} else if e := room.addDoor(ctx, door); e != nil {
				err = e
			} else {
				ret = door
			}
		}
	}
	return
}
