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
	"git.sr.ht/~ionous/tapestry/web/markup"
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
				w := print.NewLineSentences(markup.ToText(os.Stdout))
				run := NewRuntimeOptions(w, qdb, options, tapestry.AllSignatures)
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
	return
}
