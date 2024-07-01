package cmdcheck

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt/print"
	"git.sr.ht/~ionous/tapestry/support/play"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/web/markup"
	"github.com/ionous/errutil"
)

// open db, select tests, de-gob and run them each in turn.
// print the results, only error on critical errors
func CheckFile(inFile, testName string, opt qna.Options) (ret int, err error) {
	if db, e := tables.CreateRunTime(inFile); e != nil {
		err = e
	} else {
		defer db.Close()
		ret, err = checkAll(db, testName, opt, tapestry.AllSignatures)
	}
	return
}

func checkAll(db *sql.DB, actuallyJustThisOne string, options qna.Options, signatures []map[uint64]typeinfo.Instance) (ret int, err error) {
	d := query.NewDecoder(signatures)
	if q, e := qdb.NewQueries(db, d); e != nil {
		err = e
	} else if grammar, e := qdb.MakeGrammar(db); e != nil {
		err = e
	} else if checks, e := q.ReadChecks(actuallyJustThisOne); e != nil {
		err = e
	} else if len(checks) == 0 {
		err = errutil.New("no matching checks found")
	} else {
		for _, check := range checks {
			if strings.HasPrefix(check.Name, "x ") || len(check.Test) == 0 {
				log.Println("ignoring", check.Name)
			} else {
				log.Printf("-- Checking: %q\n", check.Name)
				run := qna.NewRuntimeOptions(q, options)
				run.SetWriter(print.NewLineSentences(markup.ToText(os.Stdout)))
				survey := play.MakeDefaultSurveyor(run)
				play := play.NewPlaytime(run, survey, grammar)
				if e := checkOne(play, check, &ret); e != nil {
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

func checkOne(play *play.Playtime, check query.CheckData, pret *int) (err error) {
	debug.Stepper = func(words string) (err error) {
		// FIX: errors for step are getting fmt.Println in playTime.go
		// so expect output can't test for errors ( and on error looks a bit borken )
		_, err = play.Step(words)
		return
	}
	err = RunTest(play, check)
	(*pret)++
	return
}
