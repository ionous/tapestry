package meta

const Prefix = '$' // leading character used for all internal targets

// targets for GetField
const Counter = "$count"      // sequence counter
const Domain = "$domain"      // returns whether a named domain is active
const ObjectId = "$id"        // returns the unique object id from a object name
const ObjectKind = "$kind"    // type of a game object
const ObjectKinds = "$kinds"  // ancestor of an object's type ( a text list, root at the start )
const ObjectName = "$name"    // given a noun, return the name declared by the author
const ObjectValue = "$obj"    // returns the object g.Value
const ObjectsOfKind = "$objs" // all objects of a given kind
const PatternRunning = "$run" // returns true if the named pattern is running
const Option = "$opt"         // get/set various runtime options ( see below )
const Variables = "$var"      // named values, controlled by scope, not associated with any particular object

// fields for runtime meta.Option(s)
// options are initialized at runtime startup
// new options *cannot* be added dynamically.
type Options int

//go:generate stringer -type=Options
const (
	// flag to print response names ( instead of values )
	PrintResponseNames Options = iota
	// flag to output text as json ( instead of plain text )
	JsonMode
	NumOptions
)
