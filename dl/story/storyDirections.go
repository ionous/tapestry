package story

import (
	"errors"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"

	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/imp"
)

// backwards compat
type helperNoun struct {
	name    string
	uniform string
}

func (h *helperNoun) NounName() string {
	return h.name
}

func (h *helperNoun) UniformString() string {
	return h.uniform
}

// return the name of a noun based on the name of the current noun
func (h *helperNoun) dependentNoun(name string) helperNoun {
	next := h.name + "-" + name
	return helperNoun{name: next, uniform: next}
}

// can return a noun with an empty name
func makeNoun(k *imp.Importer, name rt.TextEval) (ret helperNoun, err error) {
	if name, e := safe.GetOptionalText(nil, name, ""); e != nil {
		err = e
	} else if name := name.String(); len(name) > 0 {
		a := makeArticleName(name)
		if u, ok := eph.UniformString(a.name); ok {
			ret = helperNoun{name: a.name, uniform: u}
		}
	}
	return
}

// departing from the current room in a direction
func (op *MapHeading) PostImport(k *imp.Importer) (err error) {
	if room, e := makeNoun(k, op.RoomName); e != nil {
		err = e
	} else if len(room.name) == 0 {
		err = eph.InvalidString(room.name)
	} else if otherRoom, e := makeNoun(k, op.OtherRoomName); e != nil {
		err = e
	} else if len(otherRoom.name) == 0 {
		err = eph.InvalidString(otherRoom.name)
	} else if door, e := makeNoun(k, op.DoorName); e != nil {
		err = e // ^ note: door name is optional
	} else {
		// exit this room moving through the (optional) door
		_ = mapDirect((*storyAdapter)(k), room, otherRoom, door, op.Dir)

		// write a fact stating the general direction from one room to the other has been established.
		// ( used to detect conflicts in (the reverse directional) implications of some other statement )
		k.WriteEphemera(eph.PhaseFunction{eph.PropertyPhase,
			func(c *eph.Catalog, d *eph.Domain, at string) (err error) {
				if dir, e := op.Dir.UniformString(); e != nil {
					err = e
				} else {
					err = addDirection(d, room.uniform, dir, otherRoom.uniform, at)
				}
				return
			}})

		// reverse connect
		if op.MapConnection.isTwoWay() {
			// fix? maybe one way to sort out the ephemera phases would be to give/let
			// the the StoryStatements have a Phase() implementation directly?
			// (also maybe the functions should be allowed to be after or before a named phase?)
			k.WriteEphemera(eph.PhaseFunction{eph.FieldPhase,
				func(c *eph.Catalog, d *eph.Domain, at string) (err error) {
					if dir, e := op.Dir.UniformString(); e != nil {
						err = e
					} else if otherDir, e := d.FindOpposite(dir); e != nil {
						err = e
					} else {
						// to prioritize some other potentially more explicit definition of a door:
						// if the directional connection is newly established, lets connect these two rooms.
						// it's possible that the only way to handle all the potential conflicts
						// ( ex. an author manually specifying a door and settings its directions )
						// and explicit "assemble directions" phases would be needed.
						var missingDoor helperNoun
						if e := addDirection(d, otherRoom.uniform, otherDir, room.uniform, at); e == nil {
							// create the reverse door, etc.
							err = mapDirect(domainAdapter{d, at}, otherRoom, room, missingDoor, MapDirection{otherDir})
						} else {
							var conflict *eph.Conflict
							if !errors.As(e, &conflict) || conflict.Reason == eph.Redefined {
								err = e
							}
						}
					}
					return
				},
			})
		}
	}
	return
}

func (n MapDirection) UniformString() (ret string, err error) {
	if u, ok := eph.UniformString(n.Str); !ok {
		err = eph.InvalidString(n.Str)
	} else {
		ret = u
	}
	return
}

func addDirection(d *eph.Domain, fromRoom, inDir, toRoom, at string) error {
	// fix: a "fact" ephemera that runs immediately whenever its added? ( phase 0 or something maybe )
	return d.AddDefinition(eph.MakeKey("dir", fromRoom, inDir), at, toRoom)
}

