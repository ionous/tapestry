package jsonexp

import (
	"encoding/json"
)

type Str struct {
	Id    string `json:"id,omitempty"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Num struct {
	Id    string  `json:"id,omitempty"`
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

type Node struct {
	Id    string          `json:"id,omitempty"`
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

type Flow struct {
	Id     string `json:"id,omitempty"`
	Type   string `json:"type"`
	Fields Fields `json:"value"`
}

type Slot struct {
	Id   string `json:"id,omitempty"`
	Type string `json:"type"`
	Slat *Node  `json:"value"`
}

// note: single values are represented by node,
// repeating elements are represented by an array of nodes.
// that means that unmarshal functions have to open each message.
type Fields map[string]json.RawMessage

// NewType and Finalize are the first and last stages of reading a slot;
// sandwiched between is the de-serialization of the slat which fills it.
type Context interface {
	Source() string                            // filename, etc.
	NewType(t string) (interface{}, error)     // new by name
	Finalize(interface{}) (interface{}, error) // do something with the passed ptr
}

type DetailedMarshaler interface {
	UnmarshalDetailed(Context, []byte) error
	MarshalDetailed(Context) ([]byte, error)
}
