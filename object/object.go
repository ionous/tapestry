package object

const Prefix = '$' // leading character used for all internal targets

// internal targets for GetField
const Aspect = "$aspect"   // name of aspect for noun.trait
const Counter = "$counter" // sequence counter
const Domain = "$domain"   // returns whether a named domain is active
const Id = "$id"           // returns the unique object id from a object name
const Nouns = "$nouns"     // returns the list of active nouns for a given kind
const Running = "$running" // returns true if the named pattern is running
const Option = "$opt"
const Value = "$value"   // returns the object g.Value
const Variables = "$var" // named values, controlled by scope, not associated with any particular object

// internal fields for object
const Active = "$active" // is the noun in a valid domain
const Name = "$name"     // name of an object as declared by the user
const Kind = "$kind"     // type of a game object
const Kinds = "$kinds"   // hierarchy of an object's types ( a path )

// fields for options
type Options int

//go:generate stringer -type=Options
const (
	// a true/false flag to print response names ( instead of values )
	PrintResponseNames Options = iota + 1
)
