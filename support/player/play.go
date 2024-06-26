package player

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/game"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/parser"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/qna/raw"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/print"
	"git.sr.ht/~ionous/tapestry/support/files"
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
	d := query.NewDecoder(tapestry.AllSignatures)
	if ctx, e := createContext(mdlFile, d); e != nil {
		err = e
	} else if e := goPlay(ctx, scene, opts, testString); e != nil {
		err = e
	} else {
		ctx.q.Close()
	}
	return
}

func createContext(mdlFile string, d *query.QueryDecoder) (ret context, err error) {
	if path.Ext(mdlFile) == ".json" {
		ret, err = createRawContext(mdlFile, d)
	} else {
		ret, err = createSqlContext(mdlFile, d)
	}
	return
}

type context struct {
	grammar parser.Scanner
	q       query.Query
}

func createRawContext(mdlFile string, dec *query.QueryDecoder) (ret context, err error) {
	var data raw.Data
	if e := files.LoadJson(mdlFile, &data); e != nil {
		err = e
	} else if gram, e := readRawGrammar((*decode.Decoder)(dec), data.Grammar); e != nil {
		err = e
	} else {
		q := raw.MakeQuery(&data, dec)
		ret = context{gram, q}
	}
	return
}

func createSqlContext(mdlFile string, dec *query.QueryDecoder) (ret context, err error) {
	if db, e := tables.CreateRunTime(mdlFile); e != nil {
		err = e
	} else {
		if grammar, e := ReadGrammar(db, (*decode.Decoder)(dec)); e != nil {
			err = e
		} else if q, e := qdb.NewQueries(db, dec); e != nil {
			err = e
		} else {
			ret = context{grammar, q}
		}
		if err != nil { // otherwise query will take care of it
			defer db.Close()
		}
	}
	return
}

func goPlay(ctx context, scene string, opts qna.Options, testString string) (err error) {
	const prompt = "> "
	run := qna.NewRuntimeOptions(ctx.q, opts)
	w := print.NewLineSentences(markup.ToText(os.Stdout))
	run.SetWriter(w)
	if e := run.ActivateDomain(scene); e != nil {
		err = e
	} else {
		survey := play.MakeDefaultSurveyor(run)
		play := play.NewPlaytime(run, survey, ctx.grammar)
		if _, e := play.Call("start game", affine.None, nil, []rt.Value{survey.GetFocalObject()}); e != nil {
			err = e
		} else if len(testString) > 0 {
			for _, cmd := range strings.Split(testString, ";") {
				fmt.Println(prompt, cmd)
				if step(play, scene, cmd) {

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
					words := in[:len(in)-1]
					for _, cmd := range strings.Split(words, ";") {
						if step(play, scene, cmd) {
							break Out
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
