package main

import (
	"flag"
	"io"
	"log"

	"git.sr.ht/~ionous/tapestry/support/player"
	"github.com/ionous/errutil"
)

// go run play.go -in /Users/ionous/Documents/Tapestry/build/play.db -scene kitchenette
func main() {
	var inFile, testString, domain string
	var debugging bool
	flag.StringVar(&inFile, "in", "", "input file name (sqlite3)")
	flag.StringVar(&domain, "scene", "tapestry", "scene to start playing")
	flag.StringVar(&testString, "test", "", "optional list of commands to run (non-interactive)")
	flag.BoolVar(&debugging, "debug", false, "extra debugging output?")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.Parse()
	if !debugging {
		log.SetOutput(io.Discard)
	}
	if cnt, e := player.PlayGame(inFile, testString, domain); e != nil {
		errutil.PrintErrors(e, func(s string) { log.Println(s) })
		if errutil.Panic {
			log.Panic("mismatched")
		}
	} else {
		log.Println("done", cnt, inFile)
	}
}
