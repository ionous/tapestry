package dump

import (
	"database/sql"
	"fmt"
	"slices"

	"git.sr.ht/~ionous/tapestry/qna/decoder"
	"git.sr.ht/~ionous/tapestry/qna/raw"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/pack"
	"git.sr.ht/~ionous/tapestry/tables"
)

type valueReader struct {
	kd decoder.KindDecoder
	n  *raw.NounData
	k  *rt.Kind
	ff fieldFinder
}

func makeValueReader(n *raw.NounData, k *rt.Kind, kd decoder.KindDecoder) valueReader {
	return valueReader{kd, n, k, 0}
}

func (vr *valueReader) values(rows *sql.Rows) (err error) {
	var v struct {
		Field string
		Bytes []byte
	}
	var last string
	err = tables.ScanAll(rows, func() (_ error) {
		if last != v.Field { // the final values are stored first
			err = vr.decodeField(v.Field, v.Bytes)
			last = v.Field
		}
		return
	}, &v.Field, &v.Bytes)
	vr.ff.Reset()
	return
}

func (vr *valueReader) decodeField(field string, b []byte) (err error) {
	if ft, e := vr.ff.FindField(vr.k, field); e != nil {
		err = e
	} else if v, e := vr.kd.DecodeAssignment(ft.Affinity, b); e != nil {
		err = e
	} else {
		vr.n.Values = append(vr.n.Values, raw.EvalData{
			Field: ft.Name,
			Value: v,
		})
	}
	return
}

// records are stored in pieces *sigh*
func (vr *valueReader) records(rows *sql.Rows) (err error) {
	var next struct {
		Field string
		decoder.SparseValue
	}
	var lastField string
	var vs []decoder.SparseValue
	if e := tables.ScanAll(rows, func() (err error) {
		// FIX: also have to handle final? what does the runtime do?
		if next.Field != lastField {
			err = vr.flushRecord(lastField, vs)
			lastField = next.Field
			vs = nil
		}
		vs = append(vs, next.SparseValue)
		return
	}, &next.Field, &next.Path, &next.Bytes); e != nil {
		err = e
	} else {
		err = vr.flushRecord(lastField, vs)
	}
	vr.ff.Reset()
	return
}

func (vr *valueReader) flushRecord(field string, vs []decoder.SparseValue) (err error) {
	if len(vs) > 0 {
		if ft, e := vr.ff.FindField(vr.k, field); e != nil {
			err = e
		} else if rec, e := decoder.DecodeSparseRecord(vr.kd, ft.Type, vs); e != nil {
			err = e
		} else {
			packed := pack.PackRecord(rec)
			// println(packed)
			vr.n.Records = append(vr.n.Records, raw.RecordData{
				Field:  field,
				Packed: []byte(packed),
			})
		}
	}
	return
}

// doesnt look for aspects
type fieldFinder int

func (ff *fieldFinder) Reset() {
	(*ff) = fieldFinder(0)
}

func (ff *fieldFinder) FindField(k *rt.Kind, field string) (ret rt.Field, err error) {
	last := int(*ff)
	if i := slices.IndexFunc(k.Fields[last:], func(f rt.Field) bool {
		return f.Name == field
	}); i < 0 {
		err = fmt.Errorf("field %q not found in %q", k.Name(), field)
	} else {
		ofs := last + i
		(*ff) = fieldFinder(ofs)
		ret = k.Fields[ofs]
	}
	return
}
