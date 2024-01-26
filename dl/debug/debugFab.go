package debug

import (
	"errors"
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/game"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// a callback so fabricate can trigger a step of the parser;
// fabricate will error if this is nil/unset.
// ( see also: cmdcheck )
var Stepper func(words string) error

func (op *Fabricate) Execute(run rt.Runtime) (err error) {
	if e := op.fabricate(run); e != nil {
		err = assign.CmdError(op, e)
	}
	return
}

func (op *Fabricate) fabricate(run rt.Runtime) (err error) {
	if Stepper == nil {
		err = errutil.New("no parser set for the fabricator")
	} else if words, e := safe.GetText(run, op.Text); e != nil {
		err = e
	} else {
		words := words.String()
		if LogLevel <= C_LoggingLevel_Debug {
			log.Println("> ", words)
		}
		if len(words) > 0 {
			split := strings.Split(words, ";")
			for len(split) > 0 {
				var cmd string
				cmd, split = split[0], split[1:]
				if e := Stepper(cmd); e != nil {
					var sig game.Signal // if the game was quit, override the error if output remains
					if len(split) > 0 && errors.As(e, &sig) && sig == game.SignalQuit {
						e = errutil.New("game was quit, but input remains", strings.Join(split, ";"))
					}
					err = e
					break
				}
			}
		}
	}
	return
}
