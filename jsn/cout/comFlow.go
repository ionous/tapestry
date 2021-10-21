package cout

type comFlow struct {
	sig    Sig
	values []interface{}
}

func newComFlow(lede string) *comFlow {
	var cf comFlow
	cf.sig.WriteLede(lede)
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
		ret = map[string]interface{}{
			sig: v,
		}
	}
	return
}
