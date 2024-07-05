package qdb

import (
	"database/sql"
	"encoding"

	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/tables"
)

func (q *Query) writeValues(d query.CacheMap) (err error) {
	if tx, e := q.db.Begin(); e != nil {
		err = e
	} else {
		if s, e := tx.Prepare(runValue); e != nil {
			err = e
		} else if p, e := q.rand.Save(); e != nil {
			err = e
		} else if e := writeRandomizeer(s, p); e != nil {
			err = e
		} else if e := writeDynamicValues(s, d); e != nil {
			err = e
		}
		tx.Commit()
	}
	return
}

var runValue = tables.Insert("run_value", "domain", "noun", "field", "value")

// read noun data from the db ( rt.run_value )
// and store them as bytes; run.readNounValue will unpack them
func readDynamicValues(db *sql.DB) (ret query.CacheMap, err error) {
	out := make(query.CacheMap)
	if rows, e := db.Query(
		`select domain, noun, field, value from run_value
		where domain != ''`); e != nil {
		err = e
	} else {
		var key query.Key
		var data []byte
		if e := tables.ScanAll(rows, func() (_ error) {
			out[key] = data
			return
		}, &key.Domain, &key.Target, &key.Field, &data); e != nil {
			err = e
		} else {
			ret = out
		}
	}
	return
}

// write all dynamic values to the database using the prepared 'runValue' statement.
func writeDynamicValues(s *sql.Stmt, data query.CacheMap) (err error) {
	for key, cached := range data {
		if userVal, ok := cached.(encoding.TextMarshaler); ok {
			if b, e := userVal.MarshalText(); e != nil {
				err = e
				break
			} else if _, e := s.Exec(key.Domain, key.Target, key.Field, b); e != nil {
				err = e
				break
			}
		}
	}
	return
}

const randomizerKey = "$randomizer"

func writeRandomizeer(s *sql.Stmt, src query.RandomPersist) (err error) {
	if _, e := s.Exec("", randomizerKey, "seed", src.Seed); e != nil {
		err = e
	} else if _, e := s.Exec("", randomizerKey, "last", src.LastRandom); e != nil {
		err = e
	}
	return
}

func readRandomizer(db *sql.DB) (ret query.RandomPersist, err error) {
	var b []byte
	var last int64
	if e := db.QueryRow(
		`select value from run_value 
		where noun = $1 and field = $2`,
		randomizerKey, "seed").Scan(&b); e != nil {
		err = e
	} else if e := db.QueryRow(
		`select value from run_value 
		where noun = $1 and field = $2`,
		randomizerKey, "last").Scan(&last); e != nil {
		err = e
	} else {
		ret = query.RandomPersist{Seed: b, LastRandom: last}
	}
	return
}
