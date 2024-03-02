package jess

// names ( of properties, kinds, etc. ) that jess expects from the tapestry standard library
const (
	// fields:
	IndefiniteArticle = "indefinite article"
	PrintedName       = "printed name"
	DoorDestination   = "destination"
	// traits
	PluralNamed     = "plural named"
	ProperNameTrait = "proper named"
	CountedTrait    = "counted"
	// special names:
	PlayerSelf = "self"
	// kinds:
	Things     = "things"
	Objects    = "objects"    // the kind of a potential room or door
	Rooms      = "rooms"      // possible player locale
	Doors      = "doors"      // portals which connect rooms
	Directions = "directions" // a special kind of object representing movement of travel
)
