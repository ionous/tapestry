package play

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/decode"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/print"
	play "git.sr.ht/~ionous/tapestry/support/play/internal"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/web/js"
	"git.sr.ht/~ionous/tapestry/web/markup"
	"github.com/ionous/errutil"
)

// open db, select tests, de-gob and run them each in turn.
// print the results, only error on critical errors
func PlayGame(inFile, testString, domain string, jsonMode bool) (ret int, err error) {
	var prompt string
	if !jsonMode {
		prompt = "> "
	}
	if inFile, e := filepath.Abs(inFile); e != nil {
		err = e
	} else if db, e := sql.Open(tables.DefaultDriver, inFile); e != nil {
		err = errutil.New("couldn't create output file", inFile, e)
	} else {
		defer func() {
			_ = db.Close() // log?
		}()
		// fix: some sort of reset flag; but also: how to rejoin properly?
		if e := tables.CreateRun(db); e != nil {
			err = e
		} else if query, e := qdb.NewQueries(db, true); e != nil {
			err = e
		} else if grammar, e := play.MakeGrammar(db); e != nil {
			err = e
		} else {
			opt := qna.NewOptions()
			var w io.Writer
			var bufferedText bytes.Buffer
			if !jsonMode {
				w = print.NewLineSentences(markup.ToText(os.Stdout))
			} else {
				_ = opt.SetOption(meta.JsonMode, g.BoolOf(jsonMode))
				w = print.NewLineSentences(&bufferedText)
			}
			d := decode.NewDecoder(tapestry.AllSignatures)
			rx := qna.NewRuntimeOptions(w, query, d, opt)
			if _, e := rx.ActivateDomain(domain); e != nil {
				err = e
			} else {
				run := play.NewPlaytime(rx)

				parser := play.NewParser(run, grammar)
				//
				if len(testString) > 0 {
					for _, cmd := range strings.Split(testString, ";") {
						fmt.Println(prompt, cmd)
						step(parser, cmd, !jsonMode)
					}
				} else {
					reader := bufio.NewReader(os.Stdin)
					for {
						if len(prompt) > 0 {
							fmt.Printf(prompt)
						}
						if in, _ := reader.ReadString('\n'); len(in) <= 1 {
							break
						} else {
							words := in[:len(in)-1] // strip the newline.
							step(parser, words, !jsonMode)

							if jsonMode {
								// take buffered text and write it out
								// fix? possibly splitting newlines into array entries?
								// fix? right now serve turns this into play commands
								// but as/if we add more... we should do it in here
								// including probably hijacking the log output
								var out js.Builder
								out.Brace(js.Obj, func(inner *js.Builder) {
									inner.Q("out").R(js.Colon).Q(bufferedText.String())
								})
								_, _ = io.WriteString(os.Stdout, out.String())
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

func step(p *play.Parser, s string, pad bool) {
	if res, e := p.Step(s); e != nil {
		log.Println("error:", e)
	} else if res != nil && pad {
		fmt.Println()
	}
}