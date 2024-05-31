package generate

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/decode"
)

type MessageMap map[string]any

// turn slices of labels and args into a map of { label => arg }
func MakeMessageMap(msg compact.Message) MessageMap {
	out := make(MessageMap)
	for i, l := range msg.Labels {
		out[l] = msg.Args[i]
	}
	return out
}

// expects that v is nil, a string, or a slice of strings
func (mm MessageMap) GetStrings(key string) (ret []string, err error) {
	v := mm[key]
	switch v := v.(type) {
	case nil:
		// okay.
	case string:
		ret = []string{v}
	case []any:
		if vs, ok := compact.SliceStrings(v); !ok {
			err = fmt.Errorf("expected strings")
		} else {
			ret = vs
		}
	default:
		err = fmt.Errorf("expected strings, have %T", v)
	}
	return
}

// expects that v is nil, or a slice of messages.
func (mm MessageMap) GetMessages(key string) (ret []compact.Message, err error) {
	v := mm[key]
	if els, ok := v.([]any); !ok && v != nil {
		err = fmt.Errorf("%q expected a list, have %T %#v", key, v, v)
	} else {
		for _, el := range els {
			if msg, e := decode.ParseMessage(el); e != nil {
				err = e
				break
			} else {
				ret = append(ret, msg)
			}
		}
	}
	return
}

// returns false if there was a value, but it wasnt a bool
func (mm MessageMap) GetBool(name string) (ret bool, err error) {
	if v, ok := mm[name]; !ok {
		ret = false
	} else if a, ok := v.(bool); ok {
		ret = a
	} else {
		err = fmt.Errorf("expected a bool for %q", name)
	}
	return
}

// returns false if there was a value, but it wasnt a string
func (mm MessageMap) GetString(key, otherwise string) (ret string, err error) {
	v := mm[key]
	if v == nil && len(otherwise) > 0 {
		ret = otherwise
	} else if a, ok := v.(string); ok {
		ret = a
	} else {
		err = fmt.Errorf("expected a string for %q", key)
	}
	return
}
