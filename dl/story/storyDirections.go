package story

import (
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/imp/assert"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

var _ error = &eph.Conflict{}

// backwards compat
type helperNoun struct {
	name    string
	uniform string // fix: things that generate nouns should be be commands dependent on the generated noun names
}

// return the name of a noun based on the name of the current noun
func (h *helperNoun) dependentNoun(name string) helperNoun {
	next := h.name + "-" + name
	return helperNoun{name: next, uniform: next}
}

// can return a noun with an empty name
// fix: things that generate nouns should be be commands dependent on the generated noun names
//
//	its not valid to generate a noun name here.
func makeNoun(k *imp.Importer, name rt.TextEval) (ret helperNoun, err error) {
	if name, e := safe.GetOptionalText(k, name, ""); e != nil {
		err = e
	} else if name := name.String(); len(name) > 0 {
		a := makeArticleName(name)
		if u := lang.Underscore(a.name); len(u) > 0 {
			ret = helperNoun{name: a.name, uniform: u}
		}
	}
	return
}

// Execute - called by the macro runtime during weave.
func (op *MapHeading) Execute(macro rt.Runtime) error {
	return imp.StoryStatement(macro, op)
}

// departing from the current room in a direction
func (op *MapHeading) PostImport(k *imp.Importer) (err error) {
	if room, e := makeNoun(k, op.RoomName); e != nil {
		err = e
	} else if len(room.name) == 0 {
		err = errutil.New("missing room name")
	} else if otherRoom, e := makeNoun(k, op.OtherRoomName); e != nil {
		err = e
	} else if len(otherRoom.name) == 0 {
		err = errutil.New("missing other room name")
	} else if door, e := makeNoun(k, op.DoorName); e != nil {
		err = e // ^ note: door name is optional
	} else {
		// exit this room moving through the (optional) door
		mapDirect(k, room, otherRoom, door, op.Dir)

		// write a fact stating the general direction from one room to the other has been established.
		// ( used to detect conflicts in (the reverse directional) implications of some other statement )
		if e := k.Schedule(assert.PropertyPhase, func(m assert.World, k assert.Assertions) (err error) {
			if dir := lang.Underscore(op.Dir.Str); len(dir) == 0 {
				err = errutil.New("missing map direction")
			} else {
				err = k.AssertDefinition("dir", room.uniform, dir, otherRoom.uniform)
			}
			return
		}); e != nil {
			err = e
		} else {
			// reverse connect
			if op.MapConnection.isTwoWay() {
				// fix? maybe one way to sort out the ephemera phases would be to give/let
				// the the StoryStatements have a Phase() implementation directly?
				// (also maybe the functions should be allowed to be after or before a named phase?)
				if e := k.Schedule(assert.FieldPhase, func(m assert.World, k assert.Assertions) (err error) {
					if dir := lang.Underscore(op.Dir.Str); len(dir) == 0 {
						err = errutil.New("missing map direction")
					} else {
						otherDir := m.OppositeOf(dir)
						// to prioritize some other potentially more explicit definition of a door:
						// if the directional connection is newly established, lets connect these two rooms.
						// it's possible that the only way to handle all the potential conflicts
						// ( ex. an author manually specifying a door and settings its directions )
						// and explicit "assemble directions" phases would be needed.
						var missingDoor helperNoun
						if e := k.AssertDefinition("dir", otherRoom.uniform, otherDir, room.uniform); e == nil {
							// create the reverse door, etc.
							err = mapDirect(k, otherRoom, room, missingDoor, MapDirection{otherDir})
						} else {
							if conflict, ok := e.(*eph.Conflict); !ok || conflict.Redefined() {
								err = e
							}
							// FIX FIX FIX: this compiles with "go build" but not with "go test"?!?!?
							// var conflict eph.Conflict
							// if !errors.As(e, &conflict) || conflict.Redefined() {
							// 	err = e
							// }
						}
					}
					return
				}); e != nil {
					err = e
				}
			}
		}
	}
	return
}

// Execute - called by the macro runtime during weave.
func (op *MapDeparting) Execute(macro rt.Runtime) error {
	return imp.StoryStatement(macro, op)
}

// departing from the current room via a door
func (op *MapDeparting) PostImport(k *imp.Importer) (err error) {
	if room, e := makeNoun(k, op.RoomName); e != nil {
		err = e // ^ todo: ensure the room exists without declaring it
	} else if door, e := makeNoun(k, op.DoorName); e != nil {
		err = e
	} else if otherRoom, e := makeNoun(k, op.OtherRoomName); e != nil {
		err = e
	} else if e := k.AssertNounKind(door.name, "doors"); e != nil {
		err = e // ^ ensure the exit exists
	} else if e := k.AssertRelative("whereabouts", room.name, door.name); e != nil {
		err = e // ^ put the exit in the current room
	} else if e := assert.AssertNounValue(k, Tx(otherRoom.uniform, "rooms"), door.name, "destination"); e != nil {
		err = e // ^ set the door's target to the other room; todo:
	}
	return
}

// set the room's compass, creating an exit if needed to normalize directional travel to always involve a door.
func mapDirect(k assert.Assertions, room, otherRoom, exitDoor helperNoun, mapDir MapDirection) (err error) {
	if dir := lang.Underscore(mapDir.Str); len(dir) == 0 {
		err = errutil.New("missing map direction")
	} else {
		generateExit := len(exitDoor.name) == 0
		if generateExit { // ex. "lobby-up-door"
			exitDoor = room.dependentNoun(dir + "-door")
		}
		//  manually transform the names since we are using them as values
		exitName, otherName := exitDoor.uniform, otherRoom.uniform
		// -- Refs(nounOf(room, "rooms")))// verify the current room
		// -- Refs(nounOf(otherRoom, "rooms")))// verify the target room
		if e := k.AssertNounKind(exitDoor.name, "doors"); e != nil {
			err = e // ^ ensure the existence of the door
		} else if e := k.AssertRelative("whereabouts", room.name, exitDoor.name); e != nil {
			err = e // ^ put the exit in the current room
		} else if e := assert.AssertNounValue(k, Tx(exitName, "door"), room.name, "compass", dir); e != nil {
			err = e // ^ set the room's compass to the exit
		} else if e := assert.AssertNounValue(k, Tx(otherName, "rooms"), exitDoor.name, "destination"); e != nil {
			err = e // ^ set the door's target to the other room
		} else {
			if generateExit {
				// mark the autogenerated door as privately named scenery.
				// ( keeps it unlisted, and stops the player from being able to refer to it )
				if e := assert.AssertNounValue(k, B(true), exitName, "scenery"); e != nil {
					err = e
				} else if e := assert.AssertNounValue(k, B(true), exitName, "privately_named"); e != nil {
					err = e
				}
			}
		}
	}
	return
}

func (op *MapConnection) isTwoWay() bool {
	return op.Str == MapConnection_ConnectingTo
}
