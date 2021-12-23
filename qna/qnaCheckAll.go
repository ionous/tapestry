package qna

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/ionous/errutil"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/jsn/cin"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/tables"
)

// CheckAll tests stored in the passed db.
// It logs the results of running the tests, and only returns error on critical errors.
func CheckAll(db *sql.DB, actuallyJustThisOne string, signatures []map[uint64]interface{}) (ret int, err error) {
	var name, domain string
	var aff affine.Affinity
	var prog json.RawMessage
	var value interface{}
	//
	if len(actuallyJustThisOne) > 0 {
		actuallyJustThisOne += ";"
	}
	//
	var tests []CheckOutput
	if e := tables.QueryAll(db,
		`select mc.name, md.domain, mc.value, mc.affinity, mc.prog
		from mdl_check mc
		join mdl_domain md
			on (mc.domain=md.rowid) 
		order by mc.domain, mc.name`,
		func() (err error) {
			if len(actuallyJustThisOne) == 0 || strings.Contains(actuallyJustThisOne, name+";") {
				var act rt.Execute
				if e := cin.Decode(rt.Execute_Slot{&act}, prog, signatures); e != nil {
					err = e
				} else if str, ok := value.(string); !ok || aff != affine.Text {
					err = errutil.New("tests only compare text right now")
				} else {
					tests = append(tests, CheckOutput{
						Name:   name,
						Expect: str,
						Test:   act,
					})
				}
			}
			return
		}, &name, &domain, &value, &aff, &prog); e != nil {
		err = errutil.New("query error", e)
	} else if len(tests) == 0 {
		err = errutil.New("no matching tests found")
	} else {
		// FIX: we have to cache the statements b/c we cant use them during QueryAll
		for _, t := range tests {
			run := NewRuntime(db, iffy.AllSignatures)
			tables.Must(db, `delete from run_domain; delete from run_pair`)
			//
			if _, e := run.ActivateDomain(domain); e != nil {
				err = errutil.Append(err, e)
			} else if e := t.RunTest(run); e != nil {
				err = errutil.Append(err, e)
			}
			ret++
		}
	}
	return
}
