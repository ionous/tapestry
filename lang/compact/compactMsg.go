package compact

import (
	"fmt"
)

// A decoding of a tapestry command
// ex. {
//     "Sig:": [....]
//     "--": "comment"
//     "--markup": "markup"
// }
type Message struct {
	Signature     // (full)Key, Name, []Params
	Body      any // most often an array, sometimes a single value
	Markup    map[string]any
}

// fix? maybe beter to set these by default instead of Body.

func (op *Message) Args() (ret []any, err error) {
	switch pn := len(op.Params); pn {
	case 0:
		// FIX: ensure that msg.Args are zero?
	case 1:
		// we require that single parameters concrete values
		ret = []any{op.Body}
	default:
		if slice, ok := op.Body.([]any); !ok {
			err = fmt.Errorf("expected a slice of arguments; got %T", op.Body)
		} else if an := len(slice); an != pn {
			err = fmt.Errorf("expected a slice of %d arguments, got %d arguments instead", pn, an)
		} else {
			ret = slice
		}
	}
	return
}

func (op *Message) AddMarkup(k string, v any) {
	if op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	op.Markup[k] = v
}
