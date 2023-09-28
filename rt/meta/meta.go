package meta

const Prefix = '$' // leading character used for all internal targets

// targets for GetField
// names are ideally one word ( underscore separated if not ) and have divergent prefixes ( different first letters. )
const (
	Counter        = "$counter"   // sequence counter
	Domain         = "$scene"     // returns whether a named domain is active
	FieldsOfKind   = "$fields"    // names of the fields of a kind as a text list
	ObjectAliases  = "$alias"     // similar to object name but returns a list of names
	ObjectId       = "$object"    // returns the unique object id from a object name
	ObjectKind     = "$kind"      // type of a game object
	ObjectKinds    = "$ancestry"  // ancestor of an object's type ( a text list, root at the start )
	ObjectName     = "$name"      // given a noun, return the friendly name declared by the author
	ObjectsOfKind  = "$instances" // all objects of a given kind
	Option         = "$flags"     // get/set various runtime options ( see below )
	PatternLabels  = "$labels"    // strings of the parameter names, the last is the return
	PatternRunning = "$pattern"   // returns true if the named pattern is running
	Response       = "$response"  // returns replacements for named templates
	Variables      = "$variable"  // named values, controlled by scope, not associated with any particular object
	ValueChanged   = "$dirty"     // indicates the contents inside a value changed
)

// fields for runtime meta.Option(s)
// options are initialized at runtime startup
// new options *cannot* be added dynamically.
type Options int

//go:generate stringer -type=Options
const (
	// flag to print response names ( instead of values )
	PrintResponseNames Options = iota
	NumOptions
)
