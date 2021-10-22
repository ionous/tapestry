package tables

import (
	"bytes"
	"encoding/gob"
	r "reflect"

	"github.com/ionous/errutil"
)

// GobScanner - reads from the database directly into the targeted (composer) value
type GobScanner struct {
	Target r.Value
}

func (gs *GobScanner) Scan(val interface{}) (err error) {
	if b, ok := val.([]byte); !ok {
		err = errutil.Fmt("gob scanner received unexpected type %T", val)
	} else {
		dec := gob.NewDecoder(bytes.NewBuffer(b))
		err = dec.DecodeValue(gs.Target)
	}
	return
}

func EncodeGob(cmd interface{}) (ret []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if e := enc.Encode(cmd); e != nil {
		err = errutil.New("EncodeGob, error encoding value", e)
	} else {
		ret = buf.Bytes()
	}
	return
}

func DecodeGob(prog []byte, outPtr interface{}) (err error) {
	dec := gob.NewDecoder(bytes.NewBuffer(prog))
	if e := dec.Decode(outPtr); e != nil {
		err = errutil.New("DecodeGob, error decoding value", e)
	}
	return
}
