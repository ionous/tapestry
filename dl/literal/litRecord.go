package literal

import (
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"github.com/ionous/errutil"
)

func marshalFields(out map[string]interface{}, vs []FieldValue) (err error) {
Loop:
	for _, fv := range vs {
		var i interface{}
		field, value := fv.Field, fv.Value
		switch v := value.(type) {
		case *BoolValue:
			i = v.Value
		case *NumValue:
			i = v.Value
		case *TextValue:
			i = v.Value
		case *NumValues:
			i = v.Values
		case *TextValues:
			i = v.Values
		case *FieldList:
			next := make(map[string]interface{})
			marshalFields(next, v.Fields)
			i = next
		default:
			err = errutil.Fmt("marshalFields unhandled literal %T", value)
			break Loop
		}
		out[field] = i
	}
	return
}

// reads a literal record
// fix: when is this used? shouldn't this be a pattern call?
func unmarshalFields(msg map[string]any) (ret []FieldValue, err error) {
Loop:
	for key, v := range msg {
		var i LiteralValue
		switch val := v.(type) {
		case bool:
			i = &BoolValue{Value: val}
		case string:
			i = &TextValue{Value: val}
		case float64:
			i = &NumValue{Value: val}
		case []any:
			if vs, ok := cin.SliceFloats(val); ok {
				i = &NumValues{Values: vs}
			} else if vs, ok := cin.SliceStrings(val); ok {
				i = &TextValues{Values: vs}
			}
		case map[string]any:
			if x, e := unmarshalFields(val); e != nil {
				err = e
				break Loop
			} else {
				i = &FieldList{Fields: x}
			}
		default:
			err = errutil.Fmt("unmarshalFields unhandled literal %T", val)
			break Loop
		}
		ret = append(ret, FieldValue{Field: key, Value: i})
	}
	return
}
