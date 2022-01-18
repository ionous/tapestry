package cout

import "strings"

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
	var cmt interface{}
	if str := cf.comment; len(str) > 0 {
		if lines := strings.Split(str, "\n"); len(lines) > 1 {
			cmt = lines
		} else {
			cmt = str
		}
	}

	switch vals := cf.values; len(vals) {
	case 0:
		// note: originally i collapsed calls with zero args down to just a string
		// but in cases where commands get used to generate text --
		// there's no way to differentiate b/t a command of zero params and plain text.
		if sig != commentMarker {
			v = []interface{}{}
		} else if cmt == nil {
			v = ""
		} else {
			v, cmt = cmt, nil
		}
	case 1:
		v = vals[0]
	default:
		v = vals
	}
	m := map[string]interface{}{
		sig: v,
	}
	if cmt != nil {
		m[commentMarker] = cmt
	}
	return m
}

const commentMarker = "--"
