package jess

// names ( of properties, kinds, etc. ) that jess expects from the tapestry standard library
const (
	// special names:
	PlayerSelf    = "self"
	FactDirection = "dir"
	// kinds:
	Objects    = "objects"    // the kind of a potential room or door
	Directions = "directions" // a special kind of object representing movement of travel
	Rooms      = "rooms"      // possible player locale
	Things     = "things"
	Doors      = "doors" // portals which connect rooms
	// traits
	CountedTrait     = "counted"
	PluralNamedTrait = "plural named"
	ProperNameTrait  = "proper named"
	// fields:
	Compass           = "compass"
	DoorDestination   = "destination"
	IndefiniteArticle = "indefinite article"
	PrintedName       = "printed name"
	Private           = "privately named"
	Scenery           = "scenery"
	// relations:
	Whereabouts = "whereabouts"
)
