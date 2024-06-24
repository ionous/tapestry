package query

import (
	"encoding/json"
	"log"
)

// handle reading from db, and writing as raw json
type Bytes []byte

func (bs *Bytes) UnmarshalJSON(b []byte) error {
	(*bs) = b
	return nil
}

func (bs Bytes) MarshalJSON() (ret []byte, err error) {
	j := json.RawMessage(bs)
	return j.MarshalJSON()
}

// database scanner for null / bytes
func (bs *Bytes) Scan(v any) (_ error) {
	switch v := v.(type) {
	default:
		log.Panicf("scan %T", v)
	case string:
		(*bs) = []byte(v)
	case nil:
		*bs = nil
	}
	return
}
