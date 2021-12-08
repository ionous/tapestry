package literal

import (
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
)

func CompactEncoder(m jsn.Marshaler, flow jsn.FlowBlock) (err error) {
	switch i, typeName := flow.GetFlow(), flow.GetType(); typeName {
	default:
		err = chart.Unhandled("CustomFlow")

	case BoolValue_Type:
		var out bool = i.(*BoolValue).Bool
		err = m.MarshalValue(typeName, out)

	case NumValue_Type:
		var out float64 = i.(*NumValue).Num
		err = m.MarshalValue(typeName, out)

	case NumValues_Type:
		var out []float64 = i.(*NumValues).Values
		err = m.MarshalValue(typeName, out)

		// write text as a raw string
	case TextValue_Type:
		str := i.(*TextValue).Text
		// if the text starts with an @, add another @
		if len(str) > 0 && str[0] == '@' {
			str = "@" + str
		}
		err = m.MarshalValue(typeName, str)

	case TextValues_Type:
		var out []string = i.(*TextValues).Values
		err = m.MarshalValue(typeName, out)

	}
	return
}
