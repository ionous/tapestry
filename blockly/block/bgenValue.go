package block

import (
	"encoding/json"

	"git.sr.ht/~ionous/tapestry/lang/inspect"
)

// convert a value from the encoder to a jsonable output
// ex. a string to `"a string"`
func valueToBytes(w inspect.Iter) (ret json.RawMessage, err error) {
	v := w.NormalizedValue()
	if b, e := json.Marshal(v); e != nil {
		err = e
	} else {
		ret = b
	}
	return
}
