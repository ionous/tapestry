package jsn

// Marshaler outputs script data.
// there are two categories of data: blocks and primitive values.
// blocks are written using begin/end pairs; primitives using WriteValue.
type Marshaler interface {
	// starts a series of key-values pairs
	// the flow is closed ( written ) with a call to EndValues()
	MapValues(lede, kind string)
	// the start of a possible key:value pair inside a flow.
	// values are begin/end blocks or primitive values.
	// a new MapKey or an EndValues will cancel writing this pair.
	MapKey(sig, field string)
	// mark the one and only key:value pair for the flow
	// literal values can be written without the surrounding map in the compact format.
	// ex. `{"type":"num_value","value":{"$NUM": {"type":"number","value":3}}}`
	// can be written as just `3`.
	MapLiteral(field string)
	// writes a primitive value.
	WriteValue(kind string, value interface{})
	// writes a (potentially) enumerated value.
	WriteChoice(kind string, enum Enumeration)
	// selects one of a small set of possible choices
	// the swap is closed ( written ) with a call to EndValues()
	PickValues(kind, choice string)
	// starts a series of values ( probably hint long )
	// the repeat is closed ( written ) with a call to EndValues()
	RepeatValues(hint int)
	// ends a flow, swap, or repeat.
	EndValues()
	// sets a unique id for the next block or primitive value.
	SetCursor(id string)
	// record an error but don't terminate
	Warning(err error)
	// record an error and terminate
	Error(err error)
}

type Enumeration interface {
	String() string
	FindChoice() (ret string, found bool)
}

type Marshalee interface {
	// UnmarshalDetailed(Context, []byte)
	Marshal(Marshaler)
}
