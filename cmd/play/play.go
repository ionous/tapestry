package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/tapestry"
	play "git.sr.ht/~ionous/tapestry/cmd/play/internal"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/web/js"
	"github.com/ionous/errutil"
)

// go run play.go -in /Users/ionous/Documents/Tapestry/build/play.db
func main() {
	var inFile, testString string
	var json bool
	flag.StringVar(&inFile, "in", "", "input file name (sqlite3)")
	flag.StringVar(&testString, "test", "", "input test string")
	flag.BoolVar(&json, "json", false, "expect input/output in json (default is plain text)")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.Parse()
	debug.LogLevel = debug.LoggingLevel{debug.LoggingLevel_Warning}
	var prompt string
	if !json {
		prompt = "> "
	}
	if cnt, e := playGame(inFile, testString, prompt, json); e != nil {
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
func playGame(inFile, testString, prompt string, jsonMode bool) (ret int, err error) {
	if inFile, e := filepath.Abs(inFile); e != nil {
		err = e
	} else if db, e := sql.Open(tables.DefaultDriver, inFile); e != nil {
		err = errutil.New("couldn't create output file", inFile, e)
	} else {
		defer db.Close()
		// fix: some sort of reset flag; but also: how to rejoin properly?
		if qdb, e := qdb.NewQueries(db, true); e != nil {
			err = e
		} else if grammar, e := play.MakeGrammar(db); e != nil {
			err = e
		} else {
			opt := qna.NewOptions()
			var w io.Writer
			var bufferedText bytes.Buffer
			if !jsonMode {
				w = os.Stdout
			} else {
				opt.SetOption(meta.JsonMode, g.BoolOf(jsonMode))
				w = &bufferedText
			}

			rx := qna.NewRuntimeOptions(w, qdb, opt, tapestry.AllSignatures)
			run := play.NewPlaytime(rx, "player", "kitchen")
			if _, e := run.ActivateDomain("entire_game"); e != nil {
				err = e
			} else {
				parser := play.NewParser(run, grammar)
				//
				if len(testString) > 0 {
					for _, cmd := range strings.Split(testString, ";") {
						fmt.Println(prompt, cmd)
						step(parser, cmd)
					}
				} else {
					reader := bufio.NewReader(os.Stdin)
					for {
						fmt.Printf(prompt)
						if in, _ := reader.ReadString('\n'); len(in) <= 1 {
							break
						} else {
							words := in[:len(in)-1] // strip the newline.
							step(parser, words)

							if jsonMode {
								// take buffered text and write it out
								// fix? possibly splitting newlines into array entries?
								var out js.Builder
								out.Brace(js.Obj, func(inner *js.Builder) {
									inner.Q("out").R(js.Colon).Q(bufferedText.String())
								})
								io.WriteString(os.Stdout, out.String())
								bufferedText.Reset()
							}
						}
					}
				}
			}
		}
	}
	return
}

func step(p *play.Parser, s string) {
	if res, e := p.Step(s); e != nil {
		fmt.Println("error:", e)
	} else if res != nil {
		fmt.Println()
	}
}
