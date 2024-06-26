package qdb

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/qna/decoder"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/tables"
)

// read all the matching tests from the db.
func (q *Query) ReadChecks(actuallyJustThisOne string) (ret []query.CheckData, err error) {
	if rows, e := q.checks.Query(actuallyJustThisOne); e != nil {
		err = e
	} else {
		var check struct {
			Name   string
			Domain string
			Aff    affine.Affinity
			Prog   []byte // a serialized rt.Execute_Slice
			Value  []byte // a serialized literal.Value
		}
		err = tables.ScanAll(rows, func() (err error) {
			if act, e := q.dec.DecodeProg(check.Prog); e != nil {
				err = e
			} else if expect, e := readLegacyExpectation(q.dec, check.Value, check.Aff); e != nil {
				err = e
			} else {
				ret = append(ret, query.CheckData{
					Name:   check.Name,
					Domain: check.Domain,
					Expect: expect,
					Test:   act,
				})
			}
			return
		}, &check.Name, &check.Domain, &check.Value, &check.Aff, &check.Prog)
	}
	return
}

func readLegacyExpectation(dec decoder.Decoder, b []byte, aff affine.Affinity) (ret string, err error) {
	if len(b) > 0 {
		if v, e := dec.DecodeField(aff, b, ""); e != nil {
			err = e
		} else if expect, ok := v.(*literal.TextValue); !ok {
			err = errors.New("can only handle text values right now")
		} else {
			ret = expect.String()
		}
	}
	return
}
