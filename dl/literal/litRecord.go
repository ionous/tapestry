package literal

import "github.com/ionous/errutil"

func marshalFields(out map[string]interface{}, vs []FieldValue) (err error) {
Loop:
	for _, fv := range vs {
		var i interface{}
		field, value := fv.Field, fv.Value
		switch v := value.(type) {
		case *BoolValue:
			i = v.Bool
		case *NumValue:
			i = v.Num
		case *TextValue:
			i = v.Text
		case *NumValues:
			i = v.Values
		case *TextValues:
			i = v.Values
		case *FieldValues:
			next := make(map[string]interface{})
			marshalFields(next, v.Contains)
			i = next
		default:
			err = errutil.Fmt("marshalFields unhandled literal %T", value)
			break Loop
		}
		out[field] = i
	}
	return
}

func unmarshalFields(m map[string]interface{}) (ret []FieldValue, err error) {
Loop:
	for key, val := range m {
		var i LiteralValue
		switch v := val.(type) {
		case bool:
			i = &BoolValue{Bool: v}
		case string:
			i = &TextValue{Text: v}
		case float64:
			i = &NumValue{Num: v}
		case []float64:
			i = &NumValues{Values: v}
		case []string:
			i = &TextValues{Values: v}
		case map[string]interface{}:
			if x, e := unmarshalFields(v); e != nil {
				err = e
				break Loop
			} else {
				i = &FieldValues{Contains: x}
			}
		default:
			err = errutil.Fmt("unmarshalFields unhandled literal %T", v)
			break Loop
		}
		ret = append(ret, FieldValue{Field: key, Value: i})
	}
	return
}
