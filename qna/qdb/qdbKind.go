package qdb

import (
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
)

func (q *Query) GetKindByName(k string) (ret *rt.Kind, err error) {
	key := query.MakeKey("kinds", k, "")
	if c, e := q.constVals.Ensure(key, func() (ret any, err error) {
		if k, e := BuildKind(q, k); e != nil {
			err = e
		} else {
			ret = &k // put this into the heap
		}
		return
	}); e != nil {
		err = e
	} else {
		ret = c.(*rt.Kind)
	}
	return
}

// cached fields exclusive to a kind
func (q *Query) FieldsOf(kind string) (ret []rt.Field, err error) {
	key := query.MakeKey("fields", kind, "")
	if c, e := q.constVals.Ensure(key, func() (ret any, err error) {
		if rows, e := q.fieldsOf.Query(kind); e != nil {
			err = e
		} else {
			ret, err = ScanFields(rows, q.dec)
		}
		return
	}); e != nil {
		err = e
	} else {
		ret = c.([]rt.Field)
	}
	return
}
