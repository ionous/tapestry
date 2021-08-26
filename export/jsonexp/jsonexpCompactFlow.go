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

func (cf *CompactFlow) MarshalJSON() (ret []byte, err error) {
	sig := cf.Sig.String()
	if cnt := len(cf.Fields); cnt == 0 {
		ret, err = json.Marshal(sig)
	} else {
		var i interface{}
		if cnt == 1 {
			i = cf.Fields[0]
		} else {
			i = cf.Fields
		}
		ret, err = json.Marshal(map[string]interface{}{
			sig: i,
		})
	}
	return
}
