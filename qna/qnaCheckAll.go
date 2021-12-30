package qna

import (
	"database/sql"

	"github.com/ionous/errutil"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/dl/story"
	"git.sr.ht/~ionous/iffy/qna/qdb"
	"git.sr.ht/~ionous/iffy/rt"
)

// CheckAll tests stored in the passed db.
// It logs the results of running the checks, and only returns error on critical errors.
func CheckAll(db *sql.DB, actuallyJustThisOne string, options Options, signatures []map[uint64]interface{}) (ret int, err error) {
	if qdb, e := qdb.NewQueries(db, false); e != nil {
		err = e
	} else if checks, e := qdb.ReadChecks(actuallyJustThisOne); e != nil {
		err = e
	} else if len(checks) == 0 {
		err = errutil.New("no matching checks found")
	} else {
		for _, el := range checks {
			var act rt.Execute
			if e := story.Decode(rt.Execute_Slot{&act}, el.Prog, signatures); e != nil {
				err = e
			} else if v, e := literal.ReadLiteral(el.Aff, "", el.Value); e != nil {
				err = e
			} else if expect, ok := v.(*literal.TextValue); !ok {
				err = errutil.New("can only handle text values right now")
			} else {
				// fix? its currently necessary to activate a global domain, rather than jump straight into the check domain.
				// something about pair activation goes a bit wonky: multiple pairs can become active at once.
				qdb.ResetSavedData()
				run := NewRuntimeOptions(qdb, options, iffy.AllSignatures)
				if _, e := run.ActivateDomain("entire_game"); e != nil {
					err = e
				} else {
					t := CheckOutput{
						Name:   el.Name,
						Domain: el.Domain,
						Expect: expect.String(),
						Test:   act,
					}
					if e := t.RunTest(run); e != nil {
						err = errutil.Append(err, e)
					}
					ret++
				}
			}
		}
	}
	return
}
