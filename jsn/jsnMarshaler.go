package jsn

// Marshalee for types which fit into slots
type Marshalee interface {
	Marshal(Marshaler)
}

// Marshaler reads and writes two categories of script data: simple values, and block values.
// Block values are composed of simple values, or in some cases other block values.
// The expected implementation is a hierarchical statemachine of some sort, where
// "State" mutates during calls to other State functions.
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
	// starts one of various block types.
	// If this returns true, a matching EndBlock() must be called ( after visiting any sub values. )
	// If this returns false, EndBlock() must *not* be called ( and no sub values are allowed. )
	MarshalBlock(BlockType) bool
	// the start of a possible key:value pair inside a FlowBlock.
	// a new MarshalKey or an EndBlock will cancel writing this pair.
	// if this returns true, the corresponding value must be marshaled.
	MarshalKey(sig, field string) bool
	// specify a simple value ( or enum. )
	// returns true/false only to improve compatibility with MarshalBlock;
	// callers generally should generally not act on the return value.
	MarshalValue(typeName string, _ interface{}) bool
	// designates the end of the current block
	EndBlock()
}

type CustomizedMarshal func(Marshaler, interface{}) bool

// designation for a type which has multiple sub-values.
type BlockType interface {
	GetType() string
}

// starts a series of key-values pairs
// the flow is closed ( written ) with a call to EndBlock()
type FlowBlock interface {
	BlockType
	GetLede() string
}

// selects one of a closed set of possible values
// the swap is closed ( written ) with a call to EndBlock()
type SwapBlock interface {
	BlockType
	GetChoice() (string, bool)
	SetChoice(string) (interface{}, bool)
}
type Picker = SwapBlock

// starts a series of values
// the repeat is closed ( written ) with a call to EndBlock()
type SliceBlock interface {
	BlockType
	GetSize() int
	SetSize(int)
}
type Slicer = SliceBlock

// selects one of an unbounded set of possible values
// returns the value if it exists for future serialization
type SlotBlock interface {
	BlockType
	HasSlot() bool
	SetSlot(interface{}) bool
}
type Spotter = SlotBlock
