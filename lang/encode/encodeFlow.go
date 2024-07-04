package encode

import (
	"git.sr.ht/~ionous/tapestry/lang/compact"
)

type FlowBuilder struct {
	SigBuilder
	params []any
	markup map[string]any
}

func (out *FlowBuilder) WriteArg(label string, value any) {
	out.WriteLabel(label)
	out.params = append(out.params, value)
}

func (out *FlowBuilder) SetMarkup(markup map[string]any) {
	out.markup = markup
}

// build a map that we get serialized to json
func (out *FlowBuilder) FinalizeMap() (ret map[string]any) {
	if sig := out.String(); len(sig) > 0 {
		m := make(map[string]any)
		// is the signature only markup?
		isMarkup := compact.IsMarkup(sig)
		if !isMarkup {
			// this will most likely get overwritten by the markup loop
			switch vals := out.params; len(vals) {
			// zero parameters { "sig": true }
			case 0:
				// note: originally i collapsed calls with zero args down to just a string
				// but in cases where commands get used to generate text
				// there's no way to differentiate b/t a command of zero params and plain text.
				m[sig] = true
			// one parameter { "sig": value }
			case 1:
				m[sig] = vals[0]
			// multiple parameters { "sig": [comma,separated,values] }
			default:
				m[sig] = vals
			}
		}
		// { "color": 5 }
		for k, v := range out.markup {
			m[k] = v
		}
		// avoid writing nothing
		if len(m) == 0 && isMarkup {
			m[compact.Comment] = ""
		}
		ret = m
	}
	return
}
