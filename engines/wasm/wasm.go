// Experimental wasm version
// use `npm run wasm` to build
// or GOOS=js GOARCH=wasm go build -o main.wasm wasm.go
package main

import (
	"fmt"
	"os"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/decode"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt/print"
	"git.sr.ht/~ionous/tapestry/web/markup"
)

// go run play.go -in /Users/ionous/Documents/Tapestry/build/play.db -scene kitchenette
func main() {
	// var inFile, testString, domain string
	// var debugging bool
	// flag.StringVar(&inFile, "in", "", "input file name (sqlite3)")
	// flag.StringVar(&domain, "scene", "tapestry", "scene to start playing")
	// flag.StringVar(&testString, "test", "", "optional list of commands to run (non-interactive)")
	// flag.BoolVar(&debugging, "debug", false, "extra debugging output?")
	// flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	// flag.Parse()
	//
	//	if !debugging {
	//		log.SetOutput(io.Discard)
	//	}
	//
	//	if e := player.PlayGame(inFile, testString, domain); e != nil {
	//		errutil.PrintErrors(e, func(s string) { log.Println(s) })
	//		if errutil.Panic {
	//			log.Panic("mismatched")
	//		}
	//	}
	fmt.Println("hello")
}

func PlayGame(scene string) (err error) {
	opts := qna.NewOptions()
	return PlayWithOptions("xx", opts)
}

func PlayWithOptions(scene string, opts qna.Options) (err error) {
	const prompt = "> "

	var q query.QueryNone
	d := decode.NewDecoder(tapestry.AllSignatures)
	run := qna.NewRuntimeOptions(q, d, opts)
	w := print.NewLineSentences(markup.ToText(os.Stdout))
	run.SetWriter(w)
	if e := run.ActivateDomain(scene); e != nil {
		err = e
	} else {
		// survey := play.MakeDefaultSurveyor(run)
		// play := play.NewPlaytime(run, survey, grammar)

	}

	// if grammar, e := play.MakeGrammar(db); e != nil {
	// 	err = e
	// } else if q, e := qdb.NewQueries(db); e != nil {
	// 	err = e
	// } else {
	return
}
