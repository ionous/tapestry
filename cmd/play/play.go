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

// ex. go run play.go -in /Users/ionous/Documents/Iffy/scratch/main#entire_game::kitchen"/play.db
func main() {
	var inFile string
	flag.StringVar(&inFile, "in", "", "input file name (sqlite3)")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.Parse()
	if cnt, e := playGame(inFile); e != nil {
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
func playGame(inFile string) (ret int, err error) {
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
			reader := bufio.NewReader(os.Stdin)
			for {
				fmt.Printf("> ")
				if in, _ := reader.ReadString('\n'); len(in) <= 1 {
					break
				} else {
					words := in[:len(in)-1]
					fmt.Println(words)
					if e := parser.Parse(words); e != nil {
						fmt.Println(e)
					}
				}
			}
		}
	}
	return
}

func init() {
	iffy.RegisterGobs()
}
