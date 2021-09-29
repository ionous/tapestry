package jsn

// Marshalee for types which fit into slots
type Marshalee interface {
	Marshal(Marshaler)
}

// Marshaler outputs two categories of script data: blocks and primitive values.
// Blocks are written using begin/end pairs; primitives using SpecifyValue.
// If a block returns true, a matching EndValues() must be called ( after visiting values or sub-blocks. )
// If a block returns false, end must *not* be called ( no values or sub-blocks are allowed. )
type Marshaler interface {
	// starts a series of key-values pairs
	// the flow is closed ( written ) with a call to EndValues()
	MapValues(lede, kind string) bool
	// the start of a possible key:value pair inside a flow.
	// values are begin/end blocks or primitive values.
	// a new MapKey or an EndValues will cancel writing this pair.
	MapKey(sig, field string) bool
	// mark the one and only key:value pair for the flow
	// literal values can be written without the surrounding map in the compact format.
	// ex. `{"type":"num_value","value":{"$NUM": {"type":"number","value":3}}}`
	// can be written as just `3`.
	MapLiteral(field string) bool
	// selects one of a small set of possible choices
	// the swap is closed ( written ) with a call to EndValues()
	PickValues(Picker) bool
	// starts a series of values ( probably hint long )
	// the repeat is closed ( written ) with a call to EndValues()
	RepeatValues(hint int) bool
	// ends a flow, swap, or repeat.
	EndValues()
	// specify a true/false value.
	BoolValue(BoolMarshaler)
	// specify an enumerated value.
	EnumValue(EnumMarshaler)
	// specify a number value.
	NumValue(NumMarshaler)
	// specify a string value.
	StrValue(StrMarshaler)
	// sets a unique id for the next block or primitive value.
	SetCursor(id string)
	// record an error but don't terminate
	Warning(err error)
	// record an error and terminate
	Error(err error)
}

type BoolMarshaler interface {
	GetType() string
	GetBool() bool
	SetBool(bool)
}
type EnumMarshaler interface {
	GetType() string
	GetEnum() (key string, value string)
	SetEnum(keyOrValue string)
}
type NumMarshaler interface {
	GetType() string
	GetNum() float64
	SetNum(float64)
}
type StrMarshaler interface {
	GetType() string
	GetStr() string
	SetStr(string)
}
type Picker interface {
	GetType() string
	GetChoice() (string, bool)
	SetChoice(string) (interface{}, bool)
}
