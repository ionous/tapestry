package debug

import (
	"git.sr.ht/~ionous/tapestry/affine"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// turn a tapestry generic value into a json(ish) formatted string
func Stringify(v g.Value) (ret string) {
	var out js.Builder
	stringifyValue(&out, v)
	return out.String()
}

func stringifyValue(out *js.Builder, v g.Value) {
	switch a := v.Affinity(); a {
	case affine.None:
		out.Str(js.Null)
	case affine.Bool:
		el := v.Bool()
		out.B(el)
	case affine.Number:
		el := v.Float()
		out.F(el)
	case affine.Text:
		el := v.String()
		out.Q(el)
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
	case affine.Record:
		if rec, ok := v.Record(); !ok {
			out.Raw("null")
		} else {
			stringifyRecord(out, rec)
		}
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

func stringifyRecord(out *js.Builder, rec *g.Record) {
	out.Brace(js.Obj, func(_ *js.Builder) {
		if rec != nil {
			k := rec.Kind()
			for i, cnt := 0, k.NumField(); i < cnt; i++ {
				if rec.HasValue(i) {
					if v, e := rec.GetIndexedField(i); e != nil {
						panic(e)
					} else {
						if i > 0 {
							out.R(js.Comma)
						}
						f := k.Field(i)
						stringifyValue(out.Q(f.Name).R(js.Colon), v)
					}
				}
			}
		}
	})
}
