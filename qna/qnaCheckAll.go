package qna

import (
	"database/sql"
	"os"

	"github.com/ionous/errutil"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/print"
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
			var act rt.Execute_Slice
			if e := story.Decode(&act, el.Prog, signatures); e != nil {
				err = errutil.Append(err, e)
			} else if v, e := literal.ReadLiteral(el.Aff, "", el.Value); e != nil {
				err = errutil.Append(err, e)
			} else if expect, ok := v.(*literal.TextValue); !ok {
				e := errutil.New("can only handle text values right now")
				err = errutil.Append(err, e)
			} else {
				// fix? its necessary to reset the domain right now.
				// not entirely clear why to me.
				if _, e := qdb.ActivateDomain(""); e != nil {
					err = errutil.Append(err, e)
				} else {
					w := print.NewAutoWriter(os.Stdout)
					run := NewRuntimeOptions(w, qdb, options, tapestry.AllSignatures)
					// fix! if we dont activate "entire_game" first, we wind up with multiple pairs active
					// this is something to do with the way the pair query works
					// when there is a relation in the entire_game that is supposed to be changed by a sub-domain.
					if _, e := run.ActivateDomain("entire_game"); e != nil {
						err = errutil.Append(err, e)
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
	}
	return
}
