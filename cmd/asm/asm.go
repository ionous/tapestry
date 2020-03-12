// Package main for 'asm'.
// Generates a model database from ephemera data.
package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/assembly"
	"github.com/ionous/iffy/tables"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var inFile, outFile string
	flag.StringVar(&inFile, "in", "", "input file name (sqlite3)")
	flag.StringVar(&outFile, "out", "", "output file name (sqlite3)")
	flag.Parse()
	if e := assemble(outFile, inFile); e != nil {
		log.Fatalln(e)
	}
}

func assemble(outFile, inFile string) (err error) {
	if inFile, e := filepath.Abs(inFile); e != nil {
		err = e
	} else if outFile, e := filepath.Abs(outFile); e != nil {
		err = e
	} else if db, e := sql.Open("sqlite3", outFile); e != nil {
		err = errutil.New("couldn't create output file", outFile, e)
	} else {
		defer db.Close()
		//
		if e := tables.CreateModel(db); e != nil {
			err = e // create this in our output db
		} else if e := tables.CreateAssembly(db); e != nil {
			err = e // assembly are temporary tables used for computing the model
		} else if e := func() (err error) {
			// stat fails if there's no such file :(
			ai, _ := os.Stat(inFile)
			bi, _ := os.Stat(outFile)
			if !os.SameFile(ai, bi) {
				s := "attach database '" + inFile + "' as indb;"
				_, err = db.Exec(s)
			}
			return
		}(); e != nil {
			err = errutil.New("error attaching", e, inFile)
		} else {
			w := assembly.NewModeler(db)

			if e := assembly.DetermineAncestry(w, db, "things"); e != nil {
				err = e
			} else if e := assembly.DetermineFields(w, db); e != nil {
				err = e
			} else if e := assembly.DetermineAspects(w, db); e != nil {
				err = e
			} else if _, e := db.Exec("insert into mdl_prog select * from eph_prog;" +
				"insert into mdl_check select * from eph_check"); e != nil {
				err = e
			}
			// [-] adds relations between kinds
			// [-] creates instances
			// [-] sets instance properties
			// [-] relates instances
			// [] makes action handlers
			// [] makes event listeners
			// [] computes aliases
			// [] sets up printed name property
			// - backtracing to source:
			// ex. each "important" table entry gets an separate entry pointing back to original source
		}
	}
	return
}
