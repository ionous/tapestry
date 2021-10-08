package jsn

// Marshalee for types which fit into slots
type Marshalee interface {
	Marshal(Marshaler)
}

// Marshaler reads and writes two categories of script data: blocks and primitive values.
// Blocks are written using begin/end pairs; primitives using SpecifyValue.
// If a block returns true, a matching EndValues() must be called ( after visiting values or sub-blocks. )
// If a block returns false, end must *not* be called ( no values or sub-blocks are allowed. )
type Marshaler interface {
	// is the implementation writing json or reading it.
	IsEncoding() bool
	// starts a series of key-values pairs
	// the flow is closed ( written ) with a call to EndValues()
	MapValues(lede, typeName string) bool
	// literal values can be written without the surrounding map in the compact format.
	// ex. `{"type":"num_value","value":{"$NUM": {"type":"number","value":3}}}`
	// can be written as just `3`.
	MapLiteral(lede, typeName string) bool
	// the start of a possible key:value pair inside a flow or literal.
	// values are begin/end blocks or primitive values.
	// a new MapKey or an EndValues will cancel writing this pair.
	MapKey(sig, field string) bool
	// selects one of a small set of possible choices
	// the swap is closed ( written ) with a call to EndValues()
	PickValues(typeName string, vp Picker) bool
	// starts a series of values
	// the repeat is closed ( written ) with a call to EndValues()
	RepeatValues(typeName string, vs Slicer) bool
	// ends a flow, swap, or repeat.
	EndValues()
	// specify a single value
	GenericValue(typeName string, pv interface{})
	// sets a unique id for the next block or primitive value.
	SetCursor(id string)
	// record an error but don't terminate
	Warning(err error)
	// record an error and terminate
	Error(err error)
}

type Slicer interface {
	GetSize() int
	SetSize(int)
}

type Picker interface {
	GetChoice() (string, bool)
	SetChoice(string) (interface{}, bool)
}
