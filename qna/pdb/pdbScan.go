package pdb

import (
	"database/sql"

	"git.sr.ht/~ionous/iffy/tables"
)

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
