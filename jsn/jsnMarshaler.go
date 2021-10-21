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
	// some types may need special handling in a particular scenarios
	CustomizedMarshal(typeName string) (CustomizedMarshal, bool)
	// sets a unique id for the next block or primitive value.
	SetCursor(id string)
	// record an error
	Error(err error)
	// current state: embedded into this one
	State
}

type State interface {
	// starts a series of key-values pairs
	// the flow is closed ( written ) with a call to EndValues()
	MapValues(lede, typeName string) bool
	// the start of a possible key:value pair inside a flow or literal.
	// values are begin/end blocks or primitive values.
	// a new MapKey or an EndValues will cancel writing this pair.
	MapKey(sig, field string) bool
	// selects one of an unbounded set of possible values
	// returns the value if it exists for future serialization
	SlotValues(typeName string, slot Spotter) bool
	// selects one of a closed set of possible values
	// the swap is closed ( written ) with a call to EndValues()
	PickValues(typeName string, pick Picker) bool
	// starts a series of values
	// the repeat is closed ( written ) with a call to EndValues()
	RepeatValues(typeName string, slice Slicer) bool
	// ends a flow, swap, or repeat.
	EndValues()
	// specify a primitive value or enum.
	MarshalValue(typeName string, pv interface{}) bool
}

type CustomizedMarshal func(Marshaler, interface{}) bool

type Slicer interface {
	GetSize() int
	SetSize(int)
}

type Picker interface {
	GetChoice() (string, bool)
	SetChoice(string) (interface{}, bool)
}

type Spotter interface {
	HasSlot() bool
	SetSlot(interface{}) bool
}
