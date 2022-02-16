package block

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
	// str values depend on how we are handling things in blockly
	// for closed strings we use a dropdown, so we want $KEY
	var canBeCompact bool // note: BoxedBool doesn't implement composer...
	if c, ok := pv.(composer.Composer); ok {
		if spec := c.Compose(); spec.OpenStrings {
			canBeCompact = true
		}
	}
	// get/compact/value are implemented by enums ( see: jsnEnum.go; jsnBox.go )
	if v, ok := pv.(interface{ GetCompactValue() interface{} }); ok && canBeCompact {
		ret = v.GetCompactValue() // returns the more friendly value used in compact json
	} else if v, ok := pv.(interface{ GetValue() interface{} }); ok {
		ret = v.GetValue() // returns the $KEY value used in detailed json

	} else {
		ret = pv
	}
	return
}
