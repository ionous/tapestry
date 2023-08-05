package story

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/ionous/errutil"
)

// backwards compat
type helperNoun struct {
	name    string
	uniform string
}

// return the name of a noun based on the name of the current noun
func (h *helperNoun) dependentNoun(name string) helperNoun {
	next := h.uniform + "-" + name
	return helperNoun{name: next, uniform: next}
}

// can return a noun with an empty name
func normalizeName(w *weave.Weaver, name rt.TextEval) (ret helperNoun, err error) {
	if name, e := safe.GetOptionalText(w, name, ""); e != nil {
		err = e
	} else if name := name.String(); len(name) > 0 {
		if name, e := grok.StripArticle(name); e != nil {
			err = e
		} else if u := lang.Normalize(name); len(u) > 0 {
			ret = helperNoun{name: name, uniform: u}
		}
	}
	return
}

// Execute - called by the macro runtime during weave.
func (op *MapHeading) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// departing from the current room in a direction
func (op *MapHeading) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequirePlurals, func(w *weave.Weaver) (err error) {
		//
		if room, e := normalizeName(w, op.RoomName); e != nil {
			err = e
		} else if len(room.name) == 0 {
			err = errutil.New("empty room name")
		} else if otherRoom, e := normalizeName(w, op.OtherRoomName); e != nil {
			err = e
		} else if len(otherRoom.name) == 0 {
			err = errutil.New("empty other room name")
		} else if door, e := normalizeName(w, op.DoorName); e != nil {
			err = e // ^ note: door is optional; if not specified door names are blank
		} else if e := cat.Schedule(weave.RequireAncestry, func(w *weave.Weaver) error {
			// exit this room moving through the (optional) door
			return mapDirect(w, room, otherRoom, door, op.Dir)
		}); e != nil {
			err = e
		} else {
			pen := w.Pin()
			// write a fact stating the general direction from one room to the other has been established.
			// ( used to detect conflicts in (the reverse directional) implications of some other statement )
			if dir := lang.Normalize(op.Dir.Str); len(dir) == 0 {
				err = errutil.New("empty map direction")
			} else if e := pen.AddFact(makeKey("dir", room.uniform, dir), otherRoom.uniform); e != nil {
				err = e
			} else {
				// reverse connect
				if op.MapConnection.isTwoWay() {
					if dir := lang.Normalize(op.Dir.Str); len(dir) == 0 {
						err = errutil.New("empty map direction")
					} else {
						otherDir := w.OppositeOf(dir)
						// to prioritize some other potentially more explicit definition of a door:
						// if the directional connection is newly established, lets connect these two rooms.
						// it's possible that the only way to handle all the potential conflicts
						// ( ex. an author manually specifying a door and settings its directions )
						// and explicit "assemble directions" phases would be needed.
						var missingDoor helperNoun
						if e := pen.AddFact(makeKey("dir", otherRoom.uniform, otherDir), room.uniform); e != nil {
							err = e
						} else {
							err = cat.Schedule(weave.RequireAncestry, func(w *weave.Weaver) error {
								// create the reverse door, etc.
								return mapDirect(w, otherRoom, room, missingDoor, MapDirection{otherDir})
							})
						}
					}
				}
			}
		}
		return
	})
}

func makeKey(path ...string) (ret string) {
	return strings.Join(path, "/")
}

// Execute - called by the macro runtime during weave.
func (op *MapDeparting) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// departing from the current room via a door
func (op *MapDeparting) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequirePlurals, func(w *weave.Weaver) (err error) {
		if room, e := normalizeName(w, op.RoomName); e != nil {
			err = e // ^ todo: ensure the room exists without declaring it
		} else if door, e := normalizeName(w, op.DoorName); e != nil {
			err = e
		} else if otherRoom, e := normalizeName(w, op.OtherRoomName); e != nil {
			err = e
		} else {
			pen := w.Pin()
			if e := relateNouns(w, room, door); e != nil {
				err = e // ^ put the exit in the current room
			} else if e := pen.AddNoun(door.uniform, door.name, "doors"); e != nil {
				err = e // ^ ensure the exit exists
			} else if e := pen.AddFieldValue(door.uniform, "destination", text(otherRoom.uniform, "rooms")); e != nil {
				err = e // ^ set the door's target to the other room; todo:
			}
		}
		return
	})
}

// set the room's compass, creating an exit if needed to normalize directional travel to always involve a door.
func mapDirect(w *weave.Weaver, room, otherRoom, exitDoor helperNoun, mapDir MapDirection) (err error) {
	if dir := lang.Normalize(mapDir.Str); len(dir) == 0 {
		err = errutil.New("empty map direction")
	} else {
		pen := w.Pin()
		generateExit := len(exitDoor.name) == 0
		if generateExit { // ex. "lobby-up-door"
			exitDoor = room.dependentNoun(dir + "-door")
		}
		// -- Refs(nounOf(room, "rooms")))// verify the current room
		// -- Refs(nounOf(otherRoom, "rooms")))// verify the target room
		if e := relateNouns(w, room, exitDoor); e != nil {
			err = e // ^ put the exit in the current room
		} else if e := pen.AddNoun(exitDoor.uniform, exitDoor.name, "doors"); e != nil {
			err = e // ^ ensure the existence of the door
		} else {
			if generateExit {
				// mark the autogenerated door as privately named scenery.
				// ( keeps it unlisted, and stops the player from being able to refer to it )
				if e := pen.AddFieldValue(exitDoor.uniform, "scenery", truly()); e != nil {
					err = e
				} else if e := pen.AddFieldValue(exitDoor.uniform, "privately named", truly()); e != nil {
					err = e
				}
			}
			if err == nil {
				if e := pen.AddPathValue(room.uniform, mdl.MakePath("compass", dir), Tx(exitDoor.uniform, "door")); e != nil {
					err = e // ^ set the room's compass to the exit
				} else if e := pen.AddPathValue(exitDoor.uniform, "destination", Tx(otherRoom.uniform, "rooms")); e != nil {
					err = e // ^ set the door's target to the other room
				}
			}
		}
	}
	return
}

func (op *MapConnection) isTwoWay() bool {
	return op.Str == MapConnection_ConnectingTo
}

// queue this as its own commands helps ensure the relation gets built properly
func relateNouns(w *weave.Weaver, noun, other helperNoun) error {
	return w.Catalog.Schedule(weave.RequireNames, func(w *weave.Weaver) error {
		return w.Pin().AddPair("whereabouts", noun.uniform, other.uniform)
	})
}
