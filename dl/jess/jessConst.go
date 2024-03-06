package jess

// names ( of properties, kinds, etc. ) that jess expects from the tapestry standard library
const (
	// fields:
	IndefiniteArticle = "indefinite article"
	PrintedName       = "printed name"
	DoorDestination   = "destination"
	// traits
	PluralNamedTrait = "plural named"
	ProperNameTrait  = "proper named"
	CountedTrait     = "counted"
	// special names:
	PlayerSelf = "self"
	// kinds:
	Objects    = "objects"    // the kind of a potential room or door
	Directions = "directions" // a special kind of object representing movement of travel
	Rooms      = "rooms"      // possible player locale
	Things     = "things"
	Doors      = "doors" // portals which connect rooms
)
