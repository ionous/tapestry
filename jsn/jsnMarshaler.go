package jsn

import "github.com/ionous/errutil"

// Marshalee for types which fit into slots
type Marshalee interface {
	Marshal(Marshaler) error
}

// Marshaler reads and writes two categories of script data: simple values, and block values.
// Block values are composed of simple values, or in some cases other block values.
// The expected implementation is a hierarchical statemachine of some sort, where
// "State" mutates during calls to other State functions.
type Marshaler interface {
	// is the implementation writing json or reading it.
	IsEncoding() bool
	// sets a unique id for the next block or primitive value.
	SetComment(*string)
	// report an error
	// errors for values inside of blocks are generally reported here
	// otherwise they are returned on the stack.
	Error(error)
	// current state: embedded into this one
	State
}

// Missing provides a standard return for when a value doesnt exist to marshal.
// Most often passed on the stack ( and not reported to the Marshaler ) to know not to descend into a block
const Missing = errutil.Error("Missing")

type State interface {
	// starts one of various block types.
	// if this succeeds (returns nil), a matching EndBlock() must be called ( after visiting any sub values. )
	// If this returns an error, EndBlock() must *not* be called ( and no sub values are allowed. )
	MarshalBlock(Block) error
	// the start of a possible key:value pair inside a FlowBlock.
	// a new MarshalKey or an EndBlock will cancel writing this pair.
	// if this succeeds (returns nil), it must be followed by a MarshalValue call.
	// if this returns an error, its value must *not* be marshaled.
	MarshalKey(sig, field string) error
	// specify a simple value ( or enum. )
	// callers should generally report the error or pass it on, but not both.
	MarshalValue(typeName string, _ interface{}) error
	// designates the end of the current block
	EndBlock()
}

// fix: these could probably be moved internal to the particular machine.
type CustomizedMarshal func(Marshaler, interface{}) error

// designation for a type which has multiple sub-values.
type Block interface {
	GetType() string
}

// starts a series of key-values pairs
// the flow is closed ( written ) with a call to EndBlock()
type FlowBlock interface {
	Block
	GetLede() string
	GetFlow() interface{}
	SetFlow(interface{}) bool
}

// selects one of a closed set of possible values
// the swap is closed ( written ) with a call to EndBlock()
type SwapBlock interface {
	Block
	GetSwap() (retChoice string, retValue interface{})
	SetSwap(string) bool
}

// starts a series of values
// the repeat is closed ( written ) with a call to EndBlock()
// a size of -1 means a nil ( as opposed to an extant slice of zero size )
type SliceBlock interface {
	Block
	GetSize() int
	SetSize(int)
	MarshalEl(Marshaler, int) error
}

// selects one of an unbounded set of possible values
// returns the value if it exists for future serialization
type SlotBlock interface {
	Block
	GetSlot() (interface{}, bool)
	SetSlot(interface{}) bool
}
