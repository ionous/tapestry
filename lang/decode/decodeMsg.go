package decode

import (
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/markup"
)

// expects a map of string to value
func DecodeMessage(msg map[string]any) (ret compact.Message, err error) {
	var out compact.Message
	for k, v := range msg {
		if strings.HasPrefix(k, markup.Marker) {
			if key := k[len(markup.Marker):]; len(key) == 0 {
				out.AddMarkup(markup.Comment, v)
			} else {
				out.AddMarkup(key, v)
			}
		} else if len(out.Name) > 0 {
			err = fmt.Errorf("expected a single key, have %q and %q", out.Name, k)
			break
		} else if sig, e := DecodeSignature(k); e != nil {
			err = e
		} else if args, e := parseArgs(len(sig)-1, v); e != nil {
			err = e
		} else {
			out.Key = k
			out.Name = sig[0]
			out.Labels = sig[1:]
			out.Args = args
			continue // keep going to detect any additional (incorrect) signatures
		}
	}
	if err == nil {
		ret = out
	}
	return
}

func parseArgs(pn int, body any) (ret []any, err error) {
	switch pn {
	case 0:
		// FIX: ensure that msg.Args are zero?
	case 1:
		// we require that single parameters concrete values
		ret = []any{body}
	default:
		if slice, ok := body.([]any); !ok {
			err = fmt.Errorf("expected a slice of arguments; got %T", body)
		} else if an := len(slice); an != pn {
			err = fmt.Errorf("expected a slice of %d arguments, got %d arguments instead", pn, an)
		} else {
			ret = slice
		}
	}
	return
}
