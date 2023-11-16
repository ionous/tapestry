package literal

import (
	r "reflect"

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
func unmarshalFields(msg r.Value) (ret []FieldValue, err error) {
	if t := msg.Type(); !cin.IsValidMap(t) {
		err = errutil.Fmt("expected a map, have %s", t)
	} else {
	Loop:
		for it := msg.MapRange(); it.Next(); {
			key, val := it.Key().String(), it.Value().Elem()
			var i LiteralValue
			switch val.Kind() {
			case r.Bool:
				i = &BoolValue{Value: val.Bool()}
			case r.String:
				i = &TextValue{Value: val.String()}
			case r.Float64:
				i = &NumValue{Value: val.Float()}
			case r.Slice:
				if vs, ok := cin.SliceFloats(msg); ok {
					i = &NumValues{Values: vs}
				} else if vs, ok := cin.SliceStrings(msg); ok {
					i = &TextValues{Values: vs}
				}
			case r.Map:
				if x, e := unmarshalFields(val); e != nil {
					err = e
					break Loop
				} else {
					i = &FieldList{Fields: x}
				}
			default:
				err = errutil.Fmt("unmarshalFields unhandled literal %T", val.Type())
				break Loop
			}
			ret = append(ret, FieldValue{Field: key, Value: i})
		}
	}
	return
}
