package jsonexp

import "encoding/json"

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

// note: single values are represented by node,
// repeating elements are represented by an array of nodes.
// that means that unmarshal functions have to open each message.
type Fields map[string]json.RawMessage

type Context struct {
	Source string
}

func (ctx *Context) NewType(string) (ret DetailedMarshaler, err error) {
	return // lookup by name and return a new element
}

type DetailedMarshaler interface {
	UnmarshalDetailed(Context, []byte) error
	MarshalDetailed(Context) ([]byte, error)
}
