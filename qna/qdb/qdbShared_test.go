package qdb

import (
	"database/sql"

	"github.com/ionous/errutil"
)

type QueryTest struct {
	*Query
	scan *sql.Stmt
}

func NewQueryTest(db *sql.DB) (ret *QueryTest, err error) {
	if q, e := NewQueries(db); e != nil {
		err = e
	} else if scan, e := db.Prepare(
		`select domain || ':' || active
		from run_domain
		where active > 0
		`); e != nil {
		err = e
	} else {
		ret = &QueryTest{q, scan}
	}
	return
}

// For testing: returns a list of active domains and their activation count.
func (q *QueryTest) InnerActivate(name string) (ret []string, err error) {
	if _, _, e := q.ActivateDomains(name); e != nil {
		err = e
	} else if els, e := scanStrings(q.scan); e != nil {
		err = errutil.New("couldnt scan", name, e)
	} else {
		ret = els
	}
	return
}
