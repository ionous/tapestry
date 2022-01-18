package cout

type comFlow struct {
	sig     Sig
	values  []interface{}
	comment string
}

func newComFlow(lede, comment string) *comFlow {
	var cf comFlow
	cf.sig.WriteLede(lede)
	cf.comment = comment
	return &cf
}

func (cf *comFlow) addMsg(label string, value interface{}) {
	cf.sig.WriteLabel(label)
	cf.values = append(cf.values, value)
}

func (cf *comFlow) addMsgPair(label, choice string, value interface{}) {
	cf.sig.WriteLabelPair(label, choice)
	cf.values = append(cf.values, value)
}

func (cf *comFlow) finalize() (ret interface{}) {
	sig := cf.sig.String()
	var v interface{}
	switch vals := cf.values; len(vals) {
	case 0:
		// note: originally i collapsed calls with zero args down to just a string
		// but in cases where commands get used to generate text --
		// there's no way to differentiate b/t a command of zero params and plain text.
		if sig == commentMarker {
			v, cf.comment = cf.comment, ""
		} else {
			v = []interface{}{}
		}
	case 1:
		v = vals[0]
	default:
		v = vals
	}
	m := map[string]interface{}{
		sig: v,
	}
	if len(cf.comment) > 0 {
		m[commentMarker] = cf.comment
	}
	return m
}

const commentMarker = "--"
