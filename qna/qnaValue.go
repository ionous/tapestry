package qna

import (
	"database/sql"

	"git.sr.ht/~ionous/tapestry/tables"
)

type writeCb func(domain, noun, field string, value any) error

func writeValues(db *sql.DB, call func(writeCb) error) (err error) {
	if tx, e := db.Begin(); e != nil {
		err = e
	} else {
		if q, e := tx.Prepare(runValue); e != nil {
			err = e
		} else {
			cb := func(domain, noun, field string, value any) error {
				_, e := q.Exec(domain, noun, field, value)
				return e
			}
			err = call(cb)
			q.Close()
		}
		tx.Commit()
	}
	return
}

var runValue = tables.Insert("run_value", "domain", "noun", "field", "value")
