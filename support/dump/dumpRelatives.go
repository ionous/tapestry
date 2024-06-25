package dump

import (
	"database/sql"
	"fmt"

	"git.sr.ht/~ionous/tapestry/qna/raw"
	"git.sr.ht/~ionous/tapestry/tables"
)

func QueryRelatives(db *sql.DB, scene string) (ret []raw.RelativeData, err error) {
	if rs, e := QueryInnerRelations(db, scene); e != nil {
		err = fmt.Errorf("%w while querying relations", e)
	} else if e := QueryPairs(db, scene, rs); e != nil {
		err = fmt.Errorf("%w while querying pairs", e)
	} else {
		ret = rs
	}
	return
}

// build a list of relations
func QueryInnerRelations(db *sql.DB, scene string) (ret []raw.RelativeData, err error) {
	var rel raw.RelativeData
	if rows, e := db.Query(must("relations"), scene); e != nil {
		err = e
	} else {
		err = tables.ScanAll(rows, func() (err error) {
			ret = append(ret, rel)
			return
		}, &rel.Id, &rel.Relation, &rel.OneKind, &rel.OtherKind, &rel.Cardinality)
	}
	return
}

// build a list of pairs per relation
func QueryPairs(db *sql.DB, scene string, rs []raw.RelativeData) (err error) {
	q := must("pairs")
	for i, rel := range rs {
		if rows, e := db.Query(q, scene, rel.Id); e != nil {
			err = e
		} else if ps, e := queryPairs(rows); e != nil {
			err = e
		} else {
			rs[i].Pairs = ps
		}
	}
	return
}

func queryPairs(rows *sql.Rows) (ret []raw.Pair, err error) {
	var p raw.Pair
	err = tables.ScanAll(rows, func() (_ error) {
		ret = append(ret, p)
		return
	}, &p.One, &p.Other)
	return
}
