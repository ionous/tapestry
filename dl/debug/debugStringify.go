package debug

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// turn a tapestry generic value into a json(ish) formatted string
func Stringify(v rt.Value) (ret string) {
	var out js.Builder
	stringifyValue(&out, v)
	return out.String()
}

func stringifyValue(out *js.Builder, v rt.Value) {
	switch a := v.Affinity(); a {
	case affine.None:
		out.Str(js.Null)
	case affine.Bool:
		el := v.Bool()
		out.B(el)
	case affine.Num:
		el := v.Float()
		out.F(el)
	case affine.Text:
		el := v.String()
		out.Q(el)
	case affine.Record:
		rec := v.Record()
		stringifyRecord(out, rec)
	case affine.NumList:
		els := v.Floats()
		out.Brace(js.Array, func(_ *js.Builder) {
			for i, el := range els {
				if i > 0 {
					out.R(js.Comma)
				}
				out.F(el)
			}
		})
	case affine.TextList:
		els := v.Strings()
		out.Brace(js.Array, func(_ *js.Builder) {
			for i, el := range els {
				if i > 0 {
					out.R(js.Comma)
				}
				out.Q(el)
			}
		})
	case affine.RecordList:
		els := v.Records()
		out.Brace(js.Array, func(_ *js.Builder) {
			for i, el := range els {
				if i > 0 {
					out.R(js.Comma)
				}
				stringifyRecord(out, el)
			}
		})
	default:
		panic("unknown affinity")
	}
}

func stringifyRecord(out *js.Builder, rec *rt.Record) {
	out.Brace(js.Obj, func(_ *js.Builder) {
		if rec != nil {
			for i, cnt := 0, rec.FieldCount(); i < cnt; i++ {
				if i > 0 {
					out.R(js.Comma)
				}
				field := rec.Field(i).Name
				if v, e := rec.GetIndexedField(i); e != nil {
					if !rt.IsNilRecord(e) {
						// or panic, i suppose.
						out.Kv(field, "ERROR: "+e.Error())
					} else {
						out.Q(field).R(js.Colon).Raw("null")
					}
				} else {
					stringifyValue(out.Q(field).R(js.Colon), v)
				}
			}
		}
	})
}
