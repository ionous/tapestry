package qdb

import (
	"database/sql"
	"strings"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/dot"
	"git.sr.ht/~ionous/tapestry/tables"
)

type KindDecoder interface {
	CommandDecoder
	rt.Kinds
}

type SparseValue struct {
	Path  string
	Bytes []byte
}

func ScanSparseValues(rows *sql.Rows) (ret []SparseValue, err error) {
	var v SparseValue
	err = tables.ScanAll(rows, func() (_ error) {
		ret = append(ret, v)
		return
	}, &v.Path, &v.Bytes)
	return
}

func DecodeSparseRecord(kd KindDecoder, kind string, vs []SparseValue) (ret *rt.Record, err error) {
	if k, e := kd.GetKindByName(kind); e != nil {
		err = e
	} else {
		rr := recReader{kd, rt.NewRecord(k)}
		if e := rr.readRecord(vs); e != nil {
			err = e
		} else {
			ret = rr.rec
		}
	}
	return
}

type recordValue struct {
	rec *rt.Record
}

type recReader struct {
	kd  KindDecoder
	rec *rt.Record // target being built
}

// autocreates default sub records if need be.
func (rr *recReader) readRecord(pvs []SparseValue) (err error) {
	for _, pv := range pvs {
		if e := rr.readRecordPart(pv.Path, pv.Bytes); e != nil {
			err = e
			break
		}
	}
	return
}

func (rr *recReader) readRecordPart(dots string, bytes []byte) (err error) {
	pos := dot.MakeValueCursor(rr.kd, rt.RecordOf(rr.rec))
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
		// access the container of the inner value manually
		rec := pos.CurrentValue().Record()
		i := rec.FieldIndex(path[0])
		ft := rec.Field(i)
		if l, e := rr.kd.DecodeField(ft.Affinity, bytes, ft.Type); e != nil {
			err = e // fix: record field type?!
		} else if v, e := l.GetLiteralValue(rr.kd); e != nil {
			err = e
		} else {
			err = rec.SetIndexedField(i, v)
		}
	}
	return
}
