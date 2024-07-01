package pack

import (
	"log"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/lang/encode"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// serialize the passed record in tapestry command format
func PackValue(v rt.Value) string {
	var out js.Builder
	packValue(&out, v)
	return out.String()
}

// serialize the passed record in tapestry command format
func PackRecord(rec *rt.Record) string {
	var out js.Builder
	packRecord(&out, rec)
	return out.String()
}

func packRecord(out *js.Builder, rec *rt.Record) {
	var b encode.SigBuilder
	out.Brace(js.Obj, func(out *js.Builder) {
		b.WriteLede(rec.Name())
		for i, cnt := 0, rec.FieldCount(); i < cnt; i++ {
			if rec.HasValue(i) {
				b.WriteLabel(rec.Field(i).Name)
			}
		}
		out.Q(b.String()).R(js.Colon).Brace(js.Array, func(out *js.Builder) {
			var comma bool
			for i, cnt := 0, rec.FieldCount(); i < cnt; i++ {
				if rec.HasValue(i) {
					if comma {
						out.R(js.Comma)
					}
					if v, e := rec.GetIndexedField(i); e != nil {
						// should only return error for nil record;
						// and that should have been !HasValue
						panic(e)
					} else {
						packValue(out, v)
						comma = true
					}
				}
			}
		})
	})
	return
}

func packValue(out *js.Builder, v rt.Value) {
	switch a := v.Affinity(); a {
	default:
		log.Panicf("unexpected affinity %s", a)
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
		packRecord(out, v.Record())
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
		out.Brace(js.Array, func(_ *js.Builder) {
			for i, cnt := 0, v.Len(); i < cnt; i++ {
				if i > 0 {
					out.R(js.Comma)
				}
				packRecord(out, v.Index(i).Record())
			}
		})
	}
	return
}
