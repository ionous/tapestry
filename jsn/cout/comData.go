package cout

type flowData struct {
	sig    Sig
	values []interface{}
}

func newFlowData(lede string) *flowData {
	var cf flowData
	cf.sig.WriteLede(lede)
	return &cf
}

func (cf *flowData) addMsg(label string, value interface{}) {
	cf.sig.WriteLabel(label)
	cf.values = append(cf.values, value)
}

func (cf *flowData) finalize() (ret interface{}) {
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
