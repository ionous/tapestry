package story

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/dl/literal"
)

// departing from the current room in a direction
func (op *MapHeading) ImportPhrase(k *Importer) (noerr error) {
	// exit this room moving through the (optional) door
	_ = mapDirect((*storyAdapter)(k), op.Room, op.OtherRoom, op.Door, op.Dir)

	// write a fact stating the general direction from one room to the other has been established.
	// ( used to detect conflicts in (the reverse directional) implications of some other statement )
	k.WriteEphemera(eph.PhaseFunction{eph.PropertyPhase,
		func(c *eph.Catalog, d *eph.Domain, at string) (err error) {
			if room, e := op.Room.UniformString(); e != nil {
				err = e
			} else if dir, e := op.Dir.UniformString(); e != nil {
				err = e
			} else if toRoom, e := op.OtherRoom.UniformString(); e != nil {
				err = e
			} else {
				err = addDirection(d, room, dir, toRoom, at)
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
				} else if firstRoom, e := op.Room.UniformString(); e != nil {
					err = e
				} else if otherRoom, e := op.OtherRoom.UniformString(); e != nil {
					err = e
				} else {
					// to prioritize some other potentially more explicit definition of a door:
					// if the directional connection is newly established, lets connect these two rooms.
					// it's possible that the only way to handle all the potential conflicts
					// ( ex. an author manually specifying a door and settings its directions )
					// and explicit "assemble directions" phases would be needed.
					if e := addDirection(d, otherRoom, otherDir, firstRoom, at); e == nil {
						// create the reverse door, etc.
						err = mapDirect(domainAdapter{d, at}, op.OtherRoom, op.Room, nil, MapDirection{otherDir})
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
func (op *MapDeparting) ImportPhrase(k *Importer) (err error) {
	if exitName, e := op.Door.UniformString(); e != nil {
		err = e // ^ manually transform the names since we are using them as values
	} else {
		k.WriteEphemera(Refs(nounOf(op.Room, "rooms")))            // verify the current room
		k.WriteEphemera(nounOf(op.Door, "doors"))                  // ensure the exit
		k.WriteEphemera(relateTo(op.Room, "whereabouts", op.Door)) // put the exit in the current room
		valForField(op.Door, Tx(exitName, "rooms"), "destination") // set the door's target to the other room
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

type storyAdapter Importer

func (sa *storyAdapter) WriteEphemera(op eph.Ephemera) (noerr error) {
	k := (*Importer)(sa)
	k.WriteEphemera(op)
	return
}

// set the room's compass, creating an exit if needed to normalize directional travel to always involve a door.
func mapDirect(k ephemeraWriter, room, otherRoom SingularNoun, optionalExit SingularNoun, mapDir MapDirection) (err error) {
	if dir, ok := eph.UniformString(mapDir.Str); !ok {
		err = eph.InvalidString(mapDir.Str)
	} else {
		var exitDoor SingularNoun
		if optionalExit != nil {
			exitDoor = optionalExit
		} else {
			// ex. "lobby-up-door"
			exitDoor = dependentNoun(room, dir+"-door")
		}
		if exitName, e := exitDoor.UniformString(); e != nil {
			err = e // ^ manually transform the names since we are using them as values
		} else if otherName, e := otherRoom.UniformString(); e != nil {
			err = e // ^ manually transform the names since we are using them as values
		} else if e := k.WriteEphemera(Refs(nounOf(room, "rooms"))); e != nil {
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
		} else if optionalExit == nil {
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

// return the name of a noun based on the name of the current noun
func dependentNoun(n SingularNoun, name string) SingularNoun {
	return &CommonNoun{
		Determiner: Determiner{"our"},
		Noun:       NounNamed{Str: n.NounName() + "-" + name},
	}
}

func nounOf(n SingularNoun, kind string) *eph.EphNouns {
	return &eph.EphNouns{
		Noun: n.NounName(),
		Kind: kind,
	}
}

func relateTo(n SingularNoun, rel string, otherNoun SingularNoun) *eph.EphRelatives {
	return &eph.EphRelatives{
		Rel:       rel,
		Noun:      n.NounName(),
		OtherNoun: otherNoun.NounName()}
}

// give this noun the passed value at the named field and path
func valForField(n SingularNoun, v literal.LiteralValue, path ...string) *eph.EphValues {
	last := len(path) - 1
	field, parts := path[last], path[:last]
	return &eph.EphValues{
		Noun:  n.NounName(),
		Field: field,
		Path:  parts,
		Value: v,
	}
}
