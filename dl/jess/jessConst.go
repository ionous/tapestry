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
	Things     = "things"     // the default for named objects if nothing else is specified
	Doors      = "doors"      // portals which connect rooms
	Verbs      = "verbs"      // nouns are used to describe verbs
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
	DirectionOpposite = "opposite"
	// relations:
	Whereabouts = "whereabouts"
	// verbs:
	VerbSubject     = "subject"
	VerbAlternate   = "alternate subject"
	VerbObject      = "object"
	VerbRelation    = "relation"
	VerbImplication = "implication"
	VerbReversed    = "reversed status"
	ReversedTrait   = "reversed"
)
