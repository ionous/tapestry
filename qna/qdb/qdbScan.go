package qdb

import (
	"database/sql"

	"git.sr.ht/~ionous/tapestry/tables"
)

func scanId(q *sql.Stmt, args ...interface{}) (ret int, err error) {
	if e := q.QueryRow(args...).Scan(&ret); e != nil && e != sql.ErrNoRows {
		err = e
	}
	return
}

func scanString(q *sql.Stmt, args ...interface{}) (ret string, err error) {
	if e := q.QueryRow(args...).Scan(&ret); e != nil && e != sql.ErrNoRows {
		err = e
	}
	return
}

func scanStrings(q *sql.Stmt, args ...interface{}) (ret []string, err error) {
	if rows, e := q.Query(args...); e != nil {
		err = e
	} else {
		var one string
		err = tables.ScanAll(rows, func() (err error) {
			ret = append(ret, one)
			return
		}, &one)
	}
	return
}
