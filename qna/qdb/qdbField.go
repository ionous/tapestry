package qdb

import (
	"database/sql"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/qna/decoder"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/tables"
)

// given an query producing rows of field name, affinity, class, and init
// build an rt.Field list
func ScanFields(rows *sql.Rows, dec decoder.Decoder) (ret []rt.Field, err error) {
	var last string
	var f struct {
		Name     string
		Affinity affine.Affinity
		Class    sql.NullString
		Init     []byte
	}
	err = tables.ScanAll(rows, func() (err error) {
		// the same field might be listed twice:
		// the final value ( listed first )
		// and a non final value ( listed second )
		if last != f.Name {
			if init, e := decodeInit(dec, f.Affinity, f.Init); e != nil {
				err = e
			} else {
				ret = append(ret, rt.Field{
					Name:     f.Name,
					Affinity: f.Affinity,
					Type:     f.Class.String,
					Init:     init,
				})
				last = f.Name
			}
		}
		return
	}, &f.Name, &f.Affinity, &f.Class.String, &f.Init)
	return
}

// decode the passed assignment, if it exists.
func decodeInit(d decoder.Decoder, aff affine.Affinity, b []byte) (ret rt.Assignment, err error) {
	if len(b) > 0 {
		ret, err = d.DecodeAssignment(aff, b)
	}
	return
}
