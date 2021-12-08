package core

import (
	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
)

func CompactEncoder(m jsn.Marshaler, flow jsn.FlowBlock) (err error) {
	switch i, typeName := flow.GetFlow(), flow.GetType(); typeName {
	// write variables as a string prepended by @
	case GetVar_Type:
		ptr := i.(*GetVar)
		str := ptr.Name.Str
		// a leading ampersand would conflict with @@ escaped text serialization.
		if leadingAmp := len(str) > 0 && str[0] == '@'; !leadingAmp {
			err = m.MarshalValue(typeName, "@"+str)
		} else {
			err = chart.Unhandled(typeName)
		}
	default:
		err = literal.CompactEncoder(m, flow)
	}
	return
}
