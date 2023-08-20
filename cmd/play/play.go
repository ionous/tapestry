package main

import (
	"flag"
	"io"
	"log"

	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/support/play"
	"github.com/ionous/errutil"
)

// go run play.go -in /Users/ionous/Documents/Tapestry/build/play.db -scene kitchenette
func main() {
	var inFile, testString, domain string
	var json, debugging bool
	flag.StringVar(&inFile, "in", "", "input file name (sqlite3)")
	flag.StringVar(&domain, "scene", "tapestry", "scene to start playing")
	flag.StringVar(&testString, "test", "", "optional list of commands to run (non-interactive)")
	flag.BoolVar(&json, "json", false, "expect input/output in json (default is plain text)")
	flag.BoolVar(&debugging, "debug", false, "extra debugging output?")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.Parse()
	if !debugging {
		log.SetOutput(io.Discard)
	}
	opt := qna.NewOptions()
	opt.SetOption(meta.JsonMode, generic.BoolOf(json))

	if cnt, e := play.PlayGame(inFile, testString, domain, json); e != nil {
		errutil.PrintErrors(e, func(s string) { log.Println(s) })
		if errutil.Panic {
			log.Panic("mismatched")
		}
	} else {
		log.Println("done", cnt, inFile)
	}
}
