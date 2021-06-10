package jsonexp

import "encoding/json"

type String struct {
	Id    string `json:"id,omitempty"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Float struct {
	Id    string  `json:"id,omitempty"`
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

type Slot struct {
	Id    string          `json:"id,omitempty"`
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

type Flow struct {
	Id    string                     `json:"id,omitempty"`
	Type  string                     `json:"type"`
	Value map[string]json.RawMessage `json:"value"`
}
