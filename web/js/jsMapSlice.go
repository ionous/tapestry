package js

import (
	"bytes"
	"encoding/json"

	"github.com/ionous/errutil"
)

// https://github.com/golang/go/issues/27179
type MapItem struct {
	Key string
	Msg json.RawMessage
}

type MapSlice []MapItem

// returns a valid pointer into the slice, or nil if not found
func (om MapSlice) Find(k string) (ret *MapItem, okay bool) {
	if at := om.FindIndex(k); at >= 0 {
		ret, okay = &(om[at]), true
	}
	return
}

// returns the index of the item or -1 if not found
func (om MapSlice) FindIndex(k string) (ret int) {
	ret = -1 // provisionally
	for i, kv := range om {
		if kv.Key == k {
			ret = i
			break
		}
	}
	return
}

// expects we're unmarshaling a valid json object.
func (om *MapSlice) UnmarshalJSON(data []byte) (err error) {
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber() // so we can determine the width of the original value
	if t, e := d.Token(); e != nil {
		err = e
	} else if t != json.Delim('{') {
		err = errutil.Fmt("expected start of object, got %v", t)
	} else {
		for {
			if t, e := d.Token(); e != nil {
				err = e
				break
			} else if t == json.Delim('}') {
				break // done
			} else if str, ok := t.(string); !ok {
				err = errutil.Fmt("expected a string, got %T", t)
			} else {
				if span, e := skipValue(d, data); e != nil {
					err = e
					break
				} else {
					(*om) = append(*om, MapItem{str, span})
				}
			}
		}
	}
	return
}

// read through the passed json value until its end.
// return that isolated value, excluding any starting or ending whitespace.
func skipValue(d *json.Decoder, value []byte) (ret []byte, err error) {
	var start, depth int64
	for err == nil && ret == nil {
		if t, e := d.Token(); e != nil {
			err = e
		} else {
			switch t {
			default:
				// we havent started processing an array or object
				// so we're skipping just one thing.
				if depth == 0 {
					end := d.InputOffset()
					start := end - int64(width(t))
					ret = value[start:end]
				}
			case json.Delim('['), json.Delim('{'):
				if depth == 0 {
					start = d.InputOffset() - 1 // -1 includes the delimiter in the returned value
				}
				depth++
			case json.Delim(']'), json.Delim('}'):
				if depth = depth - 1; depth < 0 {
					// closed but never had any open tokens
					err = errutil.New("invalid end of array or object", t)
				} else if depth == 0 {
					// had an open token, and now closed it.
					// its not too picky on the type of token.
					// assumes the json is valid.
					end := d.InputOffset()
					ret = value[start:end]
				}
			}
		}
	}
	return
}

// return the length in bytes of the passed token
// requires the decoder ( see MapSlice UnmarshaJSON ) to be in "UseNumber" mode.
func width(t json.Token) (ret int) {
	switch v := t.(type) {
	case json.Delim:
		ret = 1
	case bool:
		if v {
			ret = len("true")
		} else {
			ret = len("false")
		}
	case string:
		ret = len(v) + 2 // open and close quotes
	case nil:
		ret = len("null")
	case json.Number:
		ret = len(v)
	default:
		panic("unexpected")
	}
	return
}
