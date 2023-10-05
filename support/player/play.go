package player

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/game"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/decode"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/print"
	"git.sr.ht/~ionous/tapestry/support/play"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/web/markup"
	"github.com/ionous/errutil"
)

// open db, select tests, de-gob and run them each in turn.
// print the results, only error on critical errors
func PlayGame(inFile, testString, domain string) (ret int, err error) {
	opts := qna.NewOptions()
	return PlayWithOptions(inFile, testString, domain, opts)
}

func PlayWithOptions(inFile, testString, domain string, opts qna.Options) (ret int, err error) {
	const prompt = "> "
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
			w := print.NewLineSentences(markup.ToText(os.Stdout))
			d := decode.NewDecoder(tapestry.AllSignatures)
			run := qna.NewRuntimeOptions(w, query, d, qna.Notifier{}, opts)
			if e := run.ActivateDomain(domain); e != nil {
				err = e
			} else {
				survey := play.MakeDefaultSurveyor(run)
				play := play.NewPlaytime(run, survey, grammar)
				if _, e := play.Call("start game", affine.None, nil, []g.Value{survey.GetFocalObject()}); e != nil {
					err = e
				} else if len(testString) > 0 {
					for _, cmd := range strings.Split(testString, ";") {
						fmt.Println(prompt, cmd)
						if step(play, cmd) {
							// some sort of thing about ending early if so?
							break
						}
					}
				} else {
					reader := bufio.NewReader(os.Stdin)
				Out:
					for {
						if len(prompt) > 0 {
							fmt.Print(prompt)
						}
						if in, _ := reader.ReadString('\n'); len(in) <= 1 {
							break
						} else {
							words := in[:len(in)-1] // strip the newline.
							for _, cmd := range strings.Split(words, ";") {
								if step(play, cmd) {
									break Out
								}
							}
						}
					}
				}
			}
		}
	}
	return
}

func step(p *play.Playtime, s string) (done bool) {
	var sig game.Signal
	if res, e := p.Step(s); errors.As(e, &sig) {
		switch sig {
		case game.SignalQuit:
			done = true
		default:
			log.Println("unhandled signal:", e)
		}
	} else if e != nil {
		log.Println("error:", e)
	} else if res != nil {
		fmt.Println()
	}
	return
}
