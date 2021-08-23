package jsonexp

import "encoding/json"

type CompactFlow struct {
	Sig
	Fields []json.RawMessage
}

func (cf *CompactFlow) AddMsg(label string, field json.RawMessage) {
	cf.Sig.WriteLabel(label)
	cf.Fields = append(cf.Fields, field)
}

func (cf *CompactFlow) MarshalJSON() ([]byte, error) {
	var i interface{}
	if len(cf.Fields) == 1 {
		i = cf.Fields[0]
	} else {
		i = cf.Fields
	}
	return json.Marshal(map[string]interface{}{
		cf.Sig.String(): i,
	})
}
