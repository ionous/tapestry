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
		if strings.HasPrefix(k, markupMarker) {
			if key := k[len(markupMarker):]; len(key) == 0 {
				out.AddMarkup(markup.Comment, v)
			} else {
				out.AddMarkup(key, v)
			}
		} else if len(out.Name) > 0 {
			err = fmt.Errorf("expected a single key, have %q and %q", out.Name, k)
			break
		} else if sig, e := DecodeSignature(k); e != nil {
			err = e
		} else {
			out.Signature, out.Body = sig, v
			continue // keep going to detect any additional (incorrect) signatures
		}
	}
	if err == nil {
		ret = out
	}
	return
}

const markupMarker = "--"
