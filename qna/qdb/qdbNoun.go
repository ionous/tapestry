package qdb

import (
	"database/sql"
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/qna/decoder"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/dot"
	"git.sr.ht/~ionous/tapestry/tables"
)

// interpreting the value is left to the caller ( re: field affinity )
func (q *Query) NounValue(noun, field string) (ret rt.Assignment, err error) {
	if k, e := scanString(q.nounKind, noun); e != nil {
		err = e
	} else if k, e := q.GetKindByName(k); e != nil {
		err = e
	} else if i := k.FieldIndex(field); i < 0 {
		err = fmt.Errorf("couldnt find field %q in kind %q", field, k.Name())
	} else {
		ft := k.Field(i)
		if rows, e := q.nounValues.Query(noun, field); e != nil {
			err = e
		} else {
			var v pathValue
			var vals []pathValue
			if e := tables.ScanAll(rows, func() (_ error) {
				vals = append(vals, v)
				return
			}, &v.path, &v.bytes); e != nil {
				err = e
			} else {
				if len(vals) == 0 {
					ret, err = zeroAssignment(ft, i)
				} else if ft.Affinity != affine.Record {
					ret, err = q.decodeValue(ft, vals)
				} else {
					ret, err = q.decodeRecord(ft, vals)
				}
			}
		}
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

func (q *Query) decodeValue(ft rt.Field, vals []pathValue) (ret rt.Assignment, err error) {
	if cnt := len(vals); cnt != 1 {
		err = fmt.Errorf("field %q of %s expected a single value, got %d", ft.Name, ft.Affinity, cnt)
	} else if v := vals[0]; len(v.path.String) > 0 {
		err = fmt.Errorf("field %q of %s had unexpected path %s", ft.Name, ft.Affinity, v.path.String)
	} else if a, e := q.dec.DecodeAssignment(ft.Affinity, v.bytes); e != nil {
		err = fmt.Errorf("%w decoding field %q of %s", e, ft.Name, ft.Affinity)
	} else {
		ret = a
	}
	return
}

func (q *Query) decodeRecord(ft rt.Field, vals []pathValue) (ret rt.Assignment, err error) {
	if k, e := q.GetKindByName(ft.Type); e != nil {
		err = e
	} else {
		rr := recReader{q, q.dec, rt.NewRecord(k)}
		if e := rr.readRecord(vals); e != nil {
			err = e
		} else {
			out := rt.RecordOf(rr.rec)
			ret = rt.AssignValue(out)
		}
	}
	return
}

type recordValue struct {
	rec *rt.Record
}

type recReader struct {
	ks  rt.Kinds
	dec decoder.Decoder
	rec *rt.Record // target being built
}

// autocreates default sub records if need be.
func (rr *recReader) readRecord(pvs []pathValue) (err error) {
	for _, pv := range pvs {
		if e := rr.readRecordPart(pv.path.String, pv.bytes); e != nil {
			err = e
			break
		}
	}
	return
}

func (rr *recReader) readRecordPart(dots string, bytes []byte) (err error) {
	pos := dot.MakeValueCursor(rr.ks, rt.RecordOf(rr.rec))
	// follow the dots
	path := strings.Split(dots, ".")
	for len(path) > 1 {
		part, rest := path[0], path[1:]
		if next, e := pos.GetAtField(part); e != nil {
			err = e
			break
		} else {
			pos, path = next, rest
		}
	}
	if err == nil {
		// access the container of the last value manually
		rec := pos.CurrentValue().Record()
		i := rec.FieldIndex(path[0])
		ft := rec.Field(i)
		if l, e := rr.dec.DecodeField(ft.Affinity, bytes, ft.Type); e != nil {
			err = e // fix: record field type?!
		} else if v, e := l.GetLiteralValue(rr.ks); e != nil {
			err = e
		} else {
			err = rec.SetIndexedField(i, v)
		}
	}
	return
}

type pathValue struct {
	path  sql.NullString
	bytes []byte
}
