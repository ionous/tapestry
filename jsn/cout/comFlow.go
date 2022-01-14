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
	if cnt := len(cf.values); cnt == 0 {
		ret = sig
	} else {
		var v interface{}
		if cnt == 1 {
			v = cf.values[0]
		} else {
			v = cf.values
		}
		m := map[string]interface{}{
			sig: v,
		}
		if len(cf.comment) > 0 {
			m["--"] = cf.comment
		}
		ret = m
	}
	return
}
