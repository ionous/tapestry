package cmdserve

import (
	"bytes"
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/game"
	"git.sr.ht/~ionous/tapestry/parser"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/decode"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/support/play"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

type State struct {
	Name          string
	HandleInput   func(w http.ResponseWriter, in string) (State, error)
	HandleCommand func(w http.ResponseWriter, cmd map[string]any) (State, error)
}

func unexpectedCommand(cmd map[string]any) error {
	return errutil.New("unexpected or unknown command", cmd)
}

func newServerContext(inFile string, opts qna.Options) (ret serverContext, err error) {
	if inFile, e := filepath.Abs(inFile); e != nil {
		err = e
	} else if db, e := sql.Open(tables.DefaultDriver, inFile); e != nil {
		err = errutil.New("couldn't create output file", inFile, e)
	} else if e := tables.CreateRun(db); e != nil {
		err = e
	} else if query, e := qdb.NewQueries(db, true); e != nil {
		err = e
	} else if grammar, e := play.MakeGrammar(db); e != nil {
		err = e
	} else {
		ret = serverContext{
			inFile:  inFile,
			db:      db,
			query:   query,
			grammar: grammar,
			opts:    opts,
			decoder: decode.NewDecoder(tapestry.AllSignatures),
		}
	}
	return
}

// do we really need all these things?
type serverContext struct {
	inFile  string
	db      *sql.DB
	query   *qdb.Query
	grammar parser.Scanner
	opts    qna.Options
	decoder *decode.Decoder
}

func (s *serverContext) Close() {
	log.Println("closing game db")
	s.db.Close()
}

func restart(w http.ResponseWriter, ctx serverContext, scene string) (ret State, err error) {
	log.Println("*** shuttle restart requested ***", scene)
	var buf bytes.Buffer // fix: maybe it'd be better if Step() handled the text
	run := qna.NewRuntimeOptions(&buf, ctx.query, ctx.decoder, ctx.opts)
	if e := run.ActivateDomain(scene); e != nil {
		err = e
	} else {
		survey := play.MakeDefaultSurveyor(run)
		play := play.NewPlaytime(run, survey, ctx.grammar)
		if _, e := play.Call("start game", affine.None, nil, []g.Value{survey.GetFocalObject()}); e != nil {
			err = e
		} else {
			// write the results of start game.
			writePlainText(w, buf.String())
			// return the playing state:
			ret = State{
				Name: "Playing " + scene,
				HandleInput: func(w http.ResponseWriter, words string) (ret State, err error) {
					log.Println(">", words)
					buf.Reset()
					for _, word := range strings.Split(words, ";") {
						// returns true if quit
						if step(play, word) {
							ret = stateGameOver()
						}
					}
					writePlainText(w, buf.String()) // write even on error
					return
				},
			}
		}
	}
	return
}

// a state which errors on all input and commands
func stateGameOver() State {
	log.Println("*** game over ***")
	return State{
		Name: "stateGameOver",
		HandleInput: func(w http.ResponseWriter, in string) (_ State, err error) {
			err = errutil.New("game over")
			return
		},
		HandleCommand: func(w http.ResponseWriter, cmd map[string]any) (_ State, err error) {
			err = errutil.New("game over")
			return
		},
	}
}

func writePlainText(w http.ResponseWriter, str string) {
	//log.Println("RETURNING", str)
	_, _ = io.WriteString(w, str)
	w.Header().Set("Content-Type", "plain/text")
}

func step(p *play.Playtime, s string) (done bool) {
	var sig game.Signal
	if _, e := p.Step(s); errors.As(e, &sig) {
		switch sig {
		case game.SignalQuit:
			done = true
		default:
			log.Println("unhandled signal:", e)
		}
	} else if e != nil {
		log.Println("error:", e)
	}
	return
}
