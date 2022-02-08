package bgen

import (
	"encoding/json"

	"git.sr.ht/~ionous/tapestry/dl/composer"
)

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
	// ugh. for str values whether we want the compact value or the $KEY value
	// depends on how we are handling things in blockly
	var notCompact bool
	if c, ok := pv.(composer.Composer); ok { //ugh
		if spec := c.Compose(); !spec.OpenStrings {
			notCompact = true
		}
	}
	if v, ok := pv.(interface{ GetCompactValue() interface{} }); !notCompact && ok {
		ret = v.GetCompactValue()
	} else if v, ok := pv.(interface{ GetValue() interface{} }); ok {
		ret = v.GetValue()
	} else {
		ret = pv
	}
	return
}
