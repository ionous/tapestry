package detailed

type detMap struct {
	Id     string                 `json:"id,omitempty"`
	Type   string                 `json:"type"`
	Fields map[string]interface{} `json:"value"`
}

type detValue struct {
	Id    string      `json:"id,omitempty"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}
