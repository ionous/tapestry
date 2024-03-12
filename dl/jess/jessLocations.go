package jess

import (
	"errors"
	"strings"
)

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

func (op *MapLocations) Generate(rar Registrar) (err error) {
	var links []jessLink
	if e := rar.PostProcess(GenerateNouns, func(q Query) (err error) {
		if room, e := op.Linking.BuildNoun(q, rar, nil, nil); e != nil {
			err = e
		} else {
			links = append(links, makeLink(room, ""))
			for it := op.GetOtherLocations(); it.HasNext(); {
				link := it.GetNext()
				if el, e := link.buildLink(q, rar); e != nil {
					err = e
					break
				} else {
					links = append(links, el)
				}
			}
		}
		return
	}); e != nil {
		err = e
	} else if e := rar.PostProcess(GenerateDefaultKinds, func(q Query) error {
		return generateDefaultKinds(rar, links)
	}); e != nil {
		err = e
	} else {
		err = rar.PostProcess(GenerateConnections, func(Query) error {
			return connectPlaceToPlaces(rar, links[0], links[1:])
		})
	}
	return
}

// connection behavior depends on what's room door and what's room room...
// so the generation of locations and connections needs to be after applying the default kinds
func connectPlaceToPlaces(rar Registrar, src jessLink, dst []jessLink) (err error) {
	if !src.roomLike {
		err = connectDoorToRooms(rar, src, dst)
	} else {
		err = connectRoomToPlaces(rar, src, dst)
	}
	return
}

// assuming the lhs of the phrase was a door, try the rhs as rooms.
func connectDoorToRooms(rar Registrar, door jessLink, places []jessLink) (err error) {
	for i, cnt := 0, len(places); i < cnt && err == nil; i++ {
		if p := places[i]; !p.roomLike {
			err = errors.New("both sides cant be doors")
		} else {
			err = connectDoorToRoom(rar, door, p, p.direction)
		}
	}
	return
}

// assuming the lhs of the phrase was a room, try the rhs as rooms or doors.
func connectRoomToPlaces(rar Registrar, room jessLink, places []jessLink) (err error) {
	for i, cnt := 0, len(places); i < cnt && err == nil; i++ {
		if p := places[i]; !p.roomLike {
			err = connectRoomToDoor(rar, room, p, p.direction)
		} else {
			err = connectRoomToRoom(rar, room, p, p.direction)
		}
	}
	return
}

// ex. door R is direction from room P.
// this is limited to the door being in the room
// ( i dont think there's any other interpretations
// - because directions always lead to doors, and doors are always exists. )
func connectDoorToRoom(rar Registrar, door, room jessLink, direction string) (err error) {
	if e := room.addDoor(rar, door.Noun); e != nil {
		err = e
	} else {
		// doesn't set the room to room direction because we don't have a known destination for the door.
		err = room.setCompass(rar, direction, door.Noun)
	}
	return
}

// room is direction from door:
// there are two different interpretations:
// 1. the door is inside the room on its opposite side.
// 2. the door is in some other room. in that room, moving in the specified direction
// leads to the door, and the door exits into room R.
// inform only handles the first; but it seems either are valid for tapestry.
func connectRoomToDoor(rar Registrar, room, door jessLink, direction string) (err error) {
	if back, e := reverse(rar, direction); e != nil {
		err = e
	} else if parent, e := door.getParent(rar); e != nil {
		err = e
	} else if len(parent) == 0 || parent == room.Noun {
		// put the door in the room; traveling reverse gets us there.
		err = connectDoorToRoom(rar, door, room, back)
	} else {
		otherRoom := makeRoom(parent)
		if res, e := setDirection(rar, direction, otherRoom, room); e != nil {
			err = e
		} else if res == setDirectionConflict || res == setDirectionDupe {
			err = errors.New("direction already set") // fix? might be to handle dupe in some cases
		} else if e := door.setDestination(rar, room.Noun); e != nil {
			err = e
		} else if e := otherRoom.setCompass(rar, direction, door.Noun); e != nil {
			err = e
		} else {
			_, err = createPrivateDoor(rar, back, room, otherRoom)
		}
	}
	return
}

// room R is direction from room P.
// the primary goal is to ensure there's door in room P that leads to room R;
// secondarily, try to put a door in room R leading to P.
func connectRoomToRoom(rar Registrar, room, otherRoom jessLink, direction string) (err error) {
	if door, e := createPrivateDoor(rar, direction, otherRoom, room); e != nil {
		err = e
	} else if len(door) == 0 {
		err = errors.New("room already has a door")
	} else if back, e := reverse(rar, direction); e != nil {
		err = e
	} else {
		// try to create the reverse door; dont worry if it creates nothing.
		_, err = createPrivateDoor(rar, back, room, otherRoom)
	}
	return
}

func reverse(rar Registrar, direction string) (string, error) {
	return rar.GetOpposite(direction)
}

// generate a door in the first room so that going in the specified direction leads into the other room.
// ex. SOUTH from Lhs is Rhs.
// can return the empty string if there was already room door on that side of the room.
func createPrivateDoor(rar Registrar, direction string, room, otherRoom jessLink) (ret string, err error) {
	if res, e := setDirection(rar, direction, room, otherRoom); e != nil {
		err = e
	} else if res == setDirectionConflict || res == setDirectionDupe {
		err = errors.New("direction already set") // fix? might be to handle dupe in some cases
	} else {
		// create room magic name: room-direction-door
		door := strings.Replace(room.Noun, " ", "-", -1) + "-" + strings.Replace(direction, " ", "-", -1) + "-door"
		e := room.setCompass(rar, direction, door)
		if res, e := translateError(e); e != nil {
			err = e
		} else if res == setDirectionOkay {
			if e := rar.AddNounKind(door, Doors); e != nil {
				err = e
			} else if e := rar.AddNounName(door, door, 0); e != nil {
				err = e
			} else if e := rar.AddNounTrait(door, Scenery); e != nil {
				err = e
			} else if e := rar.AddNounTrait(door, Private); e != nil {
				err = e
			} else if e := rar.AddNounValue(door, DoorDestination, text(otherRoom.Noun, Rooms)); e != nil {
				err = e
			} else if e := room.addDoor(rar, door); e != nil {
				err = e
			} else {
				ret = door
			}
		}
	}
	return
}
