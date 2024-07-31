package cmdcheck

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/print"
	"git.sr.ht/~ionous/tapestry/support/play"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/web/markup"
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

func checkAll(db *sql.DB, matching string, options qna.Options, signatures []map[uint64]typeinfo.Instance) (okay int, err error) {
	if len(matching) == 0 {
		matching = "test *"
	}
	d := query.NewDecoder(signatures)
	if q, e := qdb.NewQueries(db, d); e != nil {
		err = e
	} else if grammar, e := qdb.MakeGrammar(db); e != nil {
		err = e
	} else if checks, e := readChecks(db, matching); e != nil {
		err = e
	} else if len(checks) == 0 {
		err = fmt.Errorf("no checks matching %q found", matching)
	} else {
		var failed int
		for _, check := range checks {
			log.Printf("-- Checking: %q\n", check)
			run := qna.NewRuntimeOptions(q, options)
			// var buf strings.Builder
			run.SetWriter(print.NewLineSentences(markup.ToText(os.Stdout)))
			survey := play.MakeDefaultSurveyor(run)
			play := play.NewPlaytime(run, survey, grammar)
			//
			debug.Stepper = func(words string) (err error) {
				// FIX: errors for step are getting fmt.Println in playTime.go
				// so expect output can't test for errors ( and on error looks a bit borken )
				_, err = play.HandleTurn(words)
				return
			}
			if e := checkOne(play, check); e != nil {
				log.Printf("NG: %s failed because %s", check, e)
				failed++
			} else {
				okay++
				log.Printf("OK: %s\n", check)
			}
		}
		if failed > 0 {
			err = fmt.Errorf("failed %d tests", failed)
		}
	}
	return
}
func readChecks(db *sql.DB, matching string) (ret []string, err error) {
	return tables.QueryStrings(db, readScenes, matching)
}

func checkOne(run rt.Runtime, str string) (err error) {
	if e := run.ActivateDomain(str); e != nil && !wasQuit(e) {
		err = e
	} else if e := run.ActivateDomain(""); e != nil {
		err = fmt.Errorf("couldnt restore domain %v", e)
	}
	return
}

// fix: better might be an "Expect error:" that matches partial error return strings
// or "Expect quit/signal:" specifically for this case
func wasQuit(e error) bool {
	var sig rt.Signal // if the game was quit, override the error if output remains
	return errors.As(e, &sig) && sig == rt.SignalQuit
}

const readScenes = `select distinct domain
			from mdl_domain md
			where domain glob ?1
			order by domain`
