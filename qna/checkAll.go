package qna

import (
	"bytes"
	"database/sql"
	"encoding/gob"

	"github.com/ionous/errutil"

	"git.sr.ht/~ionous/iffy/dl/check"
	"git.sr.ht/~ionous/iffy/tables"
)

// CheckAll tests stored in the passed db.
// It logs the results of running the tests, and only returns error on critical errors.
func CheckAll(db *sql.DB, actuallyJustThisOne string) (ret int, err error) {
	var name string
	var prog []byte
	var tests []check.CheckOutput
	if e := tables.QueryAll(db,
		`select name, bytes 
		from mdl_prog pg 
		where type='CheckOutput'
		order by name`,
		func() (err error) {
			if len(actuallyJustThisOne) == 0 || actuallyJustThisOne == name {
				var curr check.CheckOutput
				dec := gob.NewDecoder(bytes.NewBuffer(prog))
				if e := dec.Decode(&curr); e != nil {
					err = e
				} else {
					tests = append(tests, curr)
				}
			}
			return
		}, &name, &prog); e != nil {
		err = errutil.New("query error", e)
	} else if len(tests) == 0 {
		err = errutil.New("no matching tests found")
	} else {
		// FIX: we have to cache the statements b/c we cant use them during QueryAll
		for _, t := range tests {
			run := NewRuntime(db)
			tables.Must(db, `delete from run_domain; delete from run_pair`)
			run.ActivateDomain("entire_game", true)
			if e := t.RunTest(run); e != nil {
				err = errutil.Append(err, e)
			}
			ret++
		}
	}
	return
}
