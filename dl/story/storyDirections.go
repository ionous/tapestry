package story

import (
	"errors"
	"strings"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/ionous/errutil"
)

func (op *MapHeading) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// departing from the current room in a direction
func (op *MapHeading) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequirePlurals, func(w *weave.Weaver) (err error) {
		//
		if room, e := safe.GetText(w, op.RoomName); e != nil {
			err = e
		} else if otherRoom, e := safe.GetText(w, op.OtherRoomName); e != nil {
			err = e
		} else if door, e := safe.GetOptionalText(w, op.DoorName, ""); e != nil {
			err = e // ^ note: door is optional; if not specified door names are blank
		} else {
			err = findClosestNouns(cat, func(w *weave.Weaver, nouns []string) (err error) {
				room, otherRoom, door := nouns[0], nouns[1], nouns[2]
				if e := cat.Schedule(weave.RequireAncestry, func(w *weave.Weaver) error {
					// exit this room moving through the (optional) door
					return mapDirect(w, room, otherRoom, door, op.Dir)
				}); e != nil {
					err = e
				} else {
					pen := w.Pin()
					// write a fact stating the general direction from one room to the other has been established.
					// ( used to detect conflicts in (the reverse directional) implications of some other statement )
					if dir := inflect.Normalize(op.Dir); len(dir) == 0 {
						err = errutil.New("empty map direction")
					} else if ok, e := pen.AddFact("dir", room, dir, otherRoom); e != nil {
						err = e
					} else if ok && op.MapConnection.isTwoWay() {
						if dir := inflect.Normalize(op.Dir); len(dir) == 0 {
							err = errutil.New("empty map direction")
						} else {
							otherDir := w.OppositeOf(dir)
							// to prioritize some other potentially more explicit definition of a door:
							// if the directional connection is newly established, lets connect these two rooms.
							// it's possible that the only way to handle all the potential conflicts
							// ( ex. an author manually specifying a door and settings its directions )
							// and explicit "assemble directions" phases would be needed.
							if ok, e := pen.AddFact("dir", otherRoom, otherDir, room); e != nil {
								err = e
							} else if ok {
								err = cat.Schedule(weave.RequireAncestry, func(w *weave.Weaver) error {
									// create the reverse door, etc.
									return mapDirect(w, otherRoom, room, "", otherDir)
								})
							}
						}
					}
				}
				return
			}, room.String(), otherRoom.String(), door.String())
		}
		return
	})
}

type afterClosestNouns func(w *weave.Weaver, nouns []string) error

// a low brow promise for translating author noun references into full noun names
type closestNouns struct {
	names []string
	nouns []string // starts empty and grows to match name[]
	next  afterClosestNouns
}

// searches for a list of existing nouns
func findClosestNouns(cat *weave.Catalog, next afterClosestNouns, names ...string) error {
	c := closestNouns{names: names, nouns: make([]string, 0, len(names)), next: next}
	return cat.Schedule(weave.RequireNouns, c.schedule)
}

// matches weave.ScheduleCallback
func (c *closestNouns) schedule(w *weave.Weaver) (err error) {
	for i := len(c.nouns); i < len(c.names); i++ {
		var noun string
		if n := c.names[i]; len(n) > 0 {
			if noun, err = w.GetClosestNoun(n); err != nil {
				break
			}
		}
		c.nouns = append(c.nouns, noun)
	}
	// all done?
	if err == nil {
		err = c.next(w, c.nouns)
	}
	return
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
		if room, e := safe.GetText(w, op.RoomName); e != nil {
			err = e
		} else if otherRoom, e := safe.GetText(w, op.OtherRoomName); e != nil {
			err = e
		} else if door, e := safe.GetText(w, op.DoorName); e != nil {
			err = e
		} else {
			err = findClosestNouns(cat, func(w *weave.Weaver, nouns []string) (err error) {
				room, otherRoom, door := nouns[0], nouns[1], nouns[2]
				pen := w.Pin()
				if e := relateNouns(w, room, door); e != nil {
					err = e // ^ put the exit in the current room
				} else if e := pen.AddNounKind(door, "doors"); e != nil && !errors.Is(e, mdl.Duplicate) {
					err = e // ^ ensure the exit exists as a door.
				} else if e := w.AddNounValue(pen, door, "destination", text(otherRoom, "rooms")); e != nil {
					err = e // ^ set the door's target to the other room; todo:
				}
				return
			}, room.String(), otherRoom.String(), door.String())
		}
		return
	})
}

// set the room's compass, creating an exit if needed to normalize directional travel to always involve a door.
func mapDirect(w *weave.Weaver, room, otherRoom, exitDoor string, mapDir string) (err error) {
	if dir := inflect.Normalize(mapDir); len(dir) == 0 {
		err = errutil.New("empty map direction")
	} else {
		pen := w.Pin()
		generateExit := len(exitDoor) == 0
		if generateExit { // ex. "lobby-up-door"
			exitDoor = strings.Replace(room, " ", "-", -1) + "-" +
				strings.Replace(dir, " ", "-", -1) + "-door"
		}
		if e := relateNouns(w, room, exitDoor); e != nil {
			err = e // ^ put the exit in the current room
		} else if e := pen.AddNounKind(exitDoor, "doors"); e != nil && !errors.Is(e, mdl.Duplicate) {
			err = e // ^ ensure the existence of the door
		} else {
			if !generateExit {
				err = mdl.AddNounNames(pen, exitDoor, mdl.MakeNames(exitDoor))
			} else {
				// mark the autogenerated door as privately named scenery.
				// ( keeps it unlisted, and stops the player from being able to refer to it )
				if e := pen.AddNounName(exitDoor, exitDoor, 0); e != nil {
					err = e
				} else if e := w.AddNounValue(pen, exitDoor, "scenery", truly()); e != nil {
					err = e
				} else if e := w.AddNounValue(pen, exitDoor, "privately named", truly()); e != nil {
					err = e
				}
			}
			if err == nil {
				if e := pen.AddNounPath(room, []string{"compass", dir}, Tx(exitDoor, "door")); e != nil {
					err = errutil.Fmt("%w going %s in %q via %q", e, dir, room, exitDoor) // ^ set the room's compass to the exit
				} else if e := pen.AddNounPath(exitDoor, []string{"destination"}, Tx(otherRoom, "rooms")); e != nil {
					err = errutil.Fmt("%w arriving at %q via %q", e, otherRoom, exitDoor) // ^ set the room's compass to the exit
				}
			}
		}
	}
	return
}

func (op MapConnection) isTwoWay() bool {
	return op == C_MapConnection_ConnectingTo
}

// queue this as its own commands helps ensure the relation gets built properly
func relateNouns(w *weave.Weaver, noun, other string) error {
	return w.Catalog.Schedule(weave.RequireNames, func(w *weave.Weaver) error {
		return w.Pin().AddNounPair("whereabouts", noun, other)
	})
}
