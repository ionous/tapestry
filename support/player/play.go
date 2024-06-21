package player

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/game"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/decode"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/print"
	"git.sr.ht/~ionous/tapestry/support/play"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/web/markup"
)

// open db, select tests, de-gob and run them each in turn.
// print the results, only error on critical errors
func PlayGame(mdlFile, testString, scene string) (err error) {
	opts := qna.NewOptions()
	return PlayWithOptions(mdlFile, testString, scene, opts)
}

func PlayWithOptions(mdlFile, testString, scene string, opts qna.Options) (err error) {
	const prompt = "> "
	if db, e := tables.CreateRunTime(mdlFile); e != nil {
		err = e
	} else {
		defer db.Close()
		if grammar, e := play.MakeGrammar(db); e != nil {
			err = e
		} else if q, e := qdb.NewQueries(db); e != nil {
			err = e
		} else {
			d := decode.NewDecoder(tapestry.AllSignatures)
			run := qna.NewRuntimeOptions(q, d, opts)
			w := print.NewLineSentences(markup.ToText(os.Stdout))
			run.SetWriter(w)
			if e := run.ActivateDomain(scene); e != nil {
				err = e
			} else {
				survey := play.MakeDefaultSurveyor(run)
				play := play.NewPlaytime(run, survey, grammar)
				if _, e := play.Call("start game", affine.None, nil, []rt.Value{survey.GetFocalObject()}); e != nil {
					err = e
				} else if len(testString) > 0 {
					for _, cmd := range strings.Split(testString, ";") {
						fmt.Println(prompt, cmd)
						if step(play, scene, cmd) {
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
								if step(play, scene, cmd) {
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

type SaveTime interface {
	LoadGame(string) (string, error)
	SaveGame(string) (string, error)
}

func step(p *play.Playtime, scene string, s string) (done bool) {
	var sig game.Signal
	if res, e := p.Step(s); errors.As(e, &sig) {
		switch sig {
		case game.SignalLoad:
			if saver, ok := p.Runtime.(SaveTime); !ok {
				log.Println("this runtime doesn't support save/load")
			} else if res, e := saver.LoadGame(scene); e != nil {
				log.Printf("couldn't load game because %v\n", e)
			} else {
				log.Printf("loaded %s from %s\n", scene, res)
			}

		case game.SignalSave:
			if saver, ok := p.Runtime.(SaveTime); !ok {
				log.Print("this runtime doesn't support save/load")
			} else if res, e := saver.SaveGame(scene); e != nil {
				log.Printf("couldn't save game because %v\n", e)
			} else {
				log.Printf("saved %s to %s\n", scene, res)
			}
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
