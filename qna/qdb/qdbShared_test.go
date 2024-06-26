package qdb

import (
	"database/sql"
	"fmt"

	"git.sr.ht/~ionous/tapestry/qna/decoder"
)

type QueryTest struct {
	*Query
	scan *sql.Stmt
}

func NewQueryTest(db *sql.DB) (ret *QueryTest, err error) {
	if q, e := NewQueries(db, decoder.DecodeNone("query test")); e != nil {
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
		err = fmt.Errorf("%w couldnt scan %q", e, name)
	} else {
		ret = els
	}
	return
}
