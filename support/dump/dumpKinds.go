package dump

import (
	"database/sql"
	"fmt"

	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/qna/raw"
	"git.sr.ht/~ionous/tapestry/tables"
)

func QueryKinds(db *sql.DB, scene string) (ret []raw.KindData, err error) {
	if ks, e := QueryInnerKinds(db, scene); e != nil {
		err = fmt.Errorf("%w while querying inner kinds", e)
	} else if e := QueryAncestors(db, ks); e != nil {
		err = fmt.Errorf("%w while querying ancestors", e)
	} else if e := QueryFields(db, ks); e != nil {
		err = fmt.Errorf("%w while querying fields", e)
	} else {
		ret = ks
	}
	return
}

// build a list of kinds active in this scene
// only includes the most basic data for the kinds: name and id.
func QueryInnerKinds(db *sql.DB, scene string) (ret []raw.KindData, err error) {
	var k raw.KindData
	if rows, e := db.Query(must("kinds"), scene); e != nil {
		err = e
	} else {
		err = tables.ScanAll(rows, func() (err error) {
			ret = append(ret, k)
			return
		}, &k.Id, &k.Domain, &k.Kind)
	}
	return
}

// get ancestors for every kind in ks
func QueryAncestors(db *sql.DB, ks []raw.KindData) (err error) {
	q := must("ancestors")
	for i, k := range ks {
		if ps, e := tables.QueryStrings(db, q, k.Id); e != nil {
			err = e
			break
		} else {
			ks[i].Ancestors = ps
		}
	}
	return
}

// get ancestors for every kind in ks
func QueryFields(db *sql.DB, ks []raw.KindData) (err error) {
	q := must("fields")
	for i, k := range ks {
		if rows, e := db.Query(q, k.Id); e != nil {
			err = e
		} else if fs, e := queryFields(rows); e != nil {
			err = e
		} else {
			ks[i].Fields = fs
		}
	}
	return
}

func queryFields(rows *sql.Rows) (ret []query.FieldData, err error) {
	var last string
	var field query.FieldData
	err = tables.ScanAll(rows, func() (_ error) {
		// the same field might be listed twice:
		// the final value ( listed first )
		// and a non final value ( listed second )
		if last != field.Name {
			ret = append(ret, field)
			last = field.Name
		}
		return
	}, &field.Name, &field.Affinity, &field.Class, &field.Init)
	return
}
