// Package player provides a console like game.
package player

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/parser"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/qna/raw"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/print"
	"git.sr.ht/~ionous/tapestry/support/play"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/web/markup"
)

// open db, select tests, de-gob and run them each in turn.
// print the results, only error on critical errors
func PlayGame(mdlFile, testString, scene string) error {
	opts := qna.NewOptions()
	return PlayWithOptions(mdlFile, testString, scene, opts)
}

func PlayWithOptions(mdlFile, testString, scene string, opts qna.Options) (err error) {
	if ctx, e := createContext(mdlFile); e != nil {
		err = e
	} else if e := goPlay(ctx, scene, opts, testString); e != nil {
		err = e
	} else {
		ctx.q.Close()
	}
	return
}

func createContext(mdlFile string) (ret context, err error) {
	if path.Ext(mdlFile) == ".gob" {
		ret, err = createRawContext(mdlFile)
	} else if path.Ext(mdlFile) == ".db" {
		ret, err = createSqlContext(mdlFile)
	} else {
		err = fmt.Errorf("expected a db or gob file")
	}
	return
}

type context struct {
	q       query.Query
	grammar []parser.Scanner
}

func createRawContext(mdlFile string) (ret context, err error) {
	var data raw.Data
	if e := LoadGob(mdlFile, &data); e != nil {
		err = e
	} else {
		q := raw.MakeQuery(&data)
		scan := make([]parser.Scanner, len(data.Grammar))
		for i, d := range data.Grammar {
			scan[i] = d.MakeScanners()
		}
		ret = context{q, scan}
	}
	return
}

// deserialize from the passed path
func LoadGob(inPath string, pd *raw.Data) (err error) {
	tapestry.Register(gob.Register)
	if fp, e := os.Open(inPath); e != nil {
		err = e
	} else {
		dec := gob.NewDecoder(fp)
		err = dec.Decode(pd)
		defer fp.Close()
	}
	return
}

func createSqlContext(mdlFile string) (ret context, err error) {
	dec := query.NewDecoder(tapestry.AllSignatures)
	if db, e := tables.CreateRunTime(mdlFile); e != nil {
		err = e
	} else {
		if grammar, e := qdb.ReadGrammar(db, dec); e != nil {
			err = e
		} else if q, e := qdb.NewQueries(db, dec); e != nil {
			err = e
		} else {
			ret = context{q, grammar}
		}
		if err != nil { // otherwise query will take care of it
			defer db.Close()
		}
	}
	return
}

func goPlay(ctx context, scene string, opts qna.Options, testString string) (err error) {
	run := qna.NewRuntimeOptions(ctx.q, opts)
	run.SetWriter(print.NewLineSentences(markup.ToText(os.Stdout)))
	survey := play.MakeDefaultSurveyor(run)
	pt := play.NewPlaytime(run, survey, ctx.grammar)
	if e := pt.ActivateDomain(scene); e != nil {
		err = e
	} else if e := pt.RunPlayerAction(play.StartGame); e != nil {
		err = e
	} else if len(testString) > 0 {
		for _, cmd := range strings.Split(testString, ";") {
			//
			if e := pt.RunPlayerAction(play.RequestingPlayerInput); e != nil {
				err = e
				break
			} else {
				fmt.Println(cmd)
				if handleTurn(pt, nil, scene, cmd) {
					break // done
				}
			}
		}
	} else {
		reader := bufio.NewReader(os.Stdin)
		persist := Persistence{run, ctx.q}
	Out:
		// loop forever
		for {
			if e := pt.RunPlayerAction(play.RequestingPlayerInput); e != nil {
				err = e
				break
			} else if in, _ := reader.ReadString('\n'); len(in) <= 1 {
				break
			} else {
				words := in[:len(in)-1]
				for _, cmd := range strings.Split(words, ";") {
					if handleTurn(pt, &persist, scene, cmd) {
						break Out
					}
				}
			}
		}
	}
	return
}

func handleTurn(pt *play.Playtime, ps *Persistence, story string, s string) (done bool) {
	var sig rt.Signal
	if okay, e := pt.HandleTurn(s); errors.As(e, &sig) {
		switch sig {
		case rt.SignalLoad:
			if ps == nil {
				log.Println("this runtime doesn't support save/load")
			} else if str, e := ps.LoadGame(story); e != nil {
				log.Printf("couldn't load game because %v\n", e)
			} else {
				log.Printf("loaded %s from %s\n", story, str)
			}

		case rt.SignalSave:
			if ps == nil {
				log.Print("this runtime doesn't support save/load")
			} else if str, e := ps.SaveGame(story); e != nil {
				log.Printf("couldn't save game because %v\n", e)
			} else {
				log.Printf("saved %s to %s\n", story, str)
			}
		case rt.SignalQuit:
			done = true
		default:
			log.Println("unhandled signal:", e)
		}
	} else if e != nil {
		log.Println("error:", e)
	} else if okay {
		fmt.Println()
	}
	return
}
