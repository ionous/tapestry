package js

import (
	"bytes"
	"encoding/json"

	"github.com/ionous/errutil"
)

// https://github.com/golang/go/issues/27179
type MapItem struct {
	Key   string
	Value json.RawMessage
}

type MapSlice []MapItem

// expects that we're unmarshaling a map
// ex. json.Unmarshal(data, om)
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

func skipValue(d *json.Decoder, data []byte) (ret []byte, err error) {
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
					ret = data[start:end]
				}
			case json.Delim('['), json.Delim('{'):
				if depth == 0 {
					start = d.InputOffset() - 1
				}
				depth++
			case json.Delim(']'), json.Delim('}'):
				if depth = depth - 1; depth < 0 {
					// never had any open tokens
					err = errutil.New("invalid end of array or object")
				} else if depth == 0 {
					// had an open token, and returned from it.
					end := d.InputOffset()
					ret = data[start:end]
				}
			}
		}
	}
	return
}

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
