package qdb

import (
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/tables"
)

func (q *Query) Random(inclusiveMin, exclusiveMax int) int {
	return q.rand.Random(inclusiveMin, exclusiveMax)
}

// read noun data from the db ( rt.run_value )
// and store them as bytes; run.readNounValue will unpack them
func (q *Query) LoadGame(outPath string) (ret query.CacheMap, err error) {
	if e := tables.LoadFile(q.db, outPath); e != nil {
		err = e
	} else if data, e := readDynamicValues(q.db); e != nil {
		err = e
	} else if p, e := readRandomizer(q.db); e != nil {
		err = e
	} else if e := q.rand.Load(p); e != nil {
		err = e
	} else {
		ret = data
	}
	return
}

func (q *Query) SaveGame(outPath string, d query.CacheMap) (err error) {
	if e := q.writeValues(d); e != nil {
		err = e
	} else {
		err = tables.SaveFile(outPath, false, q.db)
	}
	return
}
