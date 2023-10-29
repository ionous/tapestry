package shuttle

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/frame"
	"git.sr.ht/~ionous/tapestry/dl/game"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/decode"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/support/play"
)

// do we really need all these things?
type Shuttle struct {
	inFile  string
	db      *sql.DB
	query   *qdb.Query
	opts    qna.Options
	decoder *decode.Decoder
	play    *play.Playtime
	out     Collector
}

func (c *Shuttle) Close() {
	log.Println("closing game db")
	c.db.Close()
}

func (c *Shuttle) Restart(scene string) (ret *play.Playtime, err error) {
	// FIX: tear down.
	if play := c.play; play != nil {
		play.ActivateDomain("")
	}
	//
	if play, e := c.EnsurePlay(); e != nil {
		err = e
	} else if e := play.ActivateDomain(scene); e != nil {
		err = e
	} else if _, e := play.Call("start game", affine.None, nil, []g.Value{play.Survey().GetFocalObject()}); e != nil {
		err = e
	} else {
		ret = play
	}
	return
}

// create the playtime if it doesnt exist
func (c *Shuttle) EnsurePlay() (ret *play.Playtime, err error) {
	if c.play != nil {
		ret = c.play
	} else if grammar, e := play.MakeGrammar(c.db); e != nil {
		err = e
	} else {
		note := qna.Notifier{
			StartedScene:    c.out.onStartScene,
			EndedScene:      c.out.onEndScene,
			ChangedState:    c.out.onChangeState,
			ChangedRelative: c.out.onChangeRel,
		}
		run := qna.NewRuntimeOptions(&c.out.buf, c.query, c.decoder, note, c.opts)
		survey := play.MakeDefaultSurveyor(run)
		play := play.NewPlaytime(run, survey, grammar)
		c.play = play
		ret = play
		var done bool
		// fix? used for fabricate; maybe use options instead so that we can have multiple instances?
		debug.Stepper = func(words string) (err error) {
			// FIX: errors for step are getting fmt.Println in playTime.go
			// so expect output can't test for errors ( and on error looks a bit borken )
			var sig game.Signal
			if _, e := play.Step(words); !done && errors.As(e, &sig) && sig == game.SignalQuit {
				done = true // eat the quit signal on first return; fix? maybe do this on client?
			} else if e != nil {
				err = e
			}
			return
		}
	}
	return
}

type Collector struct {
	events []frame.Event
	buf    strings.Builder
}

// returns and clears all events
func (out *Collector) GetEvents() (ret []frame.Event) {
	ret, out.events = out.flush(), nil
	return
}

func (out *Collector) onStartScene(domains []string) {
	out.flush()
	out.addEvent(&frame.SceneStarted{Domains: domains})
}
func (out *Collector) onEndScene(domains []string) {
	out.flush()
	out.addEvent(&frame.SceneEnded{Domains: domains})
}
func (out *Collector) onChangeState(noun, aspect, prev, trait string) {
	out.flush()
	out.addEvent(&frame.StateChanged{Noun: noun, Aspect: aspect, Prev: prev, Trait: trait})
}
func (out *Collector) onChangeRel(a, b, rel string) {
	out.flush()
	out.addEvent(&frame.PairChanged{A: a, B: b, Rel: rel})
}
func (out *Collector) addEvent(evt frame.Event) {
	out.events = append(out.events, evt)
}
func (out *Collector) flush() []frame.Event {
	if out.buf.Len() > 0 {
		str := out.buf.String()
		out.buf.Reset()
		out.addEvent(&frame.FrameOutput{Text: str})
	}
	return out.events
}
