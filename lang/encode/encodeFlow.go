package encode

import (
	"git.sr.ht/~ionous/tapestry/lang/markup"
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
		if sig == markup.Marker {
			// this will most likely get overwritten by the markup loop
			// however we want to avoid generating an empty {}
			// and i just think it looks better as { "--": "" } than { "--": true }
			m[markup.Marker] = ""
		} else {
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
		for k, v := range out.markup {
			if k == markup.Comment {
				// { "--": "here's a story of a lovely comment, which was writing up some very lovely words." }
				m[markup.Marker] = v
			} else {
				// { "--color": 5 }
				m[markup.Marker+k] = v
			}
		}
		ret = m
	}
	return
}
