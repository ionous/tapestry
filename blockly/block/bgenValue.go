package block

import (
	"encoding/json"
	"errors"
	"strings"

	"git.sr.ht/~ionous/tapestry/lang/inspect"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// convert a value from the encoder to a jsonable output
// ex. a string to `"a string"`
func valueToBytes(w inspect.Iter) (ret json.RawMessage, err error) {
	if v, ok := getValue(w); !ok {
		err = errors.New("couldnt convert value")
	} else if b, e := json.Marshal(v); e != nil {
		err = e
	} else {
		ret = b
	}
	return
}

func getValue(w inspect.Iter) (ret any, okay bool) {
	t := w.TypeInfo()
	switch t := t.(type) {
	case *typeinfo.Num:
		ret, okay = w.Float(), true

	case *typeinfo.Str:
		s := w.String()
		if len(t.Options) > 0 {
			s = "$" + strings.ToUpper(s)
		}
		ret, okay = s, true
	}
	return
}
