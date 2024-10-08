package decode

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/lang/compact"
)

// unpack a map of plain values into a description of a tapestry command.
// doesn't operate recursively, and doesn't check to see if the command is one that exists.
func DecodeMessage(msg map[string]any) (ret compact.Message, err error) {
	defer func() {
		if err != nil {
			pos := compact.MakeSource(msg)
			err = compact.MakeSourceError(pos, err)
		}
	}()
	var out compact.Message
	for k, v := range msg {
		if compact.IsMarkup(k) {
			out.AddMarkup(k, v)
		} else if len(out.Lede) > 0 {
			err = fmt.Errorf("expected a single key, have %q and %q", out.Lede, k)
			break
		} else if sig, e := DecodeSignature(k); e != nil {
			err = e
		} else if args, e := parseArgs(len(sig)-1, v); e != nil {
			err = e
		} else {
			out.Key = k
			out.Lede = sig[0]
			out.Labels = sig[1:]
			out.Args = args
			continue // keep going to read markup and detect multiple signatures
		}
	}
	if err == nil {
		// hrm: when there's only a comment { "comment": "..." }
		// report the key as having been the comment marker
		if len(out.Key) == 0 && len(out.Markup) > 0 {
			out.Key = compact.Comment
		}
		ret = out
	}
	return
}

// same as decode; errors if the passed value is something other than a map of plain values.
func ParseMessage(v any) (ret compact.Message, err error) {
	if m, ok := v.(map[string]any); !ok {
		err = fmt.Errorf("expected json map but found %T(%v)", v, v)
	} else {
		ret, err = DecodeMessage(m)
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
