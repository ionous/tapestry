package dump

import (
	"database/sql"
	"fmt"

	"git.sr.ht/~ionous/tapestry/qna/raw"
	"git.sr.ht/~ionous/tapestry/tables"
)

func QueryNames(db *sql.DB, scene string) (ret []raw.NounName, err error) {
	var n raw.NounName
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
		}, &n.Name, &n.Noun)
	}
	return
}
