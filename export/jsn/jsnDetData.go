package jsn

type detMap struct {
	Id     string                 `json:"id,omitempty"`
	Type   string                 `json:"type"`
	Fields map[string]interface{} `json:"value"`
}

type detValueData struct {
	Id    string      `json:"id,omitempty"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}
