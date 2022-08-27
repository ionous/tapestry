package cout

type comFlow struct {
	sig    Sig
	params []any
	markup map[string]any
}

func newComFlow(lede string, markup map[string]any) *comFlow {
	var cf comFlow
	cf.sig.WriteLede(lede)
	cf.markup = markup
	return &cf
}

func (cf *comFlow) addMsg(label string, value any) {
	cf.sig.WriteLabel(label)
	cf.params = append(cf.params, value)
}

func (cf *comFlow) addMsgPair(label, choice string, value any) {
	cf.sig.WriteLabelPair(label, choice)
	cf.params = append(cf.params, value)
}

// build a map that we get serialized to json
func (cf *comFlow) finalize() map[string]any {
	m := make(map[string]any)
	if sig := cf.sig.String(); len(sig) > 0 && sig != markupMarker {
		switch vals := cf.params; len(vals) {
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
	for k, v := range cf.markup {
		if k == "comment" {
			// { "--": "here's a story of a lovely comment, which was writing up some very lovely words." }
			m[markupMarker] = v
		} else {
			// { "--color": 5 }
			m[markupMarker+k] = v
		}
	}
	return m
}

const markupMarker = "--"
