package din

import "encoding/json"

type dinMap struct {
	Id     string    `json:"id,omitempty"`
	Type   string    `json:"type"`
	Fields dinFields `json:"value"`
}

type dinValue struct {
	Id   string          `json:"id,omitempty"`
	Type string          `json:"type"`
	Msg  json.RawMessage `json:"value"`
}

type dinFields map[string]json.RawMessage
