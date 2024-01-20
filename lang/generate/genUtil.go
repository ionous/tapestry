package generate

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/decode"
)

func parseString(name string, v any, otherwise string) (ret string, err error) {
	if v == nil && len(otherwise) > 0 {
		ret = otherwise
	} else if a, ok := v.(string); ok {
		ret = a
	} else {
		err = fmt.Errorf("expected a string for %q", name)
	}
	return
}

func parseStrings(v any) (ret []string, err error) {
	switch v := v.(type) {
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

func parseMessages(v any) (ret []compact.Message, err error) {
	if els, ok := v.([]any); !ok && v != nil {
		err = fmt.Errorf("expected a list, have %T", v)
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

type MessageMap map[string]any

// turn labels, args into a map of label => arg
func messageMap(msg compact.Message) MessageMap {
	out := make(MessageMap)
	for i, l := range msg.Labels {
		out[l] = msg.Args[i]
	}
	return out
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
func (mm MessageMap) GetString(name, otherwise string) (ret string, err error) {
	return parseString(name, mm[name], otherwise)
}
