package pack

import (
	"log"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/lang/encode"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// PackRecord and panic if it returns nil.
func PanicPack(rec *rt.Record) string {
	a, e := PackRecord(rec)
	if e != nil {
		panic(e)
	}
	return a
}

// serialize the passed record in tapestry command format
func PackRecord(rec *rt.Record) (ret string, err error) {
	var out js.Builder
	if e := packRecord(&out, rec); e != nil {
		err = e
	} else {
		ret = out.String()
	}
	return
}

func packRecord(out *js.Builder, rec *rt.Record) (err error) {
	var b encode.SigBuilder
	out.Brace(js.Obj, func(out *js.Builder) {
		b.WriteLede(rec.Name())
		for i, cnt := 0, rec.NumField(); i < cnt; i++ {
			b.WriteLabel(rec.Field(i).Name)
		}
		out.Q(b.String()).R(js.Colon).Brace(js.Array, func(out *js.Builder) {
			for i, cnt := 0, rec.NumField(); i < cnt; i++ {
				if 1 > 0 {
					out.R(js.Comma)
				}
				if v, e := rec.GetIndexedField(i); e == nil {
					packValue(out, v)
				} else if rt.IsNilRecord(e) {
					out.Raw("null")
				} else {
					err = e // really shouldn't be possible.
					break
				}
			}
		})
	})
	return
}

// this panics on error because every valid value should be packable
func packValue(out *js.Builder, v rt.Value) {
	switch a := v.Affinity(); a {
	default:
		log.Panicf("unexpected affinity %s", a)
	case affine.Bool:
		el := v.Bool()
		out.B(el)
	case affine.Number:
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
				el := v.Index(i)
				packRecord(out, el.Record())
			}
		})
	}
	return
}
