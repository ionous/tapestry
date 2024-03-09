package jess

import (
	"errors"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
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
	var post []postConnect
	if e := rar.PostProcess(GenerateNouns, func(q Query) (err error) {
		if a, e := op.Linking.BuildNoun(q, rar, nil, nil); e != nil {
			err = e
		} else {
			post = append(post, makePostConnect(a, ""))
			for it := op.GetOtherLocations(); it.HasNext(); {
				link := it.GetNext()
				if b, e := link.BuildNoun(q, rar, nil, []string{""}); e != nil {
					err = e
					break
				} else {
					// direction is already normalized...
					post = append(post, makePostConnect(b, link.Direction.Text))
				}
			}
		}
		return
	}); e != nil {
		err = e
	} else if e := rar.PostProcess(GenerateDefaultKinds, func(q Query) error {
		return generateDefaultKinds(rar, post)
	}); e != nil {
		err = e
	} else {
		// connection behavior depends on what's a door and what's a room...
		// so the generation of locations and connections needs to be after applying the default kinds
		err = rar.PostProcess(GenerateConnections, func(q Query) (err error) {
			src, rest := post[0], post[1:]
			if !src.roomLike {
				err = src.genDoor(rar, rest)
			} else {
				err = src.genRoom(rar, rest)
			}
			return
		})
	}
	return
}

type postConnect struct {
	*DesiredNoun
	roomLike  bool // valid after generating default kinds
	direction string
}

func makePostConnect(n *DesiredNoun, direction string) postConnect {
	return postConnect{DesiredNoun: n, direction: direction}
}

func (p *postConnect) generateDefaultKind(rar Registrar) (err error) {
	noun := p.Noun
	// either newly stamping it as a room, or duplicating a previous room definition is okay.
	if e := rar.AddNounKind(noun, Rooms); e == nil || errors.Is(e, mdl.Duplicate) {
		p.roomLike = true
	} else {
		// some unknown error is a problem:
		if !errors.Is(e, mdl.Conflict) {
			err = e
		} else {
			// oto, if it was conflicted, maybe it was actually a door.
			if e := rar.AddNounKind(noun, Doors); e != nil && !errors.Is(e, mdl.Duplicate) {
				err = e
			}
		}
	}
	return
}

func (a *postConnect) genDoor(rar Registrar, ps []postConnect) (err error) {
	for i, cnt := 0, len(ps); i < cnt && err == nil; i++ {
		if b := ps[i]; !b.roomLike {
			err = errors.New("both sides cant be doors")
		} else {
			err = a.connectDoorToRoom(rar, b)
		}
	}
	return
}

func (a *postConnect) genRoom(rar Registrar, ps []postConnect) (err error) {
	for i, cnt := 0, len(ps); i < cnt && err == nil; i++ {
		if b := ps[i]; !b.roomLike {
			err = a.connectRoomToDoor(rar, b)
		} else {
			err = a.connectRoomToRoom(rar, b)
		}
	}
	return
}

// ex. door A is direction from room B.
// aka: direction from B is door A.
func (a *postConnect) connectDoorToRoom(rar Registrar, b postConnect) (err error) {
	room, door := b.Noun, a.Noun
	if e := addDoor(rar, room, door); e != nil {
		err = e
	} else {
		// ? setDirection(rar, b.Noun, b.direction, a.Noun) ?
		err = setCompass(rar, room, b.direction, door)
	}
	return
}

// room A is direction through door B.
func (a *postConnect) connectRoomToDoor(rar Registrar, b postConnect) (err error) {
	// 	if B were a door (ex. in some other room O), we'd want something like:
	// * `B.destination = A`
	// * `O.compass[direction] = B`
	// * `fact: 'dir -> <room>.dir -> B`
	//
	// we can also manufacture a private door in A that leads into O in the reverse direction;
	// guarding against the case where A already has a door on the reverse side.
	// to find room O, jess needs to be able to ask about B's whereabouts...
	// after the explicit phrases have been played out. ( GenerateConnections )
	panic("not implemented")
}

// room A is direction from room B.
func (a *postConnect) connectRoomToRoom(rar Registrar, b postConnect) (err error) {
	if reverse, e := rar.GetOpposite(b.direction); e != nil {
		err = e
	} else if doorb, e := createPrivateDoor(rar, b.Noun, b.direction, a.Noun); e != nil {
		err = e
	} else if len(doorb) == 0 {
		err = errors.New("room already has a door")
	} else {
		// dont worry if we couldn't create the door
		// fix: this code path needs some more love
		_, err = createPrivateDoor(rar, a.Noun, reverse, b.Noun)
	}
	return
}

// generate a door on the given side of the room
// can return the empty string if there was already a door on that side of the room.
func createPrivateDoor(rar Registrar, room, direction, otherRoom string) (ret string, err error) {
	if res, e := setDirection(rar, room, direction, otherRoom); e != nil {
		err = e
	} else if res == setDirectionConflict || res == setDirectionDupe {
		err = errors.New("direction already set") // fix? might be to handle dupe in some cases
	} else {
		// create a magic name: room-direction-door
		door := strings.Replace(room, " ", "-", -1) + "-" + strings.Replace(direction, " ", "-", -1) + "-door"
		e := setCompass(rar, room, direction, door)
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
			} else if e := setDestination(rar, door, otherRoom); e != nil {
				err = e
			} else if e := addDoor(rar, room, door); e != nil {
				err = e
			} else {
				ret = door
			}
		}
	}
	return
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

// assumes a is "room like"
func addDoor(rar Registrar, room, door string) error {
	return rar.AddNounPair(Whereabouts, room, door)
}

// FIX: do we need this!?
// create a fact which indicates the direction of movement from room to room
// these facts help with tracking and conflict detection
func setDirection(rar Registrar, room, direction, otherRoom string) (int, error) {
	e := rar.AddFact(FactDirection, room, direction, otherRoom)
	return translateError(e)
}

// set the compass on the indicated side of the room to the named door
func setCompass(rar Registrar, room, direction, door string) error {
	return rar.AddNounPath(room,
		[]string{Compass, direction},
		&literal.TextValue{Value: door, Kind: Doors},
	)
}

// set the destination of the named door
func setDestination(rar Registrar, door, otherRoom string) (err error) {
	return rar.AddNounValue(door, DoorDestination, text(otherRoom, Rooms))
}

// -
func generateDefaultKinds(rar Registrar, ps []postConnect) (err error) {
	for i, cnt := 0, len(ps); i < cnt; i++ {
		// use indexing so generateDefaultKind can properly work on the shared memory
		// range would be a copy
		if e := ps[i].generateDefaultKind(rar); e != nil {
			err = e
			break
		}
	}
	return
}