// departing from the current room via a door
func (op *MapDeparting) PostImport(k *imp.Importer) (err error) {
	if room, e := makeNoun(k, op.RoomName); e != nil {
		err = e
	} else if door, e := makeNoun(k, op.DoorName); e != nil {
		err = e
	} else {
		k.WriteEphemera(Refs(nounOf(room, "rooms")))                        // verify the current room
		k.WriteEphemera(nounOf(door, "doors"))                              // ensure the exit
		k.WriteEphemera(relateTo(room, "whereabouts", door))                // put the exit in the current room
		valForField(door, Tx(door.UniformString(), "rooms"), "destination") // set the door's target to the other room
	}
	return
}

type ephemeraWriter interface {
	WriteEphemera(eph.Ephemera) error
}

type domainAdapter struct {
	d  *eph.Domain
	at string
}

func (da domainAdapter) WriteEphemera(op eph.Ephemera) error {
	return da.d.AddEphemera(da.at, op)
}

type storyAdapter imp.Importer

func (sa *storyAdapter) WriteEphemera(op eph.Ephemera) (noerr error) {
	k := (*imp.Importer)(sa)
	k.WriteEphemera(op)
	return
}

// set the room's compass, creating an exit if needed to normalize directional travel to always involve a door.
func mapDirect(k ephemeraWriter, room, otherRoom, exitDoor helperNoun, mapDir MapDirection) (err error) {
	if dir, ok := eph.UniformString(mapDir.Str); !ok {
		err = eph.InvalidString(mapDir.Str)
	} else {
		generateExit := len(exitDoor.name) == 0
		if generateExit { // ex. "lobby-up-door"
			exitDoor = room.dependentNoun(dir + "-door")
		}
		//  manually transform the names since we are using them as values
		exitName, otherName := exitDoor.UniformString(), otherRoom.UniformString()
		if e := k.WriteEphemera(Refs(nounOf(room, "rooms"))); e != nil {
			err = e // ^ verify the current room
		} else if e := k.WriteEphemera(Refs(nounOf(otherRoom, "rooms"))); e != nil {
			err = e // ^ verify the target room
		} else if e := k.WriteEphemera(nounOf(exitDoor, "doors")); e != nil {
			err = e // ^ ensure the existence of the door
		} else if k.WriteEphemera(relateTo(room, "whereabouts", exitDoor)); e != nil {
			err = e // ^ put the exit in the current room
		} else if e := k.WriteEphemera(valForField(room, Tx(exitName, "door"), "compass", dir)); e != nil {
			err = e // ^ set the room's compass to the exit
		} else if e := k.WriteEphemera(valForField(exitDoor, Tx(otherName, "rooms"), "destination")); e != nil {
			err = e // ^ set the door's target to the other room
		} else if generateExit {
			// mark the autogenerated door as privately named scenery.
			// ( keeps it unlisted, and stops the player from being able to refer to it )
			k.WriteEphemera(&eph.EphValues{Noun: exitName, Field: "scenery", Value: B(true)})
			k.WriteEphemera(&eph.EphValues{Noun: exitName, Field: "privately_named", Value: B(true)})
		}
	}
	return
}

func (op *MapConnection) isTwoWay() bool {
	return op.Str == MapConnection_ConnectingTo
}

func nounOf(n helperNoun, kind string) *eph.EphNouns {
	return &eph.EphNouns{
		Noun: n.NounName(),
		Kind: kind,
	}
}

func relateTo(n helperNoun, rel string, otherNoun helperNoun) *eph.EphRelatives {
	return &eph.EphRelatives{
		Rel:       rel,
		Noun:      n.NounName(),
		OtherNoun: otherNoun.NounName()}
}

// give this noun the passed value at the named field and path
func valForField(n helperNoun, v literal.LiteralValue, path ...string) *eph.EphValues {
	last := len(path) - 1
	field, parts := path[last], path[:last]
	return &eph.EphValues{
		Noun:  n.NounName(),
		Field: field,
		Path:  parts,
		Value: v,
	}
}
