package cmdcheck

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strings"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/decode"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt/print"
	"git.sr.ht/~ionous/tapestry/support/play"
	"git.sr.ht/~ionous/tapestry/support/player"
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
	if query, e := qdb.NewQueries(db); e != nil {
		err = e
	} else if grammar, e := player.MakeGrammar(db); e != nil {
		err = e
	} else if checks, e := query.ReadChecks(actuallyJustThisOne); e != nil {
		err = e
	} else if len(checks) == 0 {
		err = errutil.New("no matching checks found")
	} else {
		for _, check := range checks {
			if strings.HasPrefix(check.Name, "x ") || len(check.Prog) == 0 {
				log.Println("ignoring", check.Name)
			} else {
				log.Printf("-- Checking: %q\n", check.Name)
				w := print.NewLineSentences(markup.ToText(os.Stdout))
				d := decode.NewDecoder(signatures)
				run := qna.NewRuntimeOptions(query, d, options)

				run.SetWriter(w)
				survey := play.MakeDefaultSurveyor(run)
				play := play.NewPlaytime(run, survey, grammar)
				if e := checkOne(d, play, check, &ret); e != nil {
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

func checkOne(d *decode.Decoder, play *play.Playtime, check query.CheckData, pret *int) (err error) {
	if act, e := d.DecodeProg(check.Prog); e != nil {
		err = e
	} else if expect, e := readLegacyExpectation(check); e != nil {
		err = e
	} else {
		t := CheckOutput{
			Name:   check.Name,
			Domain: check.Domain,
			Expect: expect,
			Test:   act,
		}
		debug.Stepper = func(words string) (err error) {
			// FIX: errors for step are getting fmt.Println in playTime.go
			// so expect output can't test for errors ( and on error looks a bit borken )
			_, err = play.Step(words)
			return
		}
		err = t.RunTest(play)
		(*pret)++
	}
	return
}

func readLegacyExpectation(check query.CheckData) (ret string, err error) {
	if len(check.Value) > 0 {
		var val any
		if e := json.Unmarshal(check.Value, &val); e != nil {
			err = e
		} else if v, e := literal.ReadLiteral(check.Aff, "", val); e != nil {
			err = e
		} else if expect, ok := v.(*literal.TextValue); !ok {
			err = errutil.New("can only handle text values right now")
		} else {
			ret = expect.String()
		}
	}
	return
}
