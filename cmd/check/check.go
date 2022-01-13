package main

import (
	"database/sql"
	"flag"
	"log"
	"path/filepath"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

// ex. go run check.go -in /Users/ionous/Documents/Tapestry/build/play.db
func main() {
	var inFile, testName string
	flag.StringVar(&inFile, "in", "", "input file name (.db)")
	flag.StringVar(&testName, "run", "", "optional specific test")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	responses := flag.Bool("responses", false, "print response names instead of values")
	flag.Parse()
	opt := qna.NewOptions()
	opt.SetOption(meta.PrintResponseNames, generic.BoolOf(*responses))
	if cnt, e := checkFile(inFile, lang.Underscore(testName), opt); e != nil {
		errutil.PrintErrors(e, func(s string) { log.Println(s) })
		if errutil.Panic {
			log.Panic("mismatched")
		}
	} else {
		log.Println("Checked", cnt, inFile)
	}
}

// open db, select tests, de-gob and run them each in turn.
// print the results, only error on critical errors
func checkFile(inFile, testName string, opt qna.Options) (ret int, err error) {
	if inFile, e := filepath.Abs(inFile); e != nil {
		err = e
	} else if db, e := sql.Open(tables.DefaultDriver, inFile); e != nil {
		err = errutil.New("couldn't open db", inFile, e)
	} else {
		defer db.Close()
		if e := tables.CreateRun(db); e != nil {
			err = e
		} else {
			ret, err = qna.CheckAll(db, testName, opt, tapestry.AllSignatures)
		}
	}
	return
}
