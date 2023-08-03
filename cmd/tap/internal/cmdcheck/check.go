package cmdcheck

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/decode"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/print"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/web/markup"
	"github.com/ionous/errutil"
)

// CheckAll tests stored in the passed db.
// It logs the results of running the checks, and only returns error on critical errors.
func CheckAll(db *sql.DB, actuallyJustThisOne string, options qna.Options, signatures []map[uint64]interface{}) (ret int, err error) {
	if e := tables.CreateRun(db); e != nil {
		err = e
	} else if qdb, e := qdb.NewQueries(db, true); e != nil {
		err = e
	} else if checks, e := qdb.ReadChecks(actuallyJustThisOne); e != nil {
		err = e
	} else if len(checks) == 0 {
		err = errutil.New("no matching checks found")
	} else {
		for _, check := range checks {
			if strings.HasPrefix(check.Name, "x ") || len(check.Prog) == 0 {
				log.Println("ignoring", check.Name)
			} else {
				log.Println("-- Checking:", check.Name, check.Domain)

				if e := checkOne(qdb, check, options, signatures, &ret); e != nil {
					e := errutil.New(e, "during", check.Name)
					err = errutil.Append(err, e)
					log.Println(e)
				} else {
					log.Printf("ok. test %s", check.Name)
				}
			}
		}
	}
	return
}

func checkOne(qdb *qdb.Query, check query.CheckData, options qna.Options, signatures []map[uint64]interface{}, pret *int) (err error) {
	var act rt.Execute_Slice
	if e := story.Decode(&act, check.Prog, signatures); e != nil {
		err = e
	} else if expect, e := readLegacyExpectation(check); e != nil {
		err = e
	} else {
		w := print.NewLineSentences(markup.ToText(os.Stdout))
		d := decode.NewDecoder(signatures)
		run := qna.NewRuntimeOptions(w, qdb, d, options)
		t := CheckOutput{
			Name:   check.Name,
			Domain: check.Domain,
			Expect: expect,
			Test:   act,
		}
		err = t.RunTest(run)
		(*pret)++
	}
	return
}

func readLegacyExpectation(check query.CheckData) (ret string, err error) {
	if len(check.Value) > 0 {
		if v, e := literal.ReadLiteral(check.Aff, "", check.Value); e != nil {
			err = e
		} else if expect, ok := v.(*literal.TextValue); !ok {
			err = errutil.New("can only handle text values right now")
		} else {
			ret = expect.String()
		}
	}
	return
}
