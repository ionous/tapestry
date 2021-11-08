package main

import (
	"database/sql"
	"flag"
	"log"
	"path/filepath"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/qna"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

// ex. go run check.go -in /Users/ionous/Documents/Iffy/scratch/shared/play.db
func main() {
	var inFile, testName string
	flag.StringVar(&inFile, "in", "", "input file name (sqlite3)")
	flag.StringVar(&testName, "run", "", "optional specific test ( in camelcase )")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.Parse()
	if cnt, e := checkFile(inFile, testName); e != nil {
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
func checkFile(inFile, testName string) (ret int, err error) {
	if inFile, e := filepath.Abs(inFile); e != nil {
		err = e
	} else if db, e := sql.Open(tables.DefaultDriver, inFile); e != nil {
		err = errutil.New("couldn't create output file", inFile, e)
	} else {
		defer db.Close()
		if e := tables.CreateRun(db); e != nil {
			err = e
		} else if e := tables.CreateRunViews(db); e != nil {
			err = e
		} else {
			ret, err = qna.CheckAll(db, testName, iffy.AllSignatures)
		}
	}
	return
}
