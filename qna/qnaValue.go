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

// 				m[field] = value
// 			}
// 			return
// 		}, &domain, &noun, &field, &value)
// 	}
// 	if err == nil {
// 		ret = out
// 	}
// 	return
// }

// func writeCachedValues(db *sql.DB, c cacheMap) (err error) {
// 	return writeValues(db, func(q *sql.Stmt) (err error) {
// 		for _, entry := range c {
// 			if entry.e == nil {
// 				if v := (entry.v).(cachedField); v.stored {
// 					// all stored values are variant values
// 					// and, for now at least, rely on driver serialization to nounValues
// 					raw := (v.value).(g.Value).Any()
// 					if _, e := q.Exec(v.domain, v.noun, v.field, raw); e != nil {
// 						err = e
// 						break
// 					}
// 				}
// 			}
// 		}
// 		return
// 	})
// }

func writeValues(db *sql.DB, cb func(*sql.Stmt) error) (err error) {
	if tx, e := db.Begin(); e != nil {
		err = e
	} else {
		if q, e := tx.Prepare(runValue); e != nil {
			err = e
		} else {
			err = cb(q)
			q.Close()
		}
		tx.Commit()
	}
	return
}

var runValue = tables.Insert("run_value", "domain", "noun", "field", "value")
