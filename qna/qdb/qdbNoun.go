package qdb

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/qna/decoder"
	"git.sr.ht/~ionous/tapestry/rt"
)

// interpreting the value is left to the caller ( re: field affinity )
func (q *Query) NounValue(noun, field string) (ret rt.Assignment, err error) {
	if k, e := q.GetKindByNoun(noun); e != nil {
		err = e
	} else if i := k.FieldIndex(field); i < 0 {
		err = fmt.Errorf("couldnt find field %q in kind %q", field, k.Name())
	} else {
		ft := k.Field(i)
		if rows, e := q.nounValues.Query(noun, field); e != nil {
			err = e
		} else if vals, e := decoder.ScanSparseValues(rows); e != nil {
			err = e
		} else if len(vals) == 0 {
			ret, err = zeroAssignment(ft, i)
		} else if ft.Affinity != affine.Record {
			ret, err = decodeValue(q.dec, ft, vals)
		} else if rec, e := decoder.DecodeSparseRecord(
			kindDecoder{q, q.dec},
			ft.Type, vals); e != nil {
			err = e
		} else {
			out := rt.RecordOf(rec)
			ret = rt.AssignValue(out)
		}
	}
	return
}

type kindDecoder struct {
	rt.Kinds
	decoder.Decoder
}

func (q *Query) GetKindByNoun(noun string) (ret *rt.Kind, err error) {
	if k, e := scanString(q.nounKind, noun); e != nil {
		err = e
	} else {
		ret, err = q.GetKindByName(k)
	}
	return
}

func decodeValue(dec decoder.Decoder, ft rt.Field, vals []decoder.SparseValue) (ret rt.Assignment, err error) {
	if cnt := len(vals); cnt != 1 {
		err = fmt.Errorf("field %q of %s expected a single value, got %d", ft.Name, ft.Affinity, cnt)
	} else if v := vals[0]; len(v.Path) > 0 {
		err = fmt.Errorf("field %q of %s had unexpected path %s", ft.Name, ft.Affinity, v.Path)
	} else if a, e := dec.DecodeAssignment(ft.Affinity, v.Bytes); e != nil {
		err = fmt.Errorf("%w decoding field %q of %s", e, ft.Name, ft.Affinity)
	} else {
		ret = a
	}
	return
}

// note: this doesnt properly determine the default trait for an aspect
// weave works around this by providing the correct default value in the db
func zeroAssignment(ft rt.Field, idx int) (ret rt.Assignment, err error) {
	if init := ft.Init; init != nil {
		ret = init
	} else if v, e := rt.ZeroField(ft.Affinity, ft.Type, idx); e != nil {
		err = e
	} else {
		ret = rt.AssignValue(v)
	}
	return
}
