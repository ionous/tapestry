package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"git.sr.ht/~ionous/iffy"
	play "git.sr.ht/~ionous/iffy/cmd/play/internal"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

// go run play.go -in  /Users/ionous/Documents/Iffy/scratch/shared/play.db
func main() {
	var inFile, testString string
	flag.StringVar(&inFile, "in", "", "input file name (sqlite3)")
	flag.StringVar(&testString, "test", "", "input test string")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.Parse()
	if cnt, e := playGame(inFile, testString); e != nil {
		errutil.PrintErrors(e, func(s string) { log.Println(s) })
		if errutil.Panic {
			log.Panic("mismatched")
		}
	} else {
		log.Println("done", cnt, inFile)
	}
}

// open db, select tests, de-gob and run them each in turn.
// print the results, only error on critical errors
func playGame(inFile, testString string) (ret int, err error) {
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
			run := play.NewPlaytime(db, "#entire_game::kitchen")
			parser := play.NewParser(run, nil)
			//
			if len(testString) > 0 {
				step(parser, testString)
			} else {
				reader := bufio.NewReader(os.Stdin)
				for {
					fmt.Printf("> ")
					if in, _ := reader.ReadString('\n'); len(in) <= 1 {
						break
					} else {
						words := in[:len(in)-1]
						fmt.Println(words)
						step(parser, words)
					}
				}
			}
		}
	}
	return
}

func step(p *play.Parser, s string) {
	if res, e := p.Step(s); e != nil {
		fmt.Println(e)
	} else if res != nil {
		fmt.Println()
	}
}

func init() {
	iffy.RegisterGobs()
}
