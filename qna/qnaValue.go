package qna

import (
	"database/sql"

	"git.sr.ht/~ionous/tapestry/tables"
)

// type cachedField struct {
// 	domain string
// 	noun   string
// 	field  string
// 	value  any
// 	stored bool
// }

// type storedValues struct {
// 	counter counters
// 	value   cacheMap
// }

// func readValues(db *sql.DB) (ret storedValues, err error) {
// 	out := storedValues{
// 		counter: make(counters),
// 		value:   make(cacheMap),
// 	}

// 	if rows, e := db.Query(
// 		"select domain, noun, field, value from run_value"); e != nil {
// 		err = e
// 	} else {
// 		var domain, noun, field string
// 		var value any
// 		err = tables.ScanAll(rows, func() (err error) {
// 			if noun == meta.Counter {
// 				out.counter[field] = value.(int)
// 			} else {

//					m[field] = value
//				}
//				return
//			}, &domain, &noun, &field, &value)
//		}
//		if err == nil {
//			ret = out
//		}
//		return
//	}
type writeCb func(domain, noun, field, value any) error

func writeValues(db *sql.DB, call func(writeCb) error) (err error) {
	if tx, e := db.Begin(); e != nil {
		err = e
	} else {
		if q, e := tx.Prepare(runValue); e != nil {
			err = e
		} else {
			cb := func(domain, noun, field, value any) error {
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
