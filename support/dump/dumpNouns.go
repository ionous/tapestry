package dump

import (
	"database/sql"
	"fmt"

	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/tables"
)

func QueryNames(db *sql.DB, scene string) (ret []NounName, err error) {
	var n NounName
	var last string
	if rows, e := db.Query(must("names"), scene); e != nil {
		err = fmt.Errorf("%w while querying names", e)
	} else {
		err = tables.ScanAll(rows, func() (err error) {
			if last != n.Name {
				ret = append(ret, n)
				last = n.Name
			}
			return
		}, &n.Id, &n.Name)
	}
	return
}

func QueryNouns(db *sql.DB, scene string) (ret []NounData, err error) {
	if ns, e := QueryInnerNouns(db, scene); e != nil {
		err = fmt.Errorf("%w while querying ids", e)
	} else if e := QueryValues(db, ns); e != nil {
		err = fmt.Errorf("%w while querying values", e)
	} else {
		ret = ns
	}
	return
}

func QueryInnerNouns(db *sql.DB, scene string) (ret []NounData, err error) {
	var n NounData
	if rows, e := db.Query(must("nouns"), scene); e != nil {
		err = e
	} else {
		err = tables.ScanAll(rows, func() (err error) {
			ret = append(ret, n)
			return
		}, &n.Id, &n.Domain, &n.Noun, &n.Kind)
	}
	return
}

func QueryValues(db *sql.DB, ns []NounData) (err error) {
	q := must("values")
	for i, n := range ns {
		if rows, e := db.Query(q, n.Id); e != nil {
			err = e
		} else if vs, e := queryValues(rows); e != nil {
			err = e
		} else {
			ns[i].Values = vs
		}
	}
	return
}

func queryValues(rows *sql.Rows) (ret []query.ValueData, err error) {
	var last string
	var v query.ValueData
	err = tables.ScanAll(rows, func() (_ error) {
		if last != v.Field {
			ret = append(ret, v)
			last = v.Field
		}
		return
	}, &v.Field, &v.Path, &v.Value)
	return
}
