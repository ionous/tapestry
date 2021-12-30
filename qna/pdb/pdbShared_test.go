package pdb

import (
	"database/sql"

	"github.com/ionous/errutil"
)

type QueryTest struct {
	*Query
	scan *sql.Stmt
}

func NewQueryTest(db *sql.DB) (ret *QueryTest, err error) {
	if q, e := NewQueries(db, false); e != nil {
		err = e
	} else if scan, e := db.Prepare(
		`select md.domain || ':' || rd.active
		from run_domain rd
		join mdl_domain md
			on (rd.domain=md.rowid)
		where rd.active > 0
		`); e != nil {
		err = e
	} else {
		ret = &QueryTest{q, scan}
	}
	return
}

// For testing: activate without deleting
func (q *QueryTest) InnerActivate(name string, act int) (ret []string, err error) {
	if _, e := q.domainActivation.Exec(name, act); e != nil {
		err = e
	} else if els, e := scanStrings(q.scan); e != nil {
		err = errutil.New("couldnt scan", name, e)
	} else {
		ret = els
	}
	return
}
