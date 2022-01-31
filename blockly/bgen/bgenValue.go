package bgen

import "encoding/json"

// convert a value from the encoder to a jsonable output
// ex. a string to `"a string"`
func valueToBytes(pv interface{}) (ret json.RawMessage, err error) {
	if b, e := json.Marshal(unpackValue(pv)); e != nil {
		err = e
	} else {
		ret = b
	}
	return
}

// where pv is from the encoding statemachine
func unpackValue(pv interface{}) (ret interface{}) {
	switch pv := pv.(type) {
	// case interface{ GetCompactValue() interface{} }:
	// 	ret = pv.GetCompactValue()
	case interface{ GetValue() interface{} }:
		ret = pv.GetValue()
	default:
		ret = pv
	}
	return
}
